import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import { uploadApi } from "../../api/upload";
import { subjectsApi } from "../../api/subjects";
import { testsApi, type Test } from "../../api/tests";
import type {
  Module,
  Lesson,
  Tag,
  Course,
  CreateLessonRequest,
} from "../../types/course";
import type { Subject } from "../../types/subject";

import { CourseHeader } from "../../components/course-editor/CourseHeader";
import { ModulesList } from "../../components/course-editor/ModulesList";
import { LessonEditor } from "../../components/course-editor/LessonEditor";
import { TestEditor } from "../../components/course-editor/TestEditor";
import { CourseSettings } from "../../components/course-editor/CourseSettings";
import { EmptyState } from "../../components/course-editor/EmptyState";
import { X } from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

type Tab = "curriculum" | "settings";

export const EditCoursePage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState<Tab>("curriculum");
  const [course, setCourse] = useState<Course | null>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedTags, setSelectedTags] = useState<number[]>([]);

  const [selectedLessonId, setSelectedLessonId] = useState<string | null>(null);
  const [lessonData, setLessonData] = useState<Lesson | null>(null);

  const [selectedModuleForTest, setSelectedModuleForTest] = useState<
    string | null
  >(null);
  const [testData, setTestData] = useState<Test | null>(null);
  const [isCreatingTest, setIsCreatingTest] = useState(false);

  const [isSaving, setIsSaving] = useState(false);

  const [showPreview, setShowPreview] = useState(false);
  useEffect(() => {
    if (id) {
      loadCourseData();
      loadStructure();
      loadSubjects();
      loadTags();
    }
  }, [id]);

  const loadCourseData = async () => {
    if (!id) return;
    try {
      const data = await coursesApi.getById(id);
      setCourse(data);
      if (data.tags) {
        setSelectedTags(data.tags.map((t: Tag) => t.id));
      }
    } catch (e) {
      console.error("Failed to load course info", e);
    }
  };

  const loadStructure = async () => {
    if (!id) return;
    try {
      const data = await coursesApi.getStructure(id);
      setModules(data.modules || []);
    } catch (error) {
      console.error("Failed to load structure", error);
    }
  };

  const loadSubjects = async () => {
    try {
      const list = await subjectsApi.getAll();
      setSubjects(list);
    } catch (e) {
      console.error(e);
    }
  };

  const loadTags = async () => {
    try {
      const list = await coursesApi.getAllTags();
      setTags(list);
    } catch (e) {
      console.error(e);
    }
  };

  // === ЛОГИКА УРОКОВ ===
  useEffect(() => {
    const fetchLesson = async () => {
      if (!selectedLessonId) {
        setLessonData(null);
        return;
      }
      setLessonData(null);
      try {
        const data = await coursesApi.getLesson(selectedLessonId);
        setLessonData(data);
      } catch (error) {
        console.error(error);
      }
    };
    fetchLesson();
  }, [selectedLessonId]);

  const handleAddModule = async () => {
    const title = prompt("Название модуля:");
    if (!title || !id) return;
    await coursesApi.createModule(id, title, (modules?.length || 0) + 1);
    loadStructure();
  };

  const handleDeleteModule = async (moduleId: string) => {
    if (
      !confirm("Вы уверены? Это удалит модуль и ВСЕ уроки в нем безвозвратно!")
    )
      return;
    await coursesApi.deleteModule(moduleId);
    if (lessonData && lessonData.module_id === moduleId) {
      setSelectedLessonId(null);
    }
    loadStructure();
  };

  const handleAddLesson = async (moduleId: string, count: number) => {
    const title = prompt("Название урока:");
    if (!title) return;

    const lessonData: CreateLessonRequest = {
      module_id: moduleId,
      title: title,
      order_index: count + 1,
      content_text: "",
      video_url: "",
      file_attachment_url: "",
      xp_reward: 0,
    };

    try {
      await coursesApi.createLesson(lessonData);
      loadStructure();
    } catch (error) {
      console.error("Ошибка создания урока:", error);
      alert("Не удалось создать урок");
    }
  };

  const handleDeleteLesson = async (lessonId: string) => {
    if (!confirm("Удалить урок?")) return;
    await coursesApi.deleteLesson(lessonId);
    if (selectedLessonId === lessonId) setSelectedLessonId(null);
    loadStructure();
  };

  const handleUpload = async (
    e: React.ChangeEvent<HTMLInputElement>,
    field: "video_url" | "file_attachment_url"
  ) => {
    if (e.target.files?.[0] && lessonData) {
      try {
        setIsSaving(true);
        const url = await uploadApi.uploadFile(e.target.files[0], "lesson");
        const updated = { ...lessonData, [field]: url };
        setLessonData(updated);
        await coursesApi.updateLesson(lessonData.id, updated);
      } catch {
        alert("Ошибка загрузки");
      } finally {
        setIsSaving(false);
      }
    }
  };

  const handleSaveLesson = async () => {
    if (!lessonData) return;
    setIsSaving(true);
    try {
      await coursesApi.updateLesson(lessonData.id, lessonData);
      alert("Сохранено");
      loadStructure();
    } catch {
      alert("Ошибка");
    } finally {
      setIsSaving(false);
    }
  };

  // === ЛОГИКА ТЕСТОВ ===
  const handleOpenTestEditor = async (moduleId: string) => {
    setSelectedModuleForTest(moduleId);
    setSelectedLessonId(null);
    setIsCreatingTest(false);

    try {
      const test = await testsApi.getByModuleIdWithAnswers(moduleId);
      setTestData(test);
    } catch {
      setTestData(null);
    }
  };

  const handleCreateTest = () => {
    setIsCreatingTest(true);
    setTestData({
      test_id: "",
      title: "Новый тест",
      passing_score: 70,
      questions: [
        {
          text: "",
          question_type: "single_choice",
          answers: [
            { text: "", is_correct: false },
            { text: "", is_correct: false },
          ],
        },
      ],
    });
  };

  const handleSaveTest = async () => {
    if (!testData || !selectedModuleForTest) return;
    setIsSaving(true);
    try {
      await testsApi.create({
        module_id: selectedModuleForTest,
        title: testData.title,
        passing_score: testData.passing_score,
        questions: testData.questions,
      });
      alert("Тест создан!");
      setIsCreatingTest(false);
      handleOpenTestEditor(selectedModuleForTest);
    } catch (error) {
      console.error(error);
      alert("Ошибка создания теста");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeleteTest = async () => {
    if (!testData || !confirm("Удалить тест безвозвратно?")) return;
    try {
      await testsApi.delete(testData.test_id);
      alert("Тест удален");
      setTestData(null);
      setSelectedModuleForTest(null);
    } catch (error) {
      console.error(error);
      alert("Ошибка удаления теста");
    }
  };

  const addQuestion = () => {
    if (!testData) return;
    setTestData({
      ...testData,
      questions: [
        ...testData.questions,
        {
          text: "",
          question_type: "single_choice",
          answers: [
            { text: "", is_correct: false },
            { text: "", is_correct: false },
          ],
        },
      ],
    });
  };

  const addAnswer = (qIndex: number) => {
    if (!testData) return;
    const updated = [...testData.questions];
    updated[qIndex].answers.push({ text: "", is_correct: false });
    setTestData({ ...testData, questions: updated });
  };

  const deleteAnswer = (qIndex: number, aIndex: number) => {
    if (!testData) return;
    const updated = [...testData.questions];
    updated[qIndex].answers.splice(aIndex, 1);
    setTestData({ ...testData, questions: updated });
  };

  // === ЛОГИКА НАСТРОЕК КУРСА ===
  const handleUpdateCourse = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!course) return;
    try {
      setIsSaving(true);
      await coursesApi.update(course.id, {
        title: course.title,
        description: course.description,
        subject_id: course.subject,
        difficulty_level: course.difficulty_level,
        cover_image_url: course.cover_image_url,
        tags: selectedTags,
      });
      alert("Настройки курса обновлены");
    } catch {
      alert("Ошибка");
    } finally {
      setIsSaving(false);
    }
  };

  const handlePublish = async () => {
    if (!course) return;
    const newState = !course.is_published;
    if (
      !confirm(
        newState
          ? "Опубликовать курс для студентов?"
          : "Снять курс с публикации?"
      )
    )
      return;

    try {
      await coursesApi.publish(course.id, newState);
      setCourse({ ...course, is_published: newState });
    } catch {
      alert("Ошибка смены статуса");
    }
  };

  const handleCoverUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.[0] && course) {
      const url = await uploadApi.uploadFile(e.target.files[0], "cover");
      setCourse({ ...course, cover_image_url: url });
    }
  };

  const toggleTag = (tagId: number) => {
    setSelectedTags((prev) =>
      prev.includes(tagId)
        ? prev.filter((id) => id !== tagId)
        : [...prev, tagId]
    );
  };

  const handleDeleteCourse = async () => {
    if (!course || !confirm("Удалить курс ПОЛНОСТЬЮ? Это действие необратимо!"))
      return;
    try {
      await coursesApi.delete(course.id);
      alert("Курс удален");
      navigate("/teacher/courses");
    } catch (error) {
      console.error(error);
      alert("Ошибка удаления курса");
    }
  };

  if (!course) return <div className="p-10 text-center">Загрузка курса...</div>;
  const renderPreviewModal = () => {
    if (!showPreview || !lessonData) return null;

    // Упрощенная версия плеера для превью
    return (
      <div className="fixed inset-0 z-[100] bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
        <div className="bg-white w-full max-w-4xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden">
          {/* Header */}
          <div className="p-4 border-b flex justify-between items-center bg-gray-50">
            <h3 className="font-bold text-lg text-gray-800">
              Предпросмотр: {lessonData.title}
            </h3>
            <button
              onClick={() => setShowPreview(false)}
              className="p-2 hover:bg-gray-200 rounded-full"
            >
              <X size={24} />
            </button>
          </div>

          {/* Body */}
          <div className="flex-1 overflow-y-auto p-8">
            {/* Video */}
            {lessonData.video_url && (
              <div className="mb-8 aspect-video bg-black rounded-xl overflow-hidden shadow-lg mx-auto">
                {lessonData.video_url.includes("youtu") ? (
                  <iframe
                    width="100%"
                    height="100%"
                    src={lessonData.video_url
                      .replace("watch?v=", "embed/")
                      .replace("youtu.be/", "youtube.com/embed/")}
                    title="Preview"
                    frameBorder="0"
                    allowFullScreen
                  />
                ) : (
                  <video
                    src={lessonData.video_url}
                    controls
                    className="w-full h-full"
                  />
                )}
              </div>
            )}

            {/* Markdown Content */}
            <div className="prose prose-indigo max-w-none text-gray-700 leading-relaxed mb-10">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {lessonData.content_text || "*Текст урока отсутствует*"}
              </ReactMarkdown>
            </div>

            {/* Attachments */}
            {lessonData.file_attachment_url && (
              <div className="p-4 border rounded-xl bg-gray-50 flex items-center gap-3">
                <span className="text-2xl">📎</span>
                <div>
                  <div className="font-bold text-sm text-gray-900">
                    Материалы урока
                  </div>
                  <div className="text-xs text-gray-500">
                    Доступны для скачивания
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    );
  };
  return (
    <div className="flex flex-col h-screen bg-gray-100">
      <CourseHeader
        courseTitle={course.title}
        isPublished={course.is_published}
        activeTab={activeTab}
        onTabChange={setActiveTab}
        onDelete={handleDeleteCourse}
      />
      {renderPreviewModal()}

      {activeTab === "curriculum" ? (
        <div className="flex flex-1 overflow-hidden">
          <ModulesList
            modules={modules}
            selectedLessonId={selectedLessonId}
            onLessonSelect={setSelectedLessonId}
            onLessonDelete={handleDeleteLesson}
            onModuleDelete={handleDeleteModule}
            onLessonAdd={handleAddLesson}
            onModuleAdd={handleAddModule}
            onTestOpen={handleOpenTestEditor}
          />

          <main className="flex-1 overflow-y-auto p-8 bg-gray-100">
            {lessonData ? (
              <LessonEditor
                lesson={lessonData}
                isSaving={isSaving}
                onSave={handleSaveLesson}
                onPreview={() => setShowPreview(true)}
                onChange={setLessonData}
                onUpload={handleUpload}
              />
            ) : selectedModuleForTest ? (
              <TestEditor
                test={testData}
                isCreating={isCreatingTest}
                isSaving={isSaving}
                onCreateTest={handleCreateTest}
                onSaveTest={handleSaveTest}
                onDeleteTest={handleDeleteTest}
                onChange={setTestData}
                onAddQuestion={addQuestion}
                onAddAnswer={addAnswer}
                onDeleteAnswer={deleteAnswer}
              />
            ) : (
              <EmptyState />
            )}
          </main>
        </div>
      ) : (
        <CourseSettings
          course={course}
          subjects={subjects}
          tags={tags}
          selectedTags={selectedTags}
          isSaving={isSaving}
          onCourseChange={setCourse}
          onTagToggle={toggleTag}
          onCoverUpload={handleCoverUpload}
          onSubmit={handleUpdateCourse}
          onPublish={handlePublish}
        />
      )}
    </div>
  );
};
