// Code generated by "stringer -type=Valid"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Success-0]
	_ = x[Failure-1]
	_ = x[Ignored-2]
	_ = x[Skipped-3]
}

const _Valid_name = "SuccessFailureIgnoredSkipped"

var _Valid_index = [...]uint8{0, 7, 14, 21, 28}

func (i Valid) String() string {
	if i < 0 || i >= Valid(len(_Valid_index)-1) {
		return "Valid(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Valid_name[_Valid_index[i]:_Valid_index[i+1]]
}
