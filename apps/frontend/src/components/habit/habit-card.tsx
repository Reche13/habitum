"use client";

import { useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { Habit } from "@/types/habit";
import {
  MoreVertical,
  Edit,
  Trash2,
  CheckCircle2,
  Calendar,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { getCategoryLabel, getCategoryIcon } from "@/lib/habit-utils";

interface HabitCardProps {
  habit: Habit;
  onEdit?: (habit: Habit) => void;
  onDelete?: (habitId: string) => void;
  onComplete?: (habitId: string) => void;
}

export function HabitCard({
  habit,
  onEdit,
  onDelete,
  onComplete,
}: HabitCardProps) {
  const router = useRouter();
  const streak = habit.currentStreak ?? 0;
  const completionRate = habit.completionRate ?? 0;
  const completedToday = habit.completedToday ?? false;

  const handleEdit = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    router.push(`/dashboard/habits/${habit.id}/edit`);
  };

  const handleCardClick = () => {
    router.push(`/dashboard/habits/${habit.id}`);
  };

  return (
    <div
      className="group relative rounded-xl border bg-background p-5 transition-all hover:shadow-md hover:border-primary/20 cursor-pointer"
      onClick={handleCardClick}
    >
      <div className="flex gap-4">
        {/* Icon */}
        <div
          className="h-14 w-14 rounded-xl flex items-center justify-center text-2xl shrink-0 transition-transform group-hover:scale-110"
          style={{ backgroundColor: habit.color + "20", color: habit.color }}
        >
          {habit.icon}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-base truncate">{habit.name}</h3>
              {habit.description && (
                <p className="text-sm text-muted-foreground line-clamp-2 mt-1">
                  {habit.description}
                </p>
              )}
            </div>

            {/* Actions Menu */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  className="opacity-0 group-hover:opacity-100 transition-opacity shrink-0"
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                  }}
                >
                  <MoreVertical className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align="end"
                onClick={(e) => e.stopPropagation()}
              >
                <DropdownMenuItem
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onComplete?.(habit.id);
                  }}
                  disabled={completedToday}
                >
                  <CheckCircle2 className="h-4 w-4 mr-2" />
                  {completedToday ? "Completed today" : "Mark complete"}
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    router.push(`/dashboard/habits/${habit.id}/edit`);
                  }}
                >
                  <Edit className="h-4 w-4 mr-2" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onDelete?.(habit.id);
                  }}
                  className="text-destructive focus:text-destructive"
                >
                  <Trash2 className="h-4 w-4 mr-2" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {/* Stats and Info */}
          <div className="mt-4 space-y-2">
            {/* Frequency and Category */}
            <div className="flex items-center gap-3 flex-wrap">
              <span className="text-xs text-muted-foreground">
                {habit.frequency === "daily"
                  ? "Daily"
                  : `${habit.timesPerWeek}× / week`}
              </span>
              {habit.category && (
                <>
                  <span className="text-xs text-muted-foreground">•</span>
                  <span className="text-xs rounded-full bg-muted px-2.5 py-1 flex items-center gap-1.5">
                    <span>{getCategoryIcon(habit.category)}</span>
                    {getCategoryLabel(habit.category)}
                  </span>
                </>
              )}
            </div>

            {/* Progress Indicators */}
            <div className="flex items-center gap-4 pt-2">
              {streak > 0 && (
                <div className="flex items-center gap-1.5">
                  <span className="text-lg">🔥</span>
                  <span className="text-xs font-medium">
                    {streak} day{streak !== 1 ? "s" : ""}
                  </span>
                </div>
              )}
              {completionRate > 0 && (
                <div className="flex items-center gap-1.5">
                  <Calendar className="h-3.5 w-3.5 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground">
                    {Math.round(completionRate)}% this week
                  </span>
                </div>
              )}
            </div>

            {/* Completion Status */}
            {completedToday && (
              <div className="flex items-center gap-1.5 text-xs text-green-600 dark:text-green-400 pt-1">
                <CheckCircle2 className="h-3.5 w-3.5" />
                <span>Completed today</span>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
