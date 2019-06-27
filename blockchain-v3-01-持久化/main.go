package main

import (
	"fmt"
)

//打印区块链
func main() {
	// bc := NewBlockChain()

	err := CreateBlockChain()
	fmt.Println("err:", err)

	//获取区块链实例
	bc, err := GetBlockChainInstance()
	defer bc.db.Close()

	if err != nil {
		fmt.Println("GetBlockChainInstance, err :", err)
		return
	}

	err = bc.AddBlock("hello world!!!!!")
	if err != nil {
		fmt.Println("AddBlock, err :", err)
		return
	}

	err = bc.AddBlock("hello itast!!!!!")
	if err != nil {
		fmt.Println("AddBlock, err :", err)
		return
	}

	//调用迭代器，输出blockChain
	it := bc.NewIterator()
	for {
		//调用Next方法，获取区块，游标左移
		block := it.Next()

		fmt.Printf("\n++++++++++++++++++++++\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevHash : %x\n", block.PrevHash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		//退出条件
		if block.PrevHash == nil {
			fmt.Println("区块链遍历结束!")
			break
		}
	}
}
