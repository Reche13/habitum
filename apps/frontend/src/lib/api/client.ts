import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "@/stores/auth-store";
import { authAPI } from "./auth";

// API Base URL - adjust based on your backend port
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

// Create axios instance
export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 30000, // 30 seconds
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: any) => void;
  reject: (error?: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// Request interceptor - add auth token when available
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Get token from store (works in browser context)
    if (typeof window !== "undefined") {
      const { accessToken } = useAuthStore.getState();
      if (accessToken) {
        config.headers.Authorization = `Bearer ${accessToken}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor - handle errors globally and token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // Handle 401 - try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // If already refreshing, queue this request
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            }
            return apiClient(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const { refreshToken, logout } = useAuthStore.getState();
        if (!refreshToken) {
          throw new Error("No refresh token");
        }

        const data = await authAPI.refreshToken(refreshToken);
        const { setTokens } = useAuthStore.getState();
        setTokens(data.access_token, data.refresh_token);

        processQueue(null, data.access_token);

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        }

        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        // Redirect to login on refresh failure
        if (typeof window !== "undefined") {
          const { logout } = useAuthStore.getState();
          logout();
          window.location.href = "/login";
        }
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    // Handle common errors
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data as any;

      switch (status) {
        case 401:
          // Unauthorized - handled above
          break;
        case 403:
          // Forbidden
          break;
        case 404:
          // Not found
          break;
        case 422:
          // Validation error
          break;
        case 500:
          // Server error
          break;
      }

      // Return error with message
      return Promise.reject({
        message: data?.message || error.message,
        status,
        data: data,
      });
    }

    // Network error
    if (error.request) {
      return Promise.reject({
        message: "Network error. Please check your connection.",
        status: 0,
      });
    }

    return Promise.reject(error);
  }
);

// Helper to extract error message
export function getErrorMessage(error: any): string {
  if (error?.message) {
    return error.message;
  }
  if (error?.response?.data?.message) {
    return error.response.data.message;
  }
  return "An unexpected error occurred";
}


