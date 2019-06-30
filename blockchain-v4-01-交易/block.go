package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

//定义区块结构
// 第一阶段: 先实现基础字段：前区块哈希，哈希，数据
// 第二阶段: 补充字段：Version，时间戳，难度值等
type Block struct {
	//版本号
	Version uint64

	// 前区块哈希
	PrevHash []byte

	//交易的根哈希值
	MerkleRoot []byte

	//时间戳
	TimeStamp uint64

	//难度值, 系统提供一个数据，用于计算出一个哈希值
	Bits uint64

	//随机数，挖矿要求的数值
	Nonce uint64

	// 哈希, 为了方便，我们将当前区块的哈希放入Block中
	Hash []byte

	//数据
	// Data []byte
	//一个区块可以有很多交易(交易的集合)
	Transactions []*Transaction
}

//创建一个区块（提供一个方法）
//输入：数据，前区块的哈希值
//输出：区块
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: nil, //随意写的
		TimeStamp:  uint64(time.Now().Unix()),

		Bits:  0, //随意写的, 这个值是uint64，它并不是难度值，而是可以推算出难度值
		Nonce: 0, //随意写的
		Hash:  nil,
		// Data:  []byte(data),
		Transactions: txs,
	}

	//计算哈希值
	// b.setHash()

	//将POW集成到Block中
	pow := NewProofOfWork(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce

	return &b
}

//绑定Serialize方法， gob编码
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer

	//编码器
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(b)
	if err != nil {
		fmt.Printf("Encode err:", err)
		return nil
	}

	return buffer.Bytes()
}

//反序列化，输入[]byte，返回block
func Deserialize(src []byte) *Block {
	var block Block

	//解码器
	decoder := gob.NewDecoder(bytes.NewReader(src))
	//解码
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Printf("decode err:", err)
		return nil
	}

	return &block
}
