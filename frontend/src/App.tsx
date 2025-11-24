import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { RegisterPage } from "./pages/Register";
import { LoginPage } from "./pages/Login"; // Импортируем логин
import { ChangePasswordPage } from "./pages/ChangePassword";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />

        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/change-password" element={<ChangePasswordPage />} />

        {/* Заглушка для дашборда */}
        <Route
          path="/dashboard"
          element={
            <div className="p-10 text-2xl">Добро пожаловать в OqysAI! 🚀</div>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
