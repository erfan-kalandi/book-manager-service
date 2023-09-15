package config

type Config struct {
	Database struct {
		Host     string `env:"DATABASE_HOST" env-default:"db"`
		Port     int    `env:"DATABASE_PORT" env-default:"5432"`
		Username string `env:"DATABASE_USERNAME" env-default:"postgres"`
		Password string `env:"DATABASE_PASSWORD" env-default:"5506"`
		DBName   string `env:"DATABASE_NAME" env-default:"book_manager_db"`
	}
}
