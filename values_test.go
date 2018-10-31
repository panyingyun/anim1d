// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"testing"
)

// Values

func TestSValue_Eval(t *testing.T) {
	var s SValue
	if s.Eval(23, 0) != 0 {
		t.Fail()
	}
}

func TestConst(t *testing.T) {
	if Const(2).Eval(23, 0) != 2 {
		t.Fail()
	}
}

func TestPercent(t *testing.T) {
	data := []struct {
		p        int32
		timeMS   uint32
		l        int
		expected int32
	}{
		{0, 0, 0, 0},
		{65536, 0, 0, 0},
		{65536, 1000, 0, 0},
		{65536, 0, 1000, 1000},
		{6554, 0, 1000, 100},
		{-65536, 0, 1000, -1000},
		{-6554, 0, 1000, -100},
	}
	for i, line := range data {
		if v := Percent(line.p).Eval(line.timeMS, line.l); v != line.expected {
			t.Fatalf("%d: Percent(%v).Eval(%v, %v) = %v, expected %v", i, line.p, line.timeMS, line.l, v, line.expected)
		}
	}
}

func TestOpAdd(t *testing.T) {
	if (&OpAdd{AddMS: 2}).Eval(23, 0) != 25 {
		t.Fail()
	}
	if (&OpAdd{AddMS: -2}).Eval(23, 0) != 21 {
		t.Fail()
	}
}

func TestOpMod(t *testing.T) {
	if (&OpMod{TickMS: 2}).Eval(23, 0) != 1 {
		t.Fail()
	}
}

func TestOpStep(t *testing.T) {
	if (&OpStep{TickMS: 21}).Eval(23, 0) != 21 {
		t.Fail()
	}
}

func TestRand(t *testing.T) {
	r1 := Rand{0}
	r2 := Rand{16}
	if r1.Eval(0, 0) != r2.Eval(15, 0) {
		t.Fail()
	}
	if r1.Eval(15, 0) == r2.Eval(16, 0) {
		t.Fail()
	}
	if r1.Eval(23, 0) != r2.Eval(23, 0) {
		t.Fail()
	}
}

// Scalers

func TestBell_Scale(t *testing.T) {
	half := uint16(65535 >> 1)
	data := []struct {
		i        uint16
		expected uint16
	}{
		{0, 0},
		{0x1000, 2093},
		{half, 0xffff},
		{0xefff, 2093},
		{0xffff, 0},
	}
	b := Bell{}
	for i, line := range data {
		if v := b.Scale(line.i); v != line.expected {
			t.Fatalf("%d: Bell.Scale(%v) = %v, expected %v", i, line.i, v, line.expected)
		}
	}
}

func TestCurve_limits(t *testing.T) {
	for _, v := range []Curve{Curve(""), Ease, EaseIn, EaseInOut, EaseOut, Direct} {
		if s := v.Scale(0); s != 0 {
			t.Fatalf("limit low %d != 0", s)
		}
		if s := v.Scale(65535); s != 65535 {
			t.Fatalf("limit high %d != 65535", s)
		}
	}
}

func TestCurve_Scale(t *testing.T) {
	half := uint16(65535 >> 1)
	data := []struct {
		t        Curve
		i        uint16
		expected uint16
	}{
		{Ease, half, 0xcd01},
		{EaseIn, half, 0x50df},
		{EaseInOut, half, 0x7ffe},
		{EaseOut, half, 0xaf1d},
		{Curve(""), half, 0xaf1d},
		{Curve("bleh"), half, 0xaf1d},
		{Direct, half, half},
		{StepStart, 0, 0},
		{StepStart, 255, 0},
		{StepStart, 256, 0xffff},
		{StepStart, 0xffff, 0xffff},
		{StepMiddle, 0, 0},
		{StepMiddle, 0x7fff, 0},
		{StepMiddle, 0x8000, 0xffff},
		{StepMiddle, 0xffff, 0xffff},
		{StepEnd, 0, 0},
		{StepEnd, 0xfefe, 0},
		{StepEnd, 0xfeff, 0xffff},
		{StepEnd, 0xffff, 0xffff},
	}
	for i, line := range data {
		if v := line.t.Scale(line.i); v != line.expected {
			t.Fatalf("%d: %v.Scale(%v) = %v, expected %v", i, line.t, line.i, v, line.expected)
		}
		if v := line.t.Scale8(line.i); v != uint8(line.expected>>8) {
			t.Fatalf("%d: %v.Scale8(%v) = %v, expected %v", i, line.t, line.i, v, line.expected>>8)
		}
	}
}

func TestMovePerHour(t *testing.T) {
	data := []struct {
		mps      int32
		timeMS   uint32
		cycle    int
		expected int
	}{
		{1, 0, 10, 0},
		{1, 3600000, 10, 1},
		{1, 2 * 3600000, 10, 2},
		{1, 3 * 3600000, 10, 3},
		{1, 4 * 3600000, 10, 4},
		{1, 5 * 3600000, 10, 5},
		{1, 6 * 3600000, 10, 6},
		{1, 7 * 3600000, 10, 7},
		{1, 8 * 3600000, 10, 8},
		{1, 9 * 3600000, 10, 9},
		{1, 10 * 3600000, 10, 0},
		{1, 10 * 3600000, 11, 10},
		{60, 16, 10, 0},
		{60, 1000, 9, 0},
		{60, 1000, 10, 0},
		{60, 3600000, 10, 0},
		{3600, 3600000, 10, 0},
		{3600000, 0, 10, 0},
		{3600000, 1, 10, 1},
		{3600000, 2, 10, 2},
		{2 * 3600000, 1, 10, 1},
		{2 * 3600000, 2, 10, 2},
		{2 * 3600000, 2, 0, 2},
	}
	for i, line := range data {
		m := MovePerHour{Const(line.mps)}
		if actual := m.Eval(line.timeMS, 0, line.cycle); actual != line.expected {
			t.Fatalf("%d: %d.Eval(%d, %d) = %d != %d", i, line.mps, line.timeMS, line.cycle, actual, line.expected)
		}
	}
}

func BenchmarkSetupCache(b *testing.B) {
	// Calculate how much this one-time initialization cost is.
	for i := 0; i < b.N; i++ {
		setupCache()
	}
}
