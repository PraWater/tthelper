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
const NoOfSlotsPerDay = 11

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
		slot, err := store.FindSlot(section)
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

  for i := 0; i < NoOfSlotsPerDay; i++ {
    for j:=0; j < NoOfDays; j++ {
      if filledSlots[j][i] {
        fmt.Printf("x")
      } else {
        fmt.Printf(".")
      }
    }
    fmt.Println()
  }
}
