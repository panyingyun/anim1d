package anim1d

import "github.com/maruel/fastbezier"

// Scalers

// Bell is a "good enough" approximation of a gaussian curve by using 2
// symmetrical ease-in-out bezier curves.
//
// It is not named Gaussian since it is not a gaussian curve; it really is a
// bell.
type Bell struct{}

// Scale scales input [0, 65535] to output [0, 65535] as a bell curve.
func (b *Bell) Scale(v uint16) uint16 {
	switch {
	case v == 0:
		return 0
	case v == 65535:
		return 0
	case v == 32767:
		return 65535

	case v < 32767:
		return EaseInOut.Scale(v * 2)
	default:
		return EaseInOut.Scale(65535 - v*2)
	}
}

// Curve models visually pleasing curves.
//
// They are modeled against CSS transitions.
// https://www.w3.org/TR/web-animations/#scaling-using-a-cubic-bezier-curve
type Curve string

// All the kind of known curves.
const (
	Ease       Curve = "ease"
	EaseIn     Curve = "ease-in"
	EaseInOut  Curve = "ease-in-out"
	EaseOut    Curve = "ease-out" // Recommended and default value.
	Direct     Curve = "direct"   // linear mapping
	StepStart  Curve = "steps(1,start)"
	StepMiddle Curve = "steps(1,middle)"
	StepEnd    Curve = "steps(1,end)"
)

var lutCache map[Curve]fastbezier.LUT

func setupCache() map[Curve]fastbezier.LUT {
	cache := map[Curve]fastbezier.LUT{
		Ease:      fastbezier.Make(0.25, 0.1, 0.25, 1, 18),
		EaseIn:    fastbezier.Make(0.42, 0, 1, 1, 18),
		EaseInOut: fastbezier.Make(0.42, 0, 0.58, 1, 18),
		EaseOut:   fastbezier.Make(0, 0, 0.58, 1, 18),
	}
	cache[""] = cache[EaseOut]
	return cache
}

func init() {
	lutCache = setupCache()
}

// Scale scales input [0, 65535] to output [0, 65535] using the curve
// requested.
func (c Curve) Scale(intensity uint16) uint16 {
	switch c {
	case Ease, EaseIn, EaseInOut, EaseOut, "":
		return lutCache[c].Eval(intensity)
	default:
		return lutCache[""].Eval(intensity)
	case Direct:
		return intensity
	case StepStart:
		if intensity < 256 {
			return 0
		}
		return 65535
	case StepMiddle:
		if intensity < 32768 {
			return 0
		}
		return 65535
	case StepEnd:
		if intensity >= 65535-256 {
			return 65535
		}
		return 0
	}
}

// Scale8 saves on casting.
func (c Curve) Scale8(intensity uint16) uint8 {
	return uint8(c.Scale(intensity) >> 8)
}

// Interpolation specifies a way to scales a pixel strip.
type Interpolation string

// All the kinds of interpolations.
const (
	NearestSkip Interpolation = "nearestskip" // Selects the nearest pixel but when upscaling, skips on missing pixels.
	Nearest     Interpolation = "nearest"     // Selects the nearest pixel, gives a blocky view.
	Linear      Interpolation = "linear"      // Linear interpolation, recommended and default value.
)

// Scale interpolates a frame into another using integers as much as possible
// for reasonable performance.
func (i Interpolation) Scale(in, out Frame) {
	li := len(in)
	lo := len(out)
	if li == 0 || lo == 0 {
		return
	}
	switch i {
	case NearestSkip:
		if li < lo {
			// Do not touch skipped pixels.
			for i, p := range in {
				out[(i*lo+lo/2)/li] = p
			}
			return
		}
		// When the destination is smaller than the source, Nearest and NearestSkip
		// have the same behavior.
		fallthrough
	case Nearest, "":
		fallthrough
	default:
		for i := range out {
			out[i] = in[(i*li+li/2)/lo]
		}
	case Linear:
		for i := range out {
			x := (i*li + li/2) / lo
			c := in[x]
			if x < li-1 {
				gradient := uint8(127)
				c.Mix(in[x+1], gradient)
			}
			out[i] = c
			//a := in[(i*li+li/2)/lo]
			//b := in[(i*li+li/2)/lo]
			//out[i] = (a + b) / 2
		}
	}
}