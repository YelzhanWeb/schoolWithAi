import api from "./axios";

export interface Answer {
  id?: string;
  text: string;
  is_correct: boolean;
}

export interface Question {
  text: string;
  question_type: string; // "single_choice"
  answers: Answer[];
}

export interface Test {
  id?: string;
  test_id: string;
  title: string;
  passing_score: number;
  questions: Question[];
}

export interface CreateTestRequest {
  module_id: string;
  title: string;
  passing_score: number;
  questions: Question[];
}

export interface SubmitTestResponse {
  is_passed: boolean;
  score: number;
  xp_gained: number;
}

export interface StudentAnswer {
  question_id: string;
  answer_id: string;
}

export const testsApi = {
  getByModuleId: async (moduleId: string): Promise<Test> => {
    const response = await api.get<Test>(`/modules/${moduleId}/test`);
    return response.data;
  },

  create: async (data: CreateTestRequest): Promise<{ test_id: string }> => {
    const response = await api.post<{ test_id: string }>("/tests", data);
    return response.data;
  },

  update: async (id: string, data: Test) => {
    await api.put(`/tests/${id}`, data);
  },

  delete: async (id: string) => {
    await api.delete(`/tests/${id}`);
  },

  completeLesson: async (lessonId: string): Promise<{ xp_gained: number }> => {
    const response = await api.post<{ xp_gained: number }>(
      `/student/lessons/${lessonId}/complete`
    );
    return response.data;
  },

  // Метод для отправки теста
  submitTest: async (
    testId: string,
    answers: StudentAnswer[]
  ): Promise<SubmitTestResponse> => {
    const response = await api.post<SubmitTestResponse>(
      "/student/tests/submit",
      {
        test_id: testId,
        answers: answers,
      }
    );
    return response.data;
  },
};
