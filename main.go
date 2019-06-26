package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("26号btc暴涨20%")
	bc.AddBlock("27号btc暴涨20%")
	for i, block := range bc.Blocks {
		fmt.Printf("当前区块高度: %d\n", i)
		fmt.Printf("PrevHash : %x\n", block.PrevHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
	}
}
