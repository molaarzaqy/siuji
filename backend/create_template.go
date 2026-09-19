package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Email")
	f.SetCellValue("Sheet1", "C1", "NIM")
	f.SetCellValue("Sheet1", "D1", "University")

	f.SetCellValue("Sheet1", "A2", "Budi Santoso")
	f.SetCellValue("Sheet1", "B2", "budi@example.com")
	f.SetCellValue("Sheet1", "C2", "12345678")
	f.SetCellValue("Sheet1", "D2", "Universitas Indonesia")

	if err := f.SaveAs("c:\\dev\\siuji\\frontend\\public\\template_peserta.xlsx"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Template created successfully!")
	}
}
