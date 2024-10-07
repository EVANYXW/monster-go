// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.2.0
// source: uid.proto

package pb

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

type UMFlag int32

const (
	UMFlag_UMUidMetaMin UMFlag = 0
	UMFlag_UMWorldId    UMFlag = 1
	UMFlag_UMCurRoleId  UMFlag = 2
	UMFlag_UMCurClubId  UMFlag = 3
)

// Enum value maps for UMFlag.
var (
	UMFlag_name = map[int32]string{
		0: "UMUidMetaMin",
		1: "UMWorldId",
		2: "UMCurRoleId",
		3: "UMCurClubId",
	}
	UMFlag_value = map[string]int32{
		"UMUidMetaMin": 0,
		"UMWorldId":    1,
		"UMCurRoleId":  2,
		"UMCurClubId":  3,
	}
)

func (x UMFlag) Enum() *UMFlag {
	p := new(UMFlag)
	*p = x
	return p
}

func (x UMFlag) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UMFlag) Descriptor() protoreflect.EnumDescriptor {
	return file_uid_proto_enumTypes[0].Descriptor()
}

func (UMFlag) Type() protoreflect.EnumType {
	return &file_uid_proto_enumTypes[0]
}

func (x UMFlag) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UMFlag.Descriptor instead.
func (UMFlag) EnumDescriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{0}
}

type URFlag int32

const (
	URFlag_URUidRoleIdMin URFlag = 0
	URFlag_UROpenId       URFlag = 1
	URFlag_URRoleId       URFlag = 2
)

// Enum value maps for URFlag.
var (
	URFlag_name = map[int32]string{
		0: "URUidRoleIdMin",
		1: "UROpenId",
		2: "URRoleId",
	}
	URFlag_value = map[string]int32{
		"URUidRoleIdMin": 0,
		"UROpenId":       1,
		"URRoleId":       2,
	}
)

func (x URFlag) Enum() *URFlag {
	p := new(URFlag)
	*p = x
	return p
}

func (x URFlag) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (URFlag) Descriptor() protoreflect.EnumDescriptor {
	return file_uid_proto_enumTypes[1].Descriptor()
}

func (URFlag) Type() protoreflect.EnumType {
	return &file_uid_proto_enumTypes[1]
}

func (x URFlag) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use URFlag.Descriptor instead.
func (URFlag) EnumDescriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{1}
}

type UCFlag int32

const (
	UCFlag_UCUidClubIdMin UCFlag = 0
	UCFlag_UCClubId       UCFlag = 1
	UCFlag_UCOwnerId      UCFlag = 2
)

// Enum value maps for UCFlag.
var (
	UCFlag_name = map[int32]string{
		0: "UCUidClubIdMin",
		1: "UCClubId",
		2: "UCOwnerId",
	}
	UCFlag_value = map[string]int32{
		"UCUidClubIdMin": 0,
		"UCClubId":       1,
		"UCOwnerId":      2,
	}
)

func (x UCFlag) Enum() *UCFlag {
	p := new(UCFlag)
	*p = x
	return p
}

func (x UCFlag) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UCFlag) Descriptor() protoreflect.EnumDescriptor {
	return file_uid_proto_enumTypes[2].Descriptor()
}

func (UCFlag) Type() protoreflect.EnumType {
	return &file_uid_proto_enumTypes[2]
}

func (x UCFlag) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UCFlag.Descriptor instead.
func (UCFlag) EnumDescriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{2}
}

type UidMetaDB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorldId   string `protobuf:"bytes,1,opt,name=WorldId,proto3" json:"WorldId,omitempty"`
	CurRoleId int64  `protobuf:"varint,2,opt,name=CurRoleId,proto3" json:"CurRoleId,omitempty"`
	CurClubId int64  `protobuf:"varint,3,opt,name=CurClubId,proto3" json:"CurClubId,omitempty"`
}

func (x *UidMetaDB) Reset() {
	*x = UidMetaDB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UidMetaDB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UidMetaDB) ProtoMessage() {}

func (x *UidMetaDB) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UidMetaDB.ProtoReflect.Descriptor instead.
func (*UidMetaDB) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{0}
}

func (x *UidMetaDB) GetWorldId() string {
	if x != nil {
		return x.WorldId
	}
	return ""
}

func (x *UidMetaDB) GetCurRoleId() int64 {
	if x != nil {
		return x.CurRoleId
	}
	return 0
}

func (x *UidMetaDB) GetCurClubId() int64 {
	if x != nil {
		return x.CurClubId
	}
	return 0
}

type UidRoleIdDB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpenId string `protobuf:"bytes,1,opt,name=OpenId,proto3" json:"OpenId,omitempty"`
	RoleId int64  `protobuf:"varint,2,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *UidRoleIdDB) Reset() {
	*x = UidRoleIdDB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UidRoleIdDB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UidRoleIdDB) ProtoMessage() {}

func (x *UidRoleIdDB) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UidRoleIdDB.ProtoReflect.Descriptor instead.
func (*UidRoleIdDB) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{1}
}

func (x *UidRoleIdDB) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

func (x *UidRoleIdDB) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type UidClubIdDB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClubId  int64 `protobuf:"varint,1,opt,name=ClubId,proto3" json:"ClubId,omitempty"`
	OwnerId int64 `protobuf:"varint,2,opt,name=OwnerId,proto3" json:"OwnerId,omitempty"`
}

func (x *UidClubIdDB) Reset() {
	*x = UidClubIdDB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UidClubIdDB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UidClubIdDB) ProtoMessage() {}

func (x *UidClubIdDB) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UidClubIdDB.ProtoReflect.Descriptor instead.
func (*UidClubIdDB) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{2}
}

func (x *UidClubIdDB) GetClubId() int64 {
	if x != nil {
		return x.ClubId
	}
	return 0
}

func (x *UidClubIdDB) GetOwnerId() int64 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

type AllocRoleId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AllocRoleId) Reset() {
	*x = AllocRoleId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocRoleId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocRoleId) ProtoMessage() {}

func (x *AllocRoleId) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocRoleId.ProtoReflect.Descriptor instead.
func (*AllocRoleId) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{3}
}

type AllocClubId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AllocClubId) Reset() {
	*x = AllocClubId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocClubId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocClubId) ProtoMessage() {}

func (x *AllocClubId) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocClubId.ProtoReflect.Descriptor instead.
func (*AllocClubId) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{4}
}

type AllocRoleId_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpenId string `protobuf:"bytes,1,opt,name=OpenId,proto3" json:"OpenId,omitempty"`
}

func (x *AllocRoleId_Req) Reset() {
	*x = AllocRoleId_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocRoleId_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocRoleId_Req) ProtoMessage() {}

func (x *AllocRoleId_Req) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocRoleId_Req.ProtoReflect.Descriptor instead.
func (*AllocRoleId_Req) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{3, 0}
}

func (x *AllocRoleId_Req) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

type AllocRoleId_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId int64 `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *AllocRoleId_Rsp) Reset() {
	*x = AllocRoleId_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocRoleId_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocRoleId_Rsp) ProtoMessage() {}

func (x *AllocRoleId_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocRoleId_Rsp.ProtoReflect.Descriptor instead.
func (*AllocRoleId_Rsp) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{3, 1}
}

func (x *AllocRoleId_Rsp) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type AllocClubId_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId int64 `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *AllocClubId_Req) Reset() {
	*x = AllocClubId_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocClubId_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocClubId_Req) ProtoMessage() {}

func (x *AllocClubId_Req) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocClubId_Req.ProtoReflect.Descriptor instead.
func (*AllocClubId_Req) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{4, 0}
}

func (x *AllocClubId_Req) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type AllocClubId_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClubId int64 `protobuf:"varint,1,opt,name=ClubId,proto3" json:"ClubId,omitempty"`
}

func (x *AllocClubId_Rsp) Reset() {
	*x = AllocClubId_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uid_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocClubId_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocClubId_Rsp) ProtoMessage() {}

func (x *AllocClubId_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_uid_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocClubId_Rsp.ProtoReflect.Descriptor instead.
func (*AllocClubId_Rsp) Descriptor() ([]byte, []int) {
	return file_uid_proto_rawDescGZIP(), []int{4, 1}
}

func (x *AllocClubId_Rsp) GetClubId() int64 {
	if x != nil {
		return x.ClubId
	}
	return 0
}

var File_uid_proto protoreflect.FileDescriptor

var file_uid_proto_rawDesc = []byte{
	0x0a, 0x09, 0x75, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x09, 0x55, 0x69, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x42, 0x12,
	0x2b, 0x0a, 0x07, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x11, 0x8a, 0xa6, 0x1d, 0x0d, 0x0a, 0x09, 0x55, 0x4d, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49,
	0x64, 0x10, 0x01, 0x52, 0x07, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x09,
	0x43, 0x75, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42,
	0x11, 0x8a, 0xa6, 0x1d, 0x0d, 0x0a, 0x0b, 0x55, 0x4d, 0x43, 0x75, 0x72, 0x52, 0x6f, 0x6c, 0x65,
	0x49, 0x64, 0x52, 0x09, 0x43, 0x75, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a,
	0x09, 0x43, 0x75, 0x72, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x42, 0x11, 0x8a, 0xa6, 0x1d, 0x0d, 0x0a, 0x0b, 0x55, 0x4d, 0x43, 0x75, 0x72, 0x43, 0x6c, 0x75,
	0x62, 0x49, 0x64, 0x52, 0x09, 0x43, 0x75, 0x72, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x22, 0x5f,
	0x0a, 0x0b, 0x55, 0x69, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x44, 0x42, 0x12, 0x28, 0x0a,
	0x06, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x10, 0x8a,
	0xa6, 0x1d, 0x0c, 0x0a, 0x08, 0x55, 0x52, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x10, 0x01, 0x52,
	0x06, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0e, 0x8a, 0xa6, 0x1d, 0x0a, 0x0a, 0x08, 0x55,
	0x52, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x22,
	0x62, 0x0a, 0x0b, 0x55, 0x69, 0x64, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x44, 0x42, 0x12, 0x28,
	0x0a, 0x06, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x10,
	0x8a, 0xa6, 0x1d, 0x0c, 0x0a, 0x08, 0x55, 0x43, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x10, 0x01,
	0x52, 0x06, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x07, 0x4f, 0x77, 0x6e, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0f, 0x8a, 0xa6, 0x1d, 0x0b, 0x0a,
	0x09, 0x55, 0x43, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x52, 0x07, 0x4f, 0x77, 0x6e, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x59, 0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x52, 0x6f, 0x6c, 0x65,
	0x49, 0x64, 0x1a, 0x1d, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x4f, 0x70, 0x65,
	0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4f, 0x70, 0x65, 0x6e, 0x49,
	0x64, 0x1a, 0x1d, 0x0a, 0x03, 0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64,
	0x3a, 0x0c, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0xca, 0xf3, 0x18, 0x02, 0x18, 0x01, 0x22, 0x59,
	0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x1a, 0x1d, 0x0a,
	0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x1a, 0x1d, 0x0a, 0x03,
	0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x3a, 0x0c, 0xca, 0xf3, 0x18,
	0x02, 0x0a, 0x00, 0xca, 0xf3, 0x18, 0x02, 0x18, 0x01, 0x2a, 0x4b, 0x0a, 0x06, 0x55, 0x4d, 0x46,
	0x6c, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x4d, 0x55, 0x69, 0x64, 0x4d, 0x65, 0x74, 0x61,
	0x4d, 0x69, 0x6e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x4d, 0x57, 0x6f, 0x72, 0x6c, 0x64,
	0x49, 0x64, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4d, 0x43, 0x75, 0x72, 0x52, 0x6f, 0x6c,
	0x65, 0x49, 0x64, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4d, 0x43, 0x75, 0x72, 0x43, 0x6c,
	0x75, 0x62, 0x49, 0x64, 0x10, 0x03, 0x2a, 0x38, 0x0a, 0x06, 0x55, 0x52, 0x46, 0x6c, 0x61, 0x67,
	0x12, 0x12, 0x0a, 0x0e, 0x55, 0x52, 0x55, 0x69, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x4d,
	0x69, 0x6e, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x52, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64,
	0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x52, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x10, 0x02,
	0x2a, 0x39, 0x0a, 0x06, 0x55, 0x43, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x43,
	0x55, 0x69, 0x64, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x4d, 0x69, 0x6e, 0x10, 0x00, 0x12, 0x0c,
	0x0a, 0x08, 0x55, 0x43, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09,
	0x55, 0x43, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x10, 0x02, 0x32, 0x86, 0x01, 0x0a, 0x06,
	0x55, 0x69, 0x64, 0x53, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x52,
	0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c,
	0x6c, 0x6f, 0x63, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x52, 0x6f, 0x6c, 0x65, 0x49,
	0x64, 0x2e, 0x52, 0x73, 0x70, 0x12, 0x3d, 0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x43, 0x6c,
	0x75, 0x62, 0x49, 0x64, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c, 0x6c,
	0x6f, 0x63, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64,
	0x2e, 0x52, 0x73, 0x70, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_uid_proto_rawDescOnce sync.Once
	file_uid_proto_rawDescData = file_uid_proto_rawDesc
)

func file_uid_proto_rawDescGZIP() []byte {
	file_uid_proto_rawDescOnce.Do(func() {
		file_uid_proto_rawDescData = protoimpl.X.CompressGZIP(file_uid_proto_rawDescData)
	})
	return file_uid_proto_rawDescData
}

var file_uid_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_uid_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_uid_proto_goTypes = []interface{}{
	(UMFlag)(0),             // 0: proto.UMFlag
	(URFlag)(0),             // 1: proto.URFlag
	(UCFlag)(0),             // 2: proto.UCFlag
	(*UidMetaDB)(nil),       // 3: proto.UidMetaDB
	(*UidRoleIdDB)(nil),     // 4: proto.UidRoleIdDB
	(*UidClubIdDB)(nil),     // 5: proto.UidClubIdDB
	(*AllocRoleId)(nil),     // 6: proto.AllocRoleId
	(*AllocClubId)(nil),     // 7: proto.AllocClubId
	(*AllocRoleId_Req)(nil), // 8: proto.AllocRoleId.Req
	(*AllocRoleId_Rsp)(nil), // 9: proto.AllocRoleId.Rsp
	(*AllocClubId_Req)(nil), // 10: proto.AllocClubId.Req
	(*AllocClubId_Rsp)(nil), // 11: proto.AllocClubId.Rsp
}
var file_uid_proto_depIdxs = []int32{
	8,  // 0: proto.UidSer.AllocRoleId:input_type -> proto.AllocRoleId.Req
	10, // 1: proto.UidSer.AllocClubId:input_type -> proto.AllocClubId.Req
	9,  // 2: proto.UidSer.AllocRoleId:output_type -> proto.AllocRoleId.Rsp
	11, // 3: proto.UidSer.AllocClubId:output_type -> proto.AllocClubId.Rsp
	2,  // [2:4] is the sub-list for method output_type
	0,  // [0:2] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_uid_proto_init() }
func file_uid_proto_init() {
	if File_uid_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_uid_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UidMetaDB); i {
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
		file_uid_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UidRoleIdDB); i {
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
		file_uid_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UidClubIdDB); i {
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
		file_uid_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocRoleId); i {
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
		file_uid_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocClubId); i {
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
		file_uid_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocRoleId_Req); i {
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
		file_uid_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocRoleId_Rsp); i {
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
		file_uid_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocClubId_Req); i {
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
		file_uid_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocClubId_Rsp); i {
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
			RawDescriptor: file_uid_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_uid_proto_goTypes,
		DependencyIndexes: file_uid_proto_depIdxs,
		EnumInfos:         file_uid_proto_enumTypes,
		MessageInfos:      file_uid_proto_msgTypes,
	}.Build()
	File_uid_proto = out.File
	file_uid_proto_rawDesc = nil
	file_uid_proto_goTypes = nil
	file_uid_proto_depIdxs = nil
}