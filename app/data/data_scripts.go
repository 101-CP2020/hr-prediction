package main

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func main() {
	fillStat()
}

func fillOkdtpr() {
	dsn := "user=db_user password=db_pwd dbname=hr_db host=92.63.103.157 port=7080"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	wb, err := xlsx.OpenFile("./data/codes1.xlsx")
	if err != nil {
		panic(err)
	}

	sheetStart := 0

	for i, sh := range wb.Sheets {
		fmt.Printf("Sheets number %d\n", i)

		if i < sheetStart {
			continue
		}

		_ = sh.ForEachRow(func(r *xlsx.Row) error {
			code, _ := r.GetCell(2).Int()
			title := r.GetCell(1).String()

			if code <= 0 {
				return nil
			}

			fmt.Printf("Job number %d\n", code)

			db.Table("tbl_okpdtr").Create(map[string]interface{}{
				"okpdtr": code,
				"title":  title,
			})

			return nil
		})
	}
}

func fillOkved() {
	dsn := "user=db_user password=db_pwd dbname=hr_db host=92.63.103.157 port=7080"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	wb, err := xlsx.OpenFile("./data/okved.xlsx")
	if err != nil {
		panic(err)
	}

	sheet := wb.Sheets[0]

	sheet.ForEachRow(func(r *xlsx.Row) error {
		code := r.GetCell(1).String()
		title := r.GetCell(2).String()

		if title == "" || code == "" {
			return nil
		}

		fmt.Printf("Job number %s\n", code)

		db.Table("tbl_okved").Create(map[string]interface{}{
			"okved": code,
			"title": title,
		})

		return nil
	})
}

type Okved struct {
	ID    string `gorm:"column:okved"`
	Title string `gorm:"column:title"`
}

func fillStat() {
	dsn := "user=db_user password=db_pwd dbname=hr_db host=92.63.103.157 port=7080"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	wb, err := xlsx.OpenFile("./ВакансииНаКонец2019.xlsx")
	if err != nil {
		panic(err)
	}

	sheet := wb.Sheets[0]

	currentClassName := ""
	currentClassId := ""

	sheet.ForEachRow(func(r *xlsx.Row) error {
		class := okvedTitle(r.GetCell(0).String())
		group := okpdtrFromGroup(r.GetCell(1).String())
		value, _ := r.GetCell(2).Int()

		if group < 1 {
			return nil
		}

		if class != "" && class != currentClassName {
			okved := Okved{
				Title: class,
			}
			result := db.Table("tbl_okved").Take(&okved, okved)
			if result.Error != nil {
				panic(result.Error)
			}

			fmt.Println(okved)

			currentClassId = okved.ID
			currentClassName = okved.Title
		}

		fmt.Printf("Job '%s' '%d'\n", currentClassName, group)

		db.Table("tbl_vacancies").Create(map[string]interface{}{
			"okved":      currentClassId,
			"okpdtr":     group,
			"kladr":      3900000000000,
			"number":     value,
			"created_by": 1,
			"created_at": time.Now().Unix(),
			"month_year": 122019,
		})

		return nil
	})
}

func okvedTitle(s string) string {
	if strings.Contains(s, "Итог") || !strings.Contains(s, "РАЗДЕЛ") {
		return ""
	}

	return strings.Title(strings.ToLower(s[15:19])) + strings.ToLower(s[19:])
}

func okpdtrFromGroup(s string) int {
	if !strings.Contains(s, "(") {
		return -1
	}

	s = strings.Split(strings.TrimSpace(s), "(")[1]
	s = s[:len(s)-2]
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}

	return i
}
