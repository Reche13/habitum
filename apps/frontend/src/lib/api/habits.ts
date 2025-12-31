import { apiClient } from "./client";
import type {
  HabitResponse,
  CreateHabitPayload,
  UpdateHabitPayload,
  HabitFilters,
  PaginatedResponse,
  APIResponse,
} from "./types";

// Habits API
export const habitsAPI = {
  // Get all habits with filters
  getHabits: async (filters?: HabitFilters) => {
    const params = new URLSearchParams();
    if (filters?.category) params.append("category", filters.category);
    if (filters?.search) params.append("search", filters.search);
    if (filters?.sort) params.append("sort", filters.sort);
    if (filters?.order) params.append("order", filters.order);
    if (filters?.page) params.append("page", filters.page.toString());
    if (filters?.limit) params.append("limit", filters.limit.toString());

    const response = await apiClient.get<PaginatedResponse<HabitResponse>>(
      `/habits${params.toString() ? `?${params.toString()}` : ""}`
    );
    return response.data;
  },

  // Get single habit
  getHabit: async (id: string) => {
    const response = await apiClient.get<APIResponse<HabitResponse>>(`/habits/${id}`);
    return response.data.data; // Unwrap APIResponse
  },

  // Create habit
  createHabit: async (payload: CreateHabitPayload) => {
    const response = await apiClient.post<APIResponse<HabitResponse>>("/habits", payload);
    return response.data.data; // Unwrap APIResponse
  },

  // Update habit
  updateHabit: async (id: string, payload: UpdateHabitPayload) => {
    const response = await apiClient.patch<APIResponse<HabitResponse>>(`/habits/${id}`, payload);
    return response.data.data; // Unwrap APIResponse
  },

  // Delete habit
  deleteHabit: async (id: string) => {
    await apiClient.delete(`/habits/${id}`);
  },

  // Mark habit as complete
  markComplete: async (id: string, date?: string) => {
    const response = await apiClient.post(`/habits/${id}/complete`, {
      logDate: date || new Date().toISOString().split("T")[0],
    });
    return response.data;
  },

  // Unmark completion
  unmarkComplete: async (id: string, date?: string) => {
    const dateParam = date || new Date().toISOString().split("T")[0];
    await apiClient.delete(`/habits/${id}/complete?date=${dateParam}`);
  },

  // Get completions
  getCompletions: async (id: string, startDate: string, endDate: string, limit?: number) => {
    const params = new URLSearchParams({
      startDate,
      endDate,
    });
    if (limit) params.append("limit", limit.toString());
    const response = await apiClient.get(`/habits/${id}/completions?${params.toString()}`);
    return response.data;
  },

  // Get completion history
  getCompletionHistory: async (id: string, year?: number, allTime?: boolean) => {
    const params = new URLSearchParams();
    if (year) params.append("year", year.toString());
    if (allTime) params.append("allTime", "true");
    const response = await apiClient.get<APIResponse<{ dates: string[]; totalDays: number; completedDays: number }>>(
      `/habits/${id}/completion-history?${params.toString()}`
    );
    return response.data.data; // Unwrap APIResponse
  },
};


