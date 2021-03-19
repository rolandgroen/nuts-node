// +build wasm

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
	"github.com/spf13/cobra"
)

// ServerConfig has global server settings.
type ServerConfig struct {
	Verbosity  string           `koanf:"verbosity"`
	Strictmode bool             `koanf:"strictmode"`
	Datadir    string           `koanf:"datadir"`
	HTTP       GlobalHTTPConfig `koanf:"http"`
}


// NewServerConfig creates a new config with some defaults
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Verbosity:  defaultLogLevel,
		Strictmode: defaultStrictMode,
		Datadir:    defaultDatadir,
		HTTP: GlobalHTTPConfig{
			HTTPConfig: HTTPConfig{Address: defaultAddress},
			AltBinds:   map[string]HTTPConfig{},
		},
	}
}

// Load follows the load order of configfile, env vars and then commandline param
func (ngc *ServerConfig) Load(cmd *cobra.Command) (err error) {
	// TODO
	return nil
}

// PrintConfig return the current config in string form
func (ngc *ServerConfig) PrintConfig() string {
	return "TODO"
}

// InjectIntoEngine takes the loaded config and sets the engine's config struct
func (ngc *ServerConfig) InjectIntoEngine(e Injectable) error {
	// TODO
	return nil
}
