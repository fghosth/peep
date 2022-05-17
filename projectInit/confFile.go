package projectInit

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fghosth/peep/util"
)

var (
	confYamlFile = "config.yaml"
	confFile     = "conf.go"
)

func CreateConfYaml(name string) error {
	return util.WriteWithIoutil(name+"/"+confYamlFile, CONF_YAML_TMP)
}
func CreateConf(name string, module string) error {
	ctx := map[string]interface{}{
		"module": module,
	}

	confstr, err := raymond.Render(CONF_TMP, ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return util.WriteWithIoutil(name+"/"+confFile, confstr)
}
