import api from "./axios";

export interface League {
  id: number;
  slug: string;
  name: string;
  order_index: number;
  icon_url: string;
}

export interface LeaderboardEntry {
  rank: number;
  user_id: string;
  first_name: string;
  last_name: string;
  avatar_url: string;
  xp: number;
  level: number;
  league_id?: number;
}

export interface LeaderboardResponse {
  leaderboard: LeaderboardEntry[];
  user_rank?: number;
}

export const gamificationApi = {
  getAllLeagues: async (): Promise<League[]> => {
    const response = await api.get<{ leagues: League[] }>(
      "/gamification/leagues"
    );
    return response.data.leagues;
  },

  getWeekly: async (limit: number = 50): Promise<LeaderboardResponse> => {
    const response = await api.get<LeaderboardResponse>(
      `/leaderboard/weekly?limit=${limit}`
    );
    return response.data;
  },

  getGlobal: async (limit: number = 50): Promise<LeaderboardResponse> => {
    const response = await api.get<LeaderboardResponse>(
      `/leaderboard/global?limit=${limit}`
    );
    return response.data;
  },
};
