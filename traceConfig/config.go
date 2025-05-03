package traceConfig

type dbKind int

const (
	MongoDbType dbKind = iota
)

type TraceConfig struct {
	Env         string       `json:"env"`
	DbConfig    *DbConfig    `json:"dbConfig"`
	DebugConfig *DebugConfig `json:"debugConfig"`
}

type DebugConfig struct {
	TraceId string `json:"traceId"`
}

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
