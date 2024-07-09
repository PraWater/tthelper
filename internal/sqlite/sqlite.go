package sqlite

import (
	"database/sql"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
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

func (store *DBStore) InsertDB(db *sql.DB) {
	store.db = db
}

func (store *DBStore) InitDB() error {
	_, err := store.db.Exec("DROP TABLE IF EXISTS subjects")
	if err != nil {
		return err
	}

	_, err = store.db.Exec(`CREATE TABLE IF NOT EXISTS subjects (
        subject_code TEXT PRIMARY KEY,
        subject_name TEXT
      )`)
	if err != nil {
		return err
	}

	_, err = store.db.Exec("INSERT INTO subjects (subject_code, subject_name) VALUES (?, ?), (?, ?)", "BIO", "Biology", "CS", "Computer Science")
	if err != nil {
		return err
	}

	_, err = store.db.Exec("DROP TABLE IF EXISTS courses")
	if err != nil {
		return err
	}

	_, err = store.db.Exec(`CREATE TABLE IF NOT EXISTS courses (
        subject_code TEXT NOT NULL,
        course_code TEXT NOT NULL,
        course_name TEXT,
        credits INTEGER,
        PRIMARY KEY (subject_code, course_code),
        FOREIGN KEY (subject_code) REFERENCES subjects(subject_code)
    )`)
	if err != nil {
		return err
	}

	_, err = store.db.Exec("DROP TABLE IF EXISTS sections")
	if err != nil {
		return err
	}

	_, err = store.db.Exec(`CREATE TABLE IF NOT EXISTS sections (
        subject_code TEXT NOT NULL,
        course_code TEXT NOT NULL,
        section_type INTEGER NOT NULL,
        section_no INTEGER NOT NULL,
        section_slot TEXT,
        PRIMARY KEY (subject_code, course_code, section_type, section_no),
        FOREIGN KEY (subject_code, course_code) REFERENCES courses(subject_code, course_code)
    )`)
	if err != nil {
		return err
	}

	return nil
}

func (store *DBStore) InsertCourses(courses [][]string) error {
    type parseResult struct {
        course Course
        err    error
    }

    results := make(chan parseResult, len(courses))
    sem := make(chan struct{}, 20)

    for _, course := range courses {
        sem <- struct{}{}
        go func(course []string) {
            defer func() { <-sem }()
            c, err := ParseCourse(course)
            results <- parseResult{c, err}
        }(course)
    }

    for range courses {
        result := <-results
        if result.err != nil {
            continue
        }
        _, err := store.db.Exec("INSERT INTO courses (subject_code, course_code, course_name, credits) VALUES (?, ?, ?, ?)", 
            result.course.Subject_code, result.course.Course_code, result.course.Course_name, result.course.Credits)
        if err != nil {
            return err
        }
    }

    return nil
}

func (store *DBStore) InsertSections(sections [][]string) error {
    type parseResult struct {
        section Section
        err     error
        valid   bool
    }

    results := make(chan parseResult, len(sections))
    sem := make(chan struct{}, 20)

    for _, section := range sections {
        sem <- struct{}{}
        go func(section []string) {
            defer func() { <-sem }()
            s, err := ParseSection(section)
            valid := err == nil && s.Section_type >= 0
            results <- parseResult{s, err, valid}
        }(section)
    }

    for range sections {
        result := <-results
        if !result.valid {
            continue
        }
        _, err := store.db.Exec("INSERT INTO sections (subject_code, course_code, section_type, section_no, section_slot) VALUES (?, ?, ?, ?, ?)", 
            result.section.Subject_code, result.section.Course_code, result.section.Section_type, result.section.Section_no, result.section.Section_slot)
        if err != nil {
            return err
        }
    }

    return nil
}

func ParseCourse(course []string) (Course, error) {
	codes := strings.Split(course[0], " ")
	courseName := course[1]
	credits, err := strconv.Atoi(course[2])
	if err != nil {
		return Course{}, err
	}

	return Course{Subject_code: codes[0], Course_code: codes[1], Course_name: courseName, Credits: credits}, nil
}

func ParseSection(section []string) (Section, error) {
	codes := strings.Split(section[0], " ")
	secType := 0
	switch section[1][0] {
	case 'L':
		secType = 0
	case 'T':
		secType = 1
	case 'P':
		secType = 2
	default:
		secType = -1
	}
	secNo, err := strconv.Atoi(section[1][1:])
	if err != nil {
		return Section{}, err
	}

	if len(section) == 3 {
		//Section slot in included
		return Section{Subject_code: codes[0], Course_code: codes[1], Section_type: secType, Section_no: secNo, Section_slot: section[2]}, nil
	}
	return Section{Subject_code: codes[0], Course_code: codes[1], Section_type: secType, Section_no: secNo}, nil
}

func (store *DBStore) FindSlot(section Section) (slot string, err error) {
	row := store.db.QueryRow("SELECT section_slot FROM sections WHERE subject_code = ? AND course_code = ? AND section_no = ? AND section_type = ?", section.Subject_code, section.Course_code, section.Section_no, section.Section_type)
	err = row.Scan(&slot)
	return
}

func (store *DBStore) FindSections(course Course, secType int) (sections []Section, err error) {
	rows, err := store.db.Query("SELECT * FROM sections WHERE subject_code = ? AND course_code = ? AND section_type = ?", course.Subject_code, course.Course_code, secType)
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return
	}

	for rows.Next() {
		s := Section{}
		err = rows.Scan(&s.Subject_code, &s.Course_code, &s.Section_type, &s.Section_no, &s.Section_slot)

		if err != nil {
			return
		}

		sections = append(sections, s)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (store *DBStore) AllCourses() (courses []Course, err error) {
	rows, err := store.db.Query("SELECT * FROM courses")
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return
	}

	for rows.Next() {
		c := Course{}
		err = rows.Scan(&c.Subject_code, &c.Course_code, &c.Course_name, &c.Credits)

		if err != nil {
			return
		}

		courses = append(courses, c)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}
