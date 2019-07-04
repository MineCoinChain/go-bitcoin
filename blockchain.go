package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/bolt"
	"log"
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
func CreateBlockChain(address string) error {
	if IsFileExists(blockchainDBFile) {
		fmt.Println("区块链已经存在，请直接操作")
		return nil
	}
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
			coinbase := NewCoinbaseTx(address, genesisInfo)
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
	//判断区块链是否存在
	if !IsFileExists(blockchainDBFile) {
		return nil, errors.New("当前区块链不存在，请先创建")
	}
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
func (bc *BlockChain) AddBlock(txs []*Transaction) error {
	//有效的交易会添加到区块
	tx := []*Transaction{}
	//添加区块前对交易进行校验
	for _, tx1 := range txs {
		if bc.verifyTransaction(tx1){
			log.Println("交易校验成功")
			tx = append(tx,tx1)
		}else{
			log.Println("当前交易校验失败")
		}
	}

	lastBlockHash := bc.tail
	newBlock := NewBlock(tx, lastBlockHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("Bucket should not be null")
		}
		bucket.Put(newBlock.Hash, newBlock.Serialize())
		bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)
		bc.tail = newBlock.Hash
		return nil
	})
	return err
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
func (bc *BlockChain) FindMyUTXO(PubKeyHash []byte) []UXTOInfo {
	var utxoInfos []UXTOInfo
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
				if bytes.Equal(output.ScriptPubkeyHash, PubKeyHash) {
					//fmt.Println("output.ScriptPubkeyHash",output.ScriptPubkeyHash)
					//过滤已经花费的交易
					currentTxId := string(tx.TXID)
					if _, ok := spendUtxos[currentTxId]; ok {
						//判定index是否相等
						currentIds := spendUtxos[currentTxId]
						for _, id := range currentIds {
							if outputIndex == id {
								continue LABEL
							}
						}
					}
					//如果有的话添加如utxos中
					//utxos = append(utxos, output)
					utxoinfo := UXTOInfo{Txid: tx.TXID, Index: int64(outputIndex), TXOutput: output}
					utxoInfos = append(utxoInfos, utxoinfo)
				}
			}
			//查看是否时挖矿交易，如果是则直接跳过
			if tx.IsCoinBase() {
				fmt.Println("挖矿交易，无需遍历集合")
				continue
			}
			//遍历input，添加辅助集合：
			for _, input := range tx.TXInputs {
				if bytes.Equal(GetPubKeyHashFromPubKey(input.PubKey), PubKeyHash) {
					spentKey := string(input.Txid)
					spendUtxos[spentKey] = append(spendUtxos[spentKey], int(input.Index))

				}
			}
		}
		if it.currentHash == nil {
			break
		}
	}

	return utxoInfos
}

func (bc *BlockChain) findNeedUTXO(PubKeyHash []byte, amount int) (map[string][]int64, int) {
	var retMap = make(map[string][]int64)
	var retAmount int
	//遍历账本，查找所有的UTXO
	fmt.Println("**********pubkeyHash", PubKeyHash)
	utxoInfos := bc.FindMyUTXO(PubKeyHash)
	for _, utxoinfo := range utxoInfos {
		retAmount += utxoinfo.Value
		retMap[string(utxoinfo.Txid)] = append(retMap[string(utxoinfo.Txid)], utxoinfo.Index)
		if retAmount >= amount {
			break
		}
	}
	return retMap, retAmount
}

//外部调用的签名函数
func (bc *BlockChain) signTransaction(tx *Transaction, priKey *ecdsa.PrivateKey) bool {
	fmt.Println("开始签名交易")
	//根据传递进来的tx，得到所有需要的前交易prevtxs
	prevTxs := make(map[string]*Transaction)
	//遍历账本找到所有的签名集合
	for _, input := range tx.TXInputs {
		prevTx := bc.FindTransaction(input.Txid)
		if prevTx == nil {
			log.Println("没有找到有效的引用交易")
			return false
		}
		fmt.Println("找到了引用交易")
		prevTxs[string(input.Txid)] = prevTx
	}
	//使用sign对交易进行签名
	return tx.sign(priKey, prevTxs)
}

func (bc *BlockChain) FindTransaction(txid []byte) *Transaction {
	//遍历区块，遍历账本，比较txid与id，如果相同返回交易，反之返回nil
	it := bc.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transaction {
			//如果当前对比的交易id与我们查找的交易id相同，我们就找到了目标交易
			if bytes.Equal(tx.TXID, txid) {
				return tx
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return nil
}

//校验单笔交易
func (bc *BlockChain) verifyTransaction(tx *Transaction) bool {

	fmt.Println("矿工开始校验交易")
	if tx.IsCoinBase(){
		log.Println("发现挖矿交易，无需校验")
		return true
	}
	//根据传递进来的tx，得到所有需要的前交易prevtxs
	prevTxs := make(map[string]*Transaction)
	//遍历账本找到所有的签名集合
	for _, input := range tx.TXInputs {
		prevTx := bc.FindTransaction(input.Txid)
		if prevTx == nil {
			log.Println("没有找到有效的引用交易")
			return false
		}
		fmt.Println("找到了引用交易")
		prevTxs[string(input.Txid)] = prevTx
	}
	return tx.verify(prevTxs)
}
