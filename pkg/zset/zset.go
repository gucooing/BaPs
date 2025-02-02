// Package zset is a port of t_zset.c in Redis
/*
 * Copyright (c) 2009-2012, Salvatore Sanfilippo <antirez at gmail dot com>
 * Copyright (c) 2009-2012, Pieter Noordhuis <pcnoordhuis at gmail dot com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *   * Neither the name of Redis nor the names of its contributors may be used
 *     to endorse or promote products derived from this software without
 *     specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */
package zset

import (
	"math/rand"
	"sync"

	"github.com/gucooing/BaPs/pkg/zset/cmp"
)

const zSkiplistMaxlevel = 32

type (
	// Key constraint
	Key interface {
		cmp.Ordered
	}
	skipListLevel[K Key] struct {
		forward *skipListNode[K]
		span    uint64
	}

	skipListNode[K Key] struct {
		objID    K
		score    float64
		backward *skipListNode[K]
		level    []*skipListLevel[K]
	}
	obj struct {
		score float64
	}

	skipList[K Key] struct {
		header *skipListNode[K]
		tail   *skipListNode[K]
		length int64
		level  int16
	}
	// SortedSet is the final exported sorted set we can use
	SortedSet[K Key] struct {
		dict map[K]float64
		zsl  *skipList[K]
		lock sync.RWMutex
	}
	zrangespec struct {
		min   float64
		max   float64
		minex int32
		maxex int32
	}
	zlexrangespec[K Key] struct {
		minKey K
		maxKey K
		minex  int
		maxex  int
	}
)

func zslCreateNode[K Key](level int16, score float64, id K) *skipListNode[K] {
	n := &skipListNode[K]{
		score: score,
		objID: id,
		level: make([]*skipListLevel[K], level),
	}
	for i := range n.level {
		n.level[i] = new(skipListLevel[K])
	}
	return n
}

func zslCreate[K Key]() *skipList[K] {
	return &skipList[K]{
		level:  1,
		header: zslCreateNode[K](zSkiplistMaxlevel, 0, *new(K)),
	}
}

const zSkiplistP = 0.25 /* Skiplist P = 1/4 */

/* Returns a random level for the new skiplist node we are going to create.
 * The return value of this function is between 1 and _ZSKIPLIST_MAXLEVEL
 * (both inclusive), with a powerlaw-alike distribution where higher
 * levels are less likely to be returned. */
func randomLevel() int16 {
	level := int16(1)
	for float32(rand.Int31()&0xFFFF) < (zSkiplistP * 0xFFFF) {
		level++
	}
	if level < zSkiplistMaxlevel {
		return level
	}
	return zSkiplistMaxlevel
}

/* zslInsert a new node in the skiplist. Assumes the element does not already
 * exist (up to the caller to enforce that). The skiplist takes ownership
 * of the passed SDS string 'obj'. */
func (zsl *skipList[K]) zslInsert(score float64, id K) *skipListNode[K] {
	update := make([]*skipListNode[K], zSkiplistMaxlevel)
	rank := make([]uint64, zSkiplistMaxlevel)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		/* store rank that is crossed to reach the insert position */
		if i == zsl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		if x.level[i] != nil {
			for x.level[i].forward != nil &&
				(x.level[i].forward.score < score ||
					(x.level[i].forward.score == score && x.level[i].forward.objID < id)) {
				rank[i] += x.level[i].span
				x = x.level[i].forward
			}
		}
		update[i] = x
	}
	/* we assume the element is not already inside, since we allow duplicated
	 * scores, reinserting the same element should never happen since the
	 * caller of zslInsert() should test in the hash table if the element is
	 * already inside or not. */
	level := randomLevel()
	if level > zsl.level {
		for i := zsl.level; i < level; i++ {
			rank[i] = 0
			update[i] = zsl.header
			update[i].level[i].span = uint64(zsl.length)
		}
		zsl.level = level
	}
	x = zslCreateNode(level, score, id)
	for i := int16(0); i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		/* update span covered by update[i] as x is inserted here */
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	/* increment span for untouched levels */
	for i := level; i < zsl.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == zsl.header {
		x.backward = nil
	} else {
		x.backward = update[0]

	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		zsl.tail = x
	}
	zsl.length++
	return x
}

/* Internal function used by zslDelete, zslDeleteByScore and zslDeleteByRank */
func (zsl *skipList[K]) zslDeleteNode(x *skipListNode[K], update []*skipListNode[K]) {
	for i := int16(0); i < zsl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		zsl.tail = x.backward
	}
	for zsl.level > 1 && zsl.header.level[zsl.level-1].forward == nil {
		zsl.level--
	}
	zsl.length--
}

/* Delete an element with matching score/element from the skiplist.
 * The function returns 1 if the node was found and deleted, otherwise
 * 0 is returned.
 *
 * If 'node' is NULL the deleted node is freed by zslFreeNode(), otherwise
 * it is not freed (but just unlinked) and *node is set to the node pointer,
 * so that it is possible for the caller to reuse the node (including the
 * referenced SDS string at node->obj). */
func (zsl *skipList[K]) zslDelete(score float64, id K) int {
	update := make([]*skipListNode[K], zSkiplistMaxlevel)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					x.level[i].forward.objID < id)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object. */
	x = x.level[0].forward
	if x != nil && score == x.score && x.objID == id {
		zsl.zslDeleteNode(x, update)
		return 1
	}
	return 0 /* not found */
}

func zslValueGteMin(value float64, spec *zrangespec) bool {
	if spec.minex != 0 {
		return value > spec.min
	}
	return value >= spec.min
}

func zslValueLteMax(value float64, spec *zrangespec) bool {
	if spec.maxex != 0 {
		return value < spec.max
	}
	return value <= spec.max
}

/* Returns if there is a part of the zset is in range. */
func (zsl *skipList[K]) zslIsInRange(ran *zrangespec) bool {
	/* Test for ranges that will always be empty. */
	if ran.min > ran.max ||
		(ran.min == ran.max && (ran.minex != 0 || ran.maxex != 0)) {
		return false
	}
	x := zsl.tail
	if x == nil || !zslValueGteMin(x.score, ran) {
		return false
	}
	x = zsl.header.level[0].forward
	if x == nil || !zslValueLteMax(x.score, ran) {
		return false
	}
	return true
}

/* Find the first node that is contained in the specified range.
 * Returns NULL when no element is contained in the range. */
func (zsl *skipList[K]) zslFirstInRange(ran *zrangespec) *skipListNode[K] {
	/* If everything is out of range, return early. */
	if !zsl.zslIsInRange(ran) {
		return nil
	}

	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		/* Go forward while *OUT* of range. */
		for x.level[i].forward != nil &&
			!zslValueGteMin(x.level[i].forward.score, ran) {
			x = x.level[i].forward
		}
	}
	/* This is an inner range, so the next node cannot be NULL. */
	x = x.level[0].forward
	// serverAssert(x != NULL);

	/* Check if score <= max. */
	if !zslValueLteMax(x.score, ran) {
		return nil
	}
	return x
}

/* Find the last node that is contained in the specified range.
 * Returns NULL when no element is contained in the range. */
func (zsl *skipList[K]) zslLastInRange(ran *zrangespec) *skipListNode[K] {

	/* If everything is out of range, return early. */
	if !zsl.zslIsInRange(ran) {
		return nil
	}
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		/* Go forward while *IN* range. */
		for x.level[i].forward != nil &&
			zslValueLteMax(x.level[i].forward.score, ran) {
			x = x.level[i].forward
		}
	}
	/* This is an inner range, so this node cannot be NULL. */
	// serverAssert(x != NULL);

	/* Check if score >= min. */
	if !zslValueGteMin(x.score, ran) {
		return nil
	}
	return x
}

/* Delete all the elements with score between min and max from the skiplist.
 * Min and max are inclusive, so a score >= min || score <= max is deleted.
 * Note that this function takes the reference to the hash table view of the
 * sorted set, in order to remove the elements from the hash table too. */
func (zsl *skipList[K]) zslDeleteRangeByScore(ran *zrangespec, dict map[K]float64) uint64 {
	removed := uint64(0)
	update := make([]*skipListNode[K], zSkiplistMaxlevel)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil {
			var condition bool
			if ran.minex != 0 {
				condition = x.level[i].forward.score <= ran.min
			} else {
				condition = x.level[i].forward.score < ran.min
			}
			if !condition {
				break
			}
			x = x.level[i].forward
		}
		update[i] = x
	}

	/* Current node is the last with score < or <= min. */
	x = x.level[0].forward

	/* Delete nodes while in range. */
	for x != nil {
		var condition bool
		if ran.maxex != 0 {
			condition = x.score < ran.max
		} else {
			condition = x.score <= ran.max
		}
		if !condition {
			break
		}
		next := x.level[0].forward
		zsl.zslDeleteNode(x, update)
		delete(dict, x.objID)
		// Here is where x->obj is actually released.
		// And golang has GC, don't need to free manually anymore
		// zslFreeNode(x)
		removed++
		x = next
	}
	return removed
}

func (zsl *skipList[K]) zslDeleteRangeByLex(ran *zlexrangespec[K], dict map[K]float64) uint64 {
	removed := uint64(0)

	update := make([]*skipListNode[K], zSkiplistMaxlevel)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && !zslLexValueGteMin(x.level[i].forward.objID, ran) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	/* Current node is the last with score < or <= min. */
	x = x.level[0].forward

	/* Delete nodes while in range. */
	for x != nil && zslLexValueLteMax(x.objID, ran) {
		next := x.level[0].forward
		zsl.zslDeleteNode(x, update)
		delete(dict, x.objID)
		removed++
		x = next
	}
	return removed
}

func zslLexValueGteMin[K Key](id K, spec *zlexrangespec[K]) bool {
	if spec.minex != 0 {
		return compareKey(id, spec.minKey) > 0
	}
	return compareKey(id, spec.minKey) >= 0
}

func compareKey[K Key](a, b K) int8 {
	if a == b {
		return 0
	} else if a > b {
		return 1
	}
	return -1
}

func zslLexValueLteMax[K Key](id K, spec *zlexrangespec[K]) bool {
	if spec.maxex != 0 {
		return compareKey(id, spec.maxKey) < 0
	}
	return compareKey(id, spec.maxKey) <= 0
}

/* Delete all the elements with rank between start and end from the skiplist.
 * Start and end are inclusive. Note that start and end need to be 1-based */
func (zsl *skipList[K]) zslDeleteRangeByRank(start, end uint64, dict map[K]float64) uint64 {
	update := make([]*skipListNode[K], zSkiplistMaxlevel)
	var traversed, removed uint64

	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span) < start {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	traversed++
	x = x.level[0].forward
	for x != nil && traversed <= end {
		next := x.level[0].forward
		zsl.zslDeleteNode(x, update)
		delete(dict, x.objID)
		removed++
		traversed++
		x = next
	}
	return removed
}

/* Find the rank for an element by both score and obj.
 * Returns 0 when the element cannot be found, rank otherwise.
 * Note that the rank is 1-based due to the span of zsl->header to the
 * first element. */
func (zsl *skipList[K]) zslGetRank(score float64, key K) int64 {
	rank := uint64(0)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					x.level[i].forward.objID <= key)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		/* x might be equal to zsl->header, so test if obj is non-NULL */
		if x.objID == key {
			return int64(rank)
		}
	}
	return 0
}

/* Finds an element by its rank. The rank argument needs to be 1-based. */
func (zsl *skipList[K]) zslGetElementByRank(rank uint64) *skipListNode[K] {
	traversed := uint64(0)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span) <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
}

/*-----------------------------------------------------------------------------
 * Common sorted set API
 *----------------------------------------------------------------------------*/

// New creates a new SortedSet and return its pointer
func New[K Key]() *SortedSet[K] {
	s := &SortedSet[K]{
		dict: make(map[K]float64),
		zsl:  zslCreate[K](),
		lock: sync.RWMutex{},
	}
	return s
}

// Length returns counts of elements
func (z *SortedSet[K]) Length() int64 {
	z.lock.RLock()
	defer z.lock.RUnlock()
	return z.zsl.length
}

// Set is used to add or update an element
func (z *SortedSet[K]) Set(score float64, key K) {
	z.lock.Lock()
	defer z.lock.Unlock()
	v, ok := z.dict[key]
	z.dict[key] = score
	if ok {
		/* Remove and re-insert when score changes. */
		if score != v {
			z.zsl.zslDelete(v, key)
			z.zsl.zslInsert(score, key)
		}
	} else {
		z.zsl.zslInsert(score, key)
	}
}

// IncrBy ..
func (z *SortedSet[K]) IncrBy(score float64, key K) float64 {
	z.lock.Lock()
	defer z.lock.Unlock()
	oldScore, ok := z.dict[key]
	if !ok {
		z.Set(score, key)
		return score
	}
	if score != 0 {
		z.zsl.zslDelete(oldScore, key)
		z.dict[key] += score
		z.zsl.zslInsert(z.dict[key], key)
	}
	return z.dict[key]
}

// Delete removes an element from the SortedSet
// by its key.
func (z *SortedSet[K]) Delete(key K) (ok bool) {
	z.lock.Lock()
	defer z.lock.Unlock()
	score, ok := z.dict[key]
	if ok {
		z.zsl.zslDelete(score, key)
		delete(z.dict, key)
		return true
	}
	return false
}

// GetRank returns position,score and extra data of an element which
// found by the parameter key.
// The parameter reverse determines the rank is descent or ascend，
// true means descend and false means ascend.
func (z *SortedSet[K]) GetRank(key K, reverse bool) (rank int64, score float64) {
	z.lock.RLock()
	defer z.lock.RUnlock()
	score, ok := z.dict[key]
	if !ok {
		return -1, 0
	}
	r := z.zsl.zslGetRank(score, key)
	if reverse {
		r = z.zsl.length - r
	} else {
		r--
	}
	return int64(r), score

}

// GetScore implements ZScore
func (z *SortedSet[K]) GetScore(key K) (score float64, ok bool) {
	z.lock.RLock()
	defer z.lock.RUnlock()
	score, ok = z.dict[key]
	return score, ok
}

// GetDataByRank returns the id,score and extra data of an element which
// found by position in the rank.
// The parameter rank is the position, reverse says if in the descend rank.
func (z *SortedSet[K]) GetDataByRank(rank int64, reverse bool) (key K, score float64) {
	z.lock.RLock()
	defer z.lock.RUnlock()
	if rank < 0 || rank > z.zsl.length {
		return *new(K), 0
	}
	if reverse {
		rank = z.zsl.length - rank
	} else {
		rank++
	}
	n := z.zsl.zslGetElementByRank(uint64(rank))
	if n == nil {
		return *new(K), 0
	}
	score, ok := z.dict[n.objID]
	if !ok {
		return *new(K), 0
	}
	return n.objID, score
}

// Range implements ZRANGE
func (z *SortedSet[K]) Range(start, end int64, f func(float64, K)) {
	z.snapshotRange(start, end, false, f)
}

// RevRange implements ZREVRANGE
func (z *SortedSet[K]) RevRange(start, end int64, f func(float64, K)) {
	z.snapshotRange(start, end, true, f)
}

func (z *SortedSet[K]) snapshotRange(start, end int64, reverse bool, f func(float64, K)) {
	scores := make([]float64, 0)
	keys := make([]K, 0)

	z.lock.RLock()
	z.commonRange(start, end, reverse, func(f float64, k K) {
		scores = append(scores, f)
		keys = append(keys, k)
	})
	z.lock.RUnlock()

	for i, score := range scores {
		f(score, keys[i])
	}
}

func (z *SortedSet[K]) commonRange(start, end int64, reverse bool, f func(float64, K)) {
	l := z.zsl.length
	if start < 0 {
		start += l
		if start < 0 {
			start = 0
		}
	}
	if end < 0 {
		end += l
	}

	if start > end || start >= l {
		return
	}
	if end >= l {
		end = l - 1
	}
	span := (end - start) + 1

	var node *skipListNode[K]
	if reverse {
		node = z.zsl.tail
		if start > 0 {
			node = z.zsl.zslGetElementByRank(uint64(l - start))
		}
	} else {
		node = z.zsl.header.level[0].forward
		if start > 0 {
			node = z.zsl.zslGetElementByRank(uint64(start + 1))
		}
	}
	for span > 0 {
		span--
		k := node.objID
		s := node.score
		f(s, k)
		if reverse {
			node = node.backward
		} else {
			node = node.level[0].forward
		}
	}
}
