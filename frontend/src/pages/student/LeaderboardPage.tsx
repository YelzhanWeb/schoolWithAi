import { useEffect, useState } from "react";
import { studentApi } from "../../api/student";
import {
  gamificationApi,
  type LeaderboardEntry,
  type League,
} from "../../api/gamification";
import { Trophy, Globe, Clock, ArrowUp, ArrowDown, Minus } from "lucide-react";

type Tab = "weekly" | "global";

export const LeaderboardPage = () => {
  const [activeTab, setActiveTab] = useState<Tab>("weekly");
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [timeLeft, setTimeLeft] = useState("");

  // Данные для отображения лиги
  const [currentLeague, setCurrentLeague] = useState<League | null>(null);
  const [userProfileId, setUserProfileId] = useState<string>("");

  useEffect(() => {
    loadData();
    // Запускаем таймер
    const timer = setInterval(calculateTimeLeft, 60000); // Обновляем раз в минуту
    calculateTimeLeft();
    return () => clearInterval(timer);
  }, [activeTab]);

  const loadData = async () => {
    setIsLoading(true);
    try {
      // 1. Узнаем ID текущего юзера, чтобы подсветить его в списке
      //   const me = await studentApi.getMe();
      // В идеале getMe должен возвращать ID, если нет - придется декодировать токен или брать из LS.
      // Допустим, мы сохранили ID в localStorage при логине, или добавим поле id в ответ /student/me
      // Для подсветки пока будем использовать email или имя, если ID нет.

      // А пока для демо:
      const dashboard = await studentApi.getDashboard();
      setUserProfileId(dashboard.profile.id);

      // 2. Грузим лиги, чтобы найти картинку текущей
      const leagues = await gamificationApi.getAllLeagues();
      const myLeague = leagues.find(
        (l) => l.id === dashboard.profile.current_league_id
      );
      setCurrentLeague(myLeague || null);

      // 3. Грузим сам лидерборд
      if (activeTab === "weekly") {
        const data = await gamificationApi.getWeekly();
        setEntries(data.leaderboard);
      } else {
        const data = await gamificationApi.getGlobal();
        setEntries(data.leaderboard);
      }
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  // Расчет времени до следующего понедельника 00:00 UTC
  const calculateTimeLeft = () => {
    const now = new Date();
    const nowUTC = new Date(now.getTime() + now.getTimezoneOffset() * 60000);

    // Находим следующий понедельник
    const nextMonday = new Date(nowUTC);
    nextMonday.setDate(nowUTC.getDate() + ((1 + 7 - nowUTC.getDay()) % 7 || 7));
    nextMonday.setHours(0, 0, 0, 0);

    // Если сегодня понедельник, но время уже прошло 00:00, то nextMonday будет "сегодня".
    // Нам нужна следующая неделя, поэтому если diff < 0, добавляем 7 дней.
    let diff = nextMonday.getTime() - nowUTC.getTime();
    if (diff <= 0) {
      diff += 7 * 24 * 60 * 60 * 1000;
    }

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff / (1000 * 60 * 60)) % 24);

    setTimeLeft(`${days}д ${hours}ч`);
  };

  // Логика зон (20% вверх, 20% вниз)
  const getZoneStyle = (index: number, total: number) => {
    if (activeTab === "global") return ""; // В глобальном нет вылета

    const promoteCount = Math.ceil(total * 0.2); // Топ 20%
    const demoteCount = Math.floor(total * 0.2); // Низ 20%

    if (index < promoteCount) return "bg-green-50 border-l-4 border-green-500"; // Зона повышения
    if (index >= total - demoteCount)
      return "bg-red-50 border-l-4 border-red-500"; // Зона понижения

    return "border-l-4 border-transparent"; // Обычная зона
  };

  const getZoneIcon = (index: number, total: number) => {
    if (activeTab === "global") return null;
    const promoteCount = Math.ceil(total * 0.2);
    const demoteCount = Math.floor(total * 0.2);

    if (index < promoteCount)
      return <ArrowUp size={16} className="text-green-600" />;
    if (index >= total - demoteCount)
      return <ArrowDown size={16} className="text-red-600" />;
    return <Minus size={16} className="text-gray-300" />;
  };

  return (
    <div className="max-w-4xl mx-auto p-4 md:p-8 min-h-screen">
      {/* HEADER: Картинка лиги и Таймер */}
      {activeTab === "weekly" && currentLeague && (
        <div className="bg-gradient-to-r from-indigo-600 to-purple-700 rounded-2xl p-6 text-white mb-8 flex items-center justify-between shadow-lg relative overflow-hidden">
          {/* Декор */}
          <div className="absolute -right-10 -top-10 w-40 h-40 bg-white opacity-10 rounded-full blur-2xl"></div>

          <div className="flex items-center gap-6 relative z-10">
            <div className="w-20 h-20 bg-white/20 backdrop-blur-md rounded-full flex items-center justify-center border-2 border-white/30 shadow-inner">
              <img
                src={currentLeague.icon_url}
                alt={currentLeague.name}
                className="w-14 h-14 object-contain drop-shadow-md"
              />
            </div>
            <div>
              <h1 className="text-2xl font-bold">{currentLeague.name}</h1>
              <p className="text-indigo-100 text-sm opacity-90">
                Соревнуйтесь с другими учениками!
              </p>
            </div>
          </div>

          <div className="text-right relative z-10 hidden sm:block">
            <div className="flex items-center justify-end gap-2 text-indigo-200 text-sm mb-1">
              <Clock size={16} /> До окончания
            </div>
            <div className="text-3xl font-mono font-bold tracking-wider">
              {timeLeft}
            </div>
          </div>
        </div>
      )}

      {/* TABS */}
      <div className="flex gap-4 mb-6 bg-white p-1 rounded-xl shadow-sm border w-fit">
        <button
          onClick={() => setActiveTab("weekly")}
          className={`flex items-center gap-2 px-6 py-2 rounded-lg text-sm font-medium transition-all ${
            activeTab === "weekly"
              ? "bg-indigo-100 text-indigo-700 shadow-sm"
              : "text-gray-500 hover:bg-gray-50"
          }`}
        >
          <Trophy size={18} /> Лига
        </button>
        <button
          onClick={() => setActiveTab("global")}
          className={`flex items-center gap-2 px-6 py-2 rounded-lg text-sm font-medium transition-all ${
            activeTab === "global"
              ? "bg-indigo-100 text-indigo-700 shadow-sm"
              : "text-gray-500 hover:bg-gray-50"
          }`}
        >
          <Globe size={18} /> Мировой топ
        </button>
      </div>

      {/* TABLE */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {isLoading ? (
          <div className="p-10 text-center text-gray-500">
            Загрузка рейтинга...
          </div>
        ) : entries.length === 0 ? (
          <div className="p-10 text-center text-gray-500">
            Пока никого нет. Станьте первым!
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-xs text-gray-500 uppercase tracking-wider">
                  <th className="p-4 w-16 text-center">#</th>
                  <th className="p-4">Ученик</th>
                  <th className="p-4 text-right">XP</th>
                  {activeTab === "weekly" && <th className="p-4 w-10"></th>}
                </tr>
              </thead>
              <tbody>
                {entries.map((entry, index) => {
                  const isMe = entry.user_id === userProfileId;
                  const zoneClass = getZoneStyle(index, entries.length);

                  return (
                    <tr
                      key={entry.user_id}
                      className={`
                                        border-b border-gray-50 last:border-0 transition-colors
                                        ${
                                          isMe
                                            ? "bg-indigo-50/50"
                                            : "hover:bg-gray-50"
                                        }
                                        ${zoneClass}
                                    `}
                    >
                      <td className="p-4 text-center font-bold text-gray-600">
                        {index + 1}
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-full bg-gray-200 overflow-hidden border border-gray-100">
                            {entry.avatar_url ? (
                              <img
                                src={entry.avatar_url}
                                alt=""
                                className="w-full h-full object-cover"
                              />
                            ) : (
                              <div className="w-full h-full flex items-center justify-center text-gray-400 text-xs">
                                {entry.first_name[0]}
                              </div>
                            )}
                          </div>
                          <div>
                            <div
                              className={`font-medium ${
                                isMe ? "text-indigo-700" : "text-gray-900"
                              }`}
                            >
                              {entry.first_name} {entry.last_name}{" "}
                              {isMe && "(Вы)"}
                            </div>
                            <div className="text-xs text-gray-500">
                              Уровень {entry.level}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="p-4 text-right font-mono font-bold text-orange-500">
                        {entry.xp} XP
                      </td>
                      {activeTab === "weekly" && (
                        <td className="p-4 text-center">
                          {getZoneIcon(index, entries.length)}
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Legend */}
      {activeTab === "weekly" && (
        <div className="mt-4 flex gap-6 justify-center text-xs text-gray-500">
          <div className="flex items-center gap-2">
            <ArrowUp size={14} className="text-green-600" /> Повышение
          </div>
          <div className="flex items-center gap-2">
            <Minus size={14} className="text-gray-400" /> Сохранение места
          </div>
          <div className="flex items-center gap-2">
            <ArrowDown size={14} className="text-red-600" /> Понижение
          </div>
        </div>
      )}
    </div>
  );
};
