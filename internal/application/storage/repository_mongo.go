package storage

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ Repository = (*Mongo)(nil)

type Mongo struct {
	collection *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{db.Collection("files")}
}

func (r *Mongo) Create(ctx context.Context, file *File) error {
	res, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("failed inserting file: %w", err)
	}

	file.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Mongo) Update(ctx context.Context, file *File) error {
	res, err := r.collection.ReplaceOne(ctx, bson.M{"_id": file.ID}, file)
	if err != nil {
		return fmt.Errorf("failed updating file: %w", err)
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *Mongo) GetByID(ctx context.Context, id primitive.ObjectID) (*File, error) {
	var file File
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	switch err {
	case nil:
		return &file, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New("not_found", err)
	default:
		return nil, fmt.Errorf("failed finding file: %w", err)
	}
}

func (r *Mongo) GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error) {
	var file File
	err := r.collection.FindOne(ctx, bson.M{"organizationId": organizationID, "name": name}).Decode(&file)
	switch err {
	case nil:
		return &file, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New("not_found", err)
	default:
		return nil, fmt.Errorf("failed finding file: %w", err)
	}

}

func (r *Mongo) GetByOrganizationIDAndID(ctx context.Context, organizationID, id primitive.ObjectID) (*File, error) {
	var file File
	err := r.collection.FindOne(ctx, bson.M{"organizationId": organizationID, "_id": id}).Decode(&file)
	switch err {
	case nil:
		return &file, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New("not_found", err)
	default:
		return nil, fmt.Errorf("failed finding file: %w", err)
	}
}

func (r *Mongo) DeleteByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"organizationId": organizationID, "name": name})
	if err != nil {
		return fmt.Errorf("failed deleting file: %w", err)
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *Mongo) RenameByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, oldName string, newName string) error {
	res, err := r.collection.UpdateOne(ctx, bson.M{"organizationId": organizationID, "name": oldName}, bson.M{"$set": bson.M{"name": newName}})
	if err != nil {
		return fmt.Errorf("failed renaming file: %w", err)
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *Mongo) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error) {
	var files []*File

	cursor, err := r.collection.Find(ctx, bson.M{"organizationId": organizationID})
	if err != nil {
		return nil, fmt.Errorf("failed finding files: %w", err)
	}

	if err := cursor.All(ctx, &files); err != nil {
		return nil, fmt.Errorf("failed decoding files: %w", err)
	}

	return files, nil
}
