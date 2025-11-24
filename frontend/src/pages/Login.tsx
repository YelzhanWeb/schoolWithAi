import React, { useState } from "react";
import { authApi } from "../api/auth";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import { Link, useNavigate } from "react-router-dom";
import { AxiosError } from "axios";

export const LoginPage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    if (error) setError("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const response = await authApi.login(formData);

      // Сохраняем токен, который пришел с бэкенда
      localStorage.setItem("token", response.token);

      // Можно сделать запрос за профилем юзера, но пока просто редирект
      // alert('Вход выполнен!');
      navigate("/dashboard"); // Перекидываем на главную (создадим позже)
    } catch (err) {
      const error = err as AxiosError<{ message: string }>;
      console.error(error);
      setError(error.response?.data?.message || "Неверный email или пароль");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 bg-white p-8 rounded-xl shadow-lg">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Войти в аккаунт
          </h2>
        </div>

        {error && (
          <div className="bg-red-100 text-red-700 p-3 rounded-lg text-sm text-center border border-red-200">
            {error}
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <Input
              name="email"
              type="email"
              placeholder="Email"
              required
              value={formData.email}
              onChange={handleChange}
            />

            <Input
              name="password"
              type="password"
              placeholder="Пароль"
              required
              value={formData.password}
              onChange={handleChange}
            />
          </div>

          <div className="flex items-center justify-end">
            <Link
              to="/change-password"
              className="text-sm font-medium text-indigo-600 hover:text-indigo-500"
            >
              Забыли пароль?
            </Link>
          </div>

          <Button type="submit" isLoading={isLoading}>
            Войти
          </Button>

          <div className="text-center mt-4">
            <Link
              to="/register"
              className="text-sm text-gray-600 hover:text-gray-900"
            >
              Нет аккаунта?{" "}
              <span className="text-indigo-600 font-medium">
                Зарегистрироваться
              </span>
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
};
