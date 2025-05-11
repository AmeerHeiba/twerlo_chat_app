import { Message } from "@/types";

type MessageHandler = (message: Message) => void;
type ConnectionChangeHandler = (connected: boolean) => void;

class WebSocketService {
  private socket: WebSocket | null = null;
  private messageHandlers: MessageHandler[] = [];
  private connectionHandlers: ConnectionChangeHandler[] = [];
  private reconnectTimeout: number | null = null;
  private url: string = "ws://localhost:8080/ws";

  constructor() {
    this.connect = this.connect.bind(this);
    this.disconnect = this.disconnect.bind(this);
    this.reconnect = this.reconnect.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
  }

  connect(token: string) {
    if (this.socket?.readyState === WebSocket.OPEN) return;

    try {
      this.socket = new WebSocket(`${this.url}?token=${token}`);

      this.socket.onopen = () => {
        this.notifyConnectionChange(true);
        if (this.reconnectTimeout) {
          clearTimeout(this.reconnectTimeout);
          this.reconnectTimeout = null;
        }
      };

      this.socket.onmessage = (event) => {
        try {
          const messageData = JSON.parse(event.data);
          if (
            messageData.id &&
            messageData.sender_id &&
            messageData.recipient_id
          ) {
            this.notifyMessageHandlers(messageData as Message);
          }
        } catch (error) {
          console.error("WebSocket message error:", error);
        }
      };

      this.socket.onclose = (event) => {
        this.notifyConnectionChange(false);
        this.reconnectTimeout = window.setTimeout(() => {
          this.reconnect(token);
        }, 3000);
      };

      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    } catch (error) {
      console.error("WebSocket connection error:", error);
      this.notifyConnectionChange(false);
    }
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }

    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
  }

  reconnect(token: string) {
    this.disconnect();
    this.connect(token);
  }

  sendMessage(message: string) {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(message);
      return true;
    }
    return false;
  }

  onMessage(handler: MessageHandler) {
    this.messageHandlers.push(handler);
    return () => {
      this.messageHandlers = this.messageHandlers.filter((h) => h !== handler);
    };
  }

  onConnectionChange(handler: ConnectionChangeHandler) {
    this.connectionHandlers.push(handler);
    return () => {
      this.connectionHandlers = this.connectionHandlers.filter(
        (h) => h !== handler
      );
    };
  }

  private notifyMessageHandlers(message: Message) {
    this.messageHandlers.forEach((handler) => {
      try {
        handler(message);
      } catch (error) {
        console.error("Message handler error:", error);
      }
    });
  }

  private notifyConnectionChange(connected: boolean) {
    this.connectionHandlers.forEach((handler) => handler(connected));
  }

  isConnected() {
    return this.socket?.readyState === WebSocket.OPEN;
  }
}

export const websocketService = new WebSocketService();
