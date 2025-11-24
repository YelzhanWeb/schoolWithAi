import React, { useState } from "react";
import { authApi } from "../api/auth";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import { Link, useNavigate } from "react-router-dom";
import { AxiosError } from "axios";

export const RegisterPage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState({
    email: "",
    password: "",
    firstName: "",
    lastName: "",
    role: "student",
  });
  const [avatar, setAvatar] = useState<File | null>(null);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setAvatar(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const data = new FormData();
      data.append("email", formData.email);
      data.append("password", formData.password);
      data.append("first_name", formData.firstName);
      data.append("last_name", formData.lastName);
      data.append("role", formData.role);

      if (avatar) {
        data.append("avatar", avatar);
      }

      await authApi.register(data);

      alert("Регистрация успешна! Теперь войдите.");
      navigate("/login");
    } catch (err) {
      const error = err as AxiosError<{ message: string }>;

      console.error(error);
      setError(error.response?.data?.message || "Ошибка регистрации");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 bg-white p-8 rounded-xl shadow-lg">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Создать аккаунт
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Или{" "}
            <Link
              to="/login"
              className="font-medium text-indigo-600 hover:text-indigo-500"
            >
              войти в существующий
            </Link>
          </p>
        </div>

        {error && (
          <div className="bg-red-100 text-red-700 p-3 rounded-lg text-sm text-center">
            {error}
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input
                name="firstName"
                placeholder="Имя"
                required
                onChange={handleChange}
              />
              <Input
                name="lastName"
                placeholder="Фамилия"
                required
                onChange={handleChange}
              />
            </div>

            <Input
              name="email"
              type="email"
              placeholder="Email"
              required
              onChange={handleChange}
            />

            <Input
              name="password"
              type="password"
              placeholder="Пароль (мин 8 символов)"
              required
              onChange={handleChange}
            />

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Роль
              </label>
              <select
                name="role"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 outline-none"
                onChange={handleChange}
              >
                <option value="student">Ученик</option>
                <option value="teacher">Учитель</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Аватарка (необязательно)
              </label>
              <input
                type="file"
                accept="image/*"
                onChange={handleFileChange}
                className="block w-full text-sm text-gray-500
                                file:mr-4 file:py-2 file:px-4
                                file:rounded-full file:border-0
                                file:text-sm file:font-semibold
                                file:bg-indigo-50 file:text-indigo-700
                                hover:file:bg-indigo-100"
              />
            </div>
          </div>

          <Button type="submit" isLoading={isLoading}>
            Зарегистрироваться
          </Button>
        </form>
      </div>
    </div>
  );
};
