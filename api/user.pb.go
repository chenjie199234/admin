// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: api/user.proto

//this is the proto package name,all proto in this project must use this name as the proto package name

package api

import (
	_ "github.com/chenjie199234/Corelib/pbex"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SuperAdminLoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Password string `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *SuperAdminLoginReq) Reset() {
	*x = SuperAdminLoginReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SuperAdminLoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuperAdminLoginReq) ProtoMessage() {}

func (x *SuperAdminLoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuperAdminLoginReq.ProtoReflect.Descriptor instead.
func (*SuperAdminLoginReq) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{0}
}

func (x *SuperAdminLoginReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type SuperAdminLoginResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *SuperAdminLoginResp) Reset() {
	*x = SuperAdminLoginResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SuperAdminLoginResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuperAdminLoginResp) ProtoMessage() {}

func (x *SuperAdminLoginResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuperAdminLoginResp.ProtoReflect.Descriptor instead.
func (*SuperAdminLoginResp) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{1}
}

func (x *SuperAdminLoginResp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type LoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LoginReq) Reset() {
	*x = LoginReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReq) ProtoMessage() {}

func (x *LoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReq.ProtoReflect.Descriptor instead.
func (*LoginReq) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{2}
}

type LoginResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *LoginResp) Reset() {
	*x = LoginResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResp) ProtoMessage() {}

func (x *LoginResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResp.ProtoReflect.Descriptor instead.
func (*LoginResp) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{3}
}

func (x *LoginResp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserName   string   `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Department []string `protobuf:"bytes,3,rep,name=department,proto3" json:"department,omitempty"`
	Ctime      uint32   `protobuf:"varint,4,opt,name=ctime,proto3" json:"ctime,omitempty"`
}

func (x *UserInfo) Reset() {
	*x = UserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfo) ProtoMessage() {}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfo.ProtoReflect.Descriptor instead.
func (*UserInfo) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{4}
}

func (x *UserInfo) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserInfo) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *UserInfo) GetDepartment() []string {
	if x != nil {
		return x.Department
	}
	return nil
}

func (x *UserInfo) GetCtime() uint32 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

type GetUsersReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserIds []string `protobuf:"bytes,1,rep,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
}

func (x *GetUsersReq) Reset() {
	*x = GetUsersReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersReq) ProtoMessage() {}

func (x *GetUsersReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersReq.ProtoReflect.Descriptor instead.
func (*GetUsersReq) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetUsersReq) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

type GetUsersResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*UserInfo `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"` //key userid,value username
}

func (x *GetUsersResp) Reset() {
	*x = GetUsersResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersResp) ProtoMessage() {}

func (x *GetUsersResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersResp.ProtoReflect.Descriptor instead.
func (*GetUsersResp) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{6}
}

func (x *GetUsersResp) GetUsers() []*UserInfo {
	if x != nil {
		return x.Users
	}
	return nil
}

type SearchUsersReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserName string `protobuf:"bytes,1,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
}

func (x *SearchUsersReq) Reset() {
	*x = SearchUsersReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchUsersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchUsersReq) ProtoMessage() {}

func (x *SearchUsersReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchUsersReq.ProtoReflect.Descriptor instead.
func (*SearchUsersReq) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{7}
}

func (x *SearchUsersReq) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

type SearchUsersResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*UserInfo `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"` //key userid,value username
}

func (x *SearchUsersResp) Reset() {
	*x = SearchUsersResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchUsersResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchUsersResp) ProtoMessage() {}

func (x *SearchUsersResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchUsersResp.ProtoReflect.Descriptor instead.
func (*SearchUsersResp) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{8}
}

func (x *SearchUsersResp) GetUsers() []*UserInfo {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_api_user_proto protoreflect.FileDescriptor

var file_api_user_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x1a, 0x0f, 0x70, 0x62, 0x65, 0x78, 0x2f, 0x70, 0x62,
	0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x39, 0x0a, 0x15, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65,
	0x71, 0x12, 0x20, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x04, 0xe8, 0x90, 0x4e, 0x00, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x22, 0x2e, 0x0a, 0x16, 0x73, 0x75, 0x70, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x0b, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x71,
	0x22, 0x22, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x77, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x61, 0x72,
	0x74, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x70,
	0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x30, 0x0a,
	0x0d, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x12, 0x1f,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x04, 0x90, 0x90, 0x4e, 0x00, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x22,
	0x38, 0x0a, 0x0e, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x73,
	0x70, 0x12, 0x26, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x35, 0x0a, 0x10, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x12, 0x21, 0x0a,
	0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x04, 0xe8, 0x90, 0x4e, 0x00, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x3b, 0x0a, 0x11, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x5f, 0x72, 0x65, 0x73, 0x70, 0x12, 0x26, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x32, 0xbd, 0x02,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x5a, 0x0a, 0x11, 0x73, 0x75, 0x70, 0x65, 0x72, 0x5f,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1c, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x22, 0x08, 0x8a, 0x9f, 0x49, 0x04, 0x70, 0x6f,
	0x73, 0x74, 0x12, 0x36, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x10, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x1a, 0x11, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70,
	0x22, 0x08, 0x8a, 0x9f, 0x49, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x12, 0x4b, 0x0a, 0x09, 0x67, 0x65,
	0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x14, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x1a, 0x15, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f,
	0x72, 0x65, 0x73, 0x70, 0x22, 0x11, 0x8a, 0x9f, 0x49, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x92, 0x9f,
	0x49, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x54, 0x0a, 0x0c, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x17, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x71,
	0x1a, 0x18, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x22, 0x11, 0x8a, 0x9f, 0x49, 0x04,
	0x70, 0x6f, 0x73, 0x74, 0x92, 0x9f, 0x49, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x28, 0x5a,
	0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x68, 0x65, 0x6e,
	0x6a, 0x69, 0x65, 0x31, 0x39, 0x39, 0x32, 0x33, 0x34, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_user_proto_rawDescOnce sync.Once
	file_api_user_proto_rawDescData = file_api_user_proto_rawDesc
)

func file_api_user_proto_rawDescGZIP() []byte {
	file_api_user_proto_rawDescOnce.Do(func() {
		file_api_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_user_proto_rawDescData)
	})
	return file_api_user_proto_rawDescData
}

var file_api_user_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_user_proto_goTypes = []interface{}{
	(*SuperAdminLoginReq)(nil),  // 0: admin.super_admin_login_req
	(*SuperAdminLoginResp)(nil), // 1: admin.super_admin_login_resp
	(*LoginReq)(nil),            // 2: admin.login_req
	(*LoginResp)(nil),           // 3: admin.login_resp
	(*UserInfo)(nil),            // 4: admin.user_info
	(*GetUsersReq)(nil),         // 5: admin.get_users_req
	(*GetUsersResp)(nil),        // 6: admin.get_users_resp
	(*SearchUsersReq)(nil),      // 7: admin.search_users_req
	(*SearchUsersResp)(nil),     // 8: admin.search_users_resp
}
var file_api_user_proto_depIdxs = []int32{
	4, // 0: admin.get_users_resp.users:type_name -> admin.user_info
	4, // 1: admin.search_users_resp.users:type_name -> admin.user_info
	0, // 2: admin.user.super_admin_login:input_type -> admin.super_admin_login_req
	2, // 3: admin.user.login:input_type -> admin.login_req
	5, // 4: admin.user.get_users:input_type -> admin.get_users_req
	7, // 5: admin.user.search_users:input_type -> admin.search_users_req
	1, // 6: admin.user.super_admin_login:output_type -> admin.super_admin_login_resp
	3, // 7: admin.user.login:output_type -> admin.login_resp
	6, // 8: admin.user.get_users:output_type -> admin.get_users_resp
	8, // 9: admin.user.search_users:output_type -> admin.search_users_resp
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_user_proto_init() }
func file_api_user_proto_init() {
	if File_api_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SuperAdminLoginReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SuperAdminLoginResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUsersReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUsersResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchUsersReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_user_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchUsersResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_user_proto_goTypes,
		DependencyIndexes: file_api_user_proto_depIdxs,
		MessageInfos:      file_api_user_proto_msgTypes,
	}.Build()
	File_api_user_proto = out.File
	file_api_user_proto_rawDesc = nil
	file_api_user_proto_goTypes = nil
	file_api_user_proto_depIdxs = nil
}
