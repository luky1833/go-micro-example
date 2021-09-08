package handler

import (
	"context"
	"github.com/pkg/errors"
	pb "go-todolist/task-srv/proto/task"
	"go-todolist/task-srv/repository"
)

type TaskHandler struct {
	TaskRepository repository.TaskRepository
}

func (t *TaskHandler) Create(ctx context.Context, task *pb.Task, resp *pb.EditResponse) error {
	if task.Body == "" || task.StartTime <= 0 || task.EndTime <= 0 {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.InsertOnce(ctx, task); err != nil {
		return err
	}
	resp.Msg = "success"
	return nil
}

func (t *TaskHandler) Delete(ctx context.Context, task *pb.Task, response *pb.EditResponse) error {
	if task.Id == "" {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Delete(ctx, task.Id); err != nil {
		return err
	}
	response.Msg = task.Id
	return nil
}

func (t *TaskHandler) Modify(ctx context.Context, task *pb.Task, response *pb.EditResponse) error {
	if task.Id == "" || task.Body == "" || task.StartTime <= 0 || task.EndTime <= 0 {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Modify(ctx, task); err != nil {
		return err
	}
	response.Msg = "success"
	return nil
}

func (t *TaskHandler) Finished(ctx context.Context, task *pb.Task, response *pb.EditResponse) error {
	if task.Id == "" || task.IsFinished != repository.Unfinished && task.IsFinished != repository.Finished {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Finished(ctx, task); err != nil {
		return err
	}
	response.Msg = "success"
	return nil
}

func (t *TaskHandler) Search(ctx context.Context, request *pb.SearchRequest, response *pb.SearchResponse) error {
	count, err := t.TaskRepository.Count(ctx, request.Keyword)
	if err != nil {
		return errors.WithMessage(err, "count row number")
	}
	if request.PageCode <= 0 {
		request.PageCode = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 20
	}
	if request.SortBy == "" {
		request.SortBy = "createTime"
	}
	if request.Order == 0 {
		request.Order = -1
	}
	if request.PageSize*(response.PageCode-1) > count {
		return errors.New("there`s not that much data")
	}
	rows, err2 := t.TaskRepository.Search(ctx, request)
	if err2 != nil {
		return errors.WithMessage(err2, "search data")
	}
	*response = pb.SearchResponse{
		PageCode: request.PageCode,
		PageSize: request.PageSize,
		SortBy:   request.SortBy,
		Order:    request.Order,
		Rows:     rows,
	}
	return nil

}
