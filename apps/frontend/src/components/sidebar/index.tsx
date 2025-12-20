"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { SidebarLinks } from "./sidebarLinks";

export const Sidebar = () => {
  const pathname = usePathname();

  const isActive = (href: string) => {
    if (!pathname) return false;

    if (href === "/dashboard") {
      return pathname === href;
    }

    return pathname.startsWith(href);
  };

  return (
    <div className="w-full h-full border-r bg-background flex flex-col px-3 py-6">
      <div className="text-xl font-semibold text-foreground px-4">Habitum.</div>

      <nav className="flex flex-col gap-1 mt-6">
        {SidebarLinks.map((item, index) => {
          const active = isActive(item.href);
          const newHabitHref = item.href === "/dashboard/new-habit";

          return (
            <Link
              href={item.href}
              key={`${item.href}-${index}`}
              className={cn(
                "flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-colors",
                "text-muted-foreground hover:bg-accent hover:text-foreground",
                active && "bg-accent text-foreground",
                newHabitHref &&
                  "bg-foreground hover:bg-secondary-foreground text-background hover:text-background shadow-md"
              )}
            >
              <item.icon size={18} />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
};
