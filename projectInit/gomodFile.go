package projectInit

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fghosth/peep/util"
)

var (
	gomodFile = "go.mod"
)

func CreateGomodFile(name string, module string) error {

	ctx := map[string]interface{}{
		"module": module,
	}

	str, err := raymond.Render(GOMOD_TMP, ctx)
	if err != nil {
		fmt.Println(err)
	}
	return util.WriteWithIoutil(name+"/"+gomodFile, str)
}
