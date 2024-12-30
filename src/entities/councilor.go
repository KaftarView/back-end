package entities

type Councilor struct {
	UserID       uint   `gorm:"primaryKey"`
	FirstName    string `gorm:"type:varchar(50);not null"`
	LastName     string `gorm:"type:varchar(50);not null"`
	Description  string `gorm:"text"`
	PromotedYear int    `gorm:"not null"`
	ProfilePath  string `gorm:"text"`
	User         User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
