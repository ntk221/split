package main

import (
	"os"
	"fmt"
	"flag"
	"log"
)

var (
	line_count = flag.Int("l", 1000, "行数を指定してください")
)

func main() {
	flag.Parse()

	fmt.Printf("line_count is %d\n",*line_count)
	args := flag.Args()
	fmt.Println(args)

	if len(args) == 0 {
		log.Fatal("十分な数のコマンドライン引数が与えられていない")
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println(file)
}
