export type UserRole = 'student' | 'teacher' | 'admin';

export interface User {
    id: string;
    email: string;
    role: UserRole;
    first_name: string;
    last_name: string;
    avatar_url?: string;
}

export interface LoginResponse {
    token: string;
}

export interface RegisterResponse {
    id: string;
    email: string;
    role: string;
    full_name: string;
    last_name: string;
}

export interface ResetPasswordRequest {
    email: string;
    old_password: string;
    new_password: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}