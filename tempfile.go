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
	"encoding/gob"
	"time"
	"io"
)

type DailySelection struct {
	Dataset
	RequestResult
	Time time.Time
}

func WriteBackup(fout io.Writer, selection DailySelection) error {
	encoder := gob.NewEncoder(fout)
	return encoder.Encode(selection)
}

func ReadBackup(fin io.Reader) (DailySelection, error) {
	out := DailySelection{}
	decoder := gob.NewDecoder(fin)
	err := decoder.Decode(&out)
	if err != nil {
		return DailySelection {}, err
	}
	return out, nil
}
