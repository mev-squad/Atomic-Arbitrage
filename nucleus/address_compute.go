package nucleus

import (
	"Nucleus/stateRead"
	"encoding/hex"
	"strings"

	sha3 "github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

func SortTokens(token0 string, token1 string) (string, string) {
	// Returns the tokens sorted, but with less safety checks than the sol version
	token0Number, _ := uint256.FromHex("0x" + stateRead.RemoveLeadingZeros(token0[2:])) // Sanitation so it doesn't break on 0x00 addresses
	token1Number, _ := uint256.FromHex("0x" + stateRead.RemoveLeadingZeros(token0[2:]))
	if token0Number.Lt(token1Number) {
		return token0, token1
	} else {
		return token1, token0
	}
}

func CalculatePairAddress(tokenA string, tokenB string) string {
	token0, token1 := SortTokens(tokenA, tokenB)

	/* The next few lines are replicating this clusterfuck
	pair = address(uint(keccak256(abi.encodePacked(
	      hex'ff',
	      0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f,
	      keccak256(abi.encodePacked(token0, token1)),
	      hex'96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f' // init code hash
	  ))));

	*/
	hash := sha3.NewKeccakState()
	hash.Write(decodeHex((token0 + token1[2:])[2:]))

	var buf []byte

	buf = hash.Sum(nil)

	// Works up to here
	ZeroOneHash := hex.EncodeToString(buf)

	var hashInput string = "ff" + "5757371414417b8C6CAad45bAeF941aBc7d3Ab32" + ZeroOneHash + "96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"

	hashInput = strings.ToLower(hashInput)

	finalHash := sha3.NewKeccakState()

	var buf2 []byte

	finalHash.Write(decodeHex(hashInput))

	buf2 = finalHash.Sum(nil)

	// Uint Conversion

	tempUint, _ := uint256.FromHex("0x" + stateRead.RemoveLeadingZeros(hex.EncodeToString(buf2)))

	addressBytes := tempUint.Bytes20()

	byteArray := []byte{
		addressBytes[0], addressBytes[1], addressBytes[2],
		addressBytes[3],
		addressBytes[4], addressBytes[5], addressBytes[6],
		addressBytes[7], addressBytes[8], addressBytes[9],
		addressBytes[10], addressBytes[11], addressBytes[12],
		addressBytes[13], addressBytes[14], addressBytes[15],
		addressBytes[16], addressBytes[17], addressBytes[18],
		addressBytes[19]}
	encodedString := hex.EncodeToString(byteArray)
	return "0x" + encodedString

}

func decodeHex(s string) []byte {
	b, err := hex.DecodeString(s)

	if err != nil {
		panic(err)
	}

	return b
}
