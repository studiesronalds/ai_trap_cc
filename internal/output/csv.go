package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"aisim/internal/model"
)

// WriteCSV writes a 1D sweep as CSV.
func WriteCSV(w io.Writer, pts []model.SweepPoint) error {
	cw := csv.NewWriter(w)
	if len(pts) == 0 {
		return nil
	}
	header := []string{pts[0].ParamName, "alpha_NE", "alpha_CO", "wedge", "tau_star"}
	if err := cw.Write(header); err != nil {
		return err
	}
	for _, pt := range pts {
		r := pt.Result
		row := []string{
			fmt.Sprintf("%.6f", pt.Value),
			fmt.Sprintf("%.6f", r.AlphaNE),
			fmt.Sprintf("%.6f", r.AlphaCO),
			fmt.Sprintf("%.6f", r.Wedge),
			fmt.Sprintf("%.6f", r.TaxStar),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// WriteResultCSV writes a single Result as a two-column key/value CSV.
func WriteResultCSV(w io.Writer, r model.Result) error {
	cw := csv.NewWriter(w)
	rows := [][]string{
		{"N", fmt.Sprintf("%.6f", r.Params.N)},
		{"w", fmt.Sprintf("%.6f", r.Params.W)},
		{"c", fmt.Sprintf("%.6f", r.Params.C)},
		{"k", fmt.Sprintf("%.6f", r.Params.K)},
		{"lambda", fmt.Sprintf("%.6f", r.Params.Lambda)},
		{"eta", fmt.Sprintf("%.6f", r.Params.Eta)},
		{"s", fmt.Sprintf("%.6f", r.S)},
		{"ell", fmt.Sprintf("%.6f", r.Ell)},
		{"N_star", fmt.Sprintf("%.6f", r.NStar)},
		{"alpha_NE", fmt.Sprintf("%.6f", r.AlphaNE)},
		{"alpha_CO", fmt.Sprintf("%.6f", r.AlphaCO)},
		{"alpha_SP", fmt.Sprintf("%.6f", r.AlphaSP)},
		{"wedge", fmt.Sprintf("%.6f", r.Wedge)},
		{"K_NE", fmt.Sprintf("%.6f", r.OwnerSurplus)},
		{"W_NE", fmt.Sprintf("%.6f", r.WorkerIncome)},
		{"K_CO", fmt.Sprintf("%.6f", r.OwnerSurplusCO)},
		{"surplus_loss", fmt.Sprintf("%.6f", r.SurplusLoss)},
		{"tau_star", fmt.Sprintf("%.6f", r.TaxStar)},
		{"alpha_taxed", fmt.Sprintf("%.6f", r.AlphaTaxed)},
	}
	for _, row := range rows {
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
