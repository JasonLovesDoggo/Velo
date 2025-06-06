// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: velo.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DeployRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ServiceName   string                 `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Image         string                 `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	Env           map[string]string      `protobuf:"bytes,3,rep,name=env,proto3" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeployRequest) Reset() {
	*x = DeployRequest{}
	mi := &file_velo_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeployRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployRequest) ProtoMessage() {}

func (x *DeployRequest) ProtoReflect() protoreflect.Message {
	mi := &file_velo_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployRequest.ProtoReflect.Descriptor instead.
func (*DeployRequest) Descriptor() ([]byte, []int) {
	return file_velo_proto_rawDescGZIP(), []int{0}
}

func (x *DeployRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *DeployRequest) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *DeployRequest) GetEnv() map[string]string {
	if x != nil {
		return x.Env
	}
	return nil
}

type DeployResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DeploymentId  string                 `protobuf:"bytes,1,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeployResponse) Reset() {
	*x = DeployResponse{}
	mi := &file_velo_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeployResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployResponse) ProtoMessage() {}

func (x *DeployResponse) ProtoReflect() protoreflect.Message {
	mi := &file_velo_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployResponse.ProtoReflect.Descriptor instead.
func (*DeployResponse) Descriptor() ([]byte, []int) {
	return file_velo_proto_rawDescGZIP(), []int{1}
}

func (x *DeployResponse) GetDeploymentId() string {
	if x != nil {
		return x.DeploymentId
	}
	return ""
}

func (x *DeployResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type RollbackRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DeploymentId  string                 `protobuf:"bytes,1,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RollbackRequest) Reset() {
	*x = RollbackRequest{}
	mi := &file_velo_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RollbackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RollbackRequest) ProtoMessage() {}

func (x *RollbackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_velo_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RollbackRequest.ProtoReflect.Descriptor instead.
func (*RollbackRequest) Descriptor() ([]byte, []int) {
	return file_velo_proto_rawDescGZIP(), []int{2}
}

func (x *RollbackRequest) GetDeploymentId() string {
	if x != nil {
		return x.DeploymentId
	}
	return ""
}

type GenericResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Success       bool                   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenericResponse) Reset() {
	*x = GenericResponse{}
	mi := &file_velo_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenericResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenericResponse) ProtoMessage() {}

func (x *GenericResponse) ProtoReflect() protoreflect.Message {
	mi := &file_velo_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenericResponse.ProtoReflect.Descriptor instead.
func (*GenericResponse) Descriptor() ([]byte, []int) {
	return file_velo_proto_rawDescGZIP(), []int{3}
}

func (x *GenericResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *GenericResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type StatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DeploymentId  string                 `protobuf:"bytes,1,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatusRequest) Reset() {
	*x = StatusRequest{}
	mi := &file_velo_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRequest) ProtoMessage() {}

func (x *StatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_velo_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRequest.ProtoReflect.Descriptor instead.
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return file_velo_proto_rawDescGZIP(), []int{4}
}

func (x *StatusRequest) GetDeploymentId() string {
	if x != nil {
		return x.DeploymentId
	}
	return ""
}

type StatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Logs          string                 `protobuf:"bytes,2,opt,name=logs,proto3" json:"logs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	mi := &file_velo_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_velo_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_velo_proto_rawDescGZIP(), []int{5}
}

func (x *StatusResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StatusResponse) GetLogs() string {
	if x != nil {
		return x.Logs
	}
	return ""
}

var File_velo_proto protoreflect.FileDescriptor

const file_velo_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"velo.proto\x12\x04velo\"\xb0\x01\n" +
	"\rDeployRequest\x12!\n" +
	"\fservice_name\x18\x01 \x01(\tR\vserviceName\x12\x14\n" +
	"\x05image\x18\x02 \x01(\tR\x05image\x12.\n" +
	"\x03env\x18\x03 \x03(\v2\x1c.velo.DeployRequest.EnvEntryR\x03env\x1a6\n" +
	"\bEnvEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"M\n" +
	"\x0eDeployResponse\x12#\n" +
	"\rdeployment_id\x18\x01 \x01(\tR\fdeploymentId\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\"6\n" +
	"\x0fRollbackRequest\x12#\n" +
	"\rdeployment_id\x18\x01 \x01(\tR\fdeploymentId\"E\n" +
	"\x0fGenericResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\x12\x18\n" +
	"\asuccess\x18\x02 \x01(\bR\asuccess\"4\n" +
	"\rStatusRequest\x12#\n" +
	"\rdeployment_id\x18\x01 \x01(\tR\fdeploymentId\"<\n" +
	"\x0eStatusResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12\x12\n" +
	"\x04logs\x18\x02 \x01(\tR\x04logs2\xba\x01\n" +
	"\x11DeploymentService\x123\n" +
	"\x06Deploy\x12\x13.velo.DeployRequest\x1a\x14.velo.DeployResponse\x128\n" +
	"\bRollback\x12\x15.velo.RollbackRequest\x1a\x15.velo.GenericResponse\x126\n" +
	"\tGetStatus\x12\x13.velo.StatusRequest\x1a\x14.velo.StatusResponseB\x10Z\x0evelo/api/protob\x06proto3"

var (
	file_velo_proto_rawDescOnce sync.Once
	file_velo_proto_rawDescData []byte
)

func file_velo_proto_rawDescGZIP() []byte {
	file_velo_proto_rawDescOnce.Do(func() {
		file_velo_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_velo_proto_rawDesc), len(file_velo_proto_rawDesc)))
	})
	return file_velo_proto_rawDescData
}

var file_velo_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_velo_proto_goTypes = []any{
	(*DeployRequest)(nil),   // 0: velo.DeployRequest
	(*DeployResponse)(nil),  // 1: velo.DeployResponse
	(*RollbackRequest)(nil), // 2: velo.RollbackRequest
	(*GenericResponse)(nil), // 3: velo.GenericResponse
	(*StatusRequest)(nil),   // 4: velo.StatusRequest
	(*StatusResponse)(nil),  // 5: velo.StatusResponse
	nil,                     // 6: velo.DeployRequest.EnvEntry
}
var file_velo_proto_depIdxs = []int32{
	6, // 0: velo.DeployRequest.env:type_name -> velo.DeployRequest.EnvEntry
	0, // 1: velo.DeploymentService.Deploy:input_type -> velo.DeployRequest
	2, // 2: velo.DeploymentService.Rollback:input_type -> velo.RollbackRequest
	4, // 3: velo.DeploymentService.GetStatus:input_type -> velo.StatusRequest
	1, // 4: velo.DeploymentService.Deploy:output_type -> velo.DeployResponse
	3, // 5: velo.DeploymentService.Rollback:output_type -> velo.GenericResponse
	5, // 6: velo.DeploymentService.GetStatus:output_type -> velo.StatusResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_velo_proto_init() }
func file_velo_proto_init() {
	if File_velo_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_velo_proto_rawDesc), len(file_velo_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_velo_proto_goTypes,
		DependencyIndexes: file_velo_proto_depIdxs,
		MessageInfos:      file_velo_proto_msgTypes,
	}.Build()
	File_velo_proto = out.File
	file_velo_proto_goTypes = nil
	file_velo_proto_depIdxs = nil
}
