// frontend/src/pages/ResetPassword.tsx
import React, { useState, useEffect } from "react";
import { authApi } from "../api/auth";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import { useNavigate, useSearchParams } from "react-router-dom";
import { AxiosError } from "axios";

export const ResetPasswordPage = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [formData, setFormData] = useState({
    email: "",
    code: "",
    newPassword: "",
    confirmPassword: "",
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  // Подставляем email из URL, если он там есть
  useEffect(() => {
    const emailFromUrl = searchParams.get("email");
    if (emailFromUrl) {
      setFormData((prev) => ({ ...prev, email: emailFromUrl }));
    }
  }, [searchParams]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.newPassword !== formData.confirmPassword) {
      setError("Пароли не совпадают!");
      return;
    }
    if (formData.newPassword.length < 8) {
      setError("Пароль должен быть не менее 8 символов");
      return;
    }

    setIsLoading(true);

    try {
      await authApi.resetPassword({
        email: formData.email,
        code: formData.code,
        new_password: formData.newPassword,
      });

      setSuccess(true);
      setTimeout(() => navigate("/login"), 3000);
    } catch (err) {
      const error = err as AxiosError<{ message: string }>;
      setError(
        error.response?.data?.message || "Неверный код или ошибка сервера"
      );
    } finally {
      setIsLoading(false);
    }
  };

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
        <div className="max-w-md w-full bg-white p-8 rounded-xl shadow-lg text-center">
          <div className="text-5xl mb-4">✅</div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            Пароль изменен!
          </h2>
          <p className="text-gray-600">
            Вы будете перенаправлены на страницу входа...
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 bg-white p-8 rounded-xl shadow-lg">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Новый пароль
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Введите код из письма и придумайте новый пароль
          </p>
        </div>

        {error && (
          <div className="bg-red-100 text-red-700 p-3 rounded-lg text-sm text-center">
            {error}
          </div>
        )}

        <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
          <Input
            label="Email"
            name="email"
            type="email"
            required
            value={formData.email}
            onChange={handleChange}
          />

          <Input
            label="Код из письма (6 цифр)"
            name="code"
            type="text"
            required
            placeholder="123456"
            value={formData.code}
            onChange={handleChange}
            maxLength={6}
            className="tracking-widest text-center font-mono text-lg"
          />

          <div className="pt-2">
            <Input
              label="Новый пароль"
              name="newPassword"
              type="password"
              required
              value={formData.newPassword}
              onChange={handleChange}
            />
          </div>

          <Input
            label="Подтверждение пароля"
            name="confirmPassword"
            type="password"
            required
            value={formData.confirmPassword}
            onChange={handleChange}
          />

          <Button type="submit" isLoading={isLoading}>
            Сменить пароль
          </Button>
        </form>
      </div>
    </div>
  );
};
