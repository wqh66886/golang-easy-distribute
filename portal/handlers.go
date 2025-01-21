package portal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wqh/easy/distribute/grades"
	"github.com/wqh/easy/distribute/registry"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/**
* description:
* author: wqh
* date: 2025/1/21
 */

func RegisterHandlerFunc() {
	http.Handle("/", http.RedirectHandler("/students", http.StatusPermanentRedirect))
	h := studentsHandler{}
	http.Handle("/students", h)
	http.Handle("/students/", h)
}

type studentsHandler struct{}

func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	switch len(paths) {
	case 2:
		sh.renderStudents(w, r)
	case 3:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sh.renderStudent(w, r, id)
	case 4:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if strings.ToLower(paths[3]) != "grades" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sh.renderGrades(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (sh studentsHandler) renderStudents(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error receiving students", err)
		}
	}()
	providers, err := registry.GetProviders(registry.GradeService)
	if err != nil {
		return
	}
	res, err := http.Get(providers[0] + "/students")
	if err != nil {
		return
	}
	var g grades.Students
	err = json.NewDecoder(res.Body).Decode(&g)
	if err != nil {
		return
	}
	err = rootTemplate.Lookup("students.html").Execute(w, g)
	if err != nil {
		return
	}
}

func (sh studentsHandler) renderStudent(w http.ResponseWriter, r *http.Request, id int) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error receiving student", err)
		}
	}()
	providers, err := registry.GetProviders(registry.GradeService)
	if err != nil {
		return
	}

	res, err := http.Get(fmt.Sprintf("%v/students/%v", providers[0], id))
	if err != nil {
		return
	}
	var g grades.Student
	err = json.NewDecoder(res.Body).Decode(&g)
	if err != nil {
		return
	}
	err = rootTemplate.Lookup("student.html").Execute(w, g)
	if err != nil {
		return
	}
}

func (sh studentsHandler) renderGrades(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		w.Header().Add("location", fmt.Sprintf("/students/%v", id))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()
	title := r.FormValue("Title")
	gradeType := r.FormValue("Type")
	score, err := strconv.ParseFloat(r.FormValue("Score"), 32)
	if err != nil {
		log.Println("Error parsing score", err)
		return
	}
	g := grades.Grade{
		Title: title,
		Type:  grades.GradeType(gradeType),
		Score: float32(score),
	}

	data, err := json.Marshal(g)
	if err != nil {
		log.Println("Error marshalling grade", err)
		return
	}
	providers, err := registry.GetProviders(registry.GradeService)
	if err != nil {
		return
	}
	res, err := http.Post(fmt.Sprintf("%v/students/%v/grades", providers[0], id), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.Println("Error sending grade", res.StatusCode)
		return
	}
}
