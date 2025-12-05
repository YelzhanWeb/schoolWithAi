import React from "react";
import { Button } from "../ui/Button";
import type { Lesson } from "../../types/course";
import { Input } from "../ui/Input";
import { MarkdownEditor } from "../ui/MarkdownEditor";
import { Eye, Video, Youtube, Upload, FileText } from "lucide-react";

interface LessonEditorProps {
  lesson: Lesson;
  isSaving: boolean;
  onSave: () => void;
  onPreview: () => void;
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
  onPreview,
  onChange,
  onUpload,
}) => {
  // Instead of useEffect, we use a local state that defaults based on the lesson prop,
  // but we update it ONLY when the user clicks the buttons.
  // Ideally, we infer the "initial" mode from the data.

  // Helper to determine type based on URL
  const getVideoTypeFromUrl = (url?: string) => {
    if (url && (url.includes("youtube.com") || url.includes("youtu.be"))) {
      return "youtube";
    }
    return "upload";
  };

  const [videoType, setVideoType] = React.useState<"upload" | "youtube">(() =>
    getVideoTypeFromUrl(lesson.video_url)
  );

  // Optional: If lesson changes externally (e.g. switching between lessons), update the view
  // This is safe because it only runs when lesson.id changes.
  React.useEffect(() => {
    setVideoType(getVideoTypeFromUrl(lesson.video_url));
  }, [lesson.id, lesson.video_url]);

  return (
    <div className="max-w-4xl mx-auto bg-white rounded-xl shadow-sm p-8 space-y-8">
      {/* HEADER */}
      <div className="flex justify-between items-center border-b pb-6">
        <div className="space-y-1">
          <h2 className="text-2xl font-bold text-gray-800">Редактор урока</h2>
          <p className="text-sm text-gray-500">Настройте контент и материалы</p>
        </div>
        <div className="flex gap-3">
          <Button variant="secondary" onClick={onPreview} className="w-auto">
            <Eye size={18} className="mr-2" /> Предпросмотр
          </Button>
          <Button onClick={onSave} isLoading={isSaving} className="w-auto">
            Сохранить изменения
          </Button>
        </div>
      </div>

      {/* ОСНОВНЫЕ ПОЛЯ */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2">
          <Input
            label="Название урока"
            value={lesson.title}
            onChange={(e) => onChange({ ...lesson, title: e.target.value })}
            placeholder="Введение в..."
          />
        </div>
        <div>
          <Input
            label="Награда (XP)"
            type="number"
            value={lesson.xp_reward}
            onChange={(e) =>
              onChange({ ...lesson, xp_reward: Number(e.target.value) })
            }
          />
        </div>
      </div>

      {/* БЛОК ВИДЕО */}
      <div className="bg-gray-50 p-6 rounded-xl border border-gray-200">
        <div className="flex items-center justify-between mb-4">
          <label className="block text-sm font-bold text-gray-700 flex items-center gap-2">
            <Video size={18} /> Видеоматериал
          </label>

          {/* Переключатель источника */}
          <div className="flex bg-white rounded-lg p-1 border shadow-sm">
            <button
              onClick={() => setVideoType("upload")}
              className={`px-3 py-1 text-xs font-medium rounded-md transition flex items-center gap-2 ${
                videoType === "upload"
                  ? "bg-indigo-100 text-indigo-700"
                  : "text-gray-600 hover:bg-gray-50"
              }`}
            >
              <Upload size={14} /> Файл
            </button>
            <button
              onClick={() => setVideoType("youtube")}
              className={`px-3 py-1 text-xs font-medium rounded-md transition flex items-center gap-2 ${
                videoType === "youtube"
                  ? "bg-red-100 text-red-700"
                  : "text-gray-600 hover:bg-gray-50"
              }`}
            >
              <Youtube size={14} /> YouTube
            </button>
          </div>
        </div>

        {/* Контент видео блока */}
        <div className="space-y-4">
          {videoType === "upload" ? (
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center bg-white hover:bg-gray-50 transition">
              {lesson.video_url && !lesson.video_url.includes("youtube") ? (
                <div className="mb-4">
                  <video
                    src={lesson.video_url}
                    controls
                    className="w-full max-h-64 rounded-lg bg-black mx-auto"
                  />
                  <p className="text-xs text-green-600 mt-2">Файл загружен</p>
                </div>
              ) : null}
              <input
                type="file"
                accept="video/*"
                onChange={(e) => onUpload(e, "video_url")}
                className="block w-full text-sm text-gray-500
                        file:mr-4 file:py-2 file:px-4
                        file:rounded-full file:border-0
                        file:text-sm file:font-semibold
                        file:bg-indigo-50 file:text-indigo-700
                        hover:file:bg-indigo-100"
              />
              <p className="text-xs text-gray-400 mt-2">
                MP4, WebM (макс. 100МБ)
              </p>
            </div>
          ) : (
            <div>
              <Input
                placeholder="Вставьте ссылку на YouTube (например, https://youtu.be/...)"
                value={
                  lesson.video_url && lesson.video_url.includes("http")
                    ? lesson.video_url
                    : ""
                }
                onChange={(e) =>
                  onChange({ ...lesson, video_url: e.target.value })
                }
              />
              {lesson.video_url &&
                (lesson.video_url.includes("youtube") ||
                  lesson.video_url.includes("youtu.be")) && (
                  <div className="mt-4 aspect-video bg-black rounded-lg overflow-hidden">
                    <iframe
                      width="100%"
                      height="100%"
                      src={lesson.video_url
                        .replace("watch?v=", "embed/")
                        .replace("youtu.be/", "youtube.com/embed/")}
                      title="YouTube video player"
                      frameBorder="0"
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                      allowFullScreen
                    ></iframe>
                  </div>
                )}
            </div>
          )}
        </div>
      </div>

      {/* КОНТЕНТ (Markdown) */}
      <div>
        <label className="block text-sm font-bold text-gray-700 mb-2 flex items-center gap-2">
          <FileText size={18} /> Содержание урока
        </label>
        <MarkdownEditor
          value={lesson.content_text || ""}
          onChange={(value) => onChange({ ...lesson, content_text: value })}
          placeholder="Пишите текст урока здесь. Используйте панель сверху для форматирования и добавления картинок..."
        />
      </div>

      {/* МАТЕРИАЛЫ (Внизу) */}
      <div className="border-t pt-6">
        <label className="block text-sm font-bold text-gray-700 mb-4">
          Дополнительные материалы (Скачиваемые)
        </label>
        <div className="flex items-center gap-4 bg-gray-50 p-4 rounded-lg border">
          <div className="flex-1">
            {lesson.file_attachment_url ? (
              <a
                href={lesson.file_attachment_url}
                target="_blank"
                rel="noreferrer"
                className="text-indigo-600 text-sm hover:underline font-medium flex items-center gap-2"
              >
                📎 {lesson.file_attachment_url.split("/").pop()} (Текущий файл)
              </a>
            ) : (
              <span className="text-sm text-gray-500">Файл не выбран</span>
            )}
          </div>
          <label className="cursor-pointer">
            <span className="bg-white border border-gray-300 text-gray-700 px-3 py-1.5 rounded-md text-sm hover:bg-gray-50 transition">
              Загрузить / Заменить
            </span>
            <input
              type="file"
              onChange={(e) => onUpload(e, "file_attachment_url")}
              className="hidden"
            />
          </label>
        </div>
      </div>
    </div>
  );
};
