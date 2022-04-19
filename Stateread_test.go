package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/holiman/uint256"
)

func TestGetReserves(t *testing.T) {
	// read config
	config, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	var testCases []string = []string{
		"0xE62Ec2e799305E0D367b0Cc3ee2CdA135bF89816",
		"0x1E67124681b402064CD0ABE8ed1B5c79D2e02f64",
		"0x2813D43463C374a680f235c428FB1D7f08dE0B69",
		"0xadbf1854e5883eb8aa7baf50705338739e558e5b",
		"0xdc9232e2df177d7a12fdff6ecbab114e2231198d",
		"0x853ee4b2a13f8a742d64c8f088be7ba2131f670d",
		"0xE62Ec2e799305E0D367b0Cc3ee2CdA135bF89816",
		"0x1E67124681b402064CD0ABE8ed1B5c79D2e02f64",
		"0x2813D43463C374a680f235c428FB1D7f08dE0B69",
		"0xadbf1854e5883eb8aa7baf50705338739e558e5b",
		"0xdc9232e2df177d7a12fdff6ecbab114e2231198d",
		"0x853ee4b2a13f8a742d64c8f088be7ba2131f670d",
		"0xE62Ec2e799305E0D367b0Cc3ee2CdA135bF89816",
		"0x1E67124681b402064CD0ABE8ed1B5c79D2e02f64",
		"0x2813D43463C374a680f235c428FB1D7f08dE0B69",
		"0xadbf1854e5883eb8aa7baf50705338739e558e5b",
		"0xdc9232e2df177d7a12fdff6ecbab114e2231198d",
		"0x853ee4b2a13f8a742d64c8f088be7ba2131f670d",
	}

	pipe := make(chan [2]uint256.Int, 6)
	start := time.Now()

	for index := 0; index < len(testCases)-1; index++ {
		go func() {
			reserve0, reserve1 := GetReserves(testCases[index], config.HttpURL)
			pipe <- [2]uint256.Int{reserve0, reserve1}
		}()
	}
	// Wait for routines to finish
	for index := 0; index < len(testCases)-1; index++ {
		<-pipe
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed / 18)
}

func TestDownloadBlock(t *testing.T) {
	// read config
	config, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	start := time.Now()

	DownloadBlock("latest", config.HttpURL)

	elapsed := time.Since(start)

	fmt.Println(elapsed)
}
