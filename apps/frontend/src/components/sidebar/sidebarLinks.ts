import { Calendar, ChartColumn, Home, Plus, Target } from "lucide-react";

export const SidebarLinks = [
  {
    label: "Home",
    href: "/dashboard",
    icon: Home,
  },
  {
    label: "Habits",
    href: "/dashboard/habits",
    icon: Target,
  },
  {
    label: "New Habit",
    href: "/dashboard/new-habit",
    icon: Plus,
  },
  {
    label: "Calendar",
    href: "/dashbaord/calendar",
    icon: Calendar,
  },
  {
    label: "Insights",
    href: "/dashboard/insights",
    icon: ChartColumn,
  },
];
