// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: storagepb/stats.proto

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

type GetTotalsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetTotalsRequest) Reset() {
	*x = GetTotalsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_stats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTotalsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTotalsRequest) ProtoMessage() {}

func (x *GetTotalsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_stats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTotalsRequest.ProtoReflect.Descriptor instead.
func (*GetTotalsRequest) Descriptor() ([]byte, []int) {
	return file_storagepb_stats_proto_rawDescGZIP(), []int{0}
}

type DaoStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total         int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	TotalVerified int64 `protobuf:"varint,2,opt,name=total_verified,json=totalVerified,proto3" json:"total_verified,omitempty"`
}

func (x *DaoStats) Reset() {
	*x = DaoStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_stats_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DaoStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DaoStats) ProtoMessage() {}

func (x *DaoStats) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_stats_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DaoStats.ProtoReflect.Descriptor instead.
func (*DaoStats) Descriptor() ([]byte, []int) {
	return file_storagepb_stats_proto_rawDescGZIP(), []int{1}
}

func (x *DaoStats) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *DaoStats) GetTotalVerified() int64 {
	if x != nil {
		return x.TotalVerified
	}
	return 0
}

type ProposalsStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *ProposalsStats) Reset() {
	*x = ProposalsStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_stats_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProposalsStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalsStats) ProtoMessage() {}

func (x *ProposalsStats) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_stats_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalsStats.ProtoReflect.Descriptor instead.
func (*ProposalsStats) Descriptor() ([]byte, []int) {
	return file_storagepb_stats_proto_rawDescGZIP(), []int{2}
}

func (x *ProposalsStats) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

type GetTotalsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dao       *DaoStats       `protobuf:"bytes,1,opt,name=dao,proto3" json:"dao,omitempty"`
	Proposals *ProposalsStats `protobuf:"bytes,2,opt,name=proposals,proto3" json:"proposals,omitempty"`
}

func (x *GetTotalsResponse) Reset() {
	*x = GetTotalsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_stats_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTotalsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTotalsResponse) ProtoMessage() {}

func (x *GetTotalsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_stats_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTotalsResponse.ProtoReflect.Descriptor instead.
func (*GetTotalsResponse) Descriptor() ([]byte, []int) {
	return file_storagepb_stats_proto_rawDescGZIP(), []int{3}
}

func (x *GetTotalsResponse) GetDao() *DaoStats {
	if x != nil {
		return x.Dao
	}
	return nil
}

func (x *GetTotalsResponse) GetProposals() *ProposalsStats {
	if x != nil {
		return x.Proposals
	}
	return nil
}

var File_storagepb_stats_proto protoreflect.FileDescriptor

var file_storagepb_stats_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2f, 0x73, 0x74, 0x61, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x70, 0x62, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x47, 0x0a, 0x08, 0x44, 0x61, 0x6f, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x22,
	0x26, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0x73, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x54, 0x6f,
	0x74, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x03,
	0x64, 0x61, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x44, 0x61, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x03,
	0x64, 0x61, 0x6f, 0x12, 0x37, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x32, 0x4f, 0x0a, 0x05,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x46, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x73, 0x12, 0x1b, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x47,
	0x65, 0x74, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a,
	0x0b, 0x2e, 0x3b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storagepb_stats_proto_rawDescOnce sync.Once
	file_storagepb_stats_proto_rawDescData = file_storagepb_stats_proto_rawDesc
)

func file_storagepb_stats_proto_rawDescGZIP() []byte {
	file_storagepb_stats_proto_rawDescOnce.Do(func() {
		file_storagepb_stats_proto_rawDescData = protoimpl.X.CompressGZIP(file_storagepb_stats_proto_rawDescData)
	})
	return file_storagepb_stats_proto_rawDescData
}

var file_storagepb_stats_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_storagepb_stats_proto_goTypes = []any{
	(*GetTotalsRequest)(nil),  // 0: storagepb.GetTotalsRequest
	(*DaoStats)(nil),          // 1: storagepb.DaoStats
	(*ProposalsStats)(nil),    // 2: storagepb.ProposalsStats
	(*GetTotalsResponse)(nil), // 3: storagepb.GetTotalsResponse
}
var file_storagepb_stats_proto_depIdxs = []int32{
	1, // 0: storagepb.GetTotalsResponse.dao:type_name -> storagepb.DaoStats
	2, // 1: storagepb.GetTotalsResponse.proposals:type_name -> storagepb.ProposalsStats
	0, // 2: storagepb.Stats.GetTotals:input_type -> storagepb.GetTotalsRequest
	3, // 3: storagepb.Stats.GetTotals:output_type -> storagepb.GetTotalsResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_storagepb_stats_proto_init() }
func file_storagepb_stats_proto_init() {
	if File_storagepb_stats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_storagepb_stats_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetTotalsRequest); i {
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
		file_storagepb_stats_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DaoStats); i {
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
		file_storagepb_stats_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ProposalsStats); i {
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
		file_storagepb_stats_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetTotalsResponse); i {
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
			RawDescriptor: file_storagepb_stats_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storagepb_stats_proto_goTypes,
		DependencyIndexes: file_storagepb_stats_proto_depIdxs,
		MessageInfos:      file_storagepb_stats_proto_msgTypes,
	}.Build()
	File_storagepb_stats_proto = out.File
	file_storagepb_stats_proto_rawDesc = nil
	file_storagepb_stats_proto_goTypes = nil
	file_storagepb_stats_proto_depIdxs = nil
}
