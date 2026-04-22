package cassandra

import (
	"strings"

	"github.com/gocql/gocql"
)

func NewSession(dbHost string, keyspace string, consistency string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(dbHost)
	cluster.Keyspace = keyspace
	switch strings.ToLower(consistency) {
	case "any":
		cluster.Consistency = gocql.Any
	case "one":
		cluster.Consistency = gocql.One
	case "all":
		cluster.Consistency = gocql.All
	default:
		cluster.Consistency = gocql.Quorum
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
