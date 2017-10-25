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
	"github.com/sebest/xff"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var indexTemplate *template.Template

func init() {
	var err error
	indexTemplate, err = template.New("index").Parse(
		`<!DOCTYPE HTML>
<html>
	<head>
		<title>Daily Hindsight</title>
		<style>
		body {
			font-size: 25px;
		}

		div.container {
			width: 50%;
			margin-left: auto;
			margin-right: auto;
		}

		h1 {
			text-align: center;
		}

		p.top {
			text-align: center;
		}
		</style>
	</head>

	<body>
		<div class="container">
			<h1>Daily Hindsight</h1>
			<p class="top">
				Today's Hindsight Investment: <strong>{{.description}}</strong>
			</p>
			<p class="top">
				Between {{.old_time}} and {{.new_time}},
				<strong>{{.symbol}}</strong> increased in value by
				<strong>{{.percent_increase}}</strong>%.
			</p>
			<h2>Why?</h2>
			<p>
				The purpose of this page is to demonstrate hind-sight
				bias.  It's easy to blame ourselves for not
				having taken advantage of an opportunity that seems
				obvious in retrospect, but the reality is that every
				day we miss countless investments which, unbeknownst
				to anyone at the time, are destined to increase in
				value drastically.  Every day this page will display
				one of the higher-performing securities from a year
				ago, hopefully illustrating the fact that most people
				had no idea what would come next.
			</p>
			<p>
				You can find the source code for this project at
				<a href="https://www.github.com/bieber/dailyhindsight/">
					github.com/bieber/dailyhindsight
				</a>.  If you find erroneous data displayed here,
				please let me know by writing to
				<a href="mailto:baddata@biebersprojects.com">
					baddata@biebersprojects.com
				</a>.
			</p>
		</div>
	</body>
</html>
`,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func Middleware(in http.Handler) http.Handler {
	logged := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"[%s] %s %s",
				r.Method,
				r.RemoteAddr,
				r.URL.String(),
			)
			in.ServeHTTP(w, r)
		},
	)

	xffm, err := xff.Default()
	if err != nil {
		panic(err)
	}
	return xffm.Handler(logged)
}

func IndexHandler(
	config Config,
	selection *DailySelection,
	selectionLock *sync.RWMutex,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rawIncrease := selection.NewValue - selection.OldValue
			percentIncrease := 100 * rawIncrease / selection.OldValue
			timeFormat := "January 2, 2006"

			data := map[string]string{
				"symbol":           selection.Dataset.Dataset,
				"description":      selection.Dataset.Description,
				"percent_increase": fmt.Sprintf("%.0f", percentIncrease),
				"old_time":         selection.OldTime.Format(timeFormat),
				"new_time":         selection.NewTime.Format(timeFormat),
			}

			err := indexTemplate.ExecuteTemplate(w, "index", data)
			if err != nil {
				panic(err)
			}
		},
	)
}
