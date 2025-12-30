import { Slider } from "@/components/ui/slider";

interface Props {
  value: number[];
  onChange: (v: number[]) => void;
}

export function TimesPerWeekSlider({ value, onChange }: Props) {
  return (
    <div className="space-y-3 w-full">
      <Slider
        value={value}
        onValueChange={onChange}
        min={1}
        max={7}
        step={1}
        className="w-full cursor-pointer"
      />

      <div className="flex justify-between text-xs text-muted-foreground">
        {Array.from({ length: 7 }, (_, i) => (
          <span key={i}>{i + 1}</span>
        ))}
      </div>
    </div>
  );
}
