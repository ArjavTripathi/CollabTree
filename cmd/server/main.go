package main

import (
	"SocialMedia/internal/auth"
	"SocialMedia/internal/db"
	"SocialMedia/internal/middleware"
	"SocialMedia/internal/user"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	pool, err := db.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewSessionRepository(pool)
	authSvc := auth.NewService(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), authRepo, userService)
	authHandler := auth.NewHandler(authSvc)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)
	authHandler.RegisterRoutes(mux)
	handler := middleware.Recover(middleware.Logging(mux))

	log.Fatal(http.ListenAndServe(":8080", handler))

}
