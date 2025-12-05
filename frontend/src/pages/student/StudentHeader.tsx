import { useEffect, useState } from "react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { Trophy } from "lucide-react";
import {
  BookOpen,
  LogOut,
  LayoutDashboard,
  Search,
  Flame,
  Settings,
} from "lucide-react";
import { Button } from "../../components/ui/Button";
import { studentApi, type HeaderInfo } from "../../api/student";

export const StudentHeader = () => {
  const navigate = useNavigate();
  const location = useLocation();

  // Состояние теперь типизировано нашим новым интерфейсом
  const [user, setUser] = useState<HeaderInfo | null>(null);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const isActive = (path: string) => location.pathname.startsWith(path);

  useEffect(() => {
    const loadData = async () => {
      try {
        // Вызываем новый легкий метод
        const data = await studentApi.getMe();
        setUser(data);
      } catch (e) {
        console.error("Ошибка загрузки профиля", e);
      }
    };
    loadData();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <header className="h-16 bg-white border-b border-gray-200 fixed top-0 w-full z-50 flex items-center justify-between px-4 md:px-8 shadow-sm">
      {/* ЛЕВАЯ ЧАСТЬ */}
      <div className="flex items-center gap-8">
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
          <Link to="/student/leaderboard">
            <Button
              className={`flex items-center gap-2 ${
                isActive("/student/leaderboard")
                  ? "bg-gray-100 text-gray-900"
                  : "text-gray-600"
              }`}
            >
              <Trophy size={18} /> Рейтинг
            </Button>
          </Link>
        </nav>
      </div>

      {/* ПРАВАЯ ЧАСТЬ */}
      <div className="flex items-center gap-3 md:gap-6">
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

        {/* Стрик */}
        {user && user.current_streak > 0 && (
          <div
            className="flex items-center gap-2 bg-orange-50 text-orange-600 px-3 py-1.5 rounded-full border border-orange-100 cursor-help"
            title="Дней подряд"
          >
            <Flame size={16} fill="currentColor" />
            <span className="text-sm font-bold">{user.current_streak}</span>
          </div>
        )}

        {/* Аватар / Меню */}
        <div className="relative">
          <div
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className="w-9 h-9 bg-indigo-100 rounded-full flex items-center justify-center text-indigo-700 font-bold cursor-pointer border border-indigo-200 hover:shadow-md transition overflow-hidden"
          >
            {user?.avatar_url ? (
              <img
                src={user.avatar_url}
                alt="Ava"
                className="w-full h-full object-cover"
              />
            ) : (
              <span>{user?.first_name?.[0] || "U"}</span>
            )}
          </div>

          {/* Выпадающее меню */}
          {isMenuOpen && (
            <>
              <div
                className="fixed inset-0 z-40"
                onClick={() => setIsMenuOpen(false)}
              ></div>
              <div className="absolute right-0 mt-2 w-56 bg-white rounded-xl shadow-xl border border-gray-100 py-2 z-50 animate-in fade-in zoom-in-95 duration-100">
                <div className="px-4 py-3 border-b border-gray-100 mb-1">
                  <p className="text-sm font-bold text-gray-900 truncate">
                    {user
                      ? `${user.first_name} ${user.last_name}`
                      : "Загрузка..."}
                  </p>
                  <p className="text-xs text-gray-500 truncate">
                    {user?.email}
                  </p>
                </div>

                <button
                  className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
                  onClick={() => navigate("/student/profile")}
                >
                  <Settings size={16} /> Настройки профиля
                </button>

                <button
                  onClick={handleLogout}
                  className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
                >
                  <LogOut size={16} /> Выйти
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  );
};
