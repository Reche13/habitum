import { cn } from "@/lib/utils";
import { TimesPerWeekSlider } from "./week-slider";

type Frequency = "daily" | "weekly";

interface Props {
  frequency: Frequency;
  onChange: (v: Frequency) => void;
  sliderValue: number[];
  onSliderValueChange: (v: number[]) => void;
}

export function FrequencySelector({
  frequency,
  onChange,
  sliderValue,
  onSliderValueChange,
}: Props) {
  return (
    <div className="space-y-3">
      <div className="grid grid-cols-2 rounded-xl border bg-muted p-1">
        <button
          type="button"
          onClick={() => onChange("daily")}
          className={cn(
            "rounded-lg px-4 py-2 text-sm transition",
            frequency === "daily"
              ? "bg-background shadow-sm font-medium"
              : "text-muted-foreground"
          )}
        >
          <div>Daily</div>
          <div className="text-xs text-muted-foreground">Every day</div>
        </button>

        <button
          type="button"
          onClick={() => onChange("weekly")}
          className={cn(
            "rounded-lg px-4 py-2 text-sm transition",
            frequency === "weekly"
              ? "bg-background shadow-sm font-medium"
              : "text-muted-foreground"
          )}
        >
          <div>Weekly</div>
          <div className="text-xs text-muted-foreground">Flexible</div>
        </button>
      </div>

      {frequency === "weekly" && (
        <div className="mt-8">
          <TimesPerWeekSlider
            value={sliderValue}
            onChange={onSliderValueChange}
          />
        </div>
      )}
    </div>
  );
}
