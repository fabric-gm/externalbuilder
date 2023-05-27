/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/5/27 09:19:48
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/5/27 09:19:48
 */

package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/platforms/util"
	"github.com/spf13/viper"
	"os"

	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
)

const staticLDFlagsOpts = "-ldflags \"-linkmode external -extldflags '-static'\""
const dynamicLDFlagsOpts = ""

var buildScript = `
set -e
if [ -f "/chaincode/input/src/go.mod" ] && [ -d "/chaincode/input/src/vendor" ]; then
    cd /chaincode/input/src
	go mod tidy
    GO111MODULE=on go build -v -mod=vendor %[1]s -o /chaincode/output/chaincode %[2]s
elif [ -f "/chaincode/input/src/go.mod" ]; then
    cd /chaincode/input/src
	go mod tidy
    GO111MODULE=on go build -v -mod=readonly %[1]s -o /chaincode/output/chaincode %[2]s
elif [ -f "/chaincode/input/src/%[2]s/go.mod" ] && [ -d "/chaincode/input/src/%[2]s/vendor" ]; then
    cd /chaincode/input/src/%[2]s
	go mod tidy
    GO111MODULE=on go build -v -mod=vendor %[1]s -o /chaincode/output/chaincode .
elif [ -f "/chaincode/input/src/%[2]s/go.mod" ]; then
    cd /chaincode/input/src/%[2]s
	go mod tidy
    GO111MODULE=on go build -v -mod=readonly %[1]s -o /chaincode/output/chaincode .
else
    GOPATH=/chaincode/input:$GOPATH go build -v %[1]s -o /chaincode/output/chaincode %[2]s
fi
echo Done!
`

type GolangPlatform struct {
	golang.Platform
}

func (p *GolangPlatform) DockerBuildOptions(path string) (util.DockerBuildOptions, error) {
	env := []string{}
	for _, key := range []string{"GOPROXY", "GOSUMDB"} {
		if val, ok := os.LookupEnv(key); ok {
			env = append(env, fmt.Sprintf("%s=%s", key, val))
			continue
		}
		if key == "GOPROXY" {
			env = append(env, "GOPROXY=https://goproxy.io")
		}
	}
	ldFlagOpts := getLDFlagsOpts()
	return util.DockerBuildOptions{
		Cmd: fmt.Sprintf(buildScript, ldFlagOpts, path),
		Env: env,
	}, nil
}

func getLDFlagsOpts() string {
	if viper.GetBool("chaincode.golang.dynamicLink") {
		return dynamicLDFlagsOpts
	}
	return staticLDFlagsOpts
}
