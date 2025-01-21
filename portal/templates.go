package portal

import "html/template"

/**
* description:
* author: wqh
* date: 2025/1/21
 */

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	rootTemplate, err = rootTemplate.ParseFiles(
		"./portal/students.html",
		"./portal/student.html",
	)
	if err != nil {
		return err
	}
	return nil
}
