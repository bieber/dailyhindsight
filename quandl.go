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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const timeFormat string = "2006-01-02"

var client *http.Client = nil

func init() {
	client = &http.Client{Timeout: 10 * time.Second}
}

type RequestResult struct {
	oldValue, newValue float64
	oldTime, newTime   time.Time
}

func GetRequest(
	apiKey string,
	t1 time.Time,
	dataset Dataset,
) (RequestResult, error) {
	uri := url.URL{
		Scheme: "https",
		Host:   "www.quandl.com",
		Path: fmt.Sprintf(
			"/api/v3/datasets/%s/%s/data.json",
			dataset.Database,
			dataset.Dataset,
		),
	}

	t0 := t1.Add(-1 * time.Hour * 24 * 365)

	q := uri.Query()
	q.Set("api_key", apiKey)
	q.Set("column_index", strconv.Itoa(4))
	q.Set("start_date", t0.Format(timeFormat))
	q.Set("end_date", t1.Format(timeFormat))
	uri.RawQuery = q.Encode()

	response, err := client.Get(uri.String())
	if err != nil {
		return RequestResult{}, nil
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	result := struct {
		DatasetData struct {
			Data [][]interface{} `json:"data"`
		} `json:"dataset_data"`
	}{}
	err = decoder.Decode(&result)

	if err != nil {
		return RequestResult{}, err
	}

	allData := result.DatasetData.Data
	// Data comes ordered by date descending
	oldData, newData := allData[len(allData)-1], allData[0]

	oldTime, err := time.Parse(timeFormat, oldData[0].(string))
	if err != nil {
		return RequestResult{}, err
	}
	newTime, err := time.Parse(timeFormat, newData[0].(string))
	if err != nil {
		return RequestResult{}, err
	}

	oldValue, newValue := oldData[1].(float64), newData[1].(float64)
	return RequestResult{
		oldValue, newValue,
		oldTime, newTime,
	}, nil

	return RequestResult{}, nil
}
