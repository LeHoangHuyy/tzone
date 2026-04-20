package server

import (
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/route"
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/LuuDinhTheTai/tzone/util/seed"
)

func (s *Server) MapHandlers() error {
	// Check if MongoDB is available
	if !s.HasMongoDB() {
		log.Println("⚠️ MongoDB is not available, brand and device routes will not work properly")
	}

	// AutoMigrate missing tables (Users, RefreshTokens, and the newly added RBAC tables)
	s.db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Favorite{},
		&model.Role{},
		&model.UserRole{},
		&model.Action{},
		&model.Resource{},
		&model.Permission{},
		&model.RolePermission{},
	)

	// Seed RBAC data
	seed.SeedAll(s.db)

	// Init repository
	userRepo := repository.NewUserRepository(s.db)
	tokenRepo := repository.NewRefreshTokenRepository(s.db)
	favoriteRepo := repository.NewFavoriteRepository(s.db)
	permissionRepo := repository.NewPermissionRepository(s.db)

	brandRepo := repository.NewBrandRepository()
	deviceRepo := repository.NewBrandRepository()
	if s.HasMongoDB() {
		brandRepo.SetClient(s.mongoClient)
		deviceRepo.SetClient(s.mongoClient)
		log.Printf("✅ MongoDB client set in repositories")
	}
	log.Printf("✅ Repositories initialized")

	// Init auth service
	authService := service.NewAuthService(userRepo, tokenRepo)

	// Init service
	brandService := service.NewBrandService(brandRepo)
	deviceService := service.NewDeviceService(deviceRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, deviceRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	log.Printf("✅ Services initialized")

	// Init handler
	commonHandler := handler.NewCommonHandler()
	frontendHandler := handler.NewFrontendHandler()
	brandHandler := handler.NewBrandHandler(brandService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)
	authHandler := handler.NewAuthHandler(authService)
	log.Printf("✅ Handlers initialized")

	// Init route
	route.MapCommonRoutes(s.r, commonHandler)
	route.MapFrontendRoutes(s.r, frontendHandler, permissionService)
	route.MapBrandRoutes(s.r, brandHandler, permissionService)
	route.MapDeviceRoutes(s.r, deviceHandler, permissionService)
	route.MapFavoriteRoutes(s.r, favoriteHandler)
	route.MapAuthRoutes(s.r, authHandler)
	log.Printf("✅ Routes initialized")

	return nil
}
