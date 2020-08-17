package mongodb_use

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func CreateClent() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return client
}

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	Msg   string    `bson:"msg"`
	Point TimePoint `bson:"point"`
}

//插入数据
func TestInsertMongoDb(t *testing.T) {
	var (
		database     *mongo.Database
		collection   *mongo.Collection
		err          error
		log          *LogRecord
		insertResult *mongo.InsertOneResult
		docId        primitive.ObjectID
	)

	client := CreateClent()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//1 连接
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Error(err.Error())
	}

	//mongodb数据库
	database = client.Database("runoob")

	//表
	collection = database.Collection("mycol")

	log = &LogRecord{
		Msg:   "这是一个测试结果",
		Point: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	insertResult, err = collection.InsertOne(context.TODO(), log)
	if err != nil && insertResult != nil {
		t.Error(err.Error())
		return
	}
	docId = insertResult.InsertedID.(primitive.ObjectID)

	fmt.Println(docId.String())
}

func TestFindMongodb(t *testing.T) {
	var (
		database   *mongo.Database
		collection *mongo.Collection
		err        error
		log        *LogRecord
		cursor     *mongo.Cursor
	)
	client := CreateClent()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//1 连接
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Error(err.Error())
	}

	//mongodb数据库
	database = client.Database("runoob")
	//表
	collection = database.Collection("mycol")

	log = &LogRecord{}

	var findoption = &options.FindOptions{}
	findoption.SetSkip(0)
	findoption.SetLimit(5)

	if cursor, err = collection.Find(context.TODO(), bson.D{}); err != nil {
		t.Error(err.Error())
		return
	}

	for cursor.Next(context.TODO()) {
		log = &LogRecord{}
		err := cursor.Decode(log)
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Log(log)
	}
}

func TestDeleteMongodb(t *testing.T) {
	var (
		database   *mongo.Database
		collection *mongo.Collection
		err        error
		cursor     *mongo.DeleteResult
	)
	client := CreateClent()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Error(err.Error())
	}

	//mongodb数据库
	database = client.Database("runoob")
	//表
	collection = database.Collection("mycol")

	if cursor, err = collection.DeleteMany(context.TODO(), bson.D{}); err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(cursor.DeletedCount)
}
