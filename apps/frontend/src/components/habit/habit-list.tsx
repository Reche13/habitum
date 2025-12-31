"use client";

import { useMemo, useState } from "react";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { HabitCard } from "./habit-card";
import { CATEGORIES, CategoryId } from "@/components/new-habit/select-category";
import {
  Grid3x3,
  List,
  X,
  Plus,
  Target,
  TrendingUp,
  Calendar,
  Loader2,
} from "lucide-react";
import Link from "next/link";
import { useHabits, useDeleteHabit } from "@/lib/hooks";
import { mapHabitResponsesToHabits } from "@/lib/api/mappers";
import type { HabitFilters } from "@/lib/api/types";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

type SortOption = "name" | "date" | "streak" | "completion";
type ViewMode = "grid" | "list";

export function HabitsList() {
  const [query, setQuery] = useState("");
  const [category, setCategory] = useState<CategoryId | "all">("all");
  const [sortBy, setSortBy] = useState<SortOption>("name");
  const [viewMode, setViewMode] = useState<ViewMode>("grid");
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [habitToDelete, setHabitToDelete] = useState<{
    id: string;
    name: string;
  } | null>(null);

  const deleteHabit = useDeleteHabit();

  // Build filters for API
  const filters: HabitFilters = useMemo(() => {
    const apiFilters: HabitFilters = {
      search: query || undefined,
      category: category !== "all" ? category : undefined,
      sort: sortBy === "date" ? undefined : sortBy,
      order: "asc",
    };

    // Map frontend sort to backend sort
    if (sortBy === "streak") {
      apiFilters.sort = "streak";
    } else if (sortBy === "completion") {
      apiFilters.sort = "completion";
    } else if (sortBy === "name") {
      apiFilters.sort = "name";
    }

    return apiFilters;
  }, [query, category, sortBy]);

  const { data, isLoading, error } = useHabits(filters);

  // Map backend response to frontend habits
  const habits = useMemo(() => {
    if (!data?.data) return [];
    return mapHabitResponsesToHabits(data.data);
  }, [data]);

  // Calculate statistics
  const stats = useMemo(() => {
    if (!habits.length) {
      return {
        total: 0,
        activeStreak: 0,
        avgCompletion: 0,
        completedTodayCount: 0,
      };
    }

    const total = habits.length;
    const activeStreak = Math.max(
      ...habits.map((h) => h.currentStreak ?? 0),
      0
    );
    const avgCompletion = Math.round(
      habits.reduce((sum, h) => sum + (h.completionRate ?? 0), 0) / total || 0
    );
    const completedTodayCount = habits.filter((h) => h.completedToday).length;

    return {
      total,
      activeStreak,
      avgCompletion,
      completedTodayCount,
    };
  }, [habits]);

  const hasActiveFilters = query !== "" || category !== "all";

  const handleClearFilters = () => {
    setQuery("");
    setCategory("all");
  };

  const handleDeleteClick = (habitId: string, habitName: string) => {
    setHabitToDelete({ id: habitId, name: habitName });
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!habitToDelete) return;

    try {
      await deleteHabit.mutateAsync(habitToDelete.id);
      setDeleteDialogOpen(false);
      setHabitToDelete(null);
    } catch (error) {
      console.error("Failed to delete habit:", error);
      alert("Failed to delete habit. Please try again.");
    }
  };

  return (
    <>
      {/* Statistics Dashboard */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div className="rounded-lg border bg-background p-4">
          <div className="flex items-center gap-2 mb-1">
            <Target className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">Total Habits</span>
          </div>
          <p className="text-2xl font-semibold">{stats.total}</p>
        </div>
        <div className="rounded-lg border bg-background p-4">
          <div className="flex items-center gap-2 mb-1">
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">Active Streak</span>
          </div>
          <p className="text-2xl font-semibold">{stats.activeStreak} days</p>
        </div>
        <div className="rounded-lg border bg-background p-4">
          <div className="flex items-center gap-2 mb-1">
            <Calendar className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              Avg. Completion
            </span>
          </div>
          <p className="text-2xl font-semibold">{stats.avgCompletion}%</p>
        </div>
        <div className="rounded-lg border bg-background p-4">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-sm text-muted-foreground">
              Completed Today
            </span>
          </div>
          <p className="text-2xl font-semibold">
            {stats.completedTodayCount}/{stats.total}
          </p>
        </div>
      </div>

      {/* Controls */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 mb-6">
        <div className="flex items-center gap-3 flex-1 w-full sm:w-auto">
          <Input
            placeholder="Search habits..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="max-w-sm"
          />
          {hasActiveFilters && (
            <Button
              variant="ghost"
              size="sm"
              onClick={handleClearFilters}
              className="gap-2"
            >
              <X className="h-4 w-4" />
              Clear
            </Button>
          )}
        </div>

        <div className="flex items-center gap-3">
          <Select
            value={category}
            onValueChange={(v) => setCategory(v as CategoryId | "all")}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All categories</SelectItem>
              {CATEGORIES.map((cat) => (
                <SelectItem key={cat.id} value={cat.id}>
                  <div className="flex items-center gap-2">
                    <span>{cat.icon}</span>
                    <span>{cat.label}</span>
                  </div>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={sortBy}
            onValueChange={(v) => setSortBy(v as SortOption)}
          >
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="name">Name (A-Z)</SelectItem>
              <SelectItem value="streak">Streak (High)</SelectItem>
              <SelectItem value="completion">Completion</SelectItem>
              <SelectItem value="date">Date Created</SelectItem>
            </SelectContent>
          </Select>

          <div className="flex items-center gap-1 border rounded-md p-1">
            <Button
              variant={viewMode === "grid" ? "secondary" : "ghost"}
              size="icon-sm"
              onClick={() => setViewMode("grid")}
            >
              <Grid3x3 className="h-4 w-4" />
            </Button>
            <Button
              variant={viewMode === "list" ? "secondary" : "ghost"}
              size="icon-sm"
              onClick={() => setViewMode("list")}
            >
              <List className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      {/* List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <div className="rounded-full bg-muted p-6 mb-4">
            <Target className="h-12 w-12 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold mb-2">Error loading habits</h3>
          <p className="text-sm text-muted-foreground mb-6 max-w-sm">
            {error instanceof Error ? error.message : "Failed to load habits"}
          </p>
        </div>
      ) : habits.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 text-center">
          {!hasActiveFilters ? (
            <>
              <div className="rounded-full bg-muted p-6 mb-4">
                <Target className="h-12 w-12 text-muted-foreground" />
              </div>
              <h3 className="text-lg font-semibold mb-2">No habits yet</h3>
              <p className="text-sm text-muted-foreground mb-6 max-w-sm">
                Start building better habits by creating your first one.
              </p>
              <Button asChild>
                <Link href="/dashboard/new-habit" className="gap-2">
                  <Plus className="h-4 w-4" />
                  Create your first habit
                </Link>
              </Button>
            </>
          ) : (
            <>
              <div className="rounded-full bg-muted p-6 mb-4">
                <Target className="h-12 w-12 text-muted-foreground" />
              </div>
              <h3 className="text-lg font-semibold mb-2">No habits found</h3>
              <p className="text-sm text-muted-foreground mb-6 max-w-sm">
                Try adjusting your search or filters to find what you're looking
                for.
              </p>
              <Button variant="outline" onClick={handleClearFilters}>
                Clear filters
              </Button>
            </>
          )}
        </div>
      ) : (
        <div
          className={
            viewMode === "grid"
              ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
              : "space-y-4"
          }
        >
          {habits.map((habit) => (
            <HabitCard
              key={habit.id}
              habit={habit}
              onDelete={(habitId) => handleDeleteClick(habitId, habit.name)}
            />
          ))}
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Habit</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{habitToDelete?.name}"? This
              action cannot be undone and all completion history will be
              permanently deleted.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setHabitToDelete(null)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteConfirm}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={deleteHabit.isPending}
            >
              {deleteHabit.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
