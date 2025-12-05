import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import { studentApi } from "../../api/student";
import type { Lesson, Module } from "../../types/course";
import { Button } from "../../components/ui/Button";
import {
  ChevronLeft,
  CheckCircle,
  Circle,
  FileText,
  Download,
} from "lucide-react";
import { testsApi } from "../../api/tests";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";

// Если у вас нет компонента для просмотра Markdown (MarkdownViewer),
// можно использовать тот же MDXEditor в режиме readOnly или просто отображать HTML.
// Для простоты пока используем div с стилями.

export const LessonPlayer = () => {
  const { courseId, lessonId } = useParams<{
    courseId: string;
    lessonId: string;
  }>();
  const navigate = useNavigate();

  // Структура для меню (боковая панель)
  const [modules, setModules] = useState<Module[]>([]);

  // Полные данные текущего урока (с контентом)
  const [currentLesson, setCurrentLesson] = useState<Lesson | null>(null);

  const [completedLessons, setCompletedLessons] = useState<string[]>([]);
  const [isLoadingStructure, setIsLoadingStructure] = useState(true);
  const [isLoadingLesson, setIsLoadingLesson] = useState(false);
  const [isCompleting, setIsCompleting] = useState(false);

  // 1. Загружаем структуру курса (один раз при входе)
  useEffect(() => {
    if (!courseId) return;
    const loadStructure = async () => {
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
        setIsLoadingStructure(false);
      }
    };
    loadStructure();
  }, [courseId]);

  // 2. Загружаем КОНТЕНТ конкретного урока (каждый раз при смене lessonId)
  useEffect(() => {
    if (!lessonId) return;

    const loadLessonContent = async () => {
      setIsLoadingLesson(true);
      try {
        // ВОТ ЗДЕСЬ мы берем полные данные с контентом
        const lessonData = await coursesApi.getLesson(lessonId);
        setCurrentLesson(lessonData);
      } catch (error) {
        console.error("Не удалось загрузить урок", error);
      } finally {
        setIsLoadingLesson(false);
      }
    };

    loadLessonContent();
  }, [lessonId]);

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

      // 3. Показываем XP
      if (res.xp_gained > 0) {
        // Можно заменить на красивый тост
        alert(`Урок пройден! Вы получили +${res.xp_gained} XP 🔥`);
      }

      // 4. Ищем следующий урок
      let nextLessonId = null;
      let foundCurrent = false;

      // Проходим по всем модулям и урокам по порядку
      for (const m of modules) {
        if (!m.lessons) continue;
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

  if (isLoadingStructure)
    return <div className="p-10 text-center">Загрузка структуры курса...</div>;

  return (
    <div className="flex h-screen bg-white overflow-hidden">
      {/* SIDEBAR (Список уроков) */}
      <aside className="w-80 border-r bg-gray-50 flex-col hidden md:flex h-full">
        <div className="p-4 border-b bg-white flex items-center gap-2 flex-shrink-0">
          <Button
            variant="secondary"
            onClick={() => navigate(`/student/courses/${courseId}`)}
            className="p-2 w-auto"
          >
            <ChevronLeft size={20} />
          </Button>
          <span className="font-bold text-gray-700 truncate">Содержание</span>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-6">
          {modules.map((module) => (
            <div key={module.id}>
              <h3 className="font-bold text-xs text-gray-500 mb-2 uppercase tracking-wider">
                {module.title}
              </h3>
              <div className="space-y-1">
                {module.lessons &&
                  module.lessons.map((lesson) => {
                    const isCompleted = completedLessons.includes(lesson.id);
                    const isActive = lesson.id === lessonId; // Сравниваем с URL
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
                            ? "bg-indigo-100 text-indigo-700 font-medium"
                            : "hover:bg-gray-100 text-gray-700"
                        }`}
                      >
                        {isCompleted ? (
                          <CheckCircle
                            size={16}
                            className="text-green-500 mt-0.5 flex-shrink-0"
                          />
                        ) : (
                          <Circle
                            size={16}
                            className={`mt-0.5 flex-shrink-0 ${
                              isActive ? "text-indigo-500" : "text-gray-300"
                            }`}
                          />
                        )}
                        <span className="line-clamp-2">{lesson.title}</span>
                      </button>
                    );
                  })}
              </div>
            </div>
          ))}
        </div>
      </aside>

      {/* MAIN CONTENT */}
      <main className="flex-1 flex flex-col h-full overflow-hidden relative">
        {/* Header mobile */}
        <div className="md:hidden p-4 border-b flex items-center gap-2 flex-shrink-0 bg-white">
          <Button
            variant="secondary"
            onClick={() => navigate(`/student/courses/${courseId}`)}
            className="w-auto p-2"
          >
            <ChevronLeft size={20} />
          </Button>
          <span className="font-bold truncate text-sm">
            {currentLesson?.title || "Загрузка..."}
          </span>
        </div>

        <div className="flex-1 overflow-y-auto p-6 md:p-12 w-full max-w-4xl mx-auto">
          {isLoadingLesson ? (
            <div className="flex justify-center items-center h-64">
              <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
            </div>
          ) : currentLesson ? (
            <>
              <h1 className="text-3xl font-bold text-gray-900 mb-8">
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

              {/* CONTENT (Markdown Render) */}
              <div className="prose prose-indigo max-w-none text-gray-700 leading-relaxed mb-10">
                <ReactMarkdown
                  remarkPlugins={[remarkGfm]}
                  rehypePlugins={[rehypeRaw]}
                  components={{
                    // Настройка ссылок, чтобы открывались в новой вкладке (опционально)
                    a: ({ ...props }) => (
                      <a
                        {...props}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 hover:underline"
                      />
                    ),
                    img: ({ ...props }) => (
                      <img
                        {...props}
                        className="rounded-xl shadow-sm max-w-full h-auto my-4"
                      />
                    ),
                  }}
                >
                  {currentLesson.content_text}
                </ReactMarkdown>
              </div>

              {/* ATTACHMENTS */}
              {currentLesson.file_attachment_url && (
                <div className="mb-12 p-4 border rounded-xl bg-gray-50 flex items-center gap-4 hover:bg-gray-100 transition">
                  <div className="w-12 h-12 bg-white rounded-lg flex items-center justify-center shadow-sm text-indigo-600">
                    <FileText size={24} />
                  </div>
                  <div className="flex-1">
                    <div className="font-medium text-gray-900">
                      Материалы к уроку
                    </div>
                    <div className="text-sm text-gray-500">
                      Дополнительный файл
                    </div>
                  </div>
                  <a
                    href={currentLesson.file_attachment_url}
                    target="_blank"
                    rel="noreferrer"
                    className="p-2 bg-white rounded-full border border-gray-200 text-gray-600 hover:text-indigo-600 hover:border-indigo-600 transition"
                    title="Скачать"
                  >
                    <Download size={20} />
                  </a>
                </div>
              )}

              {/* FOOTER ACTIONS */}
              <div className="border-t pt-8 pb-20 flex flex-col sm:flex-row justify-between items-center gap-4">
                <div className="text-sm text-gray-500">
                  Награда за урок:{" "}
                  <span className="font-bold text-orange-500 bg-orange-50 px-2 py-1 rounded-full ml-1">
                    +{currentLesson.xp_reward} XP
                  </span>
                </div>

                <Button
                  onClick={handleComplete}
                  isLoading={isCompleting}
                  // Если уже пройден, кнопка всё равно активна, чтобы перейти дальше, но меняется текст/стиль
                  variant={
                    completedLessons.includes(currentLesson.id)
                      ? "outline"
                      : "primary"
                  }
                  className="w-full sm:w-auto px-8 py-3 text-lg"
                >
                  {completedLessons.includes(currentLesson.id)
                    ? "Урок пройден (Следующий →)"
                    : "Завершить урок"}
                </Button>
              </div>
            </>
          ) : (
            <div className="text-center py-20 text-gray-500">
              Урок не найден
            </div>
          )}
        </div>
      </main>
    </div>
  );
};
