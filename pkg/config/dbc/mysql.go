package dbc

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (c DBConfig) MysqlDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=True%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.MysqlConfig,
	)
}
