package repl

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"aisim/internal/model"
)

const scenarioFile = "scenarios.json"

type scenarioStore map[string]model.Params

func loadScenarios() scenarioStore {
	data, err := os.ReadFile(scenarioFile)
	if err != nil {
		return make(scenarioStore)
	}
	var s scenarioStore
	if err := json.Unmarshal(data, &s); err != nil {
		return make(scenarioStore)
	}
	return s
}

func saveScenarios(s scenarioStore) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(scenarioFile, data, 0644)
}

func (r *REPL) saveScenario(name string) {
	s := loadScenarios()
	s[name] = r.params
	if err := saveScenarios(s); err != nil {
		fmt.Fprintf(r.out, "error saving: %v\n", err)
		return
	}
	fmt.Fprintf(r.out, "Saved scenario %q\n", name)
}

func (r *REPL) loadScenario(name string) {
	s := loadScenarios()
	p, ok := s[name]
	if !ok {
		fmt.Fprintf(r.out, "unknown scenario %q; use 'scenarios' to list\n", name)
		return
	}
	r.params = p
	fmt.Fprintf(r.out, "Loaded scenario %q\n", name)
}

func (r *REPL) listScenarios() {
	s := loadScenarios()
	if len(s) == 0 {
		fmt.Fprintln(r.out, "No saved scenarios.")
		return
	}
	names := make([]string, 0, len(s))
	for k := range s {
		names = append(names, k)
	}
	sort.Strings(names)
	fmt.Fprintln(r.out, "Saved scenarios:")
	for _, n := range names {
		p := s[n]
		fmt.Fprintf(r.out, "  %-20s  N=%.4g w=%.2f c=%.2f k=%.2f η=%.2f\n",
			n, p.N, p.W, p.C, p.K, p.Eta)
	}
}

func (r *REPL) diffScenarios(name1, name2 string) {
	s := loadScenarios()
	p1, ok1 := s[name1]
	p2, ok2 := s[name2]
	if !ok1 {
		fmt.Fprintf(r.out, "unknown scenario %q\n", name1)
		return
	}
	if !ok2 {
		fmt.Fprintf(r.out, "unknown scenario %q\n", name2)
		return
	}
	res1 := model.Compute(p1)
	res2 := model.Compute(p2)

	fmt.Fprintf(r.out, "\n%-22s  %-12s  %-12s  %s\n", "Metric", name1, name2, "Δ")
	fmt.Fprintf(r.out, "%s\n", "──────────────────────────────────────────────────────")
	row := func(label string, v1, v2 float64) {
		fmt.Fprintf(r.out, "%-22s  %-12.4f  %-12.4f  %+.4f\n", label, v1, v2, v2-v1)
	}
	row("N", p1.N, p2.N)
	row("w", p1.W, p2.W)
	row("c", p1.C, p2.C)
	row("k", p1.K, p2.K)
	row("λ", p1.Lambda, p2.Lambda)
	row("η", p1.Eta, p2.Eta)
	row("α_NE", res1.AlphaNE, res2.AlphaNE)
	row("α_CO", res1.AlphaCO, res2.AlphaCO)
	row("Wedge", res1.Wedge, res2.Wedge)
	row("τ*", res1.TaxStar, res2.TaxStar)
	row("K (owner surplus)", res1.OwnerSurplus, res2.OwnerSurplus)
	row("W (worker income)", res1.WorkerIncome, res2.WorkerIncome)
	fmt.Fprintln(r.out)
}
