// API Response wrapper (matches backend APIResponse structure)
export interface APIResponse<T> {
  data: T;
  message?: string;
  success?: boolean;
}

// Pagination response
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page?: number;
  limit?: number;
}

// Habit types (matching backend)
export interface HabitResponse {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  icon?: string;
  iconId?: string;
  color?: string;
  category: string;
  frequency: "daily" | "weekly";
  times_per_week?: number;
  current_streak: number;
  longest_streak: number;
  createdAt: string;
  updatedAt: string;
  archived_at?: string;
  // Computed fields
  completionRate: number;
  completedToday: boolean;
  completedTodayAt?: string;
  completedThisWeek: number;
  completionHistory?: string[];
}

export interface CreateHabitPayload {
  name: string;
  description?: string;
  icon?: string;
  color?: string;
  category: string;
  frequency: "daily" | "weekly";
  times_per_week?: number;
}

export interface UpdateHabitPayload {
  name?: string;
  description?: string;
  icon?: string;
  color?: string;
  category?: string;
  frequency?: "daily" | "weekly";
  times_per_week?: number;
}

export interface HabitFilters {
  category?: string;
  search?: string;
  sort?: "name" | "date" | "streak" | "completion";
  order?: "asc" | "desc";
  page?: number;
  limit?: number;
}

// Analytics types
export interface CompletionTrendData {
  date: string;
  completions: number;
  totalHabits: number;
  completionRate: number;
}

export interface CategoryBreakdownData {
  category: string;
  label: string;
  habitCount: number;
  avgCompletionRate: number;
  totalCompletions: number;
}

export interface DayOfWeekData {
  day: string;
  dayIndex: number;
  completions: number;
  totalHabits: number;
  completionRate: number;
}

export interface MetricsData {
  avgCompletionRate: number;
  avgStreak: number;
  totalCompletions: number;
  consistencyScore: number;
}

// Calendar types
export interface CalendarCompletionDay {
  date: string;
  habits: Array<{
    id: string;
    name: string;
    color?: string;
    icon?: string;
  }>;
  completionRate: number;
  totalHabits: number;
  completedHabits: number;
}

export interface CalendarCompletionsResponse {
  completions: CalendarCompletionDay[];
  statistics: {
    totalCompletions: number;
    daysWithCompletions: number;
    completionRate: number;
    totalDays: number;
  };
}

// Dashboard types
export interface DashboardResponse {
  today: {
    date: string;
    completionRate: number;
    completedCount: number;
    totalCount: number;
  };
  habitsToComplete: Array<{
    id: string;
    name: string;
    description?: string;
    icon?: string;
    iconId?: string;
    color?: string;
    frequency: string;
    category: string;
    currentStreak: number;
    completedToday: boolean;
    completedTodayAt?: string;
  }>;
  habitsCompleted: Array<{
    id: string;
    name: string;
    icon?: string;
    color?: string;
    completedToday: boolean;
    completedTodayAt?: string;
  }>;
  activeStreaks: Array<{
    id: string;
    name: string;
    icon?: string;
    color?: string;
    currentStreak: number;
    longestStreak: number;
  }>;
  quickStats: {
    todayRate: number;
    thisWeek: number;
    longestStreak: number;
    totalHabits: number;
  };
  recentAchievements?: Array<{
    id: string;
    name: string;
    description: string;
    icon: string;
    unlockedAt?: string;
  }>;
}


