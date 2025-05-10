export interface User {
  id: number;
  username: string;
  email: string;
  status: string;
  last_active?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  user_id: number;
  username: string;
  email: string;
}

export interface Message {
  id: number;
  content: string;
  type: "direct" | "broadcast";
  status: "sent" | "delivered" | "read";
  sender_id: number;
  recipient_id?: number;
  sent_at: string;
  delivered_at: string;
  read_at: string;
  media_url?: string;
}

export interface ConversationResponse {
  messages: Message[];
  total: number;
}

export interface ProfileUpdateRequest {
  username?: string;
  email?: string;
}

export interface PasswordUpdateRequest {
  current_password: string;
  new_password: string;
}

export interface ChatUser {
  id: number;
  username: string;
  status: string;
  lastActive: string;
}
