package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"

	"golang.org/x/crypto/ripemd160"
)

// - 结构定义
type wallet struct {
	//私钥
	PriKey *ecdsa.PrivateKey

	//公钥原型定义
	// type PublicKey struct {
	// 	elliptic.Curve
	// 	X, Y *big.Int
	// }

	// 公钥, X,Y类型一致，长度一致，我们将X和Y拼接成字节流，赋值给pubKey字段，用于传输
	// 验证时，将X，Y截取出来（类似r,s),再创建一条曲线，就可以还原公钥，进一步进行校验
	PubKey []byte
}

// - ==创建秘钥对==
func newWalletKeyPair() *wallet {
	curve := elliptic.P256()
	//创建私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("ecdsa.GenerateKey err:", err)
		return nil
	}

	//获取公钥
	pubKeyRaw := priKey.PublicKey

	//将公钥X，Y拼接到一起
	pubKey := append(pubKeyRaw.X.Bytes(), pubKeyRaw.Y.Bytes()...)

	//创建wallet结构返回
	wallet := wallet{priKey, pubKey}
	return &wallet
}

// - ==根据私钥生成地址==
func (w *wallet) getAddress() string {
	//公钥
	pubKey := w.PubKey
	hash1 := sha256.Sum256(pubKey)
	//hash160处理
	hasher := ripemd160.New()
	hasher.Write(hash1[:])

	//公钥哈希，锁定output时就是使用这值
	pubKeyHash := hasher.Sum(nil)

	//拼接version和公钥哈希，得到21字节的数据
	payload := append([]byte{byte(0x00)}, pubKeyHash...)

	//生成4字节的校验码
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节checksum
	checksum := second[0:4]

	//25字节数据
	payload = append(payload, checksum...)
	address := base58.Encode(payload)
	return address
}
