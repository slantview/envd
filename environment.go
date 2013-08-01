package main

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
	"os"
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
	envMap := make(map[string]string)

	for i := 0; i < len(env); i++ {
		eTmp := strings.Split(env[i], "=")
		envMap[eTmp[0]] = eTmp[1]
	}
	
	e.AddEnvironmentVariable(NewEnvironmentVariable("PATH", envMap["PATH"]))
	e.AddEnvironmentVariable(NewEnvironmentVariable("SHELL", envMap["SHELL"]))
	e.AddEnvironmentVariable(NewEnvironmentVariable("HOME", envMap["HOME"]))
	e.AddEnvironmentVariable(NewEnvironmentVariable("USER", envMap["USER"]))
	e.AddEnvironmentVariable(NewEnvironmentVariable("LANG", envMap["LANG"]))
	e.AddEnvironmentVariable(NewEnvironmentVariable("PWD", envMap["PWD"]))

	return &e
}

func (e *Environment) GetEnvironment() error {
	results, err := etcd.Get(e.name)
	if err != nil {
		if len(results) > 0 {
			log.Printf("Get failed with %s %s %v", results[0].Key, results[0].Value, results[0].TTL)	
		}
		
		return err
	}

	for i := 0; i < len(results); i++ {
		e.AddEnvironmentVariable(NewEnvironmentVariable(strings.Replace(results[i].Key, e.KeyName(), "", 1), results[i].Value))
	}
	return nil
}

func (e *Environment) KeyName() string {
	return fmt.Sprintf("/%s/", e.name)
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

}

func NewEnvironmentVariable(name string, value string) *EnvironmentVariable {
	ev := EnvironmentVariable{name, value}
	return &ev
}

func (ev *EnvironmentVariable) String() string {
	return ev.name + "=" + ev.value
}
