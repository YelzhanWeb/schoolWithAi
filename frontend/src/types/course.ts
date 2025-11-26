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