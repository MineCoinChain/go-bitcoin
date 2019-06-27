package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("26号btc暴涨20%")
	bc.AddBlock("27号btc暴涨20%")
	for i, block := range bc.Blocks {
		fmt.Println("*****************************************************")
		fmt.Printf("当前区块高度: %d\n", i)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("PrevHash : %x\n", block.PrevHash)
		fmt.Printf("Data : %s\n", block.Data)
		pow:=NewPOW(block)
		fmt.Println("区块合法性验证:",pow.IsValid())
	}

}
