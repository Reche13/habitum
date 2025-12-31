import type { HabitResponse } from "./types";
import type { Habit } from "@/types/habit";

// Map backend HabitResponse to frontend Habit type
export function mapHabitResponseToHabit(response: HabitResponse): Habit {
  return {
    id: response.id,
    name: response.name,
    description: response.description,
    icon: response.icon || "",
    iconId: (response.iconId || "fire") as any, // Use iconId if available, fallback to "fire"
    color: response.color || "#6366f1",
    frequency: response.frequency,
    timesPerWeek: response.times_per_week,
    category: response.category as any,
    createdAt: response.createdAt,
    currentStreak: response.current_streak,
    longestStreak: response.longest_streak,
    completionRate: response.completionRate,
    completedToday: response.completedToday,
    completedTodayAt: response.completedTodayAt,
    completedThisWeek: response.completedThisWeek,
    completionHistory: response.completionHistory,
    archivedAt: response.archived_at,
  };
}

// Map array of HabitResponse to Habit[]
export function mapHabitResponsesToHabits(responses: HabitResponse[]): Habit[] {
  return responses.map(mapHabitResponseToHabit);
}

