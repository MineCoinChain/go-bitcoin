/*
 * @modifiy by:Mine&Coin&Chain
 * @Filename:main
 * @Description:添加wallet钱包
 * @Date:2019/6/30 23:50
 * @Version:v1.0
*/
package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type wallet struct {
	//私钥
	PriKey *ecdsa.PrivateKey
	//公钥原型定义
	// type PublicKey struct {
	// 	elliptic.Curve
	// 	X, Y *big.Int
	// }
	//公钥，X，Y类型一致，长度一致，我们X和Y拼接成字节流，赋值给pubKey字段，用于传输
	//验证时，将X和Y截取出来，再创建一条曲线，就可以还原公钥，进一步校验。
	PubKey []byte
}

//创建密钥对
func newWalletKeyPair() *wallet {
	curve := elliptic.P256()
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println("生成private key错误")
		return nil
	}
	//获取公钥
	pubKeyRaw := priKey.PublicKey
	//将公钥中的X和Y进行拼接
	pubKey := append(pubKeyRaw.X.Bytes(), pubKeyRaw.Y.Bytes()...)
	//创建wallet结构
	wallet := wallet{}
	wallet.PriKey = priKey
	wallet.PubKey = pubKey
	return &wallet
}

//根据私钥生成地址
func (w *wallet) getAddress() string {
	//获得公钥
	publicKey := w.PubKey
	//生成公钥哈希（锁定output时使用的哈希）
	pubKeyHash := GetPubKeyHashFromPubKey(publicKey)
	//拼接version和公钥哈希，生成21字节的哈希
	payload := append([]byte{byte(00)}, pubKeyHash...)
	//生成四字节交易码
	checksum := checkSum(payload)
	//拼接校验值和21字节哈希，得到钱包地址
	payload = append(payload, checksum...)
	address := base58.Encode(payload)
	return address
}

//根据公钥获取公钥哈希
func GetPubKeyHashFromPubKey(PubKey []byte) []byte {
	hash1 := sha256.Sum256(PubKey)
	hasher := ripemd160.New()
	hasher.Write(hash1[:])
	//锁定output时使用的公钥哈希
	pubKeyHash := hasher.Sum(nil)
	return pubKeyHash
}

//根据钱包地址获取公钥哈希
func GetPubKeyHashFromAddress(address string) []byte {
	//base58解码
	decodeInfo := base58.Decode(address)
	//校验地址
	if len(decodeInfo) != 25 {
		log.Println(" GetPubKeyHashFromAddress 传入地址无效")
		return nil
	}
	//截取
	return decodeInfo[1 : len(decodeInfo)-4]

}

//得到4字节的校验码
func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	checkSum := second[0:4]
	return checkSum
}

//校验地址是否有效
func isValidAddress(address string) bool {
	decodeInfo := base58.Decode(address)
	//校验地址
	if len(decodeInfo) != 25 {
		log.Println(" 传入地址长度无效")
		return false
	}
	payload := decodeInfo[:len(decodeInfo)-4]
	checkSum1 := decodeInfo[len(decodeInfo)-4:]
	checkSum2 := checkSum(payload)
	return bytes.Equal(checkSum1,checkSum2)
}
