package resolvers

import (
	"database/sql"
	apivideosdk "github.com/apivideo/api.video-go-client"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/orgs"
	"github.com/les-cours/learning-service/api/users"
	"github.com/les-cours/learning-service/database"
	"go.uber.org/zap"
	"sync"
)

type Server struct {
	DB       *sql.DB
	MongoDB  database.MongoClient
	Users    users.UserServiceClient
	Orgs     orgs.OrgServiceClient
	Logger   *zap.Logger
	VideoApi *apivideosdk.Client
	learning.UnimplementedLearningServiceServer
	sync.Mutex
}
