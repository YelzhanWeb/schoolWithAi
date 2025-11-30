import React from "react";
import { Button } from "../ui/Button";
import type { Lesson } from "../../types/course";
import { Input } from "../ui/Input";
import { MarkdownEditor } from "../ui/MarkdownEditor";

interface LessonEditorProps {
  lesson: Lesson;
  isSaving: boolean;
  onSave: () => void;
  onChange: (lesson: Lesson) => void;
  onUpload: (
    e: React.ChangeEvent<HTMLInputElement>,
    field: "video_url" | "file_attachment_url"
  ) => void;
}

export const LessonEditor: React.FC<LessonEditorProps> = ({
  lesson,
  isSaving,
  onSave,
  onChange,
  onUpload,
}) => {
  return (
    <div className="max-w-3xl mx-auto bg-white rounded-xl shadow-sm p-8 space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-gray-800">
          Редактирование урока
        </h2>
        <Button onClick={onSave} isLoading={isSaving} className="w-auto">
          Сохранить
        </Button>
      </div>

      <Input
        label="Название"
        value={lesson.title}
        onChange={(e) => onChange({ ...lesson, title: e.target.value })}
      />

      <Input
        label="XP награда"
        type="number"
        value={lesson.xp_reward}
        onChange={(e) =>
          onChange({ ...lesson, xp_reward: Number(e.target.value) })
        }
      />

      {/* Видео */}
      <div className="p-4 border rounded-lg bg-gray-50">
        <label className="block text-sm font-medium mb-2">Видео</label>
        {lesson.video_url && (
          <video
            src={lesson.video_url}
            controls
            className="w-full h-48 bg-black rounded mb-2 object-contain"
          />
        )}
        <input
          key={`vid-${lesson.id}`}
          type="file"
          accept="video/*"
          onChange={(e) => onUpload(e, "video_url")}
          className="text-sm text-gray-500"
        />
      </div>

      {/* Файлы */}
      <div className="p-4 border rounded-lg bg-gray-50">
        <label className="block text-sm font-medium mb-2">Материалы</label>
        {lesson.file_attachment_url && (
          <a
            href={lesson.file_attachment_url}
            target="_blank"
            className="text-indigo-600 text-sm hover:underline block mb-2"
          >
            📎 Скачать текущий файл
          </a>
        )}
        <input
          key={`file-${lesson.id}`}
          type="file"
          onChange={(e) => onUpload(e, "file_attachment_url")}
          className="text-sm text-gray-500"
        />
      </div>

      {/* Текст контента - НОВЫЙ РЕДАКТОР */}
      <div>
        <label className="block text-sm font-medium mb-2 text-gray-700">
          Содержание урока
        </label>
        <MarkdownEditor
          value={lesson.content_text || ""}
          onChange={(value) => onChange({ ...lesson, content_text: value })}
          placeholder="Начните создавать контент урока..."
        />
      </div>
    </div>
  );
};
