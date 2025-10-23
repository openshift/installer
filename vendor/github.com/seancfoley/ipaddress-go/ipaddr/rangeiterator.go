//
// Copyright 2020-2022 Sean C Foley
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ipaddr

type sequRangeIterator[T SequentialRangeConstraint[T]] struct {
	rng                 *SequentialRange[T]
	creator             func(T, T) *SequentialRange[T]
	prefixBlockIterator Iterator[T]
	prefixLength        BitCount
	notFirst            bool
}

func (it *sequRangeIterator[T]) HasNext() bool {
	return it.prefixBlockIterator.HasNext()
}

func (it *sequRangeIterator[T]) Next() (res *SequentialRange[T]) {
	if it.HasNext() {
		next := it.prefixBlockIterator.Next()
		if !it.notFirst {
			it.notFirst = true
			// next is a prefix block
			lower := it.rng.GetLower()
			prefLen := it.prefixLength
			if it.HasNext() {
				if !lower.IncludesZeroHostLen(prefLen) {
					return it.creator(lower, next.GetUpper())
				}
			} else {
				upper := it.rng.GetUpper()
				if !lower.IncludesZeroHostLen(prefLen) || !upper.IncludesMaxHostLen(prefLen) {
					return it.creator(lower, upper)
				}
			}
		} else if !it.HasNext() {
			upper := it.rng.GetUpper()
			if !upper.IncludesMaxHostLen(it.prefixLength) {
				return it.creator(next.GetLower(), upper)
			}
		}
		lower, upper := next.getLowestHighestAddrs()
		return newSequRangeUnchecked(lower, upper, lower != upper)
	}
	return
}
