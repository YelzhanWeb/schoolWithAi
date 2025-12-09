import { useEffect, useState } from "react";
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
} from "lucide-react";
import { testsApi } from "../../api/tests";

export const CourseOverview = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [course, setCourse] = useState<Course | null>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [completedLessons, setCompletedLessons] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isFavorite, setIsFavorite] = useState(false);
  const [modulesWithTests, setModulesWithTests] = useState<Set<string>>(
    new Set()
  );

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
        setIsFavorite(courseData.is_favorite || false); // Ставим состояние из ответа
        setModules(structureData.modules || []);
        setCompletedLessons(progressData);

        const testChecks = await Promise.all(
          (structureData.modules || []).map(async (module) => {
            try {
              await testsApi.getByModuleId(module.id);
              return module.id;
            } catch {
              return null;
            }
          })
        );

        const modulesWithTestsSet = new Set(
          testChecks.filter((id): id is string => id !== null)
        );
        setModulesWithTests(modulesWithTestsSet);

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

  // --- ЛОГИКА КНОПОК ---

  const handleToggleFavorite = async () => {
    if (!course) return;
    try {
      // Оптимистичное обновление интерфейса
      setIsFavorite(!isFavorite);
      await coursesApi.toggleFavorite(course.id);
    } catch (error) {
      console.error(error);
      setIsFavorite(isFavorite); // Возвращаем как было при ошибке
    }
  };

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href);
    alert("Ссылка на курс скопирована в буфер обмена!");
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
    } else {
      alert("В этом курсе пока нет уроков.");
    }
  };

  if (isLoading)
    return (
      <div className="min-h-screen flex items-center justify-center">
        Loading...
      </div>
    );
  if (!course) return <div className="p-10 text-center">Курс не найден</div>;

  const totalLessons = modules.reduce(
    (acc, m) => acc + (m.lessons?.length || 0),
    0
  );
  const progressPercent =
    totalLessons > 0
      ? Math.round((completedLessons.length / totalLessons) * 100)
      : 0;

  return (
    <div className="min-h-screen bg-gray-50 font-sans">
      {/* HEADER */}
      <div className="bg-[#1C1D1F] text-white py-12">
        <div className="max-w-6xl mx-auto px-6">
          <div className="md:w-2/3 space-y-4">
            <h1 className="text-3xl md:text-4xl font-bold leading-tight">
              {course.title}
            </h1>
            <p className="text-lg text-gray-300 line-clamp-3">
              {course.description}
            </p>

            <div className="flex items-center gap-4 text-sm pt-4 flex-wrap">
              <div className="flex items-center gap-1 bg-yellow-500/20 text-yellow-400 px-2 py-1 rounded">
                <BarChart size={16} />
                <span>Уровень {course.difficulty_level}</span>
              </div>
              <div className="flex items-center gap-1 text-gray-300">
                <Globe size={16} />
                <span>Русский</span>
              </div>
              <div className="flex items-center gap-1 text-gray-300">
                <Clock size={16} />
                {/* Исправлено: берем реальную дату с бэкенда */}
                <span>Обновлено: {course.updated_at || "Недавно"}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-6xl mx-auto px-6 py-8 relative">
        <div className="flex flex-col md:flex-row gap-8">
          {/* LEFT COLUMN */}
          <div className="md:w-2/3 space-y-8">
            {/* БЛОК "Чему вы научитесь" УДАЛЕН ПО ПРОСЬБЕ */}

            {/* Программа курса */}
            <div>
              <h2 className="text-2xl font-bold mb-4 text-gray-900">
                Программа курса
              </h2>
              <div className="space-y-4">
                {modules.map((module) => (
                  <div
                    key={module.id}
                    className="bg-white border border-gray-200 rounded-xl overflow-hidden"
                  >
                    <button
                      onClick={() => toggleModule(module.id)}
                      className="w-full flex items-center justify-between p-4 bg-gray-50 hover:bg-gray-100 transition text-left"
                    >
                      <div className="flex items-center gap-3">
                        {expandedModules[module.id] ? (
                          <ChevronUp size={20} />
                        ) : (
                          <ChevronDown size={20} />
                        )}
                        <h3 className="font-bold text-gray-800">
                          {module.title}
                        </h3>
                      </div>
                      <span className="text-xs text-gray-500">
                        {module.lessons?.length || 0} уроков
                      </span>
                    </button>

                    {expandedModules[module.id] && (
                      <div className="divide-y divide-gray-100">
                        {module.lessons?.map((lesson) => {
                          const isCompleted = completedLessons.includes(
                            lesson.id
                          );
                          return (
                            <div
                              key={lesson.id}
                              // Клик по уроку тоже ведет в плеер
                              onClick={() =>
                                navigate(
                                  `/student/courses/${id}/lessons/${lesson.id}`
                                )
                              }
                              className="p-4 flex items-center justify-between cursor-pointer hover:bg-indigo-50 transition"
                            >
                              <div className="flex items-center gap-3">
                                {isCompleted ? (
                                  <CheckCircle
                                    className="text-green-500"
                                    size={18}
                                  />
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
                            </div>
                          );
                        })}
                        {/* ПОКАЗЫВАЕМ КНОПКУ ТЕСТА, ТОЛЬКО ЕСЛИ ОН ЕСТЬ */}
                        {modulesWithTests.has(module.id) && (
                          <div
                            onClick={() =>
                              navigate(
                                `/student/courses/${id}/modules/${module.id}/test`
                              )
                            }
                            className="p-4 flex items-center gap-3 cursor-pointer hover:bg-purple-50 transition border-t-2 border-purple-100"
                          >
                            <div className="w-5 h-5 rounded-full border-2 border-purple-500 flex items-center justify-center">
                              <div className="w-2.5 h-2.5 bg-purple-500 rounded-full"></div>
                            </div>
                            <span className="text-sm font-medium text-purple-700">
                              Итоговый тест модуля
                            </span>
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {/* ПРЕПОДАВАТЕЛЬ (ИСПРАВЛЕНО) */}
            <div>
              <h2 className="text-2xl font-bold mb-4 text-gray-900">
                Преподаватель
              </h2>
              <div className="flex items-center gap-4 bg-white p-6 rounded-xl border border-gray-100 shadow-sm">
                {course.author ? (
                  <>
                    <img
                      src={
                        course.author.avatar_url ||
                        "https://ui-avatars.com/api/?name=" +
                          course.author.full_name
                      }
                      alt={course.author.full_name}
                      className="w-16 h-16 rounded-full object-cover border-2 border-indigo-100"
                    />
                    <div>
                      <div className="font-bold text-indigo-600 text-lg">
                        {course.author.full_name}
                      </div>
                      <p className="text-gray-500 text-sm">Автор курса</p>
                    </div>
                  </>
                ) : (
                  <div className="text-gray-500">
                    Информация об авторе скрыта
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* RIGHT SIDEBAR (Sticky) */}
          <div className="md:w-1/3 relative">
            <div className="sticky top-6 space-y-6">
              <div className="bg-white border border-gray-200 rounded-xl shadow-lg overflow-hidden md:-mt-48 z-10 relative">
                {/* Картинка / Кнопка старт */}
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
                    <div className="w-full h-full flex items-center justify-center bg-gray-800">
                      <PlayCircle size={48} className="text-white opacity-80" />
                    </div>
                  )}
                  {/* Оверлей при наведении */}
                  <div className="absolute inset-0 bg-black/20 group-hover:bg-black/40 transition flex items-center justify-center">
                    <div className="bg-white/20 backdrop-blur-sm p-4 rounded-full border border-white/30">
                      <PlayCircle
                        size={48}
                        className="text-white drop-shadow-lg"
                        fill="currentColor"
                      />
                    </div>
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

                  {/* КНОПКИ ДЕЙСТВИЙ */}
                  <div className="flex gap-2 mt-6">
                    <button
                      onClick={handleToggleFavorite}
                      className={`flex-1 py-2 border rounded-lg text-sm font-medium flex items-center justify-center gap-2 transition-colors ${
                        isFavorite
                          ? "bg-red-50 border-red-200 text-red-600"
                          : "border-gray-300 text-gray-700 hover:bg-gray-50"
                      }`}
                    >
                      <Heart
                        size={16}
                        fill={isFavorite ? "currentColor" : "none"}
                      />
                      {isFavorite ? "В избранном" : "В избранное"}
                    </button>

                    <button
                      onClick={handleShare}
                      className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 flex items-center justify-center gap-2"
                    >
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
