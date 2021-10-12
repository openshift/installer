package arrays

import (
    // "sort"
    "sync"
    "reflect"
)

type Set struct {
	m map[interface{}]bool
	sync.RWMutex
}

func MakeSet(i ...interface{}) *Set {
    s := &Set{
		m: map[interface{}]bool{},
	}

    for _, d := range i {
        s.m[d] = true

    }

    return s
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true

}

// func (s *Set) Remove(item int) {
	// s.Lock()
	// defer s.Unlock()
	// delete(s.m, item)
// }

func (s *Set) Has(item interface{}) bool {
    s.RLock()
    defer s.RUnlock()
    _, ok := s.m[item]
    return ok
}

func (s *Set) Len() int {
    return len(s.List())
}

func (s *Set) Clear() {
    s.Lock()
    defer s.Unlock()
    s.m = map[interface{}]bool{}
}

func (s *Set) IsEmpty() bool {
    if s.Len() == 0 {
        return true
    }
    return false
}

func (s *Set) List() []interface{} {
    s.RLock()
    defer s.RUnlock()
    list := make([]interface{}, 0)
    for item := range s.m {
        list = append(list, item)
    }
    return list
}

func (s *Set) Strings() []string {
    s.RLock()
    defer s.RUnlock()
    list := make([]string, 0)
    for item := range s.m {
        list = append(list, reflect.ValueOf(item).String())
    }
    return list
}

func (s *Set) Ints() []int64 {
    s.RLock()
    defer s.RUnlock()
    list := make([]int64, 0)
    for item := range s.m {
        list = append(list, reflect.ValueOf(item).Int())
    }
    return list
}

func (s *Set) Floats() []float64 {
    s.RLock()
    defer s.RUnlock()
    list := make([]float64, 0)
    for item := range s.m {
        list = append(list, reflect.ValueOf(item).Float())
    }
    return list
}

// func (s *Set) SortList() []interface{} {
    // s.RLock()
    // defer s.RUnlock()
    // list := make([]interface{}, 0)
    // for item := range s.m {
        // list = append(list, item)
    // }
    // sort.Ints(list)
    // return list
// }
