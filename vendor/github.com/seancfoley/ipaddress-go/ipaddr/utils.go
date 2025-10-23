//
// Copyright 2020-2024 Sean C Foley
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

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"unsafe"
)

func nilString() string {
	return "<nil>"
}

// nilSection prints a string for sections with a nil division slice or division slice of 0 length.
// For division groupings, the division slice string is generated from using the slice, see toString() or defaultFormat() in grouping code.
func nilSection() string {
	return ""
}

func clone[T any](orig []T) []T {
	return append(make([]T, 0, len(orig)), orig...)
}

func cloneSeries[T any](addr T, orig []T) []T {
	return append(append(make([]T, 0, len(orig)+1), addr), orig...)
}

func fillDivs(orig []*AddressDivision, val *AddressDivision) {
	for i := range orig {
		orig[i] = val
	}
}

// copies cached into bytes, unless bytes is too small, in which case cached is cloned
func getBytesCopy(bytes, cached []byte) []byte {
	if bytes == nil || len(bytes) < len(cached) {
		return clone(cached)
	}
	copy(bytes, cached)
	return bytes[:len(cached)]
}

// note: only to be used when you already know the total size fits into a long
func longCount(section *AddressSection, segCount int) uint64 {
	result := getLongCount(func(index int) uint64 { return section.GetSegment(index).GetValueCount() }, segCount)
	return result
}

func getLongCount(segmentCountProvider func(index int) uint64, segCount int) uint64 {
	if segCount <= 0 {
		return 1
	}
	result := segmentCountProvider(0)
	for i := 1; i < segCount; i++ {
		result *= segmentCountProvider(i)
	}
	return result
}

// note: only to be used when you already know the total size fits into a long
func longPrefixCount(section *AddressSection, prefixLength BitCount) uint64 {
	bitsPerSegment := section.GetBitsPerSegment()
	bytesPerSegment := section.GetBytesPerSegment()
	networkSegmentIndex := getNetworkSegmentIndex(prefixLength, bytesPerSegment, bitsPerSegment)
	hostSegmentIndex := getHostSegmentIndex(prefixLength, bytesPerSegment, bitsPerSegment)
	return getLongCount(func(index int) uint64 {
		if (networkSegmentIndex == hostSegmentIndex) && index == networkSegmentIndex {
			segmentPrefixLength := getPrefixedSegmentPrefixLength(section.GetBitsPerSegment(), prefixLength, index)
			return getPrefixValueCount(section.GetSegment(index), segmentPrefixLength.bitCount())
		}
		return section.GetSegment(index).GetValueCount()
	},
		networkSegmentIndex+1)
}

func mult(currentResult *big.Int, newResult uint64) *big.Int {
	if currentResult == nil {
		return bigZero().SetUint64(newResult)
	} else if newResult == 1 {
		return currentResult
	}
	newBig := bigZero().SetUint64(newResult)
	return currentResult.Mul(currentResult, newBig)
}

// only called when isMultiple() is true, so segCount >= 1
func count(segmentCountProvider func(index int) uint64, segCount, safeMultiplies int, safeLimit uint64) *big.Int {
	if segCount <= 0 {
		return bigOne()
	}
	var result *big.Int
	i := 0
	for {
		curResult := segmentCountProvider(i)
		i++
		if i == segCount {
			return mult(result, curResult)
		}
		limit := i + safeMultiplies
		if segCount <= limit {
			// all multiplies are safe
			for i < segCount {
				curResult *= segmentCountProvider(i)
				i++
			}
			return mult(result, curResult)
		}
		// do the safe multiplies which cannot overflow
		for i < limit {
			curResult *= segmentCountProvider(i)
			i++
		}
		// do as many additional multiplies as current result allows
		for curResult <= safeLimit {
			curResult *= segmentCountProvider(i)
			i++
			if i == segCount {
				return mult(result, curResult)
			}
		}
		result = mult(result, curResult)
	}
}

func reverseUint8(b uint8) uint8 {
	x := b
	x = ((x & 0xaa) >> 1) | ((x & 0x55) << 1)
	x = ((x & 0xcc) >> 2) | ((x & 0x33) << 2)
	x = (x >> 4) | (x << 4)
	return x
}

func reverseUint16(b uint16) uint16 {
	x := b
	x = ((x & 0xaaaa) >> 1) | ((x & 0x5555) << 1)
	x = ((x & 0xcccc) >> 2) | ((x & 0x3333) << 2)
	x = ((x & 0xf0f0) >> 4) | ((x & 0x0f0f) << 4)
	return (x >> 8) | (x << 8)
}

func reverseUint32(i uint32) uint32 {
	x := i
	x = ((x & 0xaaaaaaaa) >> 1) | ((x & 0x55555555) << 1)
	x = ((x & 0xcccccccc) >> 2) | ((x & 0x33333333) << 2)
	x = ((x & 0xf0f0f0f0) >> 4) | ((x & 0x0f0f0f0f) << 4)
	x = ((x & 0xff00ff00) >> 8) | ((x & 0x00ff00ff) << 8)
	return (x >> 16) | (x << 16)
}

func flagsFromState(state fmt.State, verb rune) string {
	flags := "# +-0"
	vals := make([]rune, 0, len(flags)+5) // %, flags, width, '.', precision, verb
	vals = append(vals, '%')
	for i := 0; i < len(flags); i++ {
		b := flags[i]
		if state.Flag(int(b)) {
			vals = append(vals, rune(b))
		}
	}
	width, widthOK := state.Width()
	precision, precisionOK := state.Precision()
	if widthOK || precisionOK {
		var wpv string
		if widthOK && precisionOK {
			wpv = fmt.Sprintf("%d.%d%c", width, precision, verb)
		} else if widthOK {
			wpv = fmt.Sprintf("%d%c", width, verb)
		} else {
			wpv = fmt.Sprintf(".%d%c", precision, verb)
		}
		return string(vals) + wpv
	}
	vals = append(vals, verb)
	return string(vals)
}

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func min[T ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

//  TODO LATER generics: replace all uses of atomicStorePointer with atomic.Pointer, when we move up to 1.19

func atomicLoadPointer(dataLoc *unsafe.Pointer) unsafe.Pointer {
	return atomic.LoadPointer(dataLoc)
}

func atomicStorePointer(dataLoc *unsafe.Pointer, val unsafe.Pointer) {
	atomic.StorePointer(dataLoc, val)
}
