import { useQuery } from "@tanstack/react-query";
import { dashboardAPI } from "@/lib/api/dashboard";

// Query keys
export const dashboardKeys = {
  all: ["dashboard"] as const,
  home: () => [...dashboardKeys.all, "home"] as const,
};

// Get dashboard home data
export function useDashboardHome() {
  return useQuery({
    queryKey: dashboardKeys.home(),
    queryFn: () => dashboardAPI.getHome(),
    refetchOnWindowFocus: true, // Refetch when user comes back to tab
    staleTime: 30000, // Consider data fresh for 30 seconds
  });
}

