package main

import (
	"fmt"
	"testing"
	"time"
)

var DetectedOpportunities []ArbitrageOpportunity

func TestEndToEnd(t *testing.T) {
	// read config
	config, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	start := time.Now()
	currentBlock, _ := DownloadBlock("0x156A401", config.HttpURL)
	DetectedOpportunities = SearchBlock(currentBlock)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

func BenchmarkEndtoEnd(b *testing.B) {
	// read config
	config, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	currentBlock, _ := DownloadBlock("0x156A401", config.HttpURL)
	DetectedOpportunities = SearchBlock(currentBlock)
}
