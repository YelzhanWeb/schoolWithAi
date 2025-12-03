import React, { useState } from "react";
import { authApi } from "../api/auth";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import { Link, useNavigate } from "react-router-dom";
import { AxiosError } from "axios";

export const ChangePasswordPage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const [formData, setFormData] = useState({
    email: "",
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    if (error) setError("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.newPassword !== formData.confirmPassword) {
      setError("Новые пароли не совпадают!");
      return;
    }

    if (formData.newPassword.length < 8) {
      setError("Минимальная длина пароля — 8 символов");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      await authApi.changePassword({
        email: formData.email,
        old_password: formData.oldPassword,
        new_password: formData.newPassword,
      });

      setSuccess(true);
      setTimeout(() => navigate("/login"), 2000);
    } catch (err) {
      const error = err as AxiosError<{ message: string }>;
      console.error(error);
      setError(error.response?.data?.message || "Ошибка смены пароля");
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
            Сейчас вы будете перенаправлены на страницу входа.
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
            Смена пароля
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Введите старый пароль для подтверждения личности
          </p>
        </div>

        {error && (
          <div className="bg-red-100 text-red-700 p-3 rounded-lg text-sm text-center border border-red-200">
            {error}
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <Input
              label="Ваш Email"
              name="email"
              type="email"
              required
              value={formData.email}
              onChange={handleChange}
            />

            <Input
              label="Старый пароль"
              name="oldPassword"
              type="password"
              required
              value={formData.oldPassword}
              onChange={handleChange}
            />

            <div className="border-t pt-4 mt-4">
              <Input
                label="Новый пароль"
                name="newPassword"
                type="password"
                required
                minLength={8}
                value={formData.newPassword}
                onChange={handleChange}
              />
            </div>

            <Input
              label="Подтвердите новый пароль"
              name="confirmPassword"
              type="password"
              required
              error={
                formData.confirmPassword &&
                formData.newPassword !== formData.confirmPassword
                  ? "Пароли не совпадают"
                  : ""
              }
              value={formData.confirmPassword}
              onChange={handleChange}
            />
          </div>

          <div className="flex flex-col gap-3">
            <Button type="submit" isLoading={isLoading}>
              Сменить пароль
            </Button>
            <Link to="/login">
              <Button variant="secondary" type="button">
                Отмена
              </Button>
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
};
