package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

//定义交易结构
type Transaction struct {
	TXID      []byte     //交易id
	TXInputs  []TXInput  //可以有多个输入
	TXOutputs []TXOutput //可以有多个输出
	TimeStamp uint64     //创建交易的时间
}

type TXInput struct {
	Txid      []byte //这个input所引用的output所在的交易id
	Index     int64  //这个input所引用的output在交易中的所以
	ScriptSig string //付款人对当前交易(新交易，而不是引用的交易)的签名
}

type TXOutput struct {
	ScriptPubk string  //收款人的公钥哈希，先理解为地址
	Value      float64 //转账金额
}

// # 获取交易ID
// 对交易做哈希处理
func (tx *Transaction) setHash() error {
	//对tx做gob编码得到字节流，做sha256，赋值给TXID
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		fmt.Println("encode err:", err)
		return err
	}

	hash := sha256.Sum256(buffer.Bytes())

	//我们使用tx字节流的哈希值作为交易id
	tx.TXID = hash[:]
	return nil
}

//挖矿奖励
var reward = 12.5

// # 创建挖矿交易
func NewCoinbaseTx(miner /*挖矿人*/ string, data string) *Transaction {
	//特点：没有输入，只有一个输出，得到挖矿奖励
	//挖矿交易需要能够识别出来，没有input，所以不需要签名，
	//挖矿交易不需要签名，所以这个签名字段可以书写任意值，只有矿工有权利写
	//中本聪：写的创世语
	//现在都是由矿池来写，写自己矿池的名字
	input := TXInput{Txid: nil, Index: -1, ScriptSig: data}
	output := TXOutput{Value: reward, ScriptPubk: miner}
	timeStamp := time.Now().Unix()

	tx := Transaction{
		TXID:      nil,
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
		TimeStamp: uint64(timeStamp),
	}

	tx.setHash()
	return &tx
}

//判断一笔交易是否为挖矿交易
func (tx *Transaction) isCoinbaseTx() bool {
	inputs := tx.TXInputs
	//input个数为1，id为nil，索引为-1
	if len(inputs) == 1 && inputs[0].Txid == nil && inputs[0].Index == -1 {
		return true
	}
	return false
}

//创建普通交易
// 1. from/*付款人*/,to/*收款人*/,amount输入参数/*金额*/
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	// 2. 遍历账本，找到from满足条件utxo集合（3），返回这些utxo包含的总金额(15)

	//包含所有将要使用的utxo集合
	var spentUTXO = make(map[string][]int64)
	//这些使用utxo包含总金额
	var retValue float64

	//TODO
	//遍历账本，找到from能够使用utxo集合,以及这些utxo包含的钱
	spentUTXO, retValue = bc.findNeedUTXO(from, amount)
	// map[0x222] = []int{0}
	// map[0x333] = []int{0,1}

	// 3. 如果金额不足，创建交易失败
	if retValue < amount {
		fmt.Println("金额不足，创建交易失败!")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput

	// 4. 拼接inputs
	// > 遍历utxo集合，每一个output都要转换为一个input(3)
	for txid, indexArray := range spentUTXO {
		//遍历下标, 注意value才是我们消耗的output的下标
		for _, i := range indexArray {
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)
		}
	}

	// 5. 拼接outputs
	// > 创建一个属于to的output
	output1 := TXOutput{to, amount}
	outputs = append(outputs, output1)

	// > 如果总金额大于需要转账的金额，进行找零：给from创建一个output
	if retValue > amount {
		output2 := TXOutput{from, retValue - amount}
		outputs = append(outputs, output2)
	}

	timeStamp := time.Now().Unix()

	// 6. 设置哈希，返回
	tx := Transaction{nil, inputs, outputs, uint64(timeStamp)}
	tx.setHash()
	return &tx
}
