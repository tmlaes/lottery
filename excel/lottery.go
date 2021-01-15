package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"sort"
	"strconv"
)

const (
	filePath = "ball.xlsx"
)

var (
	PrizeBall [7]int
	LastBalls [7][2]int
	NextBalls [7][2]int
)

func Start() {
	excel, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	readLast(excel)
	readPrize(excel)
	readNext(excel)
	fmt.Println(LastBalls)
	fmt.Println(NextBalls)
}

/**
 * 读取Excel
 */
func readPrize(excel *excelize.File) {
	rows := excel.GetRows("Sheet2")
	PrizeBall[0], _ = strconv.Atoi(rows[0][0])
	PrizeBall[1], _ = strconv.Atoi(rows[0][1])
	PrizeBall[2], _ = strconv.Atoi(rows[0][2])
	PrizeBall[3], _ = strconv.Atoi(rows[0][3])
	PrizeBall[4], _ = strconv.Atoi(rows[0][4])
	PrizeBall[5], _ = strconv.Atoi(rows[0][5])
	PrizeBall[6], _ = strconv.Atoi(rows[0][6])
}

func readLast(excel *excelize.File) {
	var balls [7][]int
	rows := excel.GetRows("Sheet1")
	for _, row := range rows {
		for i, e := range row {
			b, _ := strconv.Atoi(e)
			balls[i] = append(balls[i], b)
		}
	}
	handle(balls, false)
}

func readNext(excel *excelize.File) {
	var balls [7][]int
	rows := excel.GetRows("Sheet3")
	for _, row := range rows {
		for i, e := range row {
			b, _ := strconv.Atoi(e)
			balls[i] = append(balls[i], b)
		}
	}
	handle(balls, true)
}

func handle(balls [7][]int, isNew bool) {
	for i, ball := range balls {
		temp := make(map[int]int)
		sort.Ints(ball)
		var min int
		var max int
		for _, b := range ball {
			temp[b] = temp[b] + 1
		}
		for j, b := range ball {
			if temp[b] > 2 && min == 0 {
				if j >= len(ball)-temp[b] {
					min = ball[0]
				} else {
					min = b
				}
			}
			mb := ball[len(ball)-j-1]
			if temp[mb] > 2 && max == 0 {
				if len(ball)-j-temp[mb] <= 0 {
					max = ball[len(ball)-1]
				} else {
					max = mb
				}
			}
			if min > 0 && max > 0 {
				break
			}
		}
		if min == 0 {
			min = ball[0]
		}
		if max == 0 {
			max = ball[len(ball)-1]
		}
		if isNew {
			NextBalls[i][0] = min
			NextBalls[i][1] = max
		} else {
			LastBalls[i][0] = min
			LastBalls[i][1] = max
		}
	}
}
