import React from "react";
import { Trash2, FileText } from "lucide-react";
import { Button } from "../ui/Button";
import type { Module } from "../../types/course";

interface ModulesListProps {
  modules: Module[];
  selectedLessonId: string | null;
  onLessonSelect: (lessonId: string) => void;
  onLessonDelete: (lessonId: string) => void;
  onModuleDelete: (moduleId: string) => void;
  onLessonAdd: (moduleId: string, count: number) => void;
  onModuleAdd: () => void;
  onTestOpen: (moduleId: string) => void;
}

export const ModulesList: React.FC<ModulesListProps> = ({
  modules,
  selectedLessonId,
  onLessonSelect,
  onLessonDelete,
  onModuleDelete,
  onLessonAdd,
  onModuleAdd,
  onTestOpen,
}) => {
  return (
    <aside className="w-80 bg-white border-r overflow-y-auto p-4 flex flex-col gap-6">
      {modules.map((module) => (
        <div
          key={module.id}
          className="border rounded-lg overflow-hidden bg-white shadow-sm group"
        >
          {/* Заголовок модуля */}
          <div className="bg-gray-50 p-3 font-medium text-gray-700 flex justify-between items-center">
            <span>{module.title}</span>
            <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition">
              <button
                onClick={() => onTestOpen(module.id)}
                className="text-blue-500 hover:text-blue-700"
                title="Редактировать тест"
              >
                <FileText size={16} />
              </button>
              <button
                onClick={() => onModuleDelete(module.id)}
                className="text-gray-400 hover:text-red-600"
                title="Удалить модуль"
              >
                <Trash2 size={16} />
              </button>
            </div>
          </div>

          {/* Список уроков */}
          <div className="divide-y">
            {module.lessons?.map((lesson) => (
              <div
                key={lesson.id}
                onClick={() => onLessonSelect(lesson.id)}
                className={`p-3 cursor-pointer flex justify-between items-center text-sm hover:bg-indigo-50 group/lesson
                  ${
                    selectedLessonId === lesson.id
                      ? "bg-indigo-50 text-indigo-700 border-l-2 border-indigo-600"
                      : "text-gray-600"
                  }
                `}
              >
                <span className="truncate">📄 {lesson.title}</span>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onLessonDelete(lesson.id);
                  }}
                  className="text-gray-300 hover:text-red-500 opacity-0 group-hover/lesson:opacity-100"
                >
                  <Trash2 size={14} />
                </button>
              </div>
            ))}
            <button
              onClick={() =>
                onLessonAdd(module.id, module.lessons?.length || 0)
              }
              className="w-full py-2 text-xs text-gray-500 hover:text-indigo-600 hover:bg-gray-50"
            >
              + Урок
            </button>
          </div>
        </div>
      ))}

      <Button
        onClick={onModuleAdd}
        variant="outline"
        className="w-full border-dashed"
      >
        + Новый модуль
      </Button>
    </aside>
  );
};
