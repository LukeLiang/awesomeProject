package common

type Config struct {
	Server `mapstructure:"server"`
	DB     `mapstructure:"db"`
	Jwt    `mapstructure:"jwt"`
	Api    `mapstructure:"api"`
}

type DB struct {
	Type     string
	Host     string
	Port     int64
	User     string
	Password string
	Database string
}

type Server struct {
	Port int64
}

type Jwt struct {
	Secret []byte `mapstructure:"-"`
	secret string `mapstructure:"secret"`
}

// GetSecret 获取 JWT 密钥（[]byte 格式）
func (j *Jwt) SetSecret(s string) {
	j.secret = s
	j.Secret = []byte(s)
}

type Api struct {
	Version string
	Prefix  string
}
