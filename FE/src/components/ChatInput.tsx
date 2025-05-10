import { useState, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Send, Image } from "lucide-react";
import { toast } from "sonner";
import { api } from "@/services/api";

interface ChatInputProps {
  onSendMessage: (content: string, mediaUrl?: string) => void;
  recipientId: number;
  disabled?: boolean;
}

const ChatInput = ({
  onSendMessage,
  recipientId,
  disabled = false,
}: ChatInputProps) => {
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [uploadedFileUrl, setUploadedFileUrl] = useState<string | null>(null);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!message.trim() && !uploadedFileUrl) return;

    setIsSubmitting(true);

    try {
      await onSendMessage(message.trim(), uploadedFileUrl || undefined);
      setMessage("");
      setPreviewUrl(null);
      setUploadedFileUrl(null);
    } catch (error) {
      console.error("Failed to send message:", error);
      toast.error("Failed to send message. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    // Send message on Enter without Shift
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const maxSizeInBytes = 10 * 1024 * 1024; // 10MB
    if (file.size > maxSizeInBytes) {
      toast.error("File size exceeds 10MB limit");
      return;
    }

    // Check if file is an image
    if (!file.type.startsWith("image/")) {
      toast.error("Only image files are supported");
      return;
    }

    // Show preview
    const objectUrl = URL.createObjectURL(file);
    setPreviewUrl(objectUrl);

    try {
      setIsUploading(true);

      // Upload file to server
      const response = await api.media.uploadFile(file);
      setUploadedFileUrl(response.url);

      toast.success("Image uploaded successfully");
    } catch (error) {
      console.error("Failed to upload file:", error);
      toast.error("Failed to upload image. Please try again.");
      setPreviewUrl(null);
    } finally {
      setIsUploading(false);
    }
  };

  const handleRemoveMedia = () => {
    setPreviewUrl(null);
    setUploadedFileUrl(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-4 border-t">
      {/* Media Preview */}
      {previewUrl && (
        <div className="mb-3 relative">
          <div className="relative rounded-md overflow-hidden border">
            <img
              src={previewUrl}
              alt="Upload preview"
              className="max-h-40 w-auto object-contain mx-auto"
            />
            <button
              type="button"
              onClick={handleRemoveMedia}
              className="absolute top-1 right-1 bg-black/50 text-white rounded-full p-1 hover:bg-black/70"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
          {isUploading && (
            <div className="absolute inset-0 bg-black/20 flex items-center justify-center">
              <div className="h-6 w-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
            </div>
          )}
        </div>
      )}

      <div className="flex items-end gap-2">
        {/* File Upload Button */}
        <div>
          <input
            type="file"
            ref={fileInputRef}
            onChange={handleFileChange}
            accept="image/*"
            className="hidden"
            id="file-upload"
          />
          <Button
            type="button"
            variant="outline"
            size="icon"
            className="rounded-full"
            onClick={() => fileInputRef.current?.click()}
            // disabled={isUploading || disabled}
          >
            <Image size={18} />
          </Button>
        </div>

        {/* Message Input */}
        <Textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type a message..."
          className="min-h-[20px] resize-none flex-1"
          rows={1}
          // disabled={isSubmitting || disabled}
        />

        {/* Send Button */}
        <Button
          type="submit"
          size="icon"
          className="rounded-full"
          // disabled={
          //   (!message.trim() && !uploadedFileUrl) || isSubmitting || disabled
          // }
        >
          {isSubmitting ? (
            <div className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
          ) : (
            <Send size={18} />
          )}
        </Button>
      </div>
    </form>
  );
};

export default ChatInput;
