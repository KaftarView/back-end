package enums

type DiscountType uint

const (
	Fixed DiscountType = iota + 1
	Percentage
)

func (dt DiscountType) String() string {
	switch dt {
	case Fixed:
		return "Fixed"
	case Percentage:
		return "Percentage"
	}
	return ""
}

func GetAllDiscountTypes() []DiscountType {
	return []DiscountType{
		Fixed,
		Percentage,
	}
}
