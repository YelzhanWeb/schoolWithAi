import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { RegisterPage } from "./pages/Register";
import { LoginPage } from "./pages/Login";
import { ChangePasswordPage } from "./pages/ChangePassword";
import { ProtectedRoute } from "./components/auth/ProtectedRoute";
import { CreateCoursePage } from "./pages/teacher/CreateCourse";
import { EditCoursePage } from "./pages/teacher/EditCourse";
import TeacherDashboard from "./pages/teacher/TeacherDashboard";
import { OnboardingPage } from "./pages/student/Onboarding";
import { StudentDashboard } from "./pages/student/StudentDashboard";
import { CatalogPage } from "./pages/student/Catalog";
import { CourseOverview } from "./pages/student/CourseOverview";
import { StudentLayout } from "./pages/student/StudentLayout";
import { LessonPlayer } from "./pages/student/LessonPlayer";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/change-password" element={<ChangePasswordPage />} />

        {/* TEACHER */}
        <Route
          path="/teacher/dashboard"
          element={
            <ProtectedRoute roles={["teacher", "admin"]}>
              <TeacherDashboard />
            </ProtectedRoute>
          }
        />

        {/* ADMIN
        <Route
          path="/admin/dashboard"
          element={
            <ProtectedRoute roles={["admin"]}>
              <AdminDashboard />
            </ProtectedRoute>
          }
        /> */}
        {/* Teacher pages */}
        <Route
          path="/teacher/create-course"
          element={
            <ProtectedRoute roles={["teacher", "admin"]}>
              <CreateCoursePage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/teacher/courses/:id/edit"
          element={
            <ProtectedRoute roles={["teacher", "admin"]}>
              <EditCoursePage />
            </ProtectedRoute>
          }
        />

        {/* Student pages */}
        <Route
          path="/student/onboarding"
          element={
            <ProtectedRoute roles={["student", "admin"]}>
              <OnboardingPage />
            </ProtectedRoute>
          }
        />

        <Route
          element={
            <ProtectedRoute roles={["student", "admin"]}>
              <StudentLayout />
            </ProtectedRoute>
          }
        >
          <Route
            path="/student/catalog"
            element={
              <ProtectedRoute roles={["student", "admin"]}>
                <CatalogPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/student/courses/:id"
            element={
              <ProtectedRoute roles={["student", "admin"]}>
                <CourseOverview />
              </ProtectedRoute>
            }
          />
          <Route
            path="/student/dashboard"
            element={
              <ProtectedRoute roles={["student"]}>
                <StudentDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/student/courses/:courseId/lessons/:lessonId"
            element={<LessonPlayer />}
          />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
