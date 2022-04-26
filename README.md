# Atomic Arbitrage

Atomic Arbitrage is a base example of a bare implementation of an arbitrage bot. Written in Go based off an old example bot. Starting with UNI V2 arbs before moving onto UNI v3, Balancer, and Curve arbitrage.

## Contributing?

Just tackle a To-Do task and submit a PR, ideally once done I can make some Github actions to test speed and determine if PRs actually improve efficiency. Also if you have design/architecture improvements file and issue and I will start on them, not sure how to make public right away so starting with a conceptually simple bot then gradually evolve it over time.

## To - do

- [ ] - Re-write flash-swap to Yul+ instead of Solidity 
- [ ] - Build out an example test environment (hopefully in Foundry cramming in Yul+ via FFI and yul-log)
- [ ] - Clean up code structure
- [x] - Move off hardcoded pairs from json, use event listening with Geth to reconstruct an in-memory db of all pairs for each exchange. (Recreate from factory logs at cold boot) (credit @mempooler)
- [ ] - GraphQL interface instead of relying on JSON-RPC (although this can also be skipped by just moving straight to Geth)
- [ ] - Fuzzing Tests with Go 1.18beta2
- [ ] - Better transaction signing and construction (ie call uniswap call() first then direct it to the arb contract rather than having a contract call it first)
- [ ] - e2e tests and timing
- [ ] - Detect more swaps than just swapExactTokensForTokens (or just use a better method to detect arb opportunities, and dive into the mempool)
- [ ] - Organize Code from flat repo, and fix a lot of code-style issues
- [ ] - Goroutine stuff
- [ ] - Remove interfaces in favor of generics where possible? (I think Generics are faster?)
- [ ] - Transactions is also mostly unfinished and un optimized in any way
- [ ] - Has also occured to me file names may need to be changed to be more inline with what they actually do

## Go Notes

Some good talks if you are bored, if you want to learn how to make your Go code run faster.

- [Go perfbook](https://github.com/dgryski/go-perfbook) excellent for building an optimization routine
- [Fixing Go Garbage Collection Issues](https://www.youtube.com/watch?v=NS1hmEWv4Ac)
- [Understanding the Go Scheduler before you abuse goroutines](https://www.youtube.com/watch?v=YHRO5WQGh0k)
- [Same thing with Channels](https://www.youtube.com/watch?v=KBZlN0izeiY)
- [Very in depth look at Goroutines](https://www.youtube.com/watch?v=4CrL3Ygh7S0)
- [Understanding Go memory allocation](https://www.youtube.com/watch?v=3CR4UNMK_Is)
- [Wise words](https://go-proverbs.github.io/)
