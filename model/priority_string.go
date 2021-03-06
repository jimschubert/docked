// Code generated by "stringer -type=Priority"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LowPriority-0]
	_ = x[MediumPriority-1]
	_ = x[HighPriority-2]
	_ = x[CriticalPriority-3]
}

const _Priority_name = "LowPriorityMediumPriorityHighPriorityCriticalPriority"

var _Priority_index = [...]uint8{0, 11, 25, 37, 53}

func (i Priority) String() string {
	if i < 0 || i >= Priority(len(_Priority_index)-1) {
		return "Priority(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Priority_name[_Priority_index[i]:_Priority_index[i+1]]
}
