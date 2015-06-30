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

// Package pcf8574 allows interfacing with the pcf8574 8-bit I/O Expander through I2C protocol.
package pcf8574

import (
	"github.com/kidoman/embd"
)

// PCF8574 represents a PCF8574 I/O Expander.
type PCF8574 struct {
	bus  embd.I2CBus
	addr byte
}

// New creates a new PCF8574 interface.
func New(bus embd.I2CBus, addr byte) *PCF8574 {
	return &PCF8574{
		bus:  bus,
		addr: addr,
	}
}

func (c *PCF8574) SetByte(val byte) error {
	return c.bus.WriteByte(c.addr, val)
}
