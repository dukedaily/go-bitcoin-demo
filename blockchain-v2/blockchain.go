package main

import (
	"github.com/boltdb/bolt"
)

//定义区块链结构(使用数组模拟区块链)
type BlockChain struct {
	db   *bolt.DB //用于存储数据
	tail []byte   //最后一个区块的哈希值
}

//创世语
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const blockchainDBFile = "blockchain.db"
const bucketBlock = "bucketBlock"           //装block的桶
const lastBlockHashKey = "lastBlockHashKey" //用于访问bolt数据库，得到最好一个区块的哈希值

//提供一个创建区块链的方法
func NewBlockChain() (*BlockChain, error) {
	var lastHash []byte //内存中最后一个区块的哈希值

	//两个功能：
	// 1. 如果区块链不存在，则创建，同时返回blockchain的示例
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	//不要db.Close，后续要使用这个句柄

	// 2. 如果区块链存在，则直接返回blockchain示例
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			//创建bucket
			bucket, err := tx.CreateBucket([]byte(bucketBlock))
			if err != nil {
				return err
			}
			//写入创世块
			//创建BlockChain，同时添加一个创世块
			genesisBlock := NewBlock(genesisInfo, nil)
			//key是区块的哈希值，value是block的字节流//TODO
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//更新最后区块哈希值到数据库
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

			//更新内存中的最后一个区块哈希值, 后续操作就可以基于这个哈希增加区块
			lastHash = genesisBlock.Hash
		} else {
			//直接读取特定的key，得到最后一个区块的哈希值
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	//5. 拼成BlockChain然后返回
	bc := BlockChain{db, lastHash}
	return &bc, nil
}

//提供一个向区块链中添加区块的方法
//参数：数据，不需要提供前区块的哈希值，因为bc可以通过自己的下标拿到
func (bc *BlockChain) AddBlock(data string) {
	//通过下标，得到最后一个区块
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	//最后一个区块哈希值是新区块的前哈希
	prevHash := lastBlock.Hash

	//创建block
	newBlcok := NewBlock(data, prevHash)

	//添加bc中
	bc.Blocks = append(bc.Blocks, newBlcok)
}
