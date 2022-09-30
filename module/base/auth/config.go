package auth

type Config struct {
	JwtKey string       `mapstructure:"jwt_key"`
	Admin  AdminConfig  `mapstructure:"admin"`
	Server ServerConfig `mapstructure:"server"`
	Client ClientConfig `mapstructure:"client"`
}

type ServerConfig struct {
	Enable    bool   `mapstructure:"enable"`     //是否需要鉴权
	AccessKey string `mapstructure:"access_key"` //AK
	SecretKey string `mapstructure:"secret_key"` //SK
}

type ClientConfig struct {
	Enable bool `mapstructure:"enable"` //是否需要鉴权
}

type AdminConfig struct {
	Enable bool `mapstructure:"enable"` //是否需要鉴权
}
