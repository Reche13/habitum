import { apiClient } from "./client";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  name: string;
  email: string;
  password: string;
}

export interface GoogleAuthRequest {
  token: string;
}

export interface VerifyEmailRequest {
  token: string;
}

export interface ResendVerificationRequest {
  email: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface UserResponse {
  id: string;
  name: string;
  email: string;
  email_verified: boolean;
  oauth_provider?: string;
}

export interface AuthResponse {
  user: UserResponse;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// Auth API
export const authAPI = {
  // Login with email/password
  login: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/login", {
      email,
      password,
    });
    return response.data;
  },

  // Signup with email/password
  signup: async (
    name: string,
    email: string,
    password: string
  ): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/signup", {
      name,
      email,
      password,
    });
    return response.data;
  },

  // Google OAuth
  googleAuth: async (token: string): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/google", {
      token,
    });
    return response.data;
  },

  // Verify email
  verifyEmail: async (token: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      "/auth/verify-email",
      { token }
    );
    return response.data;
  },

  // Resend verification email
  resendVerification: async (
    email: string
  ): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      "/auth/resend-verification",
      { email }
    );
    return response.data;
  },

  // Forgot password
  forgotPassword: async (email: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      "/auth/forgot-password",
      { email }
    );
    return response.data;
  },

  // Reset password
  resetPassword: async (
    token: string,
    newPassword: string
  ): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      "/auth/reset-password",
      { token, new_password: newPassword }
    );
    return response.data;
  },

  // Refresh token
  refreshToken: async (
    refreshToken: string
  ): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/refresh", {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  // Test account login
  testAccountLogin: async (): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/test-account");
    return response.data;
  },

  // Logout
  logout: async (): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>("/auth/logout");
    return response.data;
  },
};

