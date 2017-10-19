/*
 * Copyright 2017, Robert Bieber
 *
 * This file is part of dailyhindsight.
 *
 * dailyhindsight is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * dailyhindsight is distributed in the hope that it will be useful,
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with dailyhindsight.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"github.com/bieber/conflag"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"time"
)

type Config struct {
	Help     bool
	Port     int
	APIKey   string
	TempFile string
	LogFile  string
}

func main() {
	config, parser := getConfig()
	_, err := parser.Read()
	if err != nil || config.Help {
		exitCode := 0

		if err != nil {
			log.Println(err)
			exitCode = 1
		}

		if width, _, err := terminal.GetSize(0); err == nil {
			fmt.Println(parser.Usage(uint(width)))
		}
		os.Exit(exitCode)
	}

	i := 0
	RunLimited(func(v Dataset) {
		fmt.Println(i)
		i++

		_, err := GetRequest(config.APIKey, time.Now(), v)
		if err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println(GetRequest(config.APIKey, time.Now(), Datasets[0]))
}

func getConfig() (*Config, *conflag.Config) {
	config := &Config{
		Port:     8080,
		TempFile: "cache.txt",
	}

	parser, err := conflag.New(config)
	if err != nil {
		log.Fatal(err)
	}

	parser.ProgramName("dailyhindsight")
	parser.ProgramDescription("HTTP server for Daily Hindsight")
	parser.ConfigFileLongFlag("config")

	parser.Field("Help").
		ShortFlag('h').
		Description("Print usage text and exit.")

	parser.Field("Port").
		ShortFlag('p').
		Description("Port to serve HTTP traffic on.")

	parser.Field("APIKey").
		ShortFlag('k').
		FileKey("api_key").
		Required().
		Description("API key for Quandl.")

	parser.Field("TempFile").
		ShortFlag('t').
		Description("File to cache results for day in case server goes down.")

	parser.Field("LogFile").
		ShortFlag('l').
		Description("Optional log output file (logs go to stderr by default)")

	return config, parser
}
