package projectInit

import "github.com/fghosth/peep/util"

var (
	logfile = "log.go"
)

func CreateLogfile(name string) error {
	return util.WriteWithIoutil(name+"/"+logfile, LOG_TMP)
}
