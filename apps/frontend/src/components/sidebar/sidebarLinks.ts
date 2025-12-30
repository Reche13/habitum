import {
  Calendar,
  ChartColumn,
  Home,
  Plus,
  Target,
  Settings,
  CheckCircle2,
} from "lucide-react";

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
    label: "Quick Mark",
    href: "/dashboard/quick-mark",
    icon: CheckCircle2,
  },
  {
    label: "Calendar",
    href: "/dashboard/calendar",
    icon: Calendar,
  },
  {
    label: "Insights",
    href: "/dashboard/insights",
    icon: ChartColumn,
  },
  {
    label: "Settings",
    href: "/dashboard/settings",
    icon: Settings,
  },
];
