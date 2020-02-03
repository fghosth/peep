package projectInit

import (
	"github.com/k0kubun/pp"
	"testing"
)

func TestCreateDockerfile(t *testing.T) {
	err := CreateDockerfile("/Users/derek/project/demo/gomybatis/tmp")
	CreateEntrypoint("/Users/derek/project/demo/gomybatis/tmp")
	CreateReadme("/Users/derek/project/demo/gomybatis/tmp")
	CreateMakefile("/Users/derek/project/demo/gomybatis/tmp")
	pp.Println(err)
}
