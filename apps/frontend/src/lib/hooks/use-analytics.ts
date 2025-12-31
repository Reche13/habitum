import { useQuery } from "@tanstack/react-query";
import { analyticsAPI } from "@/lib/api/analytics";

// Query keys
export const analyticsKeys = {
  all: ["analytics"] as const,
  completionTrend: (period: string) => [...analyticsKeys.all, "completion-trend", period] as const,
  categoryBreakdown: () => [...analyticsKeys.all, "category-breakdown"] as const,
  dayOfWeek: (period?: string) => [...analyticsKeys.all, "day-of-week", period] as const,
  metrics: () => [...analyticsKeys.all, "metrics"] as const,
  topHabits: (limit: number, sortBy: string) => [...analyticsKeys.all, "top-habits", limit, sortBy] as const,
  streakLeaderboard: (limit: number) => [...analyticsKeys.all, "streak-leaderboard", limit] as const,
  insights: () => [...analyticsKeys.all, "insights"] as const,
};

// Get completion trend
export function useCompletionTrend(period: "7d" | "30d" | "90d" | "all" = "30d") {
  return useQuery({
    queryKey: analyticsKeys.completionTrend(period),
    queryFn: () => analyticsAPI.getCompletionTrend(period),
  });
}

// Get category breakdown
export function useCategoryBreakdown() {
  return useQuery({
    queryKey: analyticsKeys.categoryBreakdown(),
    queryFn: () => analyticsAPI.getCategoryBreakdown(),
  });
}

// Get day of week analysis
export function useDayOfWeekAnalysis(period?: string) {
  return useQuery({
    queryKey: analyticsKeys.dayOfWeek(period),
    queryFn: () => analyticsAPI.getDayOfWeekAnalysis(period),
  });
}

// Get metrics
export function useMetrics() {
  return useQuery({
    queryKey: analyticsKeys.metrics(),
    queryFn: () => analyticsAPI.getMetrics(),
  });
}

// Get top habits
export function useTopHabits(limit: number = 10, sortBy: "completion" | "streak" = "completion") {
  return useQuery({
    queryKey: analyticsKeys.topHabits(limit, sortBy),
    queryFn: () => analyticsAPI.getTopHabits(limit, sortBy),
  });
}

// Get streak leaderboard
export function useStreakLeaderboard(limit: number = 10) {
  return useQuery({
    queryKey: analyticsKeys.streakLeaderboard(limit),
    queryFn: () => analyticsAPI.getStreakLeaderboard(limit),
  });
}

// Get insights
export function useInsights() {
  return useQuery({
    queryKey: analyticsKeys.insights(),
    queryFn: () => analyticsAPI.getInsights(),
  });
}


