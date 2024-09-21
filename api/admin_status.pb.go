// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.28.0
// source: api/admin_status.proto

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

// req can be set with pbex extentions
type Pingreq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Pingreq) Reset() {
	*x = Pingreq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_admin_status_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pingreq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pingreq) ProtoMessage() {}

func (x *Pingreq) ProtoReflect() protoreflect.Message {
	mi := &file_api_admin_status_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pingreq.ProtoReflect.Descriptor instead.
func (*Pingreq) Descriptor() ([]byte, []int) {
	return file_api_admin_status_proto_rawDescGZIP(), []int{0}
}

func (x *Pingreq) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// resp's pbex extentions will be ignore
type Pingresp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientTimestamp int64   `protobuf:"varint,1,opt,name=client_timestamp,json=clientTimestamp,proto3" json:"client_timestamp,omitempty"`
	ServerTimestamp int64   `protobuf:"varint,2,opt,name=server_timestamp,json=serverTimestamp,proto3" json:"server_timestamp,omitempty"`
	TotalMem        uint64  `protobuf:"varint,3,opt,name=total_mem,json=totalMem,proto3" json:"total_mem,omitempty"`
	CurMemUsage     uint64  `protobuf:"varint,4,opt,name=cur_mem_usage,json=curMemUsage,proto3" json:"cur_mem_usage,omitempty"`
	MaxMemUsage     uint64  `protobuf:"varint,5,opt,name=max_mem_usage,json=maxMemUsage,proto3" json:"max_mem_usage,omitempty"`
	CpuNum          float64 `protobuf:"fixed64,6,opt,name=cpu_num,json=cpuNum,proto3" json:"cpu_num,omitempty"`
	CurCpuUsage     float64 `protobuf:"fixed64,7,opt,name=cur_cpu_usage,json=curCpuUsage,proto3" json:"cur_cpu_usage,omitempty"`
	AvgCpuUsage     float64 `protobuf:"fixed64,8,opt,name=avg_cpu_usage,json=avgCpuUsage,proto3" json:"avg_cpu_usage,omitempty"`
	MaxCpuUsage     float64 `protobuf:"fixed64,9,opt,name=max_cpu_usage,json=maxCpuUsage,proto3" json:"max_cpu_usage,omitempty"`
	Host            string  `protobuf:"bytes,10,opt,name=host,proto3" json:"host,omitempty"`
	Ip              string  `protobuf:"bytes,11,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *Pingresp) Reset() {
	*x = Pingresp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_admin_status_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pingresp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pingresp) ProtoMessage() {}

func (x *Pingresp) ProtoReflect() protoreflect.Message {
	mi := &file_api_admin_status_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pingresp.ProtoReflect.Descriptor instead.
func (*Pingresp) Descriptor() ([]byte, []int) {
	return file_api_admin_status_proto_rawDescGZIP(), []int{1}
}

func (x *Pingresp) GetClientTimestamp() int64 {
	if x != nil {
		return x.ClientTimestamp
	}
	return 0
}

func (x *Pingresp) GetServerTimestamp() int64 {
	if x != nil {
		return x.ServerTimestamp
	}
	return 0
}

func (x *Pingresp) GetTotalMem() uint64 {
	if x != nil {
		return x.TotalMem
	}
	return 0
}

func (x *Pingresp) GetCurMemUsage() uint64 {
	if x != nil {
		return x.CurMemUsage
	}
	return 0
}

func (x *Pingresp) GetMaxMemUsage() uint64 {
	if x != nil {
		return x.MaxMemUsage
	}
	return 0
}

func (x *Pingresp) GetCpuNum() float64 {
	if x != nil {
		return x.CpuNum
	}
	return 0
}

func (x *Pingresp) GetCurCpuUsage() float64 {
	if x != nil {
		return x.CurCpuUsage
	}
	return 0
}

func (x *Pingresp) GetAvgCpuUsage() float64 {
	if x != nil {
		return x.AvgCpuUsage
	}
	return 0
}

func (x *Pingresp) GetMaxCpuUsage() float64 {
	if x != nil {
		return x.MaxCpuUsage
	}
	return 0
}

func (x *Pingresp) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Pingresp) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

var File_api_admin_status_proto protoreflect.FileDescriptor

var file_api_admin_status_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x1a,
	0x0f, 0x70, 0x62, 0x65, 0x78, 0x2f, 0x70, 0x62, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2d, 0x0a, 0x07, 0x70, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x71, 0x12, 0x22, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x04,
	0x80, 0x92, 0x4e, 0x00, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22,
	0xee, 0x02, 0x0a, 0x08, 0x70, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x73, 0x70, 0x12, 0x29, 0x0a, 0x10,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x29, 0x0a, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x6d, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4d, 0x65, 0x6d, 0x12,
	0x22, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x5f, 0x6d, 0x65, 0x6d, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x4d, 0x65, 0x6d, 0x55, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x0d, 0x6d, 0x61, 0x78, 0x5f, 0x6d, 0x65, 0x6d, 0x5f, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x6d, 0x61, 0x78, 0x4d,
	0x65, 0x6d, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x70, 0x75, 0x5f, 0x6e,
	0x75, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x63, 0x70, 0x75, 0x4e, 0x75, 0x6d,
	0x12, 0x22, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x5f, 0x63, 0x70, 0x75, 0x5f, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x43, 0x70, 0x75, 0x55,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x0d, 0x61, 0x76, 0x67, 0x5f, 0x63, 0x70, 0x75, 0x5f,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x61, 0x76, 0x67,
	0x43, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x0d, 0x6d, 0x61, 0x78, 0x5f,
	0x63, 0x70, 0x75, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0b, 0x6d, 0x61, 0x78, 0x43, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x6f, 0x73, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70,
	0x32, 0x4a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x40, 0x0a, 0x04, 0x70, 0x69,
	0x6e, 0x67, 0x12, 0x0e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x72,
	0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x72,
	0x65, 0x73, 0x70, 0x22, 0x17, 0x8a, 0x9f, 0x49, 0x03, 0x67, 0x65, 0x74, 0x8a, 0x9f, 0x49, 0x04,
	0x63, 0x72, 0x70, 0x63, 0x8a, 0x9f, 0x49, 0x04, 0x67, 0x72, 0x70, 0x63, 0x42, 0x28, 0x5a, 0x26,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x68, 0x65, 0x6e, 0x6a,
	0x69, 0x65, 0x31, 0x39, 0x39, 0x32, 0x33, 0x34, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x61,
	0x70, 0x69, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_admin_status_proto_rawDescOnce sync.Once
	file_api_admin_status_proto_rawDescData = file_api_admin_status_proto_rawDesc
)

func file_api_admin_status_proto_rawDescGZIP() []byte {
	file_api_admin_status_proto_rawDescOnce.Do(func() {
		file_api_admin_status_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_admin_status_proto_rawDescData)
	})
	return file_api_admin_status_proto_rawDescData
}

var file_api_admin_status_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_admin_status_proto_goTypes = []interface{}{
	(*Pingreq)(nil),  // 0: admin.pingreq
	(*Pingresp)(nil), // 1: admin.pingresp
}
var file_api_admin_status_proto_depIdxs = []int32{
	0, // 0: admin.status.ping:input_type -> admin.pingreq
	1, // 1: admin.status.ping:output_type -> admin.pingresp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_admin_status_proto_init() }
func file_api_admin_status_proto_init() {
	if File_api_admin_status_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_admin_status_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pingreq); i {
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
		file_api_admin_status_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pingresp); i {
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
			RawDescriptor: file_api_admin_status_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_admin_status_proto_goTypes,
		DependencyIndexes: file_api_admin_status_proto_depIdxs,
		MessageInfos:      file_api_admin_status_proto_msgTypes,
	}.Build()
	File_api_admin_status_proto = out.File
	file_api_admin_status_proto_rawDesc = nil
	file_api_admin_status_proto_goTypes = nil
	file_api_admin_status_proto_depIdxs = nil
}
