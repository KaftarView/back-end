package entities

type Test struct {
	Name string `gorm:"type:varchar(50)"`
	Age  int    `gorm:"type:int"`
}
