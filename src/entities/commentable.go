package entities

type Commentable struct {
	ID       uint      `gorm:"primaryKey"`
	Comments []Comment `gorm:"foreignKey:CommentableID"`
}
