package main

import (
	"bufio"
	"fmt"
	"log"
	"lottery/ball"
	"lottery/spider"
	"os"
	"strconv"
	"strings"
)

const (
	THIRTY  int = 31
	FIFTY   int = 51
	HUNDRED int = 101
)

func main() {
	limit := readCmd()
	spider.Start(limit)
	ball.Start()
}

func readCmd() (limit int) {
	fmt.Println("请选择数目：")
	fmt.Println("1：30个 \t 2：50个 \t 3：100个")
	in := bufio.NewReader(os.Stdin)
	str, _, err := in.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	limits, _ := strconv.Atoi(strings.TrimSpace(string(str)))
	switch limits {
	case 1:
		limit = THIRTY
	case 2:
		limit = FIFTY
	case 3:
		limit = HUNDRED
	default:
		limit = THIRTY
	}
	return
}
