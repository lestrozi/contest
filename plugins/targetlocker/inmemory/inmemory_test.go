// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package inmemory

import (
	"testing"
	"time"

	"github.com/facebookincubator/contest/pkg/target"
	"github.com/facebookincubator/contest/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	jobID      = types.JobID(123)
	otherJobID = types.JobID(456)

	targetOne  = target.Target{Name: "target001", ID: "001"}
	targetTwo  = target.Target{Name: "target002", ID: "002"}
	oneTarget  = []*target.Target{&targetOne}
	twoTargets = []*target.Target{&targetOne, &targetTwo}
)

func TestInMemoryNew(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NotNil(t, tl)
	require.IsType(t, &InMemory{}, tl)
}

func TestInMemoryLockInvalidJobIDAndNoTargets(t *testing.T) {
	tl := New(time.Second, time.Second)
	assert.Error(t, tl.Lock(0, nil))
}

func TestInMemoryLockValidJobIDAndNoTargets(t *testing.T) {
	tl := New(time.Second, time.Second)
	assert.Error(t, tl.Lock(jobID, nil))
}

func TestInMemoryLockInvalidJobIDAndOneTarget(t *testing.T) {
	tl := New(time.Second, time.Second)
	assert.Error(t, tl.Lock(0, oneTarget))
}

func TestInMemoryLockValidJobIDAndOneTarget(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.Lock(jobID, oneTarget))
}

func TestInMemoryLockValidJobIDAndTwoTargets(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.Lock(jobID, twoTargets))
}

func TestInMemoryLockReentrantLock(t *testing.T) {
	tl := New(10*time.Second, time.Second)
	require.NoError(t, tl.Lock(jobID, twoTargets))
	require.NoError(t, tl.Lock(jobID, twoTargets))
}

func TestInMemoryLockReentrantLockDifferentJobID(t *testing.T) {
	tl := New(10*time.Second, time.Second)
	require.NoError(t, tl.Lock(jobID, twoTargets))
	require.Error(t, tl.Lock(jobID+1, twoTargets))
}

func TestInMemoryUnlockInvalidJobIDAndNoTargets(t *testing.T) {
	tl := New(time.Second, time.Second)
	assert.Error(t, tl.Unlock(jobID, nil))
}

func TestInMemoryUnlockValidJobIDAndNoTargets(t *testing.T) {
	tl := New(time.Second, time.Second)
	assert.Error(t, tl.Unlock(jobID, nil))
}

func TestInMemoryUnlockInvalidJobIDAndOneTarget(t *testing.T) {
	tl := New(time.Second, time.Second)
	assert.Error(t, tl.Unlock(0, oneTarget))
}

func TestInMemoryUnlockValidJobIDAndOneTarget(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.Error(t, tl.Unlock(jobID, oneTarget))
}

func TestInMemoryUnlockValidJobIDAndTwoTargets(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.Error(t, tl.Unlock(jobID, twoTargets))
}

func TestInMemoryUnlockUnlockTwice(t *testing.T) {
	tl := New(time.Second, time.Second)
	err := tl.Unlock(jobID, oneTarget)
	log.Print(err)
	assert.Error(t, err)
	assert.Error(t, tl.Unlock(jobID, oneTarget))
}

func TestInMemoryUnlockReentrantLockDifferentJobID(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.Error(t, tl.Unlock(jobID, twoTargets))
	assert.Error(t, tl.Unlock(jobID+1, twoTargets))
}

func TestInMemoryLockUnlockSameJobID(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.Lock(jobID, twoTargets))
	assert.NoError(t, tl.Unlock(jobID, twoTargets))
}

func TestInMemoryLockUnlockDifferentJobID(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.Lock(jobID, twoTargets))
	assert.Error(t, tl.Unlock(jobID+1, twoTargets))
}

func TestInMemoryRefreshLocks(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.RefreshLocks(jobID, twoTargets))
}

func TestInMemoryRefreshLocksTwice(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.RefreshLocks(jobID, twoTargets))
	assert.NoError(t, tl.RefreshLocks(jobID, twoTargets))
}

func TestInMemoryRefreshLocksOneThenTwo(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.RefreshLocks(jobID, oneTarget))
	assert.NoError(t, tl.RefreshLocks(jobID, twoTargets))
}

func TestInMemoryRefreshLocksTwoThenOne(t *testing.T) {
	tl := New(time.Second, time.Second)
	require.NoError(t, tl.RefreshLocks(jobID, twoTargets))
	assert.NoError(t, tl.RefreshLocks(jobID, oneTarget))
}

func TestRefreshMultiple(t *testing.T) {
	tl := New(200*time.Millisecond, 200*time.Millisecond)
	require.NoError(t, tl.Lock(jobID, twoTargets))
	time.Sleep(100 * time.Millisecond)
	// they are not expired yet, extend both
	require.NoError(t, tl.RefreshLocks(jobID, twoTargets))
	time.Sleep(150 * time.Millisecond)
	// if they were refreshed properly, they are still valid and attempts to get them must fail
	require.Error(t, tl.Lock(otherJobID, []*target.Target{&targetOne}))
	require.Error(t, tl.Lock(otherJobID, []*target.Target{&targetTwo}))
}

func TestLockingTransactional(t *testing.T) {
	tl := New(time.Second, time.Second)
	// lock the second target
	require.NoError(t, tl.Lock(jobID, []*target.Target{&targetTwo}))
	// try to lock both with another owner (this fails as expected)
	require.Error(t, tl.Lock(jobID+1, twoTargets))
	// API says target one should remain unlocked because Lock() is transactional
	// this means it can be locked by the first owner
	require.NoError(t, tl.Lock(jobID, []*target.Target{&targetOne}))
}
