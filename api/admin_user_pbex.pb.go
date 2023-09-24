// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.81
// 	protoc         v4.24.1
// source: api/admin_user.proto

package api

// return empty means pass
func (m *InviteProjectReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: invite_project_req check len eq failed"
	}
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: invite_project_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *KickProjectReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: kick_project_req check len eq failed"
	}
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: kick_project_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *SearchUsersReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: search_users_req check len eq failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateUserReq) Validate() (errstr string) {
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: update_user_req check value str len not eq failed"
	}
	if len(m.GetNewUserName()) == 0 {
		return "field: new_user_name in object: update_user_req check value str len not eq failed"
	}
	if len(m.GetNewDepartment()) == 0 {
		return "field: new_department in object: update_user_req check len not eq failed"
	}
	for _, v := range m.GetNewDepartment() {
		if len(v) == 0 {
			return "field: new_department in object: update_user_req check value str len not eq failed"
		}
	}
	return ""
}

// return empty means pass
func (m *CreateRoleReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: create_role_req check len eq failed"
	}
	if len(m.GetRoleName()) == 0 {
		return "field: role_name in object: create_role_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *SearchRolesReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: search_roles_req check len eq failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateRoleReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: update_role_req check len eq failed"
	}
	if len(m.GetRoleName()) == 0 {
		return "field: role_name in object: update_role_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *DelRolesReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: del_roles_req check len eq failed"
	}
	if len(m.GetRoleNames()) == 0 {
		return "field: role_names in object: del_roles_req check len not eq failed"
	}
	for _, v := range m.GetRoleNames() {
		if len(v) == 0 {
			return "field: role_names in object: del_roles_req check value str len not eq failed"
		}
	}
	return ""
}

// return empty means pass
func (m *AddUserRoleReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: add_user_role_req check len eq failed"
	}
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: add_user_role_req check value str len not eq failed"
	}
	if len(m.GetRoleName()) == 0 {
		return "field: role_name in object: add_user_role_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *DelUserRoleReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: del_user_role_req check len eq failed"
	}
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: del_user_role_req check value str len not eq failed"
	}
	if len(m.GetRoleName()) == 0 {
		return "field: role_name in object: del_user_role_req check value str len not eq failed"
	}
	return ""
}