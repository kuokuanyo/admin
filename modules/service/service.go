// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package service

// db.Connection(interface)ㄣΤService(interface)篈常ㄣΤNameよ猭
// Service(interface)ㄣΤdb.Connection(interface)篈常ㄣΤNameよ猭
type Service interface {
	Name() string
}

type Generator func() (Service, error)

// 盢把计kgen盢services(map[string]Generator)い
func Register(k string, gen Generator) {
	if _, ok := services[k]; ok {
		panic("service has been registered")
	}
	services[k] = gen
}

// ﹍てList(map[string]Service)Service琌interface(Nameよ猭)
func GetServices() List {
	var (
		l   = make(List)
		err error
	)

	for k, gen := range services {
		if l[k], err = gen(); err != nil {
			panic("service initialize fail")
		}
	}
	return l
}

var services = make(Generators)

type Generators map[string]Generator

type List map[string]Service

// 硓筁把计(k)眔で皌Service(interface)
func (g List) Get(k string) Service {
	if v, ok := g[k]; ok {
		return v
	}
	panic("service not found")
}

// 耞琌Τ眔才把计(k)Service
func (g List) GetOrNot(k string) (Service, bool) {
	v, ok := g[k]
	return v, ok
}

// 虑パ把计穝糤List(map[string]Service)
func (g List) Add(k string, service Service) {
	if _, ok := g[k]; ok {
		panic("service exist")
	}
	g[k] = service
}
