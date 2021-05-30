//go:generate mockgen -destination=random_mock.go -package=goswarm -source=./random.go random
package goswarm

type random interface {
	next(lower float64, upper float64) float64
}
