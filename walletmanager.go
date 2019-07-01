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

import "log"

type walletManager struct {
	Wallets map[string]*wallet
}

//创建walletManager结构
func NewWalletManager() *walletManager {
	var wm walletManager
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

	//返回cli地址
	return address
}
