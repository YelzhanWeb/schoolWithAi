import React from "react";
import { Plus, Trash2, FileText } from "lucide-react";
import type { Test, Question, Answer } from "../../api/tests";
import { Button } from "../ui/Button";
import { Input } from "../ui/Input";

interface TestEditorProps {
  test: Test | null;
  isCreating: boolean;
  isSaving: boolean;
  onCreateTest: () => void;
  onSaveTest: () => void;
  onDeleteTest?: () => void;
  onChange: (test: Test) => void;
  onAddQuestion: () => void;
  onAddAnswer: (qIndex: number) => void;
  onDeleteAnswer: (qIndex: number, aIndex: number) => void;
}

export const TestEditor: React.FC<TestEditorProps> = ({
  test,
  isCreating,
  isSaving,
  onCreateTest,
  onSaveTest,
  onDeleteTest,
  onChange,
  onAddQuestion,
  onAddAnswer,
  onDeleteAnswer,
}) => {
  const updateQuestion = (
    index: number,
    field: keyof Question,
    value: string
  ) => {
    if (!test) return;
    const updated = [...test.questions];
    updated[index] = { ...updated[index], [field]: value };
    onChange({ ...test, questions: updated });
  };

  const updateAnswer = (
    qIndex: number,
    aIndex: number,
    field: keyof Answer,
    value: string | boolean
  ) => {
    if (!test) return;
    const updated = [...test.questions];
    updated[qIndex].answers[aIndex] = {
      ...updated[qIndex].answers[aIndex],
      [field]: value,
    };
    onChange({ ...test, questions: updated });
  };

  return (
    <div className="max-w-3xl mx-auto bg-white rounded-xl shadow-sm p-8 space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-gray-800">
          {isCreating ? "Создание теста" : "Редактирование теста"}
        </h2>
        <div className="flex gap-2">
          {!test && !isCreating && (
            <Button onClick={onCreateTest} className="w-auto">
              Создать тест
            </Button>
          )}
          {isCreating && (
            <Button
              onClick={onSaveTest}
              isLoading={isSaving}
              className="w-auto"
            >
              Сохранить тест
            </Button>
          )}
          {!isCreating && test && onDeleteTest && (
            <button
              onClick={onDeleteTest}
              className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"
              title="Удалить тест"
            >
              <Trash2 size={20} />
            </button>
          )}
        </div>
      </div>

      {test ? (
        <div className="space-y-6">
          <Input
            label="Название теста"
            value={test.title}
            onChange={(e) => onChange({ ...test, title: e.target.value })}
            disabled={!isCreating}
          />
          <Input
            label="Проходной балл (%)"
            type="number"
            value={test.passing_score}
            onChange={(e) =>
              onChange({ ...test, passing_score: Number(e.target.value) })
            }
            disabled={!isCreating}
          />

          <div className="space-y-4">
            <h3 className="font-bold">Вопросы:</h3>
            {test.questions.map((q, qIndex) => (
              <div
                key={qIndex}
                className="border rounded-lg p-4 space-y-3 bg-gray-50"
              >
                <Input
                  label={`Вопрос ${qIndex + 1}`}
                  value={q.text}
                  onChange={(e) =>
                    updateQuestion(qIndex, "text", e.target.value)
                  }
                  disabled={!isCreating}
                />
                <div className="space-y-2">
                  <label className="block text-sm font-medium">Ответы:</label>
                  {q.answers.map((a, aIndex) => (
                    <div key={aIndex} className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={a.is_correct}
                        onChange={(e) =>
                          updateAnswer(
                            qIndex,
                            aIndex,
                            "is_correct",
                            e.target.checked
                          )
                        }
                        disabled={!isCreating}
                        className="w-5 h-5"
                      />
                      <Input
                        placeholder="Текст ответа"
                        value={a.text}
                        onChange={(e) =>
                          updateAnswer(qIndex, aIndex, "text", e.target.value)
                        }
                        disabled={!isCreating}
                        className="flex-1"
                      />
                      {isCreating && q.answers.length > 2 && (
                        <button
                          onClick={() => onDeleteAnswer(qIndex, aIndex)}
                          className="text-red-500 hover:text-red-700"
                        >
                          <Trash2 size={16} />
                        </button>
                      )}
                    </div>
                  ))}
                  {isCreating && (
                    <button
                      onClick={() => onAddAnswer(qIndex)}
                      className="text-indigo-600 text-sm hover:underline"
                    >
                      + Добавить ответ
                    </button>
                  )}
                </div>
              </div>
            ))}
            {isCreating && (
              <Button
                onClick={onAddQuestion}
                variant="outline"
                className="w-full"
              >
                <Plus size={16} /> Добавить вопрос
              </Button>
            )}
          </div>
        </div>
      ) : (
        <div className="text-center text-gray-400 py-10">
          <FileText size={48} className="mx-auto mb-4 opacity-50" />
          <p>Для этого модуля теста ещё нет</p>
        </div>
      )}
    </div>
  );
};
