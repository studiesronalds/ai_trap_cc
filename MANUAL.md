# aisim — User Manual

Full reference for all command-line flags, REPL commands, parameters, output formats, and model mathematics.

---

## Table of contents

1. [Invocation](#1-invocation)
2. [Arg mode reference](#2-arg-mode-reference)
3. [REPL mode reference](#3-repl-mode-reference)
4. [Parameter reference](#4-parameter-reference)
5. [Derived quantities](#5-derived-quantities)
6. [Equilibria and welfare](#6-equilibria-and-welfare)
7. [Policy instruments](#7-policy-instruments)
8. [Extensions](#8-extensions)
9. [Output formats](#9-output-formats)
10. [Presets](#10-presets)
11. [Model mathematics](#11-model-mathematics)

---

## 1. Invocation

```
aisim [flags]
```

| Condition | Mode entered |
|---|---|
| No arguments | REPL |
| `--repl` | REPL |
| Any parameter flag alone (`--N=10`) | One-shot report |
| `--sweep=<param>` | Sweep table |
| `--policy=<instrument>` | Policy report |
| `--explain=<param\|all>` | Parameter explanation |

---

## 2. Arg mode reference

### Parameter flags

These set the simulation parameters. Any unset flag takes its default value (or the chosen preset's value).

| Flag | Default | Description |
|---|---|---|
| `--N=<val>` | 7 | Number of symmetric firms |
| `--w=<val>` | 1.00 | Wage per task |
| `--c=<val>` | 0.30 | AI cost per task (must be ≤ w) |
| `--k=<val>` | 1.00 | Integration friction (≥ 0) |
| `--lambda=<val>` | 0.50 | Worker marginal propensity to consume [0,1] |
| `--eta=<val>` | 0.30 | Income replacement rate (≥ 0) |
| `--A=<val>` | 10.0 | Autonomous demand (≥ 0) |
| `--L=<val>` | 100 | Tasks per firm (> 0) |
| `--mu=<val>` | 0.50 | Social planner worker weight [0,1] |
| `--phi=<val>` | 1.00 | AI productivity multiplier (> 0) |

### Mode flags

| Flag | Description |
|---|---|
| `--repl` | Start interactive REPL |
| `--preset=<name>` | Load a named preset (see §10) |

### Command flags

| Flag | Description |
|---|---|
| `--policy=<instrument>` | Run a policy instrument: `tax`, `ubi`, `upskill`, `coalition`, `capital-tax`, `all` |
| `--sweep=<param>` | 1D parameter sweep |
| `--sweep-from=<val>` | Sweep lower bound (default 0) |
| `--sweep-to=<val>` | Sweep upper bound (default 10) |
| `--sweep-steps=<n>` | Number of steps (default 20) |
| `--sweep2=<param>` | Second param for 2D heatmap sweep |
| `--sweep2-from=<val>` | Second param lower bound (default 0) |
| `--sweep2-to=<val>` | Second param upper bound (default 1) |
| `--explain=<param\|all>` | Deep-explain a parameter or all parameters |

### Output flag

| Flag | Values | Default |
|---|---|---|
| `--output=<fmt>` | `text`, `json`, `csv` | `text` |

### Examples

```bash
# Default report (paper base)
./aisim --N=7

# Override two params
./aisim --N=10 --eta=0.6

# Load preset, override one param
./aisim --preset=competitive --k=0.5

# All policy instruments, JSON output
./aisim --policy=all --output=json

# 1D sweep of eta from 0 to 1
./aisim --sweep=eta --sweep-from=0 --sweep-to=1 --sweep-steps=20

# 2D heatmap: N (rows) vs eta (cols)
./aisim --sweep=N --sweep-from=2 --sweep-to=20 \
        --sweep2=eta --sweep2-from=0 --sweep2-to=1

# Explain income replacement rate
./aisim --explain=eta

# Full parameter reference
./aisim --explain=all
```

---

## 3. REPL mode reference

Start with `./aisim` or `./aisim --repl`. The prompt is `aisim>`.

All parameters persist across commands within a session. Use `set` to change them, `show` to inspect them, `save`/`load` to bookmark scenarios.

### Core commands

#### `run`
Compute equilibrium and print the full report for the current parameters.

```
aisim> run
```

#### `show`
Print all current parameter values plus derived quantities s and ℓ.

```
aisim> show
```

#### `set <param> <value>`
Set one parameter. Validates immediately; rejects out-of-range values.

```
aisim> set N 15
aisim> set eta 0.6
aisim> set k 0
```

Valid param names: `N  w  c  k  lambda  eta  A  L  mu  phi  etahat`

#### `??  [param]`
Deep-explain all parameters (no arg) or one specific parameter.
Shows: symbol, current value, economic meaning, model role, directional effects, paper base value.

```
aisim> ??
aisim> ?? eta
aisim> ?? lambda
aisim> ?? phi
```

Alias: `explain [param]`

---

### Analysis commands

#### `welfare [mu=<val>]`
Print welfare metrics at the given planner weight µ (or current µ if omitted).
Shows α_SP, K, W, K/K_CO, W/W_CO, and surplus loss.

```
aisim> welfare
aisim> welfare mu=0.8
aisim> welfare 0.0
```

#### `decompose`
Split the over-automation wedge into two components:
- **Demand externality**: α_NE − α_CO (pure market-failure component, µ-independent)
- **Distributional premium**: α_CO − α_SP (extra restraint wanted because of worker weighting)

```
aisim> decompose
```

#### `sensitivity`
Print ∂(wedge)/∂(param) for all core parameters using central differences.
Identifies which parameters the wedge is most sensitive to under current values.

```
aisim> sensitivity
```

---

### Policy commands

#### `policy tax`
Compute optimal Pigouvian tax τ* = ℓ(1−1/N) and the corrected equilibrium α under that tax.

```
aisim> policy tax
```

#### `policy ubi [dA]`
Show the effect of a UBI that adds `dA` to autonomous demand A (default dA=1).
α_NE is unchanged — the key result that UBI cannot close the wedge.

```
aisim> policy ubi
aisim> policy ubi 5
```

#### `policy upskill [η]`
Show the effect of raising η to the given value (default: current η + 0.20).
Demonstrates the upskilling/retraining mechanism.

```
aisim> policy upskill
aisim> policy upskill 0.8
aisim> policy upskill 1.0
```

#### `policy equity [ε]`
Show the effect of worker equity participation fraction ε (default 0.1).
Workers holding firm equity raises their effective η, narrowing (but not closing) the wedge.

```
aisim> policy equity
aisim> policy equity 0.3
```

#### `policy coalition [M]`
Show the partial coalition equilibrium: M firms coordinate, the remaining N−M play Nash.
Default M = N/2.

```
aisim> policy coalition
aisim> policy coalition 3
```

#### `policy capital-tax`
Show the no-effect result: a tax on profits does not change the per-task automation margin and therefore cannot close the wedge.

```
aisim> policy capital-tax
```

#### `policy all`
Side-by-side table of all 6 instruments: α after, residual wedge, and notes.

```
aisim> policy all
```

---

### Sweep commands

#### `sweep <param> <from> <to> [steps]`
1D parameter sweep. Prints a table of α_NE, α_CO, and wedge at each value, with a sparkline column.

```
aisim> sweep N 2 20
aisim> sweep N 2 20 30
aisim> sweep eta 0 1 20
aisim> sweep k 0.1 3 15
```

Valid sweep params: `N  w  c  k  lambda  eta  A  L  mu  phi`

#### `heatmap <param1> <param2>`
2D ASCII density heatmap of the wedge. Rows = param1, cols = param2.
Uses sensible default ranges per parameter (override with `sweep` for custom ranges).

```
aisim> heatmap N eta
aisim> heatmap N k
aisim> heatmap lambda eta
```

---

### Extensions

#### `extend phi [φ]`
AI productivity extension. φ > 1 means each automated task produces φ times the output.
Introduces a Red Queen dynamic: automation raises output and thus demand, partially offsetting labour income destruction. Equilibrium condition becomes quadratic; the positive root is used.

```
aisim> extend phi 1.2
aisim> extend phi 1.5
```

#### `extend wages [slope]`
Endogenous wages: w(ᾱ) = w₀ + slope·ᾱ where ᾱ is the industry average.
Solves for fixed-point equilibrium via monotone iteration (default slope = −0.5, wages fall as automation rises).

```
aisim> extend wages
aisim> extend wages -0.3
aisim> extend wages 0.1
```

#### `extend recycling [η̂]`
Capital income recycling. Fraction η̂ of automation profits is recycled back into product demand (through dividends, sovereign funds, etc.). Reduces effective demand destruction.

```
aisim> extend recycling 0.5
aisim> extend recycling 1.0
```

#### `extend entry [κ]`
Endogenous firm entry. Firms enter until expected profit equals entry cost κ. Returns the equilibrium N and regime (trapped or benign).

```
aisim> extend entry 5
aisim> extend entry 2
```

---

### Special modes

#### `frictionless`
Sets k=0 and runs the report. Equivalent to `set k 0` then `run`. Activates Prisoner's Dilemma mode where automation is a binary choice.

```
aisim> frictionless
```

#### `prisoners`
Explicit 2×2 Prisoner's Dilemma payoff matrix for N=2, k=0. Shows the four payoff cells (both automate / both don't / one defects) and identifies the Nash and cooperative outcomes.

```
aisim> prisoners
```

---

### Scenario management

#### `save <name>`
Save a snapshot of current parameters under a name.

```
aisim> save baseline
aisim> save high-eta
```

#### `load <name>`
Restore a previously saved scenario.

```
aisim> load baseline
```

#### `diff <name1> <name2>`
Side-by-side comparison of two saved scenarios: parameter differences and resulting equilibrium differences.

```
aisim> diff baseline high-eta
```

#### `scenarios`
List all saved scenario names in the current session.

```
aisim> scenarios
```

---

### Preset command

#### `preset <name>`
Load a built-in parameter preset (see §10).

```
aisim> preset monopolist
aisim> preset frictionless
aisim> preset competitive
```

---

### Export command

#### `export json|csv`
Print the current result in JSON or CSV format to stdout.

```
aisim> export json
aisim> export csv
```

---

### Session commands

| Command | Description |
|---|---|
| `help` | Print all commands |
| `quit` / `exit` / `q` | End the session |

---

## 4. Parameter reference

### N — Number of firms

The count of symmetric competing firms. Every firm is identical in cost structure, task count, and wage level.

**Role in model:** N is the engine of the externality. Each firm captures 100 % of its cost savings but shares the demand destruction 1/N ways. The larger N, the smaller each firm's internal cost of automating and the deeper the trap.

**Key relationships:**
- Wedge = ℓ(1−1/N)/k → grows toward ℓ/k as N→∞
- N=1 (monopolist): wedge = 0, Nash = cooperative
- Externality active when N > N* = ℓ/s

**Paper base:** 7

---

### w — Wage per task

Labour cost when a human worker performs the task. The monetary unit of the model.

**Role in model:** Sets the gross cost-saving potential (s = w−c) and simultaneously the labour income available to be destroyed (W = wLN·[1−(1−η)ᾱ]).

**Paper base:** 1.00

---

### c — AI cost per task

Per-task cost of AI automation: compute, licensing, integration, maintenance. Must satisfy c ≤ w.

**Role in model:** Together with w it determines s = w−c, the private benefit driving each firm's automation decision. Falling c (the AI price-crash scenario) is the primary mechanism that worsens the trap.

**Paper base:** 0.30 (so s = 0.70)

---

### k — Integration friction

Convexity parameter on the cost of automation. Captures reorganising workflows, retraining management, compliance overhead. Profit from automating fraction α has a −(k/2)α² term.

**Role in model:** k is the denominator of every equilibrium formula. Without it (k=0), automation is a binary Prisoner's Dilemma. Higher k slows the race to automate for all firms equally.

**Special case:** k=0 → frictionless / Prisoner's Dilemma mode.

**Paper base:** 1.00

---

### lambda (λ) — Worker marginal propensity to consume

Fraction of labour income that workers spend back into product markets. Range [0,1].

**Role in model:** Scales the demand destruction channel. When a task is automated the worker loses wage w; they would have spent λ of it. So ℓ = λ(1−η)w is the demand lost per automated task.

**Paper base:** 0.50

---

### eta (η) — Income replacement rate

Fraction of lost wage income that displaced workers recover through transfers, new employment, benefits, or retraining. η=0: no safety net. η=1: full replacement. η>1: net income gain possible.

**Role in model:** The second factor in ℓ = λ(1−η)w. Raising η directly reduces the demand destruction per task. This is the mechanism behind upskilling programmes — they work not by making AI costlier but by rebuilding worker spending power.

**Policy relevance:** The only demand-channel lever besides the Pigouvian tax. At η=1, ℓ=0 and the externality vanishes.

**Paper base:** 0.30

---

### A — Autonomous demand

Product demand that is independent of worker income: government spending, exports, investment, wealthy household consumption.

**Role in model:** Sets the demand floor. Total demand is D = A + λwLN[1−(1−η)ᾱ]. A does **not** affect α_NE or the wedge — automation incentives are determined by the per-task margin, not the demand level. UBI raises A and improves worker welfare but cannot close the wedge for this reason.

**Paper base:** 10.0

---

### L — Tasks per firm

Number of separable automatable tasks each firm controls. α·L tasks get automated; (1−α)·L remain human.

**Role in model:** A pure scale parameter. L cancels out of equilibrium conditions for α; it scales absolute welfare quantities but leaves all automation rates and the wedge unchanged.

**Paper base:** 100

---

### mu (µ) — Social planner worker weight

Weight the social planner places on worker welfare relative to owner welfare. µ=0: planner maximises profits only. µ=1: planner maximises worker income only. µ=0.5: utilitarian split.

**Role in model:** Determines α_SP. The wedge between α_NE and α_SP decomposes into the market-failure component (µ-independent) and a distributional premium (growing with µ).

**Paper base:** 0.50

---

### phi (φ) — AI productivity multiplier *(extension)*

Multiplier on output per automated task. φ=1: AI and human equally productive. φ>1: AI amplifies output (Red Queen scenario).

**Role in model:** When φ>1, automation raises total output and thus demand, partially offsetting labour income destruction. The equilibrium condition becomes quadratic; the positive root is used. See `extend phi`.

**Paper base:** 1.00 (extension inactive)

---

### etahat (η̂) — Capital income recycling *(extension)*

Fraction of automation profits recycled back into product demand through wealthy-household spending, dividends, or sovereign funds. η̂=0: profits are fully saved or invested outside the model.

**Role in model:** The capital-side counterpart to η. High η̂ can substantially offset demand destruction. See `extend recycling`.

**Paper base:** 0.00 (extension inactive)

---

## 5. Derived quantities

These are computed automatically from the parameters.

| Symbol | Formula | Meaning |
|---|---|---|
| s | w − c | Cost saving per task (private benefit of automating one task) |
| ℓ | λ(1−η)w | Demand loss per automated task (social cost firms do not fully internalise) |
| N* | ℓ/s | Externality activation threshold: trap active when N > N* |

---

## 6. Equilibria and welfare

### Equilibrium automation rates

| Name | Formula | Meaning |
|---|---|---|
| α_NE | clamp((s − ℓ/N)/k, 0, 1) | Nash equilibrium: each firm's dominant strategy |
| α_CO | clamp((s − ℓ)/k, 0, 1) | Cooperative optimum: maximises joint firm profit |
| α_SP | clamp(((1−µ)(s−ℓ) − µw(1−η))/k, 0, 1) | Social planner optimum at weight µ |
| Wedge | α_NE − α_CO = ℓ(1−1/N)/k | Over-automation gap |

**Frictionless case (k=0):** equilibria are binary (0 or 1); the game is a pure Prisoner's Dilemma.

### Welfare metrics

| Symbol | Formula | Meaning |
|---|---|---|
| K | N·[Π₀ + L(sα − ℓα − k/2·α²)] | Aggregate owner surplus at NE |
| W | wLN[1−(1−η)α] | Aggregate worker income at NE |
| S(µ) | (1−µ)K + µW | Social welfare at planner weight µ |
| K_CO | K evaluated at α_CO | Owner surplus at cooperative outcome |
| Surplus loss | K_CO − K_NE | Deadweight loss from over-automation |

Welfare ratios K/K_CO and W/W_CO are reported relative to the cooperative baseline (= 1.00).

---

## 7. Policy instruments

Six instruments from the paper, evaluated with `policy all` or individually:

| Instrument | Command | Closes wedge? | Mechanism |
|---|---|---|---|
| Pigouvian tax | `policy tax` | Yes | τ* = ℓ(1−1/N) shifts per-task margin to internalise externality |
| UBI | `policy ubi [dA]` | No | Raises A, improves W but leaves α_NE unchanged |
| Capital income tax | `policy capital-tax` | No | Taxes profit levels, not per-task margin |
| Worker equity | `policy equity [ε]` | No | Raises effective η, narrows wedge but cannot eliminate |
| Upskilling | `policy upskill [η]` | At η=1 only | Reduces ℓ directly; full closure requires η=1 |
| Coasian coalition | `policy coalition [M]` | Yes (but unstable) | All-N coalition → α_CO; not self-enforcing |

**The key result:** the Pigouvian tax τ* = ℓ(1−1/N) is the only instrument that directly corrects the per-task externality and implements the cooperative optimum as a Nash equilibrium.

---

## 8. Extensions

These go beyond the baseline model and require the `extend` command (REPL) or can be accessed via `--policy` in some cases.

### AI productivity (φ > 1) — Red Queen

When φ > 1, each automated task produces φ units of output instead of 1. The demand equation gains a positive term from increased output. Equilibrium condition is quadratic in α; the code uses the positive root. Use `extend phi <val>`.

### Endogenous wages

w(ᾱ) = w₀ + slope·ᾱ where ᾱ is the industry average automation rate. The fixed-point equilibrium is solved by monotone iteration (typically converges in < 50 iterations). A negative slope models wage depression as automation rises. Use `extend wages <slope>`.

### Capital income recycling (η̂)

Fraction η̂ of owner surplus is recycled into product demand each period. Reduces the effective demand destruction coefficient. Use `extend recycling <etahat>`.

### Endogenous firm entry

Firms enter until expected per-firm profit equals entry cost κ. Returns the equilibrium number of firms N_eq and classifies the regime as "trapped" (wedge > 0) or "benign". Use `extend entry <kappa>`.

---

## 9. Output formats

### Text (default)

Human-readable formatted report. Sections: parameters, derived quantities, equilibria, wedge decomposition, welfare, Pigouvian tax.

```bash
./aisim --N=7 --output=text
```

### JSON

Structured output suitable for piping into `jq`, Python, or any downstream tool. Contains `params`, `derived`, `equilibrium`, `welfare`, and `policy` objects.

```bash
./aisim --N=7 --output=json | jq '.equilibrium.wedge'
```

### CSV

Single-row CSV with headers. Useful for building datasets with shell loops:

```bash
for n in 2 5 10 20; do
  ./aisim --N=$n --output=csv
done
```

For sweep outputs, each row is one sweep point.

---

## 10. Presets

| Name | N | w | c | k | λ | η | Notes |
|---|---|---|---|---|---|---|---|
| `paper` | 7 | 1.00 | 0.30 | 1.00 | 0.50 | 0.30 | Paper base parameters |
| `frictionless` | 7 | 1.00 | 0.30 | **0** | 0.50 | 0.30 | k=0, Prisoner's Dilemma |
| `monopolist` | **1** | 1.00 | 0.30 | 1.00 | 0.50 | 0.30 | N=1, wedge=0 |
| `duopoly` | **2** | 1.00 | 0.30 | 1.00 | 0.50 | 0.30 | N=2 |
| `competitive` | **50** | 1.00 | 0.30 | 1.00 | 0.50 | 0.30 | Near-competitive limit |
| `optimist` | 7 | 1.00 | 0.30 | 1.00 | 0.50 | **1.10** | Full-plus income replacement |

Presets can be used as a starting point and then modified with `set` or by passing additional flags:

```bash
# Competitive market with lower friction
./aisim --preset=competitive --k=0.5

# REPL starting from monopolist
./aisim --preset=monopolist --repl
```

---

## 11. Model mathematics

### Setup

N symmetric firms, each choosing automation rate α ∈ [0,1].

Per-firm per-task profit from automation at rate α, given industry average ᾱ:

```
π(α, ᾱ) = Π₀  +  L · [s·α  −  (ℓ/N)·ᾱ·N·(1/N)  −  (k/2)·α²]
         = Π₀  +  L · [s·α  −  ℓ·ᾱ/N              −  (k/2)·α²]
```

In the symmetric Nash equilibrium all firms set α = ᾱ, so:

```
π(α) = Π₀  +  L · [(s − ℓ/N)·α  −  (k/2)·α²]
```

### Nash equilibrium (Proposition 1)

First-order condition ∂π/∂α = 0:

```
α_NE  =  clamp( (s − ℓ/N) / k ,  0, 1 )
```

### Cooperative optimum (Proposition 2)

Firms jointly maximise aggregate profit. Each task's social cost is ℓ (not ℓ/N):

```
α_CO  =  clamp( (s − ℓ) / k ,  0, 1 )
```

### Over-automation wedge (Proposition 3)

```
Wedge  =  α_NE − α_CO  =  ℓ(1 − 1/N) / k
```

The wedge:
- is zero at N=1 (monopolist)
- grows monotonically in N, approaching ℓ/k as N→∞
- is proportional to the demand loss ℓ and inversely proportional to friction k
- is zero when ℓ=0 (η=1 or λ=0 — demand destruction neutralised)

### Social planner optimum

The social planner maximises S(µ) = (1−µ)K + µW. Setting ∂S/∂α = 0:

```
α_SP  =  clamp( [(1−µ)(s−ℓ) − µw(1−η)] / k ,  0, 1 )
```

### Pigouvian tax (Proposition 5)

The optimal per-task tax that corrects the externality:

```
τ*  =  ℓ(1 − 1/N)
```

Under this tax, firm's effective saving becomes s−τ*, and Nash equilibrium reduces to:

```
α(τ*)  =  (s − τ* − ℓ/N) / k  =  (s − ℓ) / k  =  α_CO
```

The tax exactly internalises the externality and implements the cooperative outcome.

### Externality activation threshold

The externality is dormant (α_NE = α_CO) when the wedge is zero, which occurs when N ≤ N*:

```
N*  =  ℓ / s  =  λ(1−η)w / (w−c)
```

For paper base: N* = 0.35/0.70 = 0.50, so the trap is active for all N ≥ 1.

### Welfare formulas

```
K(α)   =  N · [Π₀ + L·(s − ℓ)·α − N·L·(k/2)·α²]     (owner surplus)
W(α)   =  w·L·N · [1 − (1−η)·α]                        (worker income)
S(µ,α) =  (1−µ)·K(α) + µ·W(α)                          (social welfare)
```

where Π₀ = A/N + (λ−1)·w·L is the baseline per-firm profit (before any automation).

### Wedge decomposition

```
α_NE − α_SP  =  [α_NE − α_CO]  +  [α_CO − α_SP]
              =  demand externality  +  distributional premium
```

The demand externality is the pure market-failure component (independent of µ).
The distributional premium is the extra restraint a worker-weighting planner wants beyond correcting the externality.
