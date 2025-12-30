"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { CheckCircle2, ArrowRight, Target, TrendingUp, Calendar } from "lucide-react";
import Link from "next/link";
import { cn } from "@/lib/utils";

const steps = [
  {
    id: 1,
    title: "Welcome to Habitum",
    description:
      "Your personal habit tracking companion. Let's get you started on building better habits.",
    icon: Target,
  },
  {
    id: 2,
    title: "Track Your Progress",
    description:
      "Mark habits as complete each day and watch your streaks grow. Consistency is key!",
    icon: TrendingUp,
  },
  {
    id: 3,
    title: "Visualize Your Journey",
    description:
      "Use the calendar and insights pages to see your progress over time and identify patterns.",
    icon: Calendar,
  },
];

export default function OnboardingPage() {
  const router = useRouter();
  const [currentStep, setCurrentStep] = useState(0);
  const [completedSteps, setCompletedSteps] = useState<number[]>([]);

  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
      setCompletedSteps([...completedSteps, currentStep]);
    } else {
      // Complete onboarding
      router.push("/dashboard");
    }
  };

  const handleSkip = () => {
    router.push("/dashboard");
  };

  const currentStepData = steps[currentStep];

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-primary/5 via-background to-background">
      <div className="w-full max-w-2xl">
        {/* Progress Indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            {steps.map((step, index) => {
              const isCompleted = completedSteps.includes(index) || index < currentStep;
              const isCurrent = index === currentStep;
              const Icon = step.icon;

              return (
                <div key={step.id} className="flex items-center flex-1">
                  <div className="flex flex-col items-center flex-1">
                    <div
                      className={cn(
                        "w-12 h-12 rounded-full flex items-center justify-center border-2 transition-all",
                        isCompleted
                          ? "bg-primary border-primary text-primary-foreground"
                          : isCurrent
                          ? "border-primary bg-primary/10 text-primary"
                          : "border-muted bg-background text-muted-foreground"
                      )}
                    >
                      {isCompleted ? (
                        <CheckCircle2 className="h-6 w-6" />
                      ) : (
                        <Icon className="h-6 w-6" />
                      )}
                    </div>
                    <div className="mt-2 text-xs font-medium text-center max-w-[100px]">
                      {step.title}
                    </div>
                  </div>
                  {index < steps.length - 1 && (
                    <div
                      className={cn(
                        "h-0.5 flex-1 mx-2 transition-colors",
                        isCompleted ? "bg-primary" : "bg-muted"
                      )}
                    />
                  )}
                </div>
              );
            })}
          </div>
        </div>

        {/* Content Card */}
        <div className="rounded-xl border bg-background p-8 sm:p-12 shadow-lg">
          <div className="text-center space-y-6">
            <div className="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-primary/20 mb-4">
              {currentStepData && (
                <currentStepData.icon className="h-10 w-10 text-primary" />
              )}
            </div>

            <div className="space-y-4">
              <h1 className="text-3xl sm:text-4xl font-bold">
                {currentStepData?.title}
              </h1>
              <p className="text-lg text-muted-foreground max-w-md mx-auto">
                {currentStepData?.description}
              </p>
            </div>

            {/* Feature Highlights */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-8">
              <div className="rounded-lg border p-4">
                <Target className="h-6 w-6 text-primary mb-2" />
                <h3 className="font-medium mb-1">Create Habits</h3>
                <p className="text-sm text-muted-foreground">
                  Add habits you want to build
                </p>
              </div>
              <div className="rounded-lg border p-4">
                <TrendingUp className="h-6 w-6 text-primary mb-2" />
                <h3 className="font-medium mb-1">Track Progress</h3>
                <p className="text-sm text-muted-foreground">
                  Monitor your streaks and stats
                </p>
              </div>
              <div className="rounded-lg border p-4">
                <Calendar className="h-6 w-6 text-primary mb-2" />
                <h3 className="font-medium mb-1">Visualize Data</h3>
                <p className="text-sm text-muted-foreground">
                  See your journey over time
                </p>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center justify-between mt-12 pt-6 border-t">
            <Button variant="ghost" onClick={handleSkip}>
              Skip tutorial
            </Button>
            <Button onClick={handleNext} className="gap-2">
              {currentStep < steps.length - 1 ? (
                <>
                  Next
                  <ArrowRight className="h-4 w-4" />
                </>
              ) : (
                <>
                  Get Started
                  <ArrowRight className="h-4 w-4" />
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Tips */}
        <div className="mt-8 text-center">
          <p className="text-sm text-muted-foreground">
            Tip: You can always access this tutorial from the settings page
          </p>
        </div>
      </div>
    </div>
  );
}






