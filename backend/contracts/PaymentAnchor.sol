// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title PaymentAnchor
 * @dev Smart contract for anchoring payment hashes to Sepolia blockchain
 * Each payment in the cooperative banking system gets its hash recorded on-chain
 */
contract PaymentAnchor {
    // Event emitted when a payment is anchored
    event PaymentAnchored(
        string indexed transactionId,
        bytes32 paymentHash,
        uint256 timestamp,
        address indexed anchoredBy
    );

    // Mapping from transaction ID to payment hash
    mapping(string => bytes32) public paymentHashes;
    
    // Mapping from transaction ID to timestamp
    mapping(string => uint256) public anchorTimestamps;
    
    // Owner of the contract
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    /**
     * @dev Anchor a payment hash to the blockchain
     * @param transactionId The unique transaction ID from the banking system
     * @param paymentHash The SHA-256 hash of the payment data
     */
    function anchorPayment(string memory transactionId, bytes32 paymentHash) public {
        require(bytes(transactionId).length > 0, "Transaction ID cannot be empty");
        require(paymentHash != bytes32(0), "Payment hash cannot be zero");
        require(paymentHashes[transactionId] == bytes32(0), "Payment already anchored");

        paymentHashes[transactionId] = paymentHash;
        anchorTimestamps[transactionId] = block.timestamp;

        emit PaymentAnchored(transactionId, paymentHash, block.timestamp, msg.sender);
    }

    /**
     * @dev Verify a payment hash matches what's stored on-chain
     * @param transactionId The transaction ID to verify
     * @param paymentHash The hash to verify against
     * @return bool True if the hash matches
     */
    function verifyPayment(string memory transactionId, bytes32 paymentHash) public view returns (bool) {
        return paymentHashes[transactionId] == paymentHash;
    }

    /**
     * @dev Get the anchored hash for a transaction
     * @param transactionId The transaction ID to look up
     * @return The payment hash and timestamp
     */
    function getPaymentAnchor(string memory transactionId) public view returns (bytes32, uint256) {
        return (paymentHashes[transactionId], anchorTimestamps[transactionId]);
    }
}
