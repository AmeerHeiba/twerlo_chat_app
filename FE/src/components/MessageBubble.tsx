
import { useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Message } from "@/types";
import { formatDistanceToNow } from "date-fns";
import { useAuth } from "@/contexts/AuthContext";

interface MessageBubbleProps {
  message: Message;
  senderName: string;
}

const MessageBubble = ({ message, senderName }: MessageBubbleProps) => {
  const { user } = useAuth();
  const [imageError, setImageError] = useState(false);
  
  const isCurrentUser = user?.id === message.sender_id;
  const formattedTime = formatTime(message.sent_at);
  const messageTimeAgo = formatDistanceToNow(new Date(message.sent_at), { addSuffix: true });
  
  // Format the message timestamp
  function formatTime(dateString: string) {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
  
  // Status indicator component
  const StatusIndicator = ({ status }: { status: string }) => {
    if (!isCurrentUser) return null;
    
    let statusText = "";
    
    switch (status) {
      case "sent":
        statusText = "Sent";
        break;
      case "delivered":
        statusText = "Delivered";
        break;
      case "read":
        statusText = "Read";
        break;
      default:
        statusText = "Sent";
    }
    
    return (
      <span className="text-xs text-muted-foreground">{statusText}</span>
    );
  };

  return (
    <div
      className={`flex gap-3 max-w-[80%] ${
        isCurrentUser ? "self-end flex-row-reverse" : "self-start"
      } ${isCurrentUser ? "message-out" : "message-in"}`}
    >
      {/* Avatar (only show for received messages) */}
      {!isCurrentUser && (
        <div className="flex-shrink-0">
          <Avatar className="h-8 w-8">
            <AvatarFallback>{senderName.charAt(0).toUpperCase()}</AvatarFallback>
            <AvatarImage src={`https://api.dicebear.com/7.x/initials/svg?seed=${senderName}`} />
          </Avatar>
        </div>
      )}
      
      {/* Message content */}
      <div className="flex flex-col gap-1">
        <div
          className={`rounded-2xl px-4 py-2 shadow-sm ${
            isCurrentUser
              ? "bg-primary text-primary-foreground"
              : "bg-muted"
          }`}
        >
          {/* Message text */}
          <p className="text-sm whitespace-pre-wrap">{message.content}</p>
          
          {/* Media attachment if any */}
          {message.media_url && !imageError && (
            <div className="mt-2 rounded-md overflow-hidden">
              <img
                src={message.media_url}
                alt="Attachment"
                className="max-w-full h-auto object-contain"
                onError={() => setImageError(true)}
              />
            </div>
          )}
        </div>
        
        {/* Timestamp and status */}
        <div
          className={`flex text-xs text-muted-foreground ${
            isCurrentUser ? "justify-end" : ""
          }`}
        >
          <span 
            className="text-xs text-muted-foreground" 
            title={messageTimeAgo}
          >
            {formattedTime}
          </span>
          <span className="mx-1">•</span>
          <StatusIndicator status={message.status} />
        </div>
      </div>
    </div>
  );
};

export default MessageBubble;
