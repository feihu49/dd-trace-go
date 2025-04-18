// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package appsec

import (
	"runtime"

	"github.com/DataDog/dd-trace-go/v2/internal/telemetry"
	waf "github.com/DataDog/go-libddwaf/v3"
)

// cgoEnabled is true if cgo is enabled, false otherwise.
// No way to check this at runtime, so we compute it at build time in
// telemetry_cgo.go.
var cgoEnabled bool

type appsecTelemetry struct {
	configs []telemetry.Configuration
	enabled bool
}

var (
	wafSupported, _ = waf.SupportsTarget()
	wafHealthy, _   = waf.Health()
	staticConfigs   = []telemetry.Configuration{
		{Name: "goos", Value: runtime.GOOS, Origin: telemetry.OriginCode},
		{Name: "goarch", Value: runtime.GOARCH, Origin: telemetry.OriginCode},
		{Name: "waf_supports_target", Value: wafSupported, Origin: telemetry.OriginCode},
		{Name: "waf_healthy", Value: wafHealthy, Origin: telemetry.OriginCode},
	}
)

// newAppsecTelemetry creates a new telemetry event for AppSec.
func newAppsecTelemetry() *appsecTelemetry {
	if telemetry.Disabled() {
		// If telemetry is disabled, we won't do anything...
		return nil
	}

	configs := make([]telemetry.Configuration, len(staticConfigs)+1, len(staticConfigs)+2)
	configs[0] = telemetry.Configuration{Name: "cgo_enabled", Value: cgoEnabled}
	copy(configs[1:], staticConfigs)

	return &appsecTelemetry{
		configs: configs,
	}
}

// addConfig adds a new configuration entry to this telemetry event.
func (a *appsecTelemetry) addConfig(name string, value any) {
	if a == nil {
		return
	}
	a.configs = append(a.configs, telemetry.Configuration{Name: name, Value: value})
}

// addCodeConfig adds a new configuration entry to this telemetry event.
func (a *appsecTelemetry) addCodeConfig(name string, value any) {
	if a == nil {
		return
	}
	a.configs = append(a.configs, telemetry.Configuration{Name: name, Value: value, Origin: telemetry.OriginCode})
}

// addEnvConfig adds a new envionment-sourced configuration entry to this event.
func (a *appsecTelemetry) addEnvConfig(name string, value any) {
	if a == nil {
		return
	}
	a.configs = append(a.configs, telemetry.Configuration{Name: name, Value: value, Origin: telemetry.OriginEnvVar})
}

// setEnabled makes AppSec as having effectively been enabled.
func (a *appsecTelemetry) setEnabled() {
	if a == nil {
		return
	}
	a.enabled = true
}

// emit sends the telemetry event to the telemetry.GlobalClient.
func (a *appsecTelemetry) emit() {
	if a == nil {
		return
	}

	if a.enabled {
		telemetry.ProductStarted(telemetry.NamespaceAppSec)
	} else {
		telemetry.ProductStopped(telemetry.NamespaceAppSec)
	}

	telemetry.RegisterAppConfigs(a.configs...)
}
