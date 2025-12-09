import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { studentApi } from "../../api/student";
import { gamificationApi } from "../../api/gamification";
import type { League } from "../../api/gamification";
import type { DashboardData } from "../../api/student";
import { GamificationStats } from "../../components/student/GamificationStats";
import { ActiveCourseCard } from "../../components/student/ActiveCourseCard";
import { Button } from "../../components/ui/Button";
import { Compass, BookOpen } from "lucide-react";
import type { Course } from "../../types/course";
import { coursesApi } from "../../api/courses";
import { RecommendedCourseCard } from "./RecommendedCourseCard";

export const StudentDashboard = () => {
  const navigate = useNavigate();
  const [data, setData] = useState<DashboardData | null>(null);
  const [leagues, setLeagues] = useState<League[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [recommendations, setRecommendations] = useState<Course[]>([]);

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        // 1. Критически важные данные (если упадут, то страница не имеет смысла)
        const [dashboardData, leaguesList] = await Promise.all([
          studentApi.getDashboard(),
          gamificationApi.getAllLeagues(),
        ]);

        setData(dashboardData);
        setLeagues(leaguesList);

        // 2. Рекомендации загружаем отдельно (если упадут — не страшно)
        try {
          const recsList = await coursesApi.getRecommendations();
          setRecommendations(recsList);
        } catch (recError) {
          console.warn(
            "Не удалось загрузить рекомендации, но это не критично:",
            recError
          );
          setRecommendations([]); // Просто будет пустой список
        }
      } catch (error: unknown) {
        // Обработка критических ошибок (как у тебя сейчас)
        if (
          typeof error === "object" &&
          error !== null &&
          "response" in error &&
          (error as { response?: { status?: number } }).response?.status === 404
        ) {
          navigate("/student/onboarding");
        } else {
          console.error("Критическая ошибка загрузки дашборда", error);
        }
      } finally {
        setIsLoading(false);
      }
    };
    loadDashboard();
  }, [navigate]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  if (!data) return null;

  const currentLeague = leagues.find(
    (l) => l.id === data.profile.current_league_id
  );

  return (
    <div className="min-h-screen bg-gray-50 p-6 md:p-8">
      <div className="max-w-6xl mx-auto space-y-8">
        {/* 1. Приветствие и Статистика */}
        <header>
          <h1 className="text-3xl font-bold text-gray-900 mb-6">
            Привет! Готов к новым знаниям? 👋
          </h1>
          {data.profile && (
            <GamificationStats
              profile={data.profile}
              leagueName={currentLeague?.name || "Лига"}
              leagueIcon={currentLeague?.icon_url}
            />
          )}
        </header>

        {/* 2. Активные курсы */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-800 flex items-center gap-2">
              <BookOpen className="text-indigo-600" /> Продолжить обучение
            </h2>
          </div>

          {data.active_courses.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {data.active_courses.map((course) => (
                <ActiveCourseCard key={course.course_id} course={course} />
              ))}
            </div>
          ) : (
            <div className="bg-white p-8 rounded-xl border border-dashed border-gray-300 text-center">
              <p className="text-gray-500 mb-4">
                Вы еще не начали ни одного курса.
              </p>
              <Button
                onClick={() => navigate("/student/catalog")}
                className="w-auto"
              >
                Перейти в каталог
              </Button>
            </div>
          )}
        </section>

        {/* 3. Рекомендации */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-800 flex items-center gap-2">
              <Compass className="text-purple-600" /> Рекомендовано тебе
            </h2>
            <Button onClick={() => navigate("/student/catalog")}>
              Весь каталог →
            </Button>
          </div>

          {recommendations.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {recommendations.map((course) => (
                <RecommendedCourseCard key={course.id} course={course} />
              ))}
            </div>
          ) : (
            // Фоллбэк, если ML сервис молчит или курсов мало
            <div className="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl p-8 text-white flex items-center justify-between shadow-lg">
              <div>
                <h3 className="text-2xl font-bold mb-2">
                  Ищем идеальные курсы для тебя... 🤖
                </h3>
                <p className="opacity-90">
                  Посмотри наш каталог, пока мы настраиваем алгоритмы.
                </p>
              </div>
              <Button
                className="bg-white text-indigo-600 hover:bg-indigo-50 w-auto border-none"
                onClick={() => navigate("/student/catalog")}
              >
                Открыть каталог
              </Button>
            </div>
          )}
        </section>
      </div>
    </div>
  );
};
