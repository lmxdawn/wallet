# wallet

> 虚拟币钱包服务，转账/提现/充值/归集
>
>
> 完全实现与业务服务隔离，使用http服务相互调用

# 计划支持
- [x] 以太坊（ERC20）
- [x] 波场（TRC20）
- [x] 币安（BEP20）
- [x] OKC（KIP20）
- [x] 比特币


# 接口文档

`script/api.md`

# 下载-打包

```shell
# 拉取代码
$ git clone https://github.com/lmxdawn/wallet.git
$ cd wallet

# 打包 (-tags "doc") 可选，加上可以运行swagger
$ go build [-tags "doc"]

# 直接运行示例配置
$ wallet -c config/config-example.yml

```

# 重新生成配置

```shell
# 生成配置文件
$ vim config.yml
$ wallet -c config.yml

```

# 配置文件参数解释

|  参数名   | 描述  |
|  ----  | ----  |
| coin_name  | 币种名称 |
| contract  | 合约地址（为空表示主币） |
| contract_type  | 合约类型（波场需要区分是TRC20还是TRC10） |
| protocol  | 协议名称 |
| network  | 网络名称（暂时BTC协议有用{MainNet：主网，TestNet：测试网，TestNet3：测试网3，SimNet：测试网}） |
| rpc  | rpc配置 |
| user  | rpc用户名（没有则为空） |
| pass  | rpc密码（没有则为空） |
| file  | db文件路径配置 |
| wallet_prefix  | 钱包的存储前缀 |
| hash_prefix  | 交易哈希的存储前缀 |
| block_init  | 初始块（默认读取最新块） |
| block_after_time  | 获取最新块的等待时间 |
| receipt_count  | 交易凭证worker数量 |
| receipt_after_time  | 获取交易信息的等待时间 |
| collection_after_time  | 归集等待时间 |
| collection_count  | 归集发送worker数量 |
| collection_max  | 最大的归集数量（满足多少才归集，为0表示不自动归集） |
| collection_address  | 归集地址 |
| confirms  | 确认数量 |
| recharge_notify_url  | 充值通知回调地址 |
| withdraw_notify_url  | 提现通知回调地址 |
| withdraw_private_key  | 提现的私钥地址 |

> 启动后访问： `http://localhost:10009/swagger/index.html`


# Swagger

> 把 swag cmd 包下载 `go get -u github.com/swaggo/swag/cmd/swag`

> 这时会在 bin 目录下生成一个 `swag.exe` ，把这个执行文件放到 `$GOPATH/bin` 下面

> 执行 `swag init` 注意，一定要和main.go处于同一级目录

> 启动时加上 `-tags "doc"` 才会启动swagger。 这里主要为了正式环境去掉 swagger，这样整体编译的包小一些

> 启动后访问： `http://ip:prot/swagger/index.html`

# 第三方库依赖

> log 日志 `github.com/rs/zerolog`

> 命令行工具 `github.com/urfave/cli`

> 配置文件 `github.com/jinzhu/configor`

# 环境依赖

> go 1.16+

> Redis 3

> MySQL 5.7

# 其它

> `script/Generate MyPOJOs.groovy` 生成数据库Model

# 合约相关

> `solcjs.cmd --version` 查看版本
>
> `solcjs.cmd --abi erc20.sol`
>
> `abigen --abi=erc20_sol_IERC20.abi --pkg=eth --out=erc20.go`

# 吃鸡地址

> 0xDfdf53447cA55820Ec2B3dE9EA707A31579F5c0F
>
> 定制开发请联系：https://t.me/aa333555

# 准备

要实现这些功能首先得摸清楚我们需要完成些什么东西

1. 获取最新区块
2. 获取区块内部的交易记录
3. 通过交易哈希获取交易的完成状态
4. 获取某个地址的余额
5. 创建一个地址
6. 签名并发送luo交易
7. 定义接口如下

```go
type Worker interface {
getNowBlockNum() (uint64, error)
getTransaction(uint64) ([]types.Transaction, uint64, error)
getTransactionReceipt(*types.Transaction) error
getBalance(address string) (*big.Int, error)
createWallet() (*types.Wallet, error)
sendTransaction(string, string, *big.Int) (string, error)
}
```

# 实现

> 创建一个地址后把地址和私钥保存下来

## 进

通过一个无限循环的服务不停的去获取最新块的交易数据，并且把交易数据都一一验证是否完成 ，这里判断数据的接收地址（to）是否属于本服务创建的钱包地址，如果是本服务的创建过的地址则判断为充值成功，**（这时逻辑服务里面需要做交易哈希做幂等）**

## 出

用户发起一笔提出操作，用户发起提出时通过服务配置的私钥来打包并签名luo交易。（私钥转到用户输入的提出地址），这里把提交的luo交易的哈希记录到服务 通过一个无限循环的服务不停的去获取最新块的交易数据，并且把交易数据都一一验证是否完成
，这里判断交易数据的哈希是否存在于服务，如果存在则处理**（这时逻辑服务里面需要做交易哈希做幂等）**

## 归集

通过定期循环服务创建的地址去转账到服务配置的归集地址里面，这里需要注意归集数量的限制，当满足固定的数量时才去归集（减少gas费）

# 一个简单的示例

github地址： [golang 实现加密货币的充值/提现/归集服务](https://github.com/lmxdawn/wallet)

# 特别说明

> 创建钱包的方式可以用 create2 创建合约，这样可以实现不用批量管理私钥，防止私钥丢失或者被盗。

- solidity

```solidity
// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./ERC20.sol";

contract Wallet {
    address internal token = 0xDA0bab807633f07f013f94DD0E6A4F96F8742B53;
    address internal hotWallet = 0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2;

    constructor() {
        // send all tokens from this contract to hotwallet
        IERC20(token).transfer(
            hotWallet,
            IERC20(token).balanceOf(address(this))
        );
        // selfdestruct to receive gas refund and reset nonce to 0
        selfdestruct(payable(hotWallet));
    }
}

contract Fabric {
    function createContract(uint256 salt) public returns (address newAddr){
        // get wallet init_code
        bytes memory bytecode = type(Wallet).creationCode;
        assembly {
            let codeSize := mload(bytecode) // get size of init_bytecode
            newAddr := create2(
                0, // 0 wei
                add(bytecode, 32), // the bytecode itself starts at the second slot. The first slot contains array length
                codeSize, // size of init_code
                salt // salt from function arguments
            )
        }
    }
    function getAddress(uint _salt)
        public
        view
        returns (address)
    {
        bytes memory bytecode = type(Wallet).creationCode;
        bytes32 hash = keccak256(
            abi.encodePacked(bytes1(0xff), address(this), _salt, keccak256(bytecode))
        );

        // NOTE: cast last 20 bytes of hash to address
        return address(uint160(uint(hash)));
    }

    function getBytecode() public pure returns (bytes memory) {
        bytes memory bytecode = type(Wallet).creationCode;

        return bytecode;
    }

    
    function getBytecode1() public pure returns (bytes1) {
        
        return bytes1(0xff);
    }

    
    function getBytecode3(uint256 s) public pure returns (bytes memory) {
        
        return abi.encodePacked(s);
    }
    
    
    function getBytecode2() public pure returns (bytes32) {
        bytes memory bytecode = type(Wallet).creationCode;

        return keccak256(bytecode);
    }
}
```

- go

```go
code := "6080604052600080546001600160a01b031990811673da0bab807633f07f013f94dd0e6a4f96f8742b53179091556001805490911673ab8483f64d9c6d1ecf9b849ae677dd3315835cb217905534801561005857600080fd5b506000546001546040516370a0823160e01b81523060048201526001600160a01b039283169263a9059cbb92169083906370a0823190602401602060405180830381865afa1580156100ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100d29190610150565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af115801561011d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101419190610169565b506001546001600160a01b0316ff5b60006020828403121561016257600080fd5b5051919050565b60006020828403121561017b57600080fd5b8151801515811461018b57600080fd5b939250505056fe"
codeB := common.Hex2Bytes(code)
codeHash := crypto.Keccak256Hash(codeB)
fmt.Println(codeHash)

address := common.HexToAddress("0x7EF2e0048f5bAeDe046f6BF797943daF4ED8CB47")
fmt.Println(address)

fmt.Println(common.LeftPadBytes(big.NewInt(1).Bytes(), 32))
var buffer bytes.Buffer
buffer.Write(common.FromHex("0xff"))
buffer.Write(address.Bytes())
buffer.Write(common.Hex2Bytes("0x30"))
buffer.Write(codeHash.Bytes())

hash := crypto.Keccak256Hash([]byte{0xff}, address.Bytes(), common.LeftPadBytes(big.NewInt(1).Bytes(), 32), codeHash.Bytes())

//salt := common.LeftPadBytes(big.NewInt(1).Bytes(), 32)
//crypto.CreateAddress2(address, salt, codeHash.Bytes())

fmt.Println(common.BytesToAddress(hash[12:]))
```
