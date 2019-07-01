/*
 * @modifiy by:Mine&Coin&Chain
 * @Filename:main
 * @Description:添加钱包管理类
 * @Date:2019/7/1 13:59
 * @Version:v1.0
*/

/*
	负责对外，管理生成的钱包（公钥私钥）
	私钥1->公钥-》地址1
	私钥2->公钥-》地址2
	私钥3->公钥-》地址3
	私钥4->公钥-》地址4
	实现结构为 map[钱包地址][]*wallet
*/
package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

//钱包存储文件
const walletFile = "wallet.dat"

type walletManager struct {
	Wallets map[string]*wallet
}

//创建walletManager结构
func NewWalletManager() *walletManager {
	var wm walletManager
	//为钱包分配空间
	wm.Wallets = make(map[string]*wallet)
	//从本地加载已经创建的钱包，写入wallet结构
	if !wm.loadFile() {
		return nil
	}
	return &wm
}

//创建新钱包
func (wm *walletManager) createWallet() string {
	//创建密钥对
	w := newWalletKeyPair()
	if w == nil {
		log.Println("newWalletKeyPair 失败")
		return ""
	}
	//获取地址
	address := w.getAddress()
	//将密钥写入磁盘
	wm.Wallets[address] = w
	if !wm.saveFile() {
		return ""
	}
	//返回cli地址
	return address
}

//将钱包存储至硬盘
func (wm *walletManager) saveFile() bool {
	//使用gob对wm进行编码
	var buffer bytes.Buffer
	//gob编码时需要注册接口，对于未注册的接口无法实现序列化
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wm)
	if err != nil {
		log.Println("encoder.Encode err:", err)
		return false
	}
	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println("ioutil.WriteFile error:", err)
		return false
	}
	return true
}

//从硬盘读取钱包
func (wm *walletManager) loadFile() bool {
	if !IsFileExists(walletFile) {
		log.Println("钱包文件不存在，无需加载")
		return true
	}
	//读取文件
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Println("ioutil read err:", err)
		return false
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(wm)
	if err != nil {
		log.Println("decoder.Decode err:", err)
		return false
	}
	return true
}

//将所有钱包打印(地址升序)
func (wm *walletManager) listAddress() []string {
	var addresses []string
	for  address := range wm.Wallets {
		addresses = append(addresses, address)
	}
	sort.Strings(addresses)
	return addresses
}
