import api from "./axios";
import type {
  RegisterResponse,
  ResetPasswordRequest,
  LoginResponse,
  LoginRequest,
  ChangePasswordRequest,
  ForgotPasswordRequest,
} from "../types/auth.ts";

export const authApi = {
  forgotPassword: async (email: string) => {
    await api.post("/auth/forgot-password", { email } as ForgotPasswordRequest);
  },
  resetPassword: async (data: ResetPasswordRequest) => {
    await api.post("/auth/reset-password", data);
  },
  changePassword: async (data: ChangePasswordRequest) => {
    await api.post("/auth/change-password", data);
  },
  register: async (data: FormData): Promise<RegisterResponse> => {
    const response = await api.post<RegisterResponse>("/auth/register", data, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  },
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>("/auth/login", data);
    return response.data;
  },
};
