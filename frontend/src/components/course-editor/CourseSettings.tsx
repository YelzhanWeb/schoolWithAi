import React from "react";
import type { Course, Tag } from "../../types/course";
import type { Subject } from "../../types/subject";
import { Button } from "../ui/Button";
import { Input } from "../ui/Input";

interface CourseSettingsProps {
  course: Course;
  subjects: Subject[];
  tags: Tag[];
  selectedTags: number[];
  isSaving: boolean;
  onCourseChange: (course: Course) => void;
  onTagToggle: (tagId: number) => void;
  onCoverUpload: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (e: React.FormEvent) => void;
  onPublish: () => void;
}

export const CourseSettings: React.FC<CourseSettingsProps> = ({
  course,
  subjects,
  tags,
  selectedTags,
  isSaving,
  onCourseChange,
  onTagToggle,
  onCoverUpload,
  onSubmit,
  onPublish,
}) => {
  return (
    <div className="flex-1 overflow-y-auto p-8 bg-gray-100">
      <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-sm p-8">
        <h2 className="text-2xl font-bold mb-6">Настройки курса</h2>
        <form onSubmit={onSubmit} className="space-y-6">
          {/* Обложка */}
          <div className="flex gap-6">
            <div className="w-40 h-40 bg-gray-100 rounded-lg flex-shrink-0 overflow-hidden border">
              {course.cover_image_url ? (
                <img
                  src={course.cover_image_url}
                  className="w-full h-full object-cover"
                  alt="Cover"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-gray-400">
                  Нет фото
                </div>
              )}
            </div>
            <div className="flex-1">
              <label className="block text-sm font-medium mb-2">Обложка</label>
              <input
                type="file"
                accept="image/*"
                onChange={onCoverUpload}
                className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100"
              />
            </div>
          </div>

          <Input
            label="Название"
            value={course.title}
            onChange={(e) =>
              onCourseChange({ ...course, title: e.target.value })
            }
          />

          <div>
            <label className="block text-sm font-medium mb-1">Описание</label>
            <textarea
              className="w-full h-32 p-3 border rounded-lg outline-none focus:ring-2 focus:ring-indigo-500"
              value={course.description}
              onChange={(e) =>
                onCourseChange({ ...course, description: e.target.value })
              }
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Предмет</label>
              <select
                className="w-full p-2 border rounded-lg bg-white"
                value={course.subject_id}
                onChange={(e) =>
                  onCourseChange({ ...course, subject_id: e.target.value })
                }
              >
                {subjects.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name_ru}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">
                Сложность
              </label>
              <select
                className="w-full p-2 border rounded-lg bg-white"
                value={course.difficulty_level}
                onChange={(e) =>
                  onCourseChange({
                    ...course,
                    difficulty_level: Number(e.target.value),
                  })
                }
              >
                {[1, 2, 3, 4, 5].map((l) => (
                  <option key={l} value={l}>
                    {l}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* ТЕГИ */}
          <div>
            <label className="block text-sm font-medium mb-2">Теги</label>
            <div className="flex flex-wrap gap-2">
              {tags.map((tag) => (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => onTagToggle(tag.id)}
                  className={`px-3 py-1 rounded-full text-sm font-medium transition ${
                    selectedTags.includes(tag.id)
                      ? "bg-indigo-600 text-white"
                      : "bg-gray-200 text-gray-700 hover:bg-gray-300"
                  }`}
                >
                  {tag.name}
                </button>
              ))}
            </div>
          </div>

          <div className="pt-4 border-t flex justify-between items-center">
            <button
              type="button"
              onClick={onPublish}
              className={`px-4 py-2 rounded-lg font-bold transition ${
                course.is_published
                  ? "bg-red-100 text-red-700 hover:bg-red-200"
                  : "bg-green-100 text-green-700 hover:bg-green-200"
              }`}
            >
              {course.is_published ? "Снять с публикации" : "Опубликовать курс"}
            </button>
            <Button type="submit" isLoading={isSaving} className="w-auto">
              Сохранить настройки
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
