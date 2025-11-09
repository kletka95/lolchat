package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration interface {
	GeneralParse() string
}
type generalConfig struct {
	port string
	host string
}
type config struct {
	gen *generalConfig
}

func Load(flname ...string) (*config, error) {
	err := godotenv.Load(flname...)
	if err != nil {
		return nil, err
	}
	return &config{
		gen: &generalConfig{
			host: os.Getenv("HOST"),
			port: os.Getenv("PORT"),
		},
	}, nil
}
func (c *config) GeneralParse() string {
	return fmt.Sprintf("%s:%s", c.gen.host, c.gen.port)
}
