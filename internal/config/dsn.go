package config

import (
	"fmt"
	"net/url"
)

func BuildDSN(cfg Config) (string, error) {
	if cfg.Host == "" || cfg.User == "" || cfg.Database == "" {
		return "", fmt.Errorf("host, user, and dbname are required to build a DSN")
	}

	user := url.User(cfg.User)
	if cfg.Password != "" {
		user = url.UserPassword(cfg.User, cfg.Password)
	}

	port := cfg.Port
	if port == 0 {
		port = 5432
	}

	endpoint := fmt.Sprintf("%s:%d", cfg.Host, port)

	u := url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   endpoint,
		Path:   "/" + cfg.Database,
	}

	query := url.Values{}
	if cfg.SSLMode != "" {
		query.Set("sslmode", cfg.SSLMode)
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}
