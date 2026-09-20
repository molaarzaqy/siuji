package exceltemplate

import "github.com/xuri/excelize/v2"

// GenerateParticipantTemplate builds the Excel template used for bulk
// participant import, matching the column contract expected by
// ParticipantUseCase.ImportFromExcel (Name, Email, NIM, University).
func GenerateParticipantTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "Name")
	f.SetCellValue(sheet, "B1", "Email")
	f.SetCellValue(sheet, "C1", "NIM")
	f.SetCellValue(sheet, "D1", "University")

	f.SetCellValue(sheet, "A2", "Budi Santoso")
	f.SetCellValue(sheet, "B2", "budi@example.com")
	f.SetCellValue(sheet, "C2", "12345678")
	f.SetCellValue(sheet, "D2", "Universitas Indonesia")

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}