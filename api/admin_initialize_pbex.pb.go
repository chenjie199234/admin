// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.92
// 	protoc         v4.24.4
// source: api/admin_initialize.proto

package api

// return empty means pass
func (m *InitReq) Validate() (errstr string) {
	if len(m.GetPassword()) < 10 {
		return "field: password in object: init_req check value str len gte failed"
	}
	if len(m.GetPassword()) >= 32 {
		return "field: password in object: init_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *RootLoginReq) Validate() (errstr string) {
	if len(m.GetPassword()) < 10 {
		return "field: password in object: root_login_req check value str len gte failed"
	}
	if len(m.GetPassword()) >= 32 {
		return "field: password in object: root_login_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateRootPasswordReq) Validate() (errstr string) {
	if len(m.GetOldPassword()) < 10 {
		return "field: old_password in object: update_root_password_req check value str len gte failed"
	}
	if len(m.GetOldPassword()) >= 32 {
		return "field: old_password in object: update_root_password_req check value str len lt failed"
	}
	if len(m.GetNewPassword()) < 10 {
		return "field: new_password in object: update_root_password_req check value str len gte failed"
	}
	if len(m.GetNewPassword()) >= 32 {
		return "field: new_password in object: update_root_password_req check value str len lt failed"
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

// return empty means pass
func (m *UpdateProjectReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: update_project_req check len eq failed"
	}
	if len(m.GetNewProjectName()) == 0 {
		return "field: new_project_name in object: update_project_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *GetProjectIdByNameReq) Validate() (errstr string) {
	if len(m.GetProjectName()) == 0 {
		return "field: project_name in object: get_project_id_by_name_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *DeleteProjectReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: delete_project_req check len eq failed"
	}
	return ""
}
