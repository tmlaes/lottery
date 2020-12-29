package spider

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseUrl = "http://datachart.500.com/ssq/history/newinc/history.php?"
)

func Start(limits int) {
	url := baseUrl + "limit=" + strconv.Itoa(limits) + "&sort=0"
	resp, err1 := http.Get(url)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer resp.Body.Close()

	doc, err2 := goquery.NewDocumentFromReader(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}
	xlsx := excelize.NewFile()
	sheet := xlsx.NewSheet("Sheet1")
	xlsx.NewSheet("Sheet2")
	xlsx.NewSheet("Sheet3")
	doc.Find("#tdata").Find("tr").Each(func(i int, tr *goquery.Selection) {
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			if j < 1 || j > 7 {
				return
			}
			s := string(rune(64 + j))
			ball, _ := strconv.Atoi(strings.TrimSpace(td.Text()))
			if i == 0 {
				xlsx.SetCellValue("Sheet2", s+strconv.Itoa(i+1), ball)
				xlsx.SetCellValue("Sheet3", s+strconv.Itoa(i+1), ball)
				return
			}
			if i == limits-1 {
				xlsx.SetCellValue("Sheet1", s+strconv.Itoa(i), ball)
				return
			}
			xlsx.SetCellValue("Sheet1", s+strconv.Itoa(i), ball)
			xlsx.SetCellValue("Sheet3", s+strconv.Itoa(i+1), ball)
		})
	})
	xlsx.SetActiveSheet(sheet)
	// 根据指定路径保存文件
	xlsx.SaveAs("ball.xlsx")

}
