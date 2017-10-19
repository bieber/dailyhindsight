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
	"time"
)

func RunLimited(f func(d Dataset)) {
	log := make([]time.Time, 0, len(Datasets))

	canContinue := func() bool {
		countPerLimit := make(map[time.Duration]int, len(Limits))
		for k := range Limits {
			countPerLimit[k] = 0
		}

		t := time.Now()
		for _, v := range log {
			delta := t.Sub(v)
			for k := range countPerLimit {
				if delta <= k {
					countPerLimit[k]++
				}
			}
		}

		for k, limit := range Limits {
			if countPerLimit[k] >= limit {
				return false
			}
		}

		return true
	}

	for _, v := range Datasets {
		for !canContinue() {
			time.Sleep(time.Second)
		}

		f(v)
		log = append(log, time.Now())
	}
}
