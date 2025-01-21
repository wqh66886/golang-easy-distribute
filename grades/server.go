package grades

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandleFunc() {
	handler := studentsHandler{}
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}

type studentsHandler struct{}

// 处理集合或者某个id
func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	switch len(paths) {
	case 2:
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sh.getById(w, r, id)
	case 4:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (sh studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	data, err := sh.toJson(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Printf("response data error: %v\n", err)
		return
	}
}

func (sh studentsHandler) getById(w http.ResponseWriter, r *http.Request, id int) {
	mutex.RLock()
	defer mutex.RUnlock()

	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := sh.toJson(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Printf("response data error: %v\n", err)
		return
	}
}

func (sh studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	mutex.RLock()
	defer mutex.RUnlock()
	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var grade Grade
	err = json.NewDecoder(r.Body).Decode(&grade)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("response body closed failed: %v\n", err)
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	student.Grades = append(student.Grades, grade)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJson(grade)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
}

func (sh studentsHandler) toJson(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}
