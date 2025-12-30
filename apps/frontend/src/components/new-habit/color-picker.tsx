import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useState } from "react";
import { Check } from "lucide-react";

const COLORS = [
  "#6366f1",
  "#3b82f6",
  "#06b6d4",
  "#10b981",
  "#22c55e",
  "#eab308",
  "#f97316",
  "#ef4444",
  "#ec4899",
  "#8b5cf6",
  "#64748b",
  "#0f172a",
];

export function ColorPicker({
  value,
  onChange,
}: {
  value: string;
  onChange: (value: string) => void;
}) {
  const [open, setOpen] = useState(false);
  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className="w-10 h-10 rounded-lg border-none cursor-pointer"
          style={{ background: value }}
        ></Button>
      </PopoverTrigger>

      <PopoverContent
        className="w-fit p-3 border border-zinc-200"
        align="start"
      >
        <div className="grid grid-cols-5 gap-2">
          {COLORS.map((color) => (
            <button
              key={color}
              type="button"
              onClick={() => {
                onChange(color);
                setOpen(false);
              }}
              className={cn(
                "h-10 w-10 rounded-lg transition flex items-center justify-center p-0.5"
              )}
              style={{ backgroundColor: color }}
            >
              {value === color && (
                <div className="flex items-center justify-center w-full h-full rounded-lg border border-background">
                  <Check className="text-background size-5" />
                </div>
              )}
            </button>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  );
}
