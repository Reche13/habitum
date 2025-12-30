"use client";

import { use } from "react";
import { useRouter } from "next/navigation";
import {
  format,
  startOfYear,
  endOfYear,
  eachDayOfInterval,
  startOfWeek,
  endOfWeek,
  isSameYear,
  isSameDay,
  subDays,
} from "date-fns";
import {
  ArrowLeft,
  Edit,
  Trash2,
  CheckCircle2,
  Calendar,
  TrendingUp,
  Target,
  Flame,
  Loader2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { getCategoryLabel, getCategoryIcon } from "@/lib/habit-utils";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { useState } from "react";
import { useHabit, useMarkComplete, useUnmarkComplete, useDeleteHabit, useCompletionHistory } from "@/lib/hooks";
import { mapHabitResponseToHabit } from "@/lib/api/mappers";

// Generate heatmap data for the last year
const generateHeatmapData = (habit: Habit) => {
  const endDate = new Date();
  const startDate = subDays(endDate, 364); // Last 365 days
  const days = eachDayOfInterval({ start: startDate, end: endDate });

  // Group days by week (starting Sunday)
  const weeks: Date[][] = [];
  let currentWeek: Date[] = [];
  let lastWeekDay = -1;

  days.forEach((day) => {
    const dayOfWeek = day.getDay(); // 0 = Sunday, 6 = Saturday

    if (dayOfWeek === 0 && currentWeek.length > 0) {
      weeks.push(currentWeek);
      currentWeek = [];
    }

    currentWeek.push(day);
    lastWeekDay = dayOfWeek;
  });

  if (currentWeek.length > 0) {
    weeks.push(currentWeek);
  }

  return weeks;
};

export default function HabitDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  
  const { data: habitResponse, isLoading, error } = useHabit(id);
  const { data: completionHistoryData } = useCompletionHistory(id, undefined, true);
  const markComplete = useMarkComplete();
  const unmarkComplete = useUnmarkComplete();
  const deleteHabit = useDeleteHabit();

  if (isLoading) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12 flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !habitResponse) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
        <div className="max-w-4xl mx-auto text-center py-16">
          <h1 className="text-2xl font-semibold mb-2">Habit not found</h1>
          <p className="text-muted-foreground mb-6">
            {error instanceof Error ? error.message : "The habit you're looking for doesn't exist."}
          </p>
          <Button asChild>
            <Link href="/dashboard/habits">Back to Habits</Link>
          </Button>
        </div>
      </div>
    );
  }

  const habit = mapHabitResponseToHabit(habitResponse);
  
  // Use completion history from API if available
  if (completionHistoryData?.data) {
    habit.completionHistory = completionHistoryData.data;
  }

  const weeks = generateHeatmapData(habit);
  const streak = habit.currentStreak ?? 0;
  const longestStreak = habit.longestStreak ?? 0;

  const isCompletedOnDate = (date: Date): boolean => {
    const dateStr = format(date, "yyyy-MM-dd");
    return habit.completionHistory?.includes(dateStr) || false;
  };


  // Get month labels
  const monthLabels: { weekIdx: number; month: string }[] = [];
  let lastMonth = -1;

  weeks.forEach((week, weekIdx) => {
    if (week.length > 0) {
      const firstDay = week[0];
      const month = firstDay.getMonth();
      if (month !== lastMonth) {
        monthLabels.push({ weekIdx, month: format(firstDay, "MMM") });
        lastMonth = month;
      }
    }
  });

  const handleDelete = async () => {
    if (confirm("Are you sure you want to delete this habit? This action cannot be undone.")) {
      try {
        await deleteHabit.mutateAsync(id);
        router.push("/dashboard/habits");
      } catch (error) {
        console.error("Failed to delete habit:", error);
        alert("Failed to delete habit. Please try again.");
      }
    }
  };

  const handleComplete = async () => {
    try {
      if (habit.completedToday) {
        await unmarkComplete.mutateAsync({ id: habit.id });
      } else {
        await markComplete.mutateAsync({ id: habit.id });
      }
    } catch (error) {
      console.error("Failed to toggle completion:", error);
      alert("Failed to update habit. Please try again.");
    }
  };

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Button
            variant="ghost"
            onClick={() => router.back()}
            className="mb-4"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>

          <div className="flex items-start justify-between gap-4">
            <div className="flex items-start gap-4 flex-1">
              <div
                className="h-16 w-16 rounded-xl flex items-center justify-center text-3xl shrink-0"
                style={{ backgroundColor: habit.color + "20", color: habit.color }}
              >
                {habit.icon}
              </div>
              <div className="flex-1 min-w-0">
                <h1 className="text-3xl font-semibold mb-2">{habit.name}</h1>
                {habit.description && (
                  <p className="text-muted-foreground">{habit.description}</p>
                )}
                <div className="flex items-center gap-3 mt-3 flex-wrap">
                  <span className="text-sm text-muted-foreground">
                    {habit.frequency === "daily"
                      ? "Daily"
                      : `${habit.timesPerWeek}× / week`}
                  </span>
                  {habit.category && (
                    <>
                      <span className="text-sm text-muted-foreground">•</span>
                      <span className="text-sm rounded-full bg-muted px-2.5 py-1 flex items-center gap-1.5">
                        <span>{getCategoryIcon(habit.category)}</span>
                        {getCategoryLabel(habit.category)}
                      </span>
                    </>
                  )}
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2 shrink-0">
              <Button
                variant={habit.completedToday ? "outline" : "default"}
                onClick={handleComplete}
                disabled={markComplete.isPending || unmarkComplete.isPending}
                className={cn(
                  habit.completedToday && "border-green-600 text-green-600"
                )}
              >
                {(markComplete.isPending || unmarkComplete.isPending) ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Updating...
                  </>
                ) : habit.completedToday ? (
                  <>
                    <CheckCircle2 className="h-4 w-4 mr-2" />
                    Completed Today
                  </>
                ) : (
                  <>
                    <CheckCircle2 className="h-4 w-4 mr-2" />
                    Mark Complete
                  </>
                )}
              </Button>
              <Button variant="outline" asChild>
                <Link href={`/dashboard/habits/${id}/edit`}>
                  <Edit className="h-4 w-4 mr-2" />
                  Edit
                </Link>
              </Button>
              <Button
                variant="destructive"
                onClick={handleDelete}
                disabled={deleteHabit.isPending}
              >
                {deleteHabit.isPending ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Deleting...
                  </>
                ) : (
                  <>
                    <Trash2 className="h-4 w-4 mr-2" />
                    Delete
                  </>
                )}
              </Button>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div className="rounded-lg border bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <Flame className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Current Streak</span>
            </div>
            <p className="text-2xl font-semibold">{streak} days</p>
          </div>
          <div className="rounded-lg border bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Longest Streak</span>
            </div>
            <p className="text-2xl font-semibold">{longestStreak} days</p>
          </div>
          <div className="rounded-lg border bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <Target className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Completion Rate</span>
            </div>
            <p className="text-2xl font-semibold">{habit.completionRate}%</p>
          </div>
          <div className="rounded-lg border bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">This Week</span>
            </div>
            <p className="text-2xl font-semibold">
              {habit.completedThisWeek || 0}/{habit.frequency === "daily" ? 7 : habit.timesPerWeek}
            </p>
          </div>
        </div>

        {/* Completion Heatmap */}
        <div className="rounded-lg border bg-background p-6 mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Completion History</h2>
              <div className="flex items-center gap-4 text-xs text-muted-foreground">
              <span>No activity</span>
              <div className="flex gap-1">
                <div className="w-3 h-3 rounded border bg-muted" />
                <div
                  className="w-3 h-3 rounded border"
                  style={{ backgroundColor: habit.color }}
                />
              </div>
              <span>Completed</span>
            </div>
          </div>

          <div className="space-y-2 overflow-x-auto">
            {/* Month Labels */}
            <div className="flex gap-1 mb-2 min-w-fit">
              <div className="w-12 shrink-0"></div>
              {weeks.map((week, weekIdx) => {
                const monthLabel = monthLabels.find(
                  (m) => m.weekIdx === weekIdx
                );
                return (
                  <div
                    key={weekIdx}
                    className="text-xs text-muted-foreground text-center w-3 shrink-0"
                  >
                    {monthLabel ? monthLabel.month : ""}
                  </div>
                );
              })}
            </div>

            {/* Day Labels and Heatmap */}
            <div className="flex gap-1 min-w-fit">
              {/* Day labels */}
              <div className="space-y-1 pr-2 shrink-0">
                {["", "Mon", "", "Wed", "", "Fri", ""].map((label, idx) => (
                  <div
                    key={idx}
                    className="text-xs text-muted-foreground h-3 flex items-center"
                  >
                    {label}
                  </div>
                ))}
              </div>

              {/* Heatmap cells */}
              {weeks.map((week, weekIdx) => (
                <div key={weekIdx} className="flex flex-col gap-1 shrink-0">
                  {week.map((day, dayIdx) => {
                    const isCompleted = isCompletedOnDate(day);
                    const isToday = isSameDay(day, new Date());
                    const isSelected =
                      selectedDate && isSameDay(day, selectedDate);

                    return (
                      <div
                        key={`${weekIdx}-${dayIdx}`}
                        onClick={() => setSelectedDate(day)}
                        className={cn(
                          "w-3 h-3 rounded border transition-all hover:scale-125 hover:z-10 relative cursor-pointer",
                          !isCompleted && "bg-muted",
                          isToday && "ring-2 ring-primary",
                          isSelected && "ring-2 ring-blue-500 ring-offset-1"
                        )}
                        style={
                          isCompleted
                            ? { backgroundColor: habit.color }
                            : undefined
                        }
                        title={`${format(day, "MMM d, yyyy")}: ${
                          isCompleted ? "Completed" : "Not completed"
                        }`}
                      />
                    );
                  })}
                  {/* Fill remaining days if week is incomplete */}
                  {week.length < 7 &&
                    Array.from({ length: 7 - week.length }).map((_, idx) => (
                      <div key={`empty-${idx}`} className="w-3 h-3" />
                    ))}
                </div>
              ))}
            </div>
          </div>

          {/* Selected Date Info */}
          {selectedDate && (
            <div className="mt-4 p-3 rounded-lg border bg-muted/50">
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium">
                    {format(selectedDate, "EEEE, MMMM d, yyyy")}
                  </div>
                  <div className="text-sm text-muted-foreground">
                    {isCompletedOnDate(selectedDate)
                      ? "✓ Completed"
                      : "Not completed"}
                  </div>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setSelectedDate(null)}
                >
                  Close
                </Button>
              </div>
            </div>
          )}
        </div>

        {/* Recent Completions */}
        <div className="rounded-lg border bg-background p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Completions</h2>
          {habit.completionHistory && habit.completionHistory.length > 0 ? (
            <div className="space-y-2">
              {habit.completionHistory
                .slice()
                .reverse()
                .slice(0, 10)
                .map((dateStr) => {
                  const date = new Date(dateStr);
                  const isToday = format(date, "yyyy-MM-dd") === format(new Date(), "yyyy-MM-dd");
                  return (
                    <div
                      key={dateStr}
                      className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
                    >
                      <div className="flex items-center gap-3">
                        <CheckCircle2 className="h-5 w-5 text-green-600 dark:text-green-400" />
                        <div>
                          <div className="font-medium">
                            {isToday ? "Today" : format(date, "EEEE, MMMM d, yyyy")}
                          </div>
                          {isToday && (
                            <div className="text-sm text-muted-foreground">
                              {habit.completedTodayAt &&
                                format(new Date(habit.completedTodayAt), "h:mm a")}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
            </div>
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              <p>No completions yet. Start tracking to see your progress!</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

