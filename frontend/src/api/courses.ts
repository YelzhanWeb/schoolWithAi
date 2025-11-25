import api from './axios';

export interface CreateCourseRequest {
    title: string;
    description: string;
    subject_id: string; // Пока будем вводить вручную, позже сделаем выпадающий список
    difficulty_level: number;
}

export interface CreateCourseResponse {
    id: string;
}

export const coursesApi = {
    create: async (data: CreateCourseRequest): Promise<CreateCourseResponse> => {
        const response = await api.post<CreateCourseResponse>('/courses', data);
        return response.data;
    }
};