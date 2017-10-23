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
	"math/rand"
	"os"
	"sync"
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

	if config.LogFile != "" {
		fout, err := os.OpenFile(
			config.LogFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()
		log.SetOutput(fout)
	}

	rand.Seed(time.Now().Unix())

	selection, selectionLock := DailySelection{}, sync.RWMutex{}

	if cacheFin, err := os.Open(config.TempFile); err != nil {
		log.Println("No cache file found, loading synchronously")
		SelectSynchronously(*config, &selection, &selectionLock)
	} else {
		selection, err = ReadBackup(cacheFin)
		cacheFin.Close()
		if err != nil {
			log.Println("Error reading cache file, loading synchronously")
			SelectSynchronously(*config, &selection, &selectionLock)
		} else if time.Now().After(NextLoadTime(selection.Time)) {
			log.Println("Cache file is too old, loading synchronously")
			SelectSynchronously(*config, &selection, &selectionLock)
		} else {
			log.Println("Loaded cache from", selection.Time)
		}
	}

	time.AfterFunc(
		NextLoadTime(selection.Time).Sub(time.Now()),
		func() {
			ticker := time.NewTicker(time.Hour * 24)
			go func() {
				<- ticker.C
				SelectSynchronously(*config, &selection, &selectionLock)
			}()
			SelectSynchronously(*config, &selection, &selectionLock)
		},
	)

	for {
		time.Sleep(time.Minute)
	}
}

func getConfig() (*Config, *conflag.Config) {
	config := &Config{
		Port:     8080,
		TempFile: "cache",
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
