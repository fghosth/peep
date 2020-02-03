package projectInit

import "testing"

func TestCreateCMDMainFile(t *testing.T) {
	CreateCMDMainFile("/Users/derek/project/demo/gomybatis/tmp/cmd", "gomybatis/cmd/cmd")
	CreateCMDRootFile("/Users/derek/project/demo/gomybatis/tmp/cmd/cmd")
	CreateCMDUpFile("/Users/derek/project/demo/gomybatis/tmp/cmd/cmd")
}
