package repository

import (
	"context"
	pb "go-todolist/task-srv/proto/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	Delete(ctx context.Context, task *pb.Task) error
	Modify(ctx context.Context, task *pb.Task) error
	Finished(ctx context.Context, task *pb.Task) error
	Count(ctx context.Context, task *pb.Task) error
	Search(ctx context.Context, task *pb.Task) error
}

// 数据库连接实现类
type TaskRepositoryImpl struct {
	// 需要注入一个数据库连接客户端
	Conn *mongo.Client
}

// 定义默认的操作表
func (repo TaskRepositoryImpl) collection() *mongo.Collection {
	return repo.Conn.Database(DbName).Collection(TaskCollection)
}

func (repo TaskRepositoryImpl) InsertOnce(ctx context.Context, task *pb.Task) error {
	_, err := repo.collection().InsertOne(ctx, bson.M{
		"body":       task.Body,
		"startTime":  task.StartTime,
		"endTime":    task.EndTime,
		"isFinished": task.IsFinished,
		"createTime": time.Now().Unix(),
	})
	return err
}
