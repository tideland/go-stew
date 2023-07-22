// Tideland Go Stew - Loop
//
// Copyright (C) 2017-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package loop supports the developer implementing the typical Go
// idiom for concurrent applications running a loop selecting data from
// channels in a backend goroutine. Stopping  those loops or getting aware
// of internal errors requires extra efforts. The loop package helps to
// control this kind of goroutines.
//
//	type Printer struct {
//		prints chan string
//		loop   loop.Loop
//	}
//
//	func NewPrinter(ctx context.Context) (*Printer, error) {
//		p := &Printer{
//			ctx:    ctx,
//			prints: make(chan string),
//		}
//		l, err := loop.Go(
//			p.backend,
//			loop.WithContext(ctx),
//			loop.WithFinalizer(func(err error) error {
//				...
//			})
//		if err != nil {
//			return nil, err
//		}
//		p.loop = l
//		return p, nil
//	}
//
//	func (p *printer) backend(ctx context.Context) error {
//		for {
//			select {
//			case <-ctx.Done():
//				return ctx.Err()
//			case str := <-p.prints:
//				println(str)
//		}
//	}
//
// The backend here now can be stopped with p.loop.Stop() returning
// a possible internal error. Also recovering of internal panics with
// a repairer function passed as option is possible. See the code
// examples.
package loop // import "tideland.dev/go/stew/loop"

// EOF
