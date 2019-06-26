package main

import "math/big"

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
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	tmpBigInt := new(big.Int)
	//将我们的难度值赋值给bigint
	tmpBigInt.SetString(targetStr, 16)

	pow.target = tmpBigInt
	return &pow
}
