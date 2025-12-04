package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client          *ethclient.Client
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	contractABI     abi.ABI
)

// PaymentAnchorABI is the ABI for the PaymentAnchor contract
const PaymentAnchorABI = `[
	{
		"inputs": [
			{"internalType": "string", "name": "transactionId", "type": "string"},
			{"internalType": "bytes32", "name": "paymentHash", "type": "bytes32"}
		],
		"name": "anchorPayment",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "string", "name": "transactionId", "type": "string"}
		],
		"name": "getPaymentAnchor",
		"outputs": [
			{"internalType": "bytes32", "name": "", "type": "bytes32"},
			{"internalType": "uint256", "name": "", "type": "uint256"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "string", "name": "transactionId", "type": "string"},
			{"internalType": "bytes32", "name": "paymentHash", "type": "bytes32"}
		],
		"name": "verifyPayment",
		"outputs": [
			{"internalType": "bool", "name": "", "type": "bool"}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

// InitEthereum initializes the Ethereum client and loads configuration
func InitEthereum() error {
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		return fmt.Errorf("SEPOLIA_RPC_URL not set in environment")
	}

	contractAddrStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddrStr == "" {
		return fmt.Errorf("CONTRACT_ADDRESS not set in environment")
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return fmt.Errorf("PRIVATE_KEY not set in environment")
	}

	// Connect to Sepolia
	var err error
	client, err = ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Sepolia: %w", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}
	log.Printf("Connected to Sepolia (Chain ID: %s)", chainID.String())

	// Load contract address
	contractAddress = common.HexToAddress(contractAddrStr)
	log.Printf("Using contract at: %s", contractAddress.Hex())

	// Load private key
	privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	// Load contract ABI
	contractABI, err = abi.JSON(strings.NewReader(PaymentAnchorABI))
	if err != nil {
		return fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	log.Println("Ethereum client initialized successfully")
	return nil
}

// AnchorPaymentToSepolia anchors a payment hash to the Sepolia blockchain
// Returns the Ethereum transaction hash
func AnchorPaymentToSepolia(transactionID string, paymentHash [32]byte) (string, error) {
	if client == nil {
		return "", fmt.Errorf("ethereum client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the public address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Pack the transaction data
	data, err := contractABI.Pack("anchorPayment", transactionID, paymentHash)
	if err != nil {
		return "", fmt.Errorf("failed to pack transaction data: %w", err)
	}

	// Estimate gas limit
	gasLimit := uint64(150000) // Conservative estimate for anchorPayment

	// Create the transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddress,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	txHash := signedTx.Hash().Hex()
	log.Printf("Payment anchored to Sepolia: TX=%s, TxID=%s", txHash, transactionID)

	// Wait for transaction receipt (optional - can be async)
	go waitForReceipt(signedTx.Hash(), transactionID)

	return txHash, nil
}

// waitForReceipt waits for transaction confirmation (runs in background)
func waitForReceipt(txHash common.Hash, transactionID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	receipt, err := waitForTransactionReceipt(ctx, txHash)
	if err != nil {
		log.Printf("ERROR: Transaction %s failed: %v", txHash.Hex(), err)
		return
	}

	if receipt.Status == 1 {
		log.Printf("✓ Payment %s confirmed on Sepolia (Block: %d, Gas: %d)",
			transactionID, receipt.BlockNumber.Uint64(), receipt.GasUsed)
	} else {
		log.Printf("✗ Payment %s transaction reverted on Sepolia", transactionID)
	}
}

// waitForTransactionReceipt polls for transaction receipt
func waitForTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
			// Continue polling if receipt not found yet
		}
	}
}

// VerifyPaymentOnSepolia verifies a payment hash against the Sepolia blockchain
func VerifyPaymentOnSepolia(transactionID string, expectedHash [32]byte) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("ethereum client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pack the call data
	data, err := contractABI.Pack("verifyPayment", transactionID, expectedHash)
	if err != nil {
		return false, fmt.Errorf("failed to pack call data: %w", err)
	}

	// Call the contract
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return false, fmt.Errorf("failed to call contract: %w", err)
	}

	// Unpack the result
	var isValid bool
	err = contractABI.UnpackIntoInterface(&isValid, "verifyPayment", result)
	if err != nil {
		return false, fmt.Errorf("failed to unpack result: %w", err)
	}

	return isValid, nil
}

// GetPaymentAnchor retrieves the payment hash and timestamp from Sepolia
func GetPaymentAnchor(transactionID string) (hash string, timestamp int64, err error) {
	if client == nil {
		return "", 0, fmt.Errorf("ethereum client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pack the call data
	data, err := contractABI.Pack("getPaymentAnchor", transactionID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to pack call data: %w", err)
	}

	// Call the contract
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to call contract: %w", err)
	}

	// Unpack the result
	var results []interface{}
	err = contractABI.UnpackIntoInterface(&results, "getPaymentAnchor", result)
	if err != nil {
		return "", 0, fmt.Errorf("failed to unpack result: %w", err)
	}

	if len(results) != 2 {
		return "", 0, fmt.Errorf("unexpected result length: %d", len(results))
	}

	paymentHash := results[0].([32]byte)
	blockTimestamp := results[1].(*big.Int)

	return hex.EncodeToString(paymentHash[:]), blockTimestamp.Int64(), nil
}
