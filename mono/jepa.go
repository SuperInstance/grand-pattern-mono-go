package mono

import (
	"errors"
	"math"
)

// Jepa is a weighted-history predictor for a mono vibe stream.
// It uses an exponentially-weighted moving average.
type Jepa struct {
	alpha   float64
	history []float64
}

// NewJepa creates a Jepa with the given smoothing factor.
func NewJepa(alpha float64) (*Jepa, error) {
	if alpha <= 0 || alpha > 1 {
		return nil, errors.New("alpha must be in (0, 1]")
	}
	return &Jepa{alpha: alpha}, nil
}

// Observe records a new vibe observation.
func (j *Jepa) Observe(vibe float64) {
	j.history = append(j.history, vibe)
}

// Predict returns the predicted next vibe.
func (j *Jepa) Predict() float64 {
	if len(j.history) == 0 {
		return 0
	}
	n := len(j.history)
	total := 0.0
	weightSum := 0.0
	for i, v := range j.history {
		age := float64(n - 1 - i)
		w := j.alpha * math.Pow(1-j.alpha, age)
		total += w * v
		weightSum += w
	}
	if weightSum == 0 {
		return 0
	}
	return total / weightSum
}

// Surprise returns |observed - predicted|.
func (j *Jepa) Surprise(observed float64) float64 {
	return math.Abs(observed - j.Predict())
}

// History returns a copy of the observation history.
func (j *Jepa) History() []float64 {
	out := make([]float64, len(j.history))
	copy(out, j.history)
	return out
}

// ObservationCount returns the number of observations.
func (j *Jepa) ObservationCount() int {
	return len(j.history)
}
