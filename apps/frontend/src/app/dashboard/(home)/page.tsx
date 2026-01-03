"use client";

import { format } from "date-fns";
import {
  CheckCircle2,
  Target,
  TrendingUp,
  Calendar,
  Trophy,
  Sparkles,
  Loader2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useDashboardHome, useMarkComplete } from "@/lib/hooks";
import { getIconEmoji } from "@/lib/habit-utils";
import { IconId } from "@/components/new-habit/icon-picker";

export default function Home() {
  const today = new Date();
  const { data: dashboard, isLoading, error } = useDashboardHome();
  const markComplete = useMarkComplete();

  const handleComplete = async (habitId: string) => {
    try {
      await markComplete.mutateAsync({ id: habitId });
    } catch (error) {
      console.error("Failed to mark habit as complete:", error);
    }
  };

  if (isLoading) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12 flex items-center justify-center min-h-100">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !dashboard) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
        <div className="text-center py-16">
          <p className="text-destructive">Failed to load dashboard data</p>
        </div>
      </div>
    );
  }

  const habitsToComplete = dashboard.habitsToComplete || [];
  const habitsCompleted = dashboard.habitsCompleted || [];
  const activeStreaks = dashboard.activeStreaks || [];
  const completionRate = dashboard.today.completionRate;
  const totalCompletedThisWeek = dashboard.quickStats.thisWeek;

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between gap-4 mb-2">
          <div>
            <h1 className="text-2xl sm:text-3xl font-semibold">
              Welcome back! 👋
            </h1>
            <p className="text-sm text-muted-foreground mt-1">
              {format(today, "EEEE, MMMM d, yyyy")}
            </p>
          </div>
          <div className="text-right">
            <div className="text-2xl font-semibold">
              {dashboard.today.completedCount}/{dashboard.today.totalCount}
            </div>
            <div className="text-xs text-muted-foreground">Completed today</div>
          </div>
        </div>

        {/* Today's Progress */}
        <div className="mt-6 rounded-lg border border-zinc-200 shadow-xs bg-background p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Today's Progress</h2>
            <span className="text-2xl font-bold">
              {Math.round(completionRate)}%
            </span>
          </div>
          <div className="w-full bg-muted rounded-full h-3">
            <div
              className="bg-primary h-3 rounded-full transition-all duration-500"
              style={{ width: `${completionRate}%` }}
            />
          </div>
          {completionRate === 100 && (
            <div className="mt-4 flex items-center gap-2 text-green-600 dark:text-green-400">
              <Sparkles className="h-4 w-4" />
              <span className="text-sm font-medium">
                Amazing! You've completed all habits today! 🎉
              </span>
            </div>
          )}
        </div>
      </div>

      {/* To Complete Today */}
      {habitsToComplete.length > 0 && (
        <section className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <Target className="h-5 w-5" />
              To Complete Today
            </h2>
            <span className="text-sm text-muted-foreground">
              {habitsToComplete.length} remaining
            </span>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {habitsToComplete.map((habit) => (
              <HabitQuickCard
                key={habit.id}
                habit={habit}
                onComplete={handleComplete}
                isLoading={markComplete.isPending}
              />
            ))}
          </div>
        </section>
      )}

      {/* Active Streaks */}
      {activeStreaks.length > 0 && (
        <section className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Active Streaks
            </h2>
            <Link
              href="/dashboard/habits"
              className="text-sm text-primary hover:underline"
            >
              View all
            </Link>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {activeStreaks.map((streak) => (
              <div
                key={streak.id}
                className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4 flex items-center gap-4"
              >
                <div className="h-12 w-12 rounded-xl flex items-center justify-center text-2xl shrink-0 bg-muted">
                  {streak.icon || "🔥"}
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="font-medium truncate">{streak.name}</h3>
                  <div className="flex items-center gap-1.5 mt-1">
                    <span className="text-lg">🔥</span>
                    <span className="text-sm font-medium">
                      {streak.currentStreak} day
                      {streak.currentStreak !== 1 ? "s" : ""}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Achievements
      {dashboard.recentAchievements &&
        dashboard.recentAchievements.length > 0 && (
          <section className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold flex items-center gap-2">
                <Trophy className="h-5 w-5" />
                Recent Achievements
              </h2>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {dashboard.recentAchievements.map((achievement) => (
                <div
                  key={achievement.id}
                  className="rounded-lg border bg-background p-4 flex items-center gap-4"
                >
                  <div className="h-12 w-12 rounded-xl flex items-center justify-center text-2xl shrink-0 bg-primary/10">
                    {achievement.icon}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium">{achievement.name}</h3>
                    <p className="text-sm text-muted-foreground">
                      {achievement.description}
                    </p>
                    {achievement.progress !== undefined && (
                      <div className="mt-2">
                        <div className="w-full bg-muted rounded-full h-1.5">
                          <div
                            className="bg-primary h-1.5 rounded-full"
                            style={{ width: `${achievement.progress}%` }}
                          />
                        </div>
                        <span className="text-xs text-muted-foreground mt-1">
                          {achievement.progress}% complete
                        </span>
                      </div>
                    )}
                    {achievement.unlockedAt && (
                      <span className="text-xs text-muted-foreground">
                        Unlocked{" "}
                        {format(new Date(achievement.unlockedAt), "MMM d")}
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </section>
        )} */}

      {/* Completed Today */}
      {habitsCompleted.length > 0 && (
        <section className="mb-8">
          <details className="group">
            <summary className="flex items-center justify-between cursor-pointer list-none mb-4">
              <h2 className="text-xl font-semibold flex items-center gap-2">
                <CheckCircle2 className="h-5 w-5 text-green-600 dark:text-green-400" />
                Completed Today
              </h2>
              <span className="text-sm text-muted-foreground">
                {habitsCompleted.length} completed
              </span>
            </summary>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
              {habitsCompleted.map((habit) => (
                <div
                  key={habit.id}
                  className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4 flex items-center gap-4 opacity-75"
                >
                  <div className="h-12 w-12 rounded-xl flex items-center justify-center text-2xl shrink-0 bg-muted">
                    {habit.icon || "🔥"}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <h3 className="font-medium truncate">{habit.name}</h3>
                      <CheckCircle2 className="h-4 w-4 text-green-600 dark:text-green-400 shrink-0" />
                    </div>
                    {habit.completedTodayAt && (
                      <p className="text-xs text-muted-foreground mt-1">
                        Completed at{" "}
                        {format(new Date(habit.completedTodayAt), "h:mm a")}
                      </p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </details>
        </section>
      )}

      {/* Quick Stats */}
      <section>
        <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
          <Calendar className="h-5 w-5" />
          Quick Stats
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="text-sm text-muted-foreground mb-1">
              Today's Rate
            </div>
            <div className="text-2xl font-semibold">
              {Math.round(completionRate)}%
            </div>
          </div>
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="text-sm text-muted-foreground mb-1">This Week</div>
            <div className="text-2xl font-semibold">
              {totalCompletedThisWeek}
            </div>
          </div>
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="text-sm text-muted-foreground mb-1">
              Longest Streak
            </div>
            <div className="text-2xl font-semibold">
              {dashboard.quickStats.longestStreak}
            </div>
          </div>
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="text-sm text-muted-foreground mb-1">
              Total Habits
            </div>
            <div className="text-2xl font-semibold">
              {dashboard.quickStats.totalHabits}
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}

interface HabitQuickCardProps {
  habit: {
    id: string;
    name: string;
    description?: string;
    icon?: string;
    iconId?: string;
    color?: string;
    currentStreak: number;
  };
  onComplete: (habitId: string) => void;
  isLoading?: boolean;
}

function HabitQuickCard({ habit, onComplete, isLoading }: HabitQuickCardProps) {
  const icon = habit.icon || getIconEmoji((habit.iconId as IconId) || "fire");
  const color = habit.color || "#6366f1";

  return (
    <div className="rounded-lg border border-zinc-200 bg-background shadow-xs p-4 hover:shadow-md transition-shadow">
      <div className="flex items-start gap-4">
        <div
          className="h-12 w-12 rounded-xl flex items-center justify-center text-2xl shrink-0"
          style={{
            backgroundColor: color + "20",
            color: color,
          }}
        >
          {icon}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold truncate">{habit.name}</h3>
          {habit.description && (
            <p className="text-sm text-muted-foreground line-clamp-2 mt-1">
              {habit.description}
            </p>
          )}
          {habit.currentStreak > 0 && (
            <div className="flex items-center gap-1.5 mt-2">
              <span className="text-sm">🔥</span>
              <span className="text-xs text-muted-foreground">
                {habit.currentStreak} day streak
              </span>
            </div>
          )}
        </div>
      </div>
      <Button
        onClick={() => onComplete(habit.id)}
        className="w-full mt-4"
        size="sm"
        disabled={isLoading}
        variant="default"
      >
        {isLoading ? (
          <>
            <Loader2 className="h-4 w-4 mr-2 animate-spin" />
            Marking...
          </>
        ) : (
          <>
            <CheckCircle2 className="h-4 w-4 mr-2" />
            Mark Complete
          </>
        )}
      </Button>
    </div>
  );
}
