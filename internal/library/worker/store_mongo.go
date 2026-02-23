package worker

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Store = (*MongoStore)(nil)

const DefaultJobsCollection = "jobs"

type MongoStore struct {
	collection *mongo.Collection
	limit      int64
}

// NewMongoStore creates a new MongoStore with the given MongoDB database and collection name. If the collection name is empty, it defaults to "jobs".
func NewMongoStore(db *mongo.Database, collection string) *MongoStore {
	if collection == "" {
		collection = DefaultJobsCollection
	}

	return &MongoStore{
		collection: db.Collection(collection),
		limit:      10,
	}
}

func (m *MongoStore) Create(ctx context.Context, job Job) error {
	_, err := m.collection.InsertOne(ctx, job)
	return err
}

func (m *MongoStore) Update(ctx context.Context, job Job) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": job.ID}, bson.M{"$set": job})
	return err
}

func (m *MongoStore) Fetch(ctx context.Context) ([]Job, error) {

	cursor, err := m.collection.Find(ctx, bson.M{"status": StatusPending}, options.Find().SetSort(bson.M{"created_at": 1}).SetLimit(m.limit))
	if err != nil {
		return nil, err
	}

	var jobs []Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil

}
