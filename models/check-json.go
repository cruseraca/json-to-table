package models

type CheckJsonResponse struct {
	NumberOfFields int     `json:"numberOfFields"`
	MaximumDepth   int     `json:"maximumDepth"`
	ListOfFields   []Field `json:"listOfFields"`
}

type Field struct {
	Name  string `json:"name"`
	Type string `json:"type"`
	Value any `json:"value,omitempty"`
}