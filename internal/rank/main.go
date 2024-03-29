package rank

import (
	"math"
	"sort"

	log "github.com/sirupsen/logrus"
)

// Rank does StuffRank on three int scores and returns their new scores in same order
func Rank(s [3]int) [3]int {
	sr := newStuffRank(s[0], s[1], s[2])
	sr.expected()
	sr.apply()
	sn := sr.scores
	return [3]int{sn[0].delta, sn[1].delta, sn[2].delta}
}

type score struct {
	score, rank, delta int
	e                  float64
}

type stuffRank struct {
	scores       [3]*score
	sortExpected bool
}

func newStuffRank(s1, s2, s3 int) *stuffRank {
	a := []int{s1, s2, s3}
	sr := new(stuffRank)
	for i := 0; i < 3; i++ {
		sr.scores[i] = &score{
			score: a[i],
			rank:  i,
		}
	}
	return sr
}

func (r *stuffRank) expected() {
	s := r.scores
	s1, s2, s3 := s[0].score, s[1].score, s[2].score

	log.WithFields(log.Fields{"s1": s1, "s2": s2, "s3": s3}).Debug("rank.expected")
	defer func() {
		log.WithFields(log.Fields{"s1e": s[0].e, "s2e": s[1].e, "s3e": s[2].e}).Debug("rank.expected")
	}()
	r.sortExpected = true
	r.sortExpected = true
	sort.Sort(r)
	avgs := [3]float64{avg(s2, s3), avg(s1, s3), avg(s1, s2)}
	for i, v := range r.scores {
		v.e = e(avgs[i], v.score)
	}
	r.sortExpected = false
	sort.Sort(r)
}

func (r *stuffRank) apply() {
	s := r.scores
	s1, s2, s3 := s[0], s[1], s[2]

	// Compute deltas
	s1.delta = int(30 * (1 - s1.e))
	s3.delta = int(-30 * s3.e)
	s2.delta = -(s1.delta + s3.delta)
}

func e(sOther float64, sSelf int) float64 {
	return 1.0 / (1. + math.Pow(10, (sOther-float64(sSelf))/400.))
}

func avg(a, b int) float64 {
	return float64(a+b) / 2
}

func (r *stuffRank) Len() int {
	return len(r.scores)
}

func (r *stuffRank) Less(i, j int) bool {
	if r.sortExpected {
		return r.scores[i].score < r.scores[j].score
	}
	return r.scores[i].rank < r.scores[j].rank
}

func (r *stuffRank) Swap(i, j int) {
	r.scores[i], r.scores[j] = r.scores[j], r.scores[i]
}
