package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func uintToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//使用二进制编码
	// Write(w io.Writer, order ByteOrder, data interface{}) error
	err := binary.Write(&buffer, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println("binary.Write err :", err)
		return nil
	}

	return buffer.Bytes()
}

//判断文件是否存在
func isFileExist(filename string) bool {
	// func Stat(name string) (FileInfo, error) {
	_, err := os.Stat(filename)

	//os.IsExist不要使用，不可靠
	if os.IsNotExist(err) {
		return false
	}

	return true
}
