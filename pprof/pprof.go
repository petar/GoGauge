// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pprof

import (
	"bytes"
	"io/ioutil"
	"log"
	"http"
	_ "http/pprof"
	"runtime/pprof"
	"strconv"
	"time"
)

func StartHTTP(bind string) {
	go func() {
		err := http.ListenAndServe(bind, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err.Error())
		}
	}()
}

func StartLogging(filename string, interval int64) {
	t := time.LocalTime()
	filename = filename + "-" + strconv.Itoa64(t.Year) + strconv.Itoa(t.Month) + strconv.Itoa(t.Day)
	go func() {
		for k := 0; ; k++ {
			time.Sleep(interval)
			var w bytes.Buffer
			err := pprof.WriteHeapProfile(&w)
			if err != nil {
				log.Printf("preparing pprof: %s\n", err)
				break
			}
			err = ioutil.WriteFile(filename + "-" + strconv.Itoa(k), w.Bytes(), 0666)
			if err != nil {
				log.Printf("writing pprof to file: %s\n", err)
				break
			}
		}
	}()
}
