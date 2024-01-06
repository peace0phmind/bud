package main

import (
	"fmt"
	"os"
)

func main() {
	for _, arg := range os.Args {
		fmt.Println(arg)
	}

	// 打印所有环境变量
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
}
