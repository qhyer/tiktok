package config

type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"` // 签名密钥
}

type System struct {
	JWTConfig *JWTConfig `mapstructrue:"jwt"`
}
