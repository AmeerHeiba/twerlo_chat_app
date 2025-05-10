import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/contexts/AuthContext";
import { User as UserIcon } from "lucide-react";
import { User } from "@/types";
import { toast } from "sonner";

interface ChatSidebarProps {
  contacts: User[];
  activeContactId: number | null;
  setActiveContactId: (id: number) => void;
  className?: string;
  isMobile?: boolean;
  onClose?: () => void;
}

const ChatSidebar = ({
  contacts = [],
  activeContactId,
  setActiveContactId,
  className = "",
  isMobile = false,
  onClose,
}: ChatSidebarProps) => {
  const { user, logout } = useAuth();
  const [searchQuery, setSearchQuery] = useState("");
  const navigate = useNavigate();

  // Safely filter contacts based on search query
  const filteredContacts = contacts
    .filter(contact => contact?.id) // Ensure contact has an id
    .filter((contact) => {
      const username = contact?.username ?? '';
      const query = searchQuery.toLowerCase();
      return username.toLowerCase().includes(query);
    });

  const handleContactClick = (contactId: number) => {
    setActiveContactId(contactId);
    if (isMobile && onClose) {
      onClose();
    }
  };

  const handleProfileClick = () => {
    navigate("/profile");
  };

  const getStatusColor = (status = "offline") => {
    switch (status) {
      case "online":
        return "bg-green-500";
      case "away":
        return "bg-yellow-500";
      default:
        return "bg-gray-400";
    }
  };

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  // Safe user name display
  const getDisplayName = (user?: { username?: string | null }) => {
    return user?.username?.trim() || 'Unknown';
  };

  // Safe avatar fallback
  const getAvatarFallback = (name?: string | null) => {
    return name?.charAt(0)?.toUpperCase() || '?';
  };

  return (
    <div className={`flex flex-col h-full border-r ${className}`}>
      {/* Header/Search area */}
      <div className="p-4 space-y-4 border-b">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-semibold">Chats</h2>
          <Button
            variant="ghost"
            size="sm"
            className="text-muted-foreground hover:text-foreground"
            onClick={handleProfileClick}
          >
            <UserIcon size={18} />
          </Button>
        </div>
        <Input
          placeholder="Search contacts..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="h-9"
        />
      </div>

      {/* Contacts List */}
      <ScrollArea className="flex-1 px-2 py-2">
        <div className="space-y-1">
          {filteredContacts.length > 0 ? (
            filteredContacts.map((contact) => {
              const displayName = getDisplayName(contact);
              return (
                <button
                  key={`contact-${contact.id}`} // Unique key with prefix
                  onClick={() => handleContactClick(contact.id)}
                  className={`w-full text-left p-3 rounded-md flex items-center gap-3 transition-colors ${
                    activeContactId === contact.id
                      ? "bg-accent text-accent-foreground"
                      : "hover:bg-muted"
                  }`}
                >
                  <div className="relative">
                    <Avatar className="h-10 w-10 border">
                      <AvatarFallback>
                        {getAvatarFallback(contact.username)}
                      </AvatarFallback>
                      <AvatarImage
                        src={`https://api.dicebear.com/7.x/initials/svg?seed=${displayName}`}
                      />
                    </Avatar>
                    <span
                      className={`absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-background ${getStatusColor(
                        contact.status
                      )}`}
                    ></span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-sm">
                      {displayName}
                    </div>
                    <div className="text-xs text-muted-foreground truncate">
                      {contact.status === "online"
                        ? "Online"
                        : contact.last_active
                        ? `Last seen ${new Date(contact.last_active).toLocaleTimeString()}`
                        : "Offline"}
                    </div>
                  </div>
                </button>
              );
            })
          ) : (
            <p className="text-center text-muted-foreground py-4">
              {contacts.length === 0 ? "No contacts available" : "No contacts found"}
            </p>
          )}
        </div>
      </ScrollArea>

      {/* User Profile Section */}
      <div className="p-4 mt-auto border-t">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Avatar className="h-10 w-10 border">
              <AvatarFallback>
                {getAvatarFallback(user?.username)}
              </AvatarFallback>
              <AvatarImage
                src={`https://api.dicebear.com/7.x/initials/svg?seed=${getDisplayName(user)}`}
              />
            </Avatar>
            <div>
              <div className="font-medium text-sm">
                {getDisplayName(user)}
              </div>
              <div className="text-xs text-muted-foreground">
                {user?.email || ''}
              </div>
            </div>
          </div>
          <Button variant="ghost" size="sm" onClick={handleLogout}>
            Logout
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ChatSidebar;