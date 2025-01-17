// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bmc

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/ironcore-dev/metal/internal/log"
)

type BMC interface {
	Type() string
	Tags() map[string]string
	Credentials() (Credentials, time.Time)
	EnsureInitialCredentials(ctx context.Context, defaultCreds []Credentials, tempPassword string) error
	Connect(ctx context.Context) error
	CreateUser(ctx context.Context, creds Credentials, tempPassword string) error
	DeleteUsers(ctx context.Context, regex *regexp.Regexp) error
	ReadInfo(ctx context.Context) (Info, error)
}

type LEDControl interface {
	SetLocatorLED(ctx context.Context, state string) (string, error)
}

type PowerControl interface {
	PowerOn(ctx context.Context) error
	PowerOff(ctx context.Context, immediate bool) error
}

type ResetControl interface {
	Reset(ctx context.Context, immediate bool) error
}

type newBMCFunc func(tags map[string]string, host string, port int, creds Credentials, exp time.Time) BMC

var (
	bmcs = make(map[string]newBMCFunc)
)

func registerBMC(newFunc newBMCFunc) {
	bmcs[newFunc(nil, "", 0, Credentials{}, time.Time{}).Type()] = newFunc
}

func NewBMC(typ string, tags map[string]string, host string, port int, creds Credentials, exp time.Time) (BMC, error) {
	newFunc, ok := bmcs[typ]
	if !ok {
		return nil, fmt.Errorf("BMC of type %s is not supported", typ)
	}

	return newFunc(tags, host, port, creds, exp), nil
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Info struct {
	UUID         string
	Type         string
	Capabilities []string
	SerialNumber string
	SKU          string
	Manufacturer string
	LocatorLED   string
	Power        string
	OS           string
	OSReason     string
	Console      string
	FWVersion    string
}

func must(ctx context.Context, err error) {
	if err != nil {
		log.Error(ctx, fmt.Errorf("impossible error (this should never happen lol): %w", err))
	}
}
