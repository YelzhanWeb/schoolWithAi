import React, { useState } from "react";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { coursesApi } from "../../api/courses";
import { useNavigate } from "react-router-dom";

export const CreateCoursePage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    subject_id: "math",
    difficulty_level: 1,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await coursesApi.create({
        ...formData,
        difficulty_level: Number(formData.difficulty_level),
      });
      alert("Курс создан!");
      navigate("/dashboard");
    } catch (error) {
      console.error(error);
      alert("Ошибка создания курса");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto py-10 px-4">
      <h1 className="text-3xl font-bold mb-8">Создать новый курс</h1>

      <form
        onSubmit={handleSubmit}
        className="space-y-6 bg-white p-6 rounded-xl shadow-sm"
      >
        <Input
          label="Название курса"
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          required
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Описание
          </label>
          <textarea
            className="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500 h-32"
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="ID Предмета (slug)"
            placeholder="math, physics"
            value={formData.subject_id}
            onChange={(e) =>
              setFormData({ ...formData, subject_id: e.target.value })
            }
            required
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Сложность (1-5)
            </label>
            <select
              className="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500"
              value={formData.difficulty_level}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  difficulty_level: Number(e.target.value),
                })
              }
            >
              <option value="1">1 - Легко</option>
              <option value="2">2 - Средне</option>
              <option value="3">3 - Сложно</option>
              <option value="4">4 - Хардкор</option>
              <option value="5">5 - Эксперт</option>
            </select>
          </div>
        </div>

        <Button type="submit" isLoading={isLoading}>
          Создать курс
        </Button>
      </form>
    </div>
  );
};
