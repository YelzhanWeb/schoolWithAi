import api from './axios';
import type { RegisterResponse, ResetPasswordRequest, LoginResponse, LoginRequest } from '../types/auth.ts';

export const authApi = {
    resetPassword: async (data: ResetPasswordRequest) => {
        await api.post('/auth/reset-password', data)
    },
    register: async (data: FormData): Promise<RegisterResponse> => {
        const response = await api.post<RegisterResponse>('/auth/register', data, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
        return response.data;
    },
    login: async (data: LoginRequest): Promise<LoginResponse> => {
        const response = await api.post<LoginResponse>('/auth/login', data);
        return response.data;
    },
};