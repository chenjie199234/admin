package sdk

import (
	"github.com/chenjie199234/Corelib/util/name"
)

func Init(selfp, selfg, selfa string) error {
	if name.GetSelfFullName() == "" {
		if e := name.SetSelfFullName(selfp, selfg, selfa); e != nil {
			if name.GetSelfFullName() == "" {
				return e
			}
		}
	}
	return nil
}
