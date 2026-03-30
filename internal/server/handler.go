package server

import (
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/route"
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"github.com/LuuDinhTheTai/tzone/internal/service"
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
		&model.Role{},
		&model.UserRole{},
		&model.Action{},
		&model.Resource{},
		&model.Permission{},
		&model.RolePermission{},
	)

	// Init repository
	userRepo := repository.NewUserRepository(s.db)
	tokenRepo := repository.NewRefreshTokenRepository(s.db)
	permissionRepo := repository.NewPermissionRepository(s.db)

	brandRepo := repository.NewBrandRepository()
	deviceRepo := repository.NewBrandRepository()
	if s.HasMongoDB() {
		brandRepo.SetClient(s.mongoClient)
		deviceRepo.SetClient(s.mongoClient)
		log.Printf("✅ MongoDB client set in repositories")
	}
	log.Printf("✅ Repositories initialized")

	// Init service
	brandService := service.NewBrandService(brandRepo)
	deviceService := service.NewDeviceService(deviceRepo)
	permissionService := service.NewPermissionService(permissionRepo)
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

	// Init auth service and routes
	authService := service.NewAuthService(userRepo, tokenRepo)
	authHandler := handler.NewAuthHandler(authService)
	route.MapAuthRoutes(s.r, authHandler)

	// Keep compiler happy since permissionService is meant to be used in protected routes
	_ = permissionService

	return nil
}
