package mysql

import (
	"fmt"
	"github.com/gocraft/dbr"
	"k8s.io/klog"
)

type Client struct {
	database *Database
}

func NewMySQLClient(options *Options, stopCh <-chan struct{}) (*Client, error) {
	var m Client

	conn, err := dbr.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/kube-operator?parseTime=1&multiStatements=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", options.Username, options.Password, options.Host), nil)
	if err != nil {
		klog.Error("unable to connect to mysql", err)
		return nil, err
	}

	conn.SetMaxIdleConns(options.MaxIdleConnections)
	conn.SetConnMaxLifetime(options.MaxConnectionLifeTime)
	conn.SetMaxOpenConns(options.MaxOpenConnections)

	m.database = &Database{
		Session: conn.NewSession(nil),
	}

	go func() {
		<-stopCh
		if err := conn.Close(); err != nil {
			klog.Warning("error happened during closing mysql connection", err)
		}
	}()

	return &m, nil
}

func NewMySQLClientOrDie(options *Options, stopCh <-chan struct{}) *Client {
	var m Client

	conn, err := dbr.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/kube-operator?parseTime=1&multiStatements=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", options.Username, options.Password, options.Host), nil)
	if err != nil {
		klog.Error("unable to connect to mysql", err)
		panic(err)
	}

	conn.SetMaxIdleConns(options.MaxIdleConnections)
	conn.SetConnMaxLifetime(options.MaxConnectionLifeTime)
	conn.SetMaxOpenConns(options.MaxOpenConnections)

	m.database = &Database{
		Session: conn.NewSession(nil),
	}

	go func() {
		<-stopCh
		if err := conn.Close(); err != nil {
			klog.Warning("error happened during closing mysql connection", err)
		}
	}()

	return &m
}

func (m *Client) Database() *Database {
	return m.database
}
