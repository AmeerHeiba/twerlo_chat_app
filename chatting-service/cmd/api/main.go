package main

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/config"
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/middleware"
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/routes"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/auth"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/database"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/realtime"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/storage"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func stackTraceHandler(c *fiber.Ctx, e interface{}) {
	shared.Log.Error("Recovered from panic",
		zap.Any("error", e),
		zap.String("path", c.Path()),
		zap.ByteString("stack", debug.Stack()),
	)
}

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize core dependencies
	db := initDB()
	app := fiber.New()

	shared.InitLogger(os.Getenv("APP_ENV"))
	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: stackTraceHandler,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8081",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
	}))

	app.Use(middleware.RequestContext())
	app.Use(middleware.ErrorHandler)

	// Initialize services and handlers
	deps := initDependencies(db)

	// Setup all routes
	routes.SetupRoutes(app, deps)

	// Not found route
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Endpoint not found",
			"path":    c.Path(),
		})
	})

	// Start server
	startServer(app)
}

func initDB() *gorm.DB {
	db, err := config.NewDBConnection(config.LoadDBConfig())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	return db
}

func initDependencies(db *gorm.DB) routes.Dependencies {
	// Initialize storage
	localStorage := storage.NewLocalStorage(
		os.Getenv("MEDIA_STORAGE_PATH"),
		os.Getenv("MEDIA_BASE_URL"),
	)
	mediaService := application.NewMediaService(localStorage)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	// Initialize WebSocket notifier (implements MessageNotifier interface)
	wsNotifier := realtime.NewWebSocketNotifier()

	// Repositories
	userRepo := database.NewUserRepository(db)
	messageRepo := database.NewMessageRepository(db)
	messageRecipientRepo := database.NewMessageRecipientRepository(db)

	// Services
	userService := application.NewUserService(userRepo)
	authCfg := config.LoadAuthConfig()
	jwtProvider := auth.NewJWTProvider(authCfg)
	authService := application.NewAuthService(userRepo, userService, jwtProvider)

	// Message service with WebSocket notifier
	messageService := application.NewMessageService(
		messageRepo,
		messageRecipientRepo,
		userRepo,
		wsNotifier,
		mediaService,
	)

	// WebSocket handler (for routes)
	wsHandler := handlers.NewWebSocketHandler(wsNotifier)

	return routes.Dependencies{
		DB:             db,
		UserHandler:    handlers.NewUserHandler(userService),
		AuthHandler:    handlers.NewAuthHandler(authService),
		MessageHandler: handlers.NewMessageHandler(messageService),
		MediaHandler:   mediaHandler,
		WSHandler:      wsHandler,
		JWTProvider:    jwtProvider,
	}
}

func startServer(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
