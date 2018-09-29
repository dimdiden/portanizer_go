package gorm

import (
	"fmt"

	"github.com/dimdiden/portanizer_go"
	"github.com/jinzhu/gorm"
)

type postRepo struct {
	DB *gorm.DB
}

func NewPostRepo(db *gorm.DB) portanizer.PostRepo {
	return &postRepo{DB: db}
}

func (r *postRepo) GetByID(id string) (*portanizer.Post, error) {
	var post portanizer.Post

	if r.DB.First(&post, "id = ?", id).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	r.DB.Preload("Tags").Order("ID ASC").Find(&post)

	return &post, nil
}

func (r *postRepo) GetByName(name string) (*portanizer.Post, error) {
	var post portanizer.Post

	if r.DB.First(&post, "name = ?", name).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	r.DB.Preload("Tags").Order("ID ASC").Find(&post)

	return &post, nil
}

func (r *postRepo) GetList() ([]*portanizer.Post, error) {
	var posts []*portanizer.Post
	if err := r.DB.Preload("Tags").Order("ID ASC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) Create(post portanizer.Post) (*portanizer.Post, error) {
	if !r.DB.First(&post, "name = ?", post.Name).RecordNotFound() {
		return nil, portanizer.ErrExists
	}

	newPost := portanizer.Post{Name: post.Name, Body: post.Body}
	if err := r.DB.Create(&newPost).Error; err != nil {
		return nil, err
	}

	for _, t := range post.Tags {
		r.DB.FirstOrCreate(&t, t)
		r.DB.Model(&newPost).Association("Tags").Append(t)
	}
	return &newPost, nil
}

func (r *postRepo) Update(id string, post portanizer.Post) (*portanizer.Post, error) {
	if !r.DB.First(&post, "name = ?", post.Name).RecordNotFound() && id != fmt.Sprint(post.ID) {
		return nil, portanizer.ErrExists
	}

	var updPost portanizer.Post
	if r.DB.First(&updPost, "id = ?", id).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}

	if err := r.DB.Model(&updPost).Update(portanizer.Post{Name: post.Name, Body: post.Body}).Error; err != nil {
		return nil, err
	}
	// Create tag if doesn't exist and assign tags to post
	for _, t := range post.Tags {
		r.DB.FirstOrCreate(&t, t)
		r.DB.Model(&updPost).Association("Tags").Append(t)
	}
	return &updPost, nil
}

func (r *postRepo) PutTags(pid string, tagids []string) (*portanizer.Post, error) {
	var post portanizer.Post
	if r.DB.First(&post, "id = ?", pid).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}

	for _, id := range tagids {
		var tag portanizer.Tag
		if r.DB.First(&tag, "id = ?", id).RecordNotFound() { // <= if not found then should be skipped
			return nil, portanizer.ErrNotFound
		}
		r.DB.Model(&post).Association("Tags").Append(tag)
	}
	return &post, nil
}

func (r *postRepo) Delete(id string) error {
	var post portanizer.Post
	if r.DB.First(&post, "id = ?", id).RecordNotFound() {
		return portanizer.ErrNotFound
	}
	if err := r.DB.Delete(&post).Error; err != nil {
		return err
	}
	return nil
}
