package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
	"strings"
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
	Txid  []byte //这个input所引用的output所在的交易id
	Index int64  //这个input所引用的output在交易中的所以

	// ScriptSig string //付款人对当前交易(新交易，而不是引用的交易)的签名
	ScriptSig []byte //对当前交易的签名
	PubKey    []byte //付款人的公钥
}

type TXOutput struct {
	ScriptPubKeyHash []byte  //收款人的公钥哈希
	Value            float64 //转账金额
}

//由于没有办法直接将地址赋值给TXoutput，所以需要提供一个output的方法
func newTXOutput(address string, amount float64) TXOutput {
	output := TXOutput{Value: amount}

	//通过地址获取公钥哈希值
	pubKeyHash := getPubKeyHashFromAddress(address)
	output.ScriptPubKeyHash = pubKeyHash

	return output
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
	input := TXInput{Txid: nil, Index: -1, ScriptSig: nil, PubKey: []byte(data)}

	//创建output
	// output := TXOutput{Value: reward, ScriptPubk: miner}
	output := newTXOutput(miner, reward)

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
	//钱包就是在这里使用的，from=》钱包里面找到对应的wallet-》私钥-》签名
	wm := NewWalletManager()
	if wm == nil {
		fmt.Println("打开钱包失败!")
		return nil
	}

	// 钱包里面找到对应的wallet
	wallet, ok := wm.Wallets[from]
	if !ok {
		fmt.Println("没有找到付款人地址对应的私钥!")
		return nil
	}

	fmt.Println("找到付款人的私钥和公钥，准备创建交易...")

	priKey := wallet.PriKey //私钥签名阶段使用，暂且注释掉
	pubKey := wallet.PubKey
	//我们的所有output都是由公钥哈希锁定的，所以去查找付款人能够使用的output时，也需要提供付款人的公钥哈希值
	pubKeyHash := getPubKeyHashFromPubKey(pubKey)

	// 2. 遍历账本，找到from满足条件utxo集合（3），返回这些utxo包含的总金额(15)

	//包含所有将要使用的utxo集合
	var spentUTXO = make(map[string][]int64)
	//这些使用utxo包含总金额
	var retValue float64

	//遍历账本，找到from能够使用utxo集合,以及这些utxo包含的钱
	// spentUTXO, retValue = bc.findNeedUTXO(from, amount)
	spentUTXO, retValue = bc.findNeedUTXO(pubKeyHash, amount)
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
			input := TXInput{Txid: []byte(txid), Index: i, ScriptSig: nil, PubKey: pubKey}
			inputs = append(inputs, input)
		}
	}

	// 5. 拼接outputs
	// > 创建一个属于to的output
	//创建给收款人的output
	output1 := newTXOutput(to, amount)
	outputs = append(outputs, output1)

	// > 如果总金额大于需要转账的金额，进行找零：给from创建一个output
	if retValue > amount {
		// output2 := TXOutput{from, retValue - amount}
		output2 := newTXOutput(from, retValue-amount)
		outputs = append(outputs, output2)
	}

	timeStamp := time.Now().Unix()

	// 6. 设置哈希，返回
	tx := Transaction{nil, inputs, outputs, uint64(timeStamp)}
	tx.setHash()

	if !bc.signTransaction(&tx, priKey) {
		fmt.Println("交易签名失败")
		return nil
	}
	return &tx
}

//实现具体签名动作（copy，设置为空，签名动作）
//参数1：私钥
//参数2：inputs所引用的output所在交易的集合:
// > key :交易id
// > value：交易本身
func (tx *Transaction) sign(priKey *ecdsa.PrivateKey, prevTxs map[string]*Transaction) bool {
	fmt.Println("具体对交易签名sign...")

	if tx.isCoinbaseTx() {
		fmt.Println("找到挖矿交易，无需签名!")
		return true
	}

	//1. 获取交易copy，pubKey，ScriptPubKey字段置空
	txCopy := tx.trimmedCopy()

	//2. 遍历交易的inputs for, 注意，不要遍历tx本身，而是遍历txCopy
	for i, input := range txCopy.TXInputs {
		fmt.Printf("开始对input[%d]进行签名...\n", i)

		prevTx := prevTxs[string(input.Txid)]
		if prevTx == nil {
			return false
		}

		//input引用的output
		output := prevTx.TXOutputs[input.Index]

		// > 获取引用的output的公钥哈希
		//for range是input是副本，不会影响到变量的结构
		// input.PubKey = output.ScriptPubKeyHash
		txCopy.TXInputs[i].PubKey = output.ScriptPubKeyHash

		// > 对copy交易进行签名，需要得到交易的哈希值
		txCopy.setHash()

		// > 将input的pubKey字段置位nil, 还原数据，防止干扰后面input的签名
		txCopy.TXInputs[i].PubKey = nil

		hashData := txCopy.TXID //我们去签名的具体数据

		//> 开始签名
		r, s, err := ecdsa.Sign(rand.Reader, priKey, hashData)
		if err != nil {
			fmt.Println("签名失败!")
			return false
		}
		signature := append(r.Bytes(), s.Bytes()...)

		// > 将数字签名赋值给原始tx
		tx.TXInputs[i].ScriptSig = signature
	}

	fmt.Println("交易签名成功!")
	return true
}

//trim修剪, 签名和校验时都会使用
func (tx *Transaction) trimmedCopy() *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	//创建一个交易副本，每一个input的pubKey和Sig都设置为空。
	for _, input := range tx.TXInputs {
		input := TXInput{
			Txid:      input.Txid,
			Index:     input.Index,
			ScriptSig: nil,
			PubKey:    nil,
		}
		inputs = append(inputs, input)
	}

	outputs = tx.TXOutputs

	txCopy := Transaction{tx.TXID, inputs, outputs, tx.TimeStamp}
	return &txCopy
}

//具体校验
func (tx *Transaction) verify(prevTxs map[string]*Transaction) bool {
	//1. 获取交易副本txCopy
	txCopy := tx.trimmedCopy()
	//2. 遍历交易，inputs，
	for i, input := range tx.TXInputs {
		prevTx := prevTxs[string(input.Txid)]
		if prevTx == nil {
			return false
		}

		//3. 还原数据（得到引用output的公钥哈希）获取交易的哈希值
		output := prevTx.TXOutputs[input.Index]
		txCopy.TXInputs[i].PubKey = output.ScriptPubKeyHash
		txCopy.setHash()

		//清零环境, 设置为nil
		txCopy.TXInputs[i].PubKey = nil

		//具体还原的签名数据哈希值
		hashData := txCopy.TXID
		//签名
		signature := input.ScriptSig
		//公钥的字节流
		pubKey := input.PubKey

		//开始校验
		var r, s, x, y big.Int
		//r,s 从signature截取出来
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//x, y 从pubkey截取除来，还原为公钥本身
		x.SetBytes(pubKey[:len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey)/2:])
		curve := elliptic.P256()
		pubKeyRaw := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		//进行校验
		res := ecdsa.Verify(&pubKeyRaw, hashData, &r, &s)
		if !res {
			fmt.Println("发现校验失败的input!")
			return false
		}
	}
	//4. 通过tx.ScriptSig, tx.PubKey进行校验
	fmt.Println("交易校验成功!")

	return true
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXID))

	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKeyHash))
	}

	return strings.Join(lines, "\n")
}
