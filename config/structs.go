package config

type (
	Config struct {
		AboutMessage  string
		Bot           Bot
		Database      Database
		ServerCounter ServerCounter
		Metrics       Metrics
		Redis         Redis
	}

	Bot struct {
		Token   string
		Prefix  string
		Admins  []string
		Helpers []string
		Shading Sharding
	}

	Sharding struct {
		Total  int
		Lowest int
		Max    int
	}

	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
		Pool     Pool
	}

	Pool struct {
		MaxConnections int
		MaxIdle        int
	}

	ServerCounter struct {
		Enabled bool
		BaseUrl string
		Key     string
	}

	Metrics struct {
		Statsd Statsd
	}

	Statsd struct {
		Enabled bool
		Prefix  string
		Host    string
		Port    int
	}

	Redis struct {
		Enabled bool
		Uri     string
		Threads int
	}
)
