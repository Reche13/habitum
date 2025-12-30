import { useQuery } from "@tanstack/react-query";
import { calendarAPI } from "@/lib/api/calendar";

// Query keys
export const calendarKeys = {
  all: ["calendar"] as const,
  completions: (startDate: string, endDate: string, habitIds?: string[]) =>
    [...calendarKeys.all, "completions", startDate, endDate, habitIds] as const,
  month: (year: number, month: number, habitIds?: string[]) =>
    [...calendarKeys.all, "month", year, month, habitIds] as const,
  week: (year: number, week: number, habitIds?: string[]) =>
    [...calendarKeys.all, "week", year, week, habitIds] as const,
  year: (year: number, habitIds?: string[]) =>
    [...calendarKeys.all, "year", year, habitIds] as const,
};

// Get completions
export function useCalendarCompletions(
  startDate: string,
  endDate: string,
  habitIds?: string[]
) {
  return useQuery({
    queryKey: calendarKeys.completions(startDate, endDate, habitIds),
    queryFn: () => calendarAPI.getCompletions(startDate, endDate, habitIds),
    enabled: !!startDate && !!endDate,
  });
}

// Get month view
export function useCalendarMonth(year: number, month: number, habitIds?: string[]) {
  return useQuery({
    queryKey: calendarKeys.month(year, month, habitIds),
    queryFn: () => calendarAPI.getMonth(year, month, habitIds),
    enabled: !!year && !!month,
  });
}

// Get week view
export function useCalendarWeek(year: number, week: number, habitIds?: string[]) {
  return useQuery({
    queryKey: calendarKeys.week(year, week, habitIds),
    queryFn: () => calendarAPI.getWeek(year, week, habitIds),
    enabled: !!year && !!week,
  });
}

// Get year view
export function useCalendarYear(year: number, habitIds?: string[]) {
  return useQuery({
    queryKey: calendarKeys.year(year, habitIds),
    queryFn: () => calendarAPI.getYear(year, habitIds),
    enabled: !!year,
  });
}

