
## 项目说明
本项目实现以太坊秘钥存储管理服务，在只暴露地址的情况下实现数据的签名，交易发送等功能。   
技术选型：gin + mysql
## 启动说明
1. 本项目使用的是go module包管理工具，建议使用go1.12版本。
2. 项目运行之前还要配置`.env`文件，配置方法直接参看项目一级目录下的`.env_example`文件中的内容，主要工作是修改mysql的路径为自己环境的路径。
3. 需要在代码运行环境中安装mysql个数据库，推荐安装5.0或者5.5版本的mysql


## 功能说明
1. 生成以太坊地址并加密存储
2. 传入指定的以太坊秘钥对进行加密存储
3. 分页拉取存储的地址
4. 数据签名
5. 批量发送代币交易，用于资产归集(待续)

## kms接口使用
---
#### 生成秘钥对，返回地址
###### 请求url
- `/kms/new_key`
###### 请求方式
- `GET`
##### 返回示例
```
    {
    "status": 200,
    "data": {
        "address": "0x0C3C35a3455103e1611EF896dBaD1a9B96324c61"
    },
    "msg": "创建以太坊秘钥对成功",
    "error": ""
}
```

---

#### 分页拉取地址
###### 请求url
- `/kms/batch_get_address`
###### 请求方式
- POST
###### 参数
```
{
	"startId": 3,
	"limit": 16
}
```
###### 参数说明
1. `startId`: 起始位置的上一个位置
2. `limit`: 拉取数量
###### 返回示例
```
    {
    "status": 200,
    "data": {
        "addresses": [
            "0x6919a5dB518Ec0EDbd18eaA8783883e08f50Cf66",
            "0x0928580e69A044a670ab2C3C00Ce44d2F33a6611",
            "0x00ee6E720660A03DFF31B58e7376FE650acBd419"
        ],
        "total": 6
    },
    "msg": "批量拉取地址成功",
    "error": ""
}
```
---

#### 调用签名
###### 请求url
- `/kms/sign`
###### 请求方式
- POST
###### 参数
```
{
	"address": "0x6919a5dB518Ec0EDbd18eaA8783883e08f50Cf66",
	"data": "0x1111"
}
```
###### 参数说明
1. `address`: 签名的地址
2. `data`: 需要签名的data
###### 返回示例
```
{
    "status": 200,
    "data": {
        "result": "0x48d4b18339e3a2c755086ca02e39cf3f5dc7502248751954fe4bc27faafb082a3d01e5a8369cbd541d1e64571d99aeb296a2ca916f01dcb8537871e0a0ab771700"
    },
    "msg": "签名数据成功",
    "error": ""
}
```
---

---

#### 传入密钥对进行存储
###### 请求url
- `/kms/save`
###### 请求方式
- POST
###### 参数
```
{
	"address": "0x59375A522876aB96B0ed2953D0D3b92674701Cc2",
	"private": "69F657EAF364969CCFB2531F45D9C9EFAC0A63E359CEA51E5F7D8340784168D2"
}
```
###### 参数说明
1. `address`: 保存的地址
2. `private`: 对应的私钥
###### 返回示例
```
{
    "status": 200,
    "data": null,
    "msg": "秘钥保存成功",
    "error": ""
}
```
---

---

#### get private
###### 请求url
- `kms/get_private`
###### 请求方式
- POST
###### 参数
```
{
	"address": "0x659163470514CEfdce8991c2756Cf0130f84ee31"
}
```
###### 参数说明
1. `address`: 对应地址
###### 返回示例
```
{
    "status": 200,
    "data": {
        "address": "0x659163470514CEfdce8991c2756Cf0130f84ee31",
        "private": "0x688db2b3bdb5fa6cfc75af28f354427c2ce211af509f1b8a7c71fe28d86e00ab"
    },
    "msg": "获取private成功",
    "error": ""
}
```
---

