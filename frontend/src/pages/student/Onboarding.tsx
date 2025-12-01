import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { studentApi } from "../../api/student";
import { subjectsApi } from "../../api/subjects";
import { Button } from "../../components/ui/Button";
import type { Subject } from "../../types/subject";
import { Check, ChevronRight } from "lucide-react";

export const OnboardingPage = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState(1);
  const [isLoading, setIsLoading] = useState(false);

  // Данные формы
  const [grade, setGrade] = useState<number | null>(null);
  const [selectedSubjects, setSelectedSubjects] = useState<string[]>([]);

  // Доступные для выбора
  const [allSubjects, setAllSubjects] = useState<Subject[]>([]);

  useEffect(() => {
    // Грузим предметы сразу
    subjectsApi.getAll().then(setAllSubjects).catch(console.error);
  }, []);

  const handleNext = () => setStep((p) => p + 1);

  const toggleSubject = (id: string) => {
    setSelectedSubjects((prev) =>
      prev.includes(id) ? prev.filter((s) => s !== id) : [...prev, id]
    );
  };

  const handleFinish = async () => {
    if (!grade || selectedSubjects.length === 0) return;
    setIsLoading(true);
    try {
      await studentApi.completeOnboarding(grade, selectedSubjects);
      // После успеха - на дашборд
      navigate("/student/dashboard");
    } catch (error) {
      console.error(error);
      alert("Произошла ошибка при сохранении данных.");
    } finally {
      setIsLoading(false);
    }
  };

  // --- ШАГИ ИНТЕРФЕЙСА ---

  // Шаг 1: Приветствие
  if (step === 1) {
    return (
      <div className="min-h-screen bg-indigo-600 flex flex-col items-center justify-center text-white p-6 text-center">
        <h1 className="text-4xl font-bold mb-4 animate-fade-in-up">
          Добро пожаловать в OqysAI! 🚀
        </h1>
        <p className="text-xl mb-8 max-w-lg opacity-90">
          Твоя персональная платформа для обучения. Давай настроим профиль,
          чтобы рекомендовать тебе лучшие курсы.
        </p>
        <Button
          className="bg-white text-indigo-600 hover:bg-indigo-50 w-auto px-10 py-3 text-lg font-bold"
          onClick={handleNext}
        >
          Поехали!
        </Button>
      </div>
    );
  }

  // Шаг 2: Выбор класса
  if (step === 2) {
    return (
      <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-6">
        <div className="max-w-2xl w-full bg-white rounded-2xl shadow-xl p-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">
            В каком ты классе?
          </h2>

          <div className="grid grid-cols-3 sm:grid-cols-4 gap-4 mb-8">
            {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11].map((g) => (
              <button
                key={g}
                onClick={() => setGrade(g)}
                className={`py-4 rounded-xl text-xl font-bold transition-all border-2 ${
                  grade === g
                    ? "border-indigo-600 bg-indigo-50 text-indigo-700 scale-105 shadow-md"
                    : "border-gray-200 text-gray-600 hover:border-indigo-300 hover:bg-gray-50"
                }`}
              >
                {g}
              </button>
            ))}
          </div>

          <div className="flex justify-end">
            <Button
              onClick={handleNext}
              disabled={!grade}
              className="w-auto px-8"
            >
              Дальше <ChevronRight size={20} className="ml-2" />
            </Button>
          </div>
        </div>
      </div>
    );
  }

  // Шаг 3: Выбор интересов
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-6">
      <div className="max-w-3xl w-full bg-white rounded-2xl shadow-xl p-8">
        <h2 className="text-2xl font-bold text-gray-800 mb-2 text-center">
          Что тебе интересно?
        </h2>
        <p className="text-gray-500 text-center mb-8">
          Выберите хотя бы один предмет
        </p>

        {allSubjects.length === 0 ? (
          <p className="text-center text-gray-400 py-10">
            Загрузка предметов...
          </p>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 mb-8">
            {allSubjects.map((subject) => {
              const isSelected = selectedSubjects.includes(subject.id);
              return (
                <button
                  key={subject.id}
                  onClick={() => toggleSubject(subject.id)}
                  className={`relative p-4 rounded-xl text-left transition-all border-2 ${
                    isSelected
                      ? "border-indigo-600 bg-indigo-600 text-white shadow-lg transform scale-[1.02]"
                      : "border-gray-200 text-gray-700 hover:border-indigo-300 hover:bg-gray-50"
                  }`}
                >
                  <span className="font-semibold block">{subject.name_ru}</span>
                  {isSelected && (
                    <div className="absolute top-2 right-2 bg-white text-indigo-600 rounded-full p-0.5">
                      <Check size={12} strokeWidth={4} />
                    </div>
                  )}
                </button>
              );
            })}
          </div>
        )}

        <div className="flex justify-between items-center pt-4 border-t">
          <button
            onClick={() => setStep(2)}
            className="text-gray-500 hover:text-gray-800"
          >
            Назад
          </button>
          <Button
            onClick={handleFinish}
            disabled={selectedSubjects.length === 0}
            isLoading={isLoading}
            className="w-auto px-8 py-3 text-lg"
          >
            Завершить настройку 🎉
          </Button>
        </div>
      </div>
    </div>
  );
};
