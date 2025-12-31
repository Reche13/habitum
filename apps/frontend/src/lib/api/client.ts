import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";

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

// Request interceptor - add auth token when available
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // TODO: Add auth token when auth is implemented
    // const token = getAuthToken();
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor - handle errors globally
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    // Handle common errors
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data as any;

      switch (status) {
        case 401:
          // Unauthorized - redirect to login when auth is added
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


