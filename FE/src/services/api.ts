
import { AuthResponse, ConversationResponse, Message, PasswordUpdateRequest, ProfileUpdateRequest, User } from "@/types";

const API_URL = "http://127.0.0.1:8080/api";

// Helper function to handle HTTP errors
const handleResponse = async (response: Response) => {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || `API Error: ${response.status}`);
  }
  return response.json();
};

// Helper to get auth headers
const getAuthHeaders = () => {
  const token = localStorage.getItem("access_token");
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
};

export const api = {
  auth: {
    login: async (username: string, password: string): Promise<AuthResponse> => {
      const response = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });
      return handleResponse(response);
    },

    register: async (email: string, username: string, password: string): Promise<AuthResponse> => {
      const response = await fetch(`${API_URL}/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, username, password }),
      });
      return handleResponse(response);
    },

    logout: () => {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("user_id");
      localStorage.removeItem("username");
    },
  },

  users: {
    getProfile: async (): Promise<User> => {
      const response = await fetch(`${API_URL}/users/profile`, {
        headers: getAuthHeaders(),
      });
      return handleResponse(response);
    },

    updateProfile: async (data: ProfileUpdateRequest): Promise<User> => {
      const response = await fetch(`${API_URL}/users/profile`, {
        method: "PUT",
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
      });
      return handleResponse(response);
    },

getAll: async (): Promise<User[]> => {
  const response = await fetch(`${API_URL}/users/all`, {
    headers: getAuthHeaders(),
  });
  const data = await handleResponse(response);
  
  return data.map(user => ({
    id: user.ID,
    username: user.Username,
    email: user.Email,
    status: user.Status,
    last_active: user.LastActiveAt
  }));
},

    updatePassword: async (data: PasswordUpdateRequest): Promise<void> => {
      const response = await fetch(`${API_URL}/users/profile`, {
        method: "PUT",
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
      });
      return handleResponse(response);
    },
  },

  messages: {
    getConversation: async (userId: number): Promise<ConversationResponse> => {
      const response = await fetch(`${API_URL}/messages/conversation/${userId}`, {
        headers: getAuthHeaders(),
      });
      return handleResponse(response);
    },

    sendMessage: async (
      content: string, 
      recipientId: number, 
      mediaUrl?: string
    ): Promise<Message> => {
      const payload: any = {
        content,
        recipient_id: recipientId,
        type: "direct",
      };
      
      if (mediaUrl) {
        payload.media_url = mediaUrl;
      }
      
      const response = await fetch(`${API_URL}/messages`, {
        method: "POST",
        headers: getAuthHeaders(),
        body: JSON.stringify(payload),
      });
      return handleResponse(response);
    },

    broadcastMessage: async (
      content: string, 
      recipientIds: number[]
    ): Promise<Message> => {
      const response = await fetch(`${API_URL}/messages/broadcast`, {
        method: "POST",
        headers: getAuthHeaders(),
        body: JSON.stringify({
          content,
          recipient_ids: recipientIds,
          type: "broadcast",
        }),
      });
      return handleResponse(response);
    },
  },

  media: {
    uploadFile: async (file: File): Promise<{ url: string }> => {
      const formData = new FormData();
      formData.append("file", file);

      const token = localStorage.getItem("access_token");
      
      const response = await fetch(`${API_URL}/media/upload`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          // Content-Type is automatically set when using FormData
        },
        body: formData,
      });
      
      return handleResponse(response);
    },
  },
};
