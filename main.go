package main


//func main() {
//	err := CreateBlockChain()
//	if err != nil {
//		log.Fatal(err)
//	}
//	bc, err := GetBlockChainInstance()
//	defer bc.db.Close()
//	if err != nil {
//		log.Fatal("GetBlockChainInstance err:", err)
//	}
//
//	err = bc.AddBlock("26号btc暴涨20%")
//	if err != nil {
//		log.Fatal("Add Block error:", err)
//	}
//	err = bc.AddBlock("27号btc暴涨20%")
//	if err != nil {
//		log.Fatal("Add Block error:", err)
//	}
//	it := bc.NewIterator()
//	for {
//		block := it.Next()
//		fmt.Println("*****************************************************")
//		//fmt.Printf("当前区块高度: %d\n", i)
//		fmt.Printf("Version :  %d\n", block.Version)
//		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
//		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
//		fmt.Printf("Bits : %d\n", block.Bits)
//		fmt.Printf("Nonce : %d\n", block.Nonce)
//		fmt.Printf("Hash : %x\n", block.Hash)
//		fmt.Printf("PrevHash : %x\n", block.PrevHash)
//		fmt.Printf("Data : %s\n", block.Data)
//		pow := NewPOW(block)
//		fmt.Println("区块合法性验证:", pow.IsValid())
//		if block.PrevHash == nil {
//			log.Println("the block ranges over")
//			break
//		}
//	}
//
//}
func main(){
	cli:=new(CLI)
	cli.Run()
}