package enums

type RoomType uint

const (
	Support RoomType = iota + 1
)

func (roomType RoomType) String() string {
	switch roomType {
	case Support:
		return "Support"
	}
	return ""
}

func GetAllRoomTypes() []RoomType {
	return []RoomType{
		Support,
	}
}
