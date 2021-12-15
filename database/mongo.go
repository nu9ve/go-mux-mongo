package writeon

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// config "./config/config"
	"writeon/model"
)

const DefaultDatabase = "writeon"
const UserCollection = "user"
const EssayCollection = "essay"
const TopicCollection = "topic"


// MongoHandler struct
type MongoHandler struct {
	client   *mongo.Client
	database string
}

//NewHandler MongoHandler Constructor
func NewHandler(address string) *MongoHandler {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cl, _ := mongo.Connect(ctx, options.Client().ApplyURI(address))
	mh := &MongoHandler{
		client:   cl,
		database: DefaultDatabase,
	}
	return mh
}


// USER QUERIES

// GetUser function
func (mh *MongoHandler) GetUser(u *model.User, filter interface{}) error {
	// Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection(UserCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(u)
	return err
}

// CreateUser function
func (mh *MongoHandler) CreateUser(u *model.UserTemplate) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection(UserCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, u)
	return result, err
}


// TOPIC QUERIES

// GetTopics function
func (mh *MongoHandler) GetTopics(filter interface{}) []*model.Topic {
	collection := mh.client.Database(mh.database).Collection(TopicCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*model.Topic
	for cur.Next(ctx) {
		topic := &model.Topic{}
		er := cur.Decode(topic)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, topic)
	}
	return result
}

// AddOneTopic function
func (mh *MongoHandler) AddOneTopic(t *model.TopicTemplate) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection(TopicCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, t)
	return result, err
}

// GetOneTopic function
func (mh *MongoHandler) GetOneTopic(t *model.Topic, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection(TopicCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(t)
	return err
}

// RemoveOneTopic function
func (mh *MongoHandler) RemoveOneTopic(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection(TopicCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}

// GetOneEssay function
func (mh *MongoHandler) GetOneEssay(e *model.Essay, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection(EssayCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(e)
	return err
}

// GetEssays function
func (mh *MongoHandler) GetEssays(filter interface{}) []*model.Essay {
	collection := mh.client.Database(mh.database).Collection(EssayCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*model.Essay
	for cur.Next(ctx) {
		essay := &model.Essay{}
		er := cur.Decode(essay)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, essay)
	}
	return result
}

// AddOneEssay function
func (mh *MongoHandler) AddOneEssay(e *model.EssayTemplate) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection(EssayCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, e)
	return result, err
}

// UpdateEssay function
func (mh *MongoHandler) UpdateEssay(c *model.Essay, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mh.client.Database(mh.database).Collection(EssayCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.UpdateMany(ctx, filter, update)
	return result, err
}

// RemoveOneEssay function
func (mh *MongoHandler) RemoveOneEssay(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection(EssayCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}