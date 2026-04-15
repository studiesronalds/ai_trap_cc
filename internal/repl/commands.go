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

// paramExplanation holds the structured explanation for one parameter.
type paramExplanation struct {
	symbol  string
	name    string
	current func(p model.Params) string
	body    string
}

var paramExplanations = []paramExplanation{
	{
		symbol: "N",
		name:   "Number of firms",
		current: func(p model.Params) string { return fmt.Sprintf("%.4g", p.N) },
		body: `  What it is:
    The count of symmetric competing firms in the industry. Every firm is
    identical — same wage bill, same task count, same AI cost structure.

  Why it matters:
    N is the engine of the externality. Each firm captures 100% of its own
    cost savings from automating a task, but the demand destruction that
    results is shared equally across all N firms. The bigger N is, the
    smaller each firm's share of the damage — so each has a stronger
    incentive to over-automate relative to what a cooperative planner
    would choose.

  Derived quantities that move with N:
    • Wedge = ℓ(1−1/N)/k  →  grows toward ℓ/k as N→∞, collapses to 0
      at N=1 (monopolist fully internalises the externality).
    • Optimal Pigouvian tax τ* = ℓ(1−1/N)  →  same pattern.
    • N* = ℓ/s  →  externality only activates when N > N*.

  Direction:
    ↑ N  →  larger wedge, higher α_NE, worse collective outcome.
    ↓ N  →  wedge shrinks; monopolist (N=1) means α_NE = α_CO.

  Paper base: N = 7`,
	},
	{
		symbol: "w",
		name:   "Wage per task",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.W) },
		body: `  What it is:
    The labour cost a firm pays per task when a human worker does it.
    Normalised to 1.0 in the paper base scenario — all other monetary
    quantities are expressed relative to this unit wage.

  Why it matters:
    w sets the gross cost-saving potential (s = w−c) and simultaneously
    drives the demand destruction channel: workers who lose tasks lose
    income w per task, of which they spend fraction λ(1−η) back into
    the economy as demand. A higher wage amplifies both effects.

  Direction:
    ↑ w  →  larger s (more incentive to automate) AND larger ℓ
            (more demand destruction). Net effect on wedge depends on
            which channel dominates; typically both α_NE and the wedge
            rise because the cost-saving margin grows faster than ℓ/k.
    ↓ w  →  both effects shrink together.

  Paper base: w = 1.00`,
	},
	{
		symbol: "c",
		name:   "AI cost per task",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.C) },
		body: `  What it is:
    The per-task cost of running AI instead of a human worker. Covers
    compute, licensing, integration labour, and ongoing maintenance —
    everything that replaces the wage w when a task is automated.
    Constraint: c ≤ w (otherwise automation is never privately rational).

  Why it matters:
    Together with w it determines the cost saving s = w−c. This is the
    private benefit that drives each firm's automation decision. Lowering
    c is the primary mechanism through which falling AI prices (GPT-4
    → commodity models) accelerate the externality.

  Direction:
    ↑ c  →  s shrinks → α_NE falls, wedge shrinks. High enough c
            makes automation privately unprofitable even without policy.
    ↓ c  →  s expands → α_NE rises toward 1, wedge widens. This is
            the "AI price crash" scenario.

  Paper base: c = 0.30  (s = 0.70)`,
	},
	{
		symbol: "k",
		name:   "Integration friction",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.K) },
		body: `  What it is:
    A convexity parameter capturing rising marginal costs of automation:
    reorganising workflows, retraining management, handling edge-cases,
    compliance overhead. Profit from automating fraction α is quadratic
    in α; k is the curvature coefficient.

  Why it matters:
    k is the denominator of every equilibrium expression. It is what
    keeps firms from jumping straight to α=1 even when s > 0.
    Setting k=0 removes all friction and collapses the model to a
    binary Prisoner's Dilemma (automate fully or not at all).

  Direction:
    ↑ k  →  α_NE and α_CO both shrink; wedge shrinks in absolute
            terms (Δα = ℓ(1−1/N)/k) but the gap remains proportional.
    ↓ k  →  automation accelerates; k→0 gives the frictionless PD mode.
    k = 0 →  use 'frictionless' or 'preset frictionless' to explore.

  Paper base: k = 1.00`,
	},
	{
		symbol: "lambda",
		name:   "Worker marginal propensity to consume (λ)",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.Lambda) },
		body: `  What it is:
    The fraction of labour income that workers immediately spend back
    into product markets. Ranges [0,1]. λ=0 means workers save
    everything; λ=1 means full spending passthrough (Keynesian limit).

  Why it matters:
    λ is one of two amplifiers of the demand destruction channel.
    When a task is automated, the worker loses wage w. Of that, they
    spend λ·(1−η)·w back as demand (net of replacement income η·w).
    So λ scales how hard the spending contraction hits aggregate demand.

  Direction:
    ↑ λ  →  larger ℓ → bigger wedge, larger optimal tax.
            High-λ workers (hand-to-mouth households) amplify the trap.
    ↓ λ  →  workers save more; demand destruction shrinks; wedge narrows.

  Paper base: λ = 0.50`,
	},
	{
		symbol: "eta",
		name:   "Income replacement rate (η)",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.Eta) },
		body: `  What it is:
    The fraction of lost wage income that displaced workers recover
    through transfers, benefits, or new employment. η=0 means no
    replacement; η=1 means full replacement (income-neutral displacement).
    η > 1 is possible if displacement triggers net income gains
    (e.g., very generous retraining stipends).

  Why it matters:
    η is the policy lever closest to the demand channel. It enters the
    demand-loss formula as ℓ = λ(1−η)w — raising η directly reduces
    ℓ, narrowing the wedge without touching firms' automation margins.
    This is the economic mechanism behind upskilling and retraining
    programmes: they work not by making AI more expensive but by
    rebuilding workers' spending power.

  Direction:
    ↑ η  →  ℓ shrinks → wedge narrows; at η=1, ℓ=0, wedge=0,
            externality fully neutralised (demand-channel only).
    ↓ η  →  more demand destruction per automated task; trap deepens.

  Paper base: η = 0.30
  See also: 'policy upskill' to model a targeted η increase.`,
	},
	{
		symbol: "A",
		name:   "Autonomous demand",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.A) },
		body: `  What it is:
    The component of product demand that is independent of worker
    income — government spending, export demand, investment, wealthy
    household consumption not linked to labour income. Units are
    the same as aggregate demand D.

  Why it matters:
    A sets the demand floor. In the model, total demand is:
      D = A + λ·w·L·N·[1 − (1−η)·ᾱ]
    where ᾱ is the industry average automation rate. A determines
    how large a fraction of demand can survive even if all labour
    income is destroyed. However, A does not affect α_NE or the
    wedge — firms' per-task automation incentives are independent
    of the demand level (they depend only on the margin s vs ℓ/N).

  Direction:
    ↑ A  →  higher welfare floor; workers are less harmed in absolute
            terms, but the over-automation trap is structurally unchanged.
    ↓ A  →  demand becomes more dependent on labour income; welfare
            losses from over-automation are larger in absolute terms.
    UBI effect: UBI raises A, which is why 'policy ubi' improves W
    but cannot close the wedge.

  Paper base: A = 10.0`,
	},
	{
		symbol: "L",
		name:   "Tasks per firm",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.L) },
		body: `  What it is:
    The number of separable tasks each firm can potentially automate.
    α·L tasks get automated; (1−α)·L tasks remain human-performed.
    L scales the total labour income in the economy (w·L·N total).

  Why it matters:
    L is a scale parameter. It enters welfare and demand calculations
    but cancels out of equilibrium conditions for α — the per-task
    margin that determines α_NE and α_CO does not depend on how many
    tasks there are. Changing L shifts absolute welfare numbers but
    not the equilibrium automation rates or the wedge.

  Direction:
    ↑ L  →  larger economy (more total income and demand); equilibrium
            α values unchanged; welfare gains/losses scale proportionally.
    ↓ L  →  smaller economy; same structural trap, smaller absolute stakes.

  Paper base: L = 100`,
	},
	{
		symbol: "mu",
		name:   "Social planner worker weight (µ)",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.Mu) },
		body: `  What it is:
    The weight the social planner places on worker welfare relative to
    owner (capital) welfare when computing the socially optimal automation
    rate α_SP. µ=0 means the planner only cares about profits; µ=1 means
    the planner only cares about workers; µ=0.5 (paper base) is a
    utilitarian 50/50 split.

  Why it matters:
    µ determines α_SP via:
      α_SP = max(0, (s − ℓ/N − µ·ℓ(1−1/N)) / k)
    It decompose the wedge: α_NE − α_SP splits into the externality
    component (µ-independent) and a distributional premium that grows
    with µ (a higher-µ planner wants even less automation because they
    weight worker income losses more heavily).

  Direction:
    ↑ µ  →  α_SP falls; social planner wants less automation.
            The wedge α_NE − α_SP widens.
    ↓ µ  →  α_SP rises toward α_CO; the distributional premium vanishes.
    µ = 0 →  α_SP = α_CO (planner ignores distributional concerns).
    µ = 1 →  α_SP = 0 in many parameterisations (maximum restriction).

  Paper base: µ = 0.50`,
	},
	{
		symbol: "phi",
		name:   "AI productivity multiplier (φ) — extension",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.Phi) },
		body: `  What it is:
    A multiplier on output per automated task, capturing the possibility
    that AI does the task better (faster, higher quality) than a human.
    φ=1 (paper base) means AI and human are equally productive. φ>1
    means AI amplifies output per task.

  Why it matters:
    When φ>1, automation increases total output, which raises product
    demand — partially offsetting the labour-income demand destruction.
    The equilibrium condition becomes quadratic in α, with the positive
    root giving α_NE. This is the "Red Queen" scenario: AI is so
    productive that firms race to automate even harder, and the output
    effect can amplify or dampen the wedge depending on parameter values.

  Direction:
    ↑ φ  →  output per task rises; demand destruction is partially
            offset; but private incentive to automate also rises.
            Net effect: α_NE typically rises; wedge changes sign possible.
    φ = 1 →  base model (no productivity premium).
    φ < 1 →  AI is less productive than humans (rarely modelled).

  Activate with: 'extend phi <value>'
  Paper base: φ = 1.00 (extension inactive)`,
	},
	{
		symbol: "etahat",
		name:   "Capital income recycling rate (η̂) — extension",
		current: func(p model.Params) string { return fmt.Sprintf("%.4f", p.EtaHat) },
		body: `  What it is:
    The fraction of firms' automation profits (owner surplus) that gets
    recycled back into product demand — through dividends consumed by
    wealthy households, sovereign wealth fund spending, or profit-sharing
    schemes. η̂=0 (default) means profits are saved or invested outside
    the model; η̂=1 means full recycling.

  Why it matters:
    Capital recycling is the counterpart to labour income replacement (η).
    While η plugs the hole from lost wages, η̂ asks whether automation
    profits create compensating demand. In the baseline model, profits
    are not recycled, which is why automation is unambiguously demand-
    destructive. With η̂>0, the demand destruction channel weakens because
    owners spend some of their gains.

  Direction:
    ↑ η̂  →  ℓ_eff shrinks → wedge narrows; high enough η̂ can fully
            offset demand destruction (analogous to η=1 from the labour
            side).
    η̂ = 0 →  base case; no recycling.

  Activate with: 'extend recycling <value>'
  Paper base: η̂ = 0.00 (extension inactive)`,
	},
}

// derivedExplanations covers s, ℓ, N*.
const derivedExplText = `
  Derived quantities (computed from the parameters above):

  s = w − c  (cost saving per task)
    The private net benefit of automating one task. This is the margin
    that drives every firm's automation decision. If s ≤ 0, no firm
    ever automates. All equilibrium formulas are proportional to s.

  ℓ = λ(1−η)w  (demand loss per automated task)
    The aggregate demand destroyed when one task is automated. A worker
    loses wage w; they would have spent fraction λ of it; but only the
    portion (1−η) is truly lost (η is replaced). ℓ is the social cost
    that each firm ignores 1/N of.

  N* = ℓ/s  (externality activation threshold)
    The minimum number of firms needed for the Nash equilibrium to
    exceed the cooperative equilibrium. When N ≤ N*, even the Nash
    solution is efficient. When N > N*, the trap is active. With paper
    base params: N* = 0.35/0.70 = 0.50, so the trap is always active
    for any N ≥ 1.
`

// Explain is the public entry-point used by arg mode (--explain).
func (r *REPL) Explain(args []string) { r.cmdExplain(args) }

func (r *REPL) cmdExplain(args []string) {
	if len(args) == 0 {
		// Full deep-dive: banner, all params, derived, current values
		p := r.params
		fmt.Fprintln(r.out)
		fmt.Fprintln(r.out, "╔══════════════════════════════════════════════════════════════════╗")
		fmt.Fprintln(r.out, "║          AI LAYOFF TRAP — PARAMETER REFERENCE                    ║")
		fmt.Fprintln(r.out, "╚══════════════════════════════════════════════════════════════════╝")
		fmt.Fprintln(r.out, "  The model: N symmetric firms each choose automation rate α ∈ [0,1].")
		fmt.Fprintln(r.out, "  Each captures full cost savings but bears only 1/N of demand loss")
		fmt.Fprintln(r.out, "  → dominant-strategy over-automation: the AI Layoff Trap.")
		fmt.Fprintln(r.out, "  Key formulas:")
		fmt.Fprintln(r.out, "    α_NE = min(max(0, (s − ℓ/N) / k), 1)")
		fmt.Fprintln(r.out, "    α_CO = min(max(0, (s − ℓ)   / k), 1)")
		fmt.Fprintln(r.out, "    Wedge = α_NE − α_CO = ℓ(1−1/N)/k")
		fmt.Fprintln(r.out)

		for _, e := range paramExplanations {
			cur := e.current(p)
			fmt.Fprintf(r.out, "──────────────────────────────────────────────────────────────────\n")
			fmt.Fprintf(r.out, "  %s  —  %s  (current: %s)\n", e.symbol, e.name, cur)
			fmt.Fprintln(r.out)
			fmt.Fprintln(r.out, e.body)
			fmt.Fprintln(r.out)
		}

		fmt.Fprintf(r.out, "──────────────────────────────────────────────────────────────────\n")
		fmt.Fprintln(r.out, derivedExplText)
		fmt.Fprintf(r.out, "  Tip: '?? <param>' for a single param.  e.g.:  ?? eta\n")
		fmt.Fprintln(r.out)
		return
	}

	// Single-param lookup
	key := strings.ToLower(args[0])
	for _, e := range paramExplanations {
		if strings.ToLower(e.symbol) == key {
			cur := e.current(r.params)
			fmt.Fprintln(r.out)
			fmt.Fprintf(r.out, "  %s  —  %s  (current: %s)\n", e.symbol, e.name, cur)
			fmt.Fprintln(r.out)
			fmt.Fprintln(r.out, e.body)
			fmt.Fprintln(r.out)
			return
		}
	}
	// derived
	if key == "s" || key == "ell" || key == "l" || key == "nstar" || key == "n*" {
		fmt.Fprintln(r.out, derivedExplText)
		return
	}
	fmt.Fprintf(r.out, "unknown param %q; valid: N w c k lambda eta A L mu phi etahat  (or no arg for all)\n", args[0])
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
