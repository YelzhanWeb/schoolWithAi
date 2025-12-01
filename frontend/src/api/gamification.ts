import api from "./axios";

export interface League {
  id: number;
  slug: string;
  name: string;
  order_index: number;
  icon_url: string;
}

export const gamificationApi = {
  getAllLeagues: async (): Promise<League[]> => {
    const response = await api.get<{ leagues: League[] }>(
      "/gamification/leagues"
    );
    return response.data.leagues;
  },
};
