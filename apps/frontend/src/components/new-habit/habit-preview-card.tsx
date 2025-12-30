import { Habit } from "@/types/habit";
import { IconId } from "./icon-picker";
import { CategoryId } from "./select-category";
import { getIconEmoji, getCategoryLabel, getCategoryIcon } from "@/lib/habit-utils";
import { cn } from "@/lib/utils";

interface HabitPreviewCardProps {
  name: string;
  description?: string;
  iconId: IconId;
  color: string;
  category?: CategoryId | null;
  frequency: "daily" | "weekly";
  timesPerWeek?: number;
}

export function HabitPreviewCard({
  name,
  description,
  iconId,
  color,
  category,
  frequency,
  timesPerWeek,
}: HabitPreviewCardProps) {
  const icon = getIconEmoji(iconId);
  const hasContent = name.trim().length > 0;

  return (
    <div className="rounded-xl border bg-background p-5">
      <div className="mb-4">
        <h3 className="text-sm font-medium text-muted-foreground mb-2">
          Preview
        </h3>
        <div className="h-px bg-border" />
      </div>

      {hasContent ? (
        <div className="rounded-xl border bg-background p-5 flex gap-4">
          {/* Icon */}
          <div
            className="h-14 w-14 rounded-xl flex items-center justify-center text-2xl shrink-0 transition-transform"
            style={{ backgroundColor: color + "20", color: color }}
          >
            {icon}
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <div className="flex-1 min-w-0">
                <h3 className="font-semibold text-base truncate">{name}</h3>
                {description && (
                  <p className="text-sm text-muted-foreground line-clamp-2 mt-1">
                    {description}
                  </p>
                )}
              </div>
            </div>

            {/* Info */}
            <div className="mt-4 space-y-2">
              {/* Frequency and Category */}
              <div className="flex items-center gap-3 flex-wrap">
                <span className="text-xs text-muted-foreground">
                  {frequency === "daily"
                    ? "Daily"
                    : `${timesPerWeek || 3}× / week`}
                </span>
                {category && (
                  <>
                    <span className="text-xs text-muted-foreground">•</span>
                    <span className="text-xs rounded-full bg-muted px-2.5 py-1 flex items-center gap-1.5">
                      <span>{getCategoryIcon(category)}</span>
                      {getCategoryLabel(category)}
                    </span>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className="rounded-xl border border-dashed bg-muted/30 p-8 text-center">
          <p className="text-sm text-muted-foreground">
            Start filling the form to see a preview
          </p>
        </div>
      )}
    </div>
  );
}






