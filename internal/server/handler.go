package server

import (
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/route"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	// "github.com/supabase-community/postgrest-go"
)

func (s *Server) MapHandlers() error {
	// Init repository
	mongoDBRepo := repository.NewMongoDbRepository()
	postgreRepo := repository.NewPostgreRepository(s.db)
	userRepo := repository.NewUserRepository(postgreRepo.DB)
	tokenRepo := repository.NewRefreshTokenRepository(postgreRepo.DB)
	log.Printf("✅ Repositories initialized")

	// Init service
	brandService := service.NewBrandService(mongoDBRepo)
	deviceService := service.NewDeviceService(mongoDBRepo)
	log.Printf("✅ Services initialized")

	// Init handler
	commonHandler := handler.NewCommonHandler()
	brandHandler := handler.NewBrandHandler(brandService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	log.Printf("✅ Handlers initialized")

	// Init route
	route.MapCommonRoutes(s.r, commonHandler)
	route.MapBrandRoutes(s.r, brandHandler)
	route.MapDeviceRoutes(s.r, deviceHandler)
	log.Printf("✅ Routes initialized")
	//
	authService := service.NewAuthService(userRepo, tokenRepo)
	authHandler := handler.NewAuthHandler(authService)
	route.MapAuthRoutes(s.r, authHandler)

	return nil
}
