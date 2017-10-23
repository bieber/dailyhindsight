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
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

const maxQueuedSelections = 1024
const selectTopN = 20

type selectionList []DailySelection

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
	results := []DailySelection{}
	fetchTime := time.Now()
	RunLimited(func(set Dataset) {
		result, err := GetRequest(config.APIKey, fetchTime, set)
		if err != nil {
			log.Println(err)
			return
		}

		results = append(results, DailySelection{set, result, fetchTime})
		log.Println(DailySelection{set, result, fetchTime})
	})

	sort.Sort(sort.Reverse(selectionList(results)))
	for _, s := range results {
		log.Printf("%s:%s %f", s.Database, s.Dataset.Dataset, s.NewValue/s.OldValue)
	}

	selectionLock.Lock()
	*selection = results[rand.Intn(selectTopN)]
	selectionLock.Unlock()

	cacheFout, err := os.Create(config.TempFile)
	defer cacheFout.Close()
	if err != nil {
		log.Println("Error opening selection cache file")
	}

	err = WriteBackup(cacheFout, *selection)
	if err != nil {
		log.Println("Error writing selection cache")
	}
}

func (l selectionList) Len() int {
	return len(l)
}

func (l selectionList) Less(i, j int) bool {
	return l[i].NewValue/l[i].OldValue < l[j].NewValue/l[j].OldValue
}

func (l selectionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
