package storage

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

func NewCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("scylla:9042")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	return cluster
}

func TryToCreateKeyspace(session *gocqlx.Session) error {
	q := "CREATE KEYSPACE IF NOT EXISTS audit WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}"
	return session.Query(q, nil).Exec()
}

func TryToCreateTable(session *gocqlx.Session) error {
	q := `CREATE TABLE IF NOT EXISTS audit.logs (
		id UUID PRIMARY KEY,
		user_type TEXT,
		user_id UUID,
		event TEXT,
		auditable_type TEXT,
		auditable_id UUID,
		old_values TEXT,
		new_values TEXT,
		url TEXT,
		ip_address TEXT,
		user_agent TEXT,
		tags TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)`
	return session.Query(q, nil).Exec()
}
