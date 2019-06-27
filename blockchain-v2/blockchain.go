package main

//定义区块链结构(使用数组模拟区块链)
type BlockChain struct {
	Blocks []*Block //区块链
}

//创世语
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

//提供一个创建区块链的方法
func NewBlockChain() *BlockChain {
	//创建BlockChain，同时添加一个创世块
	genesisBlock := NewBlock(genesisInfo, nil)

	bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

	return &bc
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
