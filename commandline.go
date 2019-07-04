package main

import (
	"fmt"
	"log"
)

func (cli *CLI) createBlockChain(address string) {
	if !isValidAddress(address){
		fmt.Println("传入地址无效")
		return
	}
	err := CreateBlockChain(address)
	if err != nil {
		fmt.Println("CreateBlockChain failed:", err)
		return
	}
	fmt.Println("创建区块链成功!")

}

func (cli *CLI) print() {
	bc, err := GetBlockChainInstance()
	if err != nil {
		log.Fatal("err:", err)
	}
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
	if !isValidAddress(address){
		fmt.Println("传入地址无效")
		return
	}
	bc, err := GetBlockChainInstance()
	if err != nil {
		log.Fatal("get block chain instance error:", err)
	}
	defer bc.db.Close()
	PubKeyHash:=GetPubKeyHashFromAddress(address)
	utxos := bc.FindMyUTXO(PubKeyHash)
	total := 0
	for _, txoutput := range utxos {
		total += txoutput.TXOutput.Value
	}
	fmt.Printf("address：%s 的余额为：%d\n", address, total)
}

/*由于暂时没有挖矿竞争机制，因此每次send时指定一名矿工生成一个区块，将一笔交易打包至区块*/
func (cli *CLI) Send(from, to string, amount int, miner, data string) {

	if !isValidAddress(from){
		fmt.Println("传入from地址无效")
		return
	}
	if !isValidAddress(to){
		fmt.Println("传入to地址无效")
		return
	}
	if !isValidAddress(miner){
		fmt.Println("传入miner地址无效")
		return
	}
	bc, err := GetBlockChainInstance()
	if err != nil {
		log.Println("err:", err)
		return
	}
	defer bc.db.Close()
	coinBaseTx := NewCoinbaseTx(miner, data)
	//常见txs数组有效的交易添加进来
	txs := []*Transaction{coinBaseTx}
	//创建普通交易
	tx := NewTransaction(from, to, amount, bc)
	if tx == nil {
		log.Println("这是一笔无效的交易")
	} else {
		log.Println("这是一笔有效的交易")
		txs = append(txs, tx)
	}
	err = bc.AddBlock(txs)
	if err != nil {
		log.Fatal("添加区块失败")
	}
	fmt.Println("添加区块成功，转账成功")
}
func (cli *CLI) createWallet() {
	wm := NewWalletManager()
	address := wm.createWallet()
	if len(address) == 0 {
		log.Println("创建钱包失败")
		return
	}
	fmt.Println("新钱包的地址为：", address)
}

func (cli *CLI) listAddress() {
	wm := NewWalletManager()
	if wm == nil {
		log.Println("NewWalletManager failed")
		return
	}
	addresses := wm.listAddress()
	for _, address :=range addresses{
		fmt.Printf("%s\n", address)
	}
}