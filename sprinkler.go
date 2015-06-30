/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/kidoman/embd"
	"sprinkler/embd/controller/pcf8574" // Our own custom controller

	_ "github.com/kidoman/embd/host/all"
)

type zoneslice []uint64

func (i *zoneslice) String() string {
	return fmt.Sprint(*i)
}

func (i *zoneslice) Set(value string) error {
	for _, z := range strings.Split(value, ",") {
		zone, err := strconv.ParseUint(z, 0, 64)
		if err != nil {
			return err
		}
		if zone >= 1 && zone <= 8 {
			*i = append(*i, zone)
		}
	}
	return nil
}

func main() {
	var zoneFlag zoneslice
	var timeout time.Duration
	var verbose bool
	var repeat uint64

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "amount of time to run each zone")
	flag.Uint64Var(&repeat, "repeat", 1, "how many times to repeat the zone list")
	flag.BoolVar(&verbose, "verbose", false, "whether or not to be verbose")
	flag.Var(&zoneFlag, "zone", "comma-separated list of zones to use")
	flag.Parse()

	if len(zoneFlag) == 0 {
		fmt.Fprintf(os.Stderr, "No zone specified\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if timeout == 0 && len(zoneFlag) > 1 {
		fmt.Fprintf(os.Stderr, "No timeout with multiple zones\n")
		os.Exit(1)
	}

	if repeat == 0 {
		fmt.Fprintf(os.Stderr, "Nothing to do!\n")
		os.Exit(1)
	}

	if repeat > 1 && len(zoneFlag) == 1 {
		fmt.Fprintf(os.Stderr, "Repeating only one zone\n")
		os.Exit(1)
	}

	if timeout > 0 && timeout < 5*time.Second {
		fmt.Fprintf(os.Stderr, "Increasing timeout to 5s\n")
		timeout = 5 * time.Second
	}

	if verbose {
		fmt.Printf("Zones: %v\n", zoneFlag)
		fmt.Printf("Repeating: %d\n", repeat)
		fmt.Printf("Timeout: %v\n", timeout)
	}

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)
	defer bus.Close()

	pcf8574 := pcf8574.New(bus, 0x20)
	defer pcf8574.SetByte(^byte(0))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	timer := time.Tick(timeout)

	for i := uint64(0); i < repeat; i++ {
		for _, zone := range zoneFlag {
			input := ^(byte(1) << (zone - 1))
			pcf8574.SetByte(input)
			if verbose {
				fmt.Printf("Zone: %d Input: %08b\n", zone, input)
			}
			select {
			case <-timer:
				continue
			case <-c:
				return
			}
		}
	}
}
