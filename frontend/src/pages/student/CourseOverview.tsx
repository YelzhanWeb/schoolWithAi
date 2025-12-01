import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import { studentApi } from "../../api/student";
import type { Course, Module } from "../../types/course";
import { Button } from "../../components/ui/Button";
import {
  PlayCircle,
  CheckCircle,
  Globe,
  Clock,
  Award,
  BarChart,
  ChevronDown,
  ChevronUp,
  Share2,
  Heart,
  FileText,
  Video,
  User,
} from "lucide-react";

export const CourseOverview = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [course, setCourse] = useState<Course | null>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [completedLessons, setCompletedLessons] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Состояние для сворачивания/разворачивания модулей
  const [expandedModules, setExpandedModules] = useState<
    Record<string, boolean>
  >({});

  useEffect(() => {
    if (!id) return;

    const loadData = async () => {
      try {
        setIsLoading(true);
        const [courseData, structureData, progressData] = await Promise.all([
          coursesApi.getById(id),
          coursesApi.getStructure(id),
          studentApi.getCourseProgress(id),
        ]);

        setCourse(courseData);
        setModules(structureData.modules || []);
        setCompletedLessons(progressData);

        // По умолчанию раскрываем первый модуль
        if (structureData.modules?.[0]) {
          setExpandedModules({ [structureData.modules[0].id]: true });
        }
      } catch (error) {
        console.error("Failed to load course", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [id]);

  const toggleModule = (moduleId: string) => {
    setExpandedModules((prev) => ({
      ...prev,
      [moduleId]: !prev[moduleId],
    }));
  };

  const getNextLessonId = () => {
    if (!modules.length) return null;
    for (const m of modules) {
      if (!m.lessons) continue;
      for (const l of m.lessons) {
        if (!completedLessons.includes(l.id)) {
          return l.id;
        }
      }
    }
    return modules[0]?.lessons?.[0]?.id;
  };

  const handleStart = () => {
    const nextId = getNextLessonId();
    if (nextId) {
      navigate(`/student/courses/${id}/lessons/${nextId}`);
    }
  };

  if (isLoading)
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  if (!course) return <div className="p-10 text-center">Курс не найден</div>;

  const totalLessons = modules.reduce(
    (acc, m) => acc + (m.lessons?.length || 0),
    0
  );
  const progressPercent = Math.round(
    (completedLessons.length / (totalLessons || 1)) * 100
  );
  const nextLessonId = getNextLessonId();

  return (
    <div className="min-h-screen bg-gray-50 font-sans">
      {/* --- HERO SECTION (Черный фон) --- */}
      <div className="bg-[#1C1D1F] text-white py-12">
        <div className="max-w-6xl mx-auto px-6 flex flex-col md:flex-row gap-8">
          <div className="md:w-2/3 space-y-4">
            {/* Breadcrumbs */}
            <div className="text-indigo-300 text-sm font-medium mb-2 flex items-center gap-2">
              <span
                onClick={() => navigate("/student/catalog")}
                className="cursor-pointer hover:underline"
              >
                Каталог
              </span>
              <span>/</span>
              <span>
                {course.difficulty_level === 1 ? "Начинающим" : "Продвинутым"}
              </span>
            </div>

            <h1 className="text-3xl md:text-4xl font-bold leading-tight">
              {course.title}
            </h1>
            <p className="text-lg text-gray-300 line-clamp-3">
              {course.description}
            </p>

            <div className="flex items-center gap-4 text-sm pt-4">
              {course.difficulty_level && (
                <div className="flex items-center gap-1 bg-yellow-500/20 text-yellow-400 px-2 py-1 rounded">
                  <BarChart size={16} />
                  <span>Уровень {course.difficulty_level}</span>
                </div>
              )}
              <div className="flex items-center gap-1 text-gray-300">
                <Globe size={16} />
                <span>Русский</span>
              </div>
              <div className="flex items-center gap-1 text-gray-300">
                <Clock size={16} />
                <span>
                  Последнее обновление: {new Date().toLocaleDateString()}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* --- MAIN CONTENT & SIDEBAR --- */}
      <div className="max-w-6xl mx-auto px-6 py-8 relative">
        <div className="flex flex-col md:flex-row gap-8">
          {/* LEFT COLUMN (Content) */}
          <div className="md:w-2/3 space-y-8">
            {/* What you'll learn (Заглушка, т.к. в базе нет поля) */}
            <div className="bg-white border border-gray-200 p-6 rounded-xl">
              <h2 className="text-xl font-bold mb-4 text-gray-900">
                Чему вы научитесь
              </h2>
              <ul className="grid grid-cols-1 md:grid-cols-2 gap-3">
                <li className="flex gap-2 text-sm text-gray-700">
                  <CheckCircle
                    size={18}
                    className="text-gray-400 flex-shrink-0"
                  />
                  <span>Освоите фундаментальные принципы темы</span>
                </li>
                <li className="flex gap-2 text-sm text-gray-700">
                  <CheckCircle
                    size={18}
                    className="text-gray-400 flex-shrink-0"
                  />
                  <span>Научитесь применять знания на практике</span>
                </li>
                <li className="flex gap-2 text-sm text-gray-700">
                  <CheckCircle
                    size={18}
                    className="text-gray-400 flex-shrink-0"
                  />
                  <span>Решите реальные задачи и кейсы</span>
                </li>
                <li className="flex gap-2 text-sm text-gray-700">
                  <CheckCircle
                    size={18}
                    className="text-gray-400 flex-shrink-0"
                  />
                  <span>Получите сертификат об окончании</span>
                </li>
              </ul>
            </div>

            {/* Curriculum */}
            <div>
              <h2 className="text-2xl font-bold mb-4 text-gray-900">
                Программа курса
              </h2>
              <div className="text-sm text-gray-500 mb-4">
                {modules.length} модулей • {totalLessons} уроков
              </div>

              <div className="space-y-4">
                {modules.map((module) => (
                  <div
                    key={module.id}
                    className="bg-white border border-gray-200 rounded-xl overflow-hidden"
                  >
                    {/* Module Header */}
                    <button
                      onClick={() => toggleModule(module.id)}
                      className="w-full flex items-center justify-between p-4 bg-gray-50 hover:bg-gray-100 transition text-left"
                    >
                      <div className="flex items-center gap-3">
                        {expandedModules[module.id] ? (
                          <ChevronUp size={20} className="text-gray-500" />
                        ) : (
                          <ChevronDown size={20} className="text-gray-500" />
                        )}
                        <h3 className="font-bold text-gray-800">
                          {module.title}
                        </h3>
                      </div>
                      <span className="text-xs text-gray-500">
                        {module.lessons?.length || 0} уроков
                      </span>
                    </button>

                    {/* Lessons List */}
                    {expandedModules[module.id] && (
                      <div className="divide-y divide-gray-100">
                        {module.lessons?.map((lesson) => {
                          const isCompleted = completedLessons.includes(
                            lesson.id
                          );
                          const isCurrent = lesson.id === nextLessonId;

                          return (
                            <div
                              key={lesson.id}
                              onClick={() =>
                                navigate(
                                  `/student/courses/${id}/lessons/${lesson.id}`
                                )
                              }
                              className={`p-4 flex items-center justify-between cursor-pointer transition hover:bg-indigo-50 ${
                                isCurrent ? "bg-indigo-50" : "bg-white"
                              }`}
                            >
                              <div className="flex items-center gap-3">
                                {isCompleted ? (
                                  <CheckCircle
                                    className="text-green-500 fill-green-50"
                                    size={18}
                                  />
                                ) : lesson.video_url ? (
                                  <Video className="text-gray-400" size={18} />
                                ) : (
                                  <FileText
                                    className="text-gray-400"
                                    size={18}
                                  />
                                )}
                                <span
                                  className={`text-sm ${
                                    isCompleted
                                      ? "text-gray-500 line-through"
                                      : "text-gray-700"
                                  }`}
                                >
                                  {lesson.title}
                                </span>
                              </div>

                              {lesson.video_url && (
                                <span className="text-xs text-gray-400">
                                  Видео
                                </span>
                              )}
                            </div>
                          );
                        })}
                        {(!module.lessons || module.lessons.length === 0) && (
                          <div className="p-4 text-sm text-gray-400 italic text-center">
                            Нет уроков
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {/* Instructor */}
            <div>
              <h2 className="text-2xl font-bold mb-4 text-gray-900">
                Преподаватель
              </h2>
              <div className="flex items-start gap-4">
                <div className="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                  <User size={32} className="text-gray-400" />
                </div>
                <div>
                  <div className="font-bold text-indigo-600 text-lg">
                    ID: {course.author_id}
                  </div>
                  <p className="text-gray-500 text-sm mb-2">Автор курса</p>
                  <p className="text-gray-600 text-sm leading-relaxed">
                    Опытный преподаватель, создавший этот курс специально для
                    платформы OqysAI.
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* RIGHT COLUMN (Sticky Sidebar) */}
          <div className="md:w-1/3 relative">
            <div className="sticky top-6 space-y-6">
              {/* Course Card */}
              <div className="bg-white border border-gray-200 rounded-xl shadow-lg overflow-hidden md:-mt-48 z-10 relative">
                {/* Video/Image Preview */}
                <div
                  className="h-48 bg-gray-100 relative group cursor-pointer"
                  onClick={handleStart}
                >
                  {course.cover_image_url ? (
                    <img
                      src={course.cover_image_url}
                      alt="Cover"
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-800">
                      <PlayCircle size={48} className="text-white opacity-80" />
                    </div>
                  )}
                  <div className="absolute inset-0 bg-black/20 group-hover:bg-black/40 transition flex items-center justify-center">
                    <PlayCircle
                      size={64}
                      className="text-white drop-shadow-lg scale-95 group-hover:scale-100 transition-transform"
                    />
                  </div>
                </div>

                <div className="p-6">
                  <div className="text-3xl font-bold text-gray-900 mb-4">
                    Бесплатно
                  </div>

                  <Button
                    onClick={handleStart}
                    className="mb-4 w-full font-bold text-lg py-4"
                  >
                    {completedLessons.length === 0
                      ? "Начать обучение"
                      : "Продолжить"}
                  </Button>

                  {completedLessons.length > 0 && (
                    <div className="mb-6">
                      <div className="flex justify-between text-xs text-gray-500 mb-1">
                        <span>Прогресс</span>
                        <span>{progressPercent}%</span>
                      </div>
                      <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-green-500"
                          style={{ width: `${progressPercent}%` }}
                        />
                      </div>
                    </div>
                  )}

                  <div className="space-y-3 pt-4 border-t border-gray-100">
                    <div className="flex items-center justify-between text-sm text-gray-600">
                      <span className="flex items-center gap-2">
                        <Award size={16} /> Сертификат
                      </span>
                      <span>Да</span>
                    </div>
                    <div className="flex items-center justify-between text-sm text-gray-600">
                      <span className="flex items-center gap-2">
                        <Globe size={16} /> Доступ
                      </span>
                      <span>Навсегда</span>
                    </div>
                  </div>

                  <div className="flex gap-2 mt-6">
                    <button className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 flex items-center justify-center gap-2">
                      <Heart size={16} /> В избранное
                    </button>
                    <button className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 flex items-center justify-center gap-2">
                      <Share2 size={16} /> Поделиться
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
