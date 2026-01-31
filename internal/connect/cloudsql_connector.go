package connect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type CloudSQLOptions struct {
	Instance string
	User     string
	Password string
	Database string
	SSLMode  string
	IAMAuthN bool
	IPType   string
}

func OpenCloudSQL(ctx context.Context, opts CloudSQLOptions) (*sql.DB, func(), error) {
	if opts.Instance == "" {
		return nil, nil, errors.New("cloud sql instance is required")
	}
	if opts.User == "" || opts.Database == "" {
		return nil, nil, errors.New("user and dbname are required for cloud sql")
	}

	var dialerOptions []cloudsqlconn.Option
	if opts.IAMAuthN {
		dialerOptions = append(dialerOptions, cloudsqlconn.WithIAMAuthN())
	}

	var dialOptions []cloudsqlconn.DialOption
	switch opts.IPType {
	case "", "public":
	case "private":
		dialOptions = append(dialOptions, cloudsqlconn.WithPrivateIP())
	default:
		return nil, nil, fmt.Errorf("unsupported cloudsql-ip-type: %s", opts.IPType)
	}

	dialer, err := cloudsqlconn.NewDialer(ctx, dialerOptions...)
	if err != nil {
		return nil, nil, err
	}

	dsn, err := buildLocalDSN(opts.User, opts.Password, opts.Database, opts.SSLMode)
	if err != nil {
		_ = dialer.Close()
		return nil, nil, err
	}

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		_ = dialer.Close()
		return nil, nil, err
	}

	connConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(ctx, opts.Instance, dialOptions...)
	}

	db := stdlib.OpenDB(*connConfig)
	cleanup := func() {
		_ = dialer.Close()
	}
	return db, cleanup, nil
}

func buildLocalDSN(user, password, database, sslmode string) (string, error) {
	if user == "" || database == "" {
		return "", errors.New("user and dbname are required")
	}
	cred := url.User(user)
	if password != "" {
		cred = url.UserPassword(user, password)
	}
	query := url.Values{}
	if sslmode != "" {
		query.Set("sslmode", sslmode)
	}
	return (&url.URL{
		Scheme:   "postgres",
		User:     cred,
		Host:     "localhost:5432",
		Path:     "/" + database,
		RawQuery: query.Encode(),
	}).String(), nil
}
