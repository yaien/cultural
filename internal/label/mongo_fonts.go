package label

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ FontRepository = (*MongoFonts)(nil)

type MongoFonts struct {
	collection *mongo.Collection
}

func NewMongoFonts(db *mongo.Database) *MongoFonts {
	return &MongoFonts{db.Collection("fonts")}
}

func (r *MongoFonts) GetByFamily(ctx context.Context, family string) (*Font, error) {
	var font Font
	err := r.collection.FindOne(ctx, bson.M{"family": family}).Decode(&font)
	switch err {
	case nil:
		return &font, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *MongoFonts) Find(ctx context.Context, opts *FindFontOptions) (fonts []*Font, err error) {
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

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSkip(opts.Offset).SetLimit(opts.Limit))
	if err != nil {
		return nil, fmt.Errorf("failed finding fonts: %w", err)
	}

	defer func() {
		if derr := cursor.Close(ctx); derr != nil {
			err = derr
		}
	}()

	err = cursor.All(ctx, &fonts)
	if err != nil {
		return nil, fmt.Errorf("failed decoding fonts: %w", err)
	}

	return fonts, nil
}
