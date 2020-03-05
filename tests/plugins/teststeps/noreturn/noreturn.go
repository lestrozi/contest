// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package noreturn

import (
	"github.com/facebookincubator/contest/pkg/cerrors"
	"github.com/facebookincubator/contest/pkg/event"
	"github.com/facebookincubator/contest/pkg/event/testevent"
	"github.com/facebookincubator/contest/pkg/test"
)

// Name is the name used to look this plugin up.
var Name = "NoReturn"

type noreturnStep struct {
}

// Name returns the name of the Step
func (ts *noreturnStep) Name() string {
	return Name
}

// Run executes a step which does never return.
func (ts *noreturnStep) Run(cancel, pause <-chan struct{}, ch test.TestStepChannels, params test.TestStepParameters, ev testevent.Emitter) error {
	for target := range ch.In {
		ch.Out <- target
	}
	channel := make(chan struct{})
	<-channel
	return nil
}

// ValidateParameters validates the parameters associated to the TestStep
func (ts *noreturnStep) ValidateParameters(params test.TestStepParameters) error {
	return nil
}

// Resume tries to resume a previously interrupted test step. ExampleTestStep
// cannot resume.
func (ts *noreturnStep) Resume(cancel, pause <-chan struct{}, ch test.TestStepChannels, params test.TestStepParameters, ev testevent.EmitterFetcher) error {
	return &cerrors.ErrResumeNotSupported{StepName: Name}
}

// CanResume tells whether this step is able to resume.
func (ts *noreturnStep) CanResume() bool {
	return false
}

// Factory implements test.TestStepFactory
type Factory struct{}

// New constructs and returns a "noreturn" implementation of test.TestStep
func (f *Factory) New() test.TestStep {
	return &noreturnStep{}
}

// Events defines the events that a TestStep is allow to emit
func (f *Factory) Events() []event.Name {
	return nil
}

// UniqueImplementationName returns the unique name of the implementation
func (f *Factory) UniqueImplementationName() string {
	return Name
}
