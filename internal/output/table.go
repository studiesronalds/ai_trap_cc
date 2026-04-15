package output

import (
	"fmt"
	"io"
	"math"
	"strings"

	"aisim/internal/model"
)

// Report prints the full equilibrium report to w.
func Report(w io.Writer, r model.Result) {
	p := r.Params
	kRel, wRel := model.WelfareRelative(r)
	demandExt, distrib := model.Decompose(p)

	sep := strings.Repeat("─", 57)
	fmt.Fprintf(w, "\n=== AI Layoff Trap Simulation ===\n")
	fmt.Fprintf(w, "Parameters: N=%.4g  w=%.2f  c=%.2f  k=%.2f  λ=%.2f  η=%.2f  A=%.4g  L=%.4g\n",
		p.N, p.W, p.C, p.K, p.Lambda, p.Eta, p.A, p.L)
	fmt.Fprintf(w, "%s\n", sep)

	fmt.Fprintf(w, "Derived\n")
	fmt.Fprintf(w, "  Cost saving s         = %.4f\n", r.S)
	fmt.Fprintf(w, "  Demand loss ℓ         = %.4f\n", r.Ell)
	if r.ExternalityActive {
		fmt.Fprintf(w, "  Threshold N*          = %.4f  → externality ACTIVE (N=%.4g > N*=%.4g)\n",
			r.NStar, p.N, r.NStar)
	} else {
		fmt.Fprintf(w, "  Threshold N*          = %.4f  → externality INACTIVE (N=%.4g ≤ N*=%.4g)\n",
			r.NStar, p.N, r.NStar)
	}

	fmt.Fprintf(w, "\nEquilibrium\n")
	fmt.Fprintf(w, "  Nash (α_NE)           = %.4f\n", r.AlphaNE)
	fmt.Fprintf(w, "  Cooperative (α_CO)    = %.4f\n", r.AlphaCO)
	fmt.Fprintf(w, "  Social planner (α_SP) = %.4f  (µ=%.2f)\n", r.AlphaSP, p.Mu)
	if math.IsNaN(r.WedgePct) {
		fmt.Fprintf(w, "  Over-automation wedge = %.4f\n", r.Wedge)
	} else {
		fmt.Fprintf(w, "  Over-automation wedge = %.4f   (+%.1f%% above cooperative)\n",
			r.Wedge, r.WedgePct)
	}

	fmt.Fprintf(w, "\nWedge decomposition\n")
	fmt.Fprintf(w, "  Demand externality    = %.4f\n", demandExt)
	fmt.Fprintf(w, "  Distributional (µ)    = %.4f\n", distrib)

	fmt.Fprintf(w, "\nWelfare  (relative to cooperative = 1.00)\n")
	if math.IsNaN(kRel) {
		fmt.Fprintf(w, "  Owner surplus  K      =   N/A\n")
	} else {
		fmt.Fprintf(w, "  Owner surplus  K      =  %.3f\n", kRel)
	}
	if math.IsNaN(wRel) {
		fmt.Fprintf(w, "  Worker income  W      =   N/A\n")
	} else {
		fmt.Fprintf(w, "  Worker income  W      =  %.3f\n", wRel)
	}
	if kRel < 1 && wRel < 1 {
		fmt.Fprintf(w, "  Both groups harmed — Pareto-dominated equilibrium confirmed.\n")
	}

	fmt.Fprintf(w, "\nOptimal Pigouvian tax   τ* = %.4f  (= ℓ(1−1/N))\n", r.TaxStar)
	fmt.Fprintf(w, "  Corrected α (under τ*) = %.4f\n", r.AlphaTaxed)
	fmt.Fprintf(w, "%s\n\n", sep)
}

// PolicyTable prints a comparison table of all policy instruments.
func PolicyTable(w io.Writer, policies []model.PolicyResult) {
	fmt.Fprintf(w, "\n%-22s  %-8s  %-8s  %s\n", "Instrument", "α after", "Wedge", "Note")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 72))
	for _, p := range policies {
		closes := " "
		if p.ClosesWedge {
			closes = "✓"
		}
		fmt.Fprintf(w, "%-22s  %-8.4f  %-8.4f  %s %s\n",
			p.Name, p.AlphaAfter, p.WedgeAfter, closes, p.Note)
	}
	fmt.Fprintln(w)
}

// SweepTable prints a 1D sweep as an aligned table with a sparkline.
func SweepTable(w io.Writer, pts []model.SweepPoint) {
	if len(pts) == 0 {
		return
	}
	param := pts[0].ParamName
	fmt.Fprintf(w, "\nSweep: %s\n", param)
	fmt.Fprintf(w, "%-10s  %-8s  %-8s  %-8s  %s\n", param, "α_NE", "α_CO", "Wedge", "Spark")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 52))

	// find wedge range for spark
	minW, maxW := pts[0].Result.Wedge, pts[0].Result.Wedge
	for _, pt := range pts {
		if pt.Result.Wedge < minW { minW = pt.Result.Wedge }
		if pt.Result.Wedge > maxW { maxW = pt.Result.Wedge }
	}
	spark := []rune("▁▂▃▄▅▆▇█")
	for _, pt := range pts {
		r := pt.Result
		var bar rune
		if maxW > minW {
			idx := int((r.Wedge - minW) / (maxW - minW) * float64(len(spark)-1))
			if idx < 0 { idx = 0 }
			if idx >= len(spark) { idx = len(spark) - 1 }
			bar = spark[idx]
		} else {
			bar = spark[0]
		}
		fmt.Fprintf(w, "%-10.4g  %-8.4f  %-8.4f  %-8.4f  %c\n",
			pt.Value, r.AlphaNE, r.AlphaCO, r.Wedge, bar)
	}
	fmt.Fprintln(w)
}

// SensitivityTable prints the sensitivity report.
func SensitivityTable(w io.Writer, entries []model.SensitivityEntry) {
	fmt.Fprintf(w, "\nSensitivity: ∂(wedge)/∂(param)\n")
	fmt.Fprintf(w, "%-10s  %-10s  %s\n", "Param", "Value", "∂wedge/∂param")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 36))
	for _, e := range entries {
		fmt.Fprintf(w, "%-10s  %-10.4f  %.6f\n", e.Param, e.Value, e.DWedge)
	}
	fmt.Fprintln(w)
}
