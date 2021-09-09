package repository

import (
	"context"
	"github.com/pkg/errors"
	pb "go-todolist/task-srv/proto/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

const (
	//默认数据库名
	DbName = "todolist"
	// 默认表名
	TaskCollection = "task"
	Unfinished     = 0
	Finished       = 1
)

// 定义数据库基本的增删改查的操作
type TaskRepository interface {
	InsertOnce(ctx context.Context, task *pb.Task) error
	Delete(ctx context.Context, id string) error
	Modify(ctx context.Context, task *pb.Task) error
	Finished(ctx context.Context, task *pb.Task) error
	Count(ctx context.Context, keyword string) (int64, error)
	Search(ctx context.Context, request *pb.SearchRequest) ([]*pb.Task, error)
	// 接口新增方法
	FindById(ctx context.Context, id string) (*pb.Task, error)
}

// 数据库连接实现类
type TaskRepositoryImpl struct {
	// 需要注入一个数据库连接客户端
	Conn *mongo.Client
}

// 定义默认的操作表
func (repo *TaskRepositoryImpl) collection() *mongo.Collection {
	return repo.Conn.Database(DbName).Collection(TaskCollection)
}

func (repo *TaskRepositoryImpl) InsertOnce(ctx context.Context, task *pb.Task) error {
	_, err := repo.collection().InsertOne(ctx, bson.M{
		"body":       task.Body,
		"startTime":  task.StartTime,
		"endTime":    task.EndTime,
		"isFinished": Unfinished,
		"createTime": time.Now().Unix(),
		//
	})
	return err
}

func (repo *TaskRepositoryImpl) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.collection().DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (repo *TaskRepositoryImpl) Modify(ctx context.Context, task *pb.Task) error {
	id, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		return err
	}
	_, err = repo.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"body":       task.Body,
		"startTime":  task.StartTime,
		"endTime":    task.EndTime,
		"updateTime": time.Now().Unix(),
	}})
	return err
}

func (repo *TaskRepositoryImpl) Finished(ctx context.Context, task *pb.Task) error {
	id, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	update := bson.M{
		"isFinished": int32(task.IsFinished),
		"updateTime": now,
	}
	if task.IsFinished == Finished {
		update["finishTime"] = now
	}
	log.Print(task)
	log.Println(update)
	_, err = repo.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
func (repo *TaskRepositoryImpl) Count(ctx context.Context, keyword string) (int64, error) {
	filter := bson.M{}
	if keyword != "" && strings.TrimSpace(keyword) != "" {
		filter = bson.M{"body": bson.M{"$regex": keyword}}
	}
	count, err := repo.collection().CountDocuments(ctx, filter)
	return count, err
}

func (repo *TaskRepositoryImpl) Search(ctx context.Context, request *pb.SearchRequest) ([]*pb.Task, error) {
	filter := bson.M{}
	if request.Keyword != "" && strings.TrimSpace(request.Keyword) != "" {
		filter = bson.M{"body": bson.M{"$regex": request.Keyword}}
	}
	cursor, err := repo.collection().Find(
		ctx,
		filter,
		options.Find().SetSkip((request.PageCode-1)*request.PageSize),
		options.Find().SetLimit(request.PageSize),
		options.Find().SetSort(bson.M{request.SortBy: request.Order}))
	if err != nil {
		return nil, errors.WithMessage(err, "search mongo")
	}
	var rows []*pb.Task

	if err := cursor.All(ctx, &rows); err != nil {
		return nil, errors.WithMessage(err, "parse data")
	}
	return rows, nil
}
