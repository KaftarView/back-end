package entities

type Commentable struct {
	CommentableID uint      `gorm:"primaryKey;autoIncrement"`
	Comments      []Comment `gorm:"foreignKey:CommentableID"`
}
