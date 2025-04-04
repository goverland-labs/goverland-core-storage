// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: storagepb/proposal.proto

package storagepb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProposalTimelineItem_TimelineAction int32

const (
	ProposalTimelineItem_Unspecified                 ProposalTimelineItem_TimelineAction = 0
	ProposalTimelineItem_DaoCreated                  ProposalTimelineItem_TimelineAction = 1
	ProposalTimelineItem_DaoUpdated                  ProposalTimelineItem_TimelineAction = 2
	ProposalTimelineItem_ProposalCreated             ProposalTimelineItem_TimelineAction = 3
	ProposalTimelineItem_ProposalUpdated             ProposalTimelineItem_TimelineAction = 4
	ProposalTimelineItem_ProposalVotingStartsSoon    ProposalTimelineItem_TimelineAction = 5
	ProposalTimelineItem_ProposalVotingStarted       ProposalTimelineItem_TimelineAction = 6
	ProposalTimelineItem_ProposalVotingQuorumReached ProposalTimelineItem_TimelineAction = 7
	ProposalTimelineItem_ProposalVotingEnded         ProposalTimelineItem_TimelineAction = 8
	ProposalTimelineItem_ProposalVotingEndsSoon      ProposalTimelineItem_TimelineAction = 9
)

// Enum value maps for ProposalTimelineItem_TimelineAction.
var (
	ProposalTimelineItem_TimelineAction_name = map[int32]string{
		0: "Unspecified",
		1: "DaoCreated",
		2: "DaoUpdated",
		3: "ProposalCreated",
		4: "ProposalUpdated",
		5: "ProposalVotingStartsSoon",
		6: "ProposalVotingStarted",
		7: "ProposalVotingQuorumReached",
		8: "ProposalVotingEnded",
		9: "ProposalVotingEndsSoon",
	}
	ProposalTimelineItem_TimelineAction_value = map[string]int32{
		"Unspecified":                 0,
		"DaoCreated":                  1,
		"DaoUpdated":                  2,
		"ProposalCreated":             3,
		"ProposalUpdated":             4,
		"ProposalVotingStartsSoon":    5,
		"ProposalVotingStarted":       6,
		"ProposalVotingQuorumReached": 7,
		"ProposalVotingEnded":         8,
		"ProposalVotingEndsSoon":      9,
	}
)

func (x ProposalTimelineItem_TimelineAction) Enum() *ProposalTimelineItem_TimelineAction {
	p := new(ProposalTimelineItem_TimelineAction)
	*p = x
	return p
}

func (x ProposalTimelineItem_TimelineAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProposalTimelineItem_TimelineAction) Descriptor() protoreflect.EnumDescriptor {
	return file_storagepb_proposal_proto_enumTypes[0].Descriptor()
}

func (ProposalTimelineItem_TimelineAction) Type() protoreflect.EnumType {
	return &file_storagepb_proposal_proto_enumTypes[0]
}

func (x ProposalTimelineItem_TimelineAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProposalTimelineItem_TimelineAction.Descriptor instead.
func (ProposalTimelineItem_TimelineAction) EnumDescriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{2, 0}
}

type ProposalByIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProposalId    string                 `protobuf:"bytes,1,opt,name=proposal_id,json=proposalId,proto3" json:"proposal_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalByIDRequest) Reset() {
	*x = ProposalByIDRequest{}
	mi := &file_storagepb_proposal_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalByIDRequest) ProtoMessage() {}

func (x *ProposalByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_proposal_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalByIDRequest.ProtoReflect.Descriptor instead.
func (*ProposalByIDRequest) Descriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{0}
}

func (x *ProposalByIDRequest) GetProposalId() string {
	if x != nil {
		return x.ProposalId
	}
	return ""
}

type ProposalInfo struct {
	state         protoimpl.MessageState  `protogen:"open.v1"`
	Id            string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp  `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp  `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Ipfs          string                  `protobuf:"bytes,4,opt,name=ipfs,proto3" json:"ipfs,omitempty"`
	Author        string                  `protobuf:"bytes,5,opt,name=author,proto3" json:"author,omitempty"`
	DaoId         string                  `protobuf:"bytes,6,opt,name=dao_id,json=daoId,proto3" json:"dao_id,omitempty"`
	Created       uint64                  `protobuf:"varint,7,opt,name=created,proto3" json:"created,omitempty"`
	Network       string                  `protobuf:"bytes,8,opt,name=network,proto3" json:"network,omitempty"`
	Symbol        string                  `protobuf:"bytes,9,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Type          string                  `protobuf:"bytes,10,opt,name=type,proto3" json:"type,omitempty"`
	Strategies    []*Strategy             `protobuf:"bytes,11,rep,name=strategies,proto3" json:"strategies,omitempty"`
	Title         string                  `protobuf:"bytes,12,opt,name=title,proto3" json:"title,omitempty"`
	Body          string                  `protobuf:"bytes,13,opt,name=body,proto3" json:"body,omitempty"`
	Discussion    string                  `protobuf:"bytes,14,opt,name=discussion,proto3" json:"discussion,omitempty"`
	Choices       []string                `protobuf:"bytes,15,rep,name=choices,proto3" json:"choices,omitempty"`
	Start         uint64                  `protobuf:"varint,16,opt,name=start,proto3" json:"start,omitempty"`
	End           uint64                  `protobuf:"varint,17,opt,name=end,proto3" json:"end,omitempty"`
	Quorum        float32                 `protobuf:"fixed32,18,opt,name=quorum,proto3" json:"quorum,omitempty"`
	Privacy       string                  `protobuf:"bytes,19,opt,name=privacy,proto3" json:"privacy,omitempty"`
	Snapshot      string                  `protobuf:"bytes,20,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
	State         string                  `protobuf:"bytes,21,opt,name=state,proto3" json:"state,omitempty"`
	Link          string                  `protobuf:"bytes,22,opt,name=link,proto3" json:"link,omitempty"`
	App           string                  `protobuf:"bytes,23,opt,name=app,proto3" json:"app,omitempty"`
	Scores        []float32               `protobuf:"fixed32,24,rep,packed,name=scores,proto3" json:"scores,omitempty"`
	ScoresState   string                  `protobuf:"bytes,25,opt,name=scores_state,json=scoresState,proto3" json:"scores_state,omitempty"`
	ScoresTotal   float32                 `protobuf:"fixed32,26,opt,name=scores_total,json=scoresTotal,proto3" json:"scores_total,omitempty"`
	ScoresUpdated uint64                  `protobuf:"varint,27,opt,name=scores_updated,json=scoresUpdated,proto3" json:"scores_updated,omitempty"`
	Votes         uint64                  `protobuf:"varint,28,opt,name=votes,proto3" json:"votes,omitempty"`
	Timeline      []*ProposalTimelineItem `protobuf:"bytes,29,rep,name=timeline,proto3" json:"timeline,omitempty"`
	EnsName       string                  `protobuf:"bytes,30,opt,name=ens_name,json=ensName,proto3" json:"ens_name,omitempty"`
	Spam          bool                    `protobuf:"varint,31,opt,name=spam,proto3" json:"spam,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalInfo) Reset() {
	*x = ProposalInfo{}
	mi := &file_storagepb_proposal_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalInfo) ProtoMessage() {}

func (x *ProposalInfo) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_proposal_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalInfo.ProtoReflect.Descriptor instead.
func (*ProposalInfo) Descriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{1}
}

func (x *ProposalInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ProposalInfo) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *ProposalInfo) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *ProposalInfo) GetIpfs() string {
	if x != nil {
		return x.Ipfs
	}
	return ""
}

func (x *ProposalInfo) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *ProposalInfo) GetDaoId() string {
	if x != nil {
		return x.DaoId
	}
	return ""
}

func (x *ProposalInfo) GetCreated() uint64 {
	if x != nil {
		return x.Created
	}
	return 0
}

func (x *ProposalInfo) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *ProposalInfo) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *ProposalInfo) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ProposalInfo) GetStrategies() []*Strategy {
	if x != nil {
		return x.Strategies
	}
	return nil
}

func (x *ProposalInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ProposalInfo) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *ProposalInfo) GetDiscussion() string {
	if x != nil {
		return x.Discussion
	}
	return ""
}

func (x *ProposalInfo) GetChoices() []string {
	if x != nil {
		return x.Choices
	}
	return nil
}

func (x *ProposalInfo) GetStart() uint64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *ProposalInfo) GetEnd() uint64 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *ProposalInfo) GetQuorum() float32 {
	if x != nil {
		return x.Quorum
	}
	return 0
}

func (x *ProposalInfo) GetPrivacy() string {
	if x != nil {
		return x.Privacy
	}
	return ""
}

func (x *ProposalInfo) GetSnapshot() string {
	if x != nil {
		return x.Snapshot
	}
	return ""
}

func (x *ProposalInfo) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *ProposalInfo) GetLink() string {
	if x != nil {
		return x.Link
	}
	return ""
}

func (x *ProposalInfo) GetApp() string {
	if x != nil {
		return x.App
	}
	return ""
}

func (x *ProposalInfo) GetScores() []float32 {
	if x != nil {
		return x.Scores
	}
	return nil
}

func (x *ProposalInfo) GetScoresState() string {
	if x != nil {
		return x.ScoresState
	}
	return ""
}

func (x *ProposalInfo) GetScoresTotal() float32 {
	if x != nil {
		return x.ScoresTotal
	}
	return 0
}

func (x *ProposalInfo) GetScoresUpdated() uint64 {
	if x != nil {
		return x.ScoresUpdated
	}
	return 0
}

func (x *ProposalInfo) GetVotes() uint64 {
	if x != nil {
		return x.Votes
	}
	return 0
}

func (x *ProposalInfo) GetTimeline() []*ProposalTimelineItem {
	if x != nil {
		return x.Timeline
	}
	return nil
}

func (x *ProposalInfo) GetEnsName() string {
	if x != nil {
		return x.EnsName
	}
	return ""
}

func (x *ProposalInfo) GetSpam() bool {
	if x != nil {
		return x.Spam
	}
	return false
}

type ProposalTimelineItem struct {
	state         protoimpl.MessageState              `protogen:"open.v1"`
	CreatedAt     *timestamppb.Timestamp              `protobuf:"bytes,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Action        ProposalTimelineItem_TimelineAction `protobuf:"varint,2,opt,name=action,proto3,enum=storagepb.ProposalTimelineItem_TimelineAction" json:"action,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalTimelineItem) Reset() {
	*x = ProposalTimelineItem{}
	mi := &file_storagepb_proposal_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalTimelineItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalTimelineItem) ProtoMessage() {}

func (x *ProposalTimelineItem) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_proposal_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalTimelineItem.ProtoReflect.Descriptor instead.
func (*ProposalTimelineItem) Descriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{2}
}

func (x *ProposalTimelineItem) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *ProposalTimelineItem) GetAction() ProposalTimelineItem_TimelineAction {
	if x != nil {
		return x.Action
	}
	return ProposalTimelineItem_Unspecified
}

type ProposalByIDResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Proposal      *ProposalInfo          `protobuf:"bytes,1,opt,name=proposal,proto3" json:"proposal,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalByIDResponse) Reset() {
	*x = ProposalByIDResponse{}
	mi := &file_storagepb_proposal_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalByIDResponse) ProtoMessage() {}

func (x *ProposalByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_proposal_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalByIDResponse.ProtoReflect.Descriptor instead.
func (*ProposalByIDResponse) Descriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{3}
}

func (x *ProposalByIDResponse) GetProposal() *ProposalInfo {
	if x != nil {
		return x.Proposal
	}
	return nil
}

type ProposalByFilterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dao           *string                `protobuf:"bytes,1,opt,name=dao,proto3,oneof" json:"dao,omitempty"`
	Category      *string                `protobuf:"bytes,2,opt,name=category,proto3,oneof" json:"category,omitempty"`
	Limit         *uint64                `protobuf:"varint,3,opt,name=limit,proto3,oneof" json:"limit,omitempty"`
	Offset        *uint64                `protobuf:"varint,4,opt,name=offset,proto3,oneof" json:"offset,omitempty"`
	Title         *string                `protobuf:"bytes,5,opt,name=title,proto3,oneof" json:"title,omitempty"`
	Order         *string                `protobuf:"bytes,6,opt,name=order,proto3,oneof" json:"order,omitempty"`
	Top           *bool                  `protobuf:"varint,7,opt,name=top,proto3,oneof" json:"top,omitempty"`
	ProposalIds   []string               `protobuf:"bytes,8,rep,name=proposal_ids,json=proposalIds,proto3" json:"proposal_ids,omitempty"`
	OnlyActive    *bool                  `protobuf:"varint,9,opt,name=only_active,json=onlyActive,proto3,oneof" json:"only_active,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalByFilterRequest) Reset() {
	*x = ProposalByFilterRequest{}
	mi := &file_storagepb_proposal_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalByFilterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalByFilterRequest) ProtoMessage() {}

func (x *ProposalByFilterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_proposal_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalByFilterRequest.ProtoReflect.Descriptor instead.
func (*ProposalByFilterRequest) Descriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{4}
}

func (x *ProposalByFilterRequest) GetDao() string {
	if x != nil && x.Dao != nil {
		return *x.Dao
	}
	return ""
}

func (x *ProposalByFilterRequest) GetCategory() string {
	if x != nil && x.Category != nil {
		return *x.Category
	}
	return ""
}

func (x *ProposalByFilterRequest) GetLimit() uint64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *ProposalByFilterRequest) GetOffset() uint64 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

func (x *ProposalByFilterRequest) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *ProposalByFilterRequest) GetOrder() string {
	if x != nil && x.Order != nil {
		return *x.Order
	}
	return ""
}

func (x *ProposalByFilterRequest) GetTop() bool {
	if x != nil && x.Top != nil {
		return *x.Top
	}
	return false
}

func (x *ProposalByFilterRequest) GetProposalIds() []string {
	if x != nil {
		return x.ProposalIds
	}
	return nil
}

func (x *ProposalByFilterRequest) GetOnlyActive() bool {
	if x != nil && x.OnlyActive != nil {
		return *x.OnlyActive
	}
	return false
}

type ProposalByFilterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Proposals     []*ProposalInfo        `protobuf:"bytes,1,rep,name=proposals,proto3" json:"proposals,omitempty"`
	TotalCount    uint64                 `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalByFilterResponse) Reset() {
	*x = ProposalByFilterResponse{}
	mi := &file_storagepb_proposal_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalByFilterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalByFilterResponse) ProtoMessage() {}

func (x *ProposalByFilterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_storagepb_proposal_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalByFilterResponse.ProtoReflect.Descriptor instead.
func (*ProposalByFilterResponse) Descriptor() ([]byte, []int) {
	return file_storagepb_proposal_proto_rawDescGZIP(), []int{5}
}

func (x *ProposalByFilterResponse) GetProposals() []*ProposalInfo {
	if x != nil {
		return x.Proposals
	}
	return nil
}

func (x *ProposalByFilterResponse) GetTotalCount() uint64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

var File_storagepb_proposal_proto protoreflect.FileDescriptor

var file_storagepb_proposal_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x70,
	0x6f, 0x73, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70,
	0x62, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x13,
	0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73,
	0x61, 0x6c, 0x49, 0x64, 0x22, 0x89, 0x07, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x69,
	0x70, 0x66, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x70, 0x66, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x61, 0x6f, 0x5f, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x61, 0x6f, 0x49, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x33,
	0x0a, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x53,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67,
	0x69, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x1e, 0x0a,
	0x0a, 0x64, 0x69, 0x73, 0x63, 0x75, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x63, 0x75, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x18, 0x10, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x65, 0x6e, 0x64, 0x18, 0x11, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x18, 0x12, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x06, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61,
	0x63, 0x79, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x63,
	0x79, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x16, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x70, 0x18, 0x17,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x6f,
	0x72, 0x65, 0x73, 0x18, 0x18, 0x20, 0x03, 0x28, 0x02, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x73, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x19, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x5f, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x73, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x73, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0d, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76,
	0x6f, 0x74, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x18, 0x1d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x1e, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x70, 0x61, 0x6d, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x73, 0x70, 0x61, 0x6d,
	0x22, 0x96, 0x03, 0x0a, 0x14, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x54, 0x69, 0x6d,
	0x65, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x46, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62,
	0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xfa, 0x01, 0x0a,
	0x0e, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x10, 0x00,
	0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x61, 0x6f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x10, 0x01,
	0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x61, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x10, 0x02,
	0x12, 0x13, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x10, 0x04, 0x12, 0x1c, 0x0a, 0x18, 0x50, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x73, 0x53, 0x6f, 0x6f, 0x6e, 0x10, 0x05, 0x12, 0x19, 0x0a, 0x15, 0x50, 0x72, 0x6f, 0x70,
	0x6f, 0x73, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65,
	0x64, 0x10, 0x06, 0x12, 0x1f, 0x0a, 0x1b, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56,
	0x6f, 0x74, 0x69, 0x6e, 0x67, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x52, 0x65, 0x61, 0x63, 0x68,
	0x65, 0x64, 0x10, 0x07, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c,
	0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x45, 0x6e, 0x64, 0x65, 0x64, 0x10, 0x08, 0x12, 0x1a, 0x0a,
	0x16, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x45,
	0x6e, 0x64, 0x73, 0x53, 0x6f, 0x6f, 0x6e, 0x10, 0x09, 0x22, 0x4b, 0x0a, 0x14, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x61, 0x6c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x33, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e,
	0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x70, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x22, 0xf5, 0x02, 0x0a, 0x17, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x61, 0x6c, 0x42, 0x79, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x15, 0x0a, 0x03, 0x64, 0x61, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x03, 0x64, 0x61, 0x6f, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x48, 0x02, 0x52, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x48, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x04, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a,
	0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x05, 0x52, 0x05,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x74, 0x6f, 0x70, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x08, 0x48, 0x06, 0x52, 0x03, 0x74, 0x6f, 0x70, 0x88, 0x01, 0x01, 0x12,
	0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x73, 0x18,
	0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x49,
	0x64, 0x73, 0x12, 0x24, 0x0a, 0x0b, 0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x48, 0x07, 0x52, 0x0a, 0x6f, 0x6e, 0x6c, 0x79, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x64, 0x61, 0x6f,
	0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x42, 0x08, 0x0a,
	0x06, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x42, 0x08, 0x0a, 0x06,
	0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x74, 0x6f, 0x70, 0x42, 0x0e,
	0x0a, 0x0c, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x22, 0x72,
	0x0a, 0x18, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x42, 0x79, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x09, 0x70, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73,
	0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c,
	0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x32, 0xae, 0x01, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12,
	0x4a, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x12, 0x1e, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x42, 0x79, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x42,
	0x79, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x61, 0x6c, 0x42, 0x79, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x3b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storagepb_proposal_proto_rawDescOnce sync.Once
	file_storagepb_proposal_proto_rawDescData = file_storagepb_proposal_proto_rawDesc
)

func file_storagepb_proposal_proto_rawDescGZIP() []byte {
	file_storagepb_proposal_proto_rawDescOnce.Do(func() {
		file_storagepb_proposal_proto_rawDescData = protoimpl.X.CompressGZIP(file_storagepb_proposal_proto_rawDescData)
	})
	return file_storagepb_proposal_proto_rawDescData
}

var file_storagepb_proposal_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_storagepb_proposal_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_storagepb_proposal_proto_goTypes = []any{
	(ProposalTimelineItem_TimelineAction)(0), // 0: storagepb.ProposalTimelineItem.TimelineAction
	(*ProposalByIDRequest)(nil),              // 1: storagepb.ProposalByIDRequest
	(*ProposalInfo)(nil),                     // 2: storagepb.ProposalInfo
	(*ProposalTimelineItem)(nil),             // 3: storagepb.ProposalTimelineItem
	(*ProposalByIDResponse)(nil),             // 4: storagepb.ProposalByIDResponse
	(*ProposalByFilterRequest)(nil),          // 5: storagepb.ProposalByFilterRequest
	(*ProposalByFilterResponse)(nil),         // 6: storagepb.ProposalByFilterResponse
	(*timestamppb.Timestamp)(nil),            // 7: google.protobuf.Timestamp
	(*Strategy)(nil),                         // 8: storagepb.Strategy
}
var file_storagepb_proposal_proto_depIdxs = []int32{
	7,  // 0: storagepb.ProposalInfo.created_at:type_name -> google.protobuf.Timestamp
	7,  // 1: storagepb.ProposalInfo.updated_at:type_name -> google.protobuf.Timestamp
	8,  // 2: storagepb.ProposalInfo.strategies:type_name -> storagepb.Strategy
	3,  // 3: storagepb.ProposalInfo.timeline:type_name -> storagepb.ProposalTimelineItem
	7,  // 4: storagepb.ProposalTimelineItem.created_at:type_name -> google.protobuf.Timestamp
	0,  // 5: storagepb.ProposalTimelineItem.action:type_name -> storagepb.ProposalTimelineItem.TimelineAction
	2,  // 6: storagepb.ProposalByIDResponse.proposal:type_name -> storagepb.ProposalInfo
	2,  // 7: storagepb.ProposalByFilterResponse.proposals:type_name -> storagepb.ProposalInfo
	1,  // 8: storagepb.Proposal.GetByID:input_type -> storagepb.ProposalByIDRequest
	5,  // 9: storagepb.Proposal.GetByFilter:input_type -> storagepb.ProposalByFilterRequest
	4,  // 10: storagepb.Proposal.GetByID:output_type -> storagepb.ProposalByIDResponse
	6,  // 11: storagepb.Proposal.GetByFilter:output_type -> storagepb.ProposalByFilterResponse
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_storagepb_proposal_proto_init() }
func file_storagepb_proposal_proto_init() {
	if File_storagepb_proposal_proto != nil {
		return
	}
	file_storagepb_base_proto_init()
	file_storagepb_proposal_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_storagepb_proposal_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storagepb_proposal_proto_goTypes,
		DependencyIndexes: file_storagepb_proposal_proto_depIdxs,
		EnumInfos:         file_storagepb_proposal_proto_enumTypes,
		MessageInfos:      file_storagepb_proposal_proto_msgTypes,
	}.Build()
	File_storagepb_proposal_proto = out.File
	file_storagepb_proposal_proto_rawDesc = nil
	file_storagepb_proposal_proto_goTypes = nil
	file_storagepb_proposal_proto_depIdxs = nil
}
