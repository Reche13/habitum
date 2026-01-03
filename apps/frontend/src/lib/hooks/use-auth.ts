import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { authAPI } from "@/lib/api/auth";
import { useAuthStore } from "@/stores/auth-store";

// Query keys
export const authKeys = {
  all: ["auth"] as const,
};

// Login mutation
export function useLogin() {
  const router = useRouter();
  const { login } = useAuthStore();

  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      authAPI.login(email, password),
    onSuccess: (data) => {
      login(data.user, data.access_token, data.refresh_token);
      router.push("/dashboard");
    },
  });
}

// Signup mutation
export function useSignup() {
  const router = useRouter();
  const { login } = useAuthStore();

  return useMutation({
    mutationFn: ({
      name,
      email,
      password,
    }: {
      name: string;
      email: string;
      password: string;
    }) => authAPI.signup(name, email, password),
    onSuccess: (data) => {
      login(data.user, data.access_token, data.refresh_token);
      router.push("/dashboard");
    },
  });
}

// Google OAuth mutation
export function useGoogleAuth() {
  const router = useRouter();
  const { login } = useAuthStore();

  return useMutation({
    mutationFn: (token: string) => authAPI.googleAuth(token),
    onSuccess: (data) => {
      login(data.user, data.access_token, data.refresh_token);
      router.push("/dashboard");
    },
  });
}

// Verify email mutation
export function useVerifyEmail() {
  return useMutation({
    mutationFn: (token: string) => authAPI.verifyEmail(token),
  });
}

// Resend verification mutation
export function useResendVerification() {
  return useMutation({
    mutationFn: (email: string) => authAPI.resendVerification(email),
  });
}

// Forgot password mutation
export function useForgotPassword() {
  return useMutation({
    mutationFn: (email: string) => authAPI.forgotPassword(email),
  });
}

// Reset password mutation
export function useResetPassword() {
  const router = useRouter();
  return useMutation({
    mutationFn: ({
      token,
      newPassword,
    }: {
      token: string;
      newPassword: string;
    }) => authAPI.resetPassword(token, newPassword),
    onSuccess: () => {
      router.push("/login");
    },
  });
}

// Refresh token mutation
export function useRefreshToken() {
  const { setTokens } = useAuthStore();
  return useMutation({
    mutationFn: (refreshToken: string) =>
      authAPI.refreshToken(refreshToken),
    onSuccess: (data) => {
      setTokens(data.access_token, data.refresh_token);
    },
  });
}

// Test account login mutation
export function useTestAccountLogin() {
  const router = useRouter();
  const { login, setMockAccount } = useAuthStore();

  return useMutation({
    mutationFn: () => authAPI.testAccountLogin(),
    onSuccess: (data) => {
      login(data.user, data.access_token, data.refresh_token);
      setMockAccount(true);
      router.push("/dashboard");
    },
  });
}

// Logout mutation
export function useLogout() {
  const router = useRouter();
  const { logout } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authAPI.logout(),
    onSuccess: () => {
      logout();
      queryClient.clear();
      router.push("/login");
    },
  });
}

