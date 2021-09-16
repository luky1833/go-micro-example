package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/pkg/errors"
	"go-todolist/achievement-srv/repository"
	"go-todolist/achievement-srv/subscriber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

//
// 这里是我内网的mongo地址，请根据你得实际情况配置，推荐使用dockers部署
const MNOGO_URL = "mongodb://admin:123456@10.10.10.123:27017"

func main() {
	// 在日志中打印文件路径，便于调试代码
	log.SetFlags(log.Llongfile)

	conn, err := connectMongo(MNOGO_URL, time.Second)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Disconnect(context.Background())

	// new service
	service := micro.NewService(
		micro.Name("go.micro.service.achievement"),
		micro.Version("latest"),
		// 配置etcd为注册中心，配置etcd路径，默认端口是2379
		micro.Registry(etcd.NewRegistry(
			// 地址是我本地etcd服务器地址，不要照抄
			registry.Addrs("127.0.0.1:2379"),
		)))
	// initialise service
	service.Init()
	// Register Handler
	handler := &subscriber.AchievementSub{
		Repo: &repository.AchievementRepoImpl{
			Conn: conn,
		},
	}
	// 这里的 topic注意要与 task-srv 注册的要保持一致
	if err := micro.RegisterSubscriber("go.micro.service.task.finished", service.Server(), handler); err != nil {
		log.Fatal(errors.WithMessage(err, "subscribe"))
	}
	// run service
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}

//连接到mongo
func connectMongo(url string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, errors.WithMessage(err, "create mongo connection session")
	}
	return client, nil
}
