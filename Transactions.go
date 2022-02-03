package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/0xKitsune/go-web3"
	"github.com/0xKitsune/go-web3/abi"
	"github.com/0xKitsune/go-web3/contract"
	"github.com/0xKitsune/go-web3/jsonrpc"
	"github.com/0xKitsune/go-web3/wallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// /* Mostly stolen from Unibox */

var chainID = uint64(137)

var WalletAddress string = ""

var walletKey *wallet.Key = NewWalletKey("")

var web3Client *jsonrpc.Client = initializeWeb3Client("") // Alchemy URL here

var Arbitrager = "" // Initialize Later initializeSwapContract()

// //Returns *wallet.Key which is used to sign transactions
func NewWalletKey(privateKey string) *wallet.Key {
	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	key := wallet.NewKey(ecdsaPrivateKey)

	return key
}

// // initialize an instance of the Arbitrage Contract
func initializeSwapContract() *contract.Contract {
	//initialize a web3 address with the uniswap router hex address
	SwapContract := web3.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
	//read in the uniswapRouterV2 abi from file
	abiBytes, err := ioutil.ReadFile("abi/SwapContract.json")
	if err != nil {
		fmt.Println("Error when reading SwapContract ABI.")
	}
	//create a new web3 abi
	abi, err := abi.NewABI(string(abiBytes))
	if err != nil {
		fmt.Println("Error when creating SwapContract ABI", err)
		os.Exit(1)
	}
	contractInstance := contract.NewContract(SwapContract, abi, web3Client)
	// Wallet is considered to have deployed this contract
	contractInstance.SetFrom(web3.HexToAddress(WalletAddress))

	return contractInstance
}

func arbitragePools(path []string, aToB bool, amountIn uint256.Int) {

}

func SwapExactTokensForTokens(amountIn uint, amountOutMin uint, path []common.Address, to common.Address, deadline uint, swapContract *contract.Contract) web3.Hash {
	txn := swapContract.Txn("swapExactTokensForTokens", amountIn, amountOutMin, path, to, deadline)
	err := txn.SignAndSend(walletKey, chainID)
	if err != nil {
		panic(err)
	}
	return txn.Hash
}

// //initialize web3 http client
func initializeWeb3Client(nodeURL string) *jsonrpc.Client {
	client, err := jsonrpc.NewClient(nodeURL)
	if err != nil {
		fmt.Println("Failed to connect to node", err)
		os.Exit(1)
	}
	return client
}

func log(input string) {
	go fmt.Println(input)
}
