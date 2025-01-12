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

package events

// Config holds all the configuration params
type Config struct {
	Nats NatsConfig `koanf:"nats"`
}

// NatsConfig holds the configuration for the NATS server
type NatsConfig struct {
	Port       int    `koanf:"port"`
	Timeout    int    `koanf:"timeout"`
	Hostname   string `koanf:"hostname"`
	StorageDir string `koanf:"storagedir"`
}

// DefaultConfig returns an instance of Config with the default values.
func DefaultConfig() Config {
	return Config{
		Nats: NatsConfig{
			Port:     4222,
			Timeout:  30,
			Hostname: "localhost",
		},
	}
}
