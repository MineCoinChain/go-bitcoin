/*
 * @modifiy by:Mine&Coin&Chain
 * @Filename:main
 * @Description:添加wallet钱包
 * @Date:2019/6/30 23:50
 * @Version:v1.0
*/
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	. "golang.org/x/crypto/ripemd160"
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
	//进行hash256处理
	hash1 := sha256.Sum256(publicKey)
	//hash160处理
	hasher := New()
	hasher.Write(hash1[:])
	//生成公钥哈希（锁定output时使用的哈希）
	pubKeyHash := hasher.Sum(nil)
	//拼接version和公钥哈希，生成21字节的哈希
	payload := append([]byte{byte(00)}, pubKeyHash...)
	//生成4字节的校验值
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	checksum := second[0:4]
	//拼接校验值和21字节哈希，得到钱包地址
	payload = append(payload, checksum...)
	address := base58.Encode(payload)
	return address
}
