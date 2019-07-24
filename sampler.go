package goptuna

import (
	"errors"
	"math/rand"
)

// Sampler returns the next search points
type Sampler interface {
	// Sample a parameter for a given distribution.
	Sample(*Study, FrozenTrial, string, interface{}) (float64, error)
}

var _ Sampler = &RandomSearchSampler{}

// RandomSearchSampler for random search
type RandomSearchSampler struct {
	rng *rand.Rand
}

func NewRandomSearchSampler() *RandomSearchSampler {
	return &RandomSearchSampler{
		rng: rand.New(rand.NewSource(0)),
	}
}

func (s *RandomSearchSampler) Sample(
	study *Study,
	trial FrozenTrial,
	paramName string,
	paramDistribution interface{},
) (float64, error) {
	switch d := paramDistribution.(type) {
	case UniformDistribution:
		return s.rng.Float64()*(d.Max-d.Min) + d.Min, nil
	default:
		return 0.0, errors.New("undefined distribution")
	}
}