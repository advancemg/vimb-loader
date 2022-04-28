package main

import (
	"fmt"
)

func main() {
	fmt.Println(128 << 20)
	fmt.Println((128 << 20) / (64 << 20))
	//ballast := make([]byte, 64<<20)
	//create, err := os.Create("1.txt")
	//if err != nil {
	//	panic(err)
	//}
	//defer create.Close()
	//create.Write(ballast)
}
