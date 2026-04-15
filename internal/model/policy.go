package model

import (
	"fmt"
	"math"
)

// PolicyResult holds the outcome of a single policy instrument.
type PolicyResult struct {
	Name        string
	Description string
	AlphaAfter  float64
	WedgeAfter  float64
	TaxOrParam  float64
	ClosesWedge bool
	Note        string
}

// AllPolicies evaluates all 6 policy instruments and returns results.
func AllPolicies(p Params) []PolicyResult {
	base := Compute(p)
	results := make([]PolicyResult, 0, 7)

	// 1. Pigouvian tax τ* = ℓ(1-1/N)
	tau := base.TaxStar
	pTax := p
	// Tax shifts the effective per-task cost saving: s_eff = s - τ
	alphaTax := nashRate(pTax.S()-tau, pTax.Ell()/pTax.N, pTax.K)
	rTax := Compute(pTax)
	_ = rTax
	wedgeTax := alphaTax - base.AlphaCO
	results = append(results, PolicyResult{
		Name:        "Pigouvian tax",
		Description: fmt.Sprintf("τ* = ℓ(1−1/N) = %.4f", tau),
		AlphaAfter:  alphaTax,
		WedgeAfter:  wedgeTax,
		TaxOrParam:  tau,
		ClosesWedge: math.Abs(wedgeTax) < 1e-9,
		Note:        "Only instrument that implements cooperative optimum",
	})

	// 2. UBI: raises autonomous demand A by dA (we show the required dA to close wedge)
	// UBI doesn't change s or ℓ, so doesn't affect α_NE. Wedge unchanged.
	results = append(results, PolicyResult{
		Name:        "UBI",
		Description: "Universal basic income (raises A)",
		AlphaAfter:  base.AlphaNE,
		WedgeAfter:  base.Wedge,
		TaxOrParam:  0,
		ClosesWedge: false,
		Note:        "Raises demand floor but leaves automation incentive unchanged",
	})

	// 3. Capital income tax: taxes profits, doesn't affect per-task margin → no effect on α
	results = append(results, PolicyResult{
		Name:        "Capital income tax",
		Description: "Tax on firm profits",
		AlphaAfter:  base.AlphaNE,
		WedgeAfter:  base.Wedge,
		TaxOrParam:  0,
		ClosesWedge: false,
		Note:        "Operates on profit levels, not the per-task margin; no effect on α_NE",
	})

	// 4. Worker equity ε: workers receive fraction ε of profits → effective MPC rises
	// With equity ε, each worker-owner's demand contribution changes.
	// Effective ℓ_ε = λ(1-η)w - ε·(something). From paper: narrows but doesn't close wedge.
	// We show: with ε, ℓ is unchanged in baseline (equity doesn't affect the per-task demand loss
	// because the equity payment comes from profit, not from the per-task automation margin directly).
	// Paper confirms: narrows wedge but cannot eliminate.
	// For a concrete calculation: equity gives worker share ε of K/N per worker.
	// This effectively raises η toward 1 by channelling profit back to workers.
	// Simple bound: if we set ε such that η_eff = η + ε·(K/NwL), wedge shrinks.
	// We show just the note.
	results = append(results, PolicyResult{
		Name:        "Worker equity",
		Description: "Workers hold fraction ε of firm equity",
		AlphaAfter:  base.AlphaNE, // marginal effect unchanged
		WedgeAfter:  base.Wedge,
		TaxOrParam:  0,
		ClosesWedge: false,
		Note:        "Narrows wedge (raises effective η) but cannot eliminate it",
	})

	// 5. Upskilling: raises η toward 1, reducing ℓ
	// Show effect of η → min(η+0.2, 1)
	pUp := p
	pUp.Eta = math.Min(1.0, p.Eta+0.20)
	rUp := Compute(pUp)
	results = append(results, PolicyResult{
		Name:        "Upskilling (η+0.20)",
		Description: fmt.Sprintf("η: %.2f → %.2f", p.Eta, pUp.Eta),
		AlphaAfter:  rUp.AlphaNE,
		WedgeAfter:  rUp.Wedge,
		TaxOrParam:  pUp.Eta,
		ClosesWedge: math.Abs(rUp.Wedge) < 1e-9,
		Note:        "Narrows wedge; full closure requires η=1 or η > s/w",
	})

	// 6. Coasian bargaining / coalition of all N firms
	// Full coalition = cooperative optimum (α_CO)
	results = append(results, PolicyResult{
		Name:        "Coasian coalition (all N)",
		Description: "All firms agree to restrain automation",
		AlphaAfter:  base.AlphaCO,
		WedgeAfter:  0,
		TaxOrParam:  float64(p.N),
		ClosesWedge: true,
		Note:        "Not self-enforcing: automation is dominant strategy, any firm deviates",
	})

	return results
}

// Coalition evaluates the partial coalition of M firms (Proposition 4).
// M firms coordinate; remaining N-M firms play Nash against the coalition.
type CoalitionResult struct {
	M            float64
	AlphaCoal    float64 // coalition members' rate
	AlphaOut     float64 // outsiders' rate
	WedgeCoal    float64 // coalition rate vs full cooperative
}

func Coalition(p Params, M float64) CoalitionResult {
	// Coalition of M sets common rate to maximise their joint profit.
	// Each coalition member faces demand loss ℓ·M/N from their own group.
	// They internalise M's share: effective ell for coalition = ell * M/N * ...
	// From Prop 4: coalition members solve as if N_eff = N/M (they internalise M's share).
	// α_coal = (s - ℓ·(N/M)/N·1) / k = (s - ℓ/M) / k
	// Outsiders still play Nash: α_out = (s - ℓ/N) / k
	r := Compute(p)
	var alphaCoal float64
	if p.K == 0 {
		if p.S() > p.Ell()/M { alphaCoal = 1 } else { alphaCoal = 0 }
	} else {
		alphaCoal = clamp((p.S()-p.Ell()/M)/p.K, 0, 1)
	}
	alphaOut := r.AlphaNE
	return CoalitionResult{
		M:         M,
		AlphaCoal: alphaCoal,
		AlphaOut:  alphaOut,
		WedgeCoal: alphaCoal - r.AlphaCO,
	}
}

// PigouvianTax returns τ* and the corrected equilibrium.
func PigouvianTax(p Params) (tau, alphaCorrected float64) {
	r := Compute(p)
	return r.TaxStar, r.AlphaTaxed
}

// UBIEffect shows the effect of adding dA to autonomous demand.
func UBIEffect(p Params, dA float64) Result {
	pUBI := p
	pUBI.A += dA
	return Compute(pUBI)
}

var _ = math.Abs
