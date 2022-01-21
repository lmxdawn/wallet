# wallet

> 虚拟币钱包服务，转账/提现/充值/归集
> 
> 计划支持：比特币、以太坊（ERC20）、波场（TRC20），币安（）
> 
> 完全实现与业务服务隔离，使用http服务相互调用

# 接口

`script/api.md`

# 下载-打包

```shell
# 拉取代码
$ git clone https://github.com/lmxdawn/wallet.git
$ cd wallet

# 打包 (-tags "doc") 可选，加上可以运行swagger
$ go build [-tags "doc"]

# 运行
$ wallet -c config/config-example.yml

```
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