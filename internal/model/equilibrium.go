package model

import "math"

// Result holds all equilibrium outputs for a given Params.
type Result struct {
	Params Params

	// Derived scalars
	S     float64 // cost saving per task
	Ell   float64 // demand loss per automated task
	NStar float64 // externality activation threshold

	// Equilibrium
	AlphaNE  float64 // Nash equilibrium automation rate
	AlphaCO  float64 // cooperative optimum
	AlphaSP  float64 // social planner optimum at Params.Mu
	Wedge    float64 // AlphaNE - AlphaCO
	WedgePct float64 // wedge as % above cooperative (NaN when AlphaCO=0)

	// Welfare
	OwnerSurplus float64 // K (total firm profits at NE)
	WorkerIncome float64 // W at NE
	SocialWelfare float64 // S(mu) at NE
	OwnerSurplusCO float64
	WorkerIncomeCO float64
	SurplusLoss    float64 // K_CO - K_NE (deadweight loss)

	// Pigouvian tax
	TaxStar    float64 // τ* = ℓ(1-1/N)
	AlphaTaxed float64 // equilibrium under τ*

	// Externality active?
	ExternalityActive bool
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// nashRate computes the Nash automation rate given effective cost-saving s and
// demand-loss-per-firm ell_n = ell/N, and friction k.
func nashRate(s, ellN, k float64) float64 {
	if k == 0 {
		// Frictionless: binary outcome
		if s > ellN {
			return 1
		}
		return 0
	}
	return clamp((s-ellN)/k, 0, 1)
}

// Compute fills in all equilibrium and welfare fields for p.
func Compute(p Params) Result {
	r := Result{Params: p}
	r.S = p.S()
	r.Ell = p.Ell()
	r.NStar = p.NStar()
	r.ExternalityActive = p.N > r.NStar

	// Equilibrium rates
	r.AlphaNE = nashRate(r.S, r.Ell/p.N, p.K)
	if p.K == 0 {
		if r.S > r.Ell {
			r.AlphaCO = 1
		} else {
			r.AlphaCO = 0
		}
	} else {
		r.AlphaCO = clamp((r.S-r.Ell)/p.K, 0, 1)
	}

	// Social planner: S(µ) = µW + (1-µ)K
	// dS/dα = (1-µ)·dK/dα + µ·dW/dα
	// dK/dα = L(s - ℓ(1-1/N)·? ... actually for planner choosing common α:
	// K = N·π(α) where π(α) = Π0 + L(sα - ℓα - k/2·α²)
	// dK/dα = NL(s - ℓ - kα)
	// W = wLN(1-(1-η)α)  →  dW/dα = -wLN(1-η) = -LN·ℓ/λ ... wait
	// Actually W = wLN[1-(1-η)α], dW/dα = -wLN(1-η)
	// S(µ) first-order: (1-µ)NL(s-ℓ-kα) + µ(-wLN(1-η)) = 0
	// (1-µ)(s-ℓ-kα) = µ·w(1-η)
	// kα = (1-µ)(s-ℓ) - µ·w(1-η)
	// α_SP = [(1-µ)(s-ℓ) - µ·w(1-η)] / k
	if p.K == 0 {
		r.AlphaSP = r.AlphaCO // degenerate
	} else {
		wOneMinusEta := p.W * (1 - p.Eta)
		r.AlphaSP = clamp(((1-p.Mu)*(r.S-r.Ell)-p.Mu*wOneMinusEta)/p.K, 0, 1)
	}

	r.Wedge = r.AlphaNE - r.AlphaCO
	if r.AlphaCO > 0 {
		r.WedgePct = r.Wedge / r.AlphaCO * 100
	} else {
		r.WedgePct = math.NaN()
	}

	// Welfare at NE
	r.OwnerSurplus = ownerSurplus(p, r.AlphaNE)
	r.WorkerIncome = workerIncome(p, r.AlphaNE)
	r.SocialWelfare = (1-p.Mu)*r.OwnerSurplus + p.Mu*r.WorkerIncome

	// Welfare at CO
	r.OwnerSurplusCO = ownerSurplus(p, r.AlphaCO)
	r.WorkerIncomeCO = workerIncome(p, r.AlphaCO)

	r.SurplusLoss = r.OwnerSurplusCO - r.OwnerSurplus

	// Pigouvian tax
	r.TaxStar = r.Ell * (1 - 1/p.N)
	// Under tax τ, firm's effective cost saving is s-τ, demand loss per firm still ℓ/N
	r.AlphaTaxed = nashRate(r.S-r.TaxStar, r.Ell/p.N, p.K)

	return r
}

// ownerSurplus computes aggregate firm profits K at common automation rate α.
// K = N·[Π0 + L(sα - ℓα - k/2·α²)]
// where Π0 = A/N + (λ-1)wL
func ownerSurplus(p Params, alpha float64) float64 {
	pi0 := p.A/p.N + (p.Lambda-1)*p.W*p.L
	perFirm := pi0 + p.L*(p.S()*alpha-p.Ell()*alpha-p.K/2*alpha*alpha)
	return p.N * perFirm
}

// workerIncome computes W = wLN[1-(1-η)α]
func workerIncome(p Params, alpha float64) float64 {
	return p.W * p.L * p.N * (1 - (1-p.Eta)*alpha)
}

// WedgeFormula returns the closed-form wedge ℓ(1-1/N)/k (interior case).
func WedgeFormula(p Params) float64 {
	if p.K == 0 {
		return math.Inf(1)
	}
	return p.Ell() * (1 - 1/p.N) / p.K
}
