import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import type { Course } from "../../types/course";
import { Button } from "../../components/ui/Button";

const Dashboard = () => {
  const [courses, setCourses] = useState<Course[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchCourses = async () => {
      try {
        const data = await coursesApi.getMyCourses();
        setCourses(data);
      } catch (error) {
        console.error("Ошибка загрузки курсов:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchCourses();
  }, []);

  if (isLoading) {
    return (
      <div className="p-8 text-center text-gray-500">Загрузка курсов...</div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto p-6">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Мои курсы</h1>
        <Link to="/teacher/create-course">
          <Button>+ Создать курс</Button>
        </Link>
      </div>

      {courses.length === 0 ? (
        <div className="text-center py-20 bg-gray-50 rounded-xl border-2 border-dashed border-gray-200">
          <p className="text-xl text-gray-500 mb-4">У вас пока нет курсов</p>
          <Link to="/teacher/create-course">
            <Button variant="outline">Создать первый курс</Button>
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course) => (
            <div
              key={course.id}
              className="bg-white border rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow flex flex-col"
            >
              {/* Обложка */}
              <div className="h-40 bg-gray-100 relative">
                {course.cover_image_url ? (
                  <img
                    src={course.cover_image_url}
                    alt={course.title}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-100">
                    Нет изображения
                  </div>
                )}
                <div className="absolute top-2 right-2">
                  <span
                    className={`px-2 py-1 rounded text-xs font-semibold ${
                      course.is_published
                        ? "bg-green-100 text-green-800"
                        : "bg-yellow-100 text-yellow-800"
                    }`}
                  >
                    {course.is_published ? "Опубликован" : "Черновик"}
                  </span>
                </div>
              </div>

              {/* Контент */}
              <div className="p-5 flex-1 flex flex-col">
                <h3 className="text-lg font-bold text-gray-900 mb-2 line-clamp-1">
                  {course.title}
                </h3>
                <p className="text-sm text-gray-600 mb-4 line-clamp-2 flex-1">
                  {course.description || "Нет описания"}
                </p>

                <div className="flex items-center justify-between mt-auto pt-4 border-t">
                  <span className="text-xs text-gray-500">
                    Уровень: {course.difficulty_level}/5
                  </span>
                  <Link to={`/teacher/courses/${course.id}/edit`}>
                    <Button variant="outline">Редактировать</Button>
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Dashboard;
