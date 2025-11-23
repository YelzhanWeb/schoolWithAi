package course

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewCourseRepository(connectionURL string) *CourseRepository {
	return &CourseRepository{connectionURL: connectionURL}
}

func (r *CourseRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p
	return nil
}

func (r *CourseRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *CourseRepository) Create(ctx context.Context, course *entities.Course) error {
	d := newCourseDTO(course)
	query := `
		INSERT INTO courses (id, author_id, subject_id, title, description, difficulty_level, tags, cover_image_url, is_published, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query, d.ID, d.AuthorID, d.SubjectID, d.Title, d.Description, d.DifficultyLevel, d.Tags, d.CoverImageURL, d.IsPublished, d.CreatedAt)
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}
	return nil
}

func (r *CourseRepository) GetByID(ctx context.Context, id string) (*entities.Course, error) {
	query := `SELECT id, author_id, subject_id, title, description, difficulty_level, tags, cover_image_url, is_published, created_at FROM courses WHERE id = $1`

	var d courseDTO
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.AuthorID, &d.SubjectID, &d.Title, &d.Description, &d.DifficultyLevel, &d.Tags, &d.CoverImageURL, &d.IsPublished, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("get course: %w", err)
	}
	return d.toEntity(), nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, course *entities.Course) error {
	d := newCourseDTO(course)

	query := `
        UPDATE courses 
        SET title = $2, 
            description = $3, 
            difficulty_level = $4, 
            tags = $5, 
            cover_image_url = $6, 
            is_published = $7
        WHERE id = $1
    `

	tag, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.Title,
		d.Description,
		d.DifficultyLevel,
		d.Tags,
		d.CoverImageURL,
		d.IsPublished,
	)
	if err != nil {
		return fmt.Errorf("update course: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *CourseRepository) DeleteCourse(ctx context.Context, id string) error {
	query := `DELETE FROM courses WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete course: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *CourseRepository) AddModule(ctx context.Context, module *entities.Module) error {
	d := newModuleDTO(module)
	query := `INSERT INTO modules (id, course_id, title, order_index) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, d.ID, d.CourseID, d.Title, d.OrderIndex)
	return err
}

func (r *CourseRepository) ListModulesByCourse(ctx context.Context, courseID string) ([]entities.Module, error) {
	query := `SELECT id, course_id, title, order_index FROM modules WHERE course_id = $1 ORDER BY order_index ASC`
	rows, err := r.pool.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []entities.Module
	for rows.Next() {
		var m entities.Module
		if err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.OrderIndex); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *CourseRepository) UpdateModule(ctx context.Context, module *entities.Module) error {
	d := newModuleDTO(module)

	query := `
        UPDATE modules 
        SET title = $2, 
            order_index = $3 
        WHERE id = $1
    `

	tag, err := r.pool.Exec(ctx, query, d.ID, d.Title, d.OrderIndex)
	if err != nil {
		return fmt.Errorf("update module: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *CourseRepository) DeleteModule(ctx context.Context, id string) error {
	query := `DELETE FROM modules WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete module: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *CourseRepository) AddLesson(ctx context.Context, lesson *entities.Lesson) error {
	d := newLessonDTO(lesson)
	query := `
		INSERT INTO lessons (id, module_id, title, content_text, video_url, file_attachment_url, xp_reward, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query, d.ID, d.ModuleID, d.Title, d.ContentText, d.VideoURL, d.FileAttachmentURL, d.XPReward, d.OrderIndex)
	return err
}

func (r *CourseRepository) GetModuleLessons(ctx context.Context, moduleID string) ([]entities.Lesson, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
        SELECT id, module_id, title, xp_reward, order_index
        FROM lessons
        WHERE module_id = $1
        ORDER BY order_index ASC
    `

	rows, err := r.pool.Query(ctx, query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("list module lessons: %w", err)
	}
	defer rows.Close()

	var lessons []entities.Lesson

	for rows.Next() {
		var l entities.Lesson
		err := rows.Scan(
			&l.ID,
			&l.ModuleID,
			&l.Title,
			&l.XPReward,
			&l.OrderIndex,
		)
		if err != nil {
			return nil, fmt.Errorf("scan lesson structure: %w", err)
		}
		lessons = append(lessons, l)
	}

	return lessons, nil
}

func (r *CourseRepository) GetLessonByID(ctx context.Context, lessonID string) (*entities.Lesson, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
        SELECT id, module_id, title, content_text, video_url, file_attachment_url, xp_reward, order_index
        FROM lessons
        WHERE id = $1
    `

	var d lessonDTO

	err := r.pool.QueryRow(ctx, query, lessonID).Scan(
		&d.ID,
		&d.ModuleID,
		&d.Title,
		&d.ContentText,
		&d.VideoURL,
		&d.FileAttachmentURL,
		&d.XPReward,
		&d.OrderIndex,
	)
	if err != nil {
		return nil, fmt.Errorf("get lesson full: %w", err)
	}

	return d.toEntity(), nil
}

func (r *CourseRepository) UpdateLesson(ctx context.Context, lesson *entities.Lesson) error {
	d := newLessonDTO(lesson)

	query := `
        UPDATE lessons 
        SET title = $2, 
            content_text = $3, 
            video_url = $4, 
            file_attachment_url = $5, 
            xp_reward = $6, 
            order_index = $7
        WHERE id = $1
    `

	tag, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.Title,
		d.ContentText,
		d.VideoURL,
		d.FileAttachmentURL,
		d.XPReward,
		d.OrderIndex,
	)
	if err != nil {
		return fmt.Errorf("update lesson: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *CourseRepository) DeleteLesson(ctx context.Context, id string) error {
	query := `DELETE FROM lessons WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete lesson: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

// GET ALL COURSE STRUCTURE WITH MODULES AND LESSONS
func (r *CourseRepository) GetCourseStructure(ctx context.Context, courseID string) ([]entities.Module, error) {
	modules, err := r.ListModulesByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if len(modules) == 0 {
		return modules, nil
	}

	query := `
        SELECT l.id, l.module_id, l.title, l.xp_reward, l.order_index
        FROM lessons l
        JOIN modules m ON l.module_id = m.id
        WHERE m.course_id = $1
        ORDER BY m.order_index ASC, l.order_index ASC
    `

	rows, err := r.pool.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course structure: %w", err)
	}
	defer rows.Close()

	lessonsMap := make(map[string][]entities.Lesson)
	for rows.Next() {
		var l entities.Lesson
		if err := rows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.XPReward, &l.OrderIndex); err != nil {
			return nil, err
		}
		lessonsMap[l.ModuleID] = append(lessonsMap[l.ModuleID], l)
	}

	for i := range modules {
		if lessons, ok := lessonsMap[modules[i].ID]; ok {
			modules[i].Lessons = lessons
		} else {
			modules[i].Lessons = []entities.Lesson{}
		}
	}

	return modules, nil
}

// func (r *CourseRepository) ReorderLessons(ctx context.Context, updates map[string]int) error {
// 	tx, err := r.pool.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	query := `UPDATE lessons SET order_index = $2 WHERE id = $1`

// 	for id, newOrder := range updates {
// 		_, err := tx.Exec(ctx, query, id, newOrder)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return tx.Commit(ctx)
// }
