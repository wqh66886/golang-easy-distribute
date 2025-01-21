package grades

import (
	"fmt"
	"sync"
)

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

var (
	students Students
	mutex    sync.RWMutex
)

func (s Students) GetById(id int) (*Student, error) {
	// for _,student := range s {} //student 是 s 切片的副本,修改 student 里的元素不会影响 s
	for i := range s { // i 是索引, s[i] 是 s 切片的实际元素,对 s[i]作出修改,对 s[i] 的修改会反映到 s
		if s[i].ID == id {
			return &s[i], nil
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
