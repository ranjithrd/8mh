#!/bin/bash
# Quick deployment script for PaymentAnchor contract to Sepolia

set -e

echo "🚀 PaymentAnchor Contract Deployment Script"
echo "=========================================="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found"
    echo "Please create .env with SEPOLIA_RPC_URL and PRIVATE_KEY"
    exit 1
fi

# Load environment variables
source .env

# Check required variables
if [ -z "$SEPOLIA_RPC_URL" ]; then
    echo "❌ Error: SEPOLIA_RPC_URL not set in .env"
    exit 1
fi

if [ -z "$PRIVATE_KEY" ]; then
    echo "❌ Error: PRIVATE_KEY not set in .env"
    exit 1
fi

echo "✓ Environment variables loaded"
echo ""

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing Node.js dependencies..."
    npm install --save-dev hardhat @nomicfoundation/hardhat-toolbox dotenv
    echo "✓ Dependencies installed"
    echo ""
fi

# Check if hardhat.config.js exists
if [ ! -f "hardhat.config.js" ]; then
    echo "⚙️  Creating hardhat.config.js..."
    cat > hardhat.config.js << 'EOF'
import "@nomicfoundation/hardhat-toolbox";
import dotenv from "dotenv";

dotenv.config();

export default {
  solidity: "0.8.0",
  networks: {
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL,
      accounts: [`0x${process.env.PRIVATE_KEY}`]
    }
  }
};
EOF
    echo "✓ hardhat.config.js created"
    echo ""
fi

# Create scripts directory if it doesn't exist
mkdir -p scripts

# Create deployment script
echo "📝 Creating deployment script..."
cat > scripts/deploy.js << 'EOF'
import hre from "hardhat";

async function main() {
  console.log("Deploying PaymentAnchor contract to Sepolia...");
  console.log("");
  
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);
  console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString());
  console.log("");
  
  const PaymentAnchor = await hre.ethers.getContractFactory("PaymentAnchor");
  const paymentAnchor = await PaymentAnchor.deploy();
  
  await paymentAnchor.waitForDeployment();
  
  const address = await paymentAnchor.getAddress();
  console.log("✓ PaymentAnchor deployed to:", address);
  console.log("");
  console.log("📋 Next Steps:");
  console.log("1. Add this to your .env file:");
  console.log(`   CONTRACT_ADDRESS=${address}`);
  console.log("");
  console.log("2. Verify contract on Etherscan (optional):");
  console.log(`   npx hardhat verify --network sepolia ${address}`);
  console.log("");
  console.log("3. View on Etherscan:");
  console.log(`   https://sepolia.etherscan.io/address/${address}`);
  console.log("");
  console.log("🎉 Deployment complete!");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
EOF
echo "✓ Deployment script created"
echo ""

# Deploy
echo "🚀 Deploying to Sepolia..."
echo ""
npx hardhat run scripts/deploy.js --network sepolia

echo ""
echo "=========================================="
echo "✅ Deployment script completed!"
echo ""
echo "Don't forget to update your .env file with CONTRACT_ADDRESS"
