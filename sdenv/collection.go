package sdenv

import "strings"

type Collection struct {
	keys  map[string]struct{}
	pairs [][2]string
}

func NewCollectsion() *Collection {
	return &Collection{
		keys:  map[string]struct{}{},
		pairs: [][2]string(nil),
	}
}

// Push adds pairs into collection. It keeps order of pairs.
// For equal keys it takes the last pair. It skips known keys.
// It is not thread safe.
func (c *Collection) Push(pairs [][2]string) {
	skipCount := map[string]int{}
	for _, p := range pairs {
		if _, ok := c.keys[p[0]]; ok {
			continue
		}
		skipCount[p[0]]++
	}
	for _, p := range pairs {
		if v, ok := skipCount[p[0]]; ok {
			if v > 1 {
				skipCount[p[0]] = v - 1
				continue
			}
			c.pairs = append(c.pairs, p) // add only last value for this key in sequence
		}
	}
	for x := range skipCount {
		c.keys[x] = struct{}{}
	}
}

// PushStd does the same things as Push, but accept os.Environ() formatted pairs (KEY=VALUE).
func (c *Collection) PushStd(s []string) {
	pairs := make([][2]string, len(s))
	for i, x := range s {
		a, b, _ := strings.Cut(x, "=") // consider ok somehow? check errors?
		pairs[i][0] = a
		pairs[i][1] = b
	}
	c.Push(pairs)
}

// Collection returns the collection. Do not change the resulting slice.
func (c *Collection) Collection() [][2]string {
	return c.pairs
}

// CollectionStd is a version of Collection, that returns results in os.Environ() format.
func (c *Collection) CollectionStd() []string {
	r := make([]string, len(c.pairs))
	for i, x := range c.pairs {
		r[i] = x[0] + "=" + x[1]
	}
	return r
}
