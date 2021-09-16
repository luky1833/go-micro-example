package repository

// 编写数据库操作
import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	// 默认数据库名
	DbName = "todolist"
	// 默认表名
	TaskCollection = "achievement"
)

type Achievement struct {
	Id string `bson:"_id,omitempty"`
	// 用户id
	UserId string `bson:"userId"`
	// 完成任务总数
	Total int64 `bson:"total"`
	//完成第一个任务的时间
	Finished1Time int64 `bson:"finished1Time"`
	//完成第100个任务的时间
	Finished100Time int64 `bson:"finished100Time"`
	//完成第1000个任务的时间
	Finished1000Time int64 `bson:"finished1000Time"`
	//更新时间
	UpdateTime int64 `bson:"updateTime"`
}

// 创建接口
// 因为只是演示,这里我们定义查询和保存
type AchievementRepo interface {
	FindByUserId(ctx context.Context, userId string) (*Achievement, error)
	Insert(ctx context.Context, achievement *Achievement) error
	Update(ctx context.Context, achievement *Achievement) error
}

// 创建结构体数据库连接,实现方法 AchievementRepo
type AchievementRepoImpl struct {
	Conn *mongo.Client
}

// 默认的数据库连接
func (repo *AchievementRepoImpl) collection() *mongo.Collection {
	return repo.Conn.Database(DbName).Collection(TaskCollection)
}

// 下面方法实现接口
func (repo *AchievementRepoImpl) FindByUserId(ctx context.Context, userId string) (*Achievement, error) {
	result := repo.collection().FindOne(ctx, bson.M{"userId": userId})
	// findOne 要是查不到是会报错的,这里要处理一下
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	achievement := &Achievement{}
	if err := result.Decode(achievement); err != nil {
		return nil, errors.WithMessage(err, "search mongo")
	}
	return achievement, nil
}

func (repo *AchievementRepoImpl) Insert(ctx context.Context, achievement *Achievement) error {
	_, err := repo.collection().InsertOne(ctx, achievement)
	return err
}

func (repo AchievementRepoImpl) Update(ctx context.Context, achievement *Achievement) error {
	achievement.UpdateTime = time.Now().Unix()
	oid, err := primitive.ObjectIDFromHex(achievement.Id)
	if err != nil {
		return err
	}
	achievement.Id = ""
	_, err = repo.collection().UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": achievement})
	return err
}
