import React, { useEffect, useState } from "react";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { coursesApi } from "../../api/courses";
import { subjectsApi } from "../../api/subjects";
import { uploadApi } from "../../api/upload";
import { useNavigate } from "react-router-dom";
import type { Subject } from "../../types/subject";

export const CreateCoursePage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);

  const [subjects, setSubjects] = useState<Subject[]>([]);

  const [formData, setFormData] = useState({
    title: "",
    description: "",
    subject_id: "",
    difficulty_level: 1,
    cover_image_url: "", // Ссылка на картинку
  });

  // Загружаем предметы при старте
  useEffect(() => {
    const fetchSubjects = async () => {
      try {
        const list = await subjectsApi.getAll();
        setSubjects(list);
        if (list.length > 0) {
          setFormData((prev) => ({ ...prev, subject_id: list[0].id }));
        }
      } catch (error) {
        console.error("Не удалось загрузить предметы: ", error);
      }
    };
    fetchSubjects();
  }, []);

  // Обработка загрузки обложки
  const handleCoverUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      try {
        setIsLoading(true); // Блокируем кнопку пока грузится
        const url = await uploadApi.uploadFile(e.target.files[0], "cover");
        setFormData((prev) => ({ ...prev, cover_image_url: url }));
      } catch (error) {
        alert("Ошибка загрузки изображения");
        console.error("Ошибка загрузки изображения: ", error);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await coursesApi.create({
        ...formData,
        difficulty_level: Number(formData.difficulty_level),
      });

      // После создания перекидываем в Редактор
      navigate(`/teacher/courses/${response.id}/edit`);
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
        {/* Обложка Курса */}
        <div className="flex items-center space-x-6">
          <div className="w-32 h-32 bg-gray-100 rounded-lg flex items-center justify-center overflow-hidden border-2 border-dashed border-gray-300">
            {formData.cover_image_url ? (
              <img
                src={formData.cover_image_url}
                alt="Cover"
                className="w-full h-full object-cover"
              />
            ) : (
              <span className="text-gray-400 text-xs text-center px-2">
                Нет обложки
              </span>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Обложка курса
            </label>
            <input
              type="file"
              accept="image/*"
              onChange={handleCoverUpload}
              className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100"
            />
          </div>
        </div>

        <Input
          label="Название курса"
          placeholder="Например: Основы Python"
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
            placeholder="Чему научатся студенты?"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          {/* Выбор предмета */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Предмет
            </label>
            <select
              className="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500 bg-white"
              value={formData.subject_id}
              onChange={(e) =>
                setFormData({ ...formData, subject_id: e.target.value })
              }
              required
            >
              {subjects.length === 0 && <option>Загрузка...</option>}
              {subjects.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name_ru}
                </option>
              ))}
            </select>
          </div>

          {/* Сложность */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Сложность
            </label>
            <select
              className="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500 bg-white"
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

        <div className="pt-4">
          <Button type="submit" isLoading={isLoading}>
            Создать и перейти к урокам →
          </Button>
        </div>
      </form>
    </div>
  );
};
