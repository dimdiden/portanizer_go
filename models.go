package app

type Post struct {
	ID   uint
	Name string `gorm:"unique;not null"`
	Body string
	Tags []Tag `gorm:"many2many:post_tags;"`
}

type Tag struct {
	ID   uint
	Name string `gorm:"unique;not null"`
}
