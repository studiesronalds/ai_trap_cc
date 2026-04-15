package output

import (
	"encoding/json"
	"io"

	"aisim/internal/model"
)

type jsonResult struct {
	Params struct {
		N      float64 `json:"N"`
		W      float64 `json:"w"`
		C      float64 `json:"c"`
		K      float64 `json:"k"`
		Lambda float64 `json:"lambda"`
		Eta    float64 `json:"eta"`
		A      float64 `json:"A"`
		L      float64 `json:"L"`
		Mu     float64 `json:"mu"`
	} `json:"params"`
	Derived struct {
		S     float64 `json:"s"`
		Ell   float64 `json:"ell"`
		NStar float64 `json:"N_star"`
	} `json:"derived"`
	Equilibrium struct {
		AlphaNE  float64 `json:"alpha_NE"`
		AlphaCO  float64 `json:"alpha_CO"`
		AlphaSP  float64 `json:"alpha_SP"`
		Wedge    float64 `json:"wedge"`
		WedgePct float64 `json:"wedge_pct"`
	} `json:"equilibrium"`
	Welfare struct {
		OwnerSurplus   float64 `json:"K_NE"`
		WorkerIncome   float64 `json:"W_NE"`
		OwnerSurplusCO float64 `json:"K_CO"`
		WorkerIncomeCO float64 `json:"W_CO"`
		SurplusLoss    float64 `json:"surplus_loss"`
	} `json:"welfare"`
	Policy struct {
		TaxStar    float64 `json:"tau_star"`
		AlphaTaxed float64 `json:"alpha_taxed"`
	} `json:"policy"`
}

// WriteJSON serialises a Result to JSON.
func WriteJSON(w io.Writer, r model.Result) error {
	p := r.Params
	out := jsonResult{}
	out.Params.N = p.N; out.Params.W = p.W; out.Params.C = p.C
	out.Params.K = p.K; out.Params.Lambda = p.Lambda; out.Params.Eta = p.Eta
	out.Params.A = p.A; out.Params.L = p.L; out.Params.Mu = p.Mu
	out.Derived.S = r.S; out.Derived.Ell = r.Ell; out.Derived.NStar = r.NStar
	out.Equilibrium.AlphaNE = r.AlphaNE; out.Equilibrium.AlphaCO = r.AlphaCO
	out.Equilibrium.AlphaSP = r.AlphaSP
	out.Equilibrium.Wedge = r.Wedge; out.Equilibrium.WedgePct = r.WedgePct
	out.Welfare.OwnerSurplus = r.OwnerSurplus; out.Welfare.WorkerIncome = r.WorkerIncome
	out.Welfare.OwnerSurplusCO = r.OwnerSurplusCO; out.Welfare.WorkerIncomeCO = r.WorkerIncomeCO
	out.Welfare.SurplusLoss = r.SurplusLoss
	out.Policy.TaxStar = r.TaxStar; out.Policy.AlphaTaxed = r.AlphaTaxed

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// WriteSweepJSON serialises a 1D sweep to JSON.
func WriteSweepJSON(w io.Writer, pts []model.SweepPoint) error {
	type row struct {
		Value   float64 `json:"value"`
		AlphaNE float64 `json:"alpha_NE"`
		AlphaCO float64 `json:"alpha_CO"`
		Wedge   float64 `json:"wedge"`
	}
	rows := make([]row, len(pts))
	for i, pt := range pts {
		rows[i] = row{pt.Value, pt.Result.AlphaNE, pt.Result.AlphaCO, pt.Result.Wedge}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
