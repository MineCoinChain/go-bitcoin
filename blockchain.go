package main

type BlockChain struct {
	Blocks []*Block
}

//创世语
const genesisInfo string = "This is the first block"

//提供初始化方法
func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock([]byte(genesisInfo), nil)
	return &BlockChain{
		Blocks: []*Block{genesisBlock},
	}
}

//向区块连中添加区块
func (bc *BlockChain) AddBlock(data string) {
	//查找当前最后一个区块的hash
	previousHash := bc.Blocks[len(bc.Blocks)-1].Hash
	var b = NewBlock([]byte(data), []byte(previousHash))
	bc.Blocks = append(bc.Blocks, b)
}
