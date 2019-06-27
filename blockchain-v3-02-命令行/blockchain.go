package main

import (
	"errors"
	"fmt"

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

//创建区块，从无到有：这个函数仅执行一次
func CreateBlockChain() error {
	// 1. 区块链不存在，创建
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return err
	}

	//不要db.Close，后续要使用这个句柄
	defer db.Close()

	// 2. 开始创建
	err = db.Update(func(tx *bolt.Tx) error {
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
			//key是区块的哈希值，value是block的字节流
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()) //将block序列化
			//更新最后区块哈希值到数据库
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
		}
		return nil
	})
	return err //nil
}

//获取区块链实例，用于后续操作, 每一次有业务时都会调用
func GetBlockChainInstance() (*BlockChain, error) {
	var lastHash []byte //内存中最后一个区块的哈希值

	//两个功能：
	// 1. 如果区块链不存在，则创建，同时返回blockchain的示例
	db, err := bolt.Open(blockchainDBFile, 0400, nil) //rwx  0100 => 4
	if err != nil {
		return nil, err
	}

	//不要db.Close，后续要使用这个句柄

	// 2. 如果区块链存在，则直接返回blockchain示例
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			return errors.New("bucket不应为nil")
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
func (bc *BlockChain) AddBlock(data string) error {
	lashBlockHash := bc.tail //区块链中最后一个区块的哈希值

	//1. 创建区块
	newBlock := NewBlock(data, lashBlockHash)

	//2. 写入数据库
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("AddBlock时Bucket不应为空")
		}

		//key是新区块的哈希值， value是这个区块的字节流
		bucket.Put(newBlock.Hash, newBlock.Serialize())
		bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)

		//更新bc的tail，这样后续的AddBlock才会基于我们newBlock追加
		bc.tail = newBlock.Hash
		return nil
	})

	return err
}

// +++++++++++++++++++i迭代器,遍历区块链 +++++++++++++++++++++

//定义迭代器
type Iterator struct {
	db          *bolt.DB
	currentHash []byte //不断移动的哈希值，由于访问所有区块
}

//将Iterator绑定到BlockChain
func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}

	return &it
}

//给Iterator绑定一个方法：Next
//1. 返回当前所指向的区块
//2. 向左移动（指向前一个区块）
func (it *Iterator) Next() (block *Block) {

	//读取Bucket当前哈希block
	err := it.db.View(func(tx *bolt.Tx) error {
		//读取bucket
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("Iterator Next时bucket不应为nil")
		}

		blockTmpInfo /*block的字节流*/ := bucket.Get(it.currentHash) //一定要注意，是currentHash
		block = Deserialize(blockTmpInfo)
		it.currentHash = block.PrevHash //游标左移

		return nil
	})
	//哈希游标向左移动

	if err != nil {
		fmt.Printf("iterator next err:", err)
		return nil
	}
	return
}
