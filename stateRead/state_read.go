package stateRead

import (
	"fmt"

	"github.com/holiman/uint256"
	jsonrpc "github.com/ybbus/jsonrpc/v2"
)

func DownloadBlock(blockNumber string, nodeURL string) ([]interface{}, string) {
	rpcClient := jsonrpc.NewClient(nodeURL)

	result, err := rpcClient.Call("eth_getBlockByNumber", blockNumber, true)

	if err != nil {
		fmt.Println(err)
	}

	returnedBlock := result.Result.(map[string]interface{})["number"].(string)

	txArray := result.Result.(map[string]interface{})["transactions"].([]interface{})

	return txArray, returnedBlock
}

func RemoveLeadingZeros(inputString string) string {
	for index := 0; index < len(inputString); index++ {
		// Just return the string as soon as first non zero value is detected
		if string(inputString[index]) != "0" {
			return inputString[index:]
		}
	}
	return "" // Value is only zeros
}

func deriveReservesFromSlot(slot string) (uint256.Int, uint256.Int) {

	// Yes these names are right, for some reason they are stored in reverse order
	reserve1, _ := uint256.FromHex("0x" + RemoveLeadingZeros(slot[10:38]))

	reserve0, _ := uint256.FromHex("0x" + RemoveLeadingZeros(slot[38:66]))

	return *reserve0, *reserve1
}

func GetReserves(pair string, nodeURL string) (uint256.Int, uint256.Int) {
	rpcClient := jsonrpc.NewClient(nodeURL)

	response, _ := rpcClient.Call("eth_getStorageAt", pair, "0x8", "latest")

	reserve0, reserve1 := deriveReservesFromSlot(response.Result.(string))
	return reserve0, reserve1
}
