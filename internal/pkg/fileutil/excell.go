package fileutil

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func ReadFile(filename string) {
	// 打开Excel文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// 获取工作表列表
	sheets := f.GetSheetList()

	// 读取第一个工作表的数据
	for _, sheet := range sheets {
		fmt.Printf("工作表名称: %s\n", sheet)
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Fatal(err)
		}

		for _, row := range rows {
			for _, colCell := range row {
				fmt.Print(colCell, "\t")
			}
			fmt.Println()
		}
	}
}
