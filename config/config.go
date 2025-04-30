package config

type dbKind int

const (
	MongoDbType dbKind = iota
)

type Config struct{}

type DbConfig struct {
	DbKind      dbKind       `json:"dbKind"`
	MongoConfig *MongoConfig `json:"mongoConfig"`
}

type MongoConfig struct {
	URI            string `json:"uri"`
	Database       string `json:"database"`
	MaxPoolSize    uint64 `json:"max_pool_size"`
	MinPoolSize    uint64 `json:"min_pool_size"`
	ConnectTimeout int    `json:"connect_timeout"`
	ServerAPI      string `json:"server_api"`
}
