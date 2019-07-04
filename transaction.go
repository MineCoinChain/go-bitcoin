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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
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
	ScriptSig []byte //付款人的私钥签名
	PubKey    []byte //付款人的公钥
}

//定义输出结构
type TXOutput struct {
	ScriptPubkeyHash []byte //输出地址的公钥，所定脚本使用
	Value            int
}

//封装output，使其包含output详情
type UXTOInfo struct {
	Txid  []byte
	Index int64
	TXOutput
}

//计算Transaction的哈希值获取交易TXID
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
func NewCoinbaseTx(miner string, data string) *Transaction {
	input := TXInput{
		Txid:  nil,
		Index: -1,
		//挖矿交易不需要签名，因此挖矿字段可以书写任意值
		ScriptSig: nil,
		PubKey:    []byte(data),
	}
	output := NewTXOutput(miner, int64(reward))
	timeStamp := time.Now().Unix()

	tx := Transaction{
		TXID:      nil,
		TXInputs:  []TXInput{input},
		TXOuputs:  []TXOutput{output},
		TimeStamp: uint64(timeStamp),
	}
	tx.setHash()
	return &tx
}

//创建普通交易
func NewTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	//调用钱包，找到付款人的公钥哈希
	wm := NewWalletManager()
	if wm == nil {
		log.Println("打开钱包失败")
		return nil
	}
	//钱包里面找到对应的wallet
	wallet, ok := wm.Wallets[from]
	if !ok {
		log.Println("没有找到付款人地址对应的私钥")
		return nil
	}
	fmt.Println("找到付款人的私钥和公钥")
	priKey := wallet.PriKey                       //签名使用的私钥
	pubKey := wallet.PubKey                       //查找未花费交易输出使用到的公钥
	pubKeyHash := GetPubKeyHashFromPubKey(pubKey) //得到公钥对应的公钥哈希
	var spentUTXO = make(map[string][]int64)
	var retValue int
	//遍历账本，查找from能够使用的utxo集合，以及这些UTXO的余额
	spentUTXO, retValue = bc.findNeedUTXO(pubKeyHash, amount)
	if retValue < amount {
		fmt.Println("可支付的金额不足，创建交易失败")
		return nil
	}
	var inputs []TXInput
	var output []TXOutput
	//拼接inputs
	for txid, indexArray := range spentUTXO {
		for _, index := range indexArray {
			input := TXInput{
				Txid:      []byte(txid),
				Index:     int(index),
				ScriptSig: nil,
				PubKey:    pubKey,
			}
			inputs = append(inputs, input)
		}

	}
	//拼接outputs
	output1 := NewTXOutput(to, int64(amount))
	output = append(output, output1)
	if retValue > amount {
		output2 := NewTXOutput(from, int64(retValue-amount))
		output = append(output, output2)
	}
	timeStamp := time.Now().Unix()
	tx := Transaction{
		TXID:      nil,
		TXInputs:  inputs,
		TXOuputs:  output,
		TimeStamp: uint64(timeStamp),
	}
	tx.setHash()
	if !bc.signTransaction(&tx, priKey) {
		log.Println("交易签名失败")
		return nil
	}
	return &tx
}

func (tx Transaction) IsCoinBase() bool {
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].Txid == nil && inputs[0].Index == -1 {
		return true
	}
	return false
}

//通过地址获取公钥哈希
func NewTXOutput(address string, amount int64) TXOutput {
	output := TXOutput{
		Value: int(amount),
	}
	output.ScriptPubkeyHash = GetPubKeyHashFromAddress(address)
	return output
}

//实现具体的签名动作（copy 设置为空 签名动作）
//参数1：私钥
//参数2：inputs所引用的output所在交易的集合
//>key :交易id
//>value:交易本身
func (tx *Transaction) sign(priKey *ecdsa.PrivateKey, prevTxs map[string]*Transaction) bool {
	log.Println("具体对交易的签名 sign")
	//挖矿交易不需要签名
	if tx.IsCoinBase() {
		log.Println("找到挖矿交易，不需要进行签名")
		return true
	}
	//获取交易的copy，copy时需要将交易置空
	txCopy := tx.trimmedCopy()
	//遍历交易的inputs
	for i, input := range txCopy.TXInputs {
		fmt.Printf("开始对input%d进行签名......\n", i)
		prevTx := prevTxs[string(input.Txid)]
		if prevTx == nil {
			return false
		}
		output := prevTx.TXOuputs[input.Index]
		//>>获取引用的output的公钥哈希
		txCopy.TXInputs[i].PubKey = output.ScriptPubkeyHash
		//>>对copy的交易进行签名
		txCopy.setHash()
		hashData := txCopy.TXID //我们签名的具体数据
		//>>将input的pubKey字段置为nil,还原数据，防止干扰后面input的签名
		txCopy.TXInputs[i].PubKey = nil
		//>>将数字签名赋值给tx
		r, s, err := ecdsa.Sign(rand.Reader, priKey, hashData)
		if err != nil {
			log.Println("私钥签名失败")
			return false
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXInputs[i].ScriptSig = signature

	}

	log.Println("交易签名成功")
	return true
}

func (tx *Transaction) trimmedCopy() *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		input := TXInput{
			Txid:      input.Txid,
			Index:     input.Index,
			ScriptSig: nil,
			PubKey:    nil,
		}
		inputs = append(inputs, input)
	}
	outputs = tx.TXOuputs
	txCopy := Transaction{
		TXID:      tx.TXID,
		TXInputs:  inputs,
		TXOuputs:  outputs,
		TimeStamp: tx.TimeStamp,
	}
	return &txCopy
}

//具体校验
func (tx *Transaction) verify(prevTxs map[string]*Transaction) bool {
	txCopy := tx.trimmedCopy()
	for i, input := range tx.TXInputs {
		prevTx := prevTxs[string(input.Txid)]
		if prevTx == nil {
			return false
		}
		output := prevTx.TXOuputs[input.Index]
		txCopy.TXInputs[i].PubKey = output.ScriptPubkeyHash
		txCopy.setHash()
		//清理环境，设置为nil
		txCopy.TXInputs[i].PubKey = nil
		//获取公钥验证需要的数据
		hashData := txCopy.TXID
		signature := input.ScriptSig
		pubKey := input.PubKey
		//开始校验
		var r, s, x, y big.Int
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		x.SetBytes(pubKey[:len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey)/2:])
		curve := elliptic.P256()
		pubKeyRaw := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if res := ecdsa.Verify(&pubKeyRaw, hashData, &r, &s); !res {
			log.Println("发现校验失败的input！")
			return false
		}
		log.Println("公钥校验成功")

	}
	return true
}
