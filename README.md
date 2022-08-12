## 概述

该项目介绍了比特币的基本原理，通过5个版本（da1-day5分支）的迭代，逐步介绍了：区块链、哈希、UTXO、梅克尔根、非对称加密、签名、私钥地址、ECDSA等晦涩概念，从而打下进入web3世界坚实的基础。

- 建议配套阅读《精通比特币》一书，差缺补漏，事半功倍。
- 以太坊学习请移步github：[solidity-expert](https://github.com/dukedaily/solidity-expert)

## 项目资源

- 油管视频：https://www.youtube.com/channel/UCSc6tGnLIFvVMXs-ilDyb4A/videos
- B站视频：https://space.bilibili.com/102710441/channel/seriesdetail?sid=2540263
- 配套教程：见docs目录

## 关于作者

国内第一批区块链布道者；2017年开始专注于区块链教育(btc, eth, fabric)，目前base新加坡一线交易所，专注海外defi，dex，元宇宙等业务方向。

- 微信：Adugii
- Email：[dukemecn@gmail.com](mailto:dukemecn@gmail.com)

## 乐捐Donate

- ETH：0x24DFc2e2e9b73Be1f8cCbdB66642858466e590fC
- ens：dukedu.eth
- BTC：38GqvS398GskBX39kijX8ST1tnjJhPzJ1V

## 扫码入群

<img src="assets/image-20220810134215759.png" alt="image-20220810134215759" width="200" height="230" />


## 安装

```sh
go mod init go-bitcoin
go mod tidy
```

## 编译

```sh
# ./build.sh
go build -o blockchain *.go
```

## 运行

```sh
./blockchain
```

![blockchain](assets/blockchain-0315512.gif)

## 创建钱包

```sh
./blockchain createWallet
```

![createwallet](assets/createwallet.gif)

## 打印钱包

```sh
./blockchain listAddress
```

![listaddress](assets/listaddress.gif)

## 创世块

```sh
./blockchain create 15xGXrzZqrCHjZNcZSQyjDaToPX4agz9R7
```

![getbalance](assets/getbalance.gif)

## 查询余额

```sh
./blockchain getBalance 1Q2DT2JithztxChbLhzEUTShrv78EW3duo
```

![getbalance](assets/getbalance-0317023.gif)

## 转账

```sh
./blockchain send \
	1EiLdWg278u261DNs5Vb2Wyh7opscWvV6G \
	1Q2DT2JithztxChbLhzEUTShrv78EW3duo \
	5 1NkNkQUYXWwrw3ewNw3XSdMjdv5keVK1L3 \
	"send 5 btc"
```

![send](assets/send-0317335.gif)

## 打印交易

```sh
./blockchain printTx
```

![printtx](assets/printtx.gif)
