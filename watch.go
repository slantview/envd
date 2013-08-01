package main

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

func NewWatch() {
	result, err := etcd.Watch(envName, 0, nil, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println(result)
}