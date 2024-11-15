package enums

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
