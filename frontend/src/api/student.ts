import api from "./axios";
import type { Subject } from "../types/subject";

export interface StudentProfile {
  id: string;
  grade: number;
  xp: number;
  level: number;
  current_league_id: number;
  weekly_xp: number;
  current_streak: number;
}

export interface HeaderInfo {
  first_name: string;
  last_name: string;
  email: string;
  avatar_url: string;
  current_streak: number;
  xp: number;
}

export interface DashboardData {
  profile: StudentProfile;
  interests: Subject[];
  active_courses: ActiveCourse[];
}

export interface ActiveCourse {
  course_id: string;
  title: string;
  cover_url: string;
  progress_percentage: number;
  total_lessons: number;
  completed_lessons: number;
}

export const studentApi = {
  completeOnboarding: async (grade: number, subjectIds: string[]) => {
    await api.post("/student/onboarding", {
      grade,
      subject_ids: subjectIds,
    });
  },

  getMe: async (): Promise<HeaderInfo> => {
    const response = await api.get<HeaderInfo>("/student/me");
    return response.data;
  },

  getDashboard: async (): Promise<DashboardData> => {
    const response = await api.get<DashboardData>("/student/dashboard");
    return response.data;
  },

  getCourseProgress: async (courseId: string): Promise<string[]> => {
    const response = await api.get<{ completed_lessons: string[] }>(
      `/student/courses/${courseId}/progress`
    );
    return response.data.completed_lessons || [];
  },
  getMyCourses: async (): Promise<ActiveCourse[]> => {
    const response = await api.get<{ courses: ActiveCourse[] }>(
      "/student/my-activity-courses"
    );
    return response.data.courses || [];
  },
};
