import React from "react";
import { useNavigate } from "react-router-dom";
import { ChevronLeft, BookOpen, Settings, Trash2 } from "lucide-react";

type Tab = "curriculum" | "settings";

interface CourseHeaderProps {
  courseTitle: string;
  isPublished: boolean;
  activeTab: Tab;
  onTabChange: (tab: Tab) => void;
  onDelete?: () => void;
}

export const CourseHeader: React.FC<CourseHeaderProps> = ({
  courseTitle,
  isPublished,
  activeTab,
  onTabChange,
  onDelete,
}) => {
  const navigate = useNavigate();

  return (
    <header className="bg-white border-b px-6 py-3 flex justify-between items-center">
      <div className="flex items-center space-x-4">
        <button
          onClick={() => navigate("/teacher/courses")}
          className="text-gray-500 hover:text-indigo-600"
        >
          <ChevronLeft />
        </button>
        <div>
          <h1 className="font-bold text-xl text-gray-800">{courseTitle}</h1>
          <span
            className={`text-xs px-2 py-0.5 rounded ${
              isPublished
                ? "bg-green-100 text-green-800"
                : "bg-yellow-100 text-yellow-800"
            }`}
          >
            {isPublished ? "Опубликован" : "Черновик"}
          </span>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {/* Табы */}
        <div className="flex space-x-2 bg-gray-100 p-1 rounded-lg">
          <button
            onClick={() => onTabChange("curriculum")}
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
            onClick={() => onTabChange("settings")}
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

        {/* Кнопка удаления */}
        {onDelete && (
          <button
            onClick={onDelete}
            className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"
            title="Удалить курс"
          >
            <Trash2 size={20} />
          </button>
        )}
      </div>
    </header>
  );
};
