package main

import (
	"os"
	"fmt"
	"flag"
)

var (
	line_count = flag.Int("l", 1000, "行数を指定してください")
)

func main() {
	flag.Parse()

	fmt.Printf("line_count is %d\n",*line_count)
	args := os.Args
	fmt.Println(args)

	file := args[1]
}
