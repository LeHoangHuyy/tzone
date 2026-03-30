package server

import (
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/route"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"github.com/LuuDinhTheTai/tzone/internal/service"
)

func (s *Server) MapHandlers() error {
	// Check if MongoDB is available
	if !s.HasMongoDB() {
		log.Println("⚠️ MongoDB is not available, brand and device routes will not work properly")
	}

	// Init repository
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
	return nil
}
