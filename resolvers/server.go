package resolvers

import (
	"database/sql"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/orgs"
	"github.com/les-cours/learning-service/api/users"
	"go.uber.org/zap"
	"sync"
)

type Server struct {
	DB     *sql.DB
	Users  users.UserServiceClient
	Orgs   orgs.OrgServiceClient
	Logger *zap.Logger
	learning.UnimplementedLearningServiceServer
	sync.Mutex
}
