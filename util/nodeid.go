package util

import (
	"strconv"
	"strings"
)

func ParseNodeIDstr(nodeidstr string) ([]uint32, error) {
	nodeid := make([]uint32, 0, strings.Count(nodeidstr, ",")+1)
	for {
		var s string
		index := strings.Index(nodeidstr, ",")
		if index == -1 {
			s = nodeidstr
		} else {
			s = nodeidstr[:index]
			nodeidstr = nodeidstr[index+1:]
		}
		tmp, e := strconv.ParseUint(s, 10, 64)
		if e != nil {
			return nil, e
		}
		nodeid = append(nodeid, uint32(tmp))
		if index == -1 {
			break
		}
	}
	return nodeid, nil
}
