package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 查询区块
// 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
// 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
// 输出查询结果到控制台。
// SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/_ghPu9oQgf-QhS6P4-xgW
// PRIVATE_KEY_OWNER=7f2b62396e9416d53bef6db02ce039f2f370977cbe322
// https://sepolia.infura.io/v3/2c23a3fc88434992a0277707c5f7f753
func main() {
	//blockNumberFlag := flag.Uint64("number", 0, "Block number to query")
	addrHex := flag.String("address", "", "account address (required)")
	blockNumber := flag.Int64("block", -1, "block number to query (-1 means latest)")
	flag.Parse()

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("Please set the ETH_RPC_URL environment variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatal("failed to get latest block: %v", err)
	}
	defer client.Close()

	// task 1
	// 最新区块
	latestBlock, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get latest block: %v", err)
	}

	printBlockInfo("Latest Block", latestBlock)

	// task 2
	address := common.HexToAddress(*addrHex)

	var blockNum *big.Int
	if *blockNumber >= 0 {
		blockNum = big.NewInt(*blockNumber)
	}

	balanceWei, err := client.BalanceAt(ctx, address, blockNum)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	fmt.Println("=== Account Balance ===")
	fmt.Printf("Address     : %s\n", address.Hex())

	if blockNum == nil {
		fmt.Printf("Block       : latest\n")
	} else {
		fmt.Printf("Block       : %d\n", blockNum.Uint64())
	}
	fmt.Printf("Balance Wei : %s\n", balanceWei.String())

	balanceEth := weiToEth(balanceWei)
	fmt.Printf("Balance ETH : %s\n", balanceEth.Text('f', 6))
}

func weiToEth(wei *big.Int) *big.Float {
	fWei := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(fWei, big.NewFloat(math.Pow10(18)))
	return ethValue
}

func printBlockInfo(title string, block *types.Block) {
	fmt.Println("======================================")
	fmt.Println(title)
	fmt.Println("======================================")
	fmt.Printf("Block: %+v\n", block)

	// 基本信息
	fmt.Printf("Block Number: %d\n", block.Number().Uint64())
	fmt.Printf("Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", block.ParentHash().Hex())

	// 时间信息
	blockTime := time.Unix(int64(block.Time()), 0)
	fmt.Printf("Time: %s\n", blockTime.Format(time.RFC3339))
	fmt.Printf("Timestamp (Local): %s\n", blockTime.Local().Format("2006-01-02 15:04:05 MST"))

	// Gas 信息
	gasUsed := block.GasUsed()
	gasLimit := block.GasLimit()
	gasUsagePercent := float64(gasUsed) / float64(gasLimit) * 100
	fmt.Printf("Gas Used     : %d (%.2f%%)\n", gasUsed, gasUsagePercent)
	fmt.Printf("Gas Limit    : %d\n", gasLimit)

	// 交易信息
	txCount := len(block.Transactions())
	fmt.Printf("Tx Count: %d\n", txCount)

	// 区块根信息 (Merkle 树根)
	fmt.Printf("State Root   : %s\n", block.Root().Hex())
	fmt.Printf("Tx Root      : %s\n", block.TxHash().Hex())
	fmt.Printf("Receipt Root : %s\n", block.ReceiptHash().Hex())

	// 区块大小估算（简化版，实际大小还包括其他字段）
	if txCount > 0 {
		fmt.Printf("\nFirst Tx Hash: %s\n", block.Transactions()[0].Hash().Hex())
		if txCount > 1 {
			fmt.Printf("Last Tx Hash : %s\n", block.Transactions()[txCount-1].Hash().Hex())
		}
	}

	// 难度信息（PoW 相关，PoS 后基本固定）
	fmt.Printf("Difficulty   : %s\n", block.Difficulty().String())

	// 区块奖励相关信息
	coinbase := block.Coinbase()
	if coinbase != (common.Address{}) {
		fmt.Printf("Coinbase     : %s\n", coinbase.Hex())
	}

	fmt.Println("======================================")
	fmt.Println()
}
