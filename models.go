package portanizer

type Post struct {
	ID   uint
	Name string `gorm:"unique;not null"`
	Body string
	Tags []Tag `gorm:"many2many:post_tags;"`
}

func (p *Post) IsValid() bool {
	if p.Name == "" {
		return false
	}
	return true
}

type Tag struct {
	ID   uint
	Name string `gorm:"unique;not null"`
}

func (t *Tag) IsValid() bool {
	if t.Name == "" {
		return false
	}
	return true
}

type User struct {
	ID       uint
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"unique;not null"`
}

func (u *User) IsValid() bool {
	if u.Email == "" || u.Password == "" {
		return false
	}
	return true
}
