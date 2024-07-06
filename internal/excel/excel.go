package excel

import (
	"github.com/xuri/excelize/v2"
	"log"
)

const (
	courseCode    = 1
	courseName    = 2
	courseCredits = 5
	sectionCode   = 6
	sectionSlot   = 9
)

func ReadTT(path string) (courses [][]string, sections [][]string) {
	exc, err := excelize.OpenFile(path)
	checkError(err)
	defer exc.Close()

	rows, err := exc.GetRows(exc.GetSheetName(0))
	checkError(err)

	var curCourse string
	for _, row := range rows {
		if row[courseCode] != "" && len(row) >= 6{
			curCourse = row[courseCode]
			courses = append(courses, []string{row[courseCode], row[courseName], row[courseCredits]})
		}
		if row[sectionCode] != "" && len(row) >= 10 {
			sections = append(sections, []string{curCourse, row[sectionCode], row[sectionSlot]})
		}
	}
	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
