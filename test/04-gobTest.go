package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Person struct {
	Name string
	Age  uint64
}

func main() {

	//编解码Person结构
	lily := Person{
		"Lily",
		18,
	}

	var buffer bytes.Buffer

	//编码
	//1. 创建编码器
	encoder := gob.NewEncoder(&buffer)
	//2. 编码
	err := encoder.Encode(&lily)
	if err != nil {
		fmt.Printf("encode err:", err)
		return
	}

	data := buffer.Bytes()
	fmt.Printf(" 编码后的数据: %x\n", data)
	//传输.......

	// p1 := Person{}
	var p Person

	//解码
	//1. 创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//2. 解码
	err = decoder.Decode(&p)
	if err != nil {
		fmt.Printf("decode err:", err)
		return
	}

	fmt.Printf("解码数据: %v\n", p)
}

//1. gob go语言内置编解码包
//2. 可以支持变长类型的编解码（通用）
//binary.Write，高效，短板：编解码的数据必须定长string,[]byte
//proto（不可读，高效） > binary（高效） > gob（通用） > json(可读)
