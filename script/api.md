# Swagger Example API
[toc]
## 1	环境变量

### 默认环境1
| 参数名 | 字段值 |
| ------ | ------ |



## 2	Swagger Example API

##### 说明
> 



##### 联系方式
- **联系人：**API Support
- **邮箱：**support@swagger.io
- **网址：**/http://www.swagger.io/support/

##### 文档版本
```
1.0
```


## 3	归集

> POST  /api/collection
### 请求体(Request Body)
| 参数名称 | 数据类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| address|string||false|地址|
| coinName|string||false|币种名称|
| max|string||false|最大归集数量（满足当前值才会归集）|
| protocol|string||false|协议|
### 响应体
● 200 响应数据格式：JSON
| 参数名称 | 类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| data|object||false||
|⇥ balance|string||false|实际归集的数量|
| server.Response|object||false||
|⇥ code|int32||false|错误code码|
|⇥ data|object||false|成功时返回的对象|
|⇥ message|string||false|错误信息|

##### 接口描述
> 




## 4	创建钱包地址

> POST  /api/createWallet
### 请求体(Request Body)
| 参数名称 | 数据类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| coinName|string||false|币种名称|
| protocol|string||false|协议|
### 响应体
● 200 响应数据格式：JSON
| 参数名称 | 类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| data|object||false||
|⇥ address|string||false|生成的钱包地址|
| server.Response|object||false||
|⇥ code|int32||false|错误code码|
|⇥ data|object||false|成功时返回的对象|
|⇥ message|string||false|错误信息|

##### 接口描述
> 




## 5	删除钱包地址

> POST  /api/delWallet
### 请求体(Request Body)
| 参数名称 | 数据类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| address|string||false|地址|
| coinName|string||false|币种名称|
| protocol|string||false|协议|
### 响应体
● 200 响应数据格式：JSON
| 参数名称 | 类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| code|int32||false|错误code码|
| data|object||false|成功时返回的对象|
| message|string||false|错误信息|

##### 接口描述
> 




## 6	获取交易结果

> GET  /api/getTransactionReceipt
### 请求体(Request Body)
| 参数名称 | 数据类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| coinName|string||false|币种名称|
| hash|string||false|交易哈希|
| protocol|string||false|协议|
### 响应体
● 200 响应数据格式：JSON
| 参数名称 | 类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| data|object||false||
|⇥ status|int32||false|交易状态（0：未成功，1：已成功）|
| server.Response|object||false||
|⇥ code|int32||false|错误code码|
|⇥ data|object||false|成功时返回的对象|
|⇥ message|string||false|错误信息|

##### 接口描述
> 




## 7	提现

> POST  /api/withdraw
### 请求体(Request Body)
| 参数名称 | 数据类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| address|string||false|提现地址|
| coinName|string||false|币种名称|
| orderId|string||false|订单号|
| protocol|string||false|协议|
| value|int32||false|金额|
### 响应体
● 200 响应数据格式：JSON
| 参数名称 | 类型 | 默认值 | 不为空 | 描述 |
| ------ | ------ | ------ | ------ | ------ |
| data|object||false||
|⇥ hash|string||false|生成的交易hash|
| server.Response|object||false||
|⇥ code|int32||false|错误code码|
|⇥ data|object||false|成功时返回的对象|
|⇥ message|string||false|错误信息|

##### 接口描述
> 



