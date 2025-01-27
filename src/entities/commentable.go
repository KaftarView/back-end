package entities

type Commentable struct {
	CID      uint      `gorm:"primaryKey;autoIncrement"`
	Comments []Comment `gorm:"foreignKey:CommentableID;constraint:OnDelete:CASCADE"`
}
