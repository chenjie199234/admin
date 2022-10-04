// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.1
// 	protoc         v3.21.1
// source: api/permission.proto

package api

// return empty means pass
func (m *GetUserPermissionReq) Validate() (errstr string) {
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: get_user_permission_req check value str len not eq failed"
	}
	if len(m.GetNodeId()) <= 1 {
		return "field: node_id in object: get_user_permission_req check len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateUserPermissionReq) Validate() (errstr string) {
	if len(m.GetUserId()) == 0 {
		return "field: user_id in object: update_user_permission_req check value str len not eq failed"
	}
	if len(m.GetNodeId()) <= 1 {
		return "field: node_id in object: update_user_permission_req check len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateRolePermissionReq) Validate() (errstr string) {
	if len(m.GetRoleName()) == 0 {
		return "field: role_name in object: update_role_permission_req check value str len not eq failed"
	}
	if len(m.GetNodeId()) <= 1 {
		return "field: node_id in object: update_role_permission_req check len gt failed"
	}
	return ""
}

// return empty means pass
func (m *AddNodeReq) Validate() (errstr string) {
	if len(m.GetPnodeId()) == 0 {
		return "field: pnode_id in object: add_node_req check len not eq failed"
	}
	if len(m.GetNodeName()) == 0 {
		return "field: node_name in object: add_node_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateNodeReq) Validate() (errstr string) {
	if len(m.GetNodeId()) <= 1 {
		return "field: node_id in object: update_node_req check len gt failed"
	}
	if len(m.GetNodeName()) == 0 {
		return "field: node_name in object: update_node_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *MoveNodeReq) Validate() (errstr string) {
	if len(m.GetNodeId()) <= 1 {
		return "field: node_id in object: move_node_req check len gt failed"
	}
	if len(m.GetPnodeId()) == 0 {
		return "field: pnode_id in object: move_node_req check len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *DelNodeReq) Validate() (errstr string) {
	if len(m.GetNodeId()) <= 1 {
		return "field: node_id in object: del_node_req check len gt failed"
	}
	return ""
}

// return empty means pass
func (m *ListRoleNodeReq) Validate() (errstr string) {
	if len(m.GetRoleName()) == 0 {
		return "field: role_name in object: list_role_node_req check value str len not eq failed"
	}
	return ""
}
