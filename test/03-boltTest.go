package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const testDb = "test.db"

func main() {
	//打开数据库

	// func Open(path string, mode os.FileMode, options *Options) (*DB, error) {
	// db, err := bolt.Open(testDb, 0600, nil)
	db, err := bolt.Open(testDb, 0600, nil)

	if err != nil {
		fmt.Println(" bolt Open err :", err)
		return
	}

	defer db.Close()

	//创建bucket
	err = db.Update(func(tx *bolt.Tx) error {
		//打开一个bucket
		b1 := tx.Bucket([]byte("bucket1"))

		//没有这个bucket
		if b1 == nil {
			//创建一个bucket
			b1, err = tx.CreateBucket([]byte("bucket1"))
			if err != nil {
				fmt.Printf("tx.CreateBucket err:", err)
				return err
			}

			//写入数据
			b1.Put([]byte("key1"), []byte("hello"))
			b1.Put([]byte("key2"), []byte("world"))

			//读取数据
			v1 := b1.Get([]byte("key1"))
			v2 := b1.Get([]byte("key2"))
			v3 := b1.Get([]byte("key3"))

			fmt.Printf("v1:%s\n", string(v1))
			fmt.Printf("v2:%s\n", string(v2))
			fmt.Printf("v3:%s\n", string(v3))
		}
		return nil
	})

	if err != nil {
		fmt.Printf("db.Update err:", err)
	}

	return
}
