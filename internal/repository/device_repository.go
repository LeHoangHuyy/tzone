package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DeviceRepository struct {
	client           *mongo.Client
	database         string
	deviceCollection string
}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{
		database:         "Cluster0",
		deviceCollection: "devices",
	}
}

// SetClient sets the MongoDB client for the repository
func (r *DeviceRepository) SetClient(client *mongo.Client) {
	r.client = client
}

// GetDeviceCollection returns the device collection
func (r *DeviceRepository) GetDeviceCollection() *mongo.Collection {
	return r.client.Database(r.database).Collection(r.deviceCollection)
}

// CreateDevice creates a new device in MongoDB
func (r *DeviceRepository) CreateDevice(ctx context.Context, device *model.Device) (*model.Device, error) {
	collection := r.GetDeviceCollection()

	result, err := collection.InsertOne(ctx, device)
	if err != nil {
		log.Printf("❌ Error creating device: %v", err)
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	device.ID = result.InsertedID.(bson.ObjectID)
	log.Printf("✅ Device created successfully with ID: %s", device.ID.Hex())
	return device, nil
}

// GetDeviceById retrieves a device by its ID
func (r *DeviceRepository) GetDeviceById(ctx context.Context, id string) (*model.Device, error) {
	collection := r.GetDeviceCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid device ID format: %w", err)
	}

	var device model.Device
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("device not found")
		}
		log.Printf("❌ Error finding device: %v", err)
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	log.Printf("✅ Device found: %s", device.ModelName)
	return &device, nil
}

// GetAllDevices retrieves all Devices from MongoDB
func (r *DeviceRepository) GetAllDevices(ctx context.Context) ([]model.Device, error) {
	collection := r.GetDeviceCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("❌ Error fetching devices: %v", err)
		return nil, fmt.Errorf("failed to fetch devices: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("⚠️ Error closing cursor: %v", err)
		}
	}()

	var devices []model.Device
	if err = cursor.All(ctx, &devices); err != nil {
		log.Printf("❌ Error decoding devices: %v", err)
		return nil, fmt.Errorf("failed to decode devices: %w", err)
	}

	log.Printf("✅ Retrieved %d devices", len(devices))
	return devices, nil
}

// UpdateDevice updates an existing device in MongoDB
func (r *DeviceRepository) UpdateDevice(ctx context.Context, id string, device *model.Device) (*model.Device, error) {
	collection := r.GetDeviceCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid device ID format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"model_name":     device.ModelName,
			"specifications": device.Specifications,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Printf("❌ Error updating device: %v", err)
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("device not found")
	}

	log.Printf("✅ Device updated successfully: %s", id)
	return r.GetDeviceById(ctx, id)
}

// DeleteDevice deletes a device from MongoDB
func (r *DeviceRepository) DeleteDevice(ctx context.Context, id string) error {
	collection := r.GetDeviceCollection()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid device ID format: %w", err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		log.Printf("❌ Error deleting device: %v", err)
		return fmt.Errorf("failed to delete device: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("device not found")
	}

	log.Printf("✅ Device deleted successfully: %s", id)
	return nil
}
