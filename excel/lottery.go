package excel

import (
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
	for i, ball := range balls {
		sort.Ints(ball)
		LastBalls[i][0] = ball[0]
		LastBalls[i][1] = ball[len(ball)-1]
	}
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
	for i, ball := range balls {
		sort.Ints(ball)
		NextBalls[i][0] = ball[0]
		NextBalls[i][1] = ball[len(ball)-1]
	}
}
