package rball

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"lottery/excel"
	"math/big"
	"os"
	"reflect"
	"sort"
)

const MAX = 1 << 25

var (
	count, count2 int64
	blueMap       = make(map[int]int)
)

func Start() {
	fmt.Println("开始！！！")
	excel.Start()
	cal()
	fmt.Println("遍历完结！ count -> ", count)
	win()
	fmt.Println("结束！！！")
}

func cal() {
	for {
		rb := produce(excel.LastBalls)
		count = count + 1
		if reflect.DeepEqual(rb[:6], excel.PrizeBall[:6]) {
			count2++
			fmt.Println("Second Prize!!! count->", count2, "	total ->", count)
		}
		if reflect.DeepEqual(rb, excel.PrizeBall) || count >= MAX {
			return
		}
	}
}

func produce(balls [7][2]int) [7]int {
	ball := make(map[int]int, 6)
	var rb [7]int
	for i := 0; i < 6; {
		min := balls[i][0]
		max := balls[i][1] - balls[i][0] + 1
		rad, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		r := int(rad.Int64()) + min
		if ball[r] == 0 {
			ball[r] = r
			rb[i] = r
			i++
		}
	}
	sort.Ints(rb[:6])
	min := balls[6][0]
	max := balls[6][1] - balls[6][0] + 1
	rad, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	rb[6] = int(rad.Int64()) + min
	return rb
}

func win() {
	var min = count / 7
	var max = count - min + 1
	file1, _ := os.OpenFile("./all.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	file2, _ := os.OpenFile("./ball.txt", os.O_WRONLY|os.O_CREATE, 0666)
	bufWriter1 := bufio.NewWriterSize(file1, 8192)
	bufWriter2 := bufio.NewWriterSize(file2, 4096)
	defer file1.Close()
	defer file2.Close()
	for j := 0; j < 5; j++ {
		rad, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		t := rad.Int64() + min
		var radBall [7]int
		for i := t; i > 0; i-- {
			radBall = produce(excel.NextBalls)
			s, _ := json.Marshal(radBall)
			bufWriter1.Write(s)
			bufWriter1.WriteString("\n")
		}
		radBall = checkBlue(radBall)
		s, _ := json.Marshal(radBall)
		bufWriter1.Write(s)
		bufWriter1.WriteString("\n")
		bufWriter2.Write(s)
		bufWriter2.WriteString("\n")
	}
	bufWriter1.Flush()
	bufWriter2.Flush()
}

func checkBlue(radBall [7]int) [7]int {
	if blueMap[radBall[6]] != 0 {
		radBall = produce(excel.NextBalls)
		checkBlue(radBall)
	}
	key := radBall[6]
	blueMap[key] = key
	return radBall
}
