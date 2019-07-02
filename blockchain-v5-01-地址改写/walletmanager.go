package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"sort"
)

// 负责对外，管理生成的钱包（公钥私钥）
//私钥1->公钥-》地址1
//私钥2->公钥-》地址2
//私钥3->公钥-》地址3
//私钥4->公钥-》地址4
type WalletManager struct {
	//定义一个map来管理所有的钱包
	//key:地址
	//value:wallet结构(公钥，私钥)
	Wallets map[string]*wallet
}

//创建walletManager结构
func NewWalletManager() *WalletManager {
	//创建一个, Wallets map[string]*wallet
	var wm WalletManager

	//分配空间，一定要分配，否则没有空间
	wm.Wallets = make(map[string]*wallet)

	//从本地加载已经创建的钱包,写入Wallets结构
	if !wm.loadFile() {
		return nil
	}

	//返回 walletManager
	return &wm
}

func (wm *WalletManager) createWallet() string {
	// 创建秘钥对
	w := newWalletKeyPair()
	if w == nil {
		fmt.Println("newWalletKeyPair 失败!")
		return ""
	}

	// 获取地址
	address := w.getAddress()

	//把地址和wallet写入map中: Wallets map[string]*wallet
	wm.Wallets[address] = w //<<<--------- 重要

	// 将秘钥对写入磁盘
	if !wm.saveFile() {
		return ""
	}

	// 返回给cli新地址
	return address

}

const walletFile = "wallet.dat"

func (wm *WalletManager) saveFile() bool {
	//data ????
	//使用gob对wm进行编码
	var buffer bytes.Buffer

	//未注册接口函数
	// encoder.Encode err: gob: type not registered for interface: elliptic.p256Curve
	//注册一下接口函数，这样gob才能够正确的编码
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wm)

	if err != nil {
		fmt.Println("encoder.Encode err:", err)
		return false
	}

	//将walletManager写入磁盘
	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println("ioutil.WriteFile err:", err)
		return false
	}
	return true
}

//读取wallet.dat文件，加载wm中
func (wm *WalletManager) loadFile() bool {
	//判断文件是否存在
	if !isFileExist(walletFile) {
		fmt.Println("文件不存在,无需加载!")
		return true
	}

	//读取文件
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		fmt.Println("ioutil.ReadFile err:", err)
		return false
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))

	//解密赋值，赋值给wm:====>map
	err = decoder.Decode(wm)
	if err != nil {
		fmt.Println("decoder.Decode err:", err)
		return false
	}
	return true
}

func (wm *WalletManager) listAddresses() []string {
	var addresses []string
	for address := range wm.Wallets {
		addresses = append(addresses, address)
	}

	//排序, 升序
	sort.Strings(addresses)

	return addresses
}
