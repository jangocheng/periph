// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package uartreg

import (
	"errors"
	"sort"
	"testing"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/experimental/conn/uart"
)

func TestOpen(t *testing.T) {
	defer reset()
	if _, err := Open(""); err == nil {
		t.Fatal("no bus registered")
	}
	if err := Register("a", []string{"x"}, 1, fakeBuser); err != nil {
		t.Fatal(err)
	}
	if o, err := Open(""); o == nil || err != nil {
		t.Fatal(o, err)
	}
	if o, err := Open("1"); o == nil || err != nil {
		t.Fatal(o, err)
	}
	if o, err := Open("x"); o == nil || err != nil {
		t.Fatal(o, err)
	}
	if o, err := Open("y"); o != nil || err == nil {
		t.Fatal(o, err)
	}
}

func TestDefault_NoNumber(t *testing.T) {
	defer reset()
	if err := Register("a", nil, -1, fakeBuser); err != nil {
		t.Fatal(err)
	}
	if o, err := Open(""); o == nil || err != nil {
		t.Fatal(o, err)
	}
}

func TestAll(t *testing.T) {
	defer reset()
	if a := All(); len(a) != 0 {
		t.Fatal(a)
	}
	if err := Register("a", nil, 1, fakeBuser); err != nil {
		t.Fatal(err)
	}
	if err := Register("b", nil, 2, fakeBuser); err != nil {
		t.Fatal(err)
	}
	if a := All(); len(a) != 2 {
		t.Fatal(a)
	}
}

func TestRefList(t *testing.T) {
	l := refList{&Ref{Name: "b"}, &Ref{Name: "a"}}
	sort.Sort(l)
	if l[0].Name != "a" || l[1].Name != "b" {
		t.Fatal(l)
	}
}

func TestRegister(t *testing.T) {
	defer reset()
	if err := Register("a", []string{"b"}, 42, fakeBuser); err != nil {
		t.Fatal(err)
	}
	if Register("a", nil, -1, fakeBuser) == nil {
		t.Fatal("same bus name")
	}
	if Register("b", nil, -1, fakeBuser) == nil {
		t.Fatal("same bus alias name")
	}
	if Register("c", nil, 42, fakeBuser) == nil {
		t.Fatal("same bus number")
	}
	if Register("c", []string{"a"}, -1, fakeBuser) == nil {
		t.Fatal("same bus alias")
	}
	if Register("c", []string{"b"}, -1, fakeBuser) == nil {
		t.Fatal("same bus alias")
	}
}

func TestRegister_fail(t *testing.T) {
	defer reset()
	if Register("a", nil, -1, nil) == nil {
		t.Fatal("missing Opener")
	}
	if Register("a", nil, -2, fakeBuser) == nil {
		t.Fatal("bad bus number")
	}
	if Register("", nil, 42, fakeBuser) == nil {
		t.Fatal("missing name")
	}
	if Register("1", nil, 42, fakeBuser) == nil {
		t.Fatal("numeric name")
	}
	if Register("a:b", nil, 42, fakeBuser) == nil {
		t.Fatal("':' in name")
	}
	if Register("a", []string{"a"}, 0, fakeBuser) == nil {
		t.Fatal("\"a\" is already registered")
	}
	if Register("a", []string{""}, 0, fakeBuser) == nil {
		t.Fatal("empty alias")
	}
	if Register("a", []string{"1"}, 0, fakeBuser) == nil {
		t.Fatal("numeric alias")
	}
	if Register("a", []string{"a:b"}, 0, fakeBuser) == nil {
		t.Fatal("':' in alias")
	}
	if a := All(); len(a) != 0 {
		t.Fatal(a)
	}
}

func TestUnregister(t *testing.T) {
	defer reset()
	if Unregister("") == nil {
		t.Fatal("unregister empty")
	}
	if Unregister("a") == nil {
		t.Fatal("unregister non-existing")
	}
	if err := Register("a", []string{"b"}, 0, fakeBuser); err != nil {
		t.Fatal(err)
	}
	if err := Unregister("a"); err != nil {
		t.Fatal(err)
	}
}

//

func fakeBuser() (uart.ConnCloser, error) {
	return &fakeBus{}, nil
}

type fakeBus struct {
}

func (f *fakeBus) String() string {
	return "fake"
}

func (f *fakeBus) Close() error {
	return errors.New("not implemented")
}

func (f *fakeBus) Tx(w, r []byte) error {
	return errors.New("not implemented")
}

func (f *fakeBus) Duplex() conn.Duplex {
	return conn.DuplexUnknown
}

func (f *fakeBus) Speed(baud int) error {
	return errors.New("not implemented")
}

func (f *fakeBus) Configure(stopBit uart.Stop, parity uart.Parity, bits int) error {
	return errors.New("not implemented")
}

func reset() {
	mu.Lock()
	defer mu.Unlock()
	byName = map[string]*Ref{}
	byNumber = map[int]*Ref{}
	byAlias = map[string]*Ref{}
}
