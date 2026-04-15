package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"aisim/internal/model"
	"aisim/internal/output"
	"aisim/internal/repl"
)

func main() {
	// Parameter flags
	fN := flag.Float64("N", 7, "number of firms")
	fW := flag.Float64("w", 1.0, "wage per task")
	fC := flag.Float64("c", 0.30, "AI cost per task")
	fK := flag.Float64("k", 1.0, "integration friction")
	fLambda := flag.Float64("lambda", 0.5, "worker MPC")
	fEta := flag.Float64("eta", 0.30, "income replacement rate")
	fA := flag.Float64("A", 10.0, "autonomous demand")
	fL := flag.Float64("L", 100.0, "tasks per firm")
	fMu := flag.Float64("mu", 0.5, "social planner worker weight [0,1]")
	fPhi := flag.Float64("phi", 1.0, "AI productivity multiplier")

	// Mode flags
	fREPL := flag.Bool("repl", false, "start interactive REPL")
	fPreset := flag.String("preset", "", "load preset: paper|frictionless|monopolist|duopoly|competitive|optimist")

	// Command flags
	fPolicy := flag.String("policy", "", "policy instrument: tax|ubi|upskill|coalition|capital-tax|all")
	fSweep := flag.String("sweep", "", "sweep param(s): <param> or <p1,p2>")
	fSweepFrom := flag.Float64("sweep-from", 0, "sweep lower bound")
	fSweepTo := flag.Float64("sweep-to", 10, "sweep upper bound")
	fSweepSteps := flag.Int("sweep-steps", 20, "sweep steps")
	fSweep2 := flag.String("sweep2", "", "second param for 2D sweep")
	fSweep2From := flag.Float64("sweep2-from", 0, "second param lower bound")
	fSweep2To := flag.Float64("sweep2-to", 1, "second param upper bound")

	// Output format
	fOut := flag.String("output", "text", "output format: text|json|csv")

	flag.Parse()

	// Build params from flags or preset
	var p model.Params
	if *fPreset != "" {
		p = loadPreset(*fPreset)
	} else {
		p = model.PaperBase()
	}

	// Override with explicit flags (only if they were set)
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "N":      p.N = *fN
		case "w":      p.W = *fW
		case "c":      p.C = *fC
		case "k":      p.K = *fK
		case "lambda": p.Lambda = *fLambda
		case "eta":    p.Eta = *fEta
		case "A":      p.A = *fA
		case "L":      p.L = *fL
		case "mu":     p.Mu = *fMu
		case "phi":    p.Phi = *fPhi
		}
	})

	if err := p.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid parameters: %v\n", err)
		os.Exit(1)
	}

	// REPL mode: no args or --repl
	if *fREPL || flag.NFlag() == 0 {
		r := repl.New(p)
		r.Run()
		return
	}

	// Sweep mode
	if *fSweep != "" {
		runSweep(p, *fSweep, *fSweepFrom, *fSweepTo, *fSweepSteps,
			*fSweep2, *fSweep2From, *fSweep2To, *fOut)
		return
	}

	// Policy mode
	if *fPolicy != "" {
		runPolicy(p, *fPolicy, *fOut)
		return
	}

	// Default: full report
	res := model.Compute(p)
	switch strings.ToLower(*fOut) {
	case "json":
		if err := output.WriteJSON(os.Stdout, res); err != nil {
			fmt.Fprintf(os.Stderr, "json error: %v\n", err)
		}
	case "csv":
		if err := output.WriteResultCSV(os.Stdout, res); err != nil {
			fmt.Fprintf(os.Stderr, "csv error: %v\n", err)
		}
	default:
		output.Report(os.Stdout, res)
	}
}

func loadPreset(name string) model.Params {
	switch strings.ToLower(name) {
	case "paper":
		return model.PaperBase()
	case "frictionless":
		return model.Frictionless()
	case "monopolist":
		return model.Monopolist()
	case "duopoly":
		return model.Duopoly()
	case "competitive", "competitive-limit":
		return model.CompetitiveLimit()
	case "optimist":
		return model.Optimist()
	default:
		fmt.Fprintf(os.Stderr, "unknown preset %q; valid: paper|frictionless|monopolist|duopoly|competitive|optimist\n", name)
		os.Exit(1)
		return model.Params{}
	}
}

func runPolicy(p model.Params, pol, outFmt string) {
	if pol == "all" {
		policies := model.AllPolicies(p)
		output.PolicyTable(os.Stdout, policies)
		return
	}
	// Individual policy — run full report + highlight policy
	res := model.Compute(p)
	if outFmt == "text" {
		output.Report(os.Stdout, res)
	}
	switch strings.ToLower(pol) {
	case "tax":
		tau, alphaTaxed := model.PigouvianTax(p)
		fmt.Printf("Pigouvian tax τ* = %.4f  →  corrected α = %.4f\n", tau, alphaTaxed)
	case "ubi":
		fmt.Println("UBI: raises A, does not change α_NE or the wedge.")
	case "upskill":
		p2 := p; p2.Eta = 1.0
		r2 := model.Compute(p2)
		fmt.Printf("Upskilling to η=1.0: α_NE %.4f → %.4f, wedge %.4f → %.4f\n",
			res.AlphaNE, r2.AlphaNE, res.Wedge, r2.Wedge)
	case "capital-tax":
		fmt.Printf("Capital income tax: α_NE %.4f (unchanged)\n", res.AlphaNE)
	default:
		fmt.Fprintf(os.Stderr, "unknown policy %q\n", pol)
	}
}

func runSweep(p model.Params, param string, lo, hi float64, steps int,
	param2 string, lo2, hi2 float64, outFmt string) {

	if param2 != "" {
		// 2D sweep
		grid, err := model.Sweep2D(p, param, lo, hi, steps, param2, lo2, hi2, steps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sweep error: %v\n", err)
			os.Exit(1)
		}
		output.Heatmap2D(os.Stdout, grid)
		return
	}

	// 1D sweep
	pts, err := model.Sweep1D(p, param, lo, hi, steps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sweep error: %v\n", err)
		os.Exit(1)
	}
	switch strings.ToLower(outFmt) {
	case "json":
		output.WriteSweepJSON(os.Stdout, pts) //nolint
	case "csv":
		output.WriteCSV(os.Stdout, pts) //nolint
	default:
		output.SweepTable(os.Stdout, pts)
	}
}
