package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	pb "go-todolist/task-srv/proto/task"
	"go-todolist/task-srv/repository"
	"log"
	"time"
)

func main() {
	// 在日志中打印文件路径,便于调试代码
	log.SetFlags(log.Llongfile)
	// 客户端也注册为服务
	// 客户端也注册为服务
	server := micro.NewService(
		micro.Name("go.micro.client.task"),
		// 配置etcd为注册中心，配置etcd路径，默认端口是2379
		micro.Registry(etcd.NewRegistry(
			// 地址是我本地etcd服务器地址，不要照抄
			registry.Addrs("127.0.0.1:2379"),
		)),
	)
	server.Init()
	// 发现服务
	taskService := pb.NewTaskService("go.micro.service.task", server.Client())

	// 调用服务生成三条任务
	now := time.Now()
	insertTask(taskService, "完成作业1", now.Unix(), now.Add(time.Hour*24).Unix())
	insertTask(taskService, "完成作业2", now.Add(time.Hour*24).Unix(), now.Add(time.Hour*48).Unix())
	insertTask(taskService, "完成作业3", now.Add(time.Hour*48).Unix(), now.Add(time.Hour*72).Unix())

	// 分页查询任务
	page, err := taskService.Search(context.Background(), &pb.SearchRequest{
		PageSize: 20,
		PageCode: 1,
	})
	if err != nil {
		log.Fatal("search1", err)
	}
	log.Println(page)

	// 更新第一条记录生效
	row := page.Rows[0]
	fmt.Println(row.Id)
	if _, err := taskService.Finished(context.Background(), &pb.Task{
		Id:         row.Id,
		IsFinished: repository.Finished,
	}); err != nil {
		log.Fatal("finished", row.Id, err)
	}
	// 修改查询到的第二条数据,延长截止日期
	row2 := page.Rows[1]
	if _, err := taskService.Modify(context.Background(), &pb.Task{
		Id:        row2.Id,
		Body:      row2.Body,
		StartTime: row2.StartTime,
		EndTime:   now.Add(time.Hour * 72).Unix(),
	}); err != nil {
		log.Fatal("modify", row2.Id, err)
	}
	// 删除第三条记录
	row3 := page.Rows[2]
	if _, err := taskService.Delete(context.Background(), &pb.Task{
		Id: row3.Id,
	}); err != nil {
		log.Fatal("delte", row3.Id, err)
	}
	// 再次分页查询,检查修改的结果
	page, err = taskService.Search(context.Background(), &pb.SearchRequest{})
	if err != nil {
		log.Fatal("search2", err)
	}
	log.Println(page)

}

func insertTask(service pb.TaskService, body string, start, end int64) {
	_, err := service.Create(context.Background(), &pb.Task{
		Body:      body,
		StartTime: start,
		EndTime:   end,
	})
	if err != nil {
		log.Fatal("create", err)
	}
	log.Println("create task success")
}
