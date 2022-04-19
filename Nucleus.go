package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
)

var config Config
var QuickToSushi map[string]string

func main() {
	var err error

	// read config
	config, err = readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// This is the so-called "Nucleus" of Atomic Arbitrage
	// Arbitrage broken down to it's core elements
	fmt.Println("loading pairs...")
	start := time.Now()
	QuickToSushi, err = loadQuickToSushi()
	if err != nil {
		fmt.Println(err)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("loaded %d pairs in %v\n", len(QuickToSushi), elapsed)
	for {
		currentBlock, blockNumber := DownloadBlock("latest", config.HttpURL)
		fmt.Println(blockNumber)
		DetectedOpportunities := SearchBlock(currentBlock)
		fmt.Println(DetectedOpportunities)
	}
}

const ABIS_DIR = "abis"
const QUICK_FACTORY_ADDRESS = "0x5757371414417b8C6CAad45bAeF941aBc7d3Ab32"
const SUSHI_FACTORY_ADDRESS = "0xc35DADB65012eC5796536bD9864eD8773aBc74C4"

type Config struct {
	HttpURL      string
	WebSocketURL string
}

type ArbitrageOpportunity struct {
	AmountIn    *uint256.Int // [3]uint64
	poolAddress string
	AtoB        bool
}

// read config.json
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

// read ABIs from abis folder
func readABIs() (map[string]abi.ABI, error) {
	files, err := ioutil.ReadDir(ABIS_DIR)
	if err != nil {
		return nil, err
	}
	abis := make(map[string]abi.ABI)
	for _, fileInfo := range files {
		name := fileInfo.Name()
		file, err := ioutil.ReadFile(ABIS_DIR + "/" + name)
		if err != nil {
			return nil, err
		}
		_abi, err := abi.JSON(strings.NewReader(string(file)))
		if err != nil {
			return nil, err
		}
		abiName := name[:len(name)-5]
		abis[abiName] = _abi
	}
	return abis, nil
}

// call contract method
func callContract(client *ethclient.Client, to *common.Address, ABI abi.ABI, method string, args ...interface{}) ([]interface{}, error) {
	data, err := ABI.Pack(method, args...)
	if err != nil {
		return []interface{}{}, err
	}
	msg := ethereum.CallMsg{To: to, Data: data}
	data, err = client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return []interface{}{}, err
	}
	values, err := ABI.Unpack(method, data)
	if err != nil {
		return []interface{}{}, err
	}
	return values, nil
}

// compare addresses
func compareAddresses(a common.Address, b common.Address) int {
	return bytes.Compare(a.Bytes(), b.Bytes())
}

func loadQuickToSushi() (map[string]string, error) {
	quickFactoryAddress := common.HexToAddress(QUICK_FACTORY_ADDRESS)
	sushiFactoryAddress := common.HexToAddress(SUSHI_FACTORY_ADDRESS)

	// read ABIs
	abis, err := readABIs()
	if err != nil {
		return nil, err
	}

	// connect to node
	client, err := ethclient.Dial(config.WebSocketURL)
	if err != nil {
		return nil, err
	}

	// get block number
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	// get PairCreated logs
	// batch to avoid errors
	event := abis["Factory"].Events["PairCreated"]
	query := ethereum.FilterQuery{
		Addresses: []common.Address{quickFactoryAddress, sushiFactoryAddress},
		Topics:    [][]common.Hash{[]common.Hash{event.ID}},
	}
	batchSize := uint64(1000000)
	var wg sync.WaitGroup
	var quickPairsMu sync.Mutex
	quickPairs := make(map[common.Address]map[common.Address]common.Address)
	var sushiPairsMu sync.Mutex
	sushiPairs := make(map[common.Address]map[common.Address]common.Address)
	for i := uint64(0); i <= blockNumber; i += batchSize {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			query := query
			query.FromBlock = new(big.Int).SetUint64(i)
			query.ToBlock = new(big.Int).SetUint64(i + batchSize - 1)

			logs, err := client.FilterLogs(context.Background(), query)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, vLog := range logs {
				// get factory address
				factoryAddress := vLog.Address

				// get pair address
				data, err := event.Inputs.Unpack(vLog.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				pairAddress := data[0].(common.Address)

				// get token addresses
				token0Address := common.HexToAddress(vLog.Topics[1].Hex())
				token1Address := common.HexToAddress(vLog.Topics[2].Hex())

				// store pair
				if compareAddresses(factoryAddress, quickFactoryAddress) == 0 {
					quickPairsMu.Lock()
					if _, ok := quickPairs[token0Address]; !ok {
						quickPairs[token0Address] = make(map[common.Address]common.Address)
					}
					quickPairs[token0Address][token1Address] = pairAddress
					quickPairsMu.Unlock()
				} else if compareAddresses(factoryAddress, sushiFactoryAddress) == 0 {
					sushiPairsMu.Lock()
					if _, ok := sushiPairs[token0Address]; !ok {
						sushiPairs[token0Address] = make(map[common.Address]common.Address)
					}
					sushiPairs[token0Address][token1Address] = pairAddress
					sushiPairsMu.Unlock()
				} else {
					fmt.Println("?")
					continue
				}
			}
		}()
	}
	wg.Wait()

	quickToSushi := make(map[string]string)
	for token0Address := range quickPairs {
		for token1Address := range quickPairs[token0Address] {
			quick := quickPairs[token0Address][token1Address]
			sushi, ok := sushiPairs[token0Address][token1Address]
			if !ok {
				sushi = quick
			}
			quickToSushi[strings.ToLower(quick.String())] = sushi.String()
		}
	}
	return quickToSushi, nil
}

func SearchBlock(block []interface{}) []ArbitrageOpportunity {
	var ArbitrageOpportunities []ArbitrageOpportunity
	// Go routine for each tx is too slow because 90% of transactions
	// don't require more than one check, Pareto principle states 20%
	// of the transactions will consume 80% of runtime, so we only
	// need to parallelize that 20% of transactions
	for index := 0; index < len(block); index++ {
		tx := block[index].(map[string]interface{})
		to := tx["to"]
		// handle contract creations
		if to == nil {
			continue
		}
		// Checks for calls to the Quickswap Router
		if to.(string) == "0xa5e0829caced8ffdd4de3c43696c57f7d7a678ff" {
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
			pool0Reseve0, pool0Reseve1 := GetReserves(poolAddress, config.HttpURL)
			pool1Reseve0, pool1Reseve1 := GetReserves(sushiAddress, config.HttpURL)
			aToB, result := ComputeProfitMaximizingTrade(pool0Reseve0, pool0Reseve1, pool1Reseve0, pool1Reseve1)
			ArbitrageOpportunities = append(ArbitrageOpportunities, ArbitrageOpportunity{result, poolAddress, aToB})
		}
	}
	return ArbitrageOpportunities
}
