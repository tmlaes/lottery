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
	count   int64
	index   int64
	count2  int64
	total   int64
	file1   *os.File
	file2   *os.File
	blueMap = make(map[int]int)
)

func init() {
	file1, _ = os.OpenFile("./all.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	file2, _ = os.OpenFile("./ball.txt", os.O_WRONLY|os.O_CREATE, 0666)
}
func Start() {
	defer file1.Close()
	defer file2.Close()
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
		rb := produce(lastBalls, 0)
		if reflect.DeepEqual(rb[:6], excel.PrizeBall[:6]) {
			count2++
			fmt.Println("Second Prize!!! count->", count2, "	total ->", count)
		}
		if reflect.DeepEqual(rb, excel.PrizeBall) {
			return
		}
	}
}

func produce(balls [7][]int, times int64) [7]int {
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
	if times > 0 {
		times = times - 1
		go writeFile(rb, file1)
		produce(balls, times)
	}
	return rb
}

func win() {
	var min = count / 7
	var max = count - min + 1
	nextBalls := excel.NextBalls
	for j := 0; j < 5; j++ {
		rand.Seed(time.Now().UnixNano() + int64(j))
		t := rand.Int63n(max) + min
		total = total + t
		fmt.Println(total)
		go func(t int64) {
			radBall := produce(nextBalls, t)
			radBall = checkBlue(radBall)
			save(radBall, file2)
		}(t)
	}
	for {
		select {
		default:
			if index == total {
				return
			}
		}
	}
}

func checkBlue(radBall [7]int) [7]int {
	if blueMap[radBall[6]] != 0 {
		radBall = produce(excel.NextBalls, 0)
		checkBlue(radBall)
	}
	blueMap[radBall[6]] = radBall[6]
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
