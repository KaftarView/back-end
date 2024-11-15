package enums

import (
	"errors"
)

type CategoryType uint

const (
	Public CategoryType = iota + 1
	Courses
	Events
	Contests
)

func (c CategoryType) String() string {
	switch c {
	case Public:
		return "Public"
	case Courses:
		return "Courses"
	case Events:
		return "Events"
	case Contests:
		return "Contests"
	}
	return ""
}

func GetAllCategoryTypes() []CategoryType {
	return []CategoryType{
		Public,
		Courses,
		Events,
		Contests,
	}
}

func GetAllCategoryStrings() []string {
	var categories []string
	for _, category := range GetAllCategoryTypes() {
		categories = append(categories, category.String())
	}
	return categories
}

var stringToCategoryType = map[string]CategoryType{
	"Public":   Public,
	"Courses":  Courses,
	"Events":   Events,
	"Contests": Contests,
}

func GetCategoryTypeByName(name string) (CategoryType, error) {
	category, exists := stringToCategoryType[name]
	if !exists {
		return 0, errors.New("invalid category name")
	}
	return category, nil
}
