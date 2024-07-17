package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/PraWater/tthelper/internal/excel"
	"github.com/PraWater/tthelper/internal/sqlite"
	"github.com/PraWater/tthelper/internal/timetable"
	"github.com/adrg/xdg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const NoOfDays = 6
const NoOfSlotsPerDay = 12

func main() {
	pathDb := filepath.Join(xdg.DataHome, "tthelper.db")
	db, err := sql.Open("sqlite", pathDb)
	logError(err)

	store := sqlite.DBStore{}
	store.InsertDB(db)
	logError(err)

	refreshFlag := flag.Bool("refresh", false, "Run the refresh function")
	flag.Parse()
	args := flag.Args()
	var path string
	if *refreshFlag {
		if len(args) > 0 {
			path = args[0]
		} else {
			path = filepath.Join(xdg.Home, "timetable.xlsx")
		}
		refresh(store, path)
	} else {
		if len(args) > 0 {
			path = args[0]
		} else {
			path = filepath.Join(xdg.Home, "input_tt.txt")
		}
		find(store, path)
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func refresh(store sqlite.DBStore, path string) {
	err := store.InitDB()
	logError(err)

	courses, sections := excel.ReadTT(path)
	err = store.InsertCourses(courses)
	logError(err)

	err = store.InsertSections(sections)
	logError(err)
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title+i.desc }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func find(store sqlite.DBStore, path string) {
	file, err := os.Open(path)
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
	var wg sync.WaitGroup
	resultChan := make(chan sqlite.Course, len(courses))

	for _, course := range courses {
		wg.Add(1)
		go func(course sqlite.Course) {
			defer wg.Done()

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
				resultChan <- course
			}
		}(course)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	items := []list.Item{}

	for course := range resultChan {
		items = append(items, item{title: course.Subject_code + " " + course.Course_code, desc: course.Course_name})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Available Courses"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
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
			if filledSlots[time[0]][time[1]-1] {
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
