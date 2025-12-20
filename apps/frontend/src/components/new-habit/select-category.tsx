import { useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function SelectCategory() {
  const [value, setValue] = useState<CategoryId | null>(null);

  return (
    <Select
      value={value ?? undefined}
      onValueChange={(v) => setValue(v as CategoryId)}
    >
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Select a category" className="sr-only" />
      </SelectTrigger>

      <SelectContent>
        {CATEGORIES.map((cat) => (
          <SelectItem key={cat.id} value={cat.id}>
            <div className="flex items-center gap-2">
              <span className="text-lg">{cat.icon}</span>
              <div className="font-medium">{cat.label}</div>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

export const CATEGORIES = [
  {
    id: "health",
    label: "Health",
    icon: "🫀",
  },
  {
    id: "fitness",
    label: "Fitness",
    icon: "💪",
  },
  {
    id: "learning",
    label: "Learning",
    icon: "📘",
  },
  {
    id: "mindfulness",
    label: "Mindfulness",
    icon: "🧘",
  },
  {
    id: "productivity",
    label: "Productivity",
    icon: "⚡",
  },
] as const;

export type CategoryId = (typeof CATEGORIES)[number]["id"];
