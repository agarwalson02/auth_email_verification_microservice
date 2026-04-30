package main

import (
	"auth_service/config"
	// "auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/internal/services"
	postgres "auth_service/pkg/db"
	emailclient "auth_service/pkg/email_client"
	"auth_service/pkg/logger"
	redis "auth_service/pkg/redis"
	"google.golang.org/grpc/reflection"

	// "context"
	"log"
	"net"

	grpcHandler "auth_service/internal/delivery/grpc"
	pb "auth_service/proto"

	"google.golang.org/grpc"
	// "github.com/joho/godotenv"
)

func main() {
	// if err:=godotenv.Load(); err!=nil{
	// 	log.Println("No .env file found")
	// }
	configPath := "config/config-local.yaml"

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	defer appLogger.Sync()
	appLogger.Info("Logger initialized")
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)
	appLogger.Infof("Successfully parsed config: %s", cfg.Server.AppVersion)

	//connecting to database
	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()
	appLogger.Info("Database connection established")

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}
	defer redisClient.Close()
	appLogger.Info("Redis connection established")

	//Test for db models
	repo := repository.NewUserRepository(db)

	// user := &models.User{
	// 	Email:     "test1@example.com",
	// 	FirstName: "Test",
	// 	LastName:  "User",
	// 	Password:  "hashedpassword",
	// 	Role:      "user",
	// 	Avatar:    "",
	// }

	// err = repo.CreateUser(user)
	// if err != nil {
	// 	log.Fatal("Error inserting user:", err)
	// }

	// log.Println("User created with ID:", user.ID)
	emailClient, _ := emailclient.NewEmailClient("localhost:5001")
	svc := services.NewAuthService(repo, redisClient,emailClient)

	// // user := &models.User{
	// // 	Email:     "new@example.com",
	// // 	FirstName: "New",
	// // 	LastName:  "User",
	// // 	Password:  "123456",
	// // }

	// // createdUser, err := svc.Register(user)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// ctx := context.Background()

	// user, sessionID, err := svc.Login(ctx, "new@example.com", "123456")
	// if err != nil {
	// 	log.Fatal("Error in login:", err)
	// }

	// log.Println("User ID:", user.ID)
	// log.Println("Session ID:", sessionID)
	// user, err = svc.GetMe(ctx, sessionID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("GetMe:", user.Email)

	// // logout
	// err = svc.Logout(ctx, sessionID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Logout successful")

	grpcServer := grpc.NewServer(
		
		grpc.UnaryInterceptor(grpcHandler.AuthInterceptor(appLogger)),
	)

	handler := grpcHandler.NewHandler(svc)

	pb.RegisterUserServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC server running on port", cfg.Server.Port)

	grpcServer.Serve(lis)

}
