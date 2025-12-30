import { apiClient } from "./client";
import type { CalendarCompletionsResponse } from "./types";

// Calendar API
export const calendarAPI = {
  // Get completions for date range
  getCompletions: async (startDate: string, endDate: string, habitIds?: string[]) => {
    const params = new URLSearchParams({
      startDate,
      endDate,
    });
    if (habitIds && habitIds.length > 0) {
      params.append("habitIds", habitIds.join(","));
    }
    const response = await apiClient.get<CalendarCompletionsResponse>(
      `/calendar/completions?${params.toString()}`
    );
    return response.data;
  },

  // Get month view
  getMonth: async (year: number, month: number, habitIds?: string[]) => {
    const params = new URLSearchParams({
      year: year.toString(),
      month: month.toString(),
    });
    if (habitIds && habitIds.length > 0) {
      params.append("habitIds", habitIds.join(","));
    }
    const response = await apiClient.get(`/calendar/month?${params.toString()}`);
    return response.data;
  },

  // Get week view
  getWeek: async (year: number, week: number, habitIds?: string[]) => {
    const params = new URLSearchParams({
      year: year.toString(),
      week: week.toString(),
    });
    if (habitIds && habitIds.length > 0) {
      params.append("habitIds", habitIds.join(","));
    }
    const response = await apiClient.get(`/calendar/week?${params.toString()}`);
    return response.data;
  },

  // Get year view
  getYear: async (year: number, habitIds?: string[]) => {
    const params = new URLSearchParams({
      year: year.toString(),
    });
    if (habitIds && habitIds.length > 0) {
      params.append("habitIds", habitIds.join(","));
    }
    const response = await apiClient.get(`/calendar/year?${params.toString()}`);
    return response.data;
  },
};

