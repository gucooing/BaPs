package zset

import (
	"math/rand"
	"testing"
)

var s *SortedSet[int64]

func init() {
	s = New[int64]()
}

func TestNew(t *testing.T) {
	if s == nil {
		t.Failed()
	}
	s.Set(66, 1001)
	s.Set(77, 1002)
	s.Set(88, 1003)
	s.Set(100, 1004)
	s.Set(99, 1005)
	s.Set(44, 1006)
	s.Set(44, 1001)

	rank, score := s.GetRank(1004, false)
	if rank == 5 {
		t.Log("Key:", 1004, "Rank:", rank, "Score:", score)
	} else {
		t.Error("Key:", 1004, "Rank:", rank, "Score:", score)
	}
	rank, score = s.GetRank(1001, false)
	if rank == 0 {
		t.Log("Key:", 1001, "Rank:", rank, "Score:", score)
	} else {
		t.Error("Key:", 1001, "Rank:", rank, "Score:", score)
	}
	rank, score = s.GetRank(-1, false)
	if rank == -1 {
		t.Log("Key:", -1, "Rank:", rank, "Score:", score)
	} else {
		t.Error("Key:", -1, "Rank:", rank, "Score:", score)
	}

	id, score := s.GetDataByRank(0, true)
	t.Log("GetData[REVERSE] Rank:", 0, "ID:", id, "Score:", score)
	id, score = s.GetDataByRank(0, false)
	t.Log("GetData[UNREVERSE] Rank:", 0, "ID:", id, "Score:", score)
	_, _ = s.GetDataByRank(9999, true)
	if s.Length() != 6 {
		t.Error("Rank Data Size is wrong")
	}
	s.Delete(1001)
	if s.Length() != 5 {
		t.Error("Rank Data Size is wrong")
	}
	curScore := s.IncrBy(666, 1004)
	t.Log(curScore)
}

func TestIncrBy(t *testing.T) {
	z := New[int64]()
	for i := 1000; i < 1100; i++ {
		z.Set(float64(i), int64(i))
	}
	rank, score := z.GetRank(1050, false)
	curScore := z.IncrBy(1.5, 1050)
	if score+1.5 != curScore {
		t.Error(score, curScore)
	}
	r2, score2 := z.GetRank(1050, false)
	if score2 != curScore {
		t.Fail()
	}
	if r2 != rank+1 {
		t.Error(r2, rank)
	}

}

func TestRange(t *testing.T) {
	z := New[int64]()
	z.Set(1.0, 1001)
	z.Set(2.0, 1002)
	z.Set(3.0, 1003)
	z.Set(4.0, 1004)
	z.Set(5.0, 1005)
	z.Set(6.0, 1006)

	ids := make([]int64, 0, 6)
	z.Range(0, -1, func(score float64, k int64) {
		ids = append(ids, k)
		t.Log(score, k)
	})
	if ids[0] != 1001 ||
		ids[1] != 1002 ||
		ids[2] != 1003 ||
		ids[3] != 1004 {
		t.Fail()
	}
	z.RevRange(1, 3, func(score float64, k int64) {
		t.Log(score, k)
	})

}

func BenchmarkSortedSet_Add(b *testing.B) {
	b.StopTimer()
	// data initialization
	scores := make([]float64, b.N)
	IDs := make([]int64, b.N)
	for i := range IDs {
		scores[i] = rand.Float64() + float64(rand.Int31n(99))
		IDs[i] = int64(i) + 100000
	}
	// BCE
	_ = scores[:b.N]
	_ = IDs[:b.N]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Set(scores[i], IDs[i])
	}
}

func BenchmarkSortedSet_GetRank(b *testing.B) {
	l := s.Length()
	for i := 0; i < b.N; i++ {
		s.GetRank(100000+int64(i)%l, true)
	}
}

func BenchmarkSortedSet_GetDataByRank(b *testing.B) {
	l := s.Length()
	for i := 0; i < b.N; i++ {
		s.GetDataByRank(int64(i)%l, true)
	}
}
