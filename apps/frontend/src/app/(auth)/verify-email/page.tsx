"use client";

import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, CheckCircle2, XCircle, Mail } from "lucide-react";
import { useVerifyEmail, useResendVerification } from "@/lib/hooks";
import { getErrorMessage } from "@/lib/api/client";
import { useAuthStore } from "@/stores/auth-store";
import Link from "next/link";

export default function VerifyEmailPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const token = searchParams.get("token");
  const [email, setEmail] = useState("");
  const verifyEmail = useVerifyEmail();
  const resendVerification = useResendVerification();
  const { user } = useAuthStore();

  useEffect(() => {
    if (user?.email) {
      setEmail(user.email);
    }
  }, [user]);

  useEffect(() => {
    if (token && !verifyEmail.isSuccess && !verifyEmail.isPending) {
      verifyEmail.mutate(token);
    }
  }, [token]);

  const handleResend = async () => {
    if (!email) {
      return;
    }
    try {
      await resendVerification.mutateAsync(email);
    } catch (err: any) {
      // Error handled by mutation
    }
  };

  if (verifyEmail.isPending) {
    return (
      <div className="min-h-screen flex items-center justify-center p-8">
        <div className="max-w-md w-full space-y-4 text-center">
          <Loader2 className="h-12 w-12 animate-spin mx-auto text-primary" />
          <p className="text-muted-foreground">Verifying your email...</p>
        </div>
      </div>
    );
  }

  if (verifyEmail.isSuccess) {
    return (
      <div className="min-h-screen flex items-center justify-center p-8">
        <div className="max-w-md w-full space-y-6 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/20 mb-4">
            <CheckCircle2 className="h-8 w-8 text-green-600 dark:text-green-400" />
          </div>
          <h1 className="text-3xl font-semibold">Email Verified!</h1>
          <p className="text-muted-foreground">
            Your email has been successfully verified. You can now use all features
            of Habitum.
          </p>
          <Button onClick={() => router.push("/dashboard")} className="w-full">
            Go to Dashboard
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-8">
      <div className="max-w-md w-full space-y-6">
        <div className="text-center space-y-4">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-destructive/10 mb-4">
            <XCircle className="h-8 w-8 text-destructive" />
          </div>
          <h1 className="text-3xl font-semibold">Verification Failed</h1>
          <p className="text-muted-foreground">
            {verifyEmail.error
              ? getErrorMessage(verifyEmail.error)
              : "The verification link is invalid or has expired."}
          </p>
        </div>

        {email && (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="pl-10"
                  placeholder="your@email.com"
                />
              </div>
            </div>
            <Button
              onClick={handleResend}
              disabled={resendVerification.isPending}
              className="w-full"
              variant="outline"
            >
              {resendVerification.isPending ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Sending...
                </>
              ) : (
                "Resend Verification Email"
              )}
            </Button>
          </div>
        )}

        <div className="text-center">
          <Link href="/login" className="text-sm text-primary hover:underline">
            Back to Login
          </Link>
        </div>
      </div>
    </div>
  );
}

