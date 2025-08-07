package cassandra

import (
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var session *gocql.Session

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cluster := gocql.NewCluster(os.Getenv("DB_HOST"))
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	} 
}

func GetSession() *gocql.Session {
	return session
}