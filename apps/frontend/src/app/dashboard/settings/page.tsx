"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  User,
  Bell,
  Palette,
  Database,
  Shield,
  LogOut,
  Save,
  Loader2,
} from "lucide-react";
import { cn } from "@/lib/utils";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState("profile");
  const [isSaving, setIsSaving] = useState(false);

  // Profile state
  const [name, setName] = useState("John Doe");
  const [email, setEmail] = useState("john@example.com");

  // Preferences state
  const [theme, setTheme] = useState("system");
  const [notifications, setNotifications] = useState(true);
  const [reminderTime, setReminderTime] = useState("09:00");

  const handleSave = async () => {
    setIsSaving(true);
    await new Promise((resolve) => setTimeout(resolve, 1000));
    setIsSaving(false);
    // TODO: API call to save settings
  };

  const tabs = [
    { id: "profile", label: "Profile", icon: User },
    { id: "preferences", label: "Preferences", icon: Palette },
    { id: "notifications", label: "Notifications", icon: Bell },
    { id: "data", label: "Data & Privacy", icon: Database },
  ];

  return (
    <div className="w-full px-4 sm:px-6 lg:px-20 py-8 sm:py-12">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-2xl sm:text-3xl font-semibold">Settings</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your account settings and preferences.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="space-y-1 rounded-lg border bg-background p-2">
              {tabs.map((tab) => {
                const Icon = tab.icon;
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={cn(
                      "w-full flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                      activeTab === tab.id
                        ? "bg-accent text-accent-foreground"
                        : "text-muted-foreground hover:bg-accent/50 hover:text-foreground"
                    )}
                  >
                    <Icon className="h-4 w-4" />
                    <span>{tab.label}</span>
                  </button>
                );
              })}
            </div>

            <div className="mt-6">
              <Button variant="outline" className="w-full" size="sm">
                <LogOut className="h-4 w-4 mr-2" />
                Sign out
              </Button>
            </div>
          </div>

          {/* Content */}
          <div className="lg:col-span-3">
            <div className="rounded-lg border bg-background p-6">
              {/* Profile Tab */}
              {activeTab === "profile" && (
                <div className="space-y-6">
                  <div>
                    <h2 className="text-lg font-semibold mb-1">Profile</h2>
                    <p className="text-sm text-muted-foreground">
                      Update your personal information.
                    </p>
                  </div>

                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="name">Name</Label>
                      <Input
                        id="name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="Your name"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="email">Email</Label>
                      <Input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="your@email.com"
                        disabled
                      />
                      <p className="text-xs text-muted-foreground">
                        Email cannot be changed
                      </p>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="password">Change Password</Label>
                      <Input
                        id="password"
                        type="password"
                        placeholder="Enter new password"
                      />
                    </div>
                  </div>
                </div>
              )}

              {/* Preferences Tab */}
              {activeTab === "preferences" && (
                <div className="space-y-6">
                  <div>
                    <h2 className="text-lg font-semibold mb-1">Preferences</h2>
                    <p className="text-sm text-muted-foreground">
                      Customize your app experience.
                    </p>
                  </div>

                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label>Theme</Label>
                      <div className="grid grid-cols-3 gap-2">
                        {["light", "dark", "system"].map((option) => (
                          <button
                            key={option}
                            onClick={() => setTheme(option)}
                            className={cn(
                              "rounded-lg border p-3 text-sm font-medium transition-colors",
                              theme === option
                                ? "border-primary bg-primary/10 text-primary"
                                : "border-border hover:bg-accent"
                            )}
                          >
                            {option.charAt(0).toUpperCase() + option.slice(1)}
                          </button>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Notifications Tab */}
              {activeTab === "notifications" && (
                <div className="space-y-6">
                  <div>
                    <h2 className="text-lg font-semibold mb-1">
                      Notifications
                    </h2>
                    <p className="text-sm text-muted-foreground">
                      Manage your notification preferences.
                    </p>
                  </div>

                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <Label>Enable Notifications</Label>
                        <p className="text-sm text-muted-foreground">
                          Receive reminders for your habits
                        </p>
                      </div>
                      <input
                        type="checkbox"
                        checked={notifications}
                        onChange={(e) => setNotifications(e.target.checked)}
                        className="h-4 w-4 rounded border-gray-300"
                      />
                    </div>

                    {notifications && (
                      <div className="space-y-2">
                        <Label htmlFor="reminder-time">Default Reminder Time</Label>
                        <Input
                          id="reminder-time"
                          type="time"
                          value={reminderTime}
                          onChange={(e) => setReminderTime(e.target.value)}
                        />
                      </div>
                    )}
                  </div>
                </div>
              )}

              {/* Data & Privacy Tab */}
              {activeTab === "data" && (
                <div className="space-y-6">
                  <div>
                    <h2 className="text-lg font-semibold mb-1">
                      Data & Privacy
                    </h2>
                    <p className="text-sm text-muted-foreground">
                      Manage your data and privacy settings.
                    </p>
                  </div>

                  <div className="space-y-4">
                    <div className="rounded-lg border p-4">
                      <div className="flex items-center justify-between mb-2">
                        <div>
                          <h3 className="font-medium">Export Data</h3>
                          <p className="text-sm text-muted-foreground">
                            Download all your data as JSON or CSV
                          </p>
                        </div>
                        <Button variant="outline" size="sm">
                          Export
                        </Button>
                      </div>
                    </div>

                    <div className="rounded-lg border border-destructive/50 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <div>
                          <h3 className="font-medium text-destructive">
                            Delete Account
                          </h3>
                          <p className="text-sm text-muted-foreground">
                            Permanently delete your account and all data
                          </p>
                        </div>
                        <Button variant="destructive" size="sm">
                          Delete
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Save Button */}
              <div className="mt-6 pt-6 border-t flex justify-end">
                <Button onClick={handleSave} disabled={isSaving}>
                  {isSaving ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    <>
                      <Save className="h-4 w-4 mr-2" />
                      Save changes
                    </>
                  )}
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}







