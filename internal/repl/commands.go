package repl

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aisim/internal/model"
	"aisim/internal/output"
)

// dispatch handles a single parsed command line.
func (r *REPL) dispatch(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return true
	}
	parts := strings.Fields(line)
	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case "quit", "exit", "q":
		fmt.Fprintln(r.out, "Goodbye.")
		return false

	case "help":
		r.printHelp(args)

	case "show":
		r.cmdShow()

	case "run":
		r.cmdRun()

	case "set":
		r.cmdSet(args)

	case "welfare":
		r.cmdWelfare(args)

	case "decompose":
		r.cmdDecompose()

	case "policy":
		r.cmdPolicy(args)

	case "sweep":
		r.cmdSweep(args)

	case "heatmap":
		r.cmdHeatmap(args)

	case "extend":
		r.cmdExtend(args)

	case "frictionless":
		r.params.K = 0
		fmt.Fprintln(r.out, "k set to 0 (Prisoner's Dilemma mode)")
		r.cmdRun()

	case "prisoners":
		r.cmdPrisoners()

	case "sensitivity":
		entries := model.Sensitivity(r.params)
		output.SensitivityTable(r.out, entries)

	case "save":
		if len(args) < 1 {
			fmt.Fprintln(r.out, "usage: save <name>")
		} else {
			r.saveScenario(args[0])
		}

	case "load":
		if len(args) < 1 {
			fmt.Fprintln(r.out, "usage: load <name>")
		} else {
			r.loadScenario(args[0])
		}

	case "diff":
		if len(args) < 2 {
			fmt.Fprintln(r.out, "usage: diff <name1> <name2>")
		} else {
			r.diffScenarios(args[0], args[1])
		}

	case "scenarios":
		r.listScenarios()

	case "preset":
		r.cmdPreset(args)

	case "export":
		r.cmdExport(args)

	case "??", "explain":
		r.cmdExplain(args)

	default:
		fmt.Fprintf(r.out, "unknown command %q; type 'help' for commands\n", cmd)
	}
	return true
}

func (r *REPL) cmdShow() {
	p := r.params
	fmt.Fprintf(r.out, "\nCurrent parameters:\n")
	fmt.Fprintf(r.out, "  N      = %.4g\n", p.N)
	fmt.Fprintf(r.out, "  w      = %.4f\n", p.W)
	fmt.Fprintf(r.out, "  c      = %.4f\n", p.C)
	fmt.Fprintf(r.out, "  k      = %.4f\n", p.K)
	fmt.Fprintf(r.out, "  lambda = %.4f\n", p.Lambda)
	fmt.Fprintf(r.out, "  eta    = %.4f\n", p.Eta)
	fmt.Fprintf(r.out, "  A      = %.4f\n", p.A)
	fmt.Fprintf(r.out, "  L      = %.4f\n", p.L)
	fmt.Fprintf(r.out, "  mu     = %.4f\n", p.Mu)
	fmt.Fprintf(r.out, "  phi    = %.4f\n", p.Phi)
	fmt.Fprintf(r.out, "  etahat = %.4f\n", p.EtaHat)
	fmt.Fprintf(r.out, "  s      = %.4f  (derived: w-c)\n", p.S())
	fmt.Fprintf(r.out, "  ell    = %.4f  (derived: λ(1-η)w)\n", p.Ell())
	fmt.Fprintln(r.out)
}

func (r *REPL) cmdRun() {
	if err := r.params.Validate(); err != nil {
		fmt.Fprintf(r.out, "invalid params: %v\n", err)
		return
	}
	res := model.Compute(r.params)
	output.Report(r.out, res)
}

func (r *REPL) cmdSet(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(r.out, "usage: set <param> <value>")
		return
	}
	v, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Fprintf(r.out, "invalid value %q: %v\n", args[1], err)
		return
	}
	p := r.params
	switch strings.ToLower(args[0]) {
	case "n":
		p.N = v
	case "w":
		p.W = v
	case "c":
		p.C = v
	case "k":
		p.K = v
	case "lambda":
		p.Lambda = v
	case "eta":
		p.Eta = v
	case "a":
		p.A = v
	case "l":
		p.L = v
	case "mu":
		p.Mu = v
	case "phi":
		p.Phi = v
	case "etahat":
		p.EtaHat = v
	default:
		fmt.Fprintf(r.out, "unknown param %q; valid: N w c k lambda eta A L mu phi etahat\n", args[0])
		return
	}
	if err := p.Validate(); err != nil {
		fmt.Fprintf(r.out, "invalid: %v\n", err)
		return
	}
	r.params = p
	fmt.Fprintf(r.out, "set %s = %.4g\n", args[0], v)
}

func (r *REPL) cmdWelfare(args []string) {
	p := r.params
	if len(args) > 0 {
		// parse mu=0.5 or just 0.5
		s := strings.TrimPrefix(args[0], "mu=")
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			p.Mu = v
		}
	}
	res := model.Compute(p)
	kRel, wRel := model.WelfareRelative(res)
	fmt.Fprintf(r.out, "\nWelfare at µ=%.2f\n", p.Mu)
	fmt.Fprintf(r.out, "  α_SP (social planner)  = %.4f\n", res.AlphaSP)
	fmt.Fprintf(r.out, "  K (owner surplus, NE)  = %.4f\n", res.OwnerSurplus)
	fmt.Fprintf(r.out, "  W (worker income, NE)  = %.4f\n", res.WorkerIncome)
	fmt.Fprintf(r.out, "  K/K_CO                 = %.4f\n", kRel)
	fmt.Fprintf(r.out, "  W/W_CO                 = %.4f\n", wRel)
	fmt.Fprintf(r.out, "  Surplus loss (K_CO-K)  = %.4f\n", res.SurplusLoss)
	fmt.Fprintln(r.out)
}

func (r *REPL) cmdDecompose() {
	demandExt, distrib := model.Decompose(r.params)
	res := model.Compute(r.params)
	fmt.Fprintf(r.out, "\nWedge decomposition\n")
	fmt.Fprintf(r.out, "  Total wedge (α_NE - α_CO)   = %.4f\n", res.Wedge)
	fmt.Fprintf(r.out, "  Demand externality component = %.4f  (α_NE - α_CO at µ=0)\n", demandExt)
	fmt.Fprintf(r.out, "  Distributional premium       = %.4f  (α_CO - α_SP at µ=%.2f)\n", distrib, r.params.Mu)
	fmt.Fprintln(r.out)
}

func (r *REPL) cmdPolicy(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(r.out, "usage: policy <tax|ubi <dA>|equity <ε>|upskill <η>|coalition <M>|capital-tax|all>")
		return
	}
	switch strings.ToLower(args[0]) {
	case "tax":
		tau, alphaTaxed := model.PigouvianTax(r.params)
		res := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nPigouvian tax\n")
		fmt.Fprintf(r.out, "  τ*               = %.4f  (= ℓ(1−1/N) = %.4f·%.4f)\n",
			tau, res.Ell, 1-1/r.params.N)
		fmt.Fprintf(r.out, "  α_NE (no tax)    = %.4f\n", res.AlphaNE)
		fmt.Fprintf(r.out, "  α_CO             = %.4f\n", res.AlphaCO)
		fmt.Fprintf(r.out, "  α under τ*       = %.4f  (implements cooperative)\n", alphaTaxed)
		fmt.Fprintln(r.out)

	case "ubi":
		dA := 1.0
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil {
				dA = v
			}
		}
		res := model.UBIEffect(r.params, dA)
		base := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nUBI (dA=%.4f)\n", dA)
		fmt.Fprintf(r.out, "  α_NE: %.4f → %.4f  (unchanged — automation incentive unaffected)\n",
			base.AlphaNE, res.AlphaNE)
		fmt.Fprintf(r.out, "  W (worker income): %.4f → %.4f\n", base.WorkerIncome, res.WorkerIncome)
		fmt.Fprintln(r.out)

	case "equity":
		eps := 0.1
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil {
				eps = v
			}
		}
		// Worker equity raises effective η: η_eff ≈ η + ε*(K/NE)/(w*L*N)
		res := model.Compute(r.params)
		etaEff := r.params.Eta + eps*(res.OwnerSurplus/r.params.N)/(r.params.W*r.params.L*r.params.N)
		etaEff = math.Min(etaEff, 1.0)
		pEq := r.params; pEq.Eta = etaEff
		resEq := model.Compute(pEq)
		fmt.Fprintf(r.out, "\nWorker equity (ε=%.4f)\n", eps)
		fmt.Fprintf(r.out, "  Effective η: %.4f → %.4f\n", r.params.Eta, etaEff)
		fmt.Fprintf(r.out, "  α_NE: %.4f → %.4f\n", res.AlphaNE, resEq.AlphaNE)
		fmt.Fprintf(r.out, "  Wedge: %.4f → %.4f  (narrows but cannot close)\n", res.Wedge, resEq.Wedge)
		fmt.Fprintln(r.out)

	case "upskill":
		newEta := math.Min(1.0, r.params.Eta+0.20)
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil {
				newEta = v
			}
		}
		pUp := r.params; pUp.Eta = newEta
		base := model.Compute(r.params)
		resUp := model.Compute(pUp)
		fmt.Fprintf(r.out, "\nUpskilling\n")
		fmt.Fprintf(r.out, "  η: %.4f → %.4f\n", r.params.Eta, newEta)
		fmt.Fprintf(r.out, "  ℓ: %.4f → %.4f\n", base.Ell, resUp.Ell)
		fmt.Fprintf(r.out, "  α_NE: %.4f → %.4f\n", base.AlphaNE, resUp.AlphaNE)
		fmt.Fprintf(r.out, "  Wedge: %.4f → %.4f\n", base.Wedge, resUp.Wedge)
		fmt.Fprintln(r.out)

	case "coalition":
		M := r.params.N / 2
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil {
				M = v
			}
		}
		cr := model.Coalition(r.params, M)
		base := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nPartial coalition (M=%.4g of N=%.4g)\n", M, r.params.N)
		fmt.Fprintf(r.out, "  Coalition α  = %.4f\n", cr.AlphaCoal)
		fmt.Fprintf(r.out, "  Outsiders α  = %.4f\n", cr.AlphaOut)
		fmt.Fprintf(r.out, "  Full Nash α  = %.4f\n", base.AlphaNE)
		fmt.Fprintf(r.out, "  Wedge vs CO  = %.4f\n", cr.WedgeCoal)
		fmt.Fprintln(r.out)

	case "capital-tax":
		res := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nCapital income tax\n")
		fmt.Fprintf(r.out, "  α_NE: %.4f (unchanged)\n", res.AlphaNE)
		fmt.Fprintf(r.out, "  Note: tax on profits does not affect per-task automation margin.\n")
		fmt.Fprintln(r.out)

	case "all":
		policies := model.AllPolicies(r.params)
		output.PolicyTable(r.out, policies)

	default:
		fmt.Fprintf(r.out, "unknown policy %q\n", args[0])
	}
}

func (r *REPL) cmdSweep(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(r.out, "usage: sweep <param> <from> <to> [steps=20]")
		return
	}
	param := args[0]
	lo, err1 := strconv.ParseFloat(args[1], 64)
	hi, err2 := strconv.ParseFloat(args[2], 64)
	if err1 != nil || err2 != nil {
		fmt.Fprintln(r.out, "invalid range values")
		return
	}
	steps := 20
	if len(args) > 3 {
		if v, err := strconv.Atoi(args[3]); err == nil {
			steps = v
		}
	}
	pts, err := model.Sweep1D(r.params, param, lo, hi, steps)
	if err != nil {
		fmt.Fprintf(r.out, "sweep error: %v\n", err)
		return
	}
	output.SweepTable(r.out, pts)
}

func (r *REPL) cmdHeatmap(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(r.out, "usage: heatmap <param1> <param2>  (uses default ranges per param)")
		return
	}
	p1, p2 := args[0], args[1]
	lo1, hi1 := defaultRange(p1, r.params)
	lo2, hi2 := defaultRange(p2, r.params)
	grid, err := model.Sweep2D(r.params, p1, lo1, hi1, 10, p2, lo2, hi2, 12)
	if err != nil {
		fmt.Fprintf(r.out, "heatmap error: %v\n", err)
		return
	}
	output.Heatmap2D(r.out, grid)
}

func defaultRange(param string, p model.Params) (lo, hi float64) {
	switch param {
	case "N":
		return 2, 20
	case "s":
		return 0.05, p.W * 0.95
	case "w":
		return 0.5, 2.0
	case "c":
		return 0, p.W * 0.9
	case "k":
		return 0.1, 3.0
	case "lambda":
		return 0.1, 1.0
	case "eta":
		return 0, 1.0
	default:
		return 0, 1
	}
}

func (r *REPL) cmdExtend(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(r.out, "usage: extend <phi <val>|wages <slope>|recycling <etahat>|entry <kappa>>")
		return
	}
	switch strings.ToLower(args[0]) {
	case "phi":
		phi := 1.2
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil { phi = v }
		}
		p := r.params; p.Phi = phi
		pr := model.PhiEquilibrium(p)
		base := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nAI productivity extension (φ=%.4f, Red Queen)\n", phi)
		fmt.Fprintf(r.out, "  Base wedge    = %.4f\n", pr.BaseWedge)
		fmt.Fprintf(r.out, "  φ wedge       = %.4f\n", pr.Wedge)
		fmt.Fprintf(r.out, "  Wedge change  = %+.4f  (positive = amplified)\n", pr.WedgeChange)
		fmt.Fprintf(r.out, "  α_NE (φ=1)   = %.4f → α_NE (φ=%.4f) = %.4f\n",
			base.AlphaNE, phi, pr.AlphaNE)
		fmt.Fprintln(r.out)

	case "wages":
		slope := -0.5
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil { slope = v }
		}
		wr := model.EndogenousWages(r.params, slope)
		base := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nEndogenous wages (w(ᾱ) = w₀ + %.4f·ᾱ)\n", slope)
		fmt.Fprintf(r.out, "  Equilibrium w = %.4f  (base w = %.4f)\n", wr.WageEq, r.params.W)
		fmt.Fprintf(r.out, "  α_NE          = %.4f → %.4f\n", base.AlphaNE, wr.AlphaNE)
		fmt.Fprintf(r.out, "  Iterations    = %d\n", wr.Iterations)
		fmt.Fprintln(r.out)

	case "recycling":
		etaHat := 0.5
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil { etaHat = v }
		}
		res := model.CapitalRecycling(r.params, etaHat)
		base := model.Compute(r.params)
		fmt.Fprintf(r.out, "\nCapital income recycling (η̂=%.4f)\n", etaHat)
		fmt.Fprintf(r.out, "  α_NE: %.4f → %.4f\n", base.AlphaNE, res.AlphaNE)
		fmt.Fprintf(r.out, "  Wedge: %.4f → %.4f\n", base.Wedge, res.Wedge)
		fmt.Fprintln(r.out)

	case "entry":
		kappa := 5.0
		if len(args) > 1 {
			if v, err := strconv.ParseFloat(args[1], 64); err == nil { kappa = v }
		}
		er := model.EndogenousEntry(r.params, kappa)
		fmt.Fprintf(r.out, "\nEndogenous entry (κ=%.4f)\n", kappa)
		fmt.Fprintf(r.out, "  Equilibrium N = %.4f  (regime: %s)\n", er.NEq, er.Regime)
		fmt.Fprintf(r.out, "  α_NE          = %.4f\n", er.AlphaNE)
		fmt.Fprintln(r.out)

	default:
		fmt.Fprintf(r.out, "unknown extension %q\n", args[0])
	}
}

func (r *REPL) cmdPrisoners() {
	p := r.params
	p.K = 0
	// PD payoff matrix: each firm chooses automate (1) or not (0)
	// Payoffs as multiples of baseline profit Π0
	automate := model.Compute(func() model.Params { q := p; return q }())
	noAuto := model.Compute(func() model.Params { q := p; q.N = 1; return q }())
	_ = noAuto

	// Compute 4 cells manually
	profitAt := func(alpha_i, alpha_j float64) float64 {
		// Single firm i's profit when it plays alpha_i and rival plays alpha_j
		// (2-firm PD for simplicity)
		avgAlpha := (alpha_i + alpha_j) / 2
		D := p.A + p.Lambda*p.W*p.L*2*(1-(1-p.Eta)*avgAlpha)
		price := D / (2 * p.L)
		rev := price * p.L
		cost := p.L * (p.W - (p.W-p.C)*alpha_i)
		return rev - cost
	}
	pp := profitAt(0, 0) // both don't automate
	pq := profitAt(1, 0) // i automates, j doesn't
	qp := profitAt(0, 1) // i doesn't, j does
	qq := profitAt(1, 1) // both automate

	fmt.Fprintf(r.out, "\nPrisoner's Dilemma payoff matrix (k=0, N=2, per firm)\n")
	fmt.Fprintf(r.out, "                    Rival: No Auto   Rival: Automate\n")
	fmt.Fprintf(r.out, "  Me: No Automate    %9.4f        %9.4f\n", pp, qp)
	fmt.Fprintf(r.out, "  Me: Automate       %9.4f        %9.4f\n", pq, qq)
	fmt.Fprintf(r.out, "\n  Nash: both automate (%.4f < %.4f dominates %.4f)\n",
		qq, pq, pp)
	fmt.Fprintf(r.out, "  Coop: both don't (%.4f > %.4f)\n", pp, qq)
	_ = automate
	fmt.Fprintln(r.out)
}

func (r *REPL) cmdPreset(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(r.out, "usage: preset <paper|frictionless|monopolist|duopoly|competitive|optimist>")
		return
	}
	switch strings.ToLower(args[0]) {
	case "paper":
		r.params = model.PaperBase()
	case "frictionless":
		r.params = model.Frictionless()
	case "monopolist":
		r.params = model.Monopolist()
	case "duopoly":
		r.params = model.Duopoly()
	case "competitive", "competitive-limit":
		r.params = model.CompetitiveLimit()
	case "optimist":
		r.params = model.Optimist()
	default:
		fmt.Fprintf(r.out, "unknown preset %q\n", args[0])
		return
	}
	fmt.Fprintf(r.out, "Loaded preset %q\n", args[0])
}

func (r *REPL) cmdExport(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(r.out, "usage: export <json|csv>")
		return
	}
	res := model.Compute(r.params)
	switch strings.ToLower(args[0]) {
	case "json":
		output.WriteJSON(r.out, res) //nolint
	case "csv":
		output.WriteResultCSV(r.out, res) //nolint
	default:
		fmt.Fprintf(r.out, "unknown format %q; use json or csv\n", args[0])
	}
}

func (r *REPL) printHelp(args []string) {
	if len(args) > 0 {
		// per-command help could be expanded; for now just print general help
	}
	fmt.Fprintln(r.out, `
Commands:
  set <param> <value>          Set a parameter (N w c k lambda eta A L mu phi etahat)
  show                         Show current parameters
  run                          Compute equilibrium, print full report
  welfare [mu=<µ>]             Welfare metrics at given planner weight
  decompose                    Split wedge into demand-externality + distributional
  policy tax                   Optimal Pigouvian tax τ*
  policy ubi [dA]              UBI effect (+dA autonomous demand)
  policy equity [ε]            Worker equity participation
  policy upskill [η]           Upskilling / raise η
  policy coalition [M]         Partial coalition of M firms
  policy capital-tax           Capital income tax (no-effect result)
  policy all                   Compare all instruments
  sweep <param> <from> <to>    1D sweep with sparkline
  heatmap <p1> <p2>            2D ASCII heatmap of wedge
  extend phi [φ]               AI productivity (Red Queen effect)
  extend wages [slope]         Endogenous wages w(ᾱ)=w₀+slope·ᾱ
  extend recycling [η̂]         Capital income recycling
  extend entry [κ]             Endogenous firm entry
  frictionless                 Set k=0, Prisoner's Dilemma mode
  prisoners                    Explicit PD payoff matrix (N=2, k=0)
  sensitivity                  ∂(wedge)/∂(param) for all params
  save <name>                  Save current params as named scenario
  load <name>                  Restore named scenario
  diff <name1> <name2>         Side-by-side scenario comparison
  scenarios                    List saved scenarios
  preset <name>                Load preset (paper|frictionless|monopolist|duopoly|competitive|optimist)
  export json|csv              Export current results
  ??  [param]                  Deep explanation of all params, or one param
  help                         This help text
  quit                         Exit`)
}
