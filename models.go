package portanizer

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

func (p *Post) IsValid() bool {
	if p.Name == "" {
		return false
	}
	return true
}

func (t *Tag) IsValid() bool {
	if t.Name == "" {
		return false
	}
	return true
}
