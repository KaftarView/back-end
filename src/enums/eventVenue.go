package enums

type EventVenue uint

const (
	Online EventVenue = iota + 1
	Physical
	Hybrid
)

func (v EventVenue) String() string {
	switch v {
	case Online:
		return "Online"
	case Physical:
		return "Physical"
	case Hybrid:
		return "Hybrid"
	}
	return ""
}

func GetAllEventVenues() []EventVenue {
	return []EventVenue{
		Online,
		Physical,
		Hybrid,
	}
}
