"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

import { IconId, IconPicker } from "@/components/new-habit/icon-picker";
import { ColorPicker } from "@/components/new-habit/color-picker";
import { FrequencySelector } from "@/components/new-habit/frequency-selector";
import { SelectCategory } from "@/components/new-habit/select-category";

export default function NewHabit() {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [icon, setIcon] = useState<IconId>("fire");
  const [color, setColor] = useState<string>("#6366f1");
  const [frequency, setFrequency] = useState<"daily" | "weekly">("daily");
  const [timesPerWeek, setTimesPerWeek] = useState([3]);

  const isValid = name.trim().length > 0 && frequency !== null;

  return (
    <div className="w-full max-w-2xl px-20 py-12">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold">New Habit</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Create a habit you want to practice consistently.
        </p>
      </div>

      {/* Form */}
      <div className="space-y-10">
        {/* Name */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Habit name</label>
          <Input
            className="mt-2"
            placeholder="e.g. Exercise, Read 10 pages"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        {/* Description */}
        <div className="space-y-2">
          <label className="text-sm font-medium">
            Description{" "}
            <span className="text-muted-foreground font-normal">
              (optional)
            </span>
          </label>
          <Textarea
            className="mt-2"
            placeholder="What is this habit about?"
            rows={3}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>

        <div className="flex gap-20">
          <div className="flex items-center gap-2">
            <IconPicker value={icon} onChange={setIcon} />
            <label className="text-sm font-medium">Icon</label>
          </div>

          <div className="flex items-center gap-2">
            <ColorPicker value={color} onChange={setColor} />
            <label className="text-sm font-medium">Color</label>
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
          <div className="mt-2">
            <SelectCategory />
          </div>
        </div>

        {/* Frequency */}
        <div className="space-y-3">
          <label className="text-sm font-medium">Frequency</label>
          <div className="mt-2">
            <FrequencySelector
              frequency={frequency}
              onChange={setFrequency}
              sliderValue={timesPerWeek}
              onSliderValueChange={setTimesPerWeek}
            />
          </div>
        </div>

        {/* Actions */}
        <div className="flex justify-end gap-3 pt-6">
          <Button variant="ghost">Cancel</Button>
          <Button disabled={!isValid}>Create habit</Button>
        </div>
      </div>
    </div>
  );
}
