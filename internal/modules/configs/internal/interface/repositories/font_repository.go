package repositories

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ models.FontRepository = (*FontRepository)(nil)

type FontRepository struct {
	db *mongo.Database
}

func NewFontRepository(db *mongo.Database) *FontRepository {
	return &FontRepository{db}
}

func (r *FontRepository) GetByFamily(ctx context.Context, family string) (*models.Font, error) {
	collection := r.db.Collection("fonts")

	var font models.Font

	err := collection.FindOne(ctx, bson.M{"family": family}).Decode(&font)
	if err != nil {
		return nil, translate(err)
	}

	return &font, nil
}

func (r *FontRepository) Find(ctx context.Context, opts *models.FindFontOptions) ([]*models.Font, error) {
	collection := r.db.Collection("fonts")
	filter := bson.M{}

	if opts.Family != "" {
		filter["family"] = bson.M{"$regex": opts.Family, "$options": "i"}
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	if opts.Offset < 0 {
		opts.Offset = 0
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetSkip(opts.Offset).SetLimit(opts.Limit))
	if err != nil {
		return nil, fmt.Errorf("failed finding fonts: %w", err)
	}

	defer cursor.Close(ctx)

	var fonts []*models.Font

	err = cursor.All(ctx, &fonts)
	if err != nil {
		return nil, fmt.Errorf("failed decoding fonts: %w", err)
	}

	return fonts, nil
}
