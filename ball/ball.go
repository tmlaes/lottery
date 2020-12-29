package ball

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

const (
	RED_MAX  = 33
	BLUE_MAX = 16
)

var count int64
var index int64

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
		rb := produce()
		if reflect.DeepEqual(rb, excel.PrizeBall) {
			return
		}
	}
}

func produce() [7]int {
	atomic.CompareAndSwapInt64(&count, count, count+1)
	rand.Seed(time.Now().UnixNano() + count)
	ball := make(map[int]int)
	var rb [7]int
	var reds []int
	for i := 0; len(ball) < 6; {
		r := rand.Intn(RED_MAX) + 1
		if ball[r] == 0 {
			ball[r] = r
			reds = append(reds, r)
			i++
		}
	}
	sort.Ints(reds)
	copy(rb[:], reds)
	rb[6] = rand.Intn(BLUE_MAX) + 1
	return rb
}

func win() {
	file1, _ := os.OpenFile("./all.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	file2, _ := os.OpenFile("./ball.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer file1.Close()
	defer file2.Close()
	c := make(chan [7]int)
	var total = count
	for j := 0; j < 5; j++ {
		go func() {
			var radBall [7]int
			for i := total; i > 0; i-- {
				radBall = produce()
				c <-radBall
			}
			save(radBall, file2)
		}()
	}
	for {
		select {
		case radBall := <-c:
			writeFile(radBall, file1)
		default:
			if index==5*total {
				return
			}
		}
	}
}

func writeFile(ball [7]int, file *os.File) {
	atomic.CompareAndSwapInt64(&index, index, index+1)
	buf := bytes.Buffer{}
	buf.WriteString("[")
	for _, red := range ball {
		buf.WriteString(strconv.Itoa(red))
		buf.WriteString(", ")
	}
	buf.Truncate(buf.Len()-2)
	buf.WriteString("]\n")
	file.Write(buf.Bytes())
}

func save(ball [7]int, file *os.File) {
	buf := bytes.Buffer{}
	buf.WriteString("[")
	for _, red := range ball {
		buf.WriteString(strconv.Itoa(red))
		buf.WriteString(", ")
	}
	buf.Truncate(buf.Len()-2)
	buf.WriteString("]\n")
	file.Write(buf.Bytes())
}
