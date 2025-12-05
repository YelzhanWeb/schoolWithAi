import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { testsApi, type Test, type StudentAnswer } from "../../api/tests";
import { Button } from "../../components/ui/Button";
import { CheckCircle, XCircle, AlertCircle } from "lucide-react";

export const TestPlayer = () => {
  const { moduleId } = useParams<{ moduleId: string }>();
  const navigate = useNavigate();

  const [test, setTest] = useState<Test | null>(null);
  const [answers, setAnswers] = useState<Record<string, string>>({}); // questionId -> answerId
  const [result, setResult] = useState<{
    is_passed: boolean;
    score: number;
    xp_gained: number;
  } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!moduleId) return;
    const loadTest = async () => {
      try {
        const data = await testsApi.getByModuleId(moduleId);
        setTest(data);
      } catch (error) {
        console.error("Test not found", error);
      } finally {
        setIsLoading(false);
      }
    };
    loadTest();
  }, [moduleId]);

  const handleSelectAnswer = (questionId: string, answerId: string) => {
    setAnswers((prev) => ({ ...prev, [questionId]: answerId }));
  };

  const handleSubmit = async () => {
    if (!test) return;
    setIsSubmitting(true);

    const submitData: StudentAnswer[] = Object.entries(answers).map(
      ([qId, aId]) => ({
        question_id: qId,
        answer_id: aId,
      })
    );

    try {
      const res = await testsApi.submitTest(test.test_id, submitData);
      setResult(res);
    } catch (error) {
      alert("Ошибка отправки теста");
      console.log(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading)
    return <div className="p-10 text-center">Загрузка теста...</div>;
  if (!test) return <div className="p-10 text-center">Тест не найден</div>;

  // ЭКРАН РЕЗУЛЬТАТА
  if (result) {
    return (
      <div className="max-w-2xl mx-auto p-8 bg-white rounded-xl shadow-lg text-center mt-10">
        <div className="mb-6 flex justify-center">
          {result.is_passed ? (
            <CheckCircle className="text-green-500 w-20 h-20" />
          ) : (
            <XCircle className="text-red-500 w-20 h-20" />
          )}
        </div>

        <h2 className="text-3xl font-bold mb-2">
          {result.is_passed ? "Тест сдан! 🎉" : "Тест не сдан 😔"}
        </h2>

        <p className="text-gray-500 mb-6">
          Ваш результат:{" "}
          <span className="font-bold text-gray-800">{result.score}%</span>
          (Проходной: {test.passing_score}%)
        </p>

        {result.xp_gained > 0 && (
          <div className="bg-orange-50 text-orange-600 font-bold py-3 px-6 rounded-full inline-block mb-8">
            +{result.xp_gained} XP Получено
          </div>
        )}

        <div className="flex gap-4 justify-center">
          {!result.is_passed && (
            <Button onClick={() => window.location.reload()} variant="outline">
              Попробовать снова
            </Button>
          )}
          <Button onClick={() => navigate(-1)}>Вернуться к курсу</Button>
        </div>
      </div>
    );
  }

  // ЭКРАН ПРОХОЖДЕНИЯ
  return (
    <div className="max-w-3xl mx-auto p-6 md:p-10">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">{test.title}</h1>
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <AlertCircle size={16} />
          <span>Нужно набрать {test.passing_score}% для сдачи</span>
        </div>
      </div>

      <div className="space-y-8">
        {test.questions.map((q, index) => (
          <div
            key={q.text + index}
            className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm"
          >
            <h3 className="font-medium text-lg mb-4">
              {index + 1}. {q.text}
            </h3>
            <div className="space-y-3">
              {q.answers.map((a) => (
                <label
                  key={a.id || a.text} // Используем ID если есть, иначе текст (но лучше ID)
                  className={`flex items-center p-3 rounded-lg border cursor-pointer transition-all ${
                    answers[q.id || ""] === a.id // Тут важно: q.id должен приходить с бэка. Проверьте DTO.
                      ? "border-indigo-600 bg-indigo-50 ring-1 ring-indigo-600"
                      : "border-gray-200 hover:bg-gray-50"
                  }`}
                >
                  <input
                    type="radio"
                    name={q.id || `q-${index}`} // Группировка радио-кнопок
                    value={a.id}
                    checked={answers[q.id || ""] === a.id}
                    onChange={() => handleSelectAnswer(q.id || "", a.id || "")}
                    className="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500 mr-3"
                  />
                  <span className="text-gray-700">{a.text}</span>
                </label>
              ))}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-8 flex justify-end">
        <Button
          onClick={handleSubmit}
          isLoading={isSubmitting}
          disabled={Object.keys(answers).length < test.questions.length} // Блокируем, пока не ответит на все
          className="w-full md:w-auto px-8 py-3 text-lg"
        >
          Завершить тест
        </Button>
      </div>
    </div>
  );
};
