package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`

	Port int `mapstructure:"port"`

}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo JWTConfig `mapstructure:"jwt"`
	RedisInfo RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`

	Port int `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}
