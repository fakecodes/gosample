package mongo

import (
	"context"

	"github.com/fakecodes/gosample/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTaskRepository struct {
	Conn *mongo.Database
}

func NewMongoTaskRepository(conn *mongo.Database) domain.TaskRepository {
	return &mongoTaskRepository{conn}
}

func (m *mongoTaskRepository) Create(ctx context.Context, a *domain.Task) (err error) {
	collection := m.Conn.Collection("tasks")
	res, err := collection.InsertOne(context.Background(), bson.M{"name": a.Name, "description": a.Description, "created_at": a.CreatedAt, "due_date": a.DueDate, "priority": a.Priority, "completed": a.Completed})
	if err != nil {
		return err
	}
	a.ID = res.InsertedID.(primitive.ObjectID)
	return
}
