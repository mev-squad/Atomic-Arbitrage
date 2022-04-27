package nucleus

import (
	"Nucleus/rpcClient"
	"Nucleus/stateRead"
	"fmt"
	"testing"
	"time"
)

var DetectedOpportunities []ArbitrageOpportunity

func TestEndToEnd(t *testing.T) {
	// initialize rpcClient
	rpcClient.Initialize("../rpcClient/config.json")

	start := time.Now()
	currentBlock, _ := stateRead.DownloadBlock("latest")
	DetectedOpportunities = SearchBlock(currentBlock)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

func BenchmarkEndtoEnd(b *testing.B) {
	// initialize rpcClient
	rpcClient.Initialize("../rpcClient/config.json")

	currentBlock, _ := stateRead.DownloadBlock("latest")
	DetectedOpportunities = SearchBlock(currentBlock)
}
