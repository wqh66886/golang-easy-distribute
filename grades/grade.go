package grades

import "fmt"

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, grader := range s.Grades {
		result += grader.Score
	}
	return result / float32(len(s.Grades))
}

type Students []Student

var students Students

func (s Students) GetById(id int) (*Student, error) {
	for _, student := range s {
		if student.ID == id {
			return &student, nil
		}
	}
	return nil, fmt.Errorf("student with id %d not found", id)
}

type GradeType string

const (
	Quiz = GradeType("Quiz")
	Exam = GradeType("Exam")
	Test = GradeType("Test")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
