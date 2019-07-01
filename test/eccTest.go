package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func main() {

	//go语言只提供了签名校验，未提供加解密

	//选择一条曲线
	curve := elliptic.P256()

	// func GenerateKey(c elliptic.Curve, rand io.Reader) (*PrivateKey, error) {
	//创建私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("ecdsa.GenerateKey:", err)
		return
	}

	//私钥得到公钥
	pubKey := priKey.PublicKey

	data := "hello world"

	hash := sha256.Sum256([]byte(data))

	//私钥签名
	// func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error) {
	r, s, err := ecdsa.Sign(rand.Reader, priKey, hash[:])
	if err != nil {
		fmt.Println("ecdsa.Sign:", err)
		return
	}

	//将r，s的字节流拼接起来，得到签名，进行传输
	fmt.Printf("r len: %d, r bytes: %x\n", len(r.Bytes()), r.Bytes())
	fmt.Printf("s len: %d, s bytes: %x\n", len(s.Bytes()), s.Bytes())

	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("signature:%x\n", signature)
	//传输.....

	var r1, s1 big.Int
	//在对端，将r，s从字节流中截取出来
	r1.SetBytes(signature[:len(signature)/2]) //截取前32字节作为r1
	s1.SetBytes(signature[len(signature)/2:]) //截取后32字节作为s1

	// hash1 := sha256.Sum256([]byte(data + "hello"))

	//公钥验证
	// Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool
	// res := ecdsa.Verify(&pubKey, hash1[:], r, s)
	res := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)
	fmt.Println("res:", res)
}

//签名：原文，私钥
//校验：原文，公钥，签名
