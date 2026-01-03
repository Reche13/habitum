"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, Mail, Target, CheckCircle2 } from "lucide-react";
import Link from "next/link";
import { useForgotPassword } from "@/lib/hooks";
import { getErrorMessage } from "@/lib/api/client";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const forgotPassword = useForgotPassword();
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (!email) {
      setError("Please enter your email address");
      return;
    }

    try {
      await forgotPassword.mutateAsync(email);
      setSuccess(true);
    } catch (err: any) {
      setError(getErrorMessage(err));
    }
  };

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center p-8">
        <div className="max-w-md w-full space-y-6 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/20 mb-4">
            <CheckCircle2 className="h-8 w-8 text-green-600 dark:text-green-400" />
          </div>
          <h1 className="text-3xl font-semibold">Check your email</h1>
          <p className="text-muted-foreground">
            We've sent a password reset link to <strong>{email}</strong>. Please
            check your inbox and follow the instructions.
          </p>
          <Link href="/login">
            <Button variant="outline" className="w-full">
              Back to Login
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex">
      {/* Left side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary/10 via-primary/5 to-background items-center justify-center p-12">
        <div className="max-w-md text-center space-y-6">
          <div className="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-primary/20 mb-4">
            <Target className="h-10 w-10 text-primary" />
          </div>
          <h1 className="text-4xl font-bold">Habitum</h1>
          <p className="text-lg text-muted-foreground">
            Reset your password to regain access to your account.
          </p>
        </div>
      </div>

      {/* Right side - Form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8">
        <div className="w-full max-w-md space-y-8">
          <div className="lg:hidden text-center">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-xl bg-primary/20 mb-4">
              <Target className="h-8 w-8 text-primary" />
            </div>
            <h1 className="text-3xl font-bold">Habitum</h1>
          </div>

          <div className="space-y-2">
            <h2 className="text-3xl font-semibold">Forgot password?</h2>
            <p className="text-muted-foreground">
              Enter your email address and we'll send you a link to reset your
              password.
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
                {error}
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="pl-10"
                  disabled={forgotPassword.isPending}
                  required
                />
              </div>
            </div>

            <Button
              type="submit"
              className="w-full"
              disabled={forgotPassword.isPending}
            >
              {forgotPassword.isPending ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Sending...
                </>
              ) : (
                "Send Reset Link"
              )}
            </Button>
          </form>

          <div className="text-center text-sm text-muted-foreground">
            Remember your password?{" "}
            <Link href="/login" className="text-primary hover:underline font-medium">
              Sign in
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

