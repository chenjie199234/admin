// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.1
// 	protoc         v3.21.1
// source: api/initialize.proto

package api

// return empty means pass
func (m *InitReq) Validate() (errstr string) {
	if len(m.GetPassword()) == 0 {
		return "field: password in object: init_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *RootLoginReq) Validate() (errstr string) {
	if len(m.GetPassword()) == 0 {
		return "field: password in object: root_login_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *RootPasswordReq) Validate() (errstr string) {
	if len(m.GetOldPassword()) == 0 {
		return "field: old_password in object: root_password_req check value str len not eq failed"
	}
	if len(m.GetNewPassword()) == 0 {
		return "field: new_password in object: root_password_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *CreateProjectReq) Validate() (errstr string) {
	if len(m.GetProjectName()) == 0 {
		return "field: project_name in object: create_project_req check value str len not eq failed"
	}
	return ""
}
