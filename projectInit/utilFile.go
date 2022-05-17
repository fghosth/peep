package projectInit

import "github.com/fghosth/peep/util"

var (
	uuidfile = "uuid.go"
)

func CreateUtilUUIDfile(path string) error {
	return util.WriteWithIoutil(path+"/"+uuidfile, UTIL_UUID_TMP)
}
