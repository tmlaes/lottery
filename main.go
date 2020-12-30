package main

import (
	"bufio"
	"fmt"
	"log"
	"lottery/ball"
	"lottery/rball"
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
	limit, rad := readCmd()
	spider.Start(limit)
	if rad {
		rball.Start()
	} else {
		ball.Start()
	}

}

func readCmd() (limit int, rad bool) {
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

	fmt.Println("是否随机：")
	fmt.Println("1：是 \t 2：否")
	str1, _, err1 := in.ReadLine()
	if err1 != nil {
		log.Fatal(err1)
	}
	rd, _ := strconv.Atoi(strings.TrimSpace(string(str1)))
	if rd == 2 {
		rad = false
	} else {
		rad = true
	}
	return
}
