import React from "react";
import { Trophy, Star, Flame } from "lucide-react";
import type { StudentProfile } from "../../api/student";

interface Props {
  profile: StudentProfile;
  leagueName: string;
  leagueIcon?: string;
}

export const GamificationStats: React.FC<Props> = ({
  profile,
  leagueName,
  leagueIcon,
}) => {
  // Вычисляем прогресс до следующего уровня (условно 1000 XP на уровень)
  const xpForNextLevel = 1000;
  const progressPercent =
    ((profile.xp % xpForNextLevel) / xpForNextLevel) * 100;

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex flex-col md:flex-row justify-between items-center gap-6">
      {/* Уровень и XP */}
      <div className="flex items-center gap-4 w-full md:w-auto">
        <div className="relative">
          <div className="w-16 h-16 rounded-full bg-indigo-100 flex items-center justify-center text-2xl font-bold text-indigo-600 border-4 border-white shadow-sm">
            {profile.level}
          </div>
          <div className="absolute -bottom-2 -right-2 bg-indigo-600 text-white text-xs px-2 py-0.5 rounded-full">
            LVL
          </div>
        </div>
        <div className="flex-1">
          <div className="text-sm text-gray-500 font-medium">Ваш прогресс</div>
          <div className="text-xl font-bold text-gray-800">{profile.xp} XP</div>
          <div className="w-32 h-2 bg-gray-100 rounded-full mt-2 overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full"
              style={{ width: `${progressPercent}%` }}
            />
          </div>
        </div>
      </div>

      {/* Разделитель */}
      <div className="hidden md:block w-px h-12 bg-gray-200"></div>

      {/* Статистика */}
      <div className="flex gap-8 w-full md:w-auto justify-around md:justify-start">
        <div className="text-center">
          <div className="flex items-center justify-center w-10 h-10 bg-orange-100 text-orange-500 rounded-full mb-2 mx-auto">
            <Flame size={20} />
          </div>
          <div className="font-bold text-gray-800">
            {profile.current_streak}
          </div>
          <div className="text-xs text-gray-500">Дней подряд</div>
        </div>

        {/* ЛИГА С ИКОНКОЙ */}
        <div className="text-center flex flex-col items-center">
          {leagueIcon ? (
            <img
              src={leagueIcon}
              alt={leagueName}
              className="w-10 h-10 mb-2 object-contain"
            />
          ) : (
            <div className="flex items-center justify-center w-10 h-10 bg-yellow-100 text-yellow-600 rounded-full mb-2">
              <Trophy size={20} />
            </div>
          )}
          <div className="font-bold text-gray-800">{leagueName}</div>
          <div className="text-xs text-gray-500">Текущая лига</div>
        </div>

        <div className="text-center">
          <div className="flex items-center justify-center w-10 h-10 bg-purple-100 text-purple-500 rounded-full mb-2 mx-auto">
            <Star size={20} />
          </div>
          <div className="font-bold text-gray-800">{profile.weekly_xp}</div>
          <div className="text-xs text-gray-500">XP за неделю</div>
        </div>
      </div>
    </div>
  );
};
