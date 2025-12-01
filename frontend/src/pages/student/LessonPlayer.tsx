import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import { studentApi } from "../../api/student";
import type { Lesson, Module } from "../../types/course";
import { Button } from "../../components/ui/Button";
import { ChevronLeft, CheckCircle, Circle, FileText } from "lucide-react";
import { testsApi } from "../../api/tests";

export const LessonPlayer = () => {
  const { courseId, lessonId } = useParams<{
    courseId: string;
    lessonId: string;
  }>();
  const navigate = useNavigate();

  const [modules, setModules] = useState<Module[]>([]);
  const [currentLesson, setCurrentLesson] = useState<Lesson | null>(null);
  const [completedLessons, setCompletedLessons] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCompleting, setIsCompleting] = useState(false);

  useEffect(() => {
    if (!courseId) return;
    const loadData = async () => {
      try {
        const [structData, progressData] = await Promise.all([
          coursesApi.getStructure(courseId),
          studentApi.getCourseProgress(courseId),
        ]);
        setModules(structData.modules || []);
        setCompletedLessons(progressData);
      } catch (e) {
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, [courseId]);

  // При смене lessonId в URL обновляем текущий урок
  useEffect(() => {
    if (modules.length > 0 && lessonId) {
      let found = null;
      modules.forEach((m) =>
        m.lessons.forEach((l) => {
          if (l.id === lessonId) found = l;
        })
      );
      setCurrentLesson(found);
    }
  }, [lessonId, modules]);

  const handleComplete = async () => {
    if (!currentLesson) return;
    setIsCompleting(true);
    try {
      // 1. Отправляем запрос на завершение
      const res = await testsApi.completeLesson(currentLesson.id);

      // 2. Обновляем локальный стейт (добавляем галочку)
      if (!completedLessons.includes(currentLesson.id)) {
        setCompletedLessons([...completedLessons, currentLesson.id]);
      }

      // 3. Показываем XP (можно сделать красивый тост/модалку)
      if (res.xp_gained > 0) {
        alert(`Урок пройден! Вы получили +${res.xp_gained} XP 🔥`);
      }

      // 4. Ищем следующий урок
      // (Простая логика поиска следующего по списку)
      let nextLessonId = null;
      let foundCurrent = false;

      for (const m of modules) {
        for (const l of m.lessons) {
          if (foundCurrent) {
            nextLessonId = l.id;
            break;
          }
          if (l.id === currentLesson.id) foundCurrent = true;
        }
        if (nextLessonId) break;
      }

      if (nextLessonId) {
        navigate(`/student/courses/${courseId}/lessons/${nextLessonId}`);
      } else {
        alert("Курс завершен! Поздравляем!");
        navigate(`/student/courses/${courseId}`);
      }
    } catch (e) {
      alert("Ошибка при завершении урока");
      console.log(e);
    } finally {
      setIsCompleting(false);
    }
  };

  if (isLoading || !currentLesson)
    return <div className="p-10 text-center">Загрузка...</div>;

  return (
    <div className="flex h-screen bg-white">
      {/* SIDEBAR (Список уроков) */}
      <aside className="w-80 border-r bg-gray-50 flex flex-col hidden md:flex">
        <div className="p-4 border-b bg-white flex items-center gap-2">
          <Button onClick={() => navigate(`/student/courses/${courseId}`)}>
            <ChevronLeft size={20} />
          </Button>
          <span className="font-bold text-gray-700">Содержание курса</span>
        </div>
        <div className="flex-1 overflow-y-auto p-4 space-y-6">
          {modules.map((module) => (
            <div key={module.id}>
              <h3 className="font-bold text-sm text-gray-500 mb-2 uppercase tracking-wider">
                {module.title}
              </h3>
              <div className="space-y-1">
                {module.lessons.map((lesson) => {
                  const isCompleted = completedLessons.includes(lesson.id);
                  const isActive = lesson.id === currentLesson.id;
                  return (
                    <button
                      key={lesson.id}
                      onClick={() =>
                        navigate(
                          `/student/courses/${courseId}/lessons/${lesson.id}`
                        )
                      }
                      className={`w-full text-left p-3 rounded-lg text-sm flex items-start gap-3 transition-colors ${
                        isActive
                          ? "bg-indigo-100 text-indigo-700"
                          : "hover:bg-gray-100 text-gray-700"
                      }`}
                    >
                      {isCompleted ? (
                        <CheckCircle
                          size={18}
                          className="text-green-500 mt-0.5 flex-shrink-0"
                        />
                      ) : (
                        <Circle
                          size={18}
                          className="text-gray-300 mt-0.5 flex-shrink-0"
                        />
                      )}
                      <span className={isActive ? "font-semibold" : ""}>
                        {lesson.title}
                      </span>
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </aside>

      {/* MAIN CONTENT */}
      <main className="flex-1 flex flex-col overflow-hidden">
        {/* Header mobile */}
        <div className="md:hidden p-4 border-b flex items-center gap-2">
          <Button onClick={() => navigate(`/student/courses/${courseId}`)}>
            <ChevronLeft /> Назад
          </Button>
          <span className="font-bold truncate">{currentLesson.title}</span>
        </div>

        <div className="flex-1 overflow-y-auto p-6 md:p-10 max-w-4xl mx-auto w-full">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">
            {currentLesson.title}
          </h1>

          {/* VIDEO PLAYER */}
          {currentLesson.video_url && (
            <div className="mb-8 aspect-video bg-black rounded-xl overflow-hidden shadow-lg">
              <video
                src={currentLesson.video_url}
                controls
                className="w-full h-full"
              />
            </div>
          )}

          {/* CONTENT */}
          <div className="prose prose-indigo max-w-none text-gray-700 leading-relaxed mb-10 whitespace-pre-wrap">
            {currentLesson.content_text}
          </div>

          {/* ATTACHMENTS */}
          {currentLesson.file_attachment_url && (
            <div className="mb-10 p-4 border rounded-xl bg-gray-50 flex items-center gap-4">
              <div className="w-10 h-10 bg-white rounded-lg flex items-center justify-center shadow-sm">
                <FileText className="text-indigo-600" />
              </div>
              <div>
                <div className="font-medium">Материалы к уроку</div>
                <a
                  href={currentLesson.file_attachment_url}
                  target="_blank"
                  className="text-sm text-indigo-600 hover:underline"
                >
                  Скачать файл
                </a>
              </div>
            </div>
          )}

          {/* FOOTER ACTIONS */}
          <div className="border-t pt-8 flex justify-between items-center">
            <div className="text-sm text-gray-500">
              Награда:{" "}
              <span className="font-bold text-orange-500">
                +{currentLesson.xp_reward} XP
              </span>
            </div>
            <Button
              onClick={handleComplete}
              isLoading={isCompleting}
              disabled={completedLessons.includes(currentLesson.id)}
              className={
                completedLessons.includes(currentLesson.id)
                  ? "bg-green-600 hover:bg-green-700"
                  : ""
              }
            >
              {completedLessons.includes(currentLesson.id)
                ? "Урок пройден ✅"
                : "Завершить и продолжить →"}
            </Button>
          </div>
        </div>
      </main>
    </div>
  );
};
