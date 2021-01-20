package rball

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"lottery/excel"
	"math/big"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

const MAX = 1 << 27

var (
	count, index, count2, total int64
	blueMap                     = make(map[int]int)
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
	file1, _ := os.OpenFile("./all.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	file2, _ := os.OpenFile("./ball.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer file1.Close()
	defer file2.Close()

	ticker := time.NewTicker(1 * time.Second)
	temp := make(chan [7]int)
	var min = count / 7
	var max = count - min + 1
	count = 0
	for j := 0; j < 5; j++ {
		rad, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		t := rad.Int64() + min
		total = total + t
		go func(t int64) {
			var radBall [7]int
			for i := t; i > 0; i-- {
				radBall = produce(excel.NextBalls)
				temp <- radBall
			}
			radBall = checkBlue(radBall)
			save(radBall, file2)
		}(t)
	}

	for {
		select {
		case radBall := <-temp:
			writeFile(radBall, file1)
		case <-ticker.C:
			if index == total && len(blueMap) == 5 {
				return
			}
		}
	}
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

func writeFile(ball [7]int, file *os.File) {
	atomic.CompareAndSwapInt64(&index, index, index+1)
	buf := bytes.Buffer{}
	for _, red := range ball {
		buf.WriteString(strconv.Itoa(red))
		buf.WriteString("\t")
	}
	buf.WriteString("\n")
	file.Write(buf.Bytes())
}

func save(ball [7]int, file *os.File) {
	buf := bytes.Buffer{}
	for _, red := range ball {
		buf.WriteString(strconv.Itoa(red))
		buf.WriteString("\t")
	}
	buf.WriteString("\n")
	file.Write(buf.Bytes())
}
