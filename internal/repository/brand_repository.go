package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BrandRepository struct {
	client          *mongo.Client
	database        string
	brandCollection string
}

type DeviceWithBrand struct {
	BrandID bson.ObjectID `bson:"brand_id"`
	Device  model.Device  `bson:"device"`
}

func NewBrandRepository() *BrandRepository {
	return &BrandRepository{
		database:        "Cluster0",
		brandCollection: "brands",
	}
}

// SetClient sets the MongoDB client for the repository
func (r *BrandRepository) SetClient(client *mongo.Client) {
	r.client = client
}

// GetBrandCollection returns the brand collection
func (r *BrandRepository) GetBrandCollection() *mongo.Collection {
	return r.client.Database(r.database).Collection(r.brandCollection)
}

// CreateBrand creates a new brand in MongoDB
func (r *BrandRepository) CreateBrand(ctx context.Context, brand *model.Brand) (*model.Brand, error) {
	collection := r.GetBrandCollection()

	result, err := collection.InsertOne(ctx, brand)
	if err != nil {
		log.Printf("❌ Error creating brand: %v", err)
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	brand.Id = result.InsertedID.(bson.ObjectID)
	log.Printf("✅ Brand created successfully with ID: %s", brand.Id.Hex())
	return brand, nil
}

// GetBrandById retrieves a brand by its ID
func (r *BrandRepository) GetBrandById(ctx context.Context, id string) (*model.Brand, error) {
	collection := r.GetBrandCollection()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid brand ID format: %w", err)
	}

	var brand model.Brand
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&brand)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("brand not found")
		}
		log.Printf("❌ Error finding brand: %v", err)
		return nil, fmt.Errorf("failed to find brand: %w", err)
	}

	log.Printf("✅ Brand found: %s", brand.Name)
	return &brand, nil
}

// GetAllBrands retrieves paginated brands from MongoDB
func (r *BrandRepository) GetAllBrands(ctx context.Context, page int, limit int) ([]model.Brand, int64, error) {
	collection := r.GetBrandCollection()
	skip := int64((page - 1) * limit)

	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("❌ Error counting brands: %v", err)
		return nil, 0, fmt.Errorf("failed to count brands: %w", err)
	}

	opts := options.Find().SetSkip(skip).SetLimit(int64(limit)).SetSort(bson.M{"_id": -1})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Printf("❌ Error fetching brands: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch brands: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("⚠️ Error closing cursor: %v", err)
		}
	}()

	var brands []model.Brand
	if err = cursor.All(ctx, &brands); err != nil {
		log.Printf("❌ Error decoding brands: %v", err)
		return nil, 0, fmt.Errorf("failed to decode brands: %w", err)
	}

	log.Printf("✅ Retrieved %d brands", len(brands))
	return brands, total, nil
}

// UpdateBrand updates an existing brand in MongoDB
func (r *BrandRepository) UpdateBrand(ctx context.Context, id string, brand *model.Brand) (*model.Brand, error) {
	collection := r.GetBrandCollection()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid brand ID format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"brand_name": brand.Name,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		log.Printf("❌ Error updating brand: %v", err)
		return nil, fmt.Errorf("failed to update brand: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("brand not found")
	}

	log.Printf("✅ Brand updated successfully: %s", id)
	return r.GetBrandById(ctx, id)
}

// DeleteBrand deletes a brand from MongoDB
func (r *BrandRepository) DeleteBrand(ctx context.Context, id string) error {
	collection := r.GetBrandCollection()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		log.Printf("❌ Error deleting brand: %v", err)
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("brand not found")
	}

	log.Printf("✅ Brand deleted successfully: %s", id)
	return nil
}

// ========================= Device===========================

// AddDeviceToBrand pushes a new device into the brand's devices array
func (r *BrandRepository) AddDeviceToBrand(ctx context.Context, brandID bson.ObjectID, device *model.Device) error {
	collection := r.GetBrandCollection()

	update := bson.M{
		"$push": bson.M{"devices": device},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": brandID}, update)
	if err != nil {
		log.Printf("❌ Error adding device to brand: %v", err)
		return fmt.Errorf("failed to add device to brand: %w", err)
	}

	return nil
}

// UpdateDeviceInBrand updates an existing device within the brand's device array
func (r *BrandRepository) UpdateDeviceInBrand(ctx context.Context, brandID bson.ObjectID, device *model.Device) error {
	collection := r.GetBrandCollection()

	filter := bson.M{"_id": brandID, "devices._id": device.ID}
	update := bson.M{
		"$set": bson.M{
			"devices.$": device,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("❌ Error updating device in brand: %v", err)
		return fmt.Errorf("failed to update device in brand: %w", err)
	}

	return nil
}

// RemoveDeviceFromBrand removes a device from the brand's devices array
func (r *BrandRepository) RemoveDeviceFromBrand(ctx context.Context, brandID bson.ObjectID, deviceID bson.ObjectID) error {
	collection := r.GetBrandCollection()

	update := bson.M{
		"$pull": bson.M{"devices": bson.M{"_id": deviceID}},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": brandID}, update)
	if err != nil {
		log.Printf("❌ Error removing device from brand: %v", err)
		return fmt.Errorf("failed to remove device from brand: %w", err)
	}

	return nil
}

// GetDeviceById retrieves a device by its ID
func (r *BrandRepository) GetDeviceById(ctx context.Context, id string) (*model.Device, string, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, "", fmt.Errorf("invalid device ID format")
	}

	filter := bson.M{"devices._id": objID}
	opts := options.FindOne().SetProjection(bson.M{
		"_id":     1,
		"devices": bson.M{"$elemMatch": bson.M{"_id": objID}},
	})

	var result struct {
		BrandID bson.ObjectID  `bson:"_id"`
		Devices []model.Device `bson:"devices"`
	}

	err = r.GetBrandCollection().FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", fmt.Errorf("device not found")
		}
		log.Printf("❌ Error finding device: %v", err)
		return nil, "", err
	}

	if len(result.Devices) == 0 {
		return nil, "", fmt.Errorf("device not found in array")
	}

	log.Printf("✅ Embedded device found: %s", result.Devices[0].ModelName)
	return &result.Devices[0], result.BrandID.Hex(), nil
}

// GetAllDevices retrieves paginated devices from all brand documents
func (r *BrandRepository) GetAllDevices(ctx context.Context, page int, limit int) ([]DeviceWithBrand, int64, error) {
	collection := r.GetBrandCollection()
	skip := int64((page - 1) * limit)

	countPipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$devices"}},
		{{Key: "$count", Value: "total"}},
	}

	countCursor, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		log.Printf("❌ Error counting devices: %v", err)
		return nil, 0, fmt.Errorf("failed to count devices: %w", err)
	}
	defer func() {
		if err := countCursor.Close(ctx); err != nil {
			log.Printf("⚠️ Error closing count cursor: %v", err)
		}
	}()

	var countResults []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResults); err != nil {
		log.Printf("❌ Error decoding device count: %v", err)
		return nil, 0, fmt.Errorf("failed to decode device count: %w", err)
	}

	var total int64
	if len(countResults) > 0 {
		total = countResults[0].Total
	}

	dataPipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$devices"}},
		{{Key: "$project", Value: bson.M{"brand_id": "$_id", "device": "$devices"}}},
		{{Key: "$sort", Value: bson.M{"device._id": -1}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := collection.Aggregate(ctx, dataPipeline)
	if err != nil {
		log.Printf("❌ Error fetching devices: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch devices: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("⚠️ Error closing cursor: %v", err)
		}
	}()

	var devices []DeviceWithBrand
	if err = cursor.All(ctx, &devices); err != nil {
		log.Printf("❌ Error decoding devices: %v", err)
		return nil, 0, fmt.Errorf("failed to decode devices: %w", err)
	}

	log.Printf("✅ Retrieved %d devices across all brands", len(devices))
	return devices, total, nil
}

// GetDevicesByBrandID retrieves paginated devices from a single brand document
func (r *BrandRepository) GetDevicesByBrandID(ctx context.Context, brandID bson.ObjectID, page int, limit int) ([]DeviceWithBrand, int64, error) {
	collection := r.GetBrandCollection()
	skip := int64((page - 1) * limit)

	countPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": brandID}}},
		{{Key: "$unwind", Value: "$devices"}},
		{{Key: "$count", Value: "total"}},
	}

	countCursor, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		log.Printf("❌ Error counting devices by brand: %v", err)
		return nil, 0, fmt.Errorf("failed to count devices by brand: %w", err)
	}
	defer func() {
		if err := countCursor.Close(ctx); err != nil {
			log.Printf("⚠️ Error closing count cursor: %v", err)
		}
	}()

	var countResults []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResults); err != nil {
		log.Printf("❌ Error decoding device count by brand: %v", err)
		return nil, 0, fmt.Errorf("failed to decode device count by brand: %w", err)
	}

	var total int64
	if len(countResults) > 0 {
		total = countResults[0].Total
	}

	dataPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": brandID}}},
		{{Key: "$unwind", Value: "$devices"}},
		{{Key: "$project", Value: bson.M{"brand_id": "$_id", "device": "$devices"}}},
		{{Key: "$sort", Value: bson.M{"device._id": -1}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := collection.Aggregate(ctx, dataPipeline)
	if err != nil {
		log.Printf("❌ Error fetching devices by brand: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch devices by brand: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("⚠️ Error closing cursor: %v", err)
		}
	}()

	var devices []DeviceWithBrand
	if err = cursor.All(ctx, &devices); err != nil {
		log.Printf("❌ Error decoding devices by brand: %v", err)
		return nil, 0, fmt.Errorf("failed to decode devices by brand: %w", err)
	}

	log.Printf("✅ Retrieved %d devices for brand %s", len(devices), brandID.Hex())
	return devices, total, nil
}
