package etcdv3

import (
	"crypto/tls"
	"crypto/x509"
	"go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"time"
)

type TlsConfig struct {
	TlsStatus bool
	CertPath  string
	KeyPath   string
	CaPath    string
}

type Etcd struct {
	client    *clientv3.Client
	tlsConfig *TlsConfig
	endpoints []string
	namespace string
	username  string
	password  string
	ttl       time.Duration
}

var Namespace string

func (e *Etcd) Close() error {
	return e.client.Close()
}

func (e *Etcd) GetClient() *clientv3.Client {
	return e.client
}

func (e *Etcd) Reconnect() error {
	if e.tlsConfig != nil {
		cert, err := tls.LoadX509KeyPair(e.tlsConfig.CertPath, e.tlsConfig.KeyPath)
		if err != nil {
			return err
		}
		caData, err := ioutil.ReadFile(e.tlsConfig.CaPath)
		if err != nil {
			return err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		_tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
		e.client, err = clientv3.New(clientv3.Config{
			Endpoints:   e.endpoints,
			DialTimeout: e.ttl,
			TLS:         _tlsConfig,
			Username:    e.username,
			Password:    e.password,
		})
		if err != nil {
			return err
		}
	} else {
		var err error
		e.client, err = clientv3.New(clientv3.Config{
			Endpoints:   e.endpoints,
			DialTimeout: e.ttl,
			Username:    e.username,
			Password:    e.password,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewEtcd(endpoints []string, username, password, namespace string, ttl int64, tlsConfig *TlsConfig) (*Etcd, error) {
	var err error
	etcd := &Etcd{}
	etcd.tlsConfig = tlsConfig
	etcd.ttl = time.Duration(ttl) * time.Second
	etcd.endpoints = endpoints
	etcd.namespace = namespace
	etcd.username = username
	etcd.password = password
	if etcd.tlsConfig != nil {
		cert, err := tls.LoadX509KeyPair(etcd.tlsConfig.CertPath, etcd.tlsConfig.KeyPath)
		if err != nil {
			return nil, err
		}
		caData, err := ioutil.ReadFile(etcd.tlsConfig.CaPath)
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		_tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
		etcd.client, err = clientv3.New(clientv3.Config{
			Endpoints:   etcd.endpoints,
			DialTimeout: etcd.ttl,
			TLS:         _tlsConfig,
			Username:    etcd.username,
			Password:    etcd.password,
		})
		if err != nil {
			return nil, err
		}
	} else {
		etcd.client, err = clientv3.New(clientv3.Config{
			Endpoints:   etcd.endpoints,
			DialTimeout: etcd.ttl,
			Username:    etcd.username,
			Password:    etcd.password,
		})
		if err != nil {
			return nil, err
		}
	}
	return etcd, nil
}
