import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthState {
  user: {
    id: string;
    name: string;
    email: string;
    email_verified: boolean;
    oauth_provider?: string;
  } | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  login: (
    user: {
      id: string;
      name: string;
      email: string;
      email_verified: boolean;
      oauth_provider?: string;
    },
    accessToken: string,
    refreshToken: string
  ) => void;
  setTokens: (accessToken: string, refreshToken: string) => void;
  logout: () => void;
  // Mock account flag
  isMockAccount: boolean;
  setMockAccount: (isMock: boolean) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isMockAccount: false,
      login: (user, accessToken, refreshToken) =>
        set({
          user,
          accessToken,
          refreshToken,
          isAuthenticated: true,
        }),
      setTokens: (accessToken, refreshToken) =>
        set({
          accessToken,
          refreshToken,
        }),
      logout: () =>
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
          isMockAccount: false,
        }),
      setMockAccount: (isMock) => set({ isMockAccount: isMock }),
    }),
    {
      name: "auth-storage",
    }
  )
);
