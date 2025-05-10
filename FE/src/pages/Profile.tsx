
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { TabsContent, Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuth } from "@/contexts/AuthContext";
import { api } from "@/services/api";
import { z } from "zod";
import { toast } from "sonner";

const profileSchema = z.object({
  username: z.string().min(3, "Username must be at least 3 characters"),
  email: z.string().email("Please enter a valid email"),
});

const passwordSchema = z.object({
  currentPassword: z.string().min(6, "Password must be at least 6 characters"),
  newPassword: z.string()
    .min(6, "Password must be at least 6 characters")
    .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
    .regex(/[a-z]/, "Password must contain at least one lowercase letter")
    .regex(/[0-9]/, "Password must contain at least one number")
    .regex(/[^a-zA-Z0-9]/, "Password must contain at least one special character"),
  confirmNewPassword: z.string().min(6, "Password must be at least 6 characters"),
}).refine((data) => data.newPassword === data.confirmNewPassword, {
  message: "Passwords don't match",
  path: ["confirmNewPassword"],
});

const Profile = () => {
  const { user, updateUser, logout } = useAuth();
  const navigate = useNavigate();
  
  const [profileForm, setProfileForm] = useState({
    username: user?.username || "",
    email: user?.email || "",
  });
  
  const [passwordForm, setPasswordForm] = useState({
    currentPassword: "",
    newPassword: "",
    confirmNewPassword: "",
  });
  
  const [profileErrors, setProfileErrors] = useState<Record<string, string>>({});
  const [passwordErrors, setPasswordErrors] = useState<Record<string, string>>({});
  
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false);
  const [isUpdatingPassword, setIsUpdatingPassword] = useState(false);

  const handleProfileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProfileForm((prev) => ({ ...prev, [name]: value }));
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordForm((prev) => ({ ...prev, [name]: value }));
  };

  const validateProfileForm = () => {
    try {
      profileSchema.parse(profileForm);
      setProfileErrors({});
      return true;
    } catch (error) {
      if (error instanceof z.ZodError) {
        const formattedErrors: Record<string, string> = {};
        error.errors.forEach((err) => {
          if (err.path[0]) {
            formattedErrors[err.path[0] as string] = err.message;
          }
        });
        setProfileErrors(formattedErrors);
      }
      return false;
    }
  };

  const validatePasswordForm = () => {
    try {
      passwordSchema.parse(passwordForm);
      setPasswordErrors({});
      return true;
    } catch (error) {
      if (error instanceof z.ZodError) {
        const formattedErrors: Record<string, string> = {};
        error.errors.forEach((err) => {
          if (err.path[0]) {
            formattedErrors[err.path[0] as string] = err.message;
          }
        });
        setPasswordErrors(formattedErrors);
      }
      return false;
    }
  };

  const handleProfileSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateProfileForm()) return;
    
    setIsUpdatingProfile(true);
    
    try {
      // Only update fields that have changed
      const updateData = {} as { username?: string; email?: string };
      
      if (profileForm.username !== user?.username) {
        updateData.username = profileForm.username;
      }
      
      if (profileForm.email !== user?.email) {
        updateData.email = profileForm.email;
      }
      
      // If nothing changed, don't make the API call
      if (Object.keys(updateData).length === 0) {
        toast.info("No changes to update");
        return;
      }
      
      const updatedUser = await api.users.updateProfile(updateData);
      updateUser(updatedUser);
      toast.success("Profile updated successfully");
    } catch (error) {
      console.error("Failed to update profile:", error);
      toast.error("Failed to update profile. Please try again.");
    } finally {
      setIsUpdatingProfile(false);
    }
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validatePasswordForm()) return;
    
    setIsUpdatingPassword(true);
    
    try {
      await api.users.updatePassword({
        current_password: passwordForm.currentPassword,
        new_password: passwordForm.newPassword,
      });
      
      setPasswordForm({
        currentPassword: "",
        newPassword: "",
        confirmNewPassword: "",
      });
      
      toast.success("Password updated successfully");
    } catch (error) {
      console.error("Failed to update password:", error);
      toast.error("Failed to update password. Please check your current password.");
    } finally {
      setIsUpdatingPassword(false);
    }
  };

  const handleBackToChat = () => {
    navigate("/");
  };

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-pulse flex flex-col items-center">
          <div className="h-12 w-12 rounded-full bg-primary/30 mb-4"></div>
          <div className="h-4 w-32 bg-muted rounded"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen py-8 px-4 bg-gradient-to-br from-background to-accent/20">
      <div className="max-w-2xl mx-auto animate-fade-in">
        <div className="mb-6 flex items-center justify-between">
          <Button variant="outline" onClick={handleBackToChat}>
            Back to Chat
          </Button>
          <Button variant="outline" onClick={handleLogout}>
            Logout
          </Button>
        </div>

        <div className="flex flex-col items-center mb-8">
          <Avatar className="h-24 w-24 border-4 border-background shadow-xl mb-4">
            <AvatarFallback className="text-3xl">
              {user.username.charAt(0).toUpperCase()}
            </AvatarFallback>
            <AvatarImage src={`https://api.dicebear.com/7.x/initials/svg?seed=${user.username}`} />
          </Avatar>
          <h1 className="text-2xl font-bold mb-1">{user.username}</h1>
          <p className="text-muted-foreground">{user.email}</p>
          <div className="flex items-center gap-2 mt-2">
            <span className={`h-2 w-2 rounded-full ${user.status === "online" ? "bg-green-500" : "bg-gray-400"}`}></span>
            <span className="text-sm text-muted-foreground capitalize">{user.status || "offline"}</span>
          </div>
        </div>

        <Tabs defaultValue="profile" className="w-full">
          <TabsList className="grid w-full grid-cols-2 mb-6">
            <TabsTrigger value="profile">Profile</TabsTrigger>
            <TabsTrigger value="security">Security</TabsTrigger>
          </TabsList>
          
          <TabsContent value="profile">
            <Card>
              <CardHeader>
                <CardTitle>Profile Information</CardTitle>
                <CardDescription>
                  Update your account details and personal information
                </CardDescription>
              </CardHeader>
              
              <form onSubmit={handleProfileSubmit}>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <label htmlFor="username" className="text-sm font-medium leading-none">
                      Username
                    </label>
                    <Input
                      id="username"
                      name="username"
                      value={profileForm.username}
                      onChange={handleProfileChange}
                      className={profileErrors.username ? "border-destructive" : ""}
                    />
                    {profileErrors.username && (
                      <p className="text-sm text-destructive">{profileErrors.username}</p>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    <label htmlFor="email" className="text-sm font-medium leading-none">
                      Email
                    </label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      value={profileForm.email}
                      onChange={handleProfileChange}
                      className={profileErrors.email ? "border-destructive" : ""}
                    />
                    {profileErrors.email && (
                      <p className="text-sm text-destructive">{profileErrors.email}</p>
                    )}
                  </div>
                </CardContent>
                
                <CardFooter>
                  <Button 
                    type="submit" 
                    className="ml-auto"
                    disabled={isUpdatingProfile}
                  >
                    {isUpdatingProfile ? (
                      <span className="flex items-center gap-1">
                        <span className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                        Saving...
                      </span>
                    ) : (
                      "Save Changes"
                    )}
                  </Button>
                </CardFooter>
              </form>
            </Card>
          </TabsContent>
          
          <TabsContent value="security">
            <Card>
              <CardHeader>
                <CardTitle>Change Password</CardTitle>
                <CardDescription>
                  Update your password to keep your account secure
                </CardDescription>
              </CardHeader>
              
              <form onSubmit={handlePasswordSubmit}>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <label htmlFor="currentPassword" className="text-sm font-medium leading-none">
                      Current Password
                    </label>
                    <Input
                      id="currentPassword"
                      name="currentPassword"
                      type="password"
                      value={passwordForm.currentPassword}
                      onChange={handlePasswordChange}
                      className={passwordErrors.currentPassword ? "border-destructive" : ""}
                    />
                    {passwordErrors.currentPassword && (
                      <p className="text-sm text-destructive">{passwordErrors.currentPassword}</p>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    <label htmlFor="newPassword" className="text-sm font-medium leading-none">
                      New Password
                    </label>
                    <Input
                      id="newPassword"
                      name="newPassword"
                      type="password"
                      value={passwordForm.newPassword}
                      onChange={handlePasswordChange}
                      className={passwordErrors.newPassword ? "border-destructive" : ""}
                    />
                    {passwordErrors.newPassword && (
                      <p className="text-sm text-destructive">{passwordErrors.newPassword}</p>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    <label htmlFor="confirmNewPassword" className="text-sm font-medium leading-none">
                      Confirm New Password
                    </label>
                    <Input
                      id="confirmNewPassword"
                      name="confirmNewPassword"
                      type="password"
                      value={passwordForm.confirmNewPassword}
                      onChange={handlePasswordChange}
                      className={passwordErrors.confirmNewPassword ? "border-destructive" : ""}
                    />
                    {passwordErrors.confirmNewPassword && (
                      <p className="text-sm text-destructive">{passwordErrors.confirmNewPassword}</p>
                    )}
                  </div>
                </CardContent>
                
                <CardFooter>
                  <Button 
                    type="submit" 
                    className="ml-auto"
                    disabled={isUpdatingPassword}
                  >
                    {isUpdatingPassword ? (
                      <span className="flex items-center gap-1">
                        <span className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                        Updating...
                      </span>
                    ) : (
                      "Update Password"
                    )}
                  </Button>
                </CardFooter>
              </form>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default Profile;
