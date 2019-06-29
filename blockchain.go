package main

import (
	"errors"
	"fmt"
	"github.com/bolt"
)

type BlockChain struct {
	db   *bolt.DB
	tail []byte
}

//创世语
const genesisInfo string = "This is the first block"
const blockchainDBFile = "blockchain.db"
const bucketBlock = "bucketBlock"
const lastBlockHashKey = "lastBlockHashKey"

//提供初始化方法
func CreateBlockChain() error {
	//区块不存在，创建
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte(bucketBlock))
			if err != nil {
				return err
			}
			coinbase := NewCoinbaseTx("中本聪", genesisInfo)
			txs := []*Transaction{coinbase}
			genesisBlock := NewBlock(txs, nil)
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
		}
		return nil
	})
	return err
}

//获取区块链实例，用于后续操作, 每一次有业务时都会调用
func GetBlockChainInstance() (*BlockChain, error) {
	var lastHash []byte
	db, err := bolt.Open(blockchainDBFile, 0400, nil)
	if err != nil {
		return nil, err
	}
	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("bucket should not be nil")
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})
	bc := BlockChain{db, lastHash}
	return &bc, nil
}

//向区块连中添加区块
func (bc *BlockChain) AddBlock(data string) error {
	//lastBlockHash := bc.tail
	//newBlock := NewBlock(data, lastBlockHash)
	//err := bc.db.Update(func(tx *bolt.Tx) error {
	//	bucket := tx.Bucket([]byte(bucketBlock))
	//	if bucket == nil {
	//		return errors.New("Bucket should not be null")
	//	}
	//	bucket.Put(newBlock.Hash, newBlock.Serialize())
	//	bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)
	//	bc.tail = newBlock.Hash
	//	return nil
	//})
	//return err
	return nil
}

//定义迭代器
type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}
	return &it
}
func (it *Iterator) Next() (block *Block) {
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("Iterator Next Block err")
		}
		blockTmpInfo := bucket.Get(it.currentHash)
		block = Deserialize(blockTmpInfo)
		it.currentHash = block.PrevHash
		return nil
	})
	if err != nil {
		fmt.Printf("iterator next err:", err)
		return nil
	}
	return
}

//获取指定账本的金额
func (bc *BlockChain) FindMyUTXO(address string) []TXOutput {
	var utxos []TXOutput
	//查找未花费交易的辅助集合
	var spendUtxos = make(map[string][]int)
	//遍历区块
	it := bc.NewIterator()
	for {
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transaction {
			//遍历output，判断这个output的锁定脚本是否是我们的目标地址
		LABEL:
			for outputIndex, output := range tx.TXOuputs {
				if output.ScriptPubk == address {
					//过滤已经花费的交易
					currentTxId := string(tx.TXID)
					if _, ok := spendUtxos[currentTxId]; ok {
						//判定index是否相等
						currentIds := spendUtxos[currentTxId]
						for _, id := range currentIds {
							if outputIndex == id {
								break LABEL
							}
						}
					}
					//如果有的话添加如utxos中
					utxos = append(utxos, output)
				}
			}
			//遍历input，添加辅助集合：
			for _, input := range tx.TXInputs {
				if input.ScriptSig == address {
					spentKey := string(input.Txid)
					spendUtxos[spentKey] = append(spendUtxos[spentKey], int(input.Index))
				}
			}
		}
		if it.currentHash == nil {
			break
		}
	}

	return utxos
}

