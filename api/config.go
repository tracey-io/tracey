package api

import (
	"errors"
	"fmt"
)

var ErrEnvironmentNotSupported = errors.New("environment not supported")

type Environment string

const (
	EnvironmentDev  Environment = "dev"
	EnvironmentProd Environment = "prod"
)

type ServerConfig struct {
	Address     *ServerAddress
	Environment Environment
	SecretKey   string
}

type ServerAddress struct {
	Host string
	Port string
}

func (c *ServerAddress) String() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
