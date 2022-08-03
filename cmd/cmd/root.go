package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = BaseCommand("peep", "根据数据库生成代码")

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// BaseCommand provides the basic flags vars for running a service
func BaseCommand(serviceName, shortDescription string) *cobra.Command {
	command := &cobra.Command{
		Use:   serviceName,
		Short: shortDescription,
	}

	//command.PersistentFlags().StringVar(
	//	&service.Config.Host,
	//	"grpc-host",
	//	"0.0.0.0",
	//	"gRPC service hostname",
	//)

	return command
}
