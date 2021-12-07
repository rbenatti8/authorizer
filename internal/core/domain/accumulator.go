package domain

import "time"

type Accumulator struct {
	Duration           time.Duration
	CurrentPeriodUsed  int64
	CurrentPeriodSpend int64
	PeriodEndsDate     time.Time
}

func (ac *Accumulator) AddSpend(amount int64) {
	ac.CurrentPeriodSpend = ac.CurrentPeriodSpend + amount
	ac.CurrentPeriodUsed++
}

func BuildAccumulator(
	currentTime time.Time,
	duration time.Duration,
	currentPeriodUsed int64,
	currentPeriodSpend int64,
	periodEndsDate time.Time,
) Accumulator {
	if currentTime.After(periodEndsDate) {
		return Accumulator{
			Duration:           duration,
			CurrentPeriodUsed:  0,
			CurrentPeriodSpend: 0,
			PeriodEndsDate:     currentTime.Add(duration),
		}
	}

	return Accumulator{
		Duration:           duration,
		CurrentPeriodUsed:  currentPeriodUsed,
		CurrentPeriodSpend: currentPeriodSpend,
		PeriodEndsDate:     periodEndsDate,
	}
}
