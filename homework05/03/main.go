package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	sepoliaChainID = 11155111

	// Counter 合约 ABI
	counterABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"incrementer","type":"address"},{"indexed":false,"internalType":"uint256","name":"newCount","type":"uint256"}],"name":"CountIncremented","type":"event"},{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"getCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"increment","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`
)

func main() {
	// ─────────────────────────────────────────────────────────────
	// 1. 读取配置（环境变量）
	// ─────────────────────────────────────────────────────────────
	rpcURL := mustEnv("SEPOLIA_RPC_URL",
		"请设置 SEPOLIA_RPC_URL，例如: https://sepolia.infura.io/v3/<PROJECT_ID>")
	privateKeyHex := mustEnv("PRIVATE_KEY",
		"请设置 PRIVATE_KEY（钱包私钥，不含 0x 前缀）")
	contractAddrStr := mustEnv("CONTRACT_ADDRESS",
		"请设置 CONTRACT_ADDRESS（已部署的合约地址）")

	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("  Counter 合约交互示例 - Sepolia 测试网")
	fmt.Println("═══════════════════════════════════════════════════")

	// ─────────────────────────────────────────────────────────────
	// 2. 连接到 Sepolia 测试网
	// ─────────────────────────────────────────────────────────────
	fmt.Printf("\n[1/4] 连接到 Sepolia: %s\n", rpcURL)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接 Sepolia 失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("获取 ChainID 失败: %v", err)
	}
	if chainID.Int64() != sepoliaChainID {
		log.Fatalf("Chain ID 不匹配: 期望 %d (Sepolia), 实际 %s", sepoliaChainID, chainID)
	}
	fmt.Printf("    ✓ 已连接，Chain ID: %s (Sepolia)\n", chainID)

	// ─────────────────────────────────────────────────────────────
	// 3. 加载私钥
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[2/4] 加载钱包私钥")
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}
	fromAddr := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	fmt.Printf("    ✓ 钱包地址: %s\n", fromAddr.Hex())

	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		log.Fatalf("查询余额失败: %v", err)
	}
	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("    ✓ 钱包余额: %.6f SepoliaETH\n", ethBalance)

	// ─────────────────────────────────────────────────────────────
	// 4. 解析 ABI，绑定合约地址
	// ─────────────────────────────────────────────────────────────
	fmt.Printf("\n[3/4] 解析 ABI，绑定合约: %s\n", contractAddrStr)
	parsedABI, err := abi.JSON(strings.NewReader(counterABI))
	if err != nil {
		log.Fatalf("解析 ABI 失败: %v", err)
	}
	contractAddr := common.HexToAddress(contractAddrStr)

	// 用 ABI 调用 owner()（只读）验证合约可访问
	owner, err := callView[common.Address](ctx, client, parsedABI, contractAddr, "owner")
	if err != nil {
		log.Fatalf("读取 owner 失败: %v", err)
	}
	fmt.Printf("    ✓ 合约 Owner: %s\n", owner.Hex())

	// ─────────────────────────────────────────────────────────────
	// 5. 读取初始值 → 调用 increment() → 读取结果
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[4/4] 调用合约方法")

	// 用 ABI Pack 编码 getCount() 调用数据，通过 CallContract 读取
	countBefore, err := callView[*big.Int](ctx, client, parsedABI, contractAddr, "getCount")
	if err != nil {
		log.Fatalf("读取计数器失败: %v", err)
	}
	fmt.Printf("    递增前计数值: %s\n", countBefore)

	// 用 ABI Pack 编码 increment() 调用数据，手动构造并发送交易
	tx, err := sendTx(ctx, client, privateKey, parsedABI, contractAddr, "increment")
	if err != nil {
		log.Fatalf("调用 increment 失败: %v", err)
	}
	fmt.Printf("    交易已发送，哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("    Etherscan: https://sepolia.etherscan.io/tx/%s\n", tx.Hash().Hex())

	fmt.Println("    等待交易确认...")
	if err := waitMined(ctx, client, tx.Hash()); err != nil {
		log.Fatalf("等待交易确认超时: %v", err)
	}

	countAfter, err := callView[*big.Int](ctx, client, parsedABI, contractAddr, "getCount")
	if err != nil {
		log.Fatalf("读取递增后计数器失败: %v", err)
	}

	// ─────────────────────────────────────────────────────────────
	// 6. 输出结果
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("  执行结果")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("  合约地址:      %s\n", contractAddr.Hex())
	fmt.Printf("  Etherscan:     https://sepolia.etherscan.io/address/%s\n", contractAddr.Hex())
	fmt.Printf("  递增前计数值:  %s\n", countBefore)
	fmt.Printf("  递增后计数值:  %s\n", countAfter)
	fmt.Printf("  变化量:        +%s\n", new(big.Int).Sub(countAfter, countBefore))
	fmt.Println("═══════════════════════════════════════════════════")
}

// callView 用 ABI Pack 编码只读调用，通过 eth_call 读取返回值。
func callView[T any](ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, addr common.Address, method string, args ...interface{}) (T, error) {
	var zero T

	// Pack：将方法名 + 参数编码为 ABI calldata
	callData, err := parsedABI.Pack(method, args...)
	if err != nil {
		return zero, fmt.Errorf("abi.Pack(%s): %w", method, err)
	}

	// eth_call：不发交易，直接读取链上状态
	raw, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &addr,
		Data: callData,
	}, nil /* latest block */)
	if err != nil {
		return zero, fmt.Errorf("eth_call(%s): %w", method, err)
	}

	// Unpack：将返回的字节解码为 Go 类型
	outputs, err := parsedABI.Unpack(method, raw)
	if err != nil {
		return zero, fmt.Errorf("abi.Unpack(%s): %w", method, err)
	}
	return outputs[0].(T), nil
}

// sendTx 用 ABI Pack 编码写操作，手动构造 EIP-1559 交易并发送。
func sendTx(ctx context.Context, client *ethclient.Client, privateKey *ecdsa.PrivateKey, parsedABI abi.ABI, addr common.Address, method string, args ...interface{}) (*types.Transaction, error) {
	fromAddr := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))

	// Pack：将方法名 + 参数编码为 ABI calldata
	callData, err := parsedABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack(%s): %w", method, err)
	}

	// 获取 nonce（pending 状态，包含尚未上链的交易）
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("获取 nonce 失败: %w", err)
	}

	// EIP-1559 gas 费用估算
	gasTipCap, err := client.SuggestGasTipCap(ctx) // priority fee（矿工小费）
	if err != nil {
		return nil, fmt.Errorf("获取 gasTipCap 失败: %w", err)
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("获取区块头失败: %w", err)
	}
	// gasFeeCap = 2 × baseFee + tip，预留 base fee 上涨的缓冲
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(header.BaseFee, big.NewInt(2)),
		gasTipCap,
	)

	// 构造 EIP-1559 交易
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(sepoliaChainID),
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       100000,
		To:        &addr,
		Value:     big.NewInt(0),
		Data:      callData,
	})

	// 用私钥签名交易
	signer := types.LatestSignerForChainID(big.NewInt(sepoliaChainID))
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %w", err)
	}

	// 广播交易到网络
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("发送交易失败: %w", err)
	}
	return signedTx, nil
}

// mustEnv 从环境变量读取值，若未设置则打印提示并退出。
func mustEnv(key, hint string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "错误: 环境变量 %s 未设置\n%s\n", key, hint)
		os.Exit(1)
	}
	return v
}

// waitMined 轮询等待交易被打包进区块（最长等待 2 分钟）。
func waitMined(ctx context.Context, client *ethclient.Client, txHash common.Hash) error {
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("等待超时 (2 分钟)，交易 %s 尚未确认", txHash.Hex())
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err != nil {
				fmt.Print(".")
				continue
			}
			fmt.Println()
			if receipt.Status == 0 {
				return fmt.Errorf("交易执行失败 (Status=0)，请检查 Etherscan")
			}
			fmt.Printf("    区块高度: %d，Gas 消耗: %d\n", receipt.BlockNumber, receipt.GasUsed)
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
