package cassandra

import (
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var cluster *gocql.ClusterConfig

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	
	cluster := gocql.NewCluster(os.Getenv("DB_HOST"))
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum
}

func GetSession() (*gocql.Session, error) {
	return cluster.CreateSession()
}