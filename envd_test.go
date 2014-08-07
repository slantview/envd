package main_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/coreos/go-etcd/etcd"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type EnvdSuite struct{}

var _ = Suite(&EnvdSuite{})

func (s *EnvdSuite) SetUpTest(c *C) {
	client := etcd.NewClient([]string{"http://127.0.0.1:4001/"})
	client.DeleteDir("/environments/test")
	client.CreateDir("/environments/test", 0)
	client.Set("/environments/test/VARIABLE1", "envd_var1", 0)
	client.Set("/environments/test/VARIABLE2", "envd_var2", 0)
}

func (s *EnvdSuite) Test_Environment(c *C) {
	cmd := exec.Command("./bin/envd", "-e", "/environments/test", "env")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	c.Log("%s", err)

	c.Assert(err, IsNil)
	c.Assert(strings.Replace(out.String(), "\n", " ", -1), Matches, ".*VARIABLE1=envd_var1.*")
	c.Assert(strings.Replace(out.String(), "\n", " ", -1), Matches, ".*VARIABLE2=envd_var2.*")
}
