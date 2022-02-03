package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/holiman/uint256"
)

func TestFullMul(t *testing.T) {
	z := uint256.NewInt(100)
	y := uint256.NewInt(1000)

	expectedResult := uint256.NewInt(100000)

	expectedResult2 := uint256.NewInt(0)
	result1, result2 := FullMul(z, y)
	if !result1.Eq(expectedResult) {
		t.Errorf("Multiplication was incorrect, got: %d, want: %d.", result1, expectedResult)
	}
	if !result2.Eq(expectedResult2) {
		t.Errorf("expectedResult was incorrect, got: %d, want: %d.", result2, expectedResult2)

	}
}

func TestFullDiv(t *testing.T) {
	expectedResult1, _ := uint256.FromHex("0x55555555555555555555555555555555555555555555555555555555555555FC")
	z := uint256.NewInt(1000)
	y := uint256.NewInt(100)
	x := uint256.NewInt(6)
	var result uint256.Int = FullDiv(z, y, x)
	if !result.Eq(expectedResult1) {
		t.Errorf("Division was incorrect, got: %d, want: %d.", result, expectedResult1)
	}
}

func TestMulDiv(t *testing.T) {
	var testCases [][]uint64 = [][]uint64{{500000, 600000, 100, 3000000000}, {70000, 6000, 1000, 420000}}
	for index := 0; index < len(testCases)-1; index++ {
		testCase := testCases[index]
		result, err := MulDiv(*uint256.NewInt(testCase[0]), *uint256.NewInt(testCase[1]), *uint256.NewInt(testCase[2]))
		if err != nil {
			t.Errorf("MULDIV Exited with : %d", err)
		}
		if !result.Eq(uint256.NewInt(testCase[3])) {
			t.Errorf("Division was incorrect, got: %d, want: %d.", result, testCase[2])
		}
	}
}

func TestBabylonianSqrt(t *testing.T) {
	y := uint256.NewInt(50000)

	expectedResult := uint256.NewInt(223)

	result := BabylonianSqrt(y)

	if !expectedResult.Eq(&result) {
		t.Errorf("Square Root was incorrect, got: %d, want: %d.", result, expectedResult)
	}

}

func TestComputeProfitMaximizingTrade(t *testing.T) {
	//generate random tests for 1 mil iterations
	iterations := 1000000
	sliceOfRandomTests := generateSliceOfRandomTests(iterations)
	//init the waitgroup
	var WaitGroup sync.WaitGroup

	//start the timer and begin the speed test
	start := time.Now()

	for i := 0; i < iterations-1; i++ {
		randomIntTest := sliceOfRandomTests[i]
		go wrapper(randomIntTest[0], randomIntTest[1], randomIntTest[2], randomIntTest[3], &WaitGroup)
	}
	//wait for all go routines to finish
	WaitGroup.Wait()
	elapsed := time.Since(start)
	fmt.Println("Time elapsed", elapsed)
}

func generateSliceOfRandomTests(iterations int) [][]*uint256.Int {
	fmt.Println("Generating Random Numbers...")
	sliceOfRandomTests := &[][]*uint256.Int{}
	for i := 0; i < iterations; i++ {
		generateRandomInts(sliceOfRandomTests, iterations, i)
	}
	return *sliceOfRandomTests
}

func generateRandomInts(sliceOfRandomTests *[][]*uint256.Int, iterations int, i int) {

	randomInts := []*uint256.Int{}
	for j := 0; j < 4; j++ {
		randomInt := rand.Int63n(10000) * 1000
		randomUint256, _ := uint256.FromBig(big.NewInt(randomInt))
		multiplierBigInt := big.NewInt(int64(math.Pow(10, 15)))
		multiplier256, _ := uint256.FromBig(multiplierBigInt)
		randomUint256 = randomUint256.Mul(randomUint256, multiplier256)
		randomInts = append(randomInts, randomUint256)
	}
	// fmt.Println(randomInts)
	*sliceOfRandomTests = append(*sliceOfRandomTests, randomInts)
}
func wrapper(sReserve0 *uint256.Int, sReserve1 *uint256.Int, uReserve0 *uint256.Int, uReserve1 *uint256.Int, wg *sync.WaitGroup) {
	wg.Add(1)
	// computeProfitMaximizingTrade(sReserve0,sReserve1,uReserve0,uReserve1)
	// 50000, 3000,50000, 100000
	ComputeProfitMaximizingTrade(*sReserve0, *sReserve1, *uReserve0, *uReserve1)
	// if !aToB {
	// 	fmt.Println("aToB was incorrectly identified")
	// }
	// expectedResult := uint256.NewInt(238959)
	// if !Number.Eq(expectedResult) {
	// 	fmt.Printf("Input was incorrect, got: %d, want: %d.\n", Number, expectedResult)
	// }
	wg.Done()

}
