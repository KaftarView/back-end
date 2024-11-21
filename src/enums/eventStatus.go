package enums

type EventStatus uint

const (
	Draft EventStatus = iota + 1
	Published
	Cancelled
	Completed
)

func (s EventStatus) String() string {
	switch s {
	case Draft:
		return "Draft"
	case Published:
		return "Published"
	case Cancelled:
		return "Cancelled"
	case Completed:
		return "Completed"
	}
	return ""
}

func GetAllEventStatus() []EventStatus {
	return []EventStatus{
		Draft,
		Published,
		Cancelled,
		Completed,
	}
}
