import api from './axios';
import type { CourseStructure, Lesson } from '../types/course';

export interface Course {
    id: string;
    title: string;
    description: string;
    subject_id: string;
    difficulty_level: number;
    cover_image_url?: string;
    is_published: boolean;
}
export interface CreateCourseRequest {
    title: string;
    description: string;
    subject_id: string; 
    difficulty_level: number;
    cover_image_url?: string;
}

export interface CreateCourseResponse {
    id: string;
}

export const coursesApi = {
    getMyCourses: async (): Promise<Course[]> => {
        const response = await api.get<{courses: Course[]}>('/teacher/courses'); 
        return response.data.courses;
    },

    getById: async (id: string): Promise<Course> => {
        const response = await api.get<Course>(`/courses/${id}`);
        return response.data;
    },

    create: async (data: CreateCourseRequest): Promise<CreateCourseResponse> => {
        const response = await api.post<CreateCourseResponse>('/courses', data);
        return response.data;
    },

    update: async (id: string, data: Partial<CreateCourseRequest>) => {
        await api.put(`/courses/${id}`, data);
    },

    publish: async (id: string, isPublished: boolean) => {
        await api.post(`/courses/${id}/publish`, { is_published: isPublished });
    },

    getStructure: async (courseId: string): Promise<CourseStructure> => {
        const response = await api.get<CourseStructure>(`/courses/${courseId}/structure`);
        return response.data;
    },

    createModule: async (courseId: string, title: string, orderIndex: number): Promise<{id: string}> => {
        const response = await api.post('/modules', { 
            course_id: courseId, 
            title, 
            order_index: orderIndex 
        });
        return response.data;
    },

    updateModule: async (moduleId: string, title: string, orderIndex: number) => {
        await api.put(`/modules/${moduleId}`, { title, order_index: orderIndex });
    },

    deleteModule: async (moduleId: string) => {
        await api.delete(`/modules/${moduleId}`);
    },

    createLesson: async (moduleId: string, title: string, orderIndex: number): Promise<{id: string}> => {
        const response = await api.post('/lessons', { 
            module_id: moduleId, 
            title, 
            order_index: orderIndex 
        });
        return response.data;
    },

    getLesson: async (lessonId: string): Promise<Lesson> => {
        const response = await api.get<Lesson>(`/lessons/${lessonId}`);
        return response.data;
    },

    updateLesson: async (lessonId: string, data: Partial<Lesson>) => {
        await api.put(`/lessons/${lessonId}`, data);
    },

    deleteLesson: async (lessonId: string) => {
        await api.delete(`/lessons/${lessonId}`);
    }
};