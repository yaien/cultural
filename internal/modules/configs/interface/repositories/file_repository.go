package repositories

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.FileRepository = (*FileRepository)(nil)

type FileRepository struct {
	db *mongo.Database
}

func NewFileRepository(db *mongo.Database) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, file *models.File) error {

	res, err := r.db.Collection("files").InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("failed inserting file: %w", err)
	}

	file.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *FileRepository) Get(ctx context.Context, organizationID primitive.ObjectID, name string) (*models.File, error) {
	var file models.File
	err := r.db.Collection("files").FindOne(ctx, bson.M{"organizationId": organizationID, "name": name}).Decode(&file)
	if err != nil {
		return nil, translate(err)
	}
	return &file, nil
}

func (r *FileRepository) Delete(ctx context.Context, organizationID primitive.ObjectID, name string) error {
	res, err := r.db.Collection("files").DeleteOne(ctx, bson.M{"organizationId": organizationID, "name": name})
	if err != nil {
		return fmt.Errorf("failed deleting file: %w", err)
	}
	if res.DeletedCount == 0 {
		return translate(mongo.ErrNoDocuments)
	}
	return nil
}

func (r *FileRepository) Rename(ctx context.Context, organizationID primitive.ObjectID, oldName string, newName string) error {
	res, err := r.db.Collection("files").UpdateOne(ctx, bson.M{"organizationId": organizationID, "name": oldName}, bson.M{"$set": bson.M{"name": newName}})
	if err != nil {
		return fmt.Errorf("failed renaming file: %w", err)
	}
	if res.MatchedCount == 0 {
		return translate(mongo.ErrNoDocuments)
	}
	return nil
}

func (r *FileRepository) List(ctx context.Context, organizationID primitive.ObjectID) ([]*models.File, error) {
	var files []*models.File

	cursor, err := r.db.Collection("files").Find(ctx, bson.M{"organizationId": organizationID})
	if err != nil {
		return nil, fmt.Errorf("failed finding files: %w", err)
	}

	if err := cursor.All(ctx, &files); err != nil {
		return nil, fmt.Errorf("failed decoding files: %w", err)
	}

	return files, nil
}
