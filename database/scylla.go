package database

import (
	"os"

	gocql "github.com/gocql/gocql"
)

func NewScyllaSession() *gocql.Session {
	var cluster = gocql.NewCluster(os.Getenv("MAIN_DB_HOST"))
	cluster.Keyspace = os.Getenv("MAIN_DB_KEYSPACE")
	cluster.Consistency = gocql.Quorum

	var session, err = cluster.CreateSession()
	if err != nil {
		panic("Failed to connect to cluster")
	}

	return session
}
