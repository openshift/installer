package arrays

import (
    "reflect"
)

// Contains Returns the index position of the val in array
func Contains(array interface{}, val interface{}) (index int) {
    index = -1
    switch reflect.TypeOf(array).Kind() {
        case reflect.Slice: {
            s := reflect.ValueOf(array)
            for i := 0; i < s.Len(); i++ {
                if reflect.DeepEqual(val, s.Index(i).Interface()) {
                    index = i
                    return
                }
            }
        }
    }
    return
}

// ContainsString Returns the index position of the string val in array
func ContainsString(array []string, val string) (index int) {
    index = -1
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            index = i
            return
        }
    }
    return
}

// ContainsInt Returns the index position of the int64 val in array
func ContainsInt(array []int64, val int64) (index int) {
    index = -1
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            index = i
            return
        }
    }
    return
}

// ContainsUint Returns the index position of the uint64 val in array
func ContainsUint(array []uint64, val uint64) (index int) {
    index = -1
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            index = i
            return
        }
    }
    return
}

// ContainsBool Returns the index position of the bool val in array
func ContainsBool(array []bool, val bool) (index int) {
    index = -1
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            index = i
            return
        }
    }
    return
}

// ContainsFloat Returns the index position of the float64 val in array
func ContainsFloat(array []float64, val float64) (index int) {
    index = -1
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            index = i
            return
        }
    }
    return
}

// ContainsComplex Returns the index position of the complex128 val in array
func ContainsComplex(array []complex128, val complex128) (index int) {
    index = -1
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            index = i
            return
        }
    }
    return
}

