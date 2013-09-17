/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Configuration file management
*/
package config

import (
	"errors"
	"github.com/robfig/config"
	"log"
	"os"
)

var conf Config
var isParsed bool = false

func init() {
	conf = Config{}
}

type Host struct {
	name string
	addr string
	username string
	password string
	validPassword bool
}

type Config struct {
	FilePath string
	PrivKeyPath string
	Hosts    []Host
}

// Parse the config file
// The configPath parameter is the path to the config file on the filesystem
func ParseConfig(configPath string) (Config, error) {
	conf = Config{}  /* Zero out the package-scope conf */
	var err error
	var hasPrivKey bool

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		log.Print("No configuration file: ", configPath)
		return conf, err
	}
	log.Print("Parsing configuration file: ", configPath)

	var c *config.Config
	c, err = config.ReadDefault(configPath)
	if err != nil {
		log.Print("Failed to parse config file: ", err.Error())
		return conf, err
	}

	var privKeyPath string
	if privKeyPath, err = c.String("geto", "privkey_path"); err == nil {
		conf.PrivKeyPath = privKeyPath
		hasPrivKey = true
	}

	var opts []string
	if opts, err = c.Options("hosts"); err != nil {
		log.Print("Could not find \"hosts\" section: ", err.Error())
		return conf, err
	}

	const N_MIN_REQUIRED_HOSTS = 1
	if len(opts) < N_MIN_REQUIRED_HOSTS {
		err = errors.New("Config must have at least one host")
		log.Print("Failed to parse \"hosts\" section: ", err.Error())
		return conf, err
	}

	for _, hostname := range opts {
		var addr, username, password string
		var validPassword bool
		if addr, err = c.String("hosts", hostname); err != nil {
			log.Print("Failed to parse \"hosts\" section: ", err.Error())
			return conf, err
		}
		if opts, err = c.Options(hostname); err != nil {
			log.Print("Could not find \"", hostname, "\" section: ", err.Error())
			return conf, err
		}
		if username, err = c.String(hostname, "username"); err != nil {
			log.Print("Failed to parse \"username\" option for \"",
				hostname, "\" section: ", err.Error())
			return conf, err
		}
		if password, err = c.String(hostname, "password"); err == nil {
			validPassword = true
		} else {
			if !hasPrivKey {
				log.Print("Failed to parse \"username\" option for \"",
					hostname, "\" section: ", err.Error())
				return conf, err
			}
		}
		conf.Hosts = append(
			conf.Hosts,
			Host{hostname, addr, username, password, validPassword})
	}

	conf.FilePath = configPath
	isParsed = true
	return conf, nil
}

// Return the Config object.
// ParseConfig should probably be called before this function
func GetConfig() Config {
	if !isParsed {
		log.Println("Warning: unparsed configuration")
	}
	return conf
}