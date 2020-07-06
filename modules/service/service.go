// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package service

// db.Connection(interface)�]�㦳Service(interface)���A�A�]�����㦳Name��k
// Service(interface)�]�㦳db.Connection(interface)���A�A�]�����㦳Name��k
type Service interface {
	Name() string
}

type Generator func() (Service, error)

// �N�Ѽ�k�Bgen�N�Jservices(map[string]Generator)��
func Register(k string, gen Generator) {
	if _, ok := services[k]; ok {
		panic("service has been registered")
	}
	services[k] = gen
}

// ��l��List(map[string]Service)�AService�Ointerface(Name��k)
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

// �z�L�Ѽ�(k)���o�ǰt��Service(interface)
func (g List) Get(k string) Service {
	if v, ok := g[k]; ok {
		return v
	}
	panic("service not found")
}

// �P�_�O�_�����o�ŦX�Ѽ�(k)��Service
func (g List) GetOrNot(k string) (Service, bool) {
	v, ok := g[k]
	return v, ok
}

// �ǥѰѼƷs�WList(map[string]Service)
func (g List) Add(k string, service Service) {
	if _, ok := g[k]; ok {
		panic("service exist")
	}
	g[k] = service
}
