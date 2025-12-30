import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface SelectCategoryProps {
  value?: CategoryId | null;
  onChange?: (value: CategoryId | null) => void;
}

export function SelectCategory({ value = null, onChange }: SelectCategoryProps) {
  return (
    <Select
      value={value ?? undefined}
      onValueChange={(v) => {
        if (v === "none" || v === "") {
          onChange?.(null);
        } else {
          onChange?.(v as CategoryId);
        }
      }}
    >
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Select a category (optional)" />
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
    id: "productivity",
    label: "Productivity",
    icon: "⚡",
  },
  {
    id: "learning",
    label: "Learning",
    icon: "📘",
  },
  {
    id: "work",
    label: "Work",
    icon: "💼",
  },
  {
    id: "personal",
    label: "Personal",
    icon: "👤",
  },
  {
    id: "mindfulness",
    label: "Mindfulness",
    icon: "🧘",
  },
  {
    id: "social",
    label: "Social",
    icon: "👥",
  },
  {
    id: "creative",
    label: "Creative",
    icon: "🎨",
  },
  {
    id: "finance",
    label: "Finance",
    icon: "💰",
  },
  {
    id: "other",
    label: "Other",
    icon: "📌",
  },
] as const;

export type CategoryId = (typeof CATEGORIES)[number]["id"];
