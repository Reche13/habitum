import { CategoryId } from "@/components/new-habit/select-category";
import { IconId } from "@/components/new-habit/icon-picker";

export interface Habit {
  id: string;
  name: string;
  description?: string;
  icon: string;
  iconId: IconId;
  color: string;
  frequency: "daily" | "weekly";
  timesPerWeek?: number;
  category?: CategoryId;
  createdAt?: string;
  currentStreak?: number;
  longestStreak?: number;
  completionRate?: number;
  completedToday?: boolean;
  completedTodayAt?: string;
  completedThisWeek?: number;
  completionHistory?: string[];
  archivedAt?: string;
}

export interface Achievement {
  id: string;
  name: string;
  description: string;
  icon: string;
  unlockedAt?: string; // ISO timestamp
  progress?: number; // 0-100
  target?: number; // target value for progress
}
