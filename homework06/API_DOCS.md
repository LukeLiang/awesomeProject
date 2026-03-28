# Uniswap V2 Periphery 接口文档

## 项目概述

Uniswap V2 Periphery 是 Uniswap V2 去中心化交易所的外围合约集合，提供用户与 Uniswap V2 Core（交易对、工厂合约）交互的高级接口。Solidity 版本：`0.6.6`。

### 合约架构

```
├── UniswapV2Router02 (主路由合约，推荐使用)
│   └── implements IUniswapV2Router02 (extends IUniswapV2Router01)
├── UniswapV2Router01 (旧版路由，存在已知 bug，不推荐)
│   └── implements IUniswapV2Router01
├── UniswapV2Migrator (V1→V2 流动性迁移)
│   └── implements IUniswapV2Migrator
├── libraries/
│   ├── UniswapV2Library (核心计算库)
│   ├── UniswapV2OracleLibrary (价格预言机辅助库)
│   ├── UniswapV2LiquidityMathLibrary (流动性数学库)
│   └── SafeMath (安全数学库)
└── examples/
    ├── ExampleOracleSimple (固定窗口 TWAP 预言机)
    ├── ExampleSlidingWindowOracle (滑动窗口 TWAP 预言机)
    ├── ExampleFlashSwap (闪电兑换套利)
    ├── ExampleSwapToPrice (最优价格兑换)
    └── ExampleComputeLiquidityValue (流动性价值计算)
```

---

## 一、UniswapV2Router02（核心路由合约）

> 合约文件：`contracts/UniswapV2Router02.sol`
> 接口文件：`contracts/interfaces/IUniswapV2Router02.sol`（继承 `IUniswapV2Router01.sol`）

这是用户与 Uniswap V2 交互的**主入口合约**，支持添加/移除流动性、代币兑换，以及对通缩代币（fee-on-transfer tokens）的特殊支持。

### 构造函数

```solidity
constructor(address _factory, address _WETH)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_factory` | `address` | UniswapV2Factory 合约地址 |
| `_WETH` | `address` | WETH（Wrapped ETH）合约地址 |

### 状态变量（只读）

| 函数 | 返回类型 | 说明 |
|------|----------|------|
| `factory()` | `address` | UniswapV2Factory 合约地址 |
| `WETH()` | `address` | WETH 合约地址 |

---

### 1.1 添加流动性

#### `addLiquidity`

向 ERC20/ERC20 交易对添加流动性。如果交易对不存在则自动创建。

```solidity
function addLiquidity(
    address tokenA,
    address tokenB,
    uint amountADesired,
    uint amountBDesired,
    uint amountAMin,
    uint amountBMin,
    address to,
    uint deadline
) external returns (uint amountA, uint amountB, uint liquidity)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `tokenA` | `address` | 代币 A 的合约地址 |
| `tokenB` | `address` | 代币 B 的合约地址 |
| `amountADesired` | `uint` | 期望存入的代币 A 数量 |
| `amountBDesired` | `uint` | 期望存入的代币 B 数量 |
| `amountAMin` | `uint` | 代币 A 最小可接受数量（滑点保护） |
| `amountBMin` | `uint` | 代币 B 最小可接受数量（滑点保护） |
| `to` | `address` | LP Token 接收地址 |
| `deadline` | `uint` | 交易截止时间戳（过期回滚） |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| `amountA` | `uint` | 实际存入的代币 A 数量 |
| `amountB` | `uint` | 实际存入的代币 B 数量 |
| `liquidity` | `uint` | 获得的 LP Token 数量 |

**前置条件：** 调用者需先 `approve` Router 合约使用 tokenA 和 tokenB。

**错误码：**
- `UniswapV2Router: EXPIRED` — 超过 deadline
- `UniswapV2Router: INSUFFICIENT_A_AMOUNT` — amountA 低于 amountAMin
- `UniswapV2Router: INSUFFICIENT_B_AMOUNT` — amountB 低于 amountBMin

---

#### `addLiquidityETH`

向 ERC20/ETH 交易对添加流动性。ETH 通过 `msg.value` 发送，多余的 ETH 会退还。

```solidity
function addLiquidityETH(
    address token,
    uint amountTokenDesired,
    uint amountTokenMin,
    uint amountETHMin,
    address to,
    uint deadline
) external payable returns (uint amountToken, uint amountETH, uint liquidity)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `token` | `address` | ERC20 代币合约地址 |
| `amountTokenDesired` | `uint` | 期望存入的代币数量 |
| `amountTokenMin` | `uint` | 代币最小可接受数量 |
| `amountETHMin` | `uint` | ETH 最小可接受数量 |
| `to` | `address` | LP Token 接收地址 |
| `deadline` | `uint` | 交易截止时间戳 |

**返回值：** `amountToken`（实际存入代币量）、`amountETH`（实际使用 ETH 量）、`liquidity`（LP Token 数量）

**注意：** `msg.value` 必须 >= `amountETHMin`，多余的 ETH 会自动退还给 `msg.sender`。

---

### 1.2 移除流动性

#### `removeLiquidity`

从 ERC20/ERC20 交易对移除流动性，取回两种代币。

```solidity
function removeLiquidity(
    address tokenA,
    address tokenB,
    uint liquidity,
    uint amountAMin,
    uint amountBMin,
    address to,
    uint deadline
) public returns (uint amountA, uint amountB)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `tokenA` | `address` | 代币 A 地址 |
| `tokenB` | `address` | 代币 B 地址 |
| `liquidity` | `uint` | 要销毁的 LP Token 数量 |
| `amountAMin` | `uint` | 代币 A 最小取回数量 |
| `amountBMin` | `uint` | 代币 B 最小取回数量 |
| `to` | `address` | 代币接收地址 |
| `deadline` | `uint` | 交易截止时间戳 |

**返回值：** `amountA`（取回的代币 A 数量）、`amountB`（取回的代币 B 数量）

**前置条件：** 调用者需先 `approve` Router 合约使用 LP Token。

---

#### `removeLiquidityETH`

从 ERC20/ETH 交易对移除流动性，取回代币和 ETH。

```solidity
function removeLiquidityETH(
    address token,
    uint liquidity,
    uint amountTokenMin,
    uint amountETHMin,
    address to,
    uint deadline
) public returns (uint amountToken, uint amountETH)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `token` | `address` | ERC20 代币地址 |
| `liquidity` | `uint` | 要销毁的 LP Token 数量 |
| `amountTokenMin` | `uint` | 代币最小取回数量 |
| `amountETHMin` | `uint` | ETH 最小取回数量 |
| `to` | `address` | 接收地址 |
| `deadline` | `uint` | 截止时间戳 |

**返回值：** `amountToken`、`amountETH`

---

#### `removeLiquidityWithPermit`

使用 EIP-2612 permit 签名移除流动性（无需提前 approve）。

```solidity
function removeLiquidityWithPermit(
    address tokenA,
    address tokenB,
    uint liquidity,
    uint amountAMin,
    uint amountBMin,
    address to,
    uint deadline,
    bool approveMax,
    uint8 v, bytes32 r, bytes32 s
) external returns (uint amountA, uint amountB)
```

| 额外参数 | 类型 | 说明 |
|----------|------|------|
| `approveMax` | `bool` | `true` 则授权 `uint(-1)`（最大值），`false` 则授权 `liquidity` 数量 |
| `v`, `r`, `s` | `uint8`, `bytes32`, `bytes32` | EIP-712 签名参数 |

---

#### `removeLiquidityETHWithPermit`

使用 permit 签名从 ERC20/ETH 交易对移除流动性。

```solidity
function removeLiquidityETHWithPermit(
    address token,
    uint liquidity,
    uint amountTokenMin,
    uint amountETHMin,
    address to,
    uint deadline,
    bool approveMax,
    uint8 v, bytes32 r, bytes32 s
) external returns (uint amountToken, uint amountETH)
```

---

#### `removeLiquidityETHSupportingFeeOnTransferTokens`

**Router02 新增。** 为通缩代币（交易时会扣手续费的代币）移除流动性。因为代币转移时实际到账金额可能小于发送金额，此方法通过检查合约余额而非返回值来确保正确转账。

```solidity
function removeLiquidityETHSupportingFeeOnTransferTokens(
    address token,
    uint liquidity,
    uint amountTokenMin,
    uint amountETHMin,
    address to,
    uint deadline
) public returns (uint amountETH)
```

**返回值：** 只返回 `amountETH`（ETH 不会被收取转移费）。

---

#### `removeLiquidityETHWithPermitSupportingFeeOnTransferTokens`

上述方法的 permit 版本。

```solidity
function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(
    address token,
    uint liquidity,
    uint amountTokenMin,
    uint amountETHMin,
    address to,
    uint deadline,
    bool approveMax,
    uint8 v, bytes32 r, bytes32 s
) external returns (uint amountETH)
```

---

### 1.3 代币兑换

#### `swapExactTokensForTokens`

用**精确数量的输入代币**兑换**尽可能多的输出代币**。

```solidity
function swapExactTokensForTokens(
    uint amountIn,
    uint amountOutMin,
    address[] calldata path,
    address to,
    uint deadline
) external returns (uint[] memory amounts)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `amountIn` | `uint` | 精确的输入代币数量 |
| `amountOutMin` | `uint` | 最小可接受的输出数量（滑点保护） |
| `path` | `address[]` | 兑换路径，如 `[tokenA, tokenB]` 或 `[tokenA, tokenC, tokenB]`（多跳） |
| `to` | `address` | 输出代币接收地址 |
| `deadline` | `uint` | 截止时间戳 |

**返回值：** `amounts` — 路径上每一步的代币数量数组，`amounts[0]` = 输入量，`amounts[last]` = 输出量。

**错误码：** `UniswapV2Router: INSUFFICIENT_OUTPUT_AMOUNT`

---

#### `swapTokensForExactTokens`

用**尽可能少的输入代币**兑换**精确数量的输出代币**。

```solidity
function swapTokensForExactTokens(
    uint amountOut,
    uint amountInMax,
    address[] calldata path,
    address to,
    uint deadline
) external returns (uint[] memory amounts)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `amountOut` | `uint` | 精确的期望输出数量 |
| `amountInMax` | `uint` | 最大可接受的输入数量 |
| `path` | `address[]` | 兑换路径 |
| `to` | `address` | 接收地址 |
| `deadline` | `uint` | 截止时间戳 |

**返回值：** `amounts` — 路径上每步的数量数组。

**错误码：** `UniswapV2Router: EXCESSIVE_INPUT_AMOUNT`

---

#### `swapExactETHForTokens`

用**精确数量的 ETH** 兑换代币。ETH 通过 `msg.value` 发送。

```solidity
function swapExactETHForTokens(
    uint amountOutMin,
    address[] calldata path,
    address to,
    uint deadline
) external payable returns (uint[] memory amounts)
```

**约束：** `path[0]` 必须为 WETH 地址。

---

#### `swapTokensForExactETH`

用代币兑换**精确数量的 ETH**。

```solidity
function swapTokensForExactETH(
    uint amountOut,
    uint amountInMax,
    address[] calldata path,
    address to,
    uint deadline
) external returns (uint[] memory amounts)
```

**约束：** `path[last]` 必须为 WETH 地址。

---

#### `swapExactTokensForETH`

用**精确数量的代币**兑换尽可能多的 ETH。

```solidity
function swapExactTokensForETH(
    uint amountIn,
    uint amountOutMin,
    address[] calldata path,
    address to,
    uint deadline
) external returns (uint[] memory amounts)
```

**约束：** `path[last]` 必须为 WETH 地址。

---

#### `swapETHForExactTokens`

用 ETH 兑换**精确数量的代币**。多余的 ETH 会退还。

```solidity
function swapETHForExactTokens(
    uint amountOut,
    address[] calldata path,
    address to,
    uint deadline
) external payable returns (uint[] memory amounts)
```

**约束：** `path[0]` 必须为 WETH 地址。

---

### 1.4 支持通缩代币的兑换方法（Router02 新增）

这些方法适用于交易时会自动扣除手续费的代币（如 SafeMoon 类代币），不依赖于 `getAmountsOut` 的预计算结果，而是通过比较余额差来判断实际到账金额。**这些方法不返回 amounts 数组。**

#### `swapExactTokensForTokensSupportingFeeOnTransferTokens`

```solidity
function swapExactTokensForTokensSupportingFeeOnTransferTokens(
    uint amountIn,
    uint amountOutMin,
    address[] calldata path,
    address to,
    uint deadline
) external
```

---

#### `swapExactETHForTokensSupportingFeeOnTransferTokens`

```solidity
function swapExactETHForTokensSupportingFeeOnTransferTokens(
    uint amountOutMin,
    address[] calldata path,
    address to,
    uint deadline
) external payable
```

---

#### `swapExactTokensForETHSupportingFeeOnTransferTokens`

```solidity
function swapExactTokensForETHSupportingFeeOnTransferTokens(
    uint amountIn,
    uint amountOutMin,
    address[] calldata path,
    address to,
    uint deadline
) external
```

---

### 1.5 价格计算（只读函数）

#### `quote`

根据储备量计算等价代币数量（不考虑手续费）。

```solidity
function quote(uint amountA, uint reserveA, uint reserveB) public pure returns (uint amountB)
```

**公式：** `amountB = amountA * reserveB / reserveA`

---

#### `getAmountOut`

给定输入数量和储备量，计算最大输出数量（含 0.3% 手续费）。

```solidity
function getAmountOut(uint amountIn, uint reserveIn, uint reserveOut) public pure returns (uint amountOut)
```

**公式：** `amountOut = (amountIn * 997 * reserveOut) / (reserveIn * 1000 + amountIn * 997)`

---

#### `getAmountIn`

给定期望输出数量和储备量，计算所需最小输入数量。

```solidity
function getAmountIn(uint amountOut, uint reserveIn, uint reserveOut) public pure returns (uint amountIn)
```

**公式：** `amountIn = (reserveIn * amountOut * 1000) / ((reserveOut - amountOut) * 997) + 1`

---

#### `getAmountsOut`

沿路径执行链式 `getAmountOut` 计算。

```solidity
function getAmountsOut(uint amountIn, address[] memory path) public view returns (uint[] memory amounts)
```

---

#### `getAmountsIn`

沿路径执行链式 `getAmountIn` 计算。

```solidity
function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts)
```

---

## 二、UniswapV2Migrator（V1→V2 迁移合约）

> 合约文件：`contracts/UniswapV2Migrator.sol`

将 Uniswap V1 的流动性一键迁移到 V2。

### 构造函数

```solidity
constructor(address _factoryV1, address _router)
```

### `migrate`

```solidity
function migrate(
    address token,
    uint amountTokenMin,
    uint amountETHMin,
    address to,
    uint deadline
) external
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `token` | `address` | ERC20 代币地址 |
| `amountTokenMin` | `uint` | V2 添加流动性时代币最小数量 |
| `amountETHMin` | `uint` | V2 添加流动性时 ETH 最小数量 |
| `to` | `address` | V2 LP Token 接收地址 |
| `deadline` | `uint` | 截止时间戳 |

**流程：**
1. 从 V1 Exchange 转入用户的 LP Token
2. 调用 V1 的 `removeLiquidity` 取出 ETH 和代币
3. 调用 Router 的 `addLiquidityETH` 向 V2 添加流动性
4. 退还多余的代币或 ETH

**前置条件：** 用户需先 `approve` Migrator 合约使用 V1 LP Token。

---

## 三、UniswapV2Library（核心计算库）

> 文件：`contracts/libraries/UniswapV2Library.sol`

所有函数均为 `internal`，供其他合约调用。

| 函数 | 说明 |
|------|------|
| `sortTokens(tokenA, tokenB) -> (token0, token1)` | 按地址排序代币对，确保 `token0 < token1` |
| `pairFor(factory, tokenA, tokenB) -> pair` | 使用 CREATE2 计算交易对地址（无需链上调用） |
| `getReserves(factory, tokenA, tokenB) -> (reserveA, reserveB)` | 获取并排序交易对的储备量 |
| `quote(amountA, reserveA, reserveB) -> amountB` | 等价计算（不含手续费） |
| `getAmountOut(amountIn, reserveIn, reserveOut) -> amountOut` | 含 0.3% 手续费的输出计算 |
| `getAmountIn(amountOut, reserveIn, reserveOut) -> amountIn` | 含 0.3% 手续费的反向计算 |
| `getAmountsOut(factory, amountIn, path) -> amounts[]` | 链式正向计算 |
| `getAmountsIn(factory, amountOut, path) -> amounts[]` | 链式反向计算 |

**init code hash：** `96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f`
（`pairFor` 使用此哈希值计算 CREATE2 地址，不同部署需修改此值）

---

## 四、UniswapV2OracleLibrary（价格预言机库）

> 文件：`contracts/libraries/UniswapV2OracleLibrary.sol`

| 函数 | 说明 |
|------|------|
| `currentBlockTimestamp() -> uint32` | 当前区块时间戳（取模 2^32 防溢出） |
| `currentCumulativePrices(pair) -> (price0Cumulative, price1Cumulative, blockTimestamp)` | 获取当前累计价格，含未写入链上的增量（counterfactual），节省 gas |

---

## 五、UniswapV2LiquidityMathLibrary（流动性数学库）

> 文件：`contracts/libraries/UniswapV2LiquidityMathLibrary.sol`

| 函数 | 说明 |
|------|------|
| `computeProfitMaximizingTrade(truePriceTokenA, truePriceTokenB, reserveA, reserveB) -> (aToB, amountIn)` | 计算将池子价格推到外部真实价格所需的最优套利交易方向和数量 |
| `getReservesAfterArbitrage(factory, tokenA, tokenB, truePriceTokenA, truePriceTokenB) -> (reserveA, reserveB)` | 计算套利后的储备量 |
| `computeLiquidityValue(reservesA, reservesB, totalSupply, liquidityAmount, feeOn, kLast) -> (tokenAAmount, tokenBAmount)` | 计算给定 LP Token 对应的底层代币价值（考虑协议费） |
| `getLiquidityValue(factory, tokenA, tokenB, liquidityAmount) -> (tokenAAmount, tokenBAmount)` | 获取当前流动性价值（**注意：易受三明治攻击**） |
| `getLiquidityValueAfterArbitrageToPrice(factory, tokenA, tokenB, truePriceTokenA, truePriceTokenB, liquidityAmount) -> (tokenAAmount, tokenBAmount)` | 基于外部真实价格计算流动性价值（抗操纵） |

---

## 六、SafeMath（安全数学库）

> 文件：`contracts/libraries/SafeMath.sol`

| 函数 | 说明 |
|------|------|
| `add(uint x, uint y) -> uint z` | 安全加法，溢出时回滚（`ds-math-add-overflow`） |
| `sub(uint x, uint y) -> uint z` | 安全减法，下溢时回滚（`ds-math-sub-underflow`） |
| `mul(uint x, uint y) -> uint z` | 安全乘法，溢出时回滚（`ds-math-mul-overflow`） |

---

## 七、示例合约

### 7.1 ExampleOracleSimple（固定窗口 TWAP 预言机）

> 文件：`contracts/examples/ExampleOracleSimple.sol`

基于累计价格的时间加权平均价格（TWAP）预言机，固定窗口 = 24 小时。**每个交易对需单独部署。**

```solidity
constructor(address factory, address tokenA, address tokenB)
```

| 函数 | 说明 |
|------|------|
| `update()` | 更新价格（每 24 小时至少调用一次） |
| `consult(address token, uint amountIn) -> uint amountOut` | 查询 token 的 TWAP 价格。首次 `update()` 前返回 0 |

**常量：**
- `PERIOD = 24 hours`

**状态变量：**
- `pair` — 交易对地址
- `token0` / `token1` — 排序后的代币地址
- `price0CumulativeLast` / `price1CumulativeLast` — 上次记录的累计价格
- `blockTimestampLast` — 上次更新时间戳
- `price0Average` / `price1Average` — 计算得出的平均价格（FixedPoint.uq112x112）

---

### 7.2 ExampleSlidingWindowOracle（滑动窗口 TWAP 预言机）

> 文件：`contracts/examples/ExampleSlidingWindowOracle.sol`

更精确的移动平均价格预言机，支持配置窗口大小和粒度。**单例部署，支持所有交易对。**

```solidity
constructor(address factory_, uint windowSize_, uint8 granularity_)
```

| 参数 | 说明 |
|------|------|
| `factory_` | UniswapV2Factory 地址 |
| `windowSize_` | 窗口大小（秒），如 86400（24h） |
| `granularity_` | 观测粒度（必须 > 1），如 24 表示每小时一个观测点 |

| 函数 | 说明 |
|------|------|
| `observationIndexOf(uint timestamp) -> uint8` | 返回时间戳对应的观测索引 |
| `update(address tokenA, address tokenB)` | 更新指定交易对的价格观测 |
| `consult(address tokenIn, uint amountIn, address tokenOut) -> uint amountOut` | 查询移动平均价格 |

**精度范围：** `[windowSize - periodSize * 2, windowSize]`

**数据结构：**
```solidity
struct Observation {
    uint timestamp;
    uint price0Cumulative;
    uint price1Cumulative;
}
```

---

### 7.3 ExampleFlashSwap（闪电兑换套利）

> 文件：`contracts/examples/ExampleFlashSwap.sol`

利用 V2 闪电贷在 V1 和 V2 之间套利。

```solidity
constructor(address _factory, address _factoryV1, address router)
```

| 函数 | 说明 |
|------|------|
| `uniswapV2Call(address sender, uint amount0, uint amount1, bytes calldata data)` | V2 Pair 回调函数，执行套利逻辑。`data` 编码了 V1 的滑点参数 |

**工作流程：**
1. 从 V2 闪电借出代币/WETH
2. 在 V1 上兑换获取 ETH/代币
3. 偿还 V2 闪电贷
4. 保留利润

---

### 7.4 ExampleSwapToPrice（价格收敛兑换）

> 文件：`contracts/examples/ExampleSwapToPrice.sol`

计算并执行利润最大化的兑换，将池子价格推向外部观测的真实价格。

```solidity
constructor(address factory_, IUniswapV2Router01 router_)
```

#### `swapToPrice`

```solidity
function swapToPrice(
    address tokenA,
    address tokenB,
    uint256 truePriceTokenA,
    uint256 truePriceTokenB,
    uint256 maxSpendTokenA,
    uint256 maxSpendTokenB,
    address to,
    uint256 deadline
) public
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `tokenA` | `address` | 代币 A 地址 |
| `tokenB` | `address` | 代币 B 地址 |
| `truePriceTokenA` | `uint256` | 代币 A 的真实价格（比率分子） |
| `truePriceTokenB` | `uint256` | 代币 B 的真实价格（比率分母） |
| `maxSpendTokenA` | `uint256` | 最多花费的代币 A 数量（设 0 表示不卖 A） |
| `maxSpendTokenB` | `uint256` | 最多花费的代币 B 数量（设 0 表示不卖 B） |
| `to` | `address` | 输出代币接收地址 |
| `deadline` | `uint256` | 截止时间戳 |

**前置条件：** 调用者需 `approve` 此合约使用对应代币。

---

### 7.5 ExampleComputeLiquidityValue（流动性价值计算）

> 文件：`contracts/examples/ExampleComputeLiquidityValue.sol`

`UniswapV2LiquidityMathLibrary` 的外部包装合约。

```solidity
constructor(address factory_)
```

| 函数 | 说明 |
|------|------|
| `getReservesAfterArbitrage(tokenA, tokenB, truePriceTokenA, truePriceTokenB) -> (reserveA, reserveB)` | 套利后储备量 |
| `getLiquidityValue(tokenA, tokenB, liquidityAmount) -> (tokenAAmount, tokenBAmount)` | 当前流动性价值 |
| `getLiquidityValueAfterArbitrageToPrice(tokenA, tokenB, truePriceTokenA, truePriceTokenB, liquidityAmount) -> (tokenAAmount, tokenBAmount)` | 基于真实价格的流动性价值 |
| `getGasCostOfGetLiquidityValueAfterArbitrageToPrice(...) -> uint256` | Gas 消耗测试 |

---

## 八、辅助接口

### IERC20

> 文件：`contracts/interfaces/IERC20.sol`

```solidity
interface IERC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
    function totalSupply() external view returns (uint);
    function balanceOf(address owner) external view returns (uint);
    function allowance(address owner, address spender) external view returns (uint);
    function approve(address spender, uint value) external returns (bool);
    function transfer(address to, uint value) external returns (bool);
    function transferFrom(address from, address to, uint value) external returns (bool);

    event Approval(address indexed owner, address indexed spender, uint value);
    event Transfer(address indexed from, address indexed to, uint value);
}
```

### IWETH

> 文件：`contracts/interfaces/IWETH.sol`

```solidity
interface IWETH {
    function deposit() external payable;       // ETH -> WETH
    function transfer(address to, uint value) external returns (bool);
    function withdraw(uint) external;           // WETH -> ETH
}
```

### IUniswapV1Factory

> 文件：`contracts/interfaces/V1/IUniswapV1Factory.sol`

```solidity
interface IUniswapV1Factory {
    function getExchange(address) external view returns (address);
}
```

### IUniswapV1Exchange

> 文件：`contracts/interfaces/V1/IUniswapV1Exchange.sol`

```solidity
interface IUniswapV1Exchange {
    function balanceOf(address owner) external view returns (uint);
    function transferFrom(address from, address to, uint value) external returns (bool);
    function removeLiquidity(uint, uint, uint, uint) external returns (uint, uint);
    function tokenToEthSwapInput(uint, uint, uint) external returns (uint);
    function ethToTokenSwapInput(uint, uint) external payable returns (uint);
}
```

---

## 九、错误码汇总

| 错误码 | 触发场景 |
|--------|----------|
| `UniswapV2Router: EXPIRED` | 交易超过 deadline 时间戳 |
| `UniswapV2Router: INSUFFICIENT_A_AMOUNT` | 代币 A 实际数量低于 amountAMin |
| `UniswapV2Router: INSUFFICIENT_B_AMOUNT` | 代币 B 实际数量低于 amountBMin |
| `UniswapV2Router: INSUFFICIENT_OUTPUT_AMOUNT` | 兑换输出低于 amountOutMin |
| `UniswapV2Router: EXCESSIVE_INPUT_AMOUNT` | 兑换输入超过 amountInMax |
| `UniswapV2Router: INVALID_PATH` | ETH 相关兑换路径首/尾不是 WETH |
| `UniswapV2Library: IDENTICAL_ADDRESSES` | tokenA == tokenB |
| `UniswapV2Library: ZERO_ADDRESS` | 代币地址为零 |
| `UniswapV2Library: INSUFFICIENT_AMOUNT` | quote 输入为 0 |
| `UniswapV2Library: INSUFFICIENT_LIQUIDITY` | 储备量为 0 |
| `UniswapV2Library: INSUFFICIENT_INPUT_AMOUNT` | getAmountOut 输入为 0 |
| `UniswapV2Library: INSUFFICIENT_OUTPUT_AMOUNT` | getAmountIn 输出为 0 |
| `UniswapV2Library: INVALID_PATH` | 路径长度 < 2 |
| `ds-math-add-overflow` | SafeMath 加法溢出 |
| `ds-math-sub-underflow` | SafeMath 减法下溢 |
| `ds-math-mul-overflow` | SafeMath 乘法溢出 |
| `SlidingWindowOracle: GRANULARITY` | 粒度必须 > 1 |
| `SlidingWindowOracle: WINDOW_NOT_EVENLY_DIVISIBLE` | 窗口大小不能被粒度整除 |
| `SlidingWindowOracle: MISSING_HISTORICAL_OBSERVATION` | 缺少历史观测数据 |
| `ExampleOracleSimple: PERIOD_NOT_ELAPSED` | 未达到更新周期 |
| `ExampleOracleSimple: NO_RESERVES` | 交易对无储备量 |
| `ExampleOracleSimple: INVALID_TOKEN` | 查询的代币不属于该交易对 |
| `ExampleSwapToPrice: ZERO_PRICE` | 真实价格参数为 0 |
| `ExampleSwapToPrice: ZERO_SPEND` | 两个 maxSpend 均为 0 |
| `ExampleSwapToPrice: ZERO_AMOUNT_IN` | 计算出的交易量为 0 |

---

## 十、已知问题

**UniswapV2Router01 `getAmountIn` bug：** Router01 的 `getAmountIn` 函数错误地调用了 `UniswapV2Library.getAmountOut` 而非 `UniswapV2Library.getAmountIn`，导致返回值不正确。此 bug 已在 Router02 中修复。**建议始终使用 Router02。**
