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
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	APIKey   string
	TempFile string
}

func main() {
	viper.SetDefault("port", 80)
	viper.SetDefault("temp_file", "cache")

	viper.BindEnv("port")
	viper.BindEnv("api_key")
	viper.BindEnv("temp_file")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/run/secrets")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Couldn't load config file: %s", err.Error())
	}

	config := Config{
		APIKey:   viper.GetString("api_key"),
		TempFile: viper.GetString("temp_file"),
	}

	rand.Seed(time.Now().Unix())

	selection, selectionLock := DailySelection{}, sync.RWMutex{}

	if cacheFin, err := os.Open(config.TempFile); err != nil {
		log.Println("No cache file found, loading synchronously")
		SelectSynchronously(config, &selection, &selectionLock)
	} else {
		selection, err = ReadBackup(cacheFin)
		cacheFin.Close()
		if err != nil {
			log.Println("Error reading cache file, loading synchronously")
			SelectSynchronously(config, &selection, &selectionLock)
		} else if time.Now().After(NextLoadTime(selection.Time)) {
			log.Println("Cache file is too old, loading synchronously")
			SelectSynchronously(config, &selection, &selectionLock)
		} else {
			log.Println("Loaded cache from", selection.Time)
		}
	}

	time.AfterFunc(
		NextLoadTime(selection.Time).Sub(time.Now()),
		func() {
			ticker := time.NewTicker(time.Hour * 24)
			go func() {
				for {
					<-ticker.C
					SelectSynchronously(config, &selection, &selectionLock)
				}
			}()
			SelectSynchronously(config, &selection, &selectionLock)
		},
	)

	log.Printf("Starting server on port %d\n", viper.GetInt("port"))
	http.Handle(
		"/",
		Middleware(IndexHandler(config, &selection, &selectionLock)),
	)
	http.Handle("/favicon.ico", Middleware(FaviconHandler()))
	log.Fatal(
		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil),
	)
}
