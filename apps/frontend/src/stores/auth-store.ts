import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthState {
  user: {
    id: string;
    name: string;
    email: string;
  } | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: { id: string; name: string; email: string }, token: string) => void;
  logout: () => void;
  // Mock account flag
  isMockAccount: boolean;
  setMockAccount: (isMock: boolean) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isMockAccount: false,
      login: (user, token) =>
        set({
          user,
          token,
          isAuthenticated: true,
        }),
      logout: () =>
        set({
          user: null,
          token: null,
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

