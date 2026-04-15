package model

import "math"

// PhiResult holds the result of the AI productivity extension (φ > 1, Prop 6).
// When φ > 1, each automated task produces φ units of output instead of 1,
// raising firm revenue but also distorting the equilibrium upward (Red Queen effect).
type PhiResult struct {
	Phi      float64
	AlphaNE  float64
	AlphaCO  float64
	Wedge    float64
	BaseWedge float64
	WedgeChange float64 // Wedge - BaseWedge (positive = amplified)
}

// PhiEquilibrium computes the equilibrium with AI productivity φ > 1.
// With φ, firm i's output is Yi = L·(1 + (φ-1)·αi), revenue = p·Yi.
// The equilibrium automation rate satisfies a quadratic (Prop 6).
// Returns the positive root.
func PhiEquilibrium(p Params) PhiResult {
	base := Compute(p)
	phi := p.Phi
	if math.Abs(phi-1) < 1e-10 {
		return PhiResult{
			Phi: phi, AlphaNE: base.AlphaNE, AlphaCO: base.AlphaCO,
			Wedge: base.Wedge, BaseWedge: base.Wedge, WedgeChange: 0,
		}
	}

	// With φ: firm i revenue = D(ᾱ)/N · φ·L·αi + D(ᾱ)/N·(1-αi)·L ... actually
	// Yi = L·(1+(φ-1)αi), total supply = NL(1+(φ-1)ᾱ), price p = D/supply
	// Rev_i = p·Yi = D·(1+(φ-1)αi)/(N·(1+(φ-1)ᾱ))
	// At symmetric equilibrium αi = ᾱ = α:
	// πi = D/N - Ci  (price·output - cost, the φ gain appears in D through output)
	//
	// Marginal condition (Prop 6): derivative w.r.t. αi at symmetric α gives quadratic.
	// Simplified from paper: define δ = φ-1, then:
	// k·α² - (s + δ·(A/NL + λw) - ℓ/N)·α + ... = 0
	// We solve numerically via the closed-form expression from the paper.
	//
	// For simplicity, we solve the FOC numerically (bisection) since the exact
	// quadratic coefficients require working through the full derivative carefully.

	foc := func(alpha float64) float64 {
		// D(α) = A + λwLN(1-(1-η)α)
		D := p.A + p.Lambda*p.W*p.L*p.N*(1-(1-p.Eta)*alpha)
		totalSupply := p.N * p.L * (1 + (phi-1)*alpha)
		price := D / totalSupply
		// dRev_i/dα_i at α_i = α (symmetric):
		// Rev_i = price * L * (1 + (phi-1)*alpha)
		// dRev/dα = L*(phi-1)*price + L*(1+(phi-1)*alpha) * dp/dα_i
		// dp/dα_i = (dD/dα_i * totalSupply - D * dSupply/dα_i) / totalSupply²
		// dD/dα_i = -ℓL (own contribution to demand loss)
		// dSupply/dα_i = L*(phi-1) (own contribution to supply)
		dD := -p.Ell() * p.L
		dSupply := p.L * (phi - 1)
		dPrice := (dD*totalSupply - D*dSupply) / (totalSupply * totalSupply)
		dRev := p.L*(phi-1)*price + p.L*(1+(phi-1)*alpha)*dPrice
		dCost := p.L * (p.S() - p.K*alpha) // from cost function, negated: dC/dα = -s + kα → dπ/dα via cost = s-kα
		// dπ/dα = dRev - d(Cost)/dα; dCost/dα = L(-s + kα) → dπ from cost = L(s-kα)
		_ = dCost
		return dRev + p.L*(p.S()-p.K*alpha)
	}

	// Bisect for zero of FOC in [0,1]
	lo, hi2 := 0.0, 1.0
	fLo := foc(lo)
	var alphaNE float64
	if fLo <= 0 {
		alphaNE = 0
	} else if foc(hi2) >= 0 {
		alphaNE = 1
	} else {
		for i := 0; i < 60; i++ {
			mid := (lo + hi2) / 2
			if foc(mid) > 0 {
				lo = mid
			} else {
				hi2 = mid
			}
		}
		alphaNE = (lo + hi2) / 2
	}

	// Cooperative: planner maximises K = N·π(α), same bisection
	focCO := func(alpha float64) float64 {
		D := p.A + p.Lambda*p.W*p.L*p.N*(1-(1-p.Eta)*alpha)
		totalSupply := p.N * p.L * (1 + (phi-1)*alpha)
		price := D / totalSupply
		// Planner sees full demand effect: dD/dᾱ = -ℓLN
		dD := -p.Ell() * p.L * p.N
		dSupply := p.L * p.N * (phi - 1)
		dPrice := (dD*totalSupply - D*dSupply) / (totalSupply * totalSupply)
		dRev := p.L*(phi-1)*price + p.L*(1+(phi-1)*alpha)*dPrice
		return dRev + p.L*(p.S()-p.K*alpha)
	}
	loC, hiC := 0.0, 1.0
	var alphaCO float64
	if focCO(loC) <= 0 {
		alphaCO = 0
	} else if focCO(hiC) >= 0 {
		alphaCO = 1
	} else {
		for i := 0; i < 60; i++ {
			mid := (loC + hiC) / 2
			if focCO(mid) > 0 {
				loC = mid
			} else {
				hiC = mid
			}
		}
		alphaCO = (loC + hiC) / 2
	}

	wedge := alphaNE - alphaCO
	return PhiResult{
		Phi: phi, AlphaNE: alphaNE, AlphaCO: alphaCO,
		Wedge: wedge, BaseWedge: base.Wedge, WedgeChange: wedge - base.Wedge,
	}
}

// EndogenousWages computes the fixed-point equilibrium when w adjusts with ᾱ.
// w(ᾱ) = w0 + slope·ᾱ  (from Section 5.3).
// Iterates to convergence (monotone contraction).
type EndogenousWagesResult struct {
	AlphaNE   float64
	WageEq    float64
	Iterations int
}

func EndogenousWages(p Params, slope float64) EndogenousWagesResult {
	w0 := p.W
	alpha := Compute(p).AlphaNE // initial guess
	const tol = 1e-10
	const maxIter = 1000
	for i := 0; i < maxIter; i++ {
		q := p
		q.W = w0 + slope*alpha
		if q.W < q.C {
			q.W = q.C
		}
		newAlpha := Compute(q).AlphaNE
		if math.Abs(newAlpha-alpha) < tol {
			return EndogenousWagesResult{AlphaNE: newAlpha, WageEq: q.W, Iterations: i + 1}
		}
		alpha = newAlpha
	}
	return EndogenousWagesResult{AlphaNE: alpha, WageEq: w0 + slope*alpha, Iterations: maxIter}
}

// CapitalRecycling computes the equilibrium when capital income is recycled to
// workers at rate ηHat (Section 5.4). Recycling raises effective replacement η_eff.
func CapitalRecycling(p Params, etaHat float64) Result {
	// Owner profit per firm = π_NE. Capital income recycled = etaHat * K / (NL workers).
	// This raises η effectively. We iterate.
	const tol = 1e-10
	const maxIter = 500
	q := p
	for i := 0; i < maxIter; i++ {
		r := Compute(q)
		// capital recycled per worker per task
		recyclePerWorker := etaHat * r.OwnerSurplus / (p.N * p.L * (1 - (1-p.Eta)*r.AlphaNE+1e-12))
		// effective eta: additional fraction of w replaced
		etaEff := math.Min(1.0, p.Eta + recyclePerWorker/p.W)
		qNew := p
		qNew.Eta = etaEff
		newR := Compute(qNew)
		if math.Abs(newR.AlphaNE-r.AlphaNE) < tol {
			return newR
		}
		q = qNew
	}
	return Compute(q)
}

// EndogenousEntry finds the free-entry N given per-firm entry cost kappa.
// Returns the equilibrium N (may be non-integer; floor/ceil for integer markets).
type EntryResult struct {
	NEq      float64 // equilibrium number of firms
	AlphaNE  float64
	Regime   string // "full-automation", "entry-deterrence", "high-cost"
}

func EndogenousEntry(p Params, kappa float64) EntryResult {
	// Free entry: π_NE(N) = kappa
	// π_NE = Π0 + L(s·α - ℓ·α - k/2·α²) where α = α_NE(N)
	// We search N in [1, 200].
	profit := func(n float64) float64 {
		q := p; q.N = n
		r := Compute(q)
		return ownerSurplus(q, r.AlphaNE) / n // per-firm profit
	}

	// Find N where profit = kappa (bisect)
	lo, hi2 := 1.0, 200.0
	pLo := profit(lo) - kappa
	pHi := profit(hi2) - kappa

	var regime string
	var nEq float64
	if pLo < 0 {
		// Even monopolist can't cover entry cost
		nEq = 0
		regime = "high-cost"
	} else if pHi > 0 {
		// Even at N=200 profits exceed kappa: entry keeps happening
		nEq = hi2
		regime = "entry-deterrence"
	} else {
		for i := 0; i < 60; i++ {
			mid := (lo + hi2) / 2
			if profit(mid)-kappa > 0 {
				lo = mid
			} else {
				hi2 = mid
			}
		}
		nEq = (lo + hi2) / 2
		regime = "equilibrium"
		q := p; q.N = nEq
		if Compute(q).AlphaNE > 0.999 {
			regime = "full-automation"
		}
	}

	q := p; q.N = math.Max(1, nEq)
	r := Compute(q)
	return EntryResult{NEq: nEq, AlphaNE: r.AlphaNE, Regime: regime}
}
