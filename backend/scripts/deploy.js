const hre = require("hardhat");

async function main() {
  console.log("Deploying PaymentAnchor contract to Sepolia...");
  console.log("");
  
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);
  console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString());
  console.log("");
  
  const PaymentAnchor = await hre.ethers.getContractFactory("PaymentAnchor");
  const paymentAnchor = await PaymentAnchor.deploy();
  
  await paymentAnchor.deployed();
  
  console.log("✓ PaymentAnchor deployed to:", paymentAnchor.address);
  console.log("");
  console.log("📋 Add this to your .env file:");
  console.log(`CONTRACT_ADDRESS=${paymentAnchor.address}`);
  console.log("");
  console.log("🎉 Deployment complete!");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
