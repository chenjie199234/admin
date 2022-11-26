// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.1
// 	protoc         v3.21.1
// source: api/config.proto

package api

// return empty means pass
func (m *GroupsReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: groups_req check len eq failed"
	}
	return ""
}

// return empty means pass
func (m *AppsReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: apps_req check len eq failed"
	}
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: apps_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *CreateAppReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: create_app_req check len eq failed"
	}
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: create_app_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: create_app_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: create_app_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *DelAppReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: del_app_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: del_app_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: del_app_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateAppSecretReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: update_app_secret_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: update_app_secret_req check value str len gt failed"
	}
	if len(m.GetOldSecret()) >= 32 {
		return "field: old_secret in object: update_app_secret_req check value str len lt failed"
	}
	if len(m.GetNewSecret()) >= 32 {
		return "field: new_secret in object: update_app_secret_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *KeysReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: keys_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: keys_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: keys_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *DelKeyReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: del_key_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: del_key_req check value str len gt failed"
	}
	if len(m.GetKey()) <= 0 {
		return "field: key in object: del_key_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: del_key_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *GetKeyConfigReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: get_key_config_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: get_key_config_req check value str len gt failed"
	}
	if len(m.GetKey()) <= 0 {
		return "field: key in object: get_key_config_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: get_key_config_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *SetKeyConfigReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: set_key_config_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: set_key_config_req check value str len gt failed"
	}
	if len(m.GetKey()) <= 0 {
		return "field: key in object: set_key_config_req check value str len gt failed"
	}
	if len(m.GetValue()) <= 0 {
		return "field: value in object: set_key_config_req check value str len gt failed"
	}
	if m.GetValueType() != "raw" && m.GetValueType() != "json" && m.GetValueType() != "yaml" && m.GetValueType() != "toml" {
		return "field: value_type in object: set_key_config_req check value str in failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: set_key_config_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *RollbackReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: rollback_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: rollback_req check value str len gt failed"
	}
	if len(m.GetKey()) <= 0 {
		return "field: key in object: rollback_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: rollback_req check value str len lt failed"
	}
	if m.GetIndex() <= 0 {
		return "field: index in object: rollback_req check value uint gt failed"
	}
	return ""
}

// return empty means pass
func (m *WatchReq) Validate() (errstr string) {
	if len(m.GetGroupname()) <= 0 {
		return "field: groupname in object: watch_req check value str len gt failed"
	}
	if len(m.GetAppname()) <= 0 {
		return "field: appname in object: watch_req check value str len gt failed"
	}
	if len(m.GetKeys()) <= 0 {
		return "field: keys in object: watch_req check len gt failed"
	}
	return ""
}
