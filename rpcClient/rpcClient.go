package rpcClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ybbus/jsonrpc/v2"
)

var HTTPClient jsonrpc.RPCClient
var WSClient *ethclient.Client

var config Config

type Config struct {
	HttpURL      string
	WebSocketURL string
}

func Initialize(configPath string) {
	initializeConfig(configPath)
	initializeHTTPClient()
	initializeWSClient()
}

func initializeConfig(filePath string) {
	//Read in config.json
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error when reading config.json: %s\n", err)
	}
	//create a new config
	var newConfig Config
	json.Unmarshal(file, &newConfig)
	if err != nil {
		fmt.Printf("Error when reading unmarshaling config file: %s\n", err)
	}
	//initialize the config variable
	config = newConfig
}

func initializeHTTPClient() {
	newHTTPClient := jsonrpc.NewClient(config.HttpURL)
	HTTPClient = newHTTPClient
}

func initializeWSClient() {
	wsClient, err := ethclient.Dial(config.WebSocketURL)
	if err != nil {
		fmt.Printf("Error when initializing websocket client: %s\n", err)
	}
	WSClient = wsClient
}
