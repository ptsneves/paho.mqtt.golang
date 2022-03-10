/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Allan Stockdill-Mander
 */

package mqtt

import (
	"context"
	"errors"
	"sync"
	"time"
)

type baseToken struct {
	m       sync.RWMutex
	ctx     context.Context
	cancelF context.CancelFunc
	err     error
}

func newBaseToken(ctx context.Context) baseToken {
	ctx, cancelF := context.WithCancel(ctx)
	return baseToken{
		ctx:     ctx,
		cancelF: cancelF,
	}
}

// Wait implements the Token Wait method.
func (b *baseToken) Wait() bool {
	<-b.ctx.Done()
	return true
}

// WaitTimeout implements the Token WaitTimeout method.
func (b *baseToken) WaitTimeout(d time.Duration) bool {
	b.ctx.Deadline()
	ctxTimeout, timeoutCancelF := context.WithTimeout(b.ctx, d)
	defer timeoutCancelF()
	select {
	case <-ctxTimeout.Done():
		if errors.Is(ctxTimeout.Err(), context.DeadlineExceeded) {
			return false
		}
	}
	return true
}

// Done implements the Token Done method.
func (b *baseToken) Done() <-chan struct{} {
	return b.ctx.Done()
}

func (b *baseToken) flowComplete() {
	b.cancelF()
}

func (b *baseToken) Error() error {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.err
}

func (b *baseToken) setError(e error) {
	b.m.Lock()
	b.err = e
	b.m.Unlock()
	b.flowComplete()
}
