package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct{}

const Usage = `
正确使用方法：
	./blockchain create <地址> "创建区块链"
	./blockchain print "打印区块链"
	./blockchain getBalance <地址> "获取余额"
	./blockchain send <FROM> <TO> <AMOUNT> <MINER> <DATA> 
	./blockchain createWallet "创建钱包"
	./blockchain listAddress "列举所有的钱包地址"
`

//解析输入命令的方法
func (cli *CLI) Run() {
	cmds := os.Args
	//用户至少输入两个参数
	if len(cmds) < 2 {
		fmt.Println("输入参数无效，请检查!")
		fmt.Println(Usage)
		return
	}
	switch cmds[1] {
	case "create":
		fmt.Println("创建区块被调用!")
		if len(cmds) != 3 {
			fmt.Println("输入参数无效，请检查")
			return
		}
		address := cmds[2]
		cli.createBlockChain(address)
	case "print":
		fmt.Println("打印区块被调用!")
		cli.print()
	case "getBalance":
		fmt.Println("获取余额命令被调用")
		if len(cmds) != 3 {
			fmt.Println("输入参数无效，请检查")
			return
		}
		address := cmds[2]
		cli.GetBalance(address)
	case "send":
		fmt.Println("send 命令被调用")
		if len(cmds) != 7 {
			fmt.Println("输入参数无效，请检查")
			return
		}
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.Atoi(cmds[4])
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from, to, amount, miner, data)
	case "createWallet":
		fmt.Println("创建钱包命令被调用")
		cli.createWallet()
	case "listAddress":
		fmt.Println("列举所有的钱包地址")
		cli.listAddress()
	default:
		fmt.Println("输入参数无效，请检查!")
		fmt.Println(Usage)
	}

}
