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
}

func (store *DBStore) InsertCourses(courses [][]string) {
	for _, course := range courses {
		c := ParseCourse(course)
		_, err := store.db.Exec("INSERT INTO courses (subject_code, course_code, course_name, credits) VALUES (?, ?, ?, ?)", c.Subject_code, c.Course_code, c.Course_name, c.Credits)
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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
