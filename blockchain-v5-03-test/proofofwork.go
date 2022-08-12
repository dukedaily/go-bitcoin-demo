package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//实现挖矿功能 pow

// 字段：
// ​	区块：block
// ​	目标值：target
// 方法：
// ​	run计算
// ​	功能：找到nonce，从而满足哈希币目标值小

type ProofOfWork struct {
	// ​区块：block
	block *Block
	// ​目标值：target，这个目标值要与生成哈希值比较
	target *big.Int //结构，提供了方法：比较，把哈希值设置为big.Int类型
}

//创建ProofOfWork
//block由用户提供
//target目标值由系统提供
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//难度值先写死，不去推导，后面补充推导方式
	//64位的16进制数: 64 * 4 = 256
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	tmpBigInt := new(big.Int)
	//将我们的难度值赋值给bigint
	tmpBigInt.SetString(targetStr, 16)

	pow.target = tmpBigInt
	return &pow
}

//挖矿函数，不断变化nonce，使得sha256(数据+nonce) < 难度值
//返回：区块哈希，nonce
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//定义随机数
	var nonce uint64
	var hash [32]byte
	fmt.Println("开始挖矿...")

	for {
		fmt.Printf("%x\r", hash[:])
		// 1. 拼接字符串 + nonce
		data := pow.PrepareData(nonce)
		// 2. 哈希值 = sha256(data)
		hash = sha256.Sum256(data)

		//将hash转换为bigInt类型
		tmpInt := new(big.Int)
		tmpInt.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//当前计算的哈希.Cmp(难度值)
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功,hash :%x, nonce :%d\n", hash[:], nonce)
			break
		} else {
			//如果不小于难度值
			nonce++
		}
	} //for

	// 	return 哈希，nonce
	return hash[:], nonce
}

//拼接nonce和block数据
func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {
	b := pow.block

	tmp := [][]byte{
		uintToByte(b.Version), //将uint64转换为[]byte
		b.PrevHash,
		b.MerkleRoot, //所有的交易数据计算得出梅克尔根值， 用于哈希运算
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(nonce),
		// b.Hash, //它不应该参与哈希运算
		// b.Data,
	}
	//使用join方法，将二维切片转为1维切片
	data := bytes.Join(tmp, []byte{})
	return data
}

func (pow *ProofOfWork) IsValid() bool {
	// 	1. 获取区块
	// 2. 拼装数据（block + nonce）
	data := pow.PrepareData(pow.block.Nonce)
	// 3. 计算sha256
	hash := sha256.Sum256(data)
	// 4. 与难度值比较
	tmpInt := new(big.Int)
	tmpInt.SetBytes(hash[:])

	// if tmpInt.Cmp(pow.target) == -1 {
	// 	return true
	// }
	// return false

	//满足条件，返回true
	return tmpInt.Cmp(pow.target) == -1
}
