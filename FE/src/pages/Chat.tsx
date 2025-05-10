import { useState, useEffect, useRef } from "react";
import { useParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import ChatSidebar from "@/components/ChatSidebar";
import MessageBubble from "@/components/MessageBubble";
import ChatInput from "@/components/ChatInput";
import { useAuth } from "@/contexts/AuthContext";
import { api } from "@/services/api";
import { websocketService } from "@/services/websocket";
import { Message, User } from "@/types";
import { Menu } from "lucide-react";
import { useIsMobile } from "@/hooks/use-mobile";
import { toast } from "sonner";



const Chat = () => {
  const [isLoadingContacts, setIsLoadingContacts] = useState(false);
  const { id } = useParams<{ id: string }>();
  const { user, isAuthenticated } = useAuth();
  const isMobile = useIsMobile();

  const [activeContactId, setActiveContactId] = useState<number | null>(
    id ? parseInt(id, 10) : null
  );
  const [showSidebar, setShowSidebar] = useState(!isMobile);
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [contacts, setContacts] = useState<User[]>([]);
  const [wsConnected, setWsConnected] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);

  const activeContact = contacts.find(
    (contact) => contact.id === activeContactId
  );

useEffect(() => {
  const fetchContacts = async () => {


    if (!isAuthenticated) return;

      setIsLoadingContacts(true);
    
    try {
      const users = await api.users.getAll();
      setContacts(users);
      
      // Set first contact as active if none is selected
      if (users.length > 0 && !activeContactId) {
        setActiveContactId(users[0].id);
      } else if (id) {
        // Make sure the contact from URL params exists
        const contactExists = users.some(user => user.id === parseInt(id, 10));
        if (!contactExists && users.length > 0) {
          setActiveContactId(users[0].id);
        }
      }
    } catch (error) {
      console.error("Failed to fetch contacts:", error);
      toast.error("Failed to load contacts");
    }
  };
  fetchContacts();
}, [isAuthenticated, activeContactId, id]);

  useEffect(() => {
    // Handle WebSocket connection status
    const unsubscribe = websocketService.onConnectionChange((connected) => {
      setWsConnected(connected);
      if (connected) {
        toast.success("Connected to chat server");
      } else {
        toast.error("Disconnected from chat server");
      }
    });

    return () => unsubscribe();
  }, []);

  // Handle WebSocket messages
  useEffect(() => {
    const unsubscribe = websocketService.onMessage((newMessage) => {
      if (
        (newMessage.sender_id === activeContactId &&
          newMessage.recipient_id === user?.id) ||
        (newMessage.sender_id === user?.id &&
          newMessage.recipient_id === activeContactId)
      ) {
        setMessages((prevMessages) => [...prevMessages, newMessage]);
      }
    });

    return () => unsubscribe();
  }, [activeContactId, user?.id]);

  // Fetch conversation when contact changes
  useEffect(() => {
    if (activeContactId && user) {
      fetchMessages();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeContactId, user]);

    useEffect(() => {
    const unsubscribe = websocketService.onMessage((message: Message) => {
      // Only add the message if it's from the active contact or to the current user
      if (
        (message.sender_id === activeContactId && message.recipient_id === user?.id) ||
        (message.recipient_id === activeContactId && message.sender_id === user?.id)
      ) {
        setMessages((prevMessages) => [...prevMessages, message]);
      }
    });

    return () => unsubscribe();
  }, [activeContactId, user?.id]);

    useEffect(() => {
    if (isAuthenticated) {
      const token = localStorage.getItem("access_token");
      if (token) {
        websocketService.connect(token);
      }
    }

    return () => {
      websocketService.disconnect();
    };
  }, [isAuthenticated]);

    useEffect(() => {
    const fetchConversation = async () => {
      if (activeContactId) {
        setIsLoading(true);
        try {
          const convo = await api.messages.getConversation(activeContactId);
          setMessages(convo.messages);
        } catch (error) {
          console.error("Failed to fetch conversation:", error);
          toast.error("Failed to load messages");
        } finally {
          setIsLoading(false);
        }
      }
    };

    fetchConversation();
  }, [activeContactId]);


  // Auto-scroll to bottom when messages change
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Update sidebar visibility based on screen size
  useEffect(() => {
    setShowSidebar(!isMobile);
  }, [isMobile]);

  const fetchMessages = async () => {
    if (!activeContactId) return;

    setIsLoading(true);
    try {
      const response = await api.messages.getConversation(activeContactId);
      // Sort messages by sent_at in ascending order
      const sortedMessages = [...response.messages].sort(
        (a, b) => new Date(a.sent_at).getTime() - new Date(b.sent_at).getTime()
      );
      setMessages(sortedMessages);
    } catch (error) {
      console.error("Failed to fetch messages:", error);
      toast.error("Failed to load conversation");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSendMessage = async (content: string, mediaUrl?: string) => {
    if (!activeContactId || (!content.trim() && !mediaUrl)) return;

    try {
      const newMessage = await api.messages.sendMessage(
        content,
        activeContactId,
        mediaUrl
      );
      setMessages((prevMessages) => [...prevMessages, newMessage]);
    } catch (error) {
      console.error("Failed to send message:", error);
      throw error;
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const toggleSidebar = () => {
    setShowSidebar((prev) => !prev);
  };

  const getSenderName = (senderId: number) => {
    if (senderId === user?.id) return "You";

    const contact = contacts.find((c) => c.id === senderId);
    return contact ? contact.username : "Unknown";
  };

  return (
    <div className="flex h-screen bg-background">
      {/* Mobile Sidebar */}
      {isMobile && showSidebar && (
        <div className="absolute inset-0 z-50 bg-background">
          <ChatSidebar
            contacts={contacts}
            activeContactId={activeContactId}
            setActiveContactId={setActiveContactId}
            isMobile={true}
            onClose={() => setShowSidebar(false)}
          />
        </div>
      )}

      {/* Desktop Sidebar */}
      {!isMobile && showSidebar && (
        <div className="w-80 h-full hidden md:block">
          <ChatSidebar
            contacts={contacts}
            activeContactId={activeContactId}
            setActiveContactId={setActiveContactId}
          />
        </div>
      )}

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col h-full">
        {/* Chat Header */}
        <header className="flex items-center gap-3 px-4 py-2 h-16 border-b">
          <Button
            variant="ghost"
            size="icon"
            className="md:hidden"
            onClick={toggleSidebar}
          >
            <Menu size={20} />
          </Button>

{isLoadingContacts ? (
  <div className="flex items-center gap-3">
    <div className="h-9 w-9 rounded-full bg-muted animate-pulse"></div>
    <div className="space-y-1">
      <div className="h-4 w-24 bg-muted rounded animate-pulse"></div>
      <div className="h-3 w-16 bg-muted rounded animate-pulse"></div>
    </div>
  </div>
) : activeContact ? (
  <div className="flex items-center gap-3">
    <div className="relative">
      <div className="h-9 w-9 rounded-full bg-muted flex items-center justify-center">
        {activeContact.username?.charAt(0).toUpperCase() ?? "?"}
      </div>
      <span
        className={`absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-background ${
          activeContact.status === "online"
            ? "bg-green-500"
            : activeContact.status === "away"
            ? "bg-yellow-500"
            : "bg-gray-400"
        }`}
      ></span>
    </div>
    <div>
      <h2 className="font-medium">{activeContact.username ?? "Unknown"}</h2>
      <p className="text-xs text-muted-foreground">
        {activeContact.status === "online"
          ? "Online"
          : activeContact.last_active
          ? `Last seen ${new Date(activeContact.last_active).toLocaleTimeString()}`
          : "Offline"}
      </p>
    </div>
  </div>
) : (
  <div>
    <h2 className="font-medium">Select a contact</h2>
  </div>
)}
        </header>

        {/* Messages Area */}
        <ScrollArea className="flex-1 p-4">
          {isLoading ? (
            <div className="h-full flex items-center justify-center">
              <div className="animate-pulse space-y-2">
                <div className="h-12 w-48 bg-muted rounded"></div>
                <div className="h-12 w-32 bg-muted rounded"></div>
              </div>
            </div>
          ) : activeContactId ? (
            messages.length > 0 ? (
              <div className="space-y-4">
                {messages.map((message) => (
                  <MessageBubble
                    key={message.id}
                    message={message}
                    senderName={getSenderName(message.sender_id)}
                  />
                ))}
                <div ref={messagesEndRef} />
              </div>
            ) : (
              <div className="h-full flex items-center justify-center">
                <div className="text-center space-y-2">
                  <p className="text-lg font-medium">No messages yet</p>
                  <p className="text-muted-foreground">
                    Send a message to start the conversation
                  </p>
                </div>
              </div>
            )
          ) : (
            <div className="h-full flex items-center justify-center">
              <div className="text-center space-y-2">
                <p className="text-lg font-medium">Select a contact</p>
                <p className="text-muted-foreground">
                  Choose someone to start chatting with
                </p>
              </div>
            </div>
          )}
        </ScrollArea>

        {/* Chat Input */}
        <ChatInput
          onSendMessage={handleSendMessage}
          recipientId={activeContactId || 0}
          disabled={!activeContactId || !wsConnected}
        />
      </div>
    </div>
  );
};

export default Chat;