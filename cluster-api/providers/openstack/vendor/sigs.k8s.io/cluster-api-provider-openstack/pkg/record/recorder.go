/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package record implements recording functionality.
package record

import (
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
)

var (
	initOnce        sync.Once
	defaultRecorder events.EventRecorder
)

// noopRecorder is a no-op implementation of events.EventRecorder used as the
// default before a real recorder is initialised.
type noopRecorder struct{}

func (noopRecorder) Eventf(runtime.Object, runtime.Object, string, string, string, string, ...interface{}) {
}

func init() {
	defaultRecorder = noopRecorder{}
}

// InitFromRecorder initializes the global default recorder. It can only be called once.
// Subsequent calls are considered noops.
func InitFromRecorder(recorder events.EventRecorder) {
	initOnce.Do(func() {
		defaultRecorder = recorder
	})
}

// Event constructs an event from the given information and puts it in the queue for sending.
func Event(object runtime.Object, reason, message string) {
	defaultRecorder.Eventf(object, nil, corev1.EventTypeNormal, cases.Title(language.English).String(reason), cases.Title(language.English).String(reason), message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func Eventf(object runtime.Object, reason, message string, args ...interface{}) {
	defaultRecorder.Eventf(object, nil, corev1.EventTypeNormal, cases.Title(language.English).String(reason), cases.Title(language.English).String(reason), message, args...)
}

// Warn constructs a warning event from the given information and puts it in the queue for sending.
func Warn(object runtime.Object, reason, message string) {
	defaultRecorder.Eventf(object, nil, corev1.EventTypeWarning, cases.Title(language.English).String(reason), cases.Title(language.English).String(reason), message)
}

// Warnf is just like Warn, but with Sprintf for the message field.
func Warnf(object runtime.Object, reason, message string, args ...interface{}) {
	defaultRecorder.Eventf(object, nil, corev1.EventTypeWarning, cases.Title(language.English).String(reason), cases.Title(language.English).String(reason), message, args...)
}
