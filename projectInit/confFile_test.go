package projectInit

import "testing"

func TestCreateConf(t *testing.T) {
	CreateConfYaml("/Users/derek/project/demo/gomybatis/tmp/conf")
	CreateConf("/Users/derek/project/demo/gomybatis/tmp/conf", "gomybatis")
}
