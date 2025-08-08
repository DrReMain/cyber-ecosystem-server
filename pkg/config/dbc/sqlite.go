package dbc

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zeromicro/go-zero/core/logx"
)

func (c DBConfig) SqliteDSN() string {
	if c.DBPath != "" {
		logx.Must(errors.New("the sqlite3 filepath cannot be empty"))
	}

	if _, err := os.Stat(c.DBPath); os.IsNotExist(err) {
		if f, err := os.OpenFile(c.DBPath, os.O_CREATE|os.O_RDWR, 0600); err != nil {
			logx.Must(fmt.Errorf("failed to create SQLite database file %q", c.DBPath))
		} else if err = f.Close(); err != nil {
			logx.Must(fmt.Errorf("failed to close SQLite database file %q", c.DBPath))
		}
	} else {
		if err := os.Chmod(c.DBPath, 0660); err != nil {
			logx.Must(fmt.Errorf("unable to set permission code on %s: %v", c.DBPath, err))
		}
	}

	return fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1%s", c.DBPath, c.SqliteConfig)
}
