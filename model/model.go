package model

import (
	"os"

	"github.com/chenjie199234/Corelib/util/name"
)

// Warning!!!!!!!!!!!
// This file is readonly!
// Don't modify this file!

const pkg = "github.com/chenjie199234/admin"
const Name = "admin"

var Group = os.Getenv("GROUP")
var Project = os.Getenv("PROJECT")

func init() {
	if Group == "" || Group == "<GROUP>" {
		panic("missing GROUP env")
	}
	if name.SingleCheck(Group, false) != nil {
		panic("env:GROUP format wrong")
	}
	if Project == "" || Project == "<PROJECT_ID>" {
		panic("missing PROJECT env")
	}
	if name.SingleCheck(Project, false) != nil {
		panic("env:PROJECT format wrong")
	}
}
