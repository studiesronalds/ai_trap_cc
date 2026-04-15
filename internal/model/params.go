package model

import "fmt"

// Params holds all simulation parameters. Passed by value throughout — makes
// scenario save/load trivial (copy the struct).
type Params struct {
	N      float64 // number of firms
	W      float64 // wage per task
	C      float64 // AI cost per task
	K      float64 // integration friction
	Lambda float64 // worker MPC
	Eta    float64 // income replacement rate
	A      float64 // autonomous demand
	L      float64 // tasks per firm
	Mu     float64 // social planner worker weight [0,1]
	Phi    float64 // AI productivity multiplier (extension, default 1)
	EtaHat float64 // capital income recycling rate (extension, default 0)
}

// PaperBase returns the illustrative parameters from the paper.
func PaperBase() Params {
	return Params{N: 7, W: 1.0, C: 0.30, K: 1.0, Lambda: 0.5, Eta: 0.30, A: 10.0, L: 100, Mu: 0.5, Phi: 1.0}
}

// Frictionless returns k=0 (Prisoner's Dilemma mode).
func Frictionless() Params {
	p := PaperBase()
	p.K = 0
	return p
}

// Monopolist returns N=1 (fully internalises externality).
func Monopolist() Params {
	p := PaperBase()
	p.N = 1
	return p
}

// Duopoly returns N=2.
func Duopoly() Params {
	p := PaperBase()
	p.N = 2
	return p
}

// CompetitiveLimit returns N=50 (near-competitive market).
func CompetitiveLimit() Params {
	p := PaperBase()
	p.N = 50
	return p
}

// Optimist returns eta=1.1 (full plus replacement).
func Optimist() Params {
	p := PaperBase()
	p.Eta = 1.1
	return p
}

// Validate returns an error if any parameter is out of range.
func (p Params) Validate() error {
	if p.N < 1 {
		return fmt.Errorf("N must be >= 1, got %.4g", p.N)
	}
	if p.W < 0 {
		return fmt.Errorf("w must be >= 0, got %.4g", p.W)
	}
	if p.C < 0 {
		return fmt.Errorf("c must be >= 0, got %.4g", p.C)
	}
	if p.C > p.W {
		return fmt.Errorf("c must be <= w (%.4g), got %.4g", p.W, p.C)
	}
	if p.K < 0 {
		return fmt.Errorf("k must be >= 0, got %.4g", p.K)
	}
	if p.Lambda < 0 || p.Lambda > 1 {
		return fmt.Errorf("lambda must be in [0,1], got %.4g", p.Lambda)
	}
	if p.Eta < 0 {
		return fmt.Errorf("eta must be >= 0, got %.4g", p.Eta)
	}
	if p.A < 0 {
		return fmt.Errorf("A must be >= 0, got %.4g", p.A)
	}
	if p.L <= 0 {
		return fmt.Errorf("L must be > 0, got %.4g", p.L)
	}
	if p.Mu < 0 || p.Mu > 1 {
		return fmt.Errorf("mu must be in [0,1], got %.4g", p.Mu)
	}
	if p.Phi <= 0 {
		return fmt.Errorf("phi must be > 0, got %.4g", p.Phi)
	}
	return nil
}

// S returns cost saving per task: s = w - c
func (p Params) S() float64 { return p.W - p.C }

// Ell returns demand loss per automated task: ℓ = λ(1-η)w
func (p Params) Ell() float64 { return p.Lambda * (1 - p.Eta) * p.W }

// NStar returns the externality activation threshold: N* = ℓ/s
// Returns infinity when s=0.
func (p Params) NStar() float64 {
	s := p.S()
	if s <= 0 {
		return 1e18
	}
	return p.Ell() / s
}
