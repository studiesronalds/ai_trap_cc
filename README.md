# aisim — AI Layoff Trap Simulator

A terminal simulator of the formal economic model of AI-driven over-automation described in *"The AI Layoff Trap"*. N symmetric firms each choose how much to automate. Every firm captures 100 % of its cost savings but bears only 1/N of the resulting demand destruction — a dominant-strategy externality that drives the market to over-automate relative to both the cooperative and socially optimal outcomes.

---

## Quick start

```
# interactive REPL (default when run with no args)
./aisim

# one-shot report with paper base parameters
./aisim --N=7

# explain what every parameter means
./aisim --explain=all

# sweep firms 2–20 and watch the wedge grow
./aisim --sweep=N --sweep-from=2 --sweep-to=20 --sweep-steps=18

# compare all 6 policy instruments
./aisim --policy=all
```

---

## Installation

**Pre-built binaries** are in `dist/` — pick the one for your platform:

| File | Platform |
|---|---|
| `aisim-linux-amd64` | Linux x86-64 |
| `aisim-linux-arm64` | Linux ARM64 |
| `aisim-darwin-amd64` | macOS Intel |
| `aisim-darwin-arm64` | macOS Apple Silicon |
| `aisim-windows-amd64.exe` | Windows x86-64 |
| `aisim-windows-arm64.exe` | Windows ARM64 |

**Build from source** (requires Go 1.25+):

```bash
# all platforms at once
bash build.sh          # Linux / macOS
build.cmd              # Windows

# current platform only
go build -o aisim ./cmd/aisim
```

---

## The model in one paragraph

Each firm's per-task profit from automating at rate α is:

```
π(α) = Π₀ + L·[s·α  −  (ℓ/N)·α  −  (k/2)·α²]
```

where `s = w − c` is the cost saving, `ℓ = λ(1−η)w` is the demand loss per task, and `k` is the integration friction. Firms internalise only `ℓ/N` of the demand destruction, so the Nash equilibrium sits above the cooperative optimum by the **wedge**:

```
Wedge  =  α_NE − α_CO  =  ℓ(1−1/N) / k
```

The externality is active when the number of firms exceeds the threshold `N* = ℓ/s`. The optimal Pigouvian correction tax is `τ* = ℓ(1−1/N)`.

---

## Usage modes

### Arg mode

Single-shot computations from the command line — pipe into scripts, compare with `diff`, save to JSON/CSV.

```
./aisim [--param=val ...] [--preset=name] [--command] [--output=format]
```

### REPL mode

Interactive session with persistent state, named scenarios, and a full command set.

```
./aisim          # no args → REPL
./aisim --repl   # explicit
```

See [MANUAL.md](MANUAL.md) for the complete command reference.

---

## Parameters at a glance

| Flag / REPL name | Symbol | Default | Meaning |
|---|---|---|---|
| `--N` | N | 7 | Number of symmetric firms |
| `--w` | w | 1.00 | Wage per task |
| `--c` | c | 0.30 | AI cost per task |
| `--k` | k | 1.00 | Integration friction |
| `--lambda` | λ | 0.50 | Worker marginal propensity to consume |
| `--eta` | η | 0.30 | Income replacement rate |
| `--A` | A | 10.0 | Autonomous demand |
| `--L` | L | 100 | Tasks per firm |
| `--mu` | µ | 0.50 | Social planner worker weight |
| `--phi` | φ | 1.00 | AI productivity multiplier (extension) |

Run `./aisim --explain=all` (or type `??` in the REPL) for deep economic explanations of every parameter.

---

## Presets

| Name | Description |
|---|---|
| `paper` | Paper base: N=7, w=1, c=0.30, k=1, λ=0.5, η=0.3 |
| `frictionless` | k=0 — Prisoner's Dilemma mode |
| `monopolist` | N=1 — fully internalises the externality, wedge=0 |
| `duopoly` | N=2 |
| `competitive` | N=50 — near-competitive limit |
| `optimist` | η=1.1 — generous income replacement |

```bash
./aisim --preset=monopolist
./aisim --preset=frictionless --output=json
```

---

## Output formats

| Flag | Output |
|---|---|
| `--output=text` | Human-readable report (default) |
| `--output=json` | Machine-readable JSON |
| `--output=csv` | CSV rows |

---

## Project layout

```
aisim/
├── cmd/aisim/main.go          entry point, flag parsing, mode dispatch
├── internal/
│   ├── model/
│   │   ├── params.go          Params struct, presets, validation, derived quantities
│   │   ├── equilibrium.go     Compute() — Nash, cooperative, social planner, welfare
│   │   ├── welfare.go         welfare decomposition, sensitivity analysis
│   │   ├── policy.go          all 6 policy instruments, coalition, Pigouvian tax
│   │   ├── extensions.go      φ (Red Queen), endogenous wages, capital recycling, entry
│   │   └── sweep.go           1D and 2D parameter sweep engines
│   ├── output/
│   │   ├── table.go           text report, policy table, sweep table, sensitivity table
│   │   ├── heatmap.go         2D ASCII density heatmap
│   │   ├── json.go            JSON export
│   │   └── csv.go             CSV export
│   └── repl/
│       ├── repl.go            REPL loop, I/O wiring
│       ├── commands.go        all command implementations + ?? explain system
│       └── scenario.go        named scenario save/load/diff
├── build.sh                   cross-platform build script (Linux/macOS)
├── build.cmd                  cross-platform build script (Windows)
├── dist/                      pre-built binaries (git-ignored)
├── README.md                  this file
└── MANUAL.md                  full command and parameter reference
```

---

## Further reading

- `MANUAL.md` — complete reference for every command, flag, and formula
- Type `??` in the REPL for in-terminal parameter explanations
- Type `help` in the REPL for the full command list

---

## Based on original work

This simulator implements the model from:

> **"The AI Layoff Trap"**
> Brett Hemenway Falk & Gerry Tsoukalas
> arXiv:2603.20617v1 [econ.TH], 21 March 2026
> Hemenway Falk: University of Pennsylvania (fbrett@cis.upenn.edu)
> Tsoukalas: Boston University (gerryt@bu.edu)

All economic theory, mathematical derivations, propositions, and policy analysis reproduced or referenced in this tool originate from that paper. This simulator is an independent implementation for research and educational purposes; it is not affiliated with or endorsed by the authors.

If you use this tool in work that draws on the underlying model, please cite the original paper.
