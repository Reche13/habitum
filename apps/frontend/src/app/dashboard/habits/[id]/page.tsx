"use client";

import { use } from "react";
import { useRouter } from "next/navigation";
import { format, eachDayOfInterval, isSameDay, subDays } from "date-fns";
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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import {
  getCategoryLabel,
  getCategoryIcon,
  getIconEmoji,
} from "@/lib/habit-utils";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { useState, useMemo } from "react";
import {
  useHabit,
  useMarkComplete,
  useUnmarkComplete,
  useDeleteHabit,
  useCompletionHistory,
} from "@/lib/hooks";
import { mapHabitResponseToHabit } from "@/lib/api/mappers";
import type { Habit } from "@/types/habit";

const generateHeatmapData = (habit: Habit) => {
  const endDate = new Date(2026, 6, 10);
  const startDate = subDays(endDate, 364); // Last 365 days
  const days = eachDayOfInterval({ start: startDate, end: endDate });

  // Group days by week - weeks start on Sunday (0)
  // First, find the first day and pad to Sunday if needed
  const firstDay = days[0];
  const firstDayOfWeek = firstDay.getDay(); // 0 = Sunday

  // Group days into weeks
  const weeks: Date[][] = [];
  let currentWeek: Date[] = [];

  // Pad the first week if it doesn't start on Sunday
  if (firstDayOfWeek !== 0) {
    for (let i = 0; i < firstDayOfWeek; i++) {
      currentWeek.push(new Date(0)); // Placeholder for empty days
    }
  }

  days.forEach((day) => {
    const dayOfWeek = day.getDay(); // 0 = Sunday, 6 = Saturday

    // If it's Sunday and we have a previous week, save it
    if (dayOfWeek === 0 && currentWeek.length > 0) {
      // Fill remaining days if week is incomplete
      while (currentWeek.length < 7) {
        currentWeek.push(new Date(0));
      }
      weeks.push(currentWeek);
      currentWeek = [];
    }

    currentWeek.push(day);
  });

  // Add the last week if it exists, and pad it
  if (currentWeek.length > 0) {
    while (currentWeek.length < 7) {
      currentWeek.push(new Date(0));
    }
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
  const { data: completionHistoryData } = useCompletionHistory(
    id,
    undefined,
    true
  );
  const markComplete = useMarkComplete();
  const unmarkComplete = useUnmarkComplete();
  const deleteHabit = useDeleteHabit();

  // Memoize completion history to prevent unnecessary re-renders
  // The API now returns unwrapped data, so completionHistoryData is the data object directly
  const completionHistoryKey = useMemo(() => {
    if (!completionHistoryData?.dates) return "";
    return JSON.stringify(
      Array.isArray(completionHistoryData.dates)
        ? completionHistoryData.dates
        : []
    );
  }, [completionHistoryData?.dates]);

  const completionHistoryArray = useMemo(() => {
    if (completionHistoryData?.dates) {
      return Array.isArray(completionHistoryData.dates)
        ? completionHistoryData.dates
        : [];
    }
    return [];
  }, [completionHistoryKey]);

  // Memoize habit object to prevent infinite re-renders
  // Only depend on habit ID and completion history key
  const habit = useMemo(() => {
    if (!habitResponse) return null;
    const mappedHabit = mapHabitResponseToHabit(habitResponse);
    return {
      ...mappedHabit,
      completionHistory:
        completionHistoryArray.length > 0
          ? completionHistoryArray
          : mappedHabit.completionHistory || [],
    };
  }, [habitResponse, completionHistoryKey]);

  // Get values directly from habitResponse to ensure they're fresh
  const currentStreak = habitResponse?.current_streak ?? 0;
  const longestStreak = habitResponse?.longest_streak ?? 0;
  const completionRate = habitResponse?.completionRate ?? 0;
  const completedThisWeek = habitResponse?.completedThisWeek ?? 0;

  // Memoize computed values to prevent unnecessary recalculations
  const weeks = useMemo(() => {
    if (!habit) return [];
    return generateHeatmapData(habit);
  }, [habit?.id]);

  // Use completion history directly from the array, not from habit object
  const isCompletedOnDate = useMemo(() => {
    const history = completionHistoryArray;
    const historySet = new Set(Array.isArray(history) ? history : []);
    return (date: Date): boolean => {
      const dateStr = format(date, "yyyy-MM-dd");
      return historySet.has(dateStr);
    };
  }, [completionHistoryKey]);

  // Get month labels - show month name for the first week of each month
  const monthLabels = useMemo(() => {
    const labels: { weekIdx: number; month: string }[] = [];
    let lastMonth = -1;

    weeks.forEach((week, weekIdx) => {
      if (week.length > 0) {
        // Find the first real day in the week (not placeholder)
        const firstRealDay = week.find(
          (d) => d instanceof Date && !isNaN(d.getTime())
        );
        if (firstRealDay) {
          const month = firstRealDay.getMonth();
          if (month !== lastMonth) {
            labels.push({ weekIdx, month: format(firstRealDay, "MMM") });
            lastMonth = month;
          }
        }
      }
    });
    return labels;
  }, [weeks]);

  if (isLoading) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12 flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !habitResponse || !habit) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
        <div className="max-w-4xl mx-auto text-center py-16">
          <h1 className="text-2xl font-semibold mb-2">Habit not found</h1>
          <p className="text-muted-foreground mb-6">
            {error instanceof Error
              ? error.message
              : "The habit you're looking for doesn't exist."}
          </p>
          <Button asChild>
            <Link href="/dashboard/habits">Back to Habits</Link>
          </Button>
        </div>
      </div>
    );
  }

  const handleDelete = async () => {
    try {
      await deleteHabit.mutateAsync(id);
      router.push("/dashboard/habits");
    } catch (error) {
      console.error("Failed to delete habit:", error);
      alert("Failed to delete habit. Please try again.");
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
            className="mb-4 cursor-pointer"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>

          <div className="flex items-start justify-between gap-4">
            <div className="flex items-start gap-4 flex-1">
              <div className="h-16 w-16 bg-muted rounded-xl flex items-center justify-center text-3xl shrink-0">
                {habit.icon || getIconEmoji(habit.iconId)}
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
                variant="default"
                onClick={handleComplete}
                disabled={markComplete.isPending || unmarkComplete.isPending}
                className={cn(
                  "cursor-pointer",
                  habit.completedToday
                    ? "bg-green-600 hover:bg-green-600/90"
                    : "bg-foreground"
                )}
              >
                {markComplete.isPending || unmarkComplete.isPending ? (
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
              <Button
                variant="outline"
                asChild
                className="border border-zinc-200"
              >
                <Link href={`/dashboard/habits/${id}/edit`}>
                  <Edit className="h-4 w-4 mr-2" />
                  Edit
                </Link>
              </Button>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button
                    variant="destructive"
                    className="cursor-pointer"
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
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Delete Habit</AlertDialogTitle>
                    <AlertDialogDescription>
                      Are you sure you want to delete "{habit.name}"? This
                      action cannot be undone and all completion history will be
                      permanently deleted.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel className="border border-zinc-200 cursor-pointer">
                      Cancel
                    </AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleDelete}
                      className="bg-destructive text-destructive-foreground hover:bg-destructive/90 cursor-pointer"
                    >
                      Delete
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div className="rounded-lg border bg-background p-4 border-zinc-200 shadow-xs">
            <div className="flex items-center gap-2 mb-1">
              <Flame className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">
                Current Streak
              </span>
            </div>
            <p className="text-2xl font-semibold">{currentStreak} days</p>
          </div>
          <div className="rounded-lg border bg-background p-4 border-zinc-200 shadow-xs">
            <div className="flex items-center gap-2 mb-1">
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">
                Longest Streak
              </span>
            </div>
            <p className="text-2xl font-semibold">{longestStreak} days</p>
          </div>
          <div className="rounded-lg border bg-background p-4 border-zinc-200 shadow-xs">
            <div className="flex items-center gap-2 mb-1">
              <Target className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">
                Completion Rate
              </span>
            </div>
            <p className="text-2xl font-semibold">
              {Math.round(completionRate)}%
            </p>
          </div>
          <div className="rounded-lg border bg-background p-4 border-zinc-200 shadow-xs">
            <div className="flex items-center gap-2 mb-1">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">This Week</span>
            </div>
            <p className="text-2xl font-semibold">
              {completedThisWeek}/
              {habit.frequency === "daily" ? 7 : habit.timesPerWeek || 0}
            </p>
          </div>
        </div>

        {/* Completion Heatmap */}
        <div className="rounded-lg border border-zinc-200 bg-background shadow-xs p-6 mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Completion History</h2>
            <div className="flex items-center gap-4 text-xs text-muted-foreground">
              <span>No activity</span>
              <div className="flex gap-1">
                <div className="w-3 h-3 rounded border border-zinc-200 bg-muted" />
                <div
                  className="w-3 h-3 rounded border border-zinc-200"
                  style={{ backgroundColor: habit.color }}
                />
              </div>
              <span>Completed</span>
            </div>
          </div>

          <div className="space-y-2 overflow-x-auto pb-4">
            {/* Month Labels */}
            <div className="flex gap-1 mb-3 min-w-fit px-1">
              <div className="w-10.25 shrink-0"></div>
              {weeks.map((week, weekIdx) => {
                const monthLabel = monthLabels.find(
                  (m) => m.weekIdx === weekIdx
                );
                return (
                  <div
                    key={weekIdx}
                    className="text-xs font-medium text-muted-foreground text-center w-2.75 shrink-0"
                  >
                    {monthLabel ? monthLabel.month : ""}
                  </div>
                );
              })}
            </div>

            {/* Day Labels and Heatmap */}
            <div className="flex gap-1.5 min-w-fit px-1">
              {/* Day labels */}
              <div className="space-y-1.5 pr-3 shrink-0">
                {["Sun", "Mon", "Tue", "Wed", "Thurs", "Fri", "Sat"].map(
                  (label, idx) => (
                    <div
                      key={idx}
                      className="text-[10px] text-muted-foreground/60 h-2.75 flex items-center font-medium"
                    >
                      {label}
                    </div>
                  )
                )}
              </div>

              {/* Heatmap cells */}
              <div className="flex gap-1">
                {weeks.map((week, weekIdx) => (
                  <div key={weekIdx} className="flex flex-col gap-1 shrink-0">
                    {week.map((day, dayIdx) => {
                      // Skip placeholder days
                      if (day.getTime() === 0) {
                        return (
                          <div
                            key={`${weekIdx}-${dayIdx}`}
                            className="w-2.75 h-2.75"
                          />
                        );
                      }

                      const isCompleted = isCompletedOnDate(day);
                      const isToday = isSameDay(day, new Date());
                      const isSelected =
                        selectedDate && isSameDay(day, selectedDate);

                      return (
                        <div
                          key={`${weekIdx}-${dayIdx}`}
                          onClick={() => setSelectedDate(day)}
                          className={cn(
                            "w-2.75 h-2.75 rounded-xs border box-border transition-all hover:scale-110 hover:z-10 relative cursor-pointer",
                            !isCompleted &&
                              "bg-muted border-muted-foreground/20",
                            isCompleted && "border-border",
                            isToday && "ring-2 ring-primary ring-offset-1",
                            isSelected &&
                              "ring-2 ring-green-600 ring-offset-1 shadow-lg"
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
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Selected Date Info */}
          {selectedDate && (
            <div className="mt-4 p-3 rounded-lg border border-zinc-200 shadow-xs bg-muted/50">
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-sm text-foreground/80">
                    {format(selectedDate, "EEEE, MMMM d, yyyy")}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {isCompletedOnDate(selectedDate)
                      ? "✓ Completed"
                      : "Not completed"}
                  </div>
                </div>
                <Button
                  variant="outline"
                  className="cursor-pointer border border-zinc-200"
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
        <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Completions</h2>
          {completionHistoryArray.length > 0 ? (
            <div className="space-y-2">
              {completionHistoryArray
                .slice()
                .reverse()
                .slice(0, 10)
                .map((dateStr) => {
                  const date = new Date(dateStr);
                  const isToday =
                    format(date, "yyyy-MM-dd") ===
                    format(new Date(), "yyyy-MM-dd");
                  return (
                    <div
                      key={dateStr}
                      className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
                    >
                      <div className="flex items-center gap-3">
                        <CheckCircle2 className="h-5 w-5 text-green-600 dark:text-green-400" />
                        <div>
                          <div className="font-medium text-sm text-foreground/80">
                            {isToday
                              ? "Today"
                              : format(date, "EEEE, MMMM d, yyyy")}
                          </div>
                          {isToday && (
                            <div className="text-xs text-muted-foreground">
                              {habit.completedTodayAt &&
                                format(
                                  new Date(habit.completedTodayAt),
                                  "h:mm a"
                                )}
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
