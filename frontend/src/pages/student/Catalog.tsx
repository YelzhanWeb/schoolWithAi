import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import type { Course } from "../../types/course";
import { Button } from "../../components/ui/Button";
import { Search } from "lucide-react";

export const CatalogPage = () => {
  const [courses, setCourses] = useState<Course[]>([]);
  const [filteredCourses, setFilteredCourses] = useState<Course[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await coursesApi.getCatalog();
        setCourses(data);
        setFilteredCourses(data);
      } catch (error) {
        console.error("Ошибка загрузки каталога", error);
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, []);

  // Фильтрация при поиске
  useEffect(() => {
    const filtered = courses.filter((c) =>
      c.title.toLowerCase().includes(searchQuery.toLowerCase())
    );
    setFilteredCourses(filtered);
  }, [searchQuery, courses]);

  return (
    <div className="max-w-6xl mx-auto p-6 md:p-8">
      <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <h1 className="text-3xl font-bold text-gray-900">Каталог курсов</h1>

        {/* Поиск */}
        <div className="relative w-full md:w-96">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
            <Search size={20} />
          </div>
          <input
            type="text"
            placeholder="Поиск курса..."
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 outline-none"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      {isLoading ? (
        <div className="text-center py-20">Загрузка...</div>
      ) : filteredCourses.length === 0 ? (
        <div className="text-center py-20 text-gray-500">
          {searchQuery ? "Ничего не найдено 😔" : "Курсов пока нет"}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCourses.map((course) => (
            <Link
              key={course.id}
              to={`/student/courses/${course.id}`} // Ссылка на обзор курса (сделаем следующей)
              className="bg-white border rounded-xl overflow-hidden shadow-sm hover:shadow-lg transition-all group flex flex-col h-full"
            >
              <div className="h-48 bg-gray-100 relative overflow-hidden">
                {course.cover_image_url ? (
                  <img
                    src={course.cover_image_url}
                    alt={course.title}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-400">
                    Нет изображения
                  </div>
                )}
                {/* Бейдж сложности */}
                <div className="absolute top-2 right-2 bg-white/90 backdrop-blur px-2 py-1 rounded text-xs font-bold text-gray-700 shadow-sm">
                  Уровень {course.difficulty_level}
                </div>
              </div>

              <div className="p-5 flex-1 flex flex-col">
                <h3 className="text-lg font-bold text-gray-900 mb-2 line-clamp-2">
                  {course.title}
                </h3>
                <p className="text-sm text-gray-600 mb-4 line-clamp-3 flex-1">
                  {course.description || "Без описания"}
                </p>
                <Button variant="outline" className="w-full mt-auto">
                  Подробнее
                </Button>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
};
