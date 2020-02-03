package projectInit

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fghosth/peep/util"
)

var (
	cmdMainFile = "main.go"
	cmdRootFile = "root.go"
	cmdUpFile   = "up.go"
	cmdVersionFile   = "version.go"
)

func CreateCMDMainFile(path string, pName string) error {
	ctx := map[string]interface{}{
		"cmdPName": pName,
	}

	str, err := raymond.Render(CMD_MAIN_TMP, ctx)
	if err != nil {
		fmt.Println(err)
	}
	return util.WriteWithIoutil(path+"/"+cmdMainFile, str)
}

func CreateCMDRootFile(path string) error {
	return util.WriteWithIoutil(path+"/"+cmdRootFile, CMD_ROOT_TMP)
}

func CreateCMDUpFile(path string) error {
	return util.WriteWithIoutil(path+"/"+cmdUpFile, CMD_UP_TMP)
}

func CreateVersionFile(path string) error {
	return util.WriteWithIoutil(path+"/"+cmdVersionFile, CMD_VERSION_TMP)
}
