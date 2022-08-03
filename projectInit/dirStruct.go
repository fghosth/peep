package projectInit

import (
	"fmt"
	"github.com/fghosth/peep/util"
	"os"
)

//生成目录
func CreateDirStruct(path string) {
	root := util.NewTree("root")
	root.Add(".githooks")
	root.Add("application")
	cmd := root.Add("cmd")
	cmd.Add("cmd")
	root.Add("conf")
	root.Add("bctx")
	root.Add("deploy")
	root.Add("docs")
	domain := root.Add("domain")
	domain.Add("role")
	domain.Add("object")
	infra := root.Add("infra")
	infra.Add("util")
	infra.Add("mysql")
	infra.Add("grpc_mw")
	infra.Add("swagger")
	infra.Add("dssh")
	infra.Add("apollo")
	infra.Add("apolloOpenApi")
	infra.Add("filex")
	ui := root.Add("ui")
	ui.Add("grpc")
	ui.Add("client")

	mkdirDir(root, path)
}

func mkdirDir(t util.Tree, path string) {
	for _, v := range t.Items() {
		p := path + "/" + v.Text()
		err := os.Mkdir(p, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		if len(v.Items()) > 0 {
			mkdirDir(v, p)
		}

	}
}
