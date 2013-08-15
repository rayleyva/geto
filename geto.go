/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"github.com/bgmerrell/geto/config"
	"github.com/bgmerrell/geto/server"
)

/* Variables set by command line parsing */
var configPath string

func parseCommandLine() {
	/* XXX: How to make this portable? */
	flag.StringVar(&configPath, "config-path", "/etc/geto.ini", "Configuration file path")
	flag.Parse()
}

func main() {
	parseCommandLine()
	if !config.Config(configPath) {
		return
	}
	server.Serve()
}