import { apiClient } from "./client";
import type { DashboardResponse } from "./types";

// Dashboard API
export const dashboardAPI = {
  // Get dashboard home data
  getHome: async () => {
    const response = await apiClient.get<DashboardResponse>("/dashboard/home");
    return response.data;
  },
};
