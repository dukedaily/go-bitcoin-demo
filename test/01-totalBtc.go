package main

import "fmt"

func main() {
	total := 0.0     //总量
	interval := 21.0 //区块间隔，万单位

	reward := 50.0 //奖励数量，最初50个

	for reward != 0 {
		//累加挖矿值
		amount := reward * interval

		total += amount //总量

		reward *= 0.5 //奖励减半
	}

	fmt.Println("total : ", total)
	fmt.Printf("total : %f\n", total)
}
