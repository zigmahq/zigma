package version

import (
	"strconv"
	"strings"
)

// Compare compares this version to another version. This
// returns -1, 0, or 1 if this version is smaller, equal,
// or larger than the other version, respectively.
func (v Version) Compare(o Version) int {
	vp := strings.Split(v.Number, ".")
	op := strings.Split(o.Number, ".")

	vl := len(vp)
	ol := len(op)

	l := vl
	if t := ol; t > l {
		l = t
	}
	if l > 3 {
		l = 3
	}

	for i := 0; i < l; i++ {
		var os, vs string
		if i < ol {
			os = op[i]
		}
		if i < vl {
			vs = vp[i]
		}
		on, _ := strconv.Atoi(os)
		vn, _ := strconv.Atoi(vs)
		switch {
		case on > vn:
			return -1
		case vn > on:
			return 1
		}
	}
	return 0
}

// NewerThan compares this version to another version, and returns a boolean true
// if the current version is newer than the target version
func (v Version) NewerThan(o Version) bool {
	return v.Compare(o) > 0
}

// OlderThan compares this version to another version, and returns a boolean true
// if the current version is older than the target version
func (v Version) OlderThan(o Version) bool {
	return v.Compare(o) < 0
}
