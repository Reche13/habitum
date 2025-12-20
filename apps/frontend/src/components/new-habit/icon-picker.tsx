import { useState } from "react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";

const ICONS = [
  { id: "fire", emoji: "🔥" },
  { id: "strength", emoji: "💪" },
  { id: "heart", emoji: "❤️" },
  { id: "water", emoji: "💧" },
  { id: "sun", emoji: "☀️" },
  { id: "moon", emoji: "🌙" },
  { id: "sleep", emoji: "😴" },
  { id: "energy", emoji: "⚡" },
] as const;

export type IconId = (typeof ICONS)[number]["id"];

export function IconPicker({
  value = "fire",
  onChange,
}: {
  value: IconId;
  onChange: (id: IconId) => void;
}) {
  const [open, setOpen] = useState(false);

  const selected = ICONS.find((i) => i.id === value)!;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          className="flex h-10 w-10 items-center justify-center rounded-lg border bg-background text-3xl transition hover:bg-muted"
        >
          {selected.emoji}
        </button>
      </PopoverTrigger>

      <PopoverContent className="w-72 p-2">
        <div className="grid grid-cols-6 gap-1">
          {ICONS.map((icon) => (
            <button
              key={icon.id}
              type="button"
              onClick={() => {
                onChange(icon.id);
                setOpen(false);
              }}
              className={cn(
                "flex h-10 w-10 items-center justify-center rounded-lg text-xl transition",
                value === icon.id ? "bg-primary/20" : "hover:bg-muted"
              )}
            >
              {icon.emoji}
            </button>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  );
}
