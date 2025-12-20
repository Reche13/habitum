import { Sidebar } from "@/components/sidebar";
import { redirect } from "next/navigation";
import { ReactNode } from "react";

export default function DashboardLayout({ children }: { children: ReactNode }) {
  const isAuthenticated = true;

  if (!isAuthenticated) {
    redirect("/login");
  }
  return (
    <div className="w-full min-h-screen flex relative">
      <aside className="w-64 h-screen shrink-0 sticky top-0">
        <Sidebar />
      </aside>
      <main className="flex-1">{children}</main>
    </div>
  );
}
