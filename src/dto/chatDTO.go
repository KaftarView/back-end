package dto

type RoomDetailsResponse struct {
	ID     uint                  `json:"id"`
	Tag    string                `json:"tag"`
	Admins []UserDetailsResponse `json:"admins"`
	Member UserDetailsResponse   `json:"member"`
}

type MessageDetailsResponse struct {
	Sender  UserDetailsResponse `json:"sender"`
	Content string              `json:"content"`
}
