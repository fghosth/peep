package cmd

import (
	"fmt"
	"github.com/fghosth/peep/projectInit"
	"github.com/spf13/cobra"
	"os/exec"
)

var (
	module      string
	projectPath string
)
var createProjectCmd = &cobra.Command{
	Use:     "pcreate",
	Short:   "create project struct",
	Example: "./c pcreate --module \"fghosth.net/reportssss\" --path /Users/derek/project/demo/gomybatis/tmp",
	Run: func(cmd *cobra.Command, args []string) {
		projectInit.CreateDirStruct(projectPath)
		projectInit.CreateGitIgnore(projectPath)
		projectInit.CreateGomodFile(projectPath, module)
		projectInit.CreateDockerfile(projectPath)
		projectInit.CreateReadme(projectPath)
		projectInit.CreateMakefile(projectPath, module, projectPath)
		projectInit.CreateEntrypoint(projectPath)
		projectInit.CreateCMDUpFile(projectPath+"/cmd/cmd", module)
		projectInit.CreateCMDRootFile(projectPath + "/cmd/cmd")
		projectInit.CreateVersionFile(projectPath + "/cmd/cmd")
		projectInit.CreateCMDMainFile(projectPath+"/cmd", module+"/cmd/cmd")
		projectInit.CreateConfYaml(projectPath + "/conf")
		projectInit.CreateConf(projectPath+"/conf", module)
		projectInit.CreateLogfile(projectPath + "/infra/util")
		projectInit.CreateOpentraceFile(projectPath + "/infra/grpc_mw")
		projectInit.CreateRateLimitFile(projectPath + "/infra/grpc_mw")
		projectInit.CreateUtilUUIDfile(projectPath + "/infra/util")
		projectInit.CreateDsshFile(projectPath + "/infra/dssh")
		projectInit.CreateApolloFile(projectPath+"/infra/apollo", module)
		projectInit.CreateFilexFile(projectPath + "/infra/filex")
		projectInit.CreateApolloOpenApiFile(projectPath + "/infra/apolloOpenApi")
		projectInit.CreateGithookGoCommit(projectPath + "/.githooks")
		projectInit.CreateGithookPreCommit(projectPath + "/.githooks")
		projectInit.CreateGolangCi(projectPath + "/.githooks")
		projectInit.CreateCommitMsg(projectPath + "/.githooks")

		gocmd := exec.Command("gofmt", "-w", "-s", projectPath)

		err := gocmd.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
			return
		}
	},
}

func init() {
	createProjectCmd.PersistentFlags().StringVar(
		&module,
		"module",
		"test",
		"model path",
	)
	createProjectCmd.PersistentFlags().StringVar(
		&projectPath,
		"path",
		"./",
		"mysql uri",
	)
	RootCmd.AddCommand(createProjectCmd)
}
