package sqlite

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"strconv"
	"strings"
)

type DBStore struct {
	db *sql.DB
}

type Subject struct {
	Subject_code, Subject_name string
}

type Course struct {
	Subject_code, Course_code, Course_name string
	Credits                                int
}

type Section struct {
	Subject_code, Course_code, Section_slot string
	Section_type, Section_no                int
}

func (store *DBStore) InitDB(db *sql.DB) {
	store.db = db

	_, err := store.db.Exec(`CREATE TABLE IF NOT EXISTS subjects (
        subject_code TEXT PRIMARY KEY,
        subject_name TEXT
      )`)
	checkError(err)

	_, err = store.db.Exec("INSERT INTO subjects (subject_code, subject_name) VALUES (?, ?), (?, ?)", "BIO", "Biology", "CS", "Computer Science")
	checkError(err)

	_, err = store.db.Exec(`CREATE TABLE IF NOT EXISTS courses (
        subject_code TEXT NOT NULL,
        course_code TEXT NOT NULL,
        course_name TEXT,
        credits INTEGER,
        PRIMARY KEY (subject_code, course_code),
        FOREIGN KEY (subject_code) REFERENCES subjects(subject_code)
    )`)
	checkError(err)

	_, err = store.db.Exec(`CREATE TABLE IF NOT EXISTS sections (
        subject_code TEXT NOT NULL,
        course_code TEXT NOT NULL,
        section_type INTEGER NOT NULL,
        section_no INTEGER NOT NULL,
        section_slot TEXT,
        PRIMARY KEY (subject_code, course_code, section_type, section_no),
        FOREIGN KEY (subject_code, course_code) REFERENCES courses(subject_code, course_code)
    )`)
	checkError(err)
}

func (store *DBStore) InsertCourses(courses [][]string) {
	for _, course := range courses {
		c := ParseCourse(course)
		_, err := store.db.Exec("INSERT INTO courses (subject_code, course_code, course_name, credits) VALUES (?, ?, ?, ?)", c.Subject_code, c.Course_code, c.Course_name, c.Credits)
		checkError(err)
	}
}

func (store *DBStore) InsertSections(sections [][]string) {
	for _, section := range sections {
		s := ParseSection(section)
		_, err := store.db.Exec("INSERT INTO sections (subject_code, course_code, section_type, section_no, section_slot) VALUES (?, ?, ?, ?, ?)", s.Subject_code, s.Course_code, s.Section_type, s.Section_no, s.Section_slot)
		checkError(err)
	}
}

func ParseCourse(course []string) Course {
	codes := strings.Split(course[0], " ")
	courseName := course[1]
	credits, err := strconv.Atoi(course[2])
	checkError(err)

	return Course{Subject_code: codes[0], Course_code: codes[1], Course_name: courseName, Credits: credits}
}

func ParseSection(section []string) Section {
	codes := strings.Split(section[0], " ")
	secType := 0
	if section[1][0] == 'T' {
		secType = 1
	}
	if section[1][0] == 'P' {
		secType = 2
	}
	secNo := int(section[1][1] - '0')

	return Section{Subject_code: codes[0], Course_code: codes[1], Section_type: secType, Section_no: secNo, Section_slot: section[2]}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
