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
	expect := [3]int{1015, 1000, 985}
	_testRank(t, in, expect)
}

func TestRank_TwoEqual(t *testing.T) {
	in := [3]int{1000, 1000, 1100}
	expect := [3]int{1017, 1002, 1081}
	_testRank(t, in, expect)
}

func TestRank_TwoEqualE1(t *testing.T) {
	in := [3]int{1000, 1100, 1000}
	expect := [3]int{1017, 1096, 988}
	_testRank(t, in, expect)
}

func TestRank_TwoEqualE2(t *testing.T) {
	in := [3]int{1100, 1000, 1000}
	expect := [3]int{1110, 1002, 988}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E1(t *testing.T) {
	in := [3]int{1200, 1100, 1000}
	expect := [3]int{1208, 1100, 992}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E15(t *testing.T) {
	in := [3]int{1200, 1000, 1100}
	expect := [3]int{1208, 1006, 1085}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E2(t *testing.T) {
	in := [3]int{1100, 1200, 1000}
	expect := [3]int{1115, 1194, 992}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E25(t *testing.T) {
	in := [3]int{1100, 1000, 1200}
	expect := [3]int{1115, 1006, 1179}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E3(t *testing.T) {
	in := [3]int{1000, 1100, 1200}
	expect := [3]int{1021, 1100, 1179}
	_testRank(t, in, expect)
}

func TestRank_OffBy100E35(t *testing.T) {
	in := [3]int{1000, 1200, 1100}
	expect := [3]int{1021, 1194, 1085}
	_testRank(t, in, expect)
}
