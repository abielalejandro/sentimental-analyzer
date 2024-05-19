package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abielalejandro/control/config"
	"github.com/gocql/gocql"
)

type CMessage struct {
	Id          gocql.UUID        `cql:"id"`
	msg         string            `cql:"msg"`
	msgAnalyzed SentimentalResult `cql:"msg_analyzed"`
	CreatedAt   time.Time         `cql:"created_at"`
	UpdateddAt  time.Time         `cql:"updated_at"`
	ExpiresAt   time.Time         `cql:"expires_at"`
}

type CassandraStorage struct {
	*gocql.ClusterConfig
	*gocql.Session
	*config.Config
}

func NewCassandraStorage(config *config.Config) *CassandraStorage {
	addr := config.Storage.Addr
	hosts := strings.Split(addr, ",")
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Keyspace = config.Storage.Db
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.ProtoVersion = 4

	session, err := clusterConfig.CreateSession()

	if err != nil {
		panic(err)
	}

	return &CassandraStorage{
		ClusterConfig: clusterConfig,
		Session:       session,
	}
}

func (storage *CassandraStorage) Create(
	ctx context.Context,
	msg *Message) (bool, error) {

	ttl := storage.Config.Storage.Ttl * 60
	command := fmt.Sprintf("INSERT INTO messages (id, msg,created_at,updated_at,expires_at) VALUES(?,?,?,?) USING TTL %v;", ttl)
	err := storage.Query(command,
		msg.Id,
		msg.msg,
		time.Now(),
		time.Now(),
		time.Now().Add(time.Minute*10),
	).WithContext(ctx).Exec()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *CassandraStorage) Update(ctx context.Context, id string, result *SentimentalResult) (bool, error) {
	ttl := storage.Config.Storage.Ttl * 60
	command := fmt.Sprintf("UPDATE messages USING TTL %v SET msg_analyzed=?,updated_at=? WHERE id='?'", ttl)

	err := storage.Query(command, result, time.Now(), id).WithContext(ctx).Exec()

	if err != nil {
		return false, err
	}

	return true, nil
}
