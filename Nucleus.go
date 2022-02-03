package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/holiman/uint256"
)

var QuickToSushi map[string]string

func main() {
	// This is the so-called "Nucleus" of Atomic Arbitrage
	// Arbitrage broken down to it's core elements
	QuickToSushi = readQuickTokenPairs()
	for {
		currentBlock, blockNumber := DownloadBlock("latest")
		fmt.Println(blockNumber)
		DetectedOpportunities := SearchBlock(currentBlock)
		fmt.Println(DetectedOpportunities)
	}
}

func readQuickTokenPairs() map[string]string {
	jsonFile, _ := os.Open("quickswapPairs.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var SushiToQuick map[string]string

	json.Unmarshal(byteValue, &SushiToQuick)

	return SushiToQuick

}

type ArbitrageOpportunity struct {
	AmountIn    *uint256.Int // [3]uint64
	poolAddress string
	AtoB        bool
}

func SearchBlock(block []interface{}) []ArbitrageOpportunity {
	var ArbitrageOpportunities []ArbitrageOpportunity
	// Go routine for each tx is too slow because 90% of transactions
	// don't require more than one check, Pareto principle states 20%
	// of the transactions will consume 80% of runtime, so we only
	// need to parallelize that 20% of transactions
	for index := 0; index < len(block); index++ {
		tx := block[index].(map[string]interface{})
		// Checks for calls to the Quickswap Router
		if tx["to"].(string) == "0xa5e0829caced8ffdd4de3c43696c57f7d7a678ff" {
			senderData := tx["input"]
			pathLength := ParsePathLength(senderData.(string))
			// Cap of 30 is set to stop arbitraging of arbitrage
			// Otherwise wastes time checking 30 pools guranteed to
			// explicitly not have an opportunity
			if pathLength > 0 && pathLength < 30 {

				if senderData.(string)[:10] == "0x38ed1739" {
					ArbitrageOpportunities = findSwapExactTokensForTokens(senderData.(string), pathLength)
				}
			}
		}
	}
	return ArbitrageOpportunities
}

func GetAffectedAddresses(inputData string, pathLength uint64) []string {
	affectedAddresses := make([]string, pathLength)
	for index := 0; index < int(pathLength); index++ {
		startNumber := 394 + (64 * index)
		endNumber := startNumber + 64
		affectedAddresses[index] = "0x" + (inputData[startNumber:endNumber])[24:]
	}
	return affectedAddresses
}

func ParsePathLength(tokenRoutePath string) uint64 {
	if len(tokenRoutePath) > 400 {
		pathLength := tokenRoutePath[330:394]
		hexPathLength := "0x" + RemoveLeadingZeros(pathLength)
		hexPathLength = strings.Replace(hexPathLength, "0x", "", -1)
		n, _ := strconv.ParseUint(hexPathLength, 16, 64)
		return n
	} else {
		return 0
	}
}

func findSwapExactTokensForTokens(inputData string, pathLength uint64) []ArbitrageOpportunity {
	changedAdddresses := GetAffectedAddresses(inputData, pathLength)
	var ArbitrageOpportunities []ArbitrageOpportunity
	// Next it's time to calculate some pair addresses
	for index := 0; index < len(changedAdddresses)-1; index++ {
		token0 := changedAdddresses[index]
		token1 := changedAdddresses[index+1]
		poolAddress := CalculatePairAddress(token0, token1)
		sushiAddress := QuickToSushi[poolAddress]
		if sushiAddress != "" {
			pool0Reseve0, pool0Reseve1 := GetReserves(poolAddress)
			pool1Reseve0, pool1Reseve1 := GetReserves(sushiAddress)
			aToB, result := ComputeProfitMaximizingTrade(pool0Reseve0, pool0Reseve1, pool1Reseve0, pool1Reseve1)
			ArbitrageOpportunities = append(ArbitrageOpportunities, ArbitrageOpportunity{result, poolAddress, aToB})
		}
	}
	return ArbitrageOpportunities
}
