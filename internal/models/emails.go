package models

type Email struct {
	Name         string		`json:"name"`
	Subject      string		`json:"subject"`
	Html_content string		`json:"html_content"`
}

