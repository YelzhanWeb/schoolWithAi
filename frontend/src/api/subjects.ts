import api from './axios';
import type { Subject } from '../types/subject';

export const subjectsApi = {
    getAll: async (): Promise<Subject[]> => {
        const response = await api.get<{ subjects: Subject[] }>('/subjects');
        return response.data.subjects;
    }
};