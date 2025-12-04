import { Link, useNavigate, useLocation } from "react-router-dom";
import {
  BookOpen,
  LogOut,
  LayoutDashboard,
  Search,
  Bell,
  Flame,
} from "lucide-react";
import { Button } from "../../components/ui/Button";

export const StudentHeader = () => {
  const navigate = useNavigate();
  const location = useLocation();

  // Функция для проверки активной ссылки
  const isActive = (path: string) => location.pathname.startsWith(path);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <header className="h-16 bg-white border-b border-gray-200 fixed top-0 w-full z-50 flex items-center justify-between px-4 md:px-8 shadow-sm">
      {/* ЛЕВАЯ ЧАСТЬ: Лого + Навигация */}
      <div className="flex items-center gap-8">
        {/* Логотип */}
        <Link
          to="/student/dashboard"
          className="flex items-center gap-2 hover:opacity-80 transition"
        >
          <div className="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center text-white font-bold text-lg shadow-indigo-200 shadow-lg">
            Q
          </div>
          <span className="text-xl font-bold text-gray-900 tracking-tight hidden md:block">
            OqysAI
          </span>
        </Link>

        {/* Ссылки меню */}
        <nav className="hidden md:flex items-center gap-1">
          <Link to="/student/my-learning">
            <Button
              className={`flex items-center gap-2 ${
                isActive("/student/my-learning")
                  ? "bg-gray-100 text-gray-900"
                  : "text-gray-600"
              }`}
            >
              <LayoutDashboard size={18} /> Мое обучение
            </Button>
          </Link>

          <Link to="/student/catalog">
            <Button
              className={`flex items-center gap-2 ${
                isActive("/student/catalog")
                  ? "bg-gray-100 text-gray-900"
                  : "text-gray-600"
              }`}
            >
              <BookOpen size={18} /> Каталог
            </Button>
          </Link>
        </nav>
      </div>

      {/* ПРАВАЯ ЧАСТЬ: Поиск, Геймификация, Профиль */}
      <div className="flex items-center gap-3 md:gap-6">
        {/* Поиск (на мобильных скрываем или показываем только иконку) */}
        <div className="relative hidden md:block">
          <Search
            className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
            size={16}
          />
          <input
            type="text"
            placeholder="Найти курс..."
            className="pl-9 pr-4 py-1.5 rounded-full border border-gray-200 bg-gray-50 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 w-48 transition-all focus:w-64"
          />
        </div>

        <div className="h-6 w-px bg-gray-200 hidden md:block"></div>

        {/* Геймификация (Мини-виджет) */}
        <div
          className="flex items-center gap-2 bg-orange-50 text-orange-600 px-3 py-1.5 rounded-full border border-orange-100 cursor-help"
          title="Ваш текущий стрик"
        >
          <Flame size={16} fill="currentColor" />
          <span className="text-sm font-bold">5</span>
        </div>

        <button className="text-gray-500 hover:text-indigo-600 transition relative">
          <Bell size={20} />
          <span className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
        </button>

        {/* Аватар / Меню профиля */}
        <div className="group relative">
          <div className="w-9 h-9 bg-indigo-100 rounded-full flex items-center justify-center text-indigo-700 font-bold cursor-pointer border border-indigo-200 hover:shadow-md transition">
            S
          </div>

          {/* Выпадающее меню (на ховер или клик) */}
          <div className="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-xl border border-gray-100 py-1 hidden group-hover:block animate-fade-in-up">
            <div className="px-4 py-2 border-b border-gray-100 mb-1">
              <p className="text-sm font-bold text-gray-900">Student Name</p>
              <p className="text-xs text-gray-500">student@example.com</p>
            </div>
            <Link
              to="/student/profile"
              className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 hover:text-indigo-600"
            >
              Настройки профиля
            </Link>
            <button
              onClick={handleLogout}
              className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
            >
              <LogOut size={14} /> Выйти
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};
