package dbc

import (
	"context"
	"database/sql"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

type DBConfig struct {
	Debug        bool   `json:",env=DATABASE_DEBUG"`
	Host         string `json:",env=DATABASE_HOST"`
	Port         int    `json:",env=DATABASE_PORT"`
	Username     string `json:",default=root,env=DATABASE_USERNAME"`
	Password     string `json:",optional,env=DATABASE_PASSWORD"`
	DBName       string `json:",default=simple_admin,env=DATABASE_DBNAME"`
	SSLMode      string `json:",optional,env=DATABASE_SSL_MODE"`
	Type         string `json:",default=mysql,options=[mysql,postgres,sqlite3],env=DATABASE_TYPE"`
	MaxOpenConn  int    `json:",optional,default=100,env=DATABASE_MAX_OPEN_CONN"`
	CacheTime    int    `json:",optional,default=10,env=DATABASE_CACHE_TIME"`
	DBPath       string `json:",optional,env=DATABASE_DBPATH"`
	MysqlConfig  string `json:",optional,env=DATABASE_MYSQL_CONFIG"`
	PGConfig     string `json:",optional,env=DATABASE_PG_CONFIG"`
	SqliteConfig string `json:",optional,env=DATABASE_SQLITE_CONFIG"`
}

func (c DBConfig) NewNoCacheDriver() *entsql.Driver {
	dsn := c.GetDSN()
	if dsn == "" {
		return nil
	}

	db, err := sql.Open(c.Type, dsn)
	if err != nil {
		logx.Must(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logx.Must(err)
	}

	db.SetMaxOpenConns(c.MaxOpenConn)
	return entsql.OpenDB(c.Type, db)
}

func (c DBConfig) GetDSN() string {
	switch c.Type {
	case "mysql":
		return c.MysqlDSN()
	case "postgres":
		return c.PostgresDSN()
	case "sqlite3":
		return c.SqliteDSN()
	default:
		return ""
	}
}
