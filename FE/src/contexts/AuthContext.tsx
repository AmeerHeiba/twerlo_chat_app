
import { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { User } from "@/types";
import { api } from "@/services/api";
import { toast } from "sonner";
import { websocketService } from "@/services/websocket";

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  register: (email: string, username: string, password: string) => Promise<boolean>;
  logout: () => void;
  updateUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
  login: async () => false,
  register: async () => false,
  logout: () => {},
  updateUser: () => {},
});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadUser = async () => {
      const token = localStorage.getItem("access_token");
      if (token) {
        try {
          const userData = await api.users.getProfile();
          setUser(userData);
          
          // Connect to WebSocket if we have a user and token
          websocketService.connect(token);
        } catch (error) {
          console.error("Failed to load user:", error);
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
        }
      }
      setIsLoading(false);
    };

    loadUser();

    return () => {
      // Disconnect WebSocket on unmount
      websocketService.disconnect();
    };
  }, []);

  const login = async (username: string, password: string) => {
    try {
      const response = await api.auth.login(username, password);
      
      // Save tokens and user info
      localStorage.setItem("access_token", response.access_token);
      localStorage.setItem("refresh_token", response.refresh_token);
      localStorage.setItem("user_id", response.user_id.toString());
      localStorage.setItem("username", response.username);
      
      // Get full user profile
      const userData = await api.users.getProfile();
      setUser(userData);
      
      // Connect to WebSocket
      websocketService.connect(response.access_token);
      
      toast.success(`Welcome back, ${response.username}!`);
      return true;
    } catch (error) {
      console.error("Login failed:", error);
      toast.error("Login failed: Invalid credentials");
      return false;
    }
  };

  const register = async (email: string, username: string, password: string) => {
    try {
      const response = await api.auth.register(email, username, password);
      
      // Save tokens and user info
      localStorage.setItem("access_token", response.access_token);
      localStorage.setItem("refresh_token", response.refresh_token);
      localStorage.setItem("user_id", response.user_id.toString());
      localStorage.setItem("username", response.username);
      
      // Set user directly from registration response
      setUser({
        id: response.user_id,
        username: response.username,
        email: response.email,
      });
      
      // Connect to WebSocket
      websocketService.connect(response.access_token);
      
      toast.success(`Welcome, ${response.username}!`);
      return true;
    } catch (error) {
      console.error("Registration failed:", error);
      toast.error("Registration failed. Please try again.");
      return false;
    }
  };

  const logout = () => {
    api.auth.logout();
    websocketService.disconnect();
    setUser(null);
    toast.info("You've been logged out");
  };

  const updateUser = (updatedUser: User) => {
    setUser(updatedUser);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
        updateUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
