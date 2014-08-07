package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type Environment struct {
	name      string
	variables []*EnvironmentVariable
}

type EnvironmentVariable struct {
	name  string
	value string
}

func NewEnvironment(name string, variables ...EnvironmentVariable) *Environment {
	e := Environment{name, []*EnvironmentVariable{}}
	env := os.Environ()

	for i := 0; i < len(env); i++ {
		envv := strings.Split(env[i], "=")
		e.AddEnvironmentVariable(NewEnvironmentVariable(envv[0], envv[1]))
	}

	return &e
}

func (e *Environment) GetEnvironment(hosts []string) error {
	c := etcd.NewClient(hosts)
	response, err := c.Get(e.name, true, false)
	if err != nil {
		// if len(response.Node.Nodes) > 0 {
		// 	log.Error("Get failed with %s %s %v", response.Node.Key, response.Node.Value, response.Node.TTL)
		// }

		return err
	}

	for _, node := range response.Node.Nodes {
		log.Debug("Found variable: [%s] = %s", strings.Replace(node.Key, e.KeyName(), "", 1), node.Value)
		e.AddEnvironmentVariable(NewEnvironmentVariable(strings.Replace(node.Key, e.KeyName(), "", 1), node.Value))
	}
	return nil
}

func (e *Environment) KeyName() string {
	return fmt.Sprintf("%s/", e.name)
}

func (e *Environment) AddEnvironmentVariable(ev *EnvironmentVariable) {
	e.variables = append(e.variables, ev)
}

func (e Environment) Env() []string {
	var stringArray = []string{}

	for i := 0; i < len(e.variables); i++ {
		stringArray = append(stringArray, e.variables[i].String())
	}

	return stringArray
}

func (e Environment) Update() {
	// TODO
}

func NewEnvironmentVariable(name string, value string) *EnvironmentVariable {
	ev := EnvironmentVariable{name, value}
	return &ev
}

func (ev *EnvironmentVariable) String() string {
	return ev.name + "=" + ev.value
}
