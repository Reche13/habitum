import { IconId } from "@/components/new-habit/icon-picker";
import { CategoryId, CATEGORIES } from "@/components/new-habit/select-category";

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

export function getIconEmoji(iconId: IconId): string {
  return ICONS.find((i) => i.id === iconId)?.emoji || "🔥";
}

export function getCategoryLabel(categoryId: CategoryId | string): string {
  return CATEGORIES.find((c) => c.id === categoryId)?.label || categoryId;
}

export function getCategoryIcon(categoryId: CategoryId | string): string {
  return CATEGORIES.find((c) => c.id === categoryId)?.icon || "📌";
}






