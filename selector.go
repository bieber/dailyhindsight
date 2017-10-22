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
	"log"
	"os"
	"sync"
	"time"
)

const maxQueuedSelections = 1024

func StartSelector(
	config Config,
	selection *DailySelection,
	selectionLock *sync.RWMutex,
) chan<- bool {
	c := make(chan bool, maxQueuedSelections)
	go func() {
		for {
			<-c
			SelectSynchronously(config, selection, selectionLock)
		}
	}()
	return c
}

func SelectSynchronously(
	config Config,
	selection *DailySelection,
	selectionLock *sync.RWMutex,
) {
	selectionLock.Lock()
	defer selectionLock.Unlock()

	set := Datasets[0]
	result, _ := GetRequest(config.APIKey, time.Now(), set)
	*selection = DailySelection{set, result, time.Now()}

	cacheFout, err := os.Create(config.TempFile)
	if err != nil {
		log.Println("Error opening selection cache file")
	}

	err = WriteBackup(cacheFout, *selection)
	if err != nil {
		log.Println("Error writing selection cache")
	}
}
