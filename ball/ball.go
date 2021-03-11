package ball

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

var (
	count, count2 int64
	blueMap       = make(map[int]int)
)

func Start() {
	fmt.Println("随机开始！！！")
	excel.Start()
	cal()
	fmt.Println("遍历完结！ count -> ", count)
	predict()
	fmt.Println("结束！！！")
}

func cal() {
	for {
		rb := radNew()
		count = count + 1
		if reflect.DeepEqual(rb[:6], excel.PrizeBall[:6]) {
			count2++
			fmt.Println("Second Prize!!! count->", count2, "	total ->", count)
			if rb[6] == excel.PrizeBall[6] {
				return
			}
		}
	}

}

func predict() {
	file2, _ := os.OpenFile("./ball.txt", os.O_WRONLY|os.O_CREATE, 0666)
	bufWriter2 := bufio.NewWriterSize(file2, 4096)
	for i := 0; i < 5; i++ {
		rad, _ := rand.Int(rand.Reader, big.NewInt(count))
		var rb [7]int
		for i := rad.Int64(); i > 0; i-- {
			rb = radNew()
		}
		nrb := checkBlue(rb)
		ballStr, _ := json.Marshal(nrb)
		bufWriter2.Write(ballStr)
		bufWriter2.WriteString("\n")
	}
	bufWriter2.Flush()

}

func radNew() [7]int {
	rba, bba := initBall()
	var rb [7]int
	for i := 0; i < 6; i++ {
		rad, _ := rand.Int(rand.Reader, big.NewInt(int64(len(rba))))
		redIndex := int(rad.Int64()) + 1
		rb[i] = new(rba, redIndex)
	}
	sort.Ints(rb[:6])
	rad, _ := rand.Int(rand.Reader, big.NewInt(int64(len(bba))))
	blueIndex := int(rad.Int64()) + 1
	rb[6] = new(bba, blueIndex)
	return rb
}

func initBall() (map[int]int, map[int]int) {
	rba := make(map[int]int)
	bba := make(map[int]int)
	for i := 1; i <= 33; i++ {
		rba[i] = i
	}
	for i := 1; i <= 16; i++ {
		bba[i] = i
	}
	return rba, bba
}

func new(ba map[int]int, index int) int {
	var i = 1
	for key := range ba {
		if i == index {
			delete(ba, key)
			return key
		}
		i++
	}
	return i
}

func checkBlue(radBall [7]int) [7]int {
	if blueMap[radBall[6]] == 0 {
		return radBall
	}
	for {
		newRadBall := radNew()
		if blueMap[newRadBall[6]] == 0 {
			return newRadBall
		}
	}
}
