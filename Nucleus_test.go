package main

import (
	"fmt"
	"testing"
	"time"
)

var DetectedOpportunities []ArbitrageOpportunity

func TestEndToEnd(t *testing.T) {
	start := time.Now()
	currentBlock, _ := DownloadBlock("0x156A401")
	DetectedOpportunities = SearchBlock(currentBlock)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

func BenchmarkEndtoEnd(b *testing.B) {
	currentBlock, _ := DownloadBlock("0x156A401")
	DetectedOpportunities = SearchBlock(currentBlock)
}
