package main

import (
	"fmt"
	"gofemart/internal/luhn"
)

func main() {
	fmt.Println(luhn.LuhnAlgorithm("1345776"))
}
