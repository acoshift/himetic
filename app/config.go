package app

import (
	"database/sql"

	"github.com/acoshift/session"
)

// Config is the app's config
type Config struct {
	SessionStorage session.Store
	SessionSecret  []byte
	SessionName    string
	DB             *sql.DB
}
