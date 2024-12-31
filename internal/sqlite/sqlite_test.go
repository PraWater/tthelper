package sqlite_test

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/PraWater/tthelper/internal/sqlite"
)

var store sqlite.DBStore
var db *sql.DB

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	db, err = sql.Open("sqlite", "test.db")
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	store = sqlite.DBStore{}
	store.InsertDB(db)
	err = store.InitDB()
	if err != nil {
		return -1, fmt.Errorf("Error initialising database: %w", err)
	}

	defer func() {
		for _, t := range []string{"sections", "courses", "subjects"} {
			_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", t))
		}

		db.Close()
	}()

	return m.Run(), nil
}

func TestInsertCourses(t *testing.T) {
	courses := [][]string{{"BIO F215", "BIOPHYSICS", "3", "16/03 FN1", "18/05 FN"}, {"BIO F231", "BIOLOGY PROJECT LAB", "3", "12/03 AN1", "10/05 FN"}}
	err := store.InsertCourses(courses)
	assertNoError(t, err)

	want := sqlite.Course{Subject_code: "BIO", Course_code: "F215", Course_name: "BIOPHYSICS", Credits: 3, Course_midsem: "16/03 FN1", Course_compre: "18/05 FN"}

	row := db.QueryRow("SELECT * FROM courses WHERE course_code = ?", want.Course_code)

	got := sqlite.Course{}
	err = row.Scan(&got.Subject_code, &got.Course_code, &got.Course_name, &got.Credits, &got.Course_midsem, &got.Course_compre)

	assertNoError(t, err)
	assertEqualCourse(t, got, want)
}

func TestInsertSections(t *testing.T) {
	sections := [][]string{{"BIO F215", "L1", "M W F  5"}, {"BIO F215", "T1", "Th  1"}, {"BIO F231", "L1", "S  1"}, {"BIO F231", "P1", "T  4 5  S  7 8"}}
	err := store.InsertSections(sections)
	assertNoError(t, err)

	want := sqlite.Section{Subject_code: "BIO", Course_code: "F215", Section_type: 0, Section_no: 1, Section_slot: "M W F  5"}

	row := db.QueryRow("SELECT * FROM sections WHERE course_code = ? AND section_no = ? AND section_type = ?", want.Course_code, want.Section_no, want.Section_type)

	got := sqlite.Section{}
	err = row.Scan(&got.Subject_code, &got.Course_code, &got.Section_type, &got.Section_no, &got.Section_slot)

	assertNoError(t, err)
	assertEqualSection(t, got, want)
}

func TestParseCourse(t *testing.T) {
	course := []string{"BIO F215", "BIOPHYSICS", "3", "16/03 FN1", "18/05 FN"}
	got, err := sqlite.ParseCourse(course)
	want := sqlite.Course{Subject_code: "BIO", Course_code: "F215", Course_name: "BIOPHYSICS", Credits: 3, Course_midsem: "16/03 FN1", Course_compre: "18/05 FN"}

	assertNoError(t, err)
	assertEqualCourse(t, got, want)
}

func TestParseSection(t *testing.T) {
	section := []string{"BIO F215", "L10", "M W F  5"}
	got, err := sqlite.ParseSection(section)
	want := sqlite.Section{Subject_code: "BIO", Course_code: "F215", Section_type: 0, Section_no: 10, Section_slot: "M W F  5"}

	assertNoError(t, err)
	assertEqualSection(t, got, want)
}

func TestFindSlot(t *testing.T) {
	section := sqlite.Section{Subject_code: "BIO", Course_code: "F215", Section_type: 0, Section_no: 1}
	got, err := store.FindSlot(section)
	want := "M W F  5"

	assertNoError(t, err)
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestFindSections(t *testing.T) {
	course := sqlite.Course{Subject_code: "BIO", Course_code: "F215", Course_name: "BIOPHYSICS", Credits: 3}
	got, err := store.FindSections(course, 0)
	want := []sqlite.Section{{Subject_code: "BIO", Course_code: "F215", Section_type: 0, Section_no: 1, Section_slot: "M W F  5"}}

	assertNoError(t, err)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFindExams(t *testing.T) {
  course := sqlite.Course{Subject_code: "BIO", Course_code: "F215"}
  got, err := store.FindExams(course)
  want := []string{"16/03 FN1", "18/05 FN"}

	assertNoError(t, err)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

func assertEqualCourse(t testing.TB, got, want sqlite.Course) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertEqualSection(t testing.TB, got, want sqlite.Section) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
