// Copyright 2025 s3-mcp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/diabloneo/s3-mcp/pkg/common"
	"github.com/diabloneo/s3-mcp/pkg/errors"
	"github.com/diabloneo/s3-mcp/pkg/log"
)

var (
	debug bool
)

func main() {
	app := &cli.App{
		Name:  "s3-mcp",
		Usage: "s3-mcp",
		Before: func(_ *cli.Context) error {
			level := logrus.InfoLevel
			if debug {
				level = logrus.DebugLevel
			}
			log.InitLogger(level)
			return nil
		},
		Commands: []*cli.Command{},
		Action: func(ctx *cli.Context) error {
			return cli.ShowAppHelp(ctx)
		},
		ExitErrHandler: func(_ *cli.Context, err error) {
			if err != nil {
				logrus.Debugf("Command failed stacktrace: %s", errors.StackTrace(err))
			}
			cli.HandleExitCoder(err)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug mode",
				Destination: &debug,
			},
		},
		Version: fmt.Sprintf(`%s
Git SHA: %s
Build At: %s
Go Version: %s
Go OS/Arch: %s/%s`,
			common.Version,
			common.GitSha[:7],
			common.BuildTime,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH),
	}

	if err := app.Run(os.Args); err != nil {
		log.Logger.Errorf("s3-mcp server exit error: %s", errors.StackTrace(err))
	}
}
