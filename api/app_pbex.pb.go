// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.77
// 	protoc         v4.22.3
// source: api/app.proto

package api

// return empty means pass
func (m *GetAppReq) Validate() (errstr string) {
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: get_app_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: get_app_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: get_app_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *GetAppInstancesReq) Validate() (errstr string) {
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: get_app_instances_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: get_app_instances_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: get_app_instances_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *CreateAppReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: create_app_req check len eq failed"
	}
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: create_app_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: create_app_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: create_app_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *DelAppReq) Validate() (errstr string) {
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: del_app_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: del_app_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: del_app_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateAppSecretReq) Validate() (errstr string) {
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: update_app_secret_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: update_app_secret_req check value str len gt failed"
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
func (m *DelKeyReq) Validate() (errstr string) {
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: del_key_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: del_key_req check value str len gt failed"
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
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: get_key_config_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: get_key_config_req check value str len gt failed"
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
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: set_key_config_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: set_key_config_req check value str len gt failed"
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
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: rollback_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: rollback_req check value str len gt failed"
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
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: watch_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: watch_req check value str len gt failed"
	}
	if len(m.GetKeys()) <= 0 {
		return "field: keys in object: watch_req check len gt failed"
	}
	return ""
}

// return empty means pass
func (m *SetProxyReq) Validate() (errstr string) {
	if len(m.GetGName()) == 0 {
		return "field: g_name in object: set_proxy_req check value str len not eq failed"
	}
	if len(m.GetAName()) == 0 {
		return "field: a_name in object: set_proxy_req check value str len not eq failed"
	}
	if len(m.GetPath()) == 0 {
		return "field: path in object: set_proxy_req check value str len not eq failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: set_proxy_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *DelProxyReq) Validate() (errstr string) {
	if len(m.GetGName()) == 0 {
		return "field: g_name in object: del_proxy_req check value str len not eq failed"
	}
	if len(m.GetAName()) == 0 {
		return "field: a_name in object: del_proxy_req check value str len not eq failed"
	}
	if len(m.GetPath()) == 0 {
		return "field: path in object: del_proxy_req check value str len not eq failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: del_proxy_req check value str len lt failed"
	}
	return ""
}

// return empty means pass
func (m *ProxyReq) Validate() (errstr string) {
	if len(m.GetPath()) == 0 {
		return "field: path in object: proxy_req check value str len not eq failed"
	}
	if len(m.GetGName()) == 0 {
		return "field: g_name in object: proxy_req check value str len not eq failed"
	}
	if len(m.GetAName()) == 0 {
		return "field: a_name in object: proxy_req check value str len not eq failed"
	}
	if len(m.GetData()) < 2 {
		return "field: data in object: proxy_req check value str len gte failed"
	}
	return ""
}
