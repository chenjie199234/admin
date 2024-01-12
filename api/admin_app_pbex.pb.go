// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.95
// 	protoc         v4.25.1
// source: api/admin_app.proto

package api

// return empty means pass
func (m *GetAppReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: get_app_req check len eq failed"
	}
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
func (m *SetAppReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: set_app_req check len eq failed"
	}
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: set_app_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: set_app_req check value str len gt failed"
	}
	if len(m.GetSecret()) >= 32 {
		return "field: secret in object: set_app_req check value str len lt failed"
	}
	if m.GetDiscoverMode() != "kubernetes" && m.GetDiscoverMode() != "dns" && m.GetDiscoverMode() != "static" {
		return "field: discover_mode in object: set_app_req check value str in failed"
	}
	if m.GetCrpcPort() <= 0 {
		return "field: crpc_port in object: set_app_req check value uint gt failed"
	}
	if m.GetCrpcPort() >= 65536 {
		return "field: crpc_port in object: set_app_req check value uint lt failed"
	}
	if m.GetCgrpcPort() <= 0 {
		return "field: cgrpc_port in object: set_app_req check value uint gt failed"
	}
	if m.GetCgrpcPort() >= 65536 {
		return "field: cgrpc_port in object: set_app_req check value uint lt failed"
	}
	if m.GetWebPort() <= 0 {
		return "field: web_port in object: set_app_req check value uint gt failed"
	}
	if m.GetWebPort() >= 65536 {
		return "field: web_port in object: set_app_req check value uint lt failed"
	}
	return ""
}

// return empty means pass
func (m *DelAppReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: del_app_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: update_app_secret_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: del_key_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: get_key_config_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: set_key_config_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: rollback_req check len eq failed"
	}
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
func (m *WatchConfigReq) Validate() (errstr string) {
	if len(m.GetProjectName()) <= 0 {
		return "field: project_name in object: watch_config_req check value str len gt failed"
	}
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: watch_config_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: watch_config_req check value str len gt failed"
	}
	if len(m.GetKeys()) <= 0 {
		return "field: keys in object: watch_config_req check len gt failed"
	}
	return ""
}

// return empty means pass
func (m *WatchDiscoverReq) Validate() (errstr string) {
	if len(m.GetProjectName()) <= 0 {
		return "field: project_name in object: watch_discover_req check value str len gt failed"
	}
	if len(m.GetGName()) <= 0 {
		return "field: g_name in object: watch_discover_req check value str len gt failed"
	}
	if len(m.GetAName()) <= 0 {
		return "field: a_name in object: watch_discover_req check value str len gt failed"
	}
	if m.GetCurDiscoverMode() != "kubernetes" && m.GetCurDiscoverMode() != "dns" && m.GetCurDiscoverMode() != "static" && m.GetCurDiscoverMode() != "" {
		return "field: cur_discover_mode in object: watch_discover_req check value str in failed"
	}
	if m.GetCurCrpcPort() >= 65536 {
		return "field: cur_crpc_port in object: watch_discover_req check value uint lt failed"
	}
	if m.GetCurCgrpcPort() >= 65536 {
		return "field: cur_cgrpc_port in object: watch_discover_req check value uint lt failed"
	}
	if m.GetCurWebPort() >= 65536 {
		return "field: cur_web_port in object: watch_discover_req check value uint lt failed"
	}
	return ""
}

// return empty means pass
func (m *WatchDiscoverResp) Validate() (errstr string) {
	if m.GetCrpcPort() >= 65536 {
		return "field: crpc_port in object: watch_discover_resp check value uint lt failed"
	}
	if m.GetCgrpcPort() >= 65536 {
		return "field: cgrpc_port in object: watch_discover_resp check value uint lt failed"
	}
	if m.GetWebPort() >= 65536 {
		return "field: web_port in object: watch_discover_resp check value uint lt failed"
	}
	return ""
}

// return empty means pass
func (m *GetInstancesReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: get_instances_req check len eq failed"
	}
	if len(m.GetGName()) == 0 {
		return "field: g_name in object: get_instances_req check value str len not eq failed"
	}
	if len(m.GetAName()) == 0 {
		return "field: a_name in object: get_instances_req check value str len not eq failed"
	}
	return ""
}

// return empty means pass
func (m *GetInstanceInfoReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: get_instance_info_req check len eq failed"
	}
	if len(m.GetGName()) == 0 {
		return "field: g_name in object: get_instance_info_req check value str len not eq failed"
	}
	if len(m.GetAName()) == 0 {
		return "field: a_name in object: get_instance_info_req check value str len not eq failed"
	}
	if len(m.GetAddr()) <= 0 {
		return "field: addr in object: get_instance_info_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *SetProxyReq) Validate() (errstr string) {
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: set_proxy_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: del_proxy_req check len eq failed"
	}
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
	if len(m.GetProjectId()) != 2 {
		return "field: project_id in object: proxy_req check len eq failed"
	}
	if len(m.GetGName()) == 0 {
		return "field: g_name in object: proxy_req check value str len not eq failed"
	}
	if len(m.GetAName()) == 0 {
		return "field: a_name in object: proxy_req check value str len not eq failed"
	}
	if len(m.GetPath()) == 0 {
		return "field: path in object: proxy_req check value str len not eq failed"
	}
	if len(m.GetData()) < 2 {
		return "field: data in object: proxy_req check value str len gte failed"
	}
	return ""
}
