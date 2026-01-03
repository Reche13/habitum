"use client";

import { Sidebar } from "@/components/sidebar";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { ReactNode } from "react";

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <ProtectedRoute>
      <div className="w-full min-h-screen flex relative">
        <aside className="w-64 h-screen shrink-0 sticky top-0">
          <Sidebar />
        </aside>
        <main className="flex-1">{children}</main>
      </div>
    </ProtectedRoute>
  );
}
