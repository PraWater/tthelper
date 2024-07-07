package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/PraWater/tthelper/internal/excel"
	"github.com/PraWater/tthelper/internal/sqlite"
	"github.com/PraWater/tthelper/internal/timetable"
)

const NoOfDays = 6
const NoOfSlotsPerDay = 12

func main() {
	db, err := sql.Open("sqlite", "timetable.db")
	logError(err)

	store := sqlite.DBStore{}
	store.InsertDB(db)
	logError(err)

	if len(os.Args) < 2 {
		list(store)
	} else {
		switch os.Args[1] {
		case "refresh":
			refresh(store)
		}
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func refresh(store sqlite.DBStore) {
	err := store.InitDB()
	logError(err)

	courses, sections := excel.ReadTT("timetable.xlsx")
	err = store.InsertCourses(courses)
	logError(err)

	err = store.InsertSections(sections)
	logError(err)
}

func list(store sqlite.DBStore) {
	file, err := os.Open("input_tt.txt")
	logError(err)
	defer file.Close()

	userSections, err := timetable.ReadFile(file)
	logError(err)

	filledSlots := make([][]bool, NoOfDays)
	for i := range filledSlots {
		filledSlots[i] = make([]bool, NoOfSlotsPerDay)
	}

	for _, section := range userSections {
		s, err := sqlite.ParseSection(section)
		logError(err)

		slot, err := store.FindSlot(s)
		logError(err)

		times, err := timetable.ParseSlot(slot)
		logError(err)

		for _, time := range times {
			if filledSlots[time[0]][time[1]-1] {
				log.Fatal("Conflict in User's timetable")
			}
			filledSlots[time[0]][time[1]-1] = true
		}
	}

	courses, err := store.AllCourses()
	logError(err)
	for _, course := range courses {
		fmt.Println(course)
		lSections, err := store.FindSections(course, 0)
		logError(err)
    takeLec := canTakeSections(store, lSections, filledSlots)

		tSections, err := store.FindSections(course, 1)
		logError(err)
    takeTut := canTakeSections(store, tSections, filledSlots)

		pSections, err := store.FindSections(course, 2)
		logError(err)
    takePra := canTakeSections(store, pSections, filledSlots)

		if takeLec && takeTut && takePra {
			fmt.Println(course)
		}
	}

}

func canTakeSections(store sqlite.DBStore, sections []sqlite.Section, filledSlots [][]bool) bool {
  if len(sections) == 0 {
    return true
  }
	for _, section := range sections {
		s, err := store.FindSlot(section)
		logError(err)

		times, err := timetable.ParseSlot(s)
		logError(err)

		canTake := true
		for _, time := range times {
			if filledSlots[time[0]][time[1]] {
				canTake = false
				break
			}
		}
		if canTake {
			return true
		}
	}
	return false
}
