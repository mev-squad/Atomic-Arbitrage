package main

import (
	"Nucleus/nucleus"
	"Nucleus/rpcClient"
	"Nucleus/stateRead"
	"fmt"
	"time"
)

func main() {
	//Initialize the rpcClient
	rpcClient.Initialize("./rpcClient/config.json")

	//sync pairs
	fmt.Println("loading pairs...")
	start := time.Now()
	quickToSushi, err := nucleus.LoadQuickToSushi()
	if err != nil {
		fmt.Println(err)
		return
	}

	//log time elapsed to sync pairs
	elapsed := time.Since(start)
	fmt.Printf("loaded %d pairs in %v\n", len(quickToSushi), elapsed)

	//search for arb opportunities
	for {
		currentBlock, blockNumber := stateRead.DownloadBlock("latest")
		fmt.Println(blockNumber)
		DetectedOpportunities := nucleus.SearchBlock(currentBlock)
		fmt.Println(DetectedOpportunities)
	}
}
