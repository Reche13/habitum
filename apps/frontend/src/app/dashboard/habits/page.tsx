import { HabitsList } from "@/components/habit/habit-list";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import Link from "next/link";

export default function HabitsPage() {
  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      {/* Header */}
      <div className="mb-8 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-semibold">All Habits</h1>
          <p className="text-sm text-muted-foreground mt-1">
            View and manage all your habits.
          </p>
        </div>
        <Button asChild className="gap-2 shrink-0">
          <Link href="/dashboard/new-habit">
            <Plus className="h-4 w-4" />
            <span className="hidden sm:inline">Create Habit</span>
            <span className="sm:hidden">New</span>
          </Link>
        </Button>
      </div>

      <HabitsList />
    </div>
  );
}
