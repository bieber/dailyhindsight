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

// How far into the day to refresh (using an offset to make sure we
// don't have any doubt about what day it is and DST doesn't ruin
// everything)
const dayStartOffset = time.Hour * 2

func NextLoadTime(t time.Time) time.Time {
	// I know, strictly speaking, this has some edge cases.  Hopefully
	// none of them blow up horifically
	someTimeTomorrow := t.AddDate(0, 0, 1)
	year, month, day := someTimeTomorrow.Date()
	return time.
		Date(year, month, day, 0, 0, 0, 0, t.Location()).
		Add(dayStartOffset)
}
