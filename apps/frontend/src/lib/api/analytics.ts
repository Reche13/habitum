import { apiClient } from "./client";
import type {
  CompletionTrendData,
  CategoryBreakdownData,
  DayOfWeekData,
  MetricsData,
} from "./types";

// Analytics API
export const analyticsAPI = {
  // Get completion trend
  getCompletionTrend: async (period: "7d" | "30d" | "90d" | "all" = "30d") => {
    const response = await apiClient.get<{ data: CompletionTrendData[] }>(
      `/analytics/completion-trend?period=${period}`
    );
    return response.data;
  },

  // Get category breakdown
  getCategoryBreakdown: async () => {
    const response = await apiClient.get<{ data: CategoryBreakdownData[] }>(
      "/analytics/category-breakdown"
    );
    return response.data;
  },

  // Get day of week analysis
  getDayOfWeekAnalysis: async (period?: string) => {
    const url = period
      ? `/analytics/day-of-week?period=${period}`
      : "/analytics/day-of-week";
    const response = await apiClient.get<{ data: DayOfWeekData[] }>(url);
    return response.data;
  },

  // Get metrics
  getMetrics: async () => {
    const response = await apiClient.get<MetricsData>("/analytics/metrics");
    return response.data;
  },

  // Get top habits
  getTopHabits: async (limit: number = 10, sortBy: "completion" | "streak" = "completion") => {
    const response = await apiClient.get<{ data: any[] }>(
      `/analytics/top-habits?limit=${limit}&sortBy=${sortBy}`
    );
    return response.data;
  },

  // Get streak leaderboard
  getStreakLeaderboard: async (limit: number = 10) => {
    const response = await apiClient.get<{ data: any[] }>(
      `/analytics/streak-leaderboard?limit=${limit}`
    );
    return response.data;
  },

  // Get insights
  getInsights: async () => {
    const response = await apiClient.get<{ data: any[] }>("/analytics/insights");
    return response.data;
  },
};

