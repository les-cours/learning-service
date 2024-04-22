package resolvers

import (
	"database/sql"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/orgs"
	"github.com/les-cours/learning-service/api/users"
	"sync"
)

type Server struct {
	DB    *sql.DB
	Users users.UserServiceClient
	Orgs  orgs.OrgServiceClient
	learning.UnimplementedLearningServiceServer
	sync.Mutex
}
