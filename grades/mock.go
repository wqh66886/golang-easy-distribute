package grades

func init() {
	students = []Student{
		{
			FirstName: "James",
			LastName:  "Bond",
			ID:        7,
			Grades: []Grade{
				{
					Type:  Exam,
					Title: "Midterm",
					Score: 95.0,
				},
				{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
		{
			FirstName: "Jane",
			LastName:  "De",
			ID:        9,
			Grades: []Grade{
				{
					Type:  Exam,
					Title: "Midterm",
					Score: 100.0,
				},
				{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
		{
			FirstName: "John",
			LastName:  "Do",
			ID:        10,
			Grades: []Grade{
				{
					Type:  Exam,
					Title: "Midterm",
					Score: 100.0,
				},
				{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
			ID:        11,
			Grades: []Grade{
				{
					Type:  Exam,
					Title: "Midterm",
					Score: 100.0,
				},
				{
					Type:  Exam,
					Title: "Final",
					Score: 100.0,
				},
				{
					Type:  Quiz,
					Title: "Quiz 1",
					Score: 100.0,
				},
			},
		},
	}
}
