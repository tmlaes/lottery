package rball

import (
	"bytes"
	"fmt"
	"lottery/excel"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

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
	lastBalls := excel.LastBalls
	for {
		rb := produce(lastBalls)
		if reflect.DeepEqual(rb[:6], excel.PrizeBall[:6]) {
			count2++
			fmt.Println("Second Prize!!! count->", count2, "	total ->", count)
		}
		if reflect.DeepEqual(rb, excel.PrizeBall) {
			return
		}
	}
}

func produce(balls [7][]int) [7]int {
	atomic.CompareAndSwapInt64(&count, count, count+1)
	rand.Seed(time.Now().UnixNano() + count)
	ball := make(map[int]int)
	var rb [7]int
	var reds []int
	for i := 0; len(ball) < 6; {
		j := rand.Intn(len(balls[i]))
		r := balls[i][j]
		if ball[r] == 0 {
			ball[r] = r
			reds = append(reds, r)
			i++
		}
	}
	sort.Ints(reds)
	copy(rb[:], reds)
	j := rand.Intn(len(balls[6]))
	rb[6] = balls[6][j]
	return rb
}

func win() {
	file1, _ := os.OpenFile("./all.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	file2, _ := os.OpenFile("./ball.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer file1.Close()
	defer file2.Close()
	temp := make(chan [7]int)
	var min = count / 7
	var max = count - min + 1
	nextBalls := excel.NextBalls
	for j := 0; j < 5; j++ {
		t := rand.Int63n(max) + min
		total = total + t
		go func(t int64) {
			var radBall [7]int
			for i := t; i > 0; i-- {
				radBall = produce(nextBalls)
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
		default:
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
