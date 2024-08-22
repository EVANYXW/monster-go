// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.23.4
// source: server/game.proto

package xsf_pb

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

// login --> game 玩家登录
type L_G_Login struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId    string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`          // 账户ID
	ActorId      uint32 `protobuf:"varint,2,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`               // 角色ID
	LoginTime    uint32 `protobuf:"varint,3,opt,name=login_time,json=loginTime,proto3" json:"login_time,omitempty"`         // 登录时间
	ClientId     uint32 `protobuf:"varint,4,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`            // 客户端ID
	LocalAccount string `protobuf:"bytes,5,opt,name=local_account,json=localAccount,proto3" json:"local_account,omitempty"` // 本地账号ID
}

func (x *L_G_Login) Reset() {
	*x = L_G_Login{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *L_G_Login) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*L_G_Login) ProtoMessage() {}

func (x *L_G_Login) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use L_G_Login.ProtoReflect.Descriptor instead.
func (*L_G_Login) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{0}
}

func (x *L_G_Login) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *L_G_Login) GetActorId() uint32 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *L_G_Login) GetLoginTime() uint32 {
	if x != nil {
		return x.LoginTime
	}
	return 0
}

func (x *L_G_Login) GetClientId() uint32 {
	if x != nil {
		return x.ClientId
	}
	return 0
}

func (x *L_G_Login) GetLocalAccount() string {
	if x != nil {
		return x.LocalAccount
	}
	return ""
}

// game --> login 玩家登录反馈
type G_L_LoginResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"` // 账户ID
	ClientId  uint32 `protobuf:"varint,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Result    uint32 `protobuf:"varint,3,opt,name=result,proto3" json:"result,omitempty"`                     // 登录结果
	LoginKey  uint64 `protobuf:"varint,4,opt,name=login_key,json=loginKey,proto3" json:"login_key,omitempty"` // 登录KEY
}

func (x *G_L_LoginResult) Reset() {
	*x = G_L_LoginResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G_L_LoginResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G_L_LoginResult) ProtoMessage() {}

func (x *G_L_LoginResult) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G_L_LoginResult.ProtoReflect.Descriptor instead.
func (*G_L_LoginResult) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{1}
}

func (x *G_L_LoginResult) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *G_L_LoginResult) GetClientId() uint32 {
	if x != nil {
		return x.ClientId
	}
	return 0
}

func (x *G_L_LoginResult) GetResult() uint32 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *G_L_LoginResult) GetLoginKey() uint64 {
	if x != nil {
		return x.LoginKey
	}
	return 0
}

// login --> game GM指令
type L_G_GMCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId      uint32 `protobuf:"varint,1,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"` // 角色ID
	Account      string `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	LocalAccount string `protobuf:"bytes,3,opt,name=local_account,json=localAccount,proto3" json:"local_account,omitempty"`
	ReqId        uint64 `protobuf:"varint,4,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`
	Action       string `protobuf:"bytes,5,opt,name=action,proto3" json:"action,omitempty"`
	QueryJson    string `protobuf:"bytes,6,opt,name=query_json,json=queryJson,proto3" json:"query_json,omitempty"`
}

func (x *L_G_GMCommand) Reset() {
	*x = L_G_GMCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *L_G_GMCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*L_G_GMCommand) ProtoMessage() {}

func (x *L_G_GMCommand) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use L_G_GMCommand.ProtoReflect.Descriptor instead.
func (*L_G_GMCommand) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{2}
}

func (x *L_G_GMCommand) GetActorId() uint32 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *L_G_GMCommand) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *L_G_GMCommand) GetLocalAccount() string {
	if x != nil {
		return x.LocalAccount
	}
	return ""
}

func (x *L_G_GMCommand) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

func (x *L_G_GMCommand) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *L_G_GMCommand) GetQueryJson() string {
	if x != nil {
		return x.QueryJson
	}
	return ""
}

// game --> server 玩家登录
type G_S_ActorLogin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId     uint32 `protobuf:"varint,1,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
	ClientId    uint32 `protobuf:"varint,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	LastMailId  uint32 `protobuf:"varint,3,opt,name=last_mail_id,json=lastMailId,proto3" json:"last_mail_id,omitempty"`
	IsDayChange bool   `protobuf:"varint,4,opt,name=is_day_change,json=isDayChange,proto3" json:"is_day_change,omitempty"`
}

func (x *G_S_ActorLogin) Reset() {
	*x = G_S_ActorLogin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G_S_ActorLogin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G_S_ActorLogin) ProtoMessage() {}

func (x *G_S_ActorLogin) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G_S_ActorLogin.ProtoReflect.Descriptor instead.
func (*G_S_ActorLogin) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{3}
}

func (x *G_S_ActorLogin) GetActorId() uint32 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *G_S_ActorLogin) GetClientId() uint32 {
	if x != nil {
		return x.ClientId
	}
	return 0
}

func (x *G_S_ActorLogin) GetLastMailId() uint32 {
	if x != nil {
		return x.LastMailId
	}
	return 0
}

func (x *G_S_ActorLogin) GetIsDayChange() bool {
	if x != nil {
		return x.IsDayChange
	}
	return false
}

// game --> server 玩家对象释放
type G_S_ActorRelease struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId uint32 `protobuf:"varint,1,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"` // 角色ID，uint32类型
}

func (x *G_S_ActorRelease) Reset() {
	*x = G_S_ActorRelease{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G_S_ActorRelease) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G_S_ActorRelease) ProtoMessage() {}

func (x *G_S_ActorRelease) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G_S_ActorRelease.ProtoReflect.Descriptor instead.
func (*G_S_ActorRelease) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{4}
}

func (x *G_S_ActorRelease) GetActorId() uint32 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

// mail --> game 设置玩家邮件数据
type Ml_G_ActorMailData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId    uint32 `protobuf:"varint,1,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
	LastMailId uint32 `protobuf:"varint,2,opt,name=last_mail_id,json=lastMailId,proto3" json:"last_mail_id,omitempty"`
}

func (x *Ml_G_ActorMailData) Reset() {
	*x = Ml_G_ActorMailData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ml_G_ActorMailData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ml_G_ActorMailData) ProtoMessage() {}

func (x *Ml_G_ActorMailData) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ml_G_ActorMailData.ProtoReflect.Descriptor instead.
func (*Ml_G_ActorMailData) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{5}
}

func (x *Ml_G_ActorMailData) GetActorId() uint32 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *Ml_G_ActorMailData) GetLastMailId() uint32 {
	if x != nil {
		return x.LastMailId
	}
	return 0
}

type MSG_MlItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Count uint32 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *MSG_MlItem) Reset() {
	*x = MSG_MlItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_MlItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_MlItem) ProtoMessage() {}

func (x *MSG_MlItem) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_MlItem.ProtoReflect.Descriptor instead.
func (*MSG_MlItem) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{6}
}

func (x *MSG_MlItem) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MSG_MlItem) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

// mail --> game 添加物品
type Ml_G_AddItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActorId uint32        `protobuf:"varint,1,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
	Items   []*MSG_MlItem `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	MailId  uint64        `protobuf:"varint,3,opt,name=mail_id,json=mailId,proto3" json:"mail_id,omitempty"`
}

func (x *Ml_G_AddItem) Reset() {
	*x = Ml_G_AddItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_game_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ml_G_AddItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ml_G_AddItem) ProtoMessage() {}

func (x *Ml_G_AddItem) ProtoReflect() protoreflect.Message {
	mi := &file_server_game_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ml_G_AddItem.ProtoReflect.Descriptor instead.
func (*Ml_G_AddItem) Descriptor() ([]byte, []int) {
	return file_server_game_proto_rawDescGZIP(), []int{7}
}

func (x *Ml_G_AddItem) GetActorId() uint32 {
	if x != nil {
		return x.ActorId
	}
	return 0
}

func (x *Ml_G_AddItem) GetItems() []*MSG_MlItem {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *Ml_G_AddItem) GetMailId() uint64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

var File_server_game_proto protoreflect.FileDescriptor

var file_server_game_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x4e, 0x4c, 0x44, 0x5f, 0x50, 0x42, 0x22, 0xa6, 0x01, 0x0a, 0x09,
	0x4c, 0x5f, 0x47, 0x5f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x82, 0x01, 0x0a, 0x0f, 0x47, 0x5f, 0x4c, 0x5f, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x4b, 0x65, 0x79, 0x22, 0xb7, 0x01, 0x0a, 0x0d, 0x4c, 0x5f,
	0x47, 0x5f, 0x47, 0x4d, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x61,
	0x63, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61,
	0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x23, 0x0a, 0x0d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x65, 0x71, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x65, 0x71, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x6a, 0x73,
	0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4a,
	0x73, 0x6f, 0x6e, 0x22, 0x8e, 0x01, 0x0a, 0x0e, 0x47, 0x5f, 0x53, 0x5f, 0x41, 0x63, 0x74, 0x6f,
	0x72, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x49,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20,
	0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64,
	0x12, 0x22, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x64, 0x61, 0x79, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x44, 0x61, 0x79, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x22, 0x2d, 0x0a, 0x10, 0x47, 0x5f, 0x53, 0x5f, 0x41, 0x63, 0x74, 0x6f,
	0x72, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x22, 0x51, 0x0a, 0x12, 0x4d, 0x6c, 0x5f, 0x47, 0x5f, 0x41, 0x63, 0x74, 0x6f,
	0x72, 0x4d, 0x61, 0x69, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6d, 0x61, 0x69,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74,
	0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x0a, 0x4d, 0x53, 0x47, 0x5f, 0x4d, 0x6c,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x6c, 0x0a, 0x0c, 0x4d, 0x6c,
	0x5f, 0x47, 0x5f, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x4e, 0x4c, 0x44, 0x5f, 0x50, 0x42, 0x2e, 0x4d, 0x53,
	0x47, 0x5f, 0x4d, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12,
	0x17, 0x0a, 0x07, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x62, 0x2f, 0x78,
	0x73, 0x66, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_game_proto_rawDescOnce sync.Once
	file_server_game_proto_rawDescData = file_server_game_proto_rawDesc
)

func file_server_game_proto_rawDescGZIP() []byte {
	file_server_game_proto_rawDescOnce.Do(func() {
		file_server_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_game_proto_rawDescData)
	})
	return file_server_game_proto_rawDescData
}

var file_server_game_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_server_game_proto_goTypes = []interface{}{
	(*L_G_Login)(nil),          // 0: NLD_PB.L_G_Login
	(*G_L_LoginResult)(nil),    // 1: NLD_PB.G_L_LoginResult
	(*L_G_GMCommand)(nil),      // 2: NLD_PB.L_G_GMCommand
	(*G_S_ActorLogin)(nil),     // 3: NLD_PB.G_S_ActorLogin
	(*G_S_ActorRelease)(nil),   // 4: NLD_PB.G_S_ActorRelease
	(*Ml_G_ActorMailData)(nil), // 5: NLD_PB.Ml_G_ActorMailData
	(*MSG_MlItem)(nil),         // 6: NLD_PB.MSG_MlItem
	(*Ml_G_AddItem)(nil),       // 7: NLD_PB.Ml_G_AddItem
}
var file_server_game_proto_depIdxs = []int32{
	6, // 0: NLD_PB.Ml_G_AddItem.items:type_name -> NLD_PB.MSG_MlItem
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_server_game_proto_init() }
func file_server_game_proto_init() {
	if File_server_game_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*L_G_Login); i {
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
		file_server_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G_L_LoginResult); i {
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
		file_server_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*L_G_GMCommand); i {
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
		file_server_game_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G_S_ActorLogin); i {
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
		file_server_game_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G_S_ActorRelease); i {
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
		file_server_game_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ml_G_ActorMailData); i {
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
		file_server_game_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_MlItem); i {
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
		file_server_game_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ml_G_AddItem); i {
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
			RawDescriptor: file_server_game_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_server_game_proto_goTypes,
		DependencyIndexes: file_server_game_proto_depIdxs,
		MessageInfos:      file_server_game_proto_msgTypes,
	}.Build()
	File_server_game_proto = out.File
	file_server_game_proto_rawDesc = nil
	file_server_game_proto_goTypes = nil
	file_server_game_proto_depIdxs = nil
}