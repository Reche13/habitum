import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { habitsAPI } from "@/lib/api/habits";
import type { HabitFilters, CreateHabitPayload, UpdateHabitPayload } from "@/lib/api/types";
import { dashboardKeys } from "./use-dashboard";

// Query keys
export const habitKeys = {
  all: ["habits"] as const,
  lists: () => [...habitKeys.all, "list"] as const,
  list: (filters?: HabitFilters) => [...habitKeys.lists(), filters] as const,
  details: () => [...habitKeys.all, "detail"] as const,
  detail: (id: string) => [...habitKeys.details(), id] as const,
};

// Get all habits
export function useHabits(filters?: HabitFilters) {
  return useQuery({
    queryKey: habitKeys.list(filters),
    queryFn: () => habitsAPI.getHabits(filters),
  });
}

// Get single habit
export function useHabit(id: string) {
  return useQuery({
    queryKey: habitKeys.detail(id),
    queryFn: () => habitsAPI.getHabit(id),
    enabled: !!id,
  });
}

// Create habit mutation
export function useCreateHabit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateHabitPayload) => habitsAPI.createHabit(payload),
    onSuccess: () => {
      // Invalidate habits list to refetch
      queryClient.invalidateQueries({ queryKey: habitKeys.lists() });
    },
  });
}

// Update habit mutation
export function useUpdateHabit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateHabitPayload }) =>
      habitsAPI.updateHabit(id, payload),
    onSuccess: (_, variables) => {
      // Invalidate specific habit and list
      queryClient.invalidateQueries({ queryKey: habitKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: habitKeys.lists() });
    },
  });
}

// Delete habit mutation
export function useDeleteHabit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => habitsAPI.deleteHabit(id),
    onSuccess: () => {
      // Invalidate habits list
      queryClient.invalidateQueries({ queryKey: habitKeys.lists() });
    },
  });
}

// Mark complete mutation
export function useMarkComplete() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, date }: { id: string; date?: string }) =>
      habitsAPI.markComplete(id, date),
    onSuccess: (_, variables) => {
      // Invalidate habit detail and list
      queryClient.invalidateQueries({ queryKey: habitKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: habitKeys.lists() });
      // Invalidate dashboard to refresh completion status
      queryClient.invalidateQueries({ queryKey: dashboardKeys.all });
    },
  });
}

// Unmark complete mutation
export function useUnmarkComplete() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, date }: { id: string; date?: string }) =>
      habitsAPI.unmarkComplete(id, date),
    onSuccess: (_, variables) => {
      // Invalidate habit detail and list
      queryClient.invalidateQueries({ queryKey: habitKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: habitKeys.lists() });
      // Invalidate dashboard to refresh completion status
      queryClient.invalidateQueries({ queryKey: dashboardKeys.all });
    },
  });
}

// Get completion history hook
export function useCompletionHistory(id: string, year?: number, allTime?: boolean) {
  return useQuery({
    queryKey: [...habitKeys.detail(id), "completion-history", year, allTime],
    queryFn: () => habitsAPI.getCompletionHistory(id, year, allTime),
    enabled: !!id,
  });
}

