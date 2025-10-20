-- Индексы для Collaborative Filtering
CREATE INDEX IF NOT EXISTS idx_resource_ratings_student_resource ON resource_ratings (student_id, resource_id);

CREATE INDEX IF NOT EXISTS idx_resource_ratings_rating ON resource_ratings (rating)
WHERE
    rating >= 4;

-- Индексы для Content-Based Filtering
CREATE INDEX IF NOT EXISTS idx_resource_tags_tag_resource ON resource_tags (tag_id, resource_id);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags (name);

-- Индексы для Knowledge-Based Filtering
CREATE INDEX IF NOT EXISTS idx_student_skills_student_level ON student_skills (student_id, proficiency_level);

CREATE INDEX IF NOT EXISTS idx_student_skills_name ON student_skills (skill_name);

-- Индексы для прогресса
CREATE INDEX IF NOT EXISTS idx_student_progress_student_status ON student_progress (student_id, status);

CREATE INDEX IF NOT EXISTS idx_student_progress_completed ON student_progress (student_id)
WHERE
    status = 'completed';

-- Составной индекс для модулей и ресурсов
CREATE INDEX IF NOT EXISTS idx_modules_course_order ON modules (course_id, order_index);

CREATE INDEX IF NOT EXISTS idx_resources_module_type ON resources (module_id, resource_type);