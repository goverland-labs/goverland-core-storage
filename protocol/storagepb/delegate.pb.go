// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: storagepb/delegate.proto

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

type GetDelegatesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DaoId         string   `protobuf:"bytes,1,opt,name=dao_id,json=daoId,proto3" json:"dao_id,omitempty"`
	QueryAccounts []string `protobuf:"bytes,2,rep,name=query_accounts,json=queryAccounts,proto3" json:"query_accounts,omitempty"`
	Sort          *string  `protobuf:"bytes,3,opt,name=sort,proto3,oneof" json:"sort,omitempty"`
	Limit         int32    `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset        int32    `protobuf:"varint,5,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetDelegatesRequest) Reset() {
	*x = GetDelegatesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_delegate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDelegatesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDelegatesRequest) ProtoMessage() {}

func (x *GetDelegatesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_delegate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDelegatesRequest.ProtoReflect.Descriptor instead.
func (*GetDelegatesRequest) Descriptor() ([]byte, []int) {
	return file_storagepb_delegate_proto_rawDescGZIP(), []int{0}
}

func (x *GetDelegatesRequest) GetDaoId() string {
	if x != nil {
		return x.DaoId
	}
	return ""
}

func (x *GetDelegatesRequest) GetQueryAccounts() []string {
	if x != nil {
		return x.QueryAccounts
	}
	return nil
}

func (x *GetDelegatesRequest) GetSort() string {
	if x != nil && x.Sort != nil {
		return *x.Sort
	}
	return ""
}

func (x *GetDelegatesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetDelegatesRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type GetDelegatesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Delegates []*DelegateEntry `protobuf:"bytes,1,rep,name=delegates,proto3" json:"delegates,omitempty"`
}

func (x *GetDelegatesResponse) Reset() {
	*x = GetDelegatesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_delegate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDelegatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDelegatesResponse) ProtoMessage() {}

func (x *GetDelegatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_delegate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDelegatesResponse.ProtoReflect.Descriptor instead.
func (*GetDelegatesResponse) Descriptor() ([]byte, []int) {
	return file_storagepb_delegate_proto_rawDescGZIP(), []int{1}
}

func (x *GetDelegatesResponse) GetDelegates() []*DelegateEntry {
	if x != nil {
		return x.Delegates
	}
	return nil
}

type DelegateEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address                  string  `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	EnsName                  string  `protobuf:"bytes,2,opt,name=ens_name,json=ensName,proto3" json:"ens_name,omitempty"`
	DelegatorCount           int32   `protobuf:"varint,3,opt,name=delegator_count,json=delegatorCount,proto3" json:"delegator_count,omitempty"`
	PercentOfDelegators      float64 `protobuf:"fixed64,4,opt,name=percent_of_delegators,json=percentOfDelegators,proto3" json:"percent_of_delegators,omitempty"` // in basis points
	VotingPower              float64 `protobuf:"fixed64,5,opt,name=votingPower,proto3" json:"votingPower,omitempty"`
	PercentOfVotingPower     float64 `protobuf:"fixed64,6,opt,name=percent_of_voting_power,json=percentOfVotingPower,proto3" json:"percent_of_voting_power,omitempty"` // in basis points
	About                    string  `protobuf:"bytes,7,opt,name=about,proto3" json:"about,omitempty"`
	Statement                string  `protobuf:"bytes,8,opt,name=statement,proto3" json:"statement,omitempty"`
	UserDelegatedVotingPower float64 `protobuf:"fixed64,9,opt,name=user_delegated_voting_power,json=userDelegatedVotingPower,proto3" json:"user_delegated_voting_power,omitempty"`
	VotesCount               int32   `protobuf:"varint,10,opt,name=votes_count,json=votesCount,proto3" json:"votes_count,omitempty"`
	ProposalsCount           int32   `protobuf:"varint,11,opt,name=proposals_count,json=proposalsCount,proto3" json:"proposals_count,omitempty"`
	CreateProposalsCount     int32   `protobuf:"varint,12,opt,name=create_proposals_count,json=createProposalsCount,proto3" json:"create_proposals_count,omitempty"`
}

func (x *DelegateEntry) Reset() {
	*x = DelegateEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storagepb_delegate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelegateEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelegateEntry) ProtoMessage() {}

func (x *DelegateEntry) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_delegate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelegateEntry.ProtoReflect.Descriptor instead.
func (*DelegateEntry) Descriptor() ([]byte, []int) {
	return file_storagepb_delegate_proto_rawDescGZIP(), []int{2}
}

func (x *DelegateEntry) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *DelegateEntry) GetEnsName() string {
	if x != nil {
		return x.EnsName
	}
	return ""
}

func (x *DelegateEntry) GetDelegatorCount() int32 {
	if x != nil {
		return x.DelegatorCount
	}
	return 0
}

func (x *DelegateEntry) GetPercentOfDelegators() float64 {
	if x != nil {
		return x.PercentOfDelegators
	}
	return 0
}

func (x *DelegateEntry) GetVotingPower() float64 {
	if x != nil {
		return x.VotingPower
	}
	return 0
}

func (x *DelegateEntry) GetPercentOfVotingPower() float64 {
	if x != nil {
		return x.PercentOfVotingPower
	}
	return 0
}

func (x *DelegateEntry) GetAbout() string {
	if x != nil {
		return x.About
	}
	return ""
}

func (x *DelegateEntry) GetStatement() string {
	if x != nil {
		return x.Statement
	}
	return ""
}

func (x *DelegateEntry) GetUserDelegatedVotingPower() float64 {
	if x != nil {
		return x.UserDelegatedVotingPower
	}
	return 0
}

func (x *DelegateEntry) GetVotesCount() int32 {
	if x != nil {
		return x.VotesCount
	}
	return 0
}

func (x *DelegateEntry) GetProposalsCount() int32 {
	if x != nil {
		return x.ProposalsCount
	}
	return 0
}

func (x *DelegateEntry) GetCreateProposalsCount() int32 {
	if x != nil {
		return x.CreateProposalsCount
	}
	return 0
}

var File_storagepb_delegate_proto protoreflect.FileDescriptor

var file_storagepb_delegate_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2f, 0x64, 0x65, 0x6c, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x70, 0x62, 0x22, 0xa3, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x44, 0x65, 0x6c,
	0x65, 0x67, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a,
	0x06, 0x64, 0x61, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64,
	0x61, 0x6f, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x17, 0x0a, 0x04, 0x73,
	0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x73, 0x6f, 0x72,
	0x74, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x73, 0x6f, 0x72, 0x74, 0x22, 0x4e, 0x0a, 0x14, 0x47,
	0x65, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73, 0x22, 0xed, 0x03, 0x0a, 0x0d,
	0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x73, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x73, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x64, 0x65, 0x6c,
	0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x15, 0x70,
	0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x13, 0x70, 0x65, 0x72, 0x63,
	0x65, 0x6e, 0x74, 0x4f, 0x66, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x12,
	0x20, 0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65,
	0x72, 0x12, 0x35, 0x0a, 0x17, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x66, 0x5f,
	0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x14, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x4f, 0x66, 0x56, 0x6f, 0x74,
	0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x62, 0x6f, 0x75,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x3d, 0x0a, 0x1b,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x76,
	0x6f, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x18, 0x75, 0x73, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64,
	0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x76,
	0x6f, 0x74, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x27, 0x0a, 0x0f,
	0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x16, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f,
	0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x14, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x5b, 0x0a, 0x08, 0x44,
	0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x44, 0x65,
	0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73, 0x12, 0x1e, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storagepb_delegate_proto_rawDescOnce sync.Once
	file_storagepb_delegate_proto_rawDescData = file_storagepb_delegate_proto_rawDesc
)

func file_storagepb_delegate_proto_rawDescGZIP() []byte {
	file_storagepb_delegate_proto_rawDescOnce.Do(func() {
		file_storagepb_delegate_proto_rawDescData = protoimpl.X.CompressGZIP(file_storagepb_delegate_proto_rawDescData)
	})
	return file_storagepb_delegate_proto_rawDescData
}

var file_storagepb_delegate_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_storagepb_delegate_proto_goTypes = []any{
	(*GetDelegatesRequest)(nil),  // 0: storagepb.GetDelegatesRequest
	(*GetDelegatesResponse)(nil), // 1: storagepb.GetDelegatesResponse
	(*DelegateEntry)(nil),        // 2: storagepb.DelegateEntry
}
var file_storagepb_delegate_proto_depIdxs = []int32{
	2, // 0: storagepb.GetDelegatesResponse.delegates:type_name -> storagepb.DelegateEntry
	0, // 1: storagepb.Delegate.GetDelegates:input_type -> storagepb.GetDelegatesRequest
	1, // 2: storagepb.Delegate.GetDelegates:output_type -> storagepb.GetDelegatesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_storagepb_delegate_proto_init() }
func file_storagepb_delegate_proto_init() {
	if File_storagepb_delegate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_storagepb_delegate_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetDelegatesRequest); i {
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
		file_storagepb_delegate_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetDelegatesResponse); i {
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
		file_storagepb_delegate_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DelegateEntry); i {
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
	file_storagepb_delegate_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_storagepb_delegate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storagepb_delegate_proto_goTypes,
		DependencyIndexes: file_storagepb_delegate_proto_depIdxs,
		MessageInfos:      file_storagepb_delegate_proto_msgTypes,
	}.Build()
	File_storagepb_delegate_proto = out.File
	file_storagepb_delegate_proto_rawDesc = nil
	file_storagepb_delegate_proto_goTypes = nil
	file_storagepb_delegate_proto_depIdxs = nil
}
