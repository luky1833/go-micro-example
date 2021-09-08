package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/pkg/errors"
	"go-todolist/task-srv/handler"
	pb "go-todolist/task-srv/proto/task"
	"go-todolist/task-srv/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const MNOGO_URL = "mongodb://admin:123456@10.10.10.123:27017"

func main() {
	// 在日志中打印文件路径,便于调试代码
	log.SetFlags(log.Llongfile)
	// 连接 mongo
	conn, err := connectMongo(MNOGO_URL, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	// 使用匿名函数关闭
	defer func(conn *mongo.Client, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {
			log.Println("关闭")
		}
	}(conn, context.Background())
	//new service
	service := micro.NewService(
		micro.Name("go.micro.service.task"),
		micro.Version("latest"),
		// 配置etcd为注册中心，配置etcd路径，默认端口是2379
		micro.Registry(etcd.NewRegistry(
			// 地址是我本地etcd服务器地址，不要照抄
			registry.Addrs("127.0.0.1:2379"),
		)),
	)
	// initialise service 初始化
	service.Init()
	// 接受 handler
	taskHandler := &handler.TaskHandler{
		TaskRepository: &repository.TaskRepositoryImpl{
			Conn: conn,
		},
	}
	if err := pb.RegisterTaskServiceHandler(service.Server(), taskHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}
	// run server
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}

}

func connectMongo(url string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, errors.WithMessage(err, "create mongo connection session")
	}
	return connect, nil
}
