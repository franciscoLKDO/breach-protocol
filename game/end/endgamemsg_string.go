// Code generated by "stringer -type=EndGameMsg -linecomment"; DO NOT EDIT.

package end

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Quit-0]
	_ = x[Restart-1]
}

const _EndGameMsg_name = "QuitRestart"

var _EndGameMsg_index = [...]uint8{0, 4, 11}

func (i EndGameMsg) String() string {
	if i < 0 || i >= EndGameMsg(len(_EndGameMsg_index)-1) {
		return "EndGameMsg(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EndGameMsg_name[_EndGameMsg_index[i]:_EndGameMsg_index[i+1]]
}
