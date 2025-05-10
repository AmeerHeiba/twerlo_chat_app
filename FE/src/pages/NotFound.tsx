
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";

const NotFound = () => {
  const navigate = useNavigate();
  
  const handleGoBack = () => {
    navigate(-1);
  };
  
  const handleGoHome = () => {
    navigate("/");
  };
  
  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-background to-accent/20 p-4">
      <div className="text-center max-w-md">
        <h1 className="text-6xl font-bold mb-6">404</h1>
        <h2 className="text-2xl font-semibold mb-4">Page Not Found</h2>
        <p className="text-muted-foreground mb-8">
          We couldn't find the page you were looking for. It might have been removed, renamed, or didn't exist in the first place.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Button onClick={handleGoBack} variant="outline">
            Go Back
          </Button>
          <Button onClick={handleGoHome}>
            Go Home
          </Button>
        </div>
      </div>
    </div>
  );
};

export default NotFound;
