package rank_test

import (
	"testing"

	"github.com/stuff-ai/api/internal/rank"
)

func _testRank(t *testing.T, in, expect [3]int) {
	actual := rank.Rank(in)
	t.Log(actual)
	for i, a := range actual {
		e := expect[i]
		if a != e {
			t.Fatalf("Failure in Rank %d: expected %d, got %d", i, e, a)
		}
	}
}

func TestRank_AllEqual(t *testing.T) {
	in := [3]int{1000, 1000, 1000}
	expect := [3]int{15, 0, -15}
	_testRank(t, in, expect)
}

func TestRank_TwoEqual(t *testing.T) {
	in := [3]int{1000, 1000, 1100}
	expect := [3]int{17, 2, -19}
	_testRank(t, in, expect)
}

func TestRank_TwoEqualE1(t *testing.T) {
	in := [3]int{1000, 1100, 1000}
	expect := [3]int{17, -2, -15}
	_testRank(t, in, expect)
}

func TestRank_TwoEqualE2(t *testing.T) {
	in := [3]int{1100, 1000, 1000}
	expect := [3]int{10, 2, -12}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E1(t *testing.T) {
	in := [3]int{1200, 1100, 1000}
	expect := [3]int{8, 1100, -8}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E15(t *testing.T) {
	in := [3]int{1200, 1000, 1100}
	expect := [3]int{10, 2, -12}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E2(t *testing.T) {
	in := [3]int{1100, 1200, 1000}
	expect := [3]int{12, -2, -10}
	_testRank(t, in, expect)
}

func TestRank_OffBy101E25(t *testing.T) {
	in := [3]int{1100, 1000, 1200}
	expect := [3]int{17, 4, -21}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E3(t *testing.T) {
	in := [3]int{1000, 1100, 1200}
	expect := [3]int{21, 0, -21}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E35(t *testing.T) {
	in := [3]int{1000, 1200, 1100}
	expect := [3]int{21, -4, -17}
	_testRank(t, in, expect)
}

func TestRank_OffBy200E1(t *testing.T) {
	in := [3]int{1200, 1000, 800}
	expect := [3]int{10, 0, -10}
	_testRank(t, in, expect)
}

func TestRank_OffBy400E1(t *testing.T) {
	in := [3]int{1400, 1000, 600}
	expect := [3]int{7, 0, -7}
	_testRank(t, in, expect)
}
