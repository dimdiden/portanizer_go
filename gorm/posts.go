package gorm

import (
	app "github.com/dimdiden/portanizer_sop"
	"github.com/jinzhu/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func (s *PostService) GetByID(id string) (*app.Post, error) {
	var post app.Post

	if s.DB.First(&post, "id = ?", id).RecordNotFound() {
		return nil, app.ErrNotFound
	}

	s.DB.Find(&post).Order("ID ASC").Preload("Tags")

	return &post, nil
}

func (s *PostService) GetByName(name string) (*app.Post, error) {
	var post *app.Post

	if s.DB.First(&post, "name = ?", name).RecordNotFound() {
		return nil, app.ErrNotFound
	}
	s.DB.Preload("Tags").Order("ID ASC").Find(&post)

	return nil, app.ErrNotFound
}

func (s *PostService) GetList() ([]*app.Post, error) {
	var posts []*app.Post
	if err := s.DB.Preload("Tags").Order("ID ASC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) Create(post app.Post) (*app.Post, error) {
	newPost := app.Post{Name: post.Name, Body: post.Body, Tags: []app.Tag{}}
	if err := s.DB.Create(&newPost).Error; err != nil {
		return nil, err
	}

	for _, t := range post.Tags {
		s.DB.FirstOrCreate(&t, app.Tag{Name: t.Name})
		s.DB.Model(&newPost).Association("Tags").Append(t)
	}
	return &newPost, nil
}

func (s *PostService) Update(id string, post app.Post) (*app.Post, error) {
	var updPost app.Post
	if s.DB.First(&updPost, "id = ?", id).RecordNotFound() {
		return nil, app.ErrNotFound
	}

	if err := s.DB.Model(&updPost).Update(post).Error; err != nil {
		return nil, err
	}
	// Create tag if doesn't exist and assign tags to post
	for _, t := range post.Tags {
		s.DB.FirstOrCreate(&t, t)
		s.DB.Model(&updPost).Association("Tags").Append(t)
	}
	return &updPost, nil
}

func (s *PostService) PutTags(pid string, tagids []string) (*app.Post, error) {
	var post app.Post
	if s.DB.First(&post, "id = ?", pid).RecordNotFound() {
		return nil, app.ErrNotFound
	}

	for _, id := range tagids {
		var tag app.Tag
		if s.DB.First(&tag, "id = ?", id).RecordNotFound() { // <= if not found then should be skipped
			return nil, app.ErrNotFound
		}
		s.DB.Model(&post).Association("Tags").Append(tag)
	}
	return &post, nil
}

func (s *PostService) Delete(id string) error {
	var post app.Post
	if s.DB.First(&post, "id = ?", id).RecordNotFound() {
		return app.ErrNotFound
	}
	if err := s.DB.Delete(&post).Error; err != nil {
		return err
	}
	return nil
}
