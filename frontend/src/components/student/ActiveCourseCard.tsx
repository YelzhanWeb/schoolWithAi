import React from "react";
import { Play } from "lucide-react";
import { Link } from "react-router-dom";
import type { ActiveCourse } from "../../api/student";
import { Button } from "../ui/Button";

interface Props {
  course: ActiveCourse;
}

export const ActiveCourseCard: React.FC<Props> = ({ course }) => {
  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-md transition-shadow flex flex-col h-full">
      <div className="h-32 bg-gray-100 relative">
        {course.cover_url ? (
          <img
            src={course.cover_url}
            alt={course.title}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-400 text-sm">
            Нет обложки
          </div>
        )}
        <div className="absolute inset-0 bg-black bg-opacity-10" />
      </div>

      <div className="p-5 flex-1 flex flex-col">
        <h3 className="font-bold text-lg text-gray-800 mb-2 line-clamp-1">
          {course.title}
        </h3>

        {/* Прогресс бар */}
        <div className="mt-auto">
          <div className="flex justify-between text-xs text-gray-500 mb-1">
            <span>Прогресс</span>
            <span>{course.progress_percentage}%</span>
          </div>
          <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden mb-4">
            <div
              className="h-full bg-green-500 rounded-full transition-all duration-500"
              style={{ width: `${course.progress_percentage}%` }}
            />
          </div>

          <Link to={`/student/courses/${course.course_id}`}>
            <Button className="w-full flex items-center justify-center gap-2">
              <Play size={16} fill="currentColor" /> Продолжить
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
};
