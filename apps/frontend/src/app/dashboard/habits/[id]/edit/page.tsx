"use client";

import { use, useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { AlertCircle, Loader2, ArrowLeft } from "lucide-react";
import { cn } from "@/lib/utils";
import Link from "next/link";

import { IconId, IconPicker } from "@/components/new-habit/icon-picker";
import { ColorPicker } from "@/components/new-habit/color-picker";
import { FrequencySelector } from "@/components/new-habit/frequency-selector";
import { SelectCategory, CategoryId } from "@/components/new-habit/select-category";
import { HabitPreviewCard } from "@/components/new-habit/habit-preview-card";
import { useHabit, useUpdateHabit } from "@/lib/hooks";
import { mapHabitResponseToHabit } from "@/lib/api/mappers";
import type { UpdateHabitPayload } from "@/lib/api/types";

interface FormErrors {
  name?: string;
  description?: string;
  frequency?: string;
  timesPerWeek?: string;
}

const MAX_NAME_LENGTH = 50;
const MAX_DESCRIPTION_LENGTH = 200;

export default function EditHabitPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const nameInputRef = useRef<HTMLInputElement>(null);
  
  const { data: habitResponse, isLoading, error } = useHabit(id);
  const updateHabit = useUpdateHabit();

  const habit = habitResponse ? mapHabitResponseToHabit(habitResponse) : null;

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [icon, setIcon] = useState<IconId>("fire");
  const [color, setColor] = useState<string>("#6366f1");
  const [category, setCategory] = useState<CategoryId | null>(null);
  const [frequency, setFrequency] = useState<"daily" | "weekly">("daily");
  const [timesPerWeek, setTimesPerWeek] = useState([3]);
  const [errors, setErrors] = useState<FormErrors>({});
  const [apiError, setApiError] = useState<string | null>(null);

  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description || "");
      setIcon(habit.iconId || "fire");
      setColor(habit.color || "#6366f1");
      setCategory(habit.category || null);
      setFrequency(habit.frequency);
      setTimesPerWeek([habit.timesPerWeek || 3]);
    }
    nameInputRef.current?.focus();
  }, [habit]);

  if (isLoading) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12 flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !habit) {
    return (
      <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
        <div className="max-w-4xl mx-auto text-center py-16">
          <h1 className="text-2xl font-semibold mb-2">Habit not found</h1>
          <p className="text-muted-foreground mb-6">
            {error instanceof Error ? error.message : "The habit you're trying to edit doesn't exist."}
          </p>
          <Button asChild>
            <Link href="/dashboard/habits">Back to Habits</Link>
          </Button>
        </div>
      </div>
    );
  }

  const validate = (): boolean => {
    const newErrors: FormErrors = {};

    if (!name.trim()) {
      newErrors.name = "Habit name is required";
    } else if (name.trim().length < 2) {
      newErrors.name = "Name must be at least 2 characters";
    } else if (name.length > MAX_NAME_LENGTH) {
      newErrors.name = `Name must be less than ${MAX_NAME_LENGTH} characters`;
    }

    if (description.length > MAX_DESCRIPTION_LENGTH) {
      newErrors.description = `Description must be less than ${MAX_DESCRIPTION_LENGTH} characters`;
    }

    if (!frequency) {
      newErrors.frequency = "Frequency is required";
    }

    if (frequency === "weekly") {
      const times = timesPerWeek[0];
      if (!times || times < 1 || times > 7) {
        newErrors.timesPerWeek = "Must be between 1 and 7 times per week";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setApiError(null);

    if (!validate()) {
      return;
    }

    try {
      const payload: UpdateHabitPayload = {
        name: name.trim(),
        description: description.trim() || undefined,
        icon: icon,
        color: color,
        frequency: frequency,
        times_per_week: frequency === "weekly" ? timesPerWeek[0] : undefined,
        category: category || undefined,
      };

      await updateHabit.mutateAsync({ id, payload });
      router.push(`/dashboard/habits/${id}`);
    } catch (error: any) {
      setApiError(error?.message || "Failed to update habit. Please try again.");
      console.error("Error updating habit:", error);
    }
  };

  const isValid = name.trim().length >= 2 && !Object.keys(errors).length;

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Button variant="ghost" onClick={() => router.back()} className="mb-4">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
          <h1 className="text-2xl sm:text-3xl font-semibold">Edit Habit</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Update your habit details and preferences.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Form Section */}
          <div>
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Name */}
              <div className="space-y-2">
                <label htmlFor="name" className="text-sm font-medium">
                  Habit name <span className="text-destructive">*</span>
                </label>
                <Input
                  id="name"
                  ref={nameInputRef}
                  placeholder="e.g. Exercise, Read 10 pages"
                  value={name}
                  onChange={(e) => {
                    setName(e.target.value);
                    if (errors.name) {
                      setErrors((prev) => ({ ...prev, name: undefined }));
                    }
                  }}
                  onBlur={() => validate()}
                  className={cn(errors.name && "border-destructive")}
                  maxLength={MAX_NAME_LENGTH}
                />
                <div className="flex items-center justify-between">
                  {errors.name ? (
                    <p className="text-xs text-destructive flex items-center gap-1">
                      <AlertCircle className="h-3 w-3" />
                      {errors.name}
                    </p>
                  ) : (
                    <div />
                  )}
                  <p className="text-xs text-muted-foreground">
                    {name.length}/{MAX_NAME_LENGTH}
                  </p>
                </div>
              </div>

              {/* Description */}
              <div className="space-y-2">
                <label htmlFor="description" className="text-sm font-medium">
                  Description{" "}
                  <span className="text-muted-foreground font-normal">
                    (optional)
                  </span>
                </label>
                <Textarea
                  id="description"
                  placeholder="What is this habit about?"
                  rows={3}
                  value={description}
                  onChange={(e) => {
                    setDescription(e.target.value);
                    if (errors.description) {
                      setErrors((prev) => ({ ...prev, description: undefined }));
                    }
                  }}
                  onBlur={() => validate()}
                  className={cn(errors.description && "border-destructive")}
                  maxLength={MAX_DESCRIPTION_LENGTH}
                />
                <div className="flex items-center justify-between">
                  {errors.description ? (
                    <p className="text-xs text-destructive flex items-center gap-1">
                      <AlertCircle className="h-3 w-3" />
                      {errors.description}
                    </p>
                  ) : (
                    <div />
                  )}
                  <p className="text-xs text-muted-foreground">
                    {description.length}/{MAX_DESCRIPTION_LENGTH}
                  </p>
                </div>
              </div>

              {/* Icon & Color */}
              <div className="space-y-2">
                <label className="text-sm font-medium">Appearance</label>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Icon</label>
                    <div className="flex items-center gap-3">
                      <IconPicker value={icon} onChange={setIcon} />
                      <span className="text-sm text-muted-foreground">
                        Choose an icon
                      </span>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Color</label>
                    <div className="flex items-center gap-3">
                      <ColorPicker value={color} onChange={setColor} />
                      <span className="text-sm text-muted-foreground">
                        Choose a color
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Category */}
              <div className="space-y-2">
                <label className="text-sm font-medium">
                  Category{" "}
                  <span className="text-muted-foreground font-normal">
                    (optional)
                  </span>
                </label>
                <SelectCategory value={category} onChange={setCategory} />
              </div>

              {/* Frequency */}
              <div className="space-y-3">
                <label className="text-sm font-medium">
                  Frequency <span className="text-destructive">*</span>
                </label>
                <FrequencySelector
                  frequency={frequency}
                  onChange={setFrequency}
                  sliderValue={timesPerWeek}
                  onSliderValueChange={setTimesPerWeek}
                />
                {errors.timesPerWeek && (
                  <p className="text-xs text-destructive flex items-center gap-1">
                    <AlertCircle className="h-3 w-3" />
                    {errors.timesPerWeek}
                  </p>
                )}
              </div>

              {/* Actions */}
              <div className="flex justify-end gap-3 pt-4">
                {apiError && (
                  <div className="rounded-lg border border-destructive bg-destructive/10 p-4 flex items-center gap-2">
                    <AlertCircle className="h-4 w-4 text-destructive" />
                    <p className="text-sm text-destructive">{apiError}</p>
                  </div>
                )}

                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => router.back()}
                  disabled={updateHabit.isPending}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={!isValid || updateHabit.isPending}>
                  {updateHabit.isPending ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    "Save changes"
                  )}
                </Button>
              </div>
            </form>
          </div>

          {/* Preview Section */}
          <div className="lg:sticky lg:top-8 lg:h-fit">
            <HabitPreviewCard
              name={name}
              description={description}
              iconId={icon}
              color={color}
              category={category}
              frequency={frequency}
              timesPerWeek={timesPerWeek[0]}
            />
          </div>
        </div>
      </div>
    </div>
  );
}






