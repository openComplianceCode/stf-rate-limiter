//go:build ignore

package main

import (
	"fmt"
	"github.com/alec-z/osca/auth"
)

func main() {
	str := auth.GenRand(80)
	fmt.Println("rand.Read：", str)
}
