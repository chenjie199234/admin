// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.77
// 	protoc         v4.24.1
// source: api/adminstatus.proto

package api

// return empty means pass
func (m *Pingreq) Validate() (errstr string) {
	if m.GetTimestamp() <= 0 {
		return "field: timestamp in object: pingreq check value int gt failed"
	}
	return ""
}
