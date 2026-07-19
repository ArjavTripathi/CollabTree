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
	"time"
)

func main() {
	ctx := context.Background()
	pool, err := db.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	rateLimit := make(map[string]middleware.RequestTracker)
	timeLimit := make(map[time.Time]string)

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewSessionRepository(pool)
	authSvc := auth.NewService(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), authRepo, userService)
	authHandler := auth.NewHandler(authSvc)

	authSessionStore := auth.NewSessionRepository(pool)
	authorizeRequest := middleware.NewAuthMiddleware(authSessionStore)

	CORSrequest := middleware.CORS(os.Getenv("FRONTEND_ORIGIN"))
	RateLimit := middleware.RateLimit(rateLimit, timeLimit)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux, authorizeRequest)
	authHandler.RegisterRoutes(mux)
	handler := middleware.Recover(middleware.Logging(middleware.RequestId(CORSrequest(RateLimit(mux)))))

	log.Fatal(http.ListenAndServe(":8080", handler))

}
