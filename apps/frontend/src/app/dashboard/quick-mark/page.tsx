"use client";

import { useState, useEffect } from "react";
import { CheckCircle2, Loader2, X, Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useHabits, useMarkComplete, useUnmarkComplete } from "@/lib/hooks";
import { getIconEmoji } from "@/lib/habit-utils";
import { format } from "date-fns";
import { mapHabitResponsesToHabits } from "@/lib/api/mappers";
import { cn } from "@/lib/utils";

export default function QuickMarkPage() {
  const today = new Date();
  const { data: habitsData, isLoading: habitsLoading } = useHabits();
  const markComplete = useMarkComplete();
  const unmarkComplete = useUnmarkComplete();

  const habits = habitsData?.data
    ? mapHabitResponsesToHabits(habitsData.data)
    : [];
  // Filter out archived habits - archivedAt is optional, so check if it exists
  const activeHabits = habits.filter((h) => !h.archivedAt);

  // Initialize completedIds from API data
  const initialCompletedIds = new Set(
    activeHabits.filter((h) => h.completedToday).map((h) => h.id)
  );
  const [completedIds, setCompletedIds] =
    useState<Set<string>>(initialCompletedIds);

  // Update completedIds when habits data changes
  useEffect(() => {
    const newCompletedIds = new Set(
      activeHabits.filter((h) => h.completedToday).map((h) => h.id)
    );
    setCompletedIds(newCompletedIds);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [habitsData]);

  const handleComplete = async (habitId: string) => {
    try {
      await markComplete.mutateAsync({ id: habitId });
      setCompletedIds((prev) => new Set(prev).add(habitId));
    } catch (error) {
      console.error("Failed to mark habit as complete:", error);
    }
  };

  const handleUncomplete = async (habitId: string) => {
    try {
      await unmarkComplete.mutateAsync({ id: habitId });
      setCompletedIds((prev) => {
        const newSet = new Set(prev);
        newSet.delete(habitId);
        return newSet;
      });
    } catch (error) {
      console.error("Failed to unmark habit:", error);
    }
  };

  if (habitsLoading) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12 flex items-center justify-center min-h-100">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const completedCount = completedIds.size;
  const totalCount = activeHabits.length;
  const completionRate =
    totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0;

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between gap-4 mb-2">
          <div>
            <h1 className="text-2xl sm:text-3xl font-semibold">Quick Mark</h1>
            <p className="text-sm text-muted-foreground mt-1">
              {format(today, "EEEE, MMMM d, yyyy")}
            </p>
          </div>
          <div className="text-right">
            <div className="text-2xl font-semibold">
              {completedCount}/{totalCount}
            </div>
            <div className="text-xs text-muted-foreground">Completed</div>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="mt-6 rounded-lg border border-zinc-200 shadow-xs bg-background p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Today's Progress</h2>
            <span className="text-2xl font-bold">{completionRate}%</span>
          </div>
          <div className="w-full bg-muted rounded-full h-3">
            <div
              className="bg-primary h-3 rounded-full transition-all duration-500"
              style={{ width: `${completionRate}%` }}
            />
          </div>
          {completionRate === 100 && (
            <div className="mt-4 flex items-center gap-2 text-green-600">
              <Sparkles className="h-4 w-4" />
              <span className="text-sm font-medium">
                Amazing! You've completed all habits today! 🎉
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Habits Grid */}
      {activeHabits.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-muted-foreground">No active habits to mark</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {activeHabits.map((habit) => {
            const isCompleted = completedIds.has(habit.id);
            const icon = habit.icon || getIconEmoji(habit.iconId || "fire");
            const color = habit.color || "#6366f1";

            return (
              <div
                key={habit.id}
                className={cn(
                  "rounded-lg border flex flex-col border-zinc-200 shadow-xs hover:shadow-md bg-background p-4 transition-all",
                  isCompleted ? "" : ""
                )}
              >
                <div className="flex items-start gap-3 mb-4 flex-1">
                  <div className="h-12 w-12 bg-muted rounded-xl flex items-center justify-center text-2xl shrink-0">
                    {icon}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold truncate">{habit.name}</h3>
                    {habit.description && (
                      <p className="text-xs text-muted-foreground line-clamp-2 mt-1">
                        {habit.description}
                      </p>
                    )}
                    <div className="flex items-center gap-1.5 mt-3">
                      {habit.currentStreak !== undefined &&
                      habit.currentStreak > 0 ? (
                        <>
                          <span className="text-xs">🔥</span>
                          <span className="text-xs text-muted-foreground">
                            {habit.currentStreak} day streak
                          </span>
                        </>
                      ) : (
                        <span></span>
                      )}
                    </div>
                  </div>
                </div>

                {isCompleted ? (
                  <Button
                    onClick={() => handleUncomplete(habit.id)}
                    className="w-full border border-zinc-200 cursor-pointer"
                    size="sm"
                    variant="outline"
                    disabled={
                      unmarkComplete.isPending || markComplete.isPending
                    }
                  >
                    <X className="h-4 w-4 mr-2" />
                    {unmarkComplete.isPending || markComplete.isPending
                      ? "Updating..."
                      : "Unmark"}
                  </Button>
                ) : (
                  <Button
                    onClick={() => handleComplete(habit.id)}
                    className="w-full cursor-pointer"
                    size="sm"
                    variant="default"
                    disabled={
                      markComplete.isPending || unmarkComplete.isPending
                    }
                  >
                    {markComplete.isPending || unmarkComplete.isPending ? (
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
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
