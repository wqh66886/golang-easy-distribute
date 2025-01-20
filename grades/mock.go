package grades

func init() {
	students = Students{
		Student{
			FirstName: "James",
			LastName:  "Bond",
			ID:        7,
			Grades: []Grade{
				Grade{
					Type:  Exam,
					Title: "Midterm",
					Score: 95.0,
				},
				Grade{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				Grade{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
		Student{
			FirstName: "Jane",
			LastName:  "De",
			ID:        9,
			Grades: []Grade{
				Grade{
					Type:  Exam,
					Title: "Midterm",
					Score: 100.0,
				},
				Grade{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				Grade{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
		Student{
			FirstName: "John",
			LastName:  "Do",
			ID:        10,
			Grades: []Grade{
				Grade{
					Type:  Exam,
					Title: "Midterm",
					Score: 100.0,
				},
				Grade{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				Grade{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
		Student{
			FirstName: "Jane",
			LastName:  "Doe",
			ID:        11,
			Grades: []Grade{
				Grade{
					Type:  Exam,
					Title: "Midterm",
					Score: 100.0,
				},
				Grade{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				Grade{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
	}
}
