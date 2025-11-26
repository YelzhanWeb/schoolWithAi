import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import type { Course } from "../../api/courses";
import { uploadApi } from "../../api/upload";
import { subjectsApi } from "../../api/subjects";
import type { Module, Lesson } from "../../types/course";
import type { Subject } from "../../types/subject";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { Trash2, Settings, BookOpen, ChevronLeft } from "lucide-react"; // Иконки

export const EditCoursePage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState<"curriculum" | "settings">(
    "curriculum"
  );
  const [course, setCourse] = useState<Course | null>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);

  // Состояние для редактора урока
  const [selectedLessonId, setSelectedLessonId] = useState<string | null>(null);
  const [lessonData, setLessonData] = useState<Lesson | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (id) {
      loadCourseData();
      loadStructure();
      loadSubjects();
    }
  }, [id]);

  const loadCourseData = async () => {
    if (!id) return;
    try {
      const data = await coursesApi.getById(id);
      setCourse(data);
    } catch (e) {
      console.error("Failed to load course info", e);
    }
  };

  const loadStructure = async () => {
    if (!id) return;
    try {
      const data = await coursesApi.getStructure(id);
      setModules(data.modules || []);
    } catch (error) {
      console.error("Failed to load structure", error);
    }
  };

  const loadSubjects = async () => {
    try {
      const list = await subjectsApi.getAll();
      setSubjects(list);
    } catch (e) {
      console.error(e);
    }
  };

  // --- ЛОГИКА УРОКОВ ---
  useEffect(() => {
    const fetchLesson = async () => {
      if (!selectedLessonId) {
        setLessonData(null);
        return;
      }
      // СБРОС ДАННЫХ ПРИ ПЕРЕКЛЮЧЕНИИ (Fix бага с залипанием)
      setLessonData(null);

      try {
        const data = await coursesApi.getLesson(selectedLessonId);
        setLessonData(data);
      } catch (error) {
        console.error(error);
      }
    };
    fetchLesson();
  }, [selectedLessonId]);

  // --- ДЕЙСТВИЯ (Create / Delete) ---

  const handleAddModule = async () => {
    const title = prompt("Название модуля:");
    if (!title || !id) return;
    // Безопасное получение длины (Fix красной ошибки)
    await coursesApi.createModule(id, title, (modules?.length || 0) + 1);
    loadStructure();
  };

  const handleDeleteModule = async (moduleId: string) => {
    if (
      !confirm("Вы уверены? Это удалит модуль и ВСЕ уроки в нем безвозвратно!")
    )
      return;
    await coursesApi.deleteModule(moduleId);
    if (lessonData && lessonData.module_id === moduleId) {
      setSelectedLessonId(null); // Закрываем редактор если удалили текущий модуль
    }
    loadStructure();
  };

  const handleAddLesson = async (moduleId: string, count: number) => {
    const title = prompt("Название урока:");
    if (!title) return;
    await coursesApi.createLesson(moduleId, title, count + 1);
    loadStructure();
  };

  const handleDeleteLesson = async (lessonId: string) => {
    if (!confirm("Удалить урок?")) return;
    await coursesApi.deleteLesson(lessonId);
    if (selectedLessonId === lessonId) setSelectedLessonId(null);
    loadStructure();
  };

  // --- ЗАГРУЗКА С АВТОСОХРАНЕНИЕМ (Fix бага с URL) ---
  const handleUpload = async (
    e: React.ChangeEvent<HTMLInputElement>,
    field: "video_url" | "file_attachment_url"
  ) => {
    if (e.target.files?.[0] && lessonData) {
      try {
        setIsSaving(true);
        // 1. Грузим в MinIO
        const url = await uploadApi.uploadFile(e.target.files[0], "lesson");

        // 2. Обновляем локально
        const updated = { ...lessonData, [field]: url };
        setLessonData(updated);

        // 3. Сразу пишем в базу (Auto-save)
        await coursesApi.updateLesson(lessonData.id, {
          title: updated.title,
          content_text: updated.content_text,
          video_url: updated.video_url,
          file_attachment_url: updated.file_attachment_url,
          order_index: updated.order_index,
        });
      } catch {
        alert("Ошибка загрузки");
      } finally {
        setIsSaving(false);
      }
    }
  };

  const handleSaveLesson = async () => {
    if (!lessonData) return;
    setIsSaving(true);
    try {
      await coursesApi.updateLesson(lessonData.id, lessonData);
      alert("Сохранено");
      loadStructure();
    } catch {
      alert("Ошибка");
    } finally {
      setIsSaving(false);
    }
  };

  // --- ЛОГИКА НАСТРОЕК КУРСА ---
  const handleUpdateCourse = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!course) return;
    try {
      setIsSaving(true);
      await coursesApi.update(course.id, {
        title: course.title,
        description: course.description,
        subject_id: course.subject_id,
        difficulty_level: course.difficulty_level,
        cover_image_url: course.cover_image_url,
      });
      alert("Настройки курса обновлены");
    } catch {
      alert("Ошибка");
    } finally {
      setIsSaving(false);
    }
  };

  const handlePublish = async () => {
    if (!course) return;
    const newState = !course.is_published;
    if (
      !confirm(
        newState
          ? "Опубликовать курс для студентов?"
          : "Снять курс с публикации?"
      )
    )
      return;

    try {
      await coursesApi.publish(course.id, newState);
      setCourse({ ...course, is_published: newState });
    } catch {
      alert("Ошибка смены статуса");
    }
  };

  const handleCoverUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.[0] && course) {
      const url = await uploadApi.uploadFile(e.target.files[0], "cover");
      setCourse({ ...course, cover_image_url: url });
    }
  };

  if (!course) return <div className="p-10 text-center">Загрузка курса...</div>;

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      {/* ВЕРХНЯЯ ПАНЕЛЬ */}
      <header className="bg-white border-b px-6 py-3 flex justify-between items-center">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate("/teacher/courses")}
            className="text-gray-500 hover:text-indigo-600"
          >
            <ChevronLeft />
          </button>
          <div>
            <h1 className="font-bold text-xl text-gray-800">{course.title}</h1>
            <span
              className={`text-xs px-2 py-0.5 rounded ${
                course.is_published
                  ? "bg-green-100 text-green-800"
                  : "bg-yellow-100 text-yellow-800"
              }`}
            >
              {course.is_published ? "Опубликован" : "Черновик"}
            </span>
          </div>
        </div>

        <div className="flex space-x-2 bg-gray-100 p-1 rounded-lg">
          <button
            onClick={() => setActiveTab("curriculum")}
            className={`px-4 py-2 rounded-md text-sm font-medium transition ${
              activeTab === "curriculum"
                ? "bg-white shadow text-indigo-600"
                : "text-gray-600 hover:text-gray-900"
            }`}
          >
            <div className="flex items-center gap-2">
              <BookOpen size={16} /> Программа
            </div>
          </button>
          <button
            onClick={() => setActiveTab("settings")}
            className={`px-4 py-2 rounded-md text-sm font-medium transition ${
              activeTab === "settings"
                ? "bg-white shadow text-indigo-600"
                : "text-gray-600 hover:text-gray-900"
            }`}
          >
            <div className="flex items-center gap-2">
              <Settings size={16} /> Настройки
            </div>
          </button>
        </div>
      </header>

      {/* КОНТЕНТ (Вкладки) */}
      {activeTab === "curriculum" ? (
        <div className="flex flex-1 overflow-hidden">
          {/* ЛЕВАЯ ПАНЕЛЬ */}
          <aside className="w-80 bg-white border-r overflow-y-auto p-4 flex flex-col gap-6">
            {modules.map((module) => (
              <div
                key={module.id}
                className="border rounded-lg overflow-hidden bg-white shadow-sm group"
              >
                <div className="bg-gray-50 p-3 font-medium text-gray-700 flex justify-between items-center">
                  <span>{module.title}</span>
                  <button
                    onClick={() => handleDeleteModule(module.id)}
                    className="text-gray-400 hover:text-red-600 opacity-0 group-hover:opacity-100 transition"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
                <div className="divide-y">
                  {module.lessons?.map((lesson) => (
                    <div
                      key={lesson.id}
                      onClick={() => setSelectedLessonId(lesson.id)}
                      className={`p-3 cursor-pointer flex justify-between items-center text-sm hover:bg-indigo-50 group/lesson
                                        ${
                                          selectedLessonId === lesson.id
                                            ? "bg-indigo-50 text-indigo-700 border-l-2 border-indigo-600"
                                            : "text-gray-600"
                                        }
                                    `}
                    >
                      <span className="truncate">📄 {lesson.title}</span>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteLesson(lesson.id);
                        }}
                        className="text-gray-300 hover:text-red-500 opacity-0 group-hover/lesson:opacity-100"
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  ))}
                  {/* Fix красной ошибки: || 0 */}
                  <button
                    onClick={() =>
                      handleAddLesson(module.id, module.lessons?.length || 0)
                    }
                    className="w-full py-2 text-xs text-gray-500 hover:text-indigo-600 hover:bg-gray-50"
                  >
                    + Урок
                  </button>
                </div>
              </div>
            ))}
            <Button
              onClick={handleAddModule}
              variant="outline"
              className="w-full border-dashed"
            >
              + Новый модуль
            </Button>
          </aside>

          {/* ПРАВАЯ ПАНЕЛЬ (Редактор) */}
          <main className="flex-1 overflow-y-auto p-8 bg-gray-100">
            {lessonData ? (
              <div className="max-w-3xl mx-auto bg-white rounded-xl shadow-sm p-8 space-y-6">
                <div className="flex justify-between items-center">
                  <h2 className="text-2xl font-bold text-gray-800">
                    Редактирование урока
                  </h2>
                  <Button
                    onClick={() => handleSaveLesson()}
                    isLoading={isSaving}
                    className="w-auto"
                  >
                    Сохранить
                  </Button>
                </div>
                <Input
                  label="Название"
                  value={lessonData.title}
                  onChange={(e) =>
                    setLessonData({ ...lessonData, title: e.target.value })
                  }
                />

                {/* Видео */}
                <div className="p-4 border rounded-lg bg-gray-50">
                  <label className="block text-sm font-medium mb-2">
                    Видео
                  </label>
                  {lessonData.video_url && (
                    <video
                      src={lessonData.video_url}
                      controls
                      className="w-full h-48 bg-black rounded mb-2 object-contain"
                    />
                  )}
                  {/* Fix залипания: key с id урока */}
                  <input
                    key={`vid-${lessonData.id}`}
                    type="file"
                    accept="video/*"
                    onChange={(e) => handleUpload(e, "video_url")}
                    className="text-sm text-gray-500"
                  />
                </div>

                {/* Файлы */}
                <div className="p-4 border rounded-lg bg-gray-50">
                  <label className="block text-sm font-medium mb-2">
                    Материалы
                  </label>
                  {lessonData.file_attachment_url && (
                    <a
                      href={lessonData.file_attachment_url}
                      target="_blank"
                      className="text-indigo-600 text-sm hover:underline block mb-2"
                    >
                      📎 Скачать текущий файл
                    </a>
                  )}
                  {/* Fix залипания: key с id урока */}
                  <input
                    key={`file-${lessonData.id}`}
                    type="file"
                    onChange={(e) => handleUpload(e, "file_attachment_url")}
                    className="text-sm text-gray-500"
                  />
                </div>

                {/* Текст */}
                <textarea
                  className="w-full h-64 p-4 border rounded-lg font-mono text-sm focus:ring-2 focus:ring-indigo-500 outline-none"
                  value={lessonData.content_text || ""}
                  onChange={(e) =>
                    setLessonData({
                      ...lessonData,
                      content_text: e.target.value,
                    })
                  }
                  placeholder="# Markdown контент..."
                />
              </div>
            ) : (
              <div className="h-full flex flex-col items-center justify-center text-gray-400">
                <div className="text-6xl mb-4">👈</div>
                <p className="text-xl">Выберите урок в меню слева</p>
              </div>
            )}
          </main>
        </div>
      ) : (
        // ВКЛАДКА НАСТРОЙКИ
        <div className="flex-1 overflow-y-auto p-8 bg-gray-100">
          <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-sm p-8">
            <h2 className="text-2xl font-bold mb-6">Настройки курса</h2>
            <form onSubmit={handleUpdateCourse} className="space-y-6">
              <div className="flex gap-6">
                <div className="w-40 h-40 bg-gray-100 rounded-lg flex-shrink-0 overflow-hidden border">
                  {course.cover_image_url ? (
                    <img
                      src={course.cover_image_url}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400">
                      Нет фото
                    </div>
                  )}
                </div>
                <div className="flex-1">
                  <label className="block text-sm font-medium mb-2">
                    Обложка
                  </label>
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleCoverUpload}
                    className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100"
                  />
                </div>
              </div>

              <Input
                label="Название"
                value={course.title}
                onChange={(e) =>
                  setCourse({ ...course, title: e.target.value })
                }
              />
              <div>
                <label className="block text-sm font-medium mb-1">
                  Описание
                </label>
                <textarea
                  className="w-full h-32 p-3 border rounded-lg outline-none focus:ring-2 focus:ring-indigo-500"
                  value={course.description}
                  onChange={(e) =>
                    setCourse({ ...course, description: e.target.value })
                  }
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">
                    Предмет
                  </label>
                  <select
                    className="w-full p-2 border rounded-lg bg-white"
                    value={course.subject_id}
                    onChange={(e) =>
                      setCourse({ ...course, subject_id: e.target.value })
                    }
                  >
                    {subjects.map((s) => (
                      <option key={s.id} value={s.id}>
                        {s.name_ru}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">
                    Сложность
                  </label>
                  <select
                    className="w-full p-2 border rounded-lg bg-white"
                    value={course.difficulty_level}
                    onChange={(e) =>
                      setCourse({
                        ...course,
                        difficulty_level: Number(e.target.value),
                      })
                    }
                  >
                    {[1, 2, 3, 4, 5].map((l) => (
                      <option key={l} value={l}>
                        {l}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="pt-4 border-t flex justify-between items-center">
                <button
                  type="button"
                  onClick={handlePublish}
                  className={`px-4 py-2 rounded-lg font-bold transition ${
                    course.is_published
                      ? "bg-red-100 text-red-700 hover:bg-red-200"
                      : "bg-green-100 text-green-700 hover:bg-green-200"
                  }`}
                >
                  {course.is_published
                    ? "Снять с публикации"
                    : "Опубликовать курс"}
                </button>
                <Button type="submit" isLoading={isSaving} className="w-auto">
                  Сохранить настройки
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
