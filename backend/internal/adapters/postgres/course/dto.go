package course

import (
	"time"

	"backend/internal/entities"
)

type courseDTO struct {
	ID              string
	AuthorID        string
	SubjectID       string
	Title           string
	Description     *string
	DifficultyLevel int
	CoverImageURL   *string
	IsPublished     bool
	CreatedAt       time.Time
}

type tagDTO struct {
	ID   int
	Name string
	Slug string
}

func (d *tagDTO) toEntity() entities.Tag {
	return entities.Tag{
		ID:   d.ID,
		Name: d.Name,
		Slug: d.Slug,
	}
}

func newCourseDTO(c *entities.Course) courseDTO {
	var desc *string
	if c.Description != "" {
		desc = &c.Description
	}
	var cover *string
	if c.CoverImageURL != "" {
		cover = &c.CoverImageURL
	}

	return courseDTO{
		ID:              c.ID,
		AuthorID:        c.AuthorID,
		SubjectID:       c.SubjectID,
		Title:           c.Title,
		Description:     desc,
		DifficultyLevel: c.DifficultyLevel,
		CoverImageURL:   cover,
		IsPublished:     c.IsPublished,
		CreatedAt:       c.CreatedAt,
	}
}

func (d *courseDTO) toEntity() *entities.Course {
	c := &entities.Course{
		ID:              d.ID,
		AuthorID:        d.AuthorID,
		SubjectID:       d.SubjectID,
		Title:           d.Title,
		DifficultyLevel: d.DifficultyLevel,
		Tags:            []entities.Tag{},
		IsPublished:     d.IsPublished,
		CreatedAt:       d.CreatedAt.UTC(),
	}
	if d.Description != nil {
		c.Description = *d.Description
	}
	if d.CoverImageURL != nil {
		c.CoverImageURL = *d.CoverImageURL
	}
	return c
}

type moduleDTO struct {
	ID         string
	CourseID   string
	Title      string
	OrderIndex int
}

func newModuleDTO(m *entities.Module) moduleDTO {
	return moduleDTO{
		ID:         m.ID,
		CourseID:   m.CourseID,
		Title:      m.Title,
		OrderIndex: m.OrderIndex,
	}
}

type lessonDTO struct {
	ID                string
	ModuleID          string
	Title             string
	ContentText       *string
	VideoURL          *string
	FileAttachmentURL *string
	XPReward          int
	OrderIndex        int
}

func newLessonDTO(l *entities.Lesson) lessonDTO {
	var content, video, file *string
	if l.ContentText != "" {
		content = &l.ContentText
	}
	if l.VideoURL != "" {
		video = &l.VideoURL
	}
	if l.FileAttachmentURL != "" {
		file = &l.FileAttachmentURL
	}

	return lessonDTO{
		ID:                l.ID,
		ModuleID:          l.ModuleID,
		Title:             l.Title,
		ContentText:       content,
		VideoURL:          video,
		FileAttachmentURL: file,
		XPReward:          l.XPReward,
		OrderIndex:        l.OrderIndex,
	}
}

func (l lessonDTO) toEntity() *entities.Lesson {
	return &entities.Lesson{
		ID:                l.ID,
		ModuleID:          l.ID,
		Title:             l.Title,
		ContentText:       *l.ContentText,
		VideoURL:          *l.VideoURL,
		FileAttachmentURL: *l.FileAttachmentURL,
		XPReward:          l.XPReward,
		OrderIndex:        l.OrderIndex,
	}
}

type favoriteDTO struct {
	UserID    string
	CourseID  string
	CreatedAt time.Time
}

func (d *favoriteDTO) toEntity() entities.CourseFavorite {
	return entities.CourseFavorite{
		UserID:    d.UserID,
		CourseID:  d.CourseID,
		CreatedAt: d.CreatedAt.UTC(),
	}
}
