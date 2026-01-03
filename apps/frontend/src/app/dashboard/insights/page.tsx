"use client";

import { useState, useMemo } from "react";
import { format } from "date-fns";
import {
  ChartColumn,
  TrendingUp,
  Target,
  Calendar,
  Award,
  Loader2,
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { getCategoryIcon } from "@/lib/habit-utils";
import {
  useCompletionTrend,
  useCategoryBreakdown,
  useDayOfWeekAnalysis,
  useMetrics,
  useTopHabits,
  useStreakLeaderboard,
  useInsights,
  useHabits,
} from "@/lib/hooks";
import { mapHabitResponsesToHabits } from "@/lib/api/mappers";

type TimePeriod = "7d" | "30d" | "90d" | "all";

const COLORS = [
  "#f87171", // red-400
  "#fb923c", // orange-400
  "#fbbf24", // amber-400
  "#a3e635", // lime-400
  "#4ade80", // green-400 (optional)
  "#2dd4bf", // teal-400
  "#22d3ee", // cyan-400
  "#60a5fa", // blue-400
  "#818cf8", // indigo-400
  "#a78bfa", // violet-400
];

export default function InsightsPage() {
  const [timePeriod, setTimePeriod] = useState<TimePeriod>("30d");

  const { data: completionTrend, isLoading: trendLoading } =
    useCompletionTrend(timePeriod);
  const { data: categoryBreakdown, isLoading: categoryLoading } =
    useCategoryBreakdown();
  const { data: dayOfWeek, isLoading: dayLoading } = useDayOfWeekAnalysis();
  const { data: metrics, isLoading: metricsLoading } = useMetrics();
  const { data: topHabits, isLoading: topHabitsLoading } = useTopHabits(
    10,
    "completion"
  );
  const { data: streakLeaderboard, isLoading: streakLoading } =
    useStreakLeaderboard(10);
  const { data: insights, isLoading: insightsLoading } = useInsights();

  const { data: habitsData } = useHabits();

  const isLoading =
    trendLoading ||
    categoryLoading ||
    dayLoading ||
    metricsLoading ||
    topHabitsLoading ||
    streakLoading ||
    insightsLoading;

  const completionTrendData = useMemo(() => {
    if (!completionTrend?.data) return [];
    return completionTrend.data.map((item) => ({
      date: format(new Date(item.date), "MMM d"),
      completionRate: Math.round(item.completionRate),
      completions: item.completions,
    }));
  }, [completionTrend]);

  const categoryData = useMemo(() => {
    if (!categoryBreakdown?.data) return [];
    return categoryBreakdown.data.map((item) => ({
      name: item.label,
      value: Math.round(item.avgCompletionRate),
      icon: getCategoryIcon(item.category as any),
    }));
  }, [categoryBreakdown]);

  const dayOfWeekData = useMemo(() => {
    if (!dayOfWeek?.data) return [];
    const dayNames = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
    return dayOfWeek.data
      .sort((a, b) => a.dayIndex - b.dayIndex)
      .map((item) => ({
        day: dayNames[item.dayIndex] || item.day,
        completionRate: Math.round(item.completionRate),
        completions: item.completions,
      }));
  }, [dayOfWeek]);

  const formattedMetrics = useMemo(() => {
    if (!metrics) {
      return {
        avgCompletion: 0,
        avgStreak: 0,
        totalCompletions: 0,
        consistencyScore: 0,
      };
    }
    return {
      avgCompletion: Math.round(metrics.avgCompletionRate || 0),
      avgStreak: Math.round(metrics.avgStreak || 0),
      totalCompletions: metrics.totalCompletions || 0,
      consistencyScore: Math.round(metrics.consistencyScore || 0),
    };
  }, [metrics]);

  const formattedTopHabits = useMemo(() => {
    if (!topHabits?.data || !habitsData?.data) return [];
    const habits = mapHabitResponsesToHabits(habitsData.data);
    const habitMap = new Map(habits.map((h) => [h.id, h]));

    return topHabits.data.slice(0, 5).map((item: any) => {
      const fullHabit = habitMap.get(item.habitId || item.id);
      return {
        id: item.habitId || item.id,
        name: item.name,
        completionRate: Math.round(item.completionRate || 0),
        color: fullHabit?.color || "#6366f1",
        icon: fullHabit?.icon || "🔥",
      };
    });
  }, [topHabits, habitsData]);

  const formattedStreakLeaderboard = useMemo(() => {
    if (!streakLeaderboard?.data || !habitsData?.data) return [];
    const habits = mapHabitResponsesToHabits(habitsData.data);
    const habitMap = new Map(habits.map((h) => [h.id, h]));

    return streakLeaderboard.data.slice(0, 5).map((item: any) => {
      const fullHabit = habitMap.get(item.habitId || item.id);
      return {
        id: item.habitId || item.id,
        name: item.name,
        currentStreak: item.currentStreak || 0,
        longestStreak: item.longestStreak || 0,
        color: fullHabit?.color || "#6366f1",
        icon: fullHabit?.icon || "🔥",
      };
    });
  }, [streakLeaderboard, habitsData]);

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between gap-4 mb-6">
          <div>
            <h1 className="text-2xl sm:text-3xl font-semibold">Insights</h1>
            <p className="text-sm text-muted-foreground mt-1">
              Analyze your habit patterns and performance
            </p>
          </div>
          <Select
            value={timePeriod}
            onValueChange={(v) => setTimePeriod(v as TimePeriod)}
          >
            <SelectTrigger className="w-35">
              <SelectValue />
            </SelectTrigger>
            <SelectContent className="border border-zinc-200">
              <SelectItem value="7d">Last 7 days</SelectItem>
              <SelectItem value="30d">Last 30 days</SelectItem>
              <SelectItem value="90d">Last 90 days</SelectItem>
              <SelectItem value="all">All time</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <Target className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">
                Avg. Completion
              </span>
            </div>
            <p className="text-2xl font-semibold">
              {formattedMetrics.avgCompletion}%
            </p>
          </div>
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Avg. Streak</span>
            </div>
            <p className="text-2xl font-semibold">
              {formattedMetrics.avgStreak} days
            </p>
          </div>
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">This Week</span>
            </div>
            <p className="text-2xl font-semibold">
              {formattedMetrics.totalCompletions}
            </p>
          </div>
          <div className="rounded-lg border border-zinc-200 shadow-xs bg-background p-4">
            <div className="flex items-center gap-2 mb-1">
              <Award className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Consistency</span>
            </div>
            <p className="text-2xl font-semibold">
              {formattedMetrics.consistencyScore}%
            </p>
          </div>
        </div>
      </div>

      {/* Completion Trends */}
      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Completion Trends</h2>
        <div className="rounded-lg border border-zinc-200 shadow-sm bg-background p-6">
          {completionTrendData.length === 0 ? (
            <div className="flex items-center justify-center h-75 text-muted-foreground">
              No trend data available
            </div>
          ) : (
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={completionTrendData}>
                <defs>
                  <linearGradient
                    id="colorCompletion"
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                  >
                    <stop offset="5%" stopColor="#00c951" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#00c951" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid
                  strokeDasharray="3 3"
                  className="stroke-muted-foreground/30"
                />
                <XAxis
                  dataKey="date"
                  className="text-[10px]"
                  fontSize={10}
                  color="#888888"
                  tick={{ fill: "currentColor" }}
                />
                <YAxis
                  className="text-xs"
                  fontSize={14}
                  color="#666666"
                  tick={{ fill: "currentColor" }}
                  domain={[0, 100]}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "hsl(var(--background))",
                    border: "1px solid hsl(var(--border))",
                    borderRadius: "0.5rem",
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="completionRate"
                  stroke="#00c951"
                  fillOpacity={1}
                  fill="url(#colorCompletion)"
                />
              </AreaChart>
            </ResponsiveContainer>
          )}
        </div>
      </section>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Category Analysis */}
        <section>
          <h2 className="text-xl font-semibold mb-4">Category Breakdown</h2>
          <div className="rounded-lg border border-zinc-200 shadow-sm bg-background p-6">
            {categoryData.length === 0 ? (
              <div className="flex items-center justify-center h-75 text-muted-foreground">
                No category data available
              </div>
            ) : (
              <>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={categoryData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) =>
                        `${name} ${((percent || 0) * 100).toFixed(0)}%`
                      }
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {categoryData.map((entry, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={COLORS[index % COLORS.length]}
                        />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
                <div className="mt-4 space-y-2">
                  {categoryData.map((cat, idx) => (
                    <div
                      key={cat.name}
                      className="flex items-center justify-between"
                    >
                      <div className="flex items-center gap-2">
                        <div
                          className="w-3 h-3 rounded-full"
                          style={{
                            backgroundColor: COLORS[idx % COLORS.length],
                          }}
                        />
                        <span className="text-sm">{cat.name}</span>
                      </div>
                      <span className="text-sm font-medium text-muted-foreground">
                        {cat.value}%
                      </span>
                    </div>
                  ))}
                </div>
              </>
            )}
          </div>
        </section>

        {/* Day of Week Analysis */}
        <section>
          <h2 className="text-xl font-semibold mb-4">Best Day of Week</h2>
          <div className="rounded-lg border border-zinc-200 shadow-sm bg-background p-6">
            {dayOfWeekData.length === 0 ? (
              <div className="flex items-center justify-center h-75 text-muted-foreground">
                No day of week data available
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={dayOfWeekData}>
                  <CartesianGrid
                    strokeDasharray="3 3"
                    className="stroke-muted"
                  />
                  <XAxis
                    dataKey="day"
                    className="text-xs"
                    fontSize={12}
                    color="#666666"
                    tick={{ fill: "currentColor" }}
                  />
                  <YAxis
                    className="text-xs"
                    fontSize={12}
                    color="#666666"
                    tick={{ fill: "currentColor" }}
                    domain={[0, 100]}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "hsl(var(--background))",
                      border: "1px solid hsl(var(--border))",
                      borderRadius: "0.5rem",
                    }}
                  />
                  <Bar
                    dataKey="completionRate"
                    fill="#00c951"
                    radius={[8, 8, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>
        </section>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Habit Performance */}
        <section>
          <h2 className="text-xl font-semibold mb-4">Habit Performance</h2>
          <div className="rounded-lg border border-zinc-200 shadow-sm bg-background p-6 space-y-4">
            {formattedTopHabits.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                No habit data available
              </div>
            ) : (
              formattedTopHabits.map((habit: any, idx) => (
                <div key={habit.id} className="flex items-center gap-4">
                  <div className="flex items-center gap-3 flex-1">
                    <div className="h-10 w-10 rounded-lg flex items-center justify-center text-xl shrink-0 bg-muted">
                      {habit.icon}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between mb-1">
                        <span className="font-medium text-sm text-foreground truncate">
                          {habit.name}
                        </span>
                        <span className="text-sm font-medium">
                          {habit.completionRate}%
                        </span>
                      </div>
                      <div className="w-full bg-muted rounded-full h-2">
                        <div
                          className="h-2 rounded-full transition-all"
                          style={{
                            width: `${habit.completionRate}%`,
                            backgroundColor: habit.color,
                          }}
                        />
                      </div>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </section>

        {/* Streak Leaderboard */}
        <section>
          <h2 className="text-xl font-semibold mb-4">Streak Leaderboard</h2>
          <div className="rounded-lg border border-zinc-200 shadow-sm bg-background p-6 space-y-4">
            {formattedStreakLeaderboard.length > 0 ? (
              formattedStreakLeaderboard.map((habit: any, idx) => (
                <div
                  key={habit.id}
                  className="flex items-center gap-4 p-3 rounded-lg bg-muted/50"
                >
                  <div className="flex items-center justify-center w-8 h-8 rounded-full bg-primary/10 text-primary font-semibold text-sm shrink-0">
                    {idx + 1}
                  </div>
                  <div className="h-10 w-10 rounded-lg flex items-center justify-center text-xl shrink-0 bg-muted">
                    {habit.icon}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium truncate">{habit.name}</div>
                    <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                      <span>
                        {habit.currentStreak} day
                        {habit.currentStreak !== 1 ? "s" : ""}
                      </span>
                    </div>
                  </div>
                  {habit.longestStreak && (
                    <div className="text-sm text-muted-foreground">
                      best : {habit.longestStreak}
                    </div>
                  )}
                </div>
              ))
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                No active streaks
              </div>
            )}
          </div>
        </section>
      </div>

      {/* Insights & Recommendations */}
      {insights?.data && insights.data.length > 0 && (
        <section className="mt-8">
          <h2 className="text-xl font-semibold mb-4">
            Insights & Recommendations
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {insights.data.map((insight: any, idx: number) => {
              const getIcon = () => {
                switch (insight.type) {
                  case "positive":
                    return (
                      <TrendingUp className="h-4 w-4 text-green-600 dark:text-green-400" />
                    );
                  case "suggestion":
                    return (
                      <Target className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                    );
                  case "achievement":
                    return (
                      <Award className="h-4 w-4 text-orange-600 dark:text-orange-400" />
                    );
                  default:
                    return (
                      <ChartColumn className="h-4 w-4 text-purple-600 dark:text-purple-400" />
                    );
                }
              };

              const getBgColor = () => {
                switch (insight.type) {
                  case "positive":
                    return "bg-green-500/10";
                  case "suggestion":
                    return "bg-blue-500/10";
                  case "achievement":
                    return "bg-orange-500/10";
                  default:
                    return "bg-purple-500/10";
                }
              };

              return (
                <div
                  key={idx}
                  className="rounded-lg border border-zinc-200 shadow-sm bg-background p-4"
                >
                  <div className="flex items-start gap-3">
                    <div className={`rounded-full ${getBgColor()} p-2`}>
                      {getIcon()}
                    </div>
                    <div>
                      <h3 className="font-medium mb-1">{insight.title}</h3>
                      <p className="text-sm text-muted-foreground">
                        {insight.description}
                      </p>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </section>
      )}
    </div>
  );
}
