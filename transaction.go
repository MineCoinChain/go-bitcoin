/*
 * @modifiy by:Mine&Coin&Chain
 * @Filename:main
 * @Description:添加交易结构
 * @Date:2019/6/29 12:18
 * @Version:v1.0
*/
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)
var reward = 1250000000
//定义交易结构
type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOuputs  []TXOutput
	TimeStamp uint64
}

//定义输入结构
type TXInput struct {
	Txid      []byte
	Index     int
	ScriptSig string
}

//定义输出结构
type TXOutput struct {
	ScriptPubk string
	Value      int
}

//获取交易ID
func (tx *Transaction) setHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Fatal("encode err:", err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

//创建挖矿交易
func NewCoinbaseTx(miner string,data string) *Transaction{
	input := TXInput{
		Txid:nil,
		Index:-1,
		//挖矿交易不需要签名，因此挖矿字段可以书写任意值
		ScriptSig:data,
	}
	output:=TXOutput{
		Value:reward,
		ScriptPubk:miner,
	}
	timeStamp:=time.Now().Unix()

	tx:=Transaction{
		TXID:nil,
		TXInputs:[]TXInput{input},
		TXOuputs:[]TXOutput{output},
		TimeStamp:uint64(timeStamp),
	}
	tx.setHash()
	return &tx
}

