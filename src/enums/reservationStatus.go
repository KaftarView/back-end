package enums

type ReservationStatus uint

const (
	Pending ReservationStatus = iota + 1
	Expired
	Confirmed
)

func (status ReservationStatus) String() string {
	switch status {
	case Pending:
		return "Pending"
	case Expired:
		return "Expired"
	case Confirmed:
		return "Confirmed"
	}
	return ""
}

func GetAllReservationStatus() []ReservationStatus {
	return []ReservationStatus{
		Pending,
		Expired,
		Confirmed,
	}
}
