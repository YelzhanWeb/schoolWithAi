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
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	d := newCourseDTO(course)
	query := `
		INSERT INTO courses (id, author_id, subject_id, title, description, difficulty_level, cover_image_url, is_published, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(
		ctx,
		query,
		d.ID,
		d.AuthorID,
		d.SubjectID,
		d.Title,
		d.Description,
		d.DifficultyLevel,
		d.CoverImageURL,
		d.IsPublished,
		d.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}

	if err := r.updateCourseTags(ctx, tx, course.ID, course.Tags); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *CourseRepository) GetByID(ctx context.Context, id string) (*entities.Course, error) {
	query := `SELECT id, author_id, subject_id, title, description, difficulty_level, cover_image_url, is_published, created_at FROM courses WHERE id = $1`

	var d courseDTO
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.AuthorID, &d.SubjectID, &d.Title, &d.Description, &d.DifficultyLevel, &d.CoverImageURL, &d.IsPublished, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("get course: %w", err)
	}
	course := d.toEntity()

	tags, err := r.getTagsByCourseID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}
	course.Tags = tags

	return course, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, course *entities.Course) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	d := newCourseDTO(course)

	query := `
        UPDATE courses 
        SET title = $2, 
            description = $3, 
            difficulty_level = $4, 
            cover_image_url = $5, 
            is_published = $6
        WHERE id = $1
    `

	tag, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.Title,
		d.Description,
		d.DifficultyLevel,
		d.CoverImageURL,
		d.IsPublished,
	)
	if err != nil {
		return fmt.Errorf("update course: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	if err := r.updateCourseTags(ctx, tx, course.ID, course.Tags); err != nil {
		return err
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

func (r *CourseRepository) CreateTag(ctx context.Context, name, slug string) (*entities.Tag, error) {
	query := `
		INSERT INTO tags (name, slug) 
		VALUES ($1, $2) 
		ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, slug
	`

	var t tagDTO
	err := r.pool.QueryRow(ctx, query, name, slug).Scan(&t.ID, &t.Name, &t.Slug)
	if err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}

	tag := t.toEntity()
	return &tag, nil
}

func (r *CourseRepository) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	query := `SELECT id, name, slug FROM tags ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all tags: %w", err)
	}
	defer rows.Close()

	var tags []entities.Tag
	for rows.Next() {
		var t tagDTO
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, t.toEntity())
	}
	return tags, nil
}

func (r *CourseRepository) updateCourseTags(
	ctx context.Context,
	tx pgx.Tx,
	courseID string,
	tags []entities.Tag,
) error {
	_, err := tx.Exec(ctx, `DELETE FROM course_tags WHERE course_id = $1`, courseID)
	if err != nil {
		return fmt.Errorf("clear tags: %w", err)
	}

	if len(tags) == 0 {
		return nil
	}

	query := `INSERT INTO course_tags (course_id, tag_id) VALUES ($1, $2)`
	for _, t := range tags {
		_, err := tx.Exec(ctx, query, courseID, t.ID)
		if err != nil {
			return fmt.Errorf("link tag %d: %w", t.ID, err)
		}
	}
	return nil
}

func (r *CourseRepository) getTagsByCourseID(ctx context.Context, courseID string) ([]entities.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug 
		FROM tags t
		JOIN course_tags ct ON t.id = ct.tag_id
		WHERE ct.course_id = $1
	`
	rows, err := r.pool.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []entities.Tag
	for rows.Next() {
		var t tagDTO
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, t.toEntity())
	}
	return tags, nil
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
	_, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.ModuleID,
		d.Title,
		d.ContentText,
		d.VideoURL,
		d.FileAttachmentURL,
		d.XPReward,
		d.OrderIndex,
	)
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

func (r *CourseRepository) ToggleFavorite(ctx context.Context, userID, courseID string) (bool, error) {
	existsQuery := `SELECT EXISTS(SELECT 1 FROM course_favorites WHERE user_id = $1 AND course_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, existsQuery, userID, courseID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check favorite exists: %w", err)
	}

	if exists {
		deleteQuery := `DELETE FROM course_favorites WHERE user_id = $1 AND course_id = $2`
		_, err := r.pool.Exec(ctx, deleteQuery, userID, courseID)
		if err != nil {
			return true, fmt.Errorf("remove favorite: %w", err)
		}
		return false, nil
	} else {
		insertQuery := `INSERT INTO course_favorites (user_id, course_id) VALUES ($1, $2)`
		_, err := r.pool.Exec(ctx, insertQuery, userID, courseID)
		if err != nil {
			return false, fmt.Errorf("add favorite: %w", err)
		}
		return true, nil
	}
}

func (r *CourseRepository) GetUserFavorites(ctx context.Context, userID string) ([]entities.Course, error) {
	query := `
		SELECT c.id, c.author_id, c.subject_id, c.title, c.description, 
		       c.difficulty_level, c.cover_image_url, c.is_published, c.created_at
		FROM courses c
		JOIN course_favorites cf ON c.id = cf.course_id
		WHERE cf.user_id = $1
		ORDER BY cf.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user favorites: %w", err)
	}
	defer rows.Close()

	var courses []entities.Course
	for rows.Next() {
		var d courseDTO
		err := rows.Scan(
			&d.ID, &d.AuthorID, &d.SubjectID, &d.Title, &d.Description,
			&d.DifficultyLevel, &d.CoverImageURL, &d.IsPublished, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan favorite course: %w", err)
		}

		entity := d.toEntity()
		courses = append(courses, *entity)
	}

	return courses, nil
}

// IsFavorite: Check favorite or not for button
func (r *CourseRepository) IsFavorite(ctx context.Context, userID, courseID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM course_favorites WHERE user_id = $1 AND course_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, courseID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
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
