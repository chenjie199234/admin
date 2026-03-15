package util

import (
	"strconv"
	"strings"

	"github.com/chenjie199234/Corelib/util/common"
)

func ParseID(id string) ([]uint32, error) {
	nodeid := make([]uint32, 0, strings.Count(id, ",")+1)
	for {
		var s string
		var ok bool
		s, id, ok = strings.Cut(id, ",")
		tmp, e := strconv.ParseUint(s, 10, 64)
		if e != nil {
			return nil, e
		}
		nodeid = append(nodeid, uint32(tmp))
		if !ok {
			break
		}
	}
	return nodeid, nil
}
func FormID(id []uint32) string {
	buf := make([]byte, 0, 256)
	for i, v := range id {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	return common.BTS(buf)
}
