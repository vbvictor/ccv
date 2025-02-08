package plot

import (
	"errors"
	"fmt"

	"github.com/go-echarts/go-echarts/v2/opts"
)

type RiskLevel struct {
	Name  string
	Color string
	Min   uint
	Max   uint
}

var ErrInvalidRiskThreshold = errors.New("invalid risk threshold")

// Need to make it more general: TODO refactor.
func ValidateRiskThresholds() error {
	if VeryLowRisk >= LowRisk {
		return fmt.Errorf("%w: Very Low Risk threshold (%d) must be less than Low Risk threshold (%d)",
			ErrInvalidRiskThreshold, VeryLowRisk, LowRisk)
	}

	if LowRisk >= MediumRisk {
		return fmt.Errorf("%w: Low Risk threshold (%d) must be less than Medium Risk threshold (%d)",
			ErrInvalidRiskThreshold, LowRisk, MediumRisk)
	}

	if MediumRisk >= HighRisk {
		return fmt.Errorf("%w: Medium Risk threshold (%d) must be less than High Risk threshold (%d)",
			ErrInvalidRiskThreshold, MediumRisk, HighRisk)
	}

	if HighRisk >= VeryHighRisk {
		return fmt.Errorf("%w: High Risk threshold (%d) must be less than Very High Risk threshold (%d)",
			ErrInvalidRiskThreshold, HighRisk, VeryHighRisk)
	}

	if VeryHighRisk >= CriticalRisk {
		return fmt.Errorf("%w: Very High Risk threshold (%d) must be less than Critical Risk threshold (%d)",
			ErrInvalidRiskThreshold, VeryHighRisk, CriticalRisk)
	}

	return nil
}

func getRiskLevels() []RiskLevel {
	return []RiskLevel{
		{Name: "Very Low Risk", Color: "#90EE90", Min: VeryLowRisk, Max: LowRisk - 1},
		{Name: "Low Risk", Color: "#47d147", Min: LowRisk, Max: MediumRisk - 1},
		{Name: "Medium Risk", Color: "#ffd700", Min: MediumRisk, Max: HighRisk - 1},
		{Name: "High Risk", Color: "#ffa64d", Min: HighRisk, Max: VeryHighRisk - 1},
		{Name: "Very High Risk", Color: "#ff4d4d", Min: VeryHighRisk, Max: CriticalRisk - 1},
		{Name: "Critical Risk", Color: "#8b0000", Min: CriticalRisk, Max: ^uint(0)},
	}
}

type RisksMapper struct {
	levels []RiskLevel
}

func NewRisksMapper() *RisksMapper {
	return &RisksMapper{
		levels: getRiskLevels(),
	}
}

var _ EntryMapper = (*RisksMapper)(nil)

func (rm *RisksMapper) Map(data ScatterData) Category {
	riskScore := data.Complexity + float64(data.Churn)

	for _, level := range rm.levels {
		if riskScore >= float64(level.Min) && riskScore <= float64(level.Max) {
			return level.Name
		}
	}

	return "Unknown"
}

func (rm *RisksMapper) Style(category Category) opts.ItemStyle {
	for _, level := range rm.levels {
		if level.Name == category {
			return opts.ItemStyle{
				Color: level.Color,
			}
		}
	}

	return opts.ItemStyle{}
}
