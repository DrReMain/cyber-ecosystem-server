package dbc

import (
	"fmt"

	_ "github.com/lib/pq"
)

func (c DBConfig) PostgresDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
		c.PGConfig,
	)
}
