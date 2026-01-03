"use client";

import { useState, useMemo } from "react";
import {
  format,
  startOfMonth,
  endOfMonth,
  eachDayOfInterval,
  isSameMonth,
  isSameDay,
  addMonths,
  subMonths,
  addWeeks,
  subWeeks,
  addYears,
  subYears,
  startOfWeek,
  endOfWeek,
  startOfYear,
  endOfYear,
  getISOWeek,
  isSameYear,
} from "date-fns";
import {
  ChevronLeft,
  ChevronRight,
  Calendar as CalendarIcon,
  Filter,
  CheckCircle2,
  X,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useHabits,
  useCalendarMonth,
  useCalendarWeek,
  useCalendarYear,
  useCalendarCompletions,
} from "@/lib/hooks";
import { mapHabitResponsesToHabits } from "@/lib/api/mappers";

type ViewMode = "month" | "week" | "year";

export default function CalendarPage() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [viewMode, setViewMode] = useState<ViewMode>("month");
  const [selectedHabitIds, setSelectedHabitIds] = useState<string[]>([]);

  const { data: habitsData } = useHabits();
  const habits = habitsData?.data
    ? mapHabitResponsesToHabits(habitsData.data)
    : [];

  const effectiveHabitIds =
    selectedHabitIds.length > 0 ? selectedHabitIds : undefined;

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth() + 1;
  const week = getISOWeek(currentDate);

  const { data: monthData, isLoading: monthLoading } = useCalendarMonth(
    year,
    month,
    effectiveHabitIds
  );
  const { data: weekData, isLoading: weekLoading } = useCalendarWeek(
    year,
    week,
    effectiveHabitIds
  );
  const { data: yearData, isLoading: yearLoading } = useCalendarYear(
    year,
    effectiveHabitIds
  );

  const calendarDays = useMemo(() => {
    if (viewMode === "week") {
      const weekStart = startOfWeek(currentDate, { weekStartsOn: 0 });
      const weekEnd = endOfWeek(currentDate, { weekStartsOn: 0 });
      return eachDayOfInterval({ start: weekStart, end: weekEnd });
    } else if (viewMode === "year") {
      const yearStart = startOfYear(currentDate);
      const yearEnd = endOfYear(currentDate);
      return eachDayOfInterval({ start: yearStart, end: yearEnd });
    } else {
      // Month view
      const monthStart = startOfMonth(currentDate);
      const monthEnd = endOfMonth(currentDate);
      const calendarStart = startOfWeek(monthStart, { weekStartsOn: 0 });
      const calendarEnd = endOfWeek(monthEnd, { weekStartsOn: 0 });
      return eachDayOfInterval({ start: calendarStart, end: calendarEnd });
    }
  }, [currentDate, viewMode]);

  const toggleHabit = (habitId: string) => {
    setSelectedHabitIds((prev) => {
      if (prev.includes(habitId)) {
        return prev.filter((id) => id !== habitId);
      }
      return [...prev, habitId];
    });
  };

  const getCompletionsForDate = (date: Date) => {
    const dateStr = format(date, "yyyy-MM-dd");
    let completions: Array<{
      id: string;
      name: string;
      color?: string;
      icon?: string;
    }> = [];

    if (viewMode === "month" && monthData?.days) {
      const dayData = monthData.days.find((d: { date: string }) => d.date === dateStr);
      if (dayData) {
        completions = dayData.completions.map((id: string) => {
          const habit = habits.find((h) => h.id === id);
          return {
            id,
            name: habit?.name || "",
            color: habit?.color,
            icon: habit?.icon,
          };
        });
      }
    } else if (viewMode === "week" && weekData?.days) {
      const dayData = weekData.days.find((d: { date: string }) => d.date === dateStr);
      if (dayData) {
        completions = dayData.completions.map((id: string) => {
          const habit = habits.find((h) => h.id === id);
          return {
            id,
            name: habit?.name || "",
            color: habit?.color,
            icon: habit?.icon,
          };
        });
      }
    }
    // Year view completions are handled separately in YearHeatmapView component

    return completions;
  };

  const getCompletionRateForDate = (date: Date) => {
    const dateStr = format(date, "yyyy-MM-dd");

    if (viewMode === "month" && monthData?.days) {
      const dayData = monthData.days.find((d: { date: string }) => d.date === dateStr);
      return dayData?.completionRate || 0;
    } else if (viewMode === "week" && weekData?.days) {
      const dayData = weekData.days.find((d: { date: string }) => d.date === dateStr);
      return dayData?.completionRate || 0;
    } else if (viewMode === "year" && yearData?.heatmap) {
      const dayData = yearData.heatmap.find((d: { date: string }) => d.date === dateStr);
      return dayData?.completionRate || 0;
    }

    return 0;
  };

  const periodStats = useMemo(() => {
    if (viewMode === "month" && monthData?.statistics) {
      return {
        totalCompletions: monthData.statistics.totalCompletions,
        daysWithCompletions: monthData.statistics.daysWithCompletions,
        completionRate: Math.round(monthData.statistics.completionRate),
        totalDays: monthData.statistics.totalDays,
      };
    } else if (viewMode === "week" && weekData?.statistics) {
      return {
        totalCompletions: weekData.statistics.totalCompletions,
        daysWithCompletions: weekData.statistics.daysWithCompletions,
        completionRate: Math.round(weekData.statistics.completionRate),
        totalDays: weekData.statistics.totalDays,
      };
    } else if (viewMode === "year" && yearData?.statistics) {
      return {
        totalCompletions: yearData.statistics.totalCompletions,
        daysWithCompletions: yearData.statistics.daysWithCompletions,
        completionRate: Math.round(yearData.statistics.completionRate),
        totalDays: yearData.statistics.totalDays,
      };
    }

    return {
      totalCompletions: 0,
      daysWithCompletions: 0,
      completionRate: 0,
      totalDays: 0,
    };
  }, [viewMode, monthData, weekData, yearData]);

  const navigateDate = (direction: "prev" | "next") => {
    setCurrentDate((prev) => {
      if (viewMode === "week") {
        return direction === "prev" ? subWeeks(prev, 1) : addWeeks(prev, 1);
      } else if (viewMode === "year") {
        return direction === "prev" ? subYears(prev, 1) : addYears(prev, 1);
      } else {
        return direction === "prev" ? subMonths(prev, 1) : addMonths(prev, 1);
      }
    });
  };

  const getDateLabel = () => {
    if (viewMode === "week") {
      const weekStart = startOfWeek(currentDate, { weekStartsOn: 0 });
      const weekEnd = endOfWeek(currentDate, { weekStartsOn: 0 });
      return `${format(weekStart, "MMM d")} - ${format(
        weekEnd,
        "MMM d, yyyy"
      )}`;
    } else if (viewMode === "year") {
      return format(currentDate, "yyyy");
    } else {
      return format(currentDate, "MMMM yyyy");
    }
  };

  const goToToday = () => {
    setCurrentDate(new Date());
  };

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between gap-4 mb-6">
          <div>
            <h1 className="text-2xl sm:text-3xl font-semibold">Calendar</h1>
            <p className="text-sm text-muted-foreground mt-1">
              Track your habit completion over time
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Select
              value={viewMode}
              onValueChange={(v) => setViewMode(v as ViewMode)}
            >
              <SelectTrigger className="w-30">
                <SelectValue />
              </SelectTrigger>
              <SelectContent className="border border-zinc-200">
                <SelectItem value="month">Month</SelectItem>
                <SelectItem value="week">Week</SelectItem>
                <SelectItem value="year">Year</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Navigation */}
        <div className="flex items-center justify-between mb-6 bg-background rounded-lg border border-zinc-200 shadow-xs p-4">
          <div className="flex items-center gap-4">
            <Button
              variant="outline"
              size="icon"
              onClick={() => navigateDate("prev")}
              className="h-9 w-9 border-zinc-200 shadow-sm cursor-pointer"
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <h2 className="text-xl font-medium min-w-50 text-center">
              {getDateLabel()}
            </h2>
            <Button
              variant="outline"
              size="icon"
              onClick={() => navigateDate("next")}
              className="h-9 w-9 border-zinc-200 shadow-sm cursor-pointer"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
          <Button
            variant="default"
            onClick={goToToday}
            className="font-medium cursor-pointer"
          >
            Today
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Calendar Grid */}
        <div className="lg:col-span-3">
          <div className="rounded-xl border-2 border-zinc-200 bg-background p-6 shadow-sm">
            {viewMode === "year" ? (
              <YearHeatmapView
                calendarDays={calendarDays}
                currentDate={currentDate}
                getCompletionsForDate={getCompletionsForDate}
                getCompletionRateForDate={getCompletionRateForDate}
                effectiveHabitIds={effectiveHabitIds}
              />
            ) : viewMode === "week" ? (
              <WeekView
                calendarDays={calendarDays}
                getCompletionsForDate={getCompletionsForDate}
                getCompletionRateForDate={getCompletionRateForDate}
              />
            ) : (
              <>
                {/* Day Headers */}
                <div className="grid grid-cols-7 gap-2 mb-3">
                  {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map(
                    (day) => (
                      <div
                        key={day}
                        className="text-center text-sm font-semibold text-muted-foreground py-2"
                      >
                        {day}
                      </div>
                    )
                  )}
                </div>

                {/* Calendar Days */}
                <div className="grid grid-cols-7 gap-2">
                  {calendarDays.map((day, idx) => {
                    const isCurrentMonth = isSameMonth(day, currentDate);
                    const isToday = isSameDay(day, new Date());
                    const completions = getCompletionsForDate(day);
                    const completionRate = getCompletionRateForDate(day);
                    const totalHabits =
                      selectedHabitIds.length > 0
                        ? selectedHabitIds.length
                        : habits.length;
                    const hasManyHabits = totalHabits > 10;

                    return (
                      <div
                        key={idx}
                        className={cn(
                          "aspect-square rounded-xl border-2 p-3 transition-all cursor-pointer hover:shadow-md hover:scale-105",
                          !isCurrentMonth && "opacity-30 border-muted",
                          isCurrentMonth && "border-border",
                          isToday &&
                            "ring-2 ring-primary ring-offset-1 bg-primary/5",
                          completions.length > 0 &&
                            !isToday &&
                            "bg-green-50 dark:bg-green-950/20 border-green-200 dark:border-green-800"
                        )}
                        onClick={() => {}}
                      >
                        <div className="flex flex-col h-full">
                          <span
                            className={cn(
                              "text-sm font-medium mb-1",
                              isToday && "text-primary font-semibold"
                            )}
                          >
                            {format(day, "d")}
                          </span>

                          {hasManyHabits ? (
                            // Compact view for many habits
                            <>
                              <div className="flex-1 flex items-center justify-center">
                                <div className="text-center">
                                  <div className="text-lg font-semibold">
                                    {Math.round(completionRate)}%
                                  </div>
                                  <div className="text-xs text-muted-foreground">
                                    {completions.length}/{totalHabits}
                                  </div>
                                </div>
                              </div>
                              {completionRate > 0 && (
                                <div className="mt-1">
                                  <div className="w-full bg-muted rounded-full h-1.5">
                                    <div
                                      className="bg-primary h-1.5 rounded-full transition-all"
                                      style={{ width: `${completionRate}%` }}
                                    />
                                  </div>
                                </div>
                              )}
                            </>
                          ) : (
                            // Detailed view for few habits
                            <>
                              <div className="flex-1 flex flex-wrap gap-0.5 items-start">
                                {completions.slice(0, 4).map((habit) => (
                                  <div
                                    key={habit.id}
                                    className="w-2 h-2 rounded-full"
                                    style={{ backgroundColor: habit.color }}
                                    title={habit.name}
                                  />
                                ))}
                                {completions.length > 4 && (
                                  <span className="text-xs text-muted-foreground">
                                    +{completions.length - 4}
                                  </span>
                                )}
                              </div>
                              {completionRate > 0 && (
                                <div className="mt-1">
                                  <div className="w-full bg-muted rounded-full h-1">
                                    <div
                                      className="bg-primary h-1 rounded-full"
                                      style={{ width: `${completionRate}%` }}
                                    />
                                  </div>
                                </div>
                              )}
                            </>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Statistics */}
          <div className="rounded-lg border border-zinc-200 bg-background p-4 shadow-xs">
            <h3 className="font-semibold mb-4 flex items-center gap-2">
              <CalendarIcon className="h-4 w-4" />
              {viewMode === "week"
                ? "This Week"
                : viewMode === "year"
                ? "This Year"
                : "This Month"}
            </h3>
            <div className="space-y-3">
              <div>
                <div className="text-sm text-muted-foreground mb-1">
                  Completion Rate
                </div>
                <div className="text-2xl font-semibold">
                  {periodStats.completionRate}%
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground mb-1">
                  Total Completions
                </div>
                <div className="text-2xl font-semibold">
                  {periodStats.totalCompletions}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground mb-1">
                  Active Days
                </div>
                <div className="text-2xl font-semibold">
                  {periodStats.daysWithCompletions}
                </div>
              </div>
            </div>
          </div>

          {/* Habit Filter */}
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <h3 className="font-semibold mb-4 flex items-center gap-2">
              <Filter className="h-4 w-4" />
              Habits
            </h3>
            <div className="space-y-2">
              {habits.map((habit) => (
                <label
                  key={habit.id}
                  className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={
                      selectedHabitIds.length === 0 ||
                      selectedHabitIds.includes(habit.id)
                    }
                    onChange={() => toggleHabit(habit.id)}
                    className="rounded border-gray-300 accent-foreground"
                  />
                  <div
                    className="w-4 h-4 rounded-full shrink-0"
                    style={{ backgroundColor: habit.color }}
                  />
                  <span className="text-sm font-medium flex-1">
                    {habit.name}
                  </span>
                </label>
              ))}
            </div>
            <div className="mt-4 flex gap-2">
              <Button
                variant="outline"
                size="sm"
                className="flex-1 border-zinc-200 cursor-pointer"
                onClick={() => setSelectedHabitIds([])}
              >
                Show All
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="flex-1 border-zinc-200 cursor-pointer"
                onClick={() => setSelectedHabitIds(habits.map((h) => h.id))}
              >
                Clear
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Week View Component
function WeekView({
  calendarDays,
  getCompletionsForDate,
  getCompletionRateForDate,
}: {
  calendarDays: Date[];
  getCompletionsForDate: (
    date: Date
  ) => Array<{ id: string; name: string; color?: string; icon?: string }>;
  getCompletionRateForDate: (date: Date) => number;
}) {
  const [expandedDays, setExpandedDays] = useState<Set<string>>(new Set());

  const toggleDay = (dayStr: string) => {
    const newExpanded = new Set(expandedDays);
    if (newExpanded.has(dayStr)) {
      newExpanded.delete(dayStr);
    } else {
      newExpanded.add(dayStr);
    }
    setExpandedDays(newExpanded);
  };

  return (
    <div className="space-y-4">
      {/* Day Headers */}
      <div className="grid grid-cols-7 gap-2">
        {calendarDays.map((day) => {
          const dayStr = format(day, "yyyy-MM-dd");
          const isToday = isSameDay(day, new Date());
          const completions = getCompletionsForDate(day);
          const completionRate = getCompletionRateForDate(day);
          const isExpanded = expandedDays.has(dayStr);
          const showCompact = completions.length > 10 && !isExpanded;

          return (
            <div
              key={dayStr}
              className={cn(
                "rounded-lg border border-zinc-200 shadow-xs p-3 flex flex-col h-125",
                isToday && "ring-2 ring-primary bg-primary/5"
              )}
            >
              {/* Date Header */}
              <div className="text-center mb-2">
                <div
                  className={cn(
                    "text-sm font-medium",
                    isToday && "text-primary font-semibold"
                  )}
                >
                  {format(day, "EEE")}
                </div>
                <div
                  className={cn(
                    "text-lg font-semibold mt-1",
                    isToday && "text-primary"
                  )}
                >
                  {format(day, "d")}
                </div>
              </div>

              {/* Completion Summary */}
              <div className="text-center mb-3">
                <div className="text-2xl font-bold">
                  {Math.round(completionRate)}%
                </div>
                <div className="text-xs text-muted-foreground">
                  {completions.length}
                </div>
                {completionRate > 0 && (
                  <div className="mt-2">
                    <div className="w-full bg-muted rounded-full h-2">
                      <div
                        className="bg-primary h-2 rounded-full transition-all"
                        style={{ width: `${completionRate}%` }}
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Habits List - Fixed height scrollable area */}
              <div className="flex-1 flex flex-col min-h-0">
                {showCompact ? (
                  <>
                    <div className="space-y-1 overflow-hidden shrink-0">
                      {completions.slice(0, 5).map((habit) => {
                        return (
                          <div
                            key={habit.id}
                            className={cn(
                              "flex items-center gap-1.5 p-1 rounded text-xs",
                              "bg-green-500/10"
                            )}
                          >
                            <div
                              className="w-1.5 h-1.5 rounded-full shrink-0"
                              style={{
                                backgroundColor: habit.color || "#6366f1",
                              }}
                            />
                            <span className="truncate flex-1 text-[10px]">
                              {habit.name}
                            </span>
                            <span className="text-green-600 dark:text-green-400 text-[10px] shrink-0">
                              ✓
                            </span>
                          </div>
                        );
                      })}
                    </div>
                    {completions.length > 5 && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="w-full mt-2 text-xs h-7 shrink-0"
                        onClick={() => toggleDay(dayStr)}
                      >
                        +{completions.length - 5} more
                      </Button>
                    )}
                  </>
                ) : (
                  // Expanded/Full view - scrollable
                  <div className="space-y-1 flex-1 overflow-y-auto pr-1 min-h-0">
                    {completions.map((habit) => {
                      return (
                        <div
                          key={habit.id}
                          className={cn(
                            "flex items-center gap-2 p-1.5 rounded text-xs",
                            "bg-green-500/10 border border-green-500/20"
                          )}
                        >
                          <div
                            className="w-2 h-2 rounded-full shrink-0"
                            style={{
                              backgroundColor: habit.color || "#6366f1",
                            }}
                          />
                          <span className="truncate flex-1">{habit.name}</span>
                          <span className="text-green-600 dark:text-green-400 shrink-0">
                            ✓
                          </span>
                        </div>
                      );
                    })}
                    {completions.length > 10 && isExpanded && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="w-full mt-2 text-xs h-7"
                        onClick={() => toggleDay(dayStr)}
                      >
                        Show less
                      </Button>
                    )}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// Year Heatmap View Component
function YearHeatmapView({
  calendarDays,
  currentDate,
  getCompletionsForDate,
  getCompletionRateForDate,
  effectiveHabitIds,
}: {
  calendarDays: Date[];
  currentDate: Date;
  getCompletionsForDate: (
    date: Date
  ) => Array<{ id: string; name: string; color?: string; icon?: string }>;
  getCompletionRateForDate: (date: Date) => number;
  effectiveHabitIds?: string[];
}) {
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  
  // Fetch completions for the selected date
  const selectedDateStr = selectedDate ? format(selectedDate, "yyyy-MM-dd") : "";
  const { data: selectedDateCompletions } = useCalendarCompletions(
    selectedDateStr,
    selectedDateStr,
    effectiveHabitIds
  );

  const firstDay = calendarDays[0];
  const firstDayOfWeek = firstDay.getDay(); // 0 = Sunday

  const weeks: Date[][] = [];
  let currentWeek: Date[] = [];

  if (firstDayOfWeek !== 0) {
    for (let i = 0; i < firstDayOfWeek; i++) {
      currentWeek.push(new Date(0)); // Placeholder for empty days
    }
  }

  calendarDays.forEach((day) => {
    const dayOfWeek = day.getDay(); // 0 = Sunday, 6 = Saturday

    if (dayOfWeek === 0 && currentWeek.length > 0) {
      while (currentWeek.length < 7) {
        currentWeek.push(new Date(0));
      }
      weeks.push(currentWeek);
      currentWeek = [];
    }

    currentWeek.push(day);
  });

  if (currentWeek.length > 0) {
    while (currentWeek.length < 7) {
      currentWeek.push(new Date(0));
    }
    weeks.push(currentWeek);
  }

  const monthLabels: { weekIdx: number; month: string }[] = [];
  let lastMonth = -1;

  weeks.forEach((week, weekIdx) => {
    const firstRealDay = week.find((d) => d.getTime() !== 0);
    if (firstRealDay) {
      const month = firstRealDay.getMonth();
      if (month !== lastMonth) {
        monthLabels.push({ weekIdx, month: format(firstRealDay, "MMM") });
        lastMonth = month;
      }
    }
  });

  const getIntensity = (completionRate: number) => {
    if (completionRate === 0) return 0;
    if (completionRate < 25) return 1;
    if (completionRate < 50) return 2;
    if (completionRate < 75) return 3;
    return 4;
  };

  const getColor = (intensity: number) => {
    const colors = [
      "bg-muted border-muted-foreground/20", // 0 - no completions
      "bg-green-100 dark:bg-green-900/30 border-green-300 dark:border-green-700", // 1 - low
      "bg-green-300 dark:bg-green-800/50 border-green-400 dark:border-green-600", // 2 - medium
      "bg-green-500 dark:bg-green-700 border-green-600 dark:border-green-500", // 3 - high
      "bg-green-600 dark:bg-green-600 border-green-700 dark:border-green-400", // 4 - very high
    ];
    return colors[intensity] || colors[0];
  };

  // Get completions for selected date - use API data if available, otherwise fallback
  const selectedCompletions = useMemo(() => {
    if (!selectedDate) return [];
    
    // If we have API data for the selected date, use it
    if (selectedDateCompletions?.completions) {
      const dateStr = format(selectedDate, "yyyy-MM-dd");
      const dayData = selectedDateCompletions.completions.find(
        (c: { date: string }) => c.date === dateStr
      );
      if (dayData && dayData.habits) {
        return dayData.habits;
      }
    }
    
    // Fallback to empty array for year view (since heatmap doesn't include habit details)
    return [];
  }, [selectedDate, selectedDateCompletions]);

  return (
    <div className="space-y-6">
      <div className="space-y-4 overflow-x-auto pb-4">
        {/* Month Labels */}
        <div className="flex gap-1 mb-3 min-w-fit px-1">
          <div className="w-10.25 shrink-0"></div>
          {weeks.map((week, weekIdx) => {
            const monthLabel = monthLabels.find((m) => m.weekIdx === weekIdx);
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
        <div className="flex gap-1 min-w-fit px-1">
          {/* Day labels */}
          <div className="space-y-1 pr-3 shrink-0">
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

                  const completionRate = getCompletionRateForDate(day);
                  const intensity = getIntensity(completionRate);
                  const isToday = isSameDay(day, new Date());
                  const isCurrentYear = isSameYear(day, currentDate);
                  const isSelected =
                    selectedDate && isSameDay(day, selectedDate);

                  return (
                    <div
                      key={`${weekIdx}-${dayIdx}`}
                      onClick={() => setSelectedDate(day)}
                      className={cn(
                        "w-2.75 h-2.75 rounded-xs border transition-all hover:scale-110 hover:z-10 relative cursor-pointer",
                        getColor(intensity),
                        !isCurrentYear && "opacity-20",
                        isToday && "ring-2 ring-primary ring-offset-1",
                        isSelected &&
                          "ring-2 ring-blue-500 ring-offset-1 shadow-lg"
                      )}
                      title={`${format(
                        day,
                        "MMM d, yyyy"
                      )}: ${completionRate}% completion`}
                    />
                  );
                })}
              </div>
            ))}
          </div>
        </div>

        {/* Legend */}
        <div className="flex items-center justify-center gap-3 mt-6 text-xs text-muted-foreground">
          <span className="font-medium">Less</span>
          <div className="flex gap-1.5">
            {[0, 1, 2, 3, 4].map((intensity) => (
              <div
                key={intensity}
                className={cn(
                  "w-2.75 h-2.75 rounded-sm border",
                  getColor(intensity)
                )}
                title={`${
                  intensity === 0
                    ? "No"
                    : intensity === 1
                    ? "Low"
                    : intensity === 2
                    ? "Medium"
                    : intensity === 3
                    ? "High"
                    : "Very High"
                } activity`}
              />
            ))}
          </div>
          <span className="font-medium">More</span>
        </div>
      </div>

      {/* Selected Date Details */}
      {selectedDate && (
        <div className="rounded-xl border-2 border-zinc-200 bg-background p-6 mt-6 shadow-md">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-lg font-semibold">
                {format(selectedDate, "EEEE, MMMM d, yyyy")}
              </h3>
              <p className="text-sm text-muted-foreground mt-1">
                {getCompletionRateForDate(selectedDate)}% completion rate
              </p>
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setSelectedDate(null)}
              className="h-8 w-8 p-0 cursor-pointer"
            >
              <X className="h-4 w-4" />
            </Button>
          </div>

          {selectedCompletions.length > 0 ? (
            <div className="space-y-4">
              <p className="text-sm font-medium text-muted-foreground">
                {selectedCompletions.length} habit
                {selectedCompletions.length !== 1 ? "s" : ""} completed
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {selectedCompletions.map((habit) => {
                  return (
                    <div
                      key={habit.id}
                      className={cn(
                        "flex items-center gap-3 p-4 rounded-lg border-2",
                        "bg-green-500/10 border-green-500/30 hover:bg-green-500/15 transition-colors"
                      )}
                    >
                      <div
                        className="h-12 w-12 rounded-xl flex items-center justify-center text-2xl shrink-0 shadow-sm"
                        style={{
                          backgroundColor: (habit.color || "#6366f1") + "20",
                          color: habit.color || "#6366f1",
                        }}
                      >
                        {habit.icon || "🔥"}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="font-semibold text-sm">
                          {habit.name}
                        </div>
                      </div>
                      <div className="flex items-center gap-1.5 text-green-600 dark:text-green-400 shrink-0">
                        <CheckCircle2 className="h-5 w-5" />
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              <p className="mb-2">No habits completed on this day</p>
              <p className="text-sm">Try completing some habits!</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
