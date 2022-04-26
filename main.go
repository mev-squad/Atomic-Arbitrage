package main

import (
	"Nucleus/nucleus"
	"Nucleus/stateRead"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	var err error

	// read config
	config, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

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
		currentBlock, blockNumber := stateRead.DownloadBlock("latest", config.HttpURL)
		fmt.Println(blockNumber)
		DetectedOpportunities := nucleus.SearchBlock(currentBlock)
		fmt.Println(DetectedOpportunities)
	}
}

func readConfig() (Config, error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}
	var config Config
	json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

type Config struct {
	HttpURL      string
	WebSocketURL string
}
