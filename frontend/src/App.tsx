import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { RegisterPage } from "./pages/Register";
import { LoginPage } from "./pages/Login";
import { ChangePasswordPage } from "./pages/ChangePassword";
import { ProtectedRoute } from "./components/auth/ProtectedRoute";
import { CreateCoursePage } from "./pages/teacher/CreateCourse";
import { EditCoursePage } from "./pages/teacher/EditCourse";
import Dashboard from "./pages/teacher/Dashboard";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />

        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/change-password" element={<ChangePasswordPage />} />

        {/* Заглушка для дашборда */}
        <Route path="/dashboard" element={<Dashboard />} />
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
      </Routes>
    </BrowserRouter>
  );
}

export default App;
