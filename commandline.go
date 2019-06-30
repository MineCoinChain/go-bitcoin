package main

import (
	"fmt"
	"log"
)

func (cli *CLI) addBlock(data string) {
	//fmt.Println("添加区块被调用!")
	//bc, _ := GetBlockChainInstance()
	//err := bc.AddBlock(data)
	//if err != nil {
	//	fmt.Println("AddBlock failed:", err)
	//	return
	//}
	fmt.Println("添加区块成功!")
}

func (cli *CLI) createBlockChain() {
	err := CreateBlockChain()
	if err != nil {
		fmt.Println("CreateBlockChain failed:", err)
		return
	}
	fmt.Println("创建区块链成功!")

}

func (cli *CLI) print() {
	bc, _ := GetBlockChainInstance()
	//调用迭代器，输出blockChain
	it := bc.NewIterator()
	for {
		//调用Next方法，获取区块，游标左移
		block := it.Next()

		fmt.Printf("\n++++++++++++++++++++++\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevHash : %x\n", block.PrevHash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Transaction[0].TXInputs[0].ScriptSig)

		pow := NewPOW(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		//退出条件
		if block.PrevHash == nil {
			fmt.Println("区块链遍历结束!")
			break
		}
	}

}
func (cli *CLI) GetBalance(address string) {
	bc, err := GetBlockChainInstance()
	if err != nil {
		log.Fatal("get block chain instance error:", err)
	}
	utxos := bc.FindMyUTXO(address)
	total := 0
	for _, txoutput := range utxos {
		total += txoutput.TXOutput.Value
	}
	fmt.Printf("address：%s 的余额为：%d\n", address, total)
}

/*由于暂时没有挖矿竞争机制，因此每次send时指定一名矿工生成一个区块，将交易打包至区块*/
func (cli *CLI) Send(from, to string, amount int, miner, data string) {
	fmt.Println("from:",from)
	fmt.Println("to:",to)
	fmt.Println("amount:",amount)
	fmt.Println("miner:",miner)
	fmt.Println("data:",miner)
}
