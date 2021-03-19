/*
 * Nuts node
 * Copyright (C) 2021 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package core

import (
	"github.com/spf13/pflag"
)

const defaultConfigFile = "nuts.yaml"
const configFileFlag = "configfile"
const addressFlag = "address"
const datadirFlag = "datadir"
const defaultAddress = "localhost:1323"
const strictModeFlag = "strictmode"
const defaultStrictMode = false
const defaultDatadir = "./data"
const defaultLogLevel = "info"
const loggerLevelFlag = "verbosity"

// GlobalHTTPConfig is the top-level config struct for HTTP interfaces.
type GlobalHTTPConfig struct {
	// HTTPConfig contains the config for the default HTTP interface.
	HTTPConfig `koanf:"default"`
	// AltBinds contains binds for alternative HTTP interfaces. The key of the map is the first part of the path
	// of the URL (e.g. `/internal/some-api` -> `internal`), the value is the HTTP interface it must be bound to.
	AltBinds map[string]HTTPConfig `koanf:"alt"`
}

// HTTPConfig contains configuration for an HTTP interface, e.g. address.
// It will probably contain security related properties in the future (TLS configuration, user/pwd requirements).
type HTTPConfig struct {
	// Address holds the interface address the HTTP service must be bound to, in the format of `interface:port` (e.g. localhost:5555).
	Address string `koanf:"address"`
}

// FlagSet returns the default server flags
func FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("server", pflag.ContinueOnError)
	flagSet.String(configFileFlag, defaultConfigFile, "Nuts config file")
	flagSet.String(loggerLevelFlag, defaultLogLevel, "Log level (trace, debug, info, warn, error)")
	flagSet.String(addressFlag, defaultAddress, "Address and port the server will be listening to")
	flagSet.Bool(strictModeFlag, defaultStrictMode, "When set, insecure settings are forbidden.")
	flagSet.String(datadirFlag, defaultDatadir, "Directory where the node stores its files.")
	return flagSet
}