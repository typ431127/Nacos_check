package fmtd

import (
	"fmt"
	"os"
)

// Fatalln 自定义退出
func Fatalln(a ...any) {
	fmt.Println(a)
	os.Exit(2)
}
