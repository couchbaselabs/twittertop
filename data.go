package main

import (
	"github.com/dustin/gomemcached"
	"github.com/dustin/gomemcached/client"
)

type persister interface {
	CAS(k string, f memcached.CasFunc,
		initexp int) (rv *gomemcached.MCResponse, err error)
	Incr(key string,
		amt, def uint64, exp int) (uint64, error)
	Get(key string) (*gomemcached.MCResponse, error)
	GetBulk(keys []string) (map[string]*gomemcached.MCResponse, error)

	Close() error
}

type mcAdaptor struct {
	mc *memcached.Client
}

func (m mcAdaptor) CAS(k string, f memcached.CasFunc,
	initexp int) (rv *gomemcached.MCResponse, err error) {

	return m.mc.CAS(0, k, f, initexp)
}

func (m mcAdaptor) Incr(key string,
	amt, def uint64, exp int) (uint64, error) {

	return m.mc.Incr(0, key, amt, def, exp)
}

func (m mcAdaptor) Get(key string) (*gomemcached.MCResponse, error) {
	return m.mc.Get(0, key)
}

func (m mcAdaptor) GetBulk(keys []string) (map[string]*gomemcached.MCResponse, error) {
	return m.mc.GetBulk(0, keys)
}

func (m mcAdaptor) Close() error {
	m.mc.Close()
	return nil
}

func getPersister(u string) (persister, error) {
	client, err := memcached.Connect("tcp", u)
	if err != nil {
		return nil, err
	}
	return &mcAdaptor{client}, nil
}
