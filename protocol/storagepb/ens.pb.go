// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: storagepb/ens.proto

package storagepb

import (
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

type EnsByAddressesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses []string `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *EnsByAddressesRequest) Reset() {
	*x = EnsByAddressesRequest{}
	mi := &file_storagepb_ens_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnsByAddressesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnsByAddressesRequest) ProtoMessage() {}

func (x *EnsByAddressesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_ens_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnsByAddressesRequest.ProtoReflect.Descriptor instead.
func (*EnsByAddressesRequest) Descriptor() ([]byte, []int) {
	return file_storagepb_ens_proto_rawDescGZIP(), []int{0}
}

func (x *EnsByAddressesRequest) GetAddresses() []string {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type EnsName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *EnsName) Reset() {
	*x = EnsName{}
	mi := &file_storagepb_ens_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnsName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnsName) ProtoMessage() {}

func (x *EnsName) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_ens_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnsName.ProtoReflect.Descriptor instead.
func (*EnsName) Descriptor() ([]byte, []int) {
	return file_storagepb_ens_proto_rawDescGZIP(), []int{1}
}

func (x *EnsName) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *EnsName) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type EnsByAddressesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnsNames []*EnsName `protobuf:"bytes,1,rep,name=ens_names,json=ensNames,proto3" json:"ens_names,omitempty"`
}

func (x *EnsByAddressesResponse) Reset() {
	*x = EnsByAddressesResponse{}
	mi := &file_storagepb_ens_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnsByAddressesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnsByAddressesResponse) ProtoMessage() {}

func (x *EnsByAddressesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_ens_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnsByAddressesResponse.ProtoReflect.Descriptor instead.
func (*EnsByAddressesResponse) Descriptor() ([]byte, []int) {
	return file_storagepb_ens_proto_rawDescGZIP(), []int{2}
}

func (x *EnsByAddressesResponse) GetEnsNames() []*EnsName {
	if x != nil {
		return x.EnsNames
	}
	return nil
}

var File_storagepb_ens_proto protoreflect.FileDescriptor

var file_storagepb_ens_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2f, 0x65, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62,
	0x22, 0x35, 0x0a, 0x15, 0x45, 0x6e, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0x37, 0x0a, 0x07, 0x45, 0x6e, 0x73, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x49, 0x0a, 0x16, 0x45, 0x6e, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x09, 0x65, 0x6e,
	0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x45, 0x6e, 0x73, 0x4e, 0x61, 0x6d,
	0x65, 0x52, 0x08, 0x65, 0x6e, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x32, 0x5f, 0x0a, 0x03, 0x45,
	0x6e, 0x73, 0x12, 0x58, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x73, 0x42, 0x79, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x70, 0x62, 0x2e, 0x45, 0x6e, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x45, 0x6e, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b,
	0x2e, 0x3b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_storagepb_ens_proto_rawDescOnce sync.Once
	file_storagepb_ens_proto_rawDescData = file_storagepb_ens_proto_rawDesc
)

func file_storagepb_ens_proto_rawDescGZIP() []byte {
	file_storagepb_ens_proto_rawDescOnce.Do(func() {
		file_storagepb_ens_proto_rawDescData = protoimpl.X.CompressGZIP(file_storagepb_ens_proto_rawDescData)
	})
	return file_storagepb_ens_proto_rawDescData
}

var file_storagepb_ens_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_storagepb_ens_proto_goTypes = []any{
	(*EnsByAddressesRequest)(nil),  // 0: storagepb.EnsByAddressesRequest
	(*EnsName)(nil),                // 1: storagepb.EnsName
	(*EnsByAddressesResponse)(nil), // 2: storagepb.EnsByAddressesResponse
}
var file_storagepb_ens_proto_depIdxs = []int32{
	1, // 0: storagepb.EnsByAddressesResponse.ens_names:type_name -> storagepb.EnsName
	0, // 1: storagepb.Ens.GetEnsByAddresses:input_type -> storagepb.EnsByAddressesRequest
	2, // 2: storagepb.Ens.GetEnsByAddresses:output_type -> storagepb.EnsByAddressesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_storagepb_ens_proto_init() }
func file_storagepb_ens_proto_init() {
	if File_storagepb_ens_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_storagepb_ens_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storagepb_ens_proto_goTypes,
		DependencyIndexes: file_storagepb_ens_proto_depIdxs,
		MessageInfos:      file_storagepb_ens_proto_msgTypes,
	}.Build()
	File_storagepb_ens_proto = out.File
	file_storagepb_ens_proto_rawDesc = nil
	file_storagepb_ens_proto_goTypes = nil
	file_storagepb_ens_proto_depIdxs = nil
}
