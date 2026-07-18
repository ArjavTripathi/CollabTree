package main

import (
	"SocialMedia/internal/db"
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

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))

}
