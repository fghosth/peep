package projectInit

import (
	"github.com/fghosth/peep/util"
)

var (
	opentraceServerFile = "opentracingServer.go"
	opentraceClientFile = "opentracingClient.go"

	limiterFile   = "limit.go"
	ratelimitFile = "ratelimit.go"
)

func CreateOpentraceFile(path string) error {
	err := util.WriteWithIoutil(path+"/"+opentraceClientFile, GRPC_OPENTRACE_CLIENT_TMP)
	if err != nil {
		return err
	}
	err = util.WriteWithIoutil(path+"/"+opentraceServerFile, GRPC_OPENTRACE_SERVER_TMP)
	if err != nil {
		return err
	}
	return nil
}

func CreateRateLimitFile(path string) error {
	err := util.WriteWithIoutil(path+"/"+limiterFile, GRPC_LIMITER_TMP)
	if err != nil {
		return err
	}
	err = util.WriteWithIoutil(path+"/"+ratelimitFile, GRPC_RATE_LIMITER_TMP)
	if err != nil {
		return err
	}
	return nil
}
