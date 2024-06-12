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
	store.InitDB(db)

	defer func() {
		for _, t := range []string{"sections", "courses", "subjects"} {
			_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", t))
		}

		db.Close()
	}()

	return m.Run(), nil
}

func TestInsertCourses(t *testing.T) {
	courses := [][]string{{"BIO F215", "BIOPHYSICS", "3"}, {"BIO F231", "BIOLOGY PROJECT LAB", "3"}}
	store.InsertCourses(courses)

	want := sqlite.Course{Subject_code: "BIO", Course_code: "F215", Course_name: "BIOPHYSICS", Credits: 3}

	row := db.QueryRow("SELECT * FROM courses WHERE course_code = ?", want.Course_code)

	got := sqlite.Course{}
	err := row.Scan(&got.Subject_code, &got.Course_code, &got.Course_name, &got.Credits)

	assertNoError(t, err)

	assertEqualCourse(t, got, want)
}

func TestInsertSections(t *testing.T) {
	sections := [][]string{{"BIO F215", "L1", "M W F  5"}, {"BIO F215", "T1", "Th  1"}, {"BIO F231", "L1", "S  1"}, {"BIO F231", "P1", "T  4 5  S  7 8"}}
	store.InsertSections(sections)

	want := sqlite.Section{Subject_code: "BIO", Course_code: "F215", Section_type: 0, Section_no: 1, Section_slot: "M W F  5"}

	row := db.QueryRow("SELECT * FROM sections WHERE course_code = ? AND section_no = ? AND section_type = ?", want.Course_code, want.Section_no, want.Section_type)

	got := sqlite.Section{}
	err := row.Scan(&got.Subject_code, &got.Course_code, &got.Section_type, &got.Section_no, &got.Section_slot)

	assertNoError(t, err)

	assertEqualSection(t, got, want)
}

func TestParseCourse(t *testing.T) {
	course := []string{"BIO F215", "BIOPHYSICS", "3"}
	got := sqlite.ParseCourse(course)
	want := sqlite.Course{Subject_code: "BIO", Course_code: "F215", Course_name: "BIOPHYSICS", Credits: 3}

	assertEqualCourse(t, got, want)
}

func TestParseSection(t *testing.T) {
	section := []string{"BIO F215", "L1", "M W F  5"}
	got := sqlite.ParseSection(section)
	want := sqlite.Section{Subject_code: "BIO", Course_code: "F215", Section_type: 0, Section_no: 1, Section_slot: "M W F  5"}

	assertEqualSection(t, got, want)
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

func assertEqualCourse(t testing.TB, got, want sqlite.Course) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertEqualSection(t testing.TB, got, want sqlite.Section) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
