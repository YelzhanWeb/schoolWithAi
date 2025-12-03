import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { studentApi, type ActiveCourse } from "../../api/student";
import { coursesApi } from "../../api/courses";
import type { Course } from "../../types/course";
import { ActiveCourseCard } from "../../components/student/ActiveCourseCard"; // Используем существующую карточку
import { BookOpen, Heart, Loader2 } from "lucide-react";
import { Button } from "../../components/ui/Button";

type Tab = "active" | "favorites";

export const MyLearningPage = () => {
  const [activeTab, setActiveTab] = useState<Tab>("active");
  const [isLoading, setIsLoading] = useState(true);

  // Данные
  const [activeCourses, setActiveCourses] = useState<ActiveCourse[]>([]);
  const [favoriteCourses, setFavoriteCourses] = useState<Course[]>([]);

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async () => {
    setIsLoading(true);
    try {
      if (activeTab === "active") {
        const data = await studentApi.getMyCourses();
        setActiveCourses(data);
      } else {
        const data = await coursesApi.getFavorites();
        setFavoriteCourses(data);
      }
    } catch (error) {
      console.error("Ошибка загрузки:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-6xl mx-auto p-6 md:p-8 min-h-screen">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Мое обучение</h1>

      {/* TABS */}
      <div className="flex border-b border-gray-200 mb-8">
        <button
          onClick={() => setActiveTab("active")}
          className={`flex items-center gap-2 px-6 py-3 font-medium text-sm transition-all border-b-2 ${
            activeTab === "active"
              ? "border-indigo-600 text-indigo-600"
              : "border-transparent text-gray-500 hover:text-gray-700"
          }`}
        >
          <BookOpen size={18} /> Текущие курсы
        </button>
        <button
          onClick={() => setActiveTab("favorites")}
          className={`flex items-center gap-2 px-6 py-3 font-medium text-sm transition-all border-b-2 ${
            activeTab === "favorites"
              ? "border-indigo-600 text-indigo-600"
              : "border-transparent text-gray-500 hover:text-gray-700"
          }`}
        >
          <Heart size={18} /> Избранное
        </button>
      </div>

      {/* CONTENT */}
      {isLoading ? (
        <div className="flex justify-center py-20">
          <Loader2 className="animate-spin text-indigo-600" size={32} />
        </div>
      ) : (
        <div>
          {/* Вкладка: Текущие */}
          {activeTab === "active" && (
            <>
              {activeCourses.length > 0 ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {activeCourses.map((course) => (
                    <ActiveCourseCard key={course.course_id} course={course} />
                  ))}
                </div>
              ) : (
                <EmptyState
                  message="Вы еще не начали ни одного курса."
                  actionText="Перейти в каталог"
                  link="/student/catalog"
                />
              )}
            </>
          )}

          {/* Вкладка: Избранное */}
          {activeTab === "favorites" && (
            <>
              {favoriteCourses.length > 0 ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {favoriteCourses.map((course) => (
                    <Link
                      key={course.id}
                      to={`/student/courses/${course.id}`}
                      className="bg-white border rounded-xl overflow-hidden shadow-sm hover:shadow-lg transition-all group flex flex-col h-full"
                    >
                      <div className="h-44 bg-gray-100 relative">
                        {course.cover_image_url ? (
                          <img
                            src={course.cover_image_url}
                            alt={course.title}
                            className="w-full h-full object-cover"
                          />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center text-gray-400">
                            Нет фото
                          </div>
                        )}
                        <div className="absolute top-2 right-2 bg-white/90 px-2 py-1 rounded text-xs font-bold shadow-sm">
                          Lvl {course.difficulty_level}
                        </div>
                      </div>
                      <div className="p-5 flex-1 flex flex-col">
                        <h3 className="font-bold text-gray-900 mb-2 line-clamp-2">
                          {course.title}
                        </h3>
                        <Button variant="outline" className="mt-auto w-full">
                          Подробнее
                        </Button>
                      </div>
                    </Link>
                  ))}
                </div>
              ) : (
                <EmptyState
                  message="В избранном пока пусто."
                  actionText="Найти интересное"
                  link="/student/catalog"
                />
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
};

const EmptyState = ({
  message,
  actionText,
  link,
}: {
  message: string;
  actionText: string;
  link: string;
}) => (
  <div className="text-center py-20 bg-white rounded-xl border border-dashed border-gray-200">
    <p className="text-gray-500 text-lg mb-4">{message}</p>
    <Link to={link}>
      <Button className="w-auto">{actionText}</Button>
    </Link>
  </div>
);
