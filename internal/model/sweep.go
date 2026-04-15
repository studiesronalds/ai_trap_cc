package model

import "fmt"

// SweepPoint is one data point in a parameter sweep.
type SweepPoint struct {
	ParamName string
	Value     float64
	Result    Result
}

// Sweep1D sweeps a single parameter from lo to hi in nSteps steps.
func Sweep1D(base Params, param string, lo, hi float64, nSteps int) ([]SweepPoint, error) {
	if nSteps < 2 {
		nSteps = 2
	}
	set, err := setter(param)
	if err != nil {
		return nil, err
	}
	step := (hi - lo) / float64(nSteps-1)
	pts := make([]SweepPoint, nSteps)
	for i := range pts {
		v := lo + float64(i)*step
		p := base
		set(&p, v)
		pts[i] = SweepPoint{ParamName: param, Value: v, Result: Compute(p)}
	}
	return pts, nil
}

// Sweep2DPoint is one cell in a 2D sweep.
type Sweep2DPoint struct {
	P1Name string
	P1     float64
	P2Name string
	P2     float64
	Wedge  float64
	AlphaNE float64
}

// Sweep2D sweeps two parameters, returning a grid.
func Sweep2D(base Params, p1Name string, p1Lo, p1Hi float64, p1Steps int,
	p2Name string, p2Lo, p2Hi float64, p2Steps int) ([][]Sweep2DPoint, error) {
	set1, err := setter(p1Name)
	if err != nil {
		return nil, err
	}
	set2, err := setter(p2Name)
	if err != nil {
		return nil, err
	}
	if p1Steps < 2 { p1Steps = 2 }
	if p2Steps < 2 { p2Steps = 2 }
	step1 := (p1Hi - p1Lo) / float64(p1Steps-1)
	step2 := (p2Hi - p2Lo) / float64(p2Steps-1)

	grid := make([][]Sweep2DPoint, p1Steps)
	for i := range grid {
		grid[i] = make([]Sweep2DPoint, p2Steps)
		v1 := p1Lo + float64(i)*step1
		for j := range grid[i] {
			v2 := p2Lo + float64(j)*step2
			p := base
			set1(&p, v1)
			set2(&p, v2)
			r := Compute(p)
			grid[i][j] = Sweep2DPoint{
				P1Name: p1Name, P1: v1,
				P2Name: p2Name, P2: v2,
				Wedge: r.Wedge, AlphaNE: r.AlphaNE,
			}
		}
	}
	return grid, nil
}

// setter returns a function that sets the named parameter on a *Params.
func setter(param string) (func(*Params, float64), error) {
	switch param {
	case "N":
		return func(p *Params, v float64) { p.N = v }, nil
	case "w":
		return func(p *Params, v float64) { p.W = v }, nil
	case "c":
		return func(p *Params, v float64) { p.C = v }, nil
	case "s":
		// sweep s = w-c by adjusting c (hold w fixed)
		return func(p *Params, v float64) { p.C = p.W - v }, nil
	case "k":
		return func(p *Params, v float64) { p.K = v }, nil
	case "lambda":
		return func(p *Params, v float64) { p.Lambda = v }, nil
	case "eta":
		return func(p *Params, v float64) { p.Eta = v }, nil
	case "A":
		return func(p *Params, v float64) { p.A = v }, nil
	case "L":
		return func(p *Params, v float64) { p.L = v }, nil
	case "mu":
		return func(p *Params, v float64) { p.Mu = v }, nil
	case "phi":
		return func(p *Params, v float64) { p.Phi = v }, nil
	default:
		return nil, fmt.Errorf("unknown parameter %q; valid: N w c s k lambda eta A L mu phi", param)
	}
}

// ParamNames returns all sweepable parameter names.
func ParamNames() []string {
	return []string{"N", "w", "c", "s", "k", "lambda", "eta", "A", "L", "mu", "phi"}
}
