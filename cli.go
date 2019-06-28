package main

import (
	"os"
	"fmt"
)

type CLI struct{}

const Usage = `
正确使用方法：
	./blockchain create "创建区块链"
	./blockchain addBlock <需要写入的的数据> "添加区块"
	./blockchain print "打印区块链"
`
//解析输入命令的方法
func (cli *CLI) Run(){
	cmds := os.Args
	//用户至少输入两个参数
	if len(cmds)<2 {
		fmt.Println("输入参数无效，请检查!")
		fmt.Println(Usage)
		return
	}
	switch cmds[1]{
	case "create":
		fmt.Println("创建区块被调用!")
		cli.createBlockChain()
		data := cmds[2] //需要检验个数
		cli.addBlock(data)
	case "print":
		fmt.Println("打印区块被调用!")
		cli.print()
	default:
		fmt.Println("输入参数无效，请检查!")
		fmt.Println(Usage)
	}

}


