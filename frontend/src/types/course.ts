export interface Lesson {
  id: string;
  module_id: string;
  title: string;
  content_text?: string;
  video_url?: string;
  file_attachment_url?: string;
  xp_reward: number;
  order_index: number;
}

export interface Module {
  id: string;
  course_id: string;
  title: string;
  order_index: number;
  lessons: Lesson[];
}

export interface CourseStructure {
  modules: Module[];
}

export interface Course {
  id: string;
  author_id: string;
  subject_id: string;
  title: string;
  description: string;
  difficulty_level: number;
  cover_image_url?: string;
  is_published: boolean;
  tags: Tag[];
  subject?: { name_ru: string };
}

export interface CourseStructure {
  modules: Module[];
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
}
export interface CreateCourseRequest {
  title: string;
  description: string;
  subject_id: string;
  difficulty_level: number;
  cover_image_url?: string;
  tags: number[];
}

export interface CreateCourseResponse {
  id: string;
}
export interface UpdateCourseRequest {
  title: string;
  description: string;
  difficulty_level: number;
  cover_image_url: string;
  tags: number[];
}

export interface CreateLessonRequest {
  module_id: string;
  title: string;
  order_index: number;
  content_text: string;
  video_url: string;
  file_attachment_url: string;
  xp_reward: number;
}

export interface UpdateLessonRequest {
  title: string;
  order_index: number;
  content_text: string;
  video_url: string;
  file_attachment_url: string;
  xp_reward: number;
}
