# wallet

> 虚拟币钱包服务，转账/提现/充值/归集
>
>
> 完全实现与业务服务隔离，使用http服务相互调用

# sdk列表

> go示例 [https://github.com/lmxdawn/wallet-v2-sdk-go](https://github.com/lmxdawn/wallet-v2-sdk-go)
> 
> 客户是PHP后端，特定增加一个php示例 [https://github.com/lmxdawn/wallet-v2-sdk-php](https://github.com/lmxdawn/wallet-v2-sdk-php)

# 目前有 Pro 版本，可以在 Issues 联系我

> 有 Pro 版本，更高级的钱包服务，可以wx 联系 appeth，添加时备注GitHub

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

