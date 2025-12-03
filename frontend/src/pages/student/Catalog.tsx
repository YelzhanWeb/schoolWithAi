import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { coursesApi } from "../../api/courses";
import { subjectsApi } from "../../api/subjects";
import type { Course, Tag } from "../../types/course";
import type { Subject } from "../../types/subject";
import { Button } from "../../components/ui/Button";
import { Search, Filter } from "lucide-react";

export const CatalogPage = () => {
  // Данные
  const [courses, setCourses] = useState<Course[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [allTags, setAllTags] = useState<Tag[]>([]);

  // Состояние отображения
  const [filteredCourses, setFilteredCourses] = useState<Course[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Состояние фильтров
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedSubject, setSelectedSubject] = useState<string>("");
  const [selectedDifficulty, setSelectedDifficulty] = useState<number | "">("");
  const [selectedTags, setSelectedTags] = useState<number[]>([]);

  // Мобильное меню фильтров
  const [showMobileFilters, setShowMobileFilters] = useState(false);

  // 1. Загрузка справочников и курсов
  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true);
        const [coursesData, subjectsData, tagsData] = await Promise.all([
          coursesApi.getCatalog(),
          subjectsApi.getAll(),
          coursesApi.getAllTags(),
        ]);
        setCourses(coursesData);
        setSubjects(subjectsData);
        setAllTags(tagsData);
        setFilteredCourses(coursesData);
      } catch (error) {
        console.error("Ошибка загрузки данных", error);
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, []);

  // 2. Логика фильтрации (Client-Side)
  useEffect(() => {
    const filtered = courses.filter((c) => {
      // Поиск по названию
      const matchesSearch = c.title
        .toLowerCase()
        .includes(searchQuery.toLowerCase());

      // Фильтр по предмету
      const matchesSubject = selectedSubject
        ? c.subject_id === selectedSubject
        : true;

      // Фильтр по сложности
      const matchesDifficulty = selectedDifficulty
        ? c.difficulty_level === selectedDifficulty
        : true;

      // Фильтр по тегам (если выбран хоть один тег, курс должен содержать его)
      // Логика OR (ИЛИ): если выбран тег A и B, покажет курсы где есть A ИЛИ B.
      // Можно сделать AND, если нужно строгое совпадение.
      const matchesTags =
        selectedTags.length > 0
          ? c.tags?.some((tag) => selectedTags.includes(tag.id))
          : true;

      return (
        matchesSearch && matchesSubject && matchesDifficulty && matchesTags
      );
    });

    setFilteredCourses(filtered);
  }, [searchQuery, selectedSubject, selectedDifficulty, selectedTags, courses]);

  // Хендлер для тегов (toggle)
  const toggleTag = (tagId: number) => {
    setSelectedTags((prev) =>
      prev.includes(tagId)
        ? prev.filter((id) => id !== tagId)
        : [...prev, tagId]
    );
  };

  // Сброс фильтров
  const resetFilters = () => {
    setSearchQuery("");
    setSelectedSubject("");
    setSelectedDifficulty("");
    setSelectedTags([]);
  };

  return (
    <div className="max-w-7xl mx-auto p-4 md:p-8">
      <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <h1 className="text-3xl font-bold text-gray-900">Каталог курсов</h1>

        {/* Поиск (Верхний) */}
        <div className="relative w-full md:w-96 flex gap-2">
          <div className="relative flex-1">
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
          {/* Кнопка фильтров для мобилок */}
          <button
            className="md:hidden p-2 border rounded-lg bg-white"
            onClick={() => setShowMobileFilters(!showMobileFilters)}
          >
            <Filter size={20} />
          </button>
        </div>
      </div>

      <div className="flex flex-col md:flex-row gap-8">
        {/* --- ЛЕВАЯ КОЛОНКА: ФИЛЬТРЫ --- */}
        <aside
          className={`md:w-64 flex-shrink-0 space-y-6 ${
            showMobileFilters ? "block" : "hidden md:block"
          }`}
        >
          <div className="bg-white p-5 rounded-xl border border-gray-200 shadow-sm">
            <div className="flex justify-between items-center mb-4">
              <h3 className="font-bold text-gray-800">Фильтры</h3>
              {(selectedSubject ||
                selectedDifficulty ||
                selectedTags.length > 0) && (
                <button
                  onClick={resetFilters}
                  className="text-xs text-red-500 hover:underline"
                >
                  Сбросить
                </button>
              )}
            </div>

            {/* Предмет */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Предмет
              </label>
              <select
                className="w-full p-2 border rounded-lg bg-gray-50 text-sm outline-none focus:border-indigo-500"
                value={selectedSubject}
                onChange={(e) => setSelectedSubject(e.target.value)}
              >
                <option value="">Все предметы</option>
                {subjects.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name_ru}
                  </option>
                ))}
              </select>
            </div>

            {/* Сложность */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Сложность
              </label>
              <div className="space-y-2">
                {[1, 2, 3, 4, 5].map((level) => (
                  <label
                    key={level}
                    className="flex items-center cursor-pointer"
                  >
                    <input
                      type="radio"
                      name="difficulty"
                      className="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
                      checked={selectedDifficulty === level}
                      onChange={() => setSelectedDifficulty(level)}
                    />
                    <span className="ml-2 text-sm text-gray-600">
                      {level === 1
                        ? "1 - Новичок"
                        : level === 5
                        ? "5 - Эксперт"
                        : `${level} уровень`}
                    </span>
                  </label>
                ))}
                <label className="flex items-center cursor-pointer">
                  <input
                    type="radio"
                    name="difficulty"
                    className="w-4 h-4 text-indigo-600"
                    checked={selectedDifficulty === ""}
                    onChange={() => setSelectedDifficulty("")}
                  />
                  <span className="ml-2 text-sm text-gray-600">Любая</span>
                </label>
              </div>
            </div>

            {/* Теги */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Теги
              </label>
              <div className="flex flex-wrap gap-2">
                {allTags.map((tag) => (
                  <button
                    key={tag.id}
                    onClick={() => toggleTag(tag.id)}
                    className={`text-xs px-2 py-1 rounded-full border transition-colors ${
                      selectedTags.includes(tag.id)
                        ? "bg-indigo-100 border-indigo-200 text-indigo-700"
                        : "bg-white border-gray-200 text-gray-600 hover:border-gray-300"
                    }`}
                  >
                    {tag.name}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </aside>

        {/* --- ПРАВАЯ КОЛОНКА: СПИСОК КУРСОВ --- */}
        <div className="flex-1">
          {isLoading ? (
            <div className="flex justify-center py-20">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
            </div>
          ) : filteredCourses.length === 0 ? (
            <div className="text-center py-20 bg-white rounded-xl border border-dashed border-gray-300">
              <p className="text-gray-500 text-lg">Ничего не найдено 😔</p>
              <Button
                variant="outline"
                onClick={resetFilters}
                className="mt-4 w-auto"
              >
                Сбросить фильтры
              </Button>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredCourses.map((course) => (
                <Link
                  key={course.id}
                  to={`/student/courses/${course.id}`}
                  className="bg-white border rounded-xl overflow-hidden shadow-sm hover:shadow-lg transition-all group flex flex-col h-full"
                >
                  <div className="h-44 bg-gray-100 relative overflow-hidden">
                    {course.cover_image_url ? (
                      <img
                        src={course.cover_image_url}
                        alt={course.title}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-50">
                        Нет изображения
                      </div>
                    )}

                    {/* Бейджи поверх картинки */}
                    <div className="absolute top-2 right-2 flex flex-col gap-1 items-end">
                      <span className="bg-white/95 backdrop-blur px-2 py-1 rounded text-xs font-bold text-gray-700 shadow-sm">
                        Lvl {course.difficulty_level}
                      </span>
                    </div>
                  </div>

                  <div className="p-5 flex-1 flex flex-col">
                    {/* Теги в карточке */}
                    <div className="flex flex-wrap gap-1 mb-2">
                      {course.tags?.slice(0, 3).map((t) => (
                        <span
                          key={t.id}
                          className="text-[10px] bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded"
                        >
                          #{t.name}
                        </span>
                      ))}
                    </div>

                    <h3 className="text-lg font-bold text-gray-900 mb-2 line-clamp-2 leading-tight">
                      {course.title}
                    </h3>
                    <p className="text-sm text-gray-600 mb-4 line-clamp-2 flex-1">
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
      </div>
    </div>
  );
};
