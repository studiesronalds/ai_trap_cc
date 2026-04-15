package model

import "math"

// WelfareAt computes K, W, S(mu) at an arbitrary alpha (not necessarily NE).
func WelfareAt(p Params, alpha float64) (K, W, S float64) {
	K = ownerSurplus(p, alpha)
	W = workerIncome(p, alpha)
	S = (1-p.Mu)*K + p.Mu*W
	return
}

// Decompose splits the over-automation wedge into two components.
// demandExt  = α_NE − α_CO  (the pure demand externality at µ=0)
// distributional = α_CO − α_SP  (extra restraint wanted by a worker-weighting planner)
func Decompose(p Params) (demandExt, distributional float64) {
	r := Compute(p)
	demandExt = r.Wedge
	distributional = r.AlphaCO - r.AlphaSP
	return
}

// SurplusLossFraction returns (K_CO - K_NE) / K_CO.
func SurplusLossFraction(r Result) float64 {
	if r.OwnerSurplusCO == 0 {
		return math.NaN()
	}
	return r.SurplusLoss / r.OwnerSurplusCO
}

// WelfareRelative returns K_NE/K_CO and W_NE/W_CO (cooperative = 1.0).
func WelfareRelative(r Result) (kRel, wRel float64) {
	if r.OwnerSurplusCO != 0 {
		kRel = r.OwnerSurplus / r.OwnerSurplusCO
	} else {
		kRel = math.NaN()
	}
	if r.WorkerIncomeCO != 0 {
		wRel = r.WorkerIncome / r.WorkerIncomeCO
	} else {
		wRel = math.NaN()
	}
	return
}

// SensitivityEntry is one row of the sensitivity report.
type SensitivityEntry struct {
	Param  string
	Value  float64
	DWedge float64 // ∂wedge/∂param (central difference)
}

// Sensitivity returns ∂(wedge)/∂(each param) numerically for all core params.
func Sensitivity(p Params) []SensitivityEntry {
	wedge := func(q Params) float64 { return Compute(q).Wedge }

	entries := []struct {
		name string
		get  func(Params) float64
		perturb func(Params, float64) (Params, Params)
	}{
		{"N", func(q Params) float64 { return q.N },
			func(q Params, h float64) (Params, Params) {
				hi, lo := q, q; hi.N += h; lo.N = math.Max(1, lo.N-h); return hi, lo
			}},
		{"w", func(q Params) float64 { return q.W },
			func(q Params, h float64) (Params, Params) {
				hi, lo := q, q; hi.W += h; lo.W = math.Max(q.C, lo.W-h); return hi, lo
			}},
		{"c", func(q Params) float64 { return q.C },
			func(q Params, h float64) (Params, Params) {
				hi, lo := q, q
				hi.C = math.Min(q.W, hi.C+h); lo.C = math.Max(0, lo.C-h); return hi, lo
			}},
		{"k", func(q Params) float64 { return q.K },
			func(q Params, h float64) (Params, Params) {
				hi, lo := q, q; hi.K += h; lo.K = math.Max(0, lo.K-h); return hi, lo
			}},
		{"lambda", func(q Params) float64 { return q.Lambda },
			func(q Params, h float64) (Params, Params) {
				hi, lo := q, q
				hi.Lambda = math.Min(1, hi.Lambda+h); lo.Lambda = math.Max(0, lo.Lambda-h); return hi, lo
			}},
		{"eta", func(q Params) float64 { return q.Eta },
			func(q Params, h float64) (Params, Params) {
				hi, lo := q, q; hi.Eta += h; lo.Eta -= h; return hi, lo
			}},
	}

	out := make([]SensitivityEntry, 0, len(entries))
	for _, e := range entries {
		v := e.get(p)
		h := math.Abs(v) * 1e-4
		if h < 1e-8 {
			h = 1e-6
		}
		hi, lo := e.perturb(p, h)
		denom := e.get(hi) - e.get(lo)
		var dw float64
		if denom != 0 {
			dw = (wedge(hi) - wedge(lo)) / denom
		}
		out = append(out, SensitivityEntry{Param: e.name, Value: v, DWedge: dw})
	}
	return out
}
