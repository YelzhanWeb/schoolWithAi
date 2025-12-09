import { useNavigate } from "react-router-dom";
import { Sparkles, BarChart, User } from "lucide-react";
import type { Course } from "../../types/course";
import { Button } from "../../components/ui/Button";

interface RecommendedCourseCardProps {
  course: Course;
}

export const RecommendedCourseCard = ({
  course,
}: RecommendedCourseCardProps) => {
  const navigate = useNavigate();

  return (
    <div className="bg-white rounded-xl shadow-sm hover:shadow-md transition-shadow border border-gray-100 overflow-hidden flex flex-col h-full group">
      {/* Изображение + Бейдж AI */}
      <div className="relative h-40 overflow-hidden">
        <img
          src={course.cover_image_url || "/assets/default_course.jpg"} // Заглушка, если нет картинки
          alt={course.title}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
        />
        {/* Тот самый бейдж "AI Pick" */}
        <div className="absolute top-3 right-3 bg-gradient-to-r from-violet-600 to-indigo-600 text-white text-xs font-bold px-2 py-1 rounded-full flex items-center gap-1 shadow-lg">
          <Sparkles className="w-3 h-3" />
          AI Pick
        </div>
      </div>

      {/* Контент */}
      <div className="p-4 flex flex-col flex-grow">
        <div className="mb-2">
          <span className="text-xs font-semibold text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-md">
            {course.subject || "Общий"}
          </span>
        </div>

        <h3 className="font-bold text-gray-900 mb-2 line-clamp-2 min-h-[3rem]">
          {course.title}
        </h3>

        {/* Мета-информация */}
        <div className="mt-auto space-y-3">
          <div className="flex items-center justify-between text-xs text-gray-500">
            <div className="flex items-center gap-1">
              <BarChart className="w-3 h-3" />
              <span>Уровень {course.difficulty_level}</span>
            </div>
            {course.author && (
              <div className="flex items-center gap-1">
                <User className="w-3 h-3" />
                <span className="truncate max-w-[80px]">
                  {course.author.full_name}
                </span>
              </div>
            )}
          </div>

          <Button
            variant="outline"
            className="w-full justify-center border-indigo-200 text-indigo-700 hover:bg-indigo-50 hover:border-indigo-300"
            onClick={() => navigate(`/student/courses/${course.id}`)}
          >
            Подробнее
          </Button>
        </div>
      </div>
    </div>
  );
};
