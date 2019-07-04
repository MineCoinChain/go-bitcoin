package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	//版本号
	Version uint64
	//交易的根哈希值
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值, 系统提供一个数据，用于计算出一个哈希值
	Bits uint64
	//随机数，挖矿要求的数值
	Nonce uint64
	// 前区块哈希
	PrevHash []byte
	// 哈希, 为了方便，我们将当前区块的哈希放入Block中
	Hash []byte
	//交易数据
	Transaction []*Transaction
}

func NewBlock(transaction []*Transaction, PrevHash []byte) *Block {
	b := Block{
		Version:     0,
		MerkleRoot:  nil,
		TimeStamp:   uint64(time.Now().Unix()),
		Bits:        0,
		Nonce:       0,
		PrevHash:    PrevHash,
		Hash:        nil,
		Transaction: transaction,
	}
	//添加默克尔树
	b.HashTransactionMerkleRoot()
	fmt.Println("merkleroot is", b.MerkleRoot)
	var pow = NewPOW(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce
	return &b
}

func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatal("encode err:", err)
	}
	return buffer.Bytes()
}

func Deserialize(src []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(src))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal("Deserialize error1111:", err)
	}
	return &block
}

//添加简易默克尔树
func (block *Block) HashTransactionMerkleRoot() {
	//遍历所有交易，求出交易哈希值
	var info [][]byte
	for _, tx := range block.Transaction {
		info = append(info, tx.TXID)
	}
	value := bytes.Join(info, []byte{})
	hash := sha256.Sum256(value)
	block.MerkleRoot = hash[:]
}
