package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildAccumulator2(t *testing.T) {
	accumulator := Accumulator{
		Duration:           2 * time.Minute,
		CurrentPeriodUsed:  3,
		CurrentPeriodSpend: 500,
		PeriodEndsDate:     time.Date(2021, 10, 26, 10, 4, 0, 0, time.Local),
	}

	accumulator.AddSpend(500)

	assert.Equal(t, int64(1000), accumulator.CurrentPeriodSpend)
	assert.Equal(t, int64(4), accumulator.CurrentPeriodUsed)
}

type args struct {
	currentTime        time.Time
	duration           time.Duration
	currentPeriodUsed  int64
	currentPeriodSpend int64
	periodEndsDate     time.Time
}

func TestBuildAccumulator(t *testing.T) {
	testCases := []struct {
		name           string
		args           args
		expectedResult Accumulator
	}{
		{
			name: "buildando um acumulador que não precisa ser resetado",
			args: args{
				currentTime:        time.Date(2021, 10, 26, 10, 2, 30, 0, time.Local),
				duration:           2 * time.Minute,
				currentPeriodUsed:  3,
				currentPeriodSpend: 500,
				periodEndsDate:     time.Date(2021, 10, 26, 10, 4, 0, 0, time.Local),
			},
			expectedResult: Accumulator{
				Duration:           2 * time.Minute,
				CurrentPeriodUsed:  3,
				CurrentPeriodSpend: 500,
				PeriodEndsDate:     time.Date(2021, 10, 26, 10, 4, 0, 0, time.Local),
			},
		},
		{
			name: "buildando um acumulador que precisa ser resetado",
			args: args{
				currentTime:        time.Date(2021, 10, 26, 10, 2, 30, 0, time.Local),
				duration:           2 * time.Minute,
				currentPeriodUsed:  3,
				currentPeriodSpend: 500,
				periodEndsDate:     time.Date(2021, 10, 26, 10, 2, 0, 0, time.Local),
			},
			expectedResult: Accumulator{
				Duration:           2 * time.Minute,
				CurrentPeriodUsed:  0,
				CurrentPeriodSpend: 0,
				PeriodEndsDate:     time.Date(2021, 10, 26, 10, 4, 30, 0, time.Local),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accumulator := BuildAccumulator(
				tt.args.currentTime,
				tt.args.duration,
				tt.args.currentPeriodUsed,
				tt.args.currentPeriodSpend,
				tt.args.periodEndsDate,
			)

			assert.Equal(t, tt.expectedResult, accumulator)
		})
	}
}
