// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: api/account/v1/account.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type TriggerType int32

const (
	TriggerType_TRIGGER_UNKNOWN   TriggerType = 0 // 未知触发方式（默认值）
	TriggerType_TRIGGER_MANUAL    TriggerType = 1 // 手动触发（如管理员点击按钮）
	TriggerType_TRIGGER_SCHEDULED TriggerType = 2 // 定时任务触发（如每天凌晨2点自动同步）
)

// Enum value maps for TriggerType.
var (
	TriggerType_name = map[int32]string{
		0: "TRIGGER_UNKNOWN",
		1: "TRIGGER_MANUAL",
		2: "TRIGGER_SCHEDULED",
	}
	TriggerType_value = map[string]int32{
		"TRIGGER_UNKNOWN":   0,
		"TRIGGER_MANUAL":    1,
		"TRIGGER_SCHEDULED": 2,
	}
)

func (x TriggerType) Enum() *TriggerType {
	p := new(TriggerType)
	*p = x
	return p
}

func (x TriggerType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TriggerType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_account_v1_account_proto_enumTypes[0].Descriptor()
}

func (TriggerType) Type() protoreflect.EnumType {
	return &file_api_account_v1_account_proto_enumTypes[0]
}

func (x TriggerType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TriggerType.Descriptor instead.
func (TriggerType) EnumDescriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{0}
}

type SyncType int32

const (
	SyncType_FULL        SyncType = 0 // 全量同步
	SyncType_INCREMENTAL SyncType = 1 // 增量同步
)

// Enum value maps for SyncType.
var (
	SyncType_name = map[int32]string{
		0: "FULL",
		1: "INCREMENTAL",
	}
	SyncType_value = map[string]int32{
		"FULL":        0,
		"INCREMENTAL": 1,
	}
)

func (x SyncType) Enum() *SyncType {
	p := new(SyncType)
	*p = x
	return p
}

func (x SyncType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SyncType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_account_v1_account_proto_enumTypes[1].Descriptor()
}

func (SyncType) Type() protoreflect.EnumType {
	return &file_api_account_v1_account_proto_enumTypes[1]
}

func (x SyncType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SyncType.Descriptor instead.
func (SyncType) EnumDescriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{1}
}

type GetSyncAccountReply_Status int32

const (
	GetSyncAccountReply_PENDING GetSyncAccountReply_Status = 0 // 待执行
	GetSyncAccountReply_RUNNING GetSyncAccountReply_Status = 1 // 执行中
	GetSyncAccountReply_SUCCESS GetSyncAccountReply_Status = 2 // 成功
	GetSyncAccountReply_FAILED  GetSyncAccountReply_Status = 3 // 失败
)

// Enum value maps for GetSyncAccountReply_Status.
var (
	GetSyncAccountReply_Status_name = map[int32]string{
		0: "PENDING",
		1: "RUNNING",
		2: "SUCCESS",
		3: "FAILED",
	}
	GetSyncAccountReply_Status_value = map[string]int32{
		"PENDING": 0,
		"RUNNING": 1,
		"SUCCESS": 2,
		"FAILED":  3,
	}
)

func (x GetSyncAccountReply_Status) Enum() *GetSyncAccountReply_Status {
	p := new(GetSyncAccountReply_Status)
	*p = x
	return p
}

func (x GetSyncAccountReply_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GetSyncAccountReply_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_api_account_v1_account_proto_enumTypes[2].Descriptor()
}

func (GetSyncAccountReply_Status) Type() protoreflect.EnumType {
	return &file_api_account_v1_account_proto_enumTypes[2]
}

func (x GetSyncAccountReply_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GetSyncAccountReply_Status.Descriptor instead.
func (GetSyncAccountReply_Status) EnumDescriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{5, 0}
}

type UploadRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileData      []byte                 `protobuf:"bytes,1,opt,name=file_data,json=fileData,proto3" json:"file_data,omitempty"`
	FileName      string                 `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadRequest) Reset() {
	*x = UploadRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadRequest) ProtoMessage() {}

func (x *UploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadRequest.ProtoReflect.Descriptor instead.
func (*UploadRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{0}
}

func (x *UploadRequest) GetFileData() []byte {
	if x != nil {
		return x.FileData
	}
	return nil
}

func (x *UploadRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type UploadReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileUrl       string                 `protobuf:"bytes,1,opt,name=file_url,json=fileUrl,proto3" json:"file_url,omitempty"`
	FileSize      int64                  `protobuf:"varint,2,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadReply) Reset() {
	*x = UploadReply{}
	mi := &file_api_account_v1_account_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadReply) ProtoMessage() {}

func (x *UploadReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadReply.ProtoReflect.Descriptor instead.
func (*UploadReply) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{1}
}

func (x *UploadReply) GetFileUrl() string {
	if x != nil {
		return x.FileUrl
	}
	return ""
}

func (x *UploadReply) GetFileSize() int64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

// 创建同步请求
type CreateSyncAccountRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TriggerType   TriggerType            `protobuf:"varint,1,opt,name=trigger_type,json=triggerType,proto3,enum=api.account.v1.TriggerType" json:"trigger_type,omitempty"` // 触发类型
	SyncType      SyncType               `protobuf:"varint,2,opt,name=sync_type,json=syncType,proto3,enum=api.account.v1.SyncType" json:"sync_type,omitempty"`             // 同步类型
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateSyncAccountRequest) Reset() {
	*x = CreateSyncAccountRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSyncAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSyncAccountRequest) ProtoMessage() {}

func (x *CreateSyncAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSyncAccountRequest.ProtoReflect.Descriptor instead.
func (*CreateSyncAccountRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{2}
}

func (x *CreateSyncAccountRequest) GetTriggerType() TriggerType {
	if x != nil {
		return x.TriggerType
	}
	return TriggerType_TRIGGER_UNKNOWN
}

func (x *CreateSyncAccountRequest) GetSyncType() SyncType {
	if x != nil {
		return x.SyncType
	}
	return SyncType_FULL
}

// 创建同步响应
type CreateSyncAccountReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`             // 生成的任务ID
	CreateTime    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"` // 任务创建时间
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateSyncAccountReply) Reset() {
	*x = CreateSyncAccountReply{}
	mi := &file_api_account_v1_account_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSyncAccountReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSyncAccountReply) ProtoMessage() {}

func (x *CreateSyncAccountReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSyncAccountReply.ProtoReflect.Descriptor instead.
func (*CreateSyncAccountReply) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{3}
}

func (x *CreateSyncAccountReply) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *CreateSyncAccountReply) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

// 查询同步请求
type GetSyncAccountRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"` // 要查询的任务ID
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSyncAccountRequest) Reset() {
	*x = GetSyncAccountRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSyncAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSyncAccountRequest) ProtoMessage() {}

func (x *GetSyncAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSyncAccountRequest.ProtoReflect.Descriptor instead.
func (*GetSyncAccountRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{4}
}

func (x *GetSyncAccountRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

// 查询同步响应
type GetSyncAccountReply struct {
	state                       protoimpl.MessageState     `protogen:"open.v1"`
	Status                      GetSyncAccountReply_Status `protobuf:"varint,1,opt,name=status,proto3,enum=api.account.v1.GetSyncAccountReply_Status" json:"status,omitempty"`
	UserCount                   int64                      `protobuf:"varint,2,opt,name=user_count,json=userCount,proto3" json:"user_count,omitempty"`
	DepartmentCount             int64                      `protobuf:"varint,3,opt,name=department_count,json=departmentCount,proto3" json:"department_count,omitempty"`
	UserDepartmentRelationCount int64                      `protobuf:"varint,4,opt,name=user_department_relation_count,json=userDepartmentRelationCount,proto3" json:"user_department_relation_count,omitempty"`
	LatestSyncTime              *timestamppb.Timestamp     `protobuf:"bytes,5,opt,name=latest_sync_time,json=latestSyncTime,proto3" json:"latest_sync_time,omitempty"`
	unknownFields               protoimpl.UnknownFields
	sizeCache                   protoimpl.SizeCache
}

func (x *GetSyncAccountReply) Reset() {
	*x = GetSyncAccountReply{}
	mi := &file_api_account_v1_account_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSyncAccountReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSyncAccountReply) ProtoMessage() {}

func (x *GetSyncAccountReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSyncAccountReply.ProtoReflect.Descriptor instead.
func (*GetSyncAccountReply) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{5}
}

func (x *GetSyncAccountReply) GetStatus() GetSyncAccountReply_Status {
	if x != nil {
		return x.Status
	}
	return GetSyncAccountReply_PENDING
}

func (x *GetSyncAccountReply) GetUserCount() int64 {
	if x != nil {
		return x.UserCount
	}
	return 0
}

func (x *GetSyncAccountReply) GetDepartmentCount() int64 {
	if x != nil {
		return x.DepartmentCount
	}
	return 0
}

func (x *GetSyncAccountReply) GetUserDepartmentRelationCount() int64 {
	if x != nil {
		return x.UserDepartmentRelationCount
	}
	return 0
}

func (x *GetSyncAccountReply) GetLatestSyncTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LatestSyncTime
	}
	return nil
}

type CancelSyncAccountRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"` // 要删除的任务ID
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CancelSyncAccountRequest) Reset() {
	*x = CancelSyncAccountRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelSyncAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelSyncAccountRequest) ProtoMessage() {}

func (x *CancelSyncAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelSyncAccountRequest.ProtoReflect.Descriptor instead.
func (*CancelSyncAccountRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{6}
}

func (x *CancelSyncAccountRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

type GetAccessTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAccessTokenRequest) Reset() {
	*x = GetAccessTokenRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessTokenRequest) ProtoMessage() {}

func (x *GetAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*GetAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{7}
}

func (x *GetAccessTokenRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type GetAccessTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken  string                 `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	ExpiresIn     int64                  `protobuf:"varint,3,opt,name=expires_in,json=expiresIn,proto3" json:"expires_in,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAccessTokenResponse) Reset() {
	*x = GetAccessTokenResponse{}
	mi := &file_api_account_v1_account_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAccessTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessTokenResponse) ProtoMessage() {}

func (x *GetAccessTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessTokenResponse.ProtoReflect.Descriptor instead.
func (*GetAccessTokenResponse) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{8}
}

func (x *GetAccessTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *GetAccessTokenResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *GetAccessTokenResponse) GetExpiresIn() int64 {
	if x != nil {
		return x.ExpiresIn
	}
	return 0
}

type GetUserInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserInfoRequest) Reset() {
	*x = GetUserInfoRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoRequest) ProtoMessage() {}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoRequest.ProtoReflect.Descriptor instead.
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{9}
}

func (x *GetUserInfoRequest) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type GetUserInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UnionId       string                 `protobuf:"bytes,1,opt,name=union_id,json=unionId,proto3" json:"union_id,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email         string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Avatar        string                 `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserInfoResponse) Reset() {
	*x = GetUserInfoResponse{}
	mi := &file_api_account_v1_account_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoResponse) ProtoMessage() {}

func (x *GetUserInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoResponse.ProtoReflect.Descriptor instead.
func (*GetUserInfoResponse) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{10}
}

func (x *GetUserInfoResponse) GetUnionId() string {
	if x != nil {
		return x.UnionId
	}
	return ""
}

func (x *GetUserInfoResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetUserInfoResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetUserInfoResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GetUserInfoResponse) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

type CallbackRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`   // OAuth2 授权码
	State         string                 `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"` // 防止 CSRF 的随机字符串
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CallbackRequest) Reset() {
	*x = CallbackRequest{}
	mi := &file_api_account_v1_account_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CallbackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallbackRequest) ProtoMessage() {}

func (x *CallbackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallbackRequest.ProtoReflect.Descriptor instead.
func (*CallbackRequest) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{11}
}

func (x *CallbackRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *CallbackRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

// 定义回调响应
type CallbackResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`   // 例如 "success" 或 "error"
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // 可选描述信息
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CallbackResponse) Reset() {
	*x = CallbackResponse{}
	mi := &file_api_account_v1_account_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CallbackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallbackResponse) ProtoMessage() {}

func (x *CallbackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_v1_account_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallbackResponse.ProtoReflect.Descriptor instead.
func (*CallbackResponse) Descriptor() ([]byte, []int) {
	return file_api_account_v1_account_proto_rawDescGZIP(), []int{12}
}

func (x *CallbackResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *CallbackResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_api_account_v1_account_proto protoreflect.FileDescriptor

const file_api_account_v1_account_proto_rawDesc = "" +
	"\n" +
	"\x1capi/account/v1/account.proto\x12\x0eapi.account.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1bgoogle/protobuf/empty.proto\"I\n" +
	"\rUploadRequest\x12\x1b\n" +
	"\tfile_data\x18\x01 \x01(\fR\bfileData\x12\x1b\n" +
	"\tfile_name\x18\x02 \x01(\tR\bfileName\"E\n" +
	"\vUploadReply\x12\x19\n" +
	"\bfile_url\x18\x01 \x01(\tR\afileUrl\x12\x1b\n" +
	"\tfile_size\x18\x02 \x01(\x03R\bfileSize\"\x91\x01\n" +
	"\x18CreateSyncAccountRequest\x12>\n" +
	"\ftrigger_type\x18\x01 \x01(\x0e2\x1b.api.account.v1.TriggerTypeR\vtriggerType\x125\n" +
	"\tsync_type\x18\x02 \x01(\x0e2\x18.api.account.v1.SyncTypeR\bsyncType\"n\n" +
	"\x16CreateSyncAccountReply\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\x12;\n" +
	"\vcreate_time\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"createTime\"0\n" +
	"\x15GetSyncAccountRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\"\xeb\x02\n" +
	"\x13GetSyncAccountReply\x12B\n" +
	"\x06status\x18\x01 \x01(\x0e2*.api.account.v1.GetSyncAccountReply.StatusR\x06status\x12\x1d\n" +
	"\n" +
	"user_count\x18\x02 \x01(\x03R\tuserCount\x12)\n" +
	"\x10department_count\x18\x03 \x01(\x03R\x0fdepartmentCount\x12C\n" +
	"\x1euser_department_relation_count\x18\x04 \x01(\x03R\x1buserDepartmentRelationCount\x12D\n" +
	"\x10latest_sync_time\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\x0elatestSyncTime\";\n" +
	"\x06Status\x12\v\n" +
	"\aPENDING\x10\x00\x12\v\n" +
	"\aRUNNING\x10\x01\x12\v\n" +
	"\aSUCCESS\x10\x02\x12\n" +
	"\n" +
	"\x06FAILED\x10\x03\"3\n" +
	"\x18CancelSyncAccountRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\"+\n" +
	"\x15GetAccessTokenRequest\x12\x12\n" +
	"\x04code\x18\x02 \x01(\tR\x04code\"\x7f\n" +
	"\x16GetAccessTokenResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\x12#\n" +
	"\rrefresh_token\x18\x02 \x01(\tR\frefreshToken\x12\x1d\n" +
	"\n" +
	"expires_in\x18\x03 \x01(\x03R\texpiresIn\"7\n" +
	"\x12GetUserInfoRequest\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\"\x8b\x01\n" +
	"\x13GetUserInfoResponse\x12\x19\n" +
	"\bunion_id\x18\x01 \x01(\tR\aunionId\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x14\n" +
	"\x05email\x18\x04 \x01(\tR\x05email\x12\x16\n" +
	"\x06avatar\x18\x05 \x01(\tR\x06avatar\";\n" +
	"\x0fCallbackRequest\x12\x12\n" +
	"\x04code\x18\x01 \x01(\tR\x04code\x12\x14\n" +
	"\x05state\x18\x02 \x01(\tR\x05state\"D\n" +
	"\x10CallbackResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage*M\n" +
	"\vTriggerType\x12\x13\n" +
	"\x0fTRIGGER_UNKNOWN\x10\x00\x12\x12\n" +
	"\x0eTRIGGER_MANUAL\x10\x01\x12\x15\n" +
	"\x11TRIGGER_SCHEDULED\x10\x02*%\n" +
	"\bSyncType\x12\b\n" +
	"\x04FULL\x10\x00\x12\x0f\n" +
	"\vINCREMENTAL\x10\x012\xac\x06\n" +
	"\aAccount\x12}\n" +
	"\x11CreateSyncAccount\x12(.api.account.v1.CreateSyncAccountRequest\x1a&.api.account.v1.CreateSyncAccountReply\"\x16\x82\xd3\xe4\x93\x02\x10:\x01*\"\v/v1/account\x12q\n" +
	"\x0eGetSyncAccount\x12%.api.account.v1.GetSyncAccountRequest\x1a#.api.account.v1.GetSyncAccountReply\"\x13\x82\xd3\xe4\x93\x02\r\x12\v/v1/account\x12g\n" +
	"\x0eCancelSyncTask\x12(.api.account.v1.CancelSyncAccountRequest\x1a\x16.google.protobuf.Empty\"\x13\x82\xd3\xe4\x93\x02\r*\v/v1/account\x12u\n" +
	"\vGetUserInfo\x12\".api.account.v1.GetUserInfoRequest\x1a#.api.account.v1.GetUserInfoResponse\"\x1d\x82\xd3\xe4\x93\x02\x17\x12\x15/v1/oauth/userinfo/me\x12\x82\x01\n" +
	"\x0eGetAccessToken\x12%.api.account.v1.GetAccessTokenRequest\x1a&.api.account.v1.GetAccessTokenResponse\"!\x82\xd3\xe4\x93\x02\x1b\x12\x19/v1/oauth/userAccessToken\x12i\n" +
	"\bCallback\x12\x1f.api.account.v1.CallbackRequest\x1a .api.account.v1.CallbackResponse\"\x1a\x82\xd3\xe4\x93\x02\x14\x12\x12/v1/oauth/callback\x12_\n" +
	"\n" +
	"UploadFile\x12\x1d.api.account.v1.UploadRequest\x1a\x1b.api.account.v1.UploadReply\"\x15\x82\xd3\xe4\x93\x02\x0f:\x01*\"\n" +
	"/v1/uploadB?\n" +
	"\x0eapi.account.v1B\x0eAccountProtoV1P\x01Z\x1bnancalacc/api/account/v1;v1b\x06proto3"

var (
	file_api_account_v1_account_proto_rawDescOnce sync.Once
	file_api_account_v1_account_proto_rawDescData []byte
)

func file_api_account_v1_account_proto_rawDescGZIP() []byte {
	file_api_account_v1_account_proto_rawDescOnce.Do(func() {
		file_api_account_v1_account_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_account_v1_account_proto_rawDesc), len(file_api_account_v1_account_proto_rawDesc)))
	})
	return file_api_account_v1_account_proto_rawDescData
}

var file_api_account_v1_account_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_api_account_v1_account_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_api_account_v1_account_proto_goTypes = []any{
	(TriggerType)(0),                 // 0: api.account.v1.TriggerType
	(SyncType)(0),                    // 1: api.account.v1.SyncType
	(GetSyncAccountReply_Status)(0),  // 2: api.account.v1.GetSyncAccountReply.Status
	(*UploadRequest)(nil),            // 3: api.account.v1.UploadRequest
	(*UploadReply)(nil),              // 4: api.account.v1.UploadReply
	(*CreateSyncAccountRequest)(nil), // 5: api.account.v1.CreateSyncAccountRequest
	(*CreateSyncAccountReply)(nil),   // 6: api.account.v1.CreateSyncAccountReply
	(*GetSyncAccountRequest)(nil),    // 7: api.account.v1.GetSyncAccountRequest
	(*GetSyncAccountReply)(nil),      // 8: api.account.v1.GetSyncAccountReply
	(*CancelSyncAccountRequest)(nil), // 9: api.account.v1.CancelSyncAccountRequest
	(*GetAccessTokenRequest)(nil),    // 10: api.account.v1.GetAccessTokenRequest
	(*GetAccessTokenResponse)(nil),   // 11: api.account.v1.GetAccessTokenResponse
	(*GetUserInfoRequest)(nil),       // 12: api.account.v1.GetUserInfoRequest
	(*GetUserInfoResponse)(nil),      // 13: api.account.v1.GetUserInfoResponse
	(*CallbackRequest)(nil),          // 14: api.account.v1.CallbackRequest
	(*CallbackResponse)(nil),         // 15: api.account.v1.CallbackResponse
	(*timestamppb.Timestamp)(nil),    // 16: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),            // 17: google.protobuf.Empty
}
var file_api_account_v1_account_proto_depIdxs = []int32{
	0,  // 0: api.account.v1.CreateSyncAccountRequest.trigger_type:type_name -> api.account.v1.TriggerType
	1,  // 1: api.account.v1.CreateSyncAccountRequest.sync_type:type_name -> api.account.v1.SyncType
	16, // 2: api.account.v1.CreateSyncAccountReply.create_time:type_name -> google.protobuf.Timestamp
	2,  // 3: api.account.v1.GetSyncAccountReply.status:type_name -> api.account.v1.GetSyncAccountReply.Status
	16, // 4: api.account.v1.GetSyncAccountReply.latest_sync_time:type_name -> google.protobuf.Timestamp
	5,  // 5: api.account.v1.Account.CreateSyncAccount:input_type -> api.account.v1.CreateSyncAccountRequest
	7,  // 6: api.account.v1.Account.GetSyncAccount:input_type -> api.account.v1.GetSyncAccountRequest
	9,  // 7: api.account.v1.Account.CancelSyncTask:input_type -> api.account.v1.CancelSyncAccountRequest
	12, // 8: api.account.v1.Account.GetUserInfo:input_type -> api.account.v1.GetUserInfoRequest
	10, // 9: api.account.v1.Account.GetAccessToken:input_type -> api.account.v1.GetAccessTokenRequest
	14, // 10: api.account.v1.Account.Callback:input_type -> api.account.v1.CallbackRequest
	3,  // 11: api.account.v1.Account.UploadFile:input_type -> api.account.v1.UploadRequest
	6,  // 12: api.account.v1.Account.CreateSyncAccount:output_type -> api.account.v1.CreateSyncAccountReply
	8,  // 13: api.account.v1.Account.GetSyncAccount:output_type -> api.account.v1.GetSyncAccountReply
	17, // 14: api.account.v1.Account.CancelSyncTask:output_type -> google.protobuf.Empty
	13, // 15: api.account.v1.Account.GetUserInfo:output_type -> api.account.v1.GetUserInfoResponse
	11, // 16: api.account.v1.Account.GetAccessToken:output_type -> api.account.v1.GetAccessTokenResponse
	15, // 17: api.account.v1.Account.Callback:output_type -> api.account.v1.CallbackResponse
	4,  // 18: api.account.v1.Account.UploadFile:output_type -> api.account.v1.UploadReply
	12, // [12:19] is the sub-list for method output_type
	5,  // [5:12] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_account_v1_account_proto_init() }
func file_api_account_v1_account_proto_init() {
	if File_api_account_v1_account_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_account_v1_account_proto_rawDesc), len(file_api_account_v1_account_proto_rawDesc)),
			NumEnums:      3,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_account_v1_account_proto_goTypes,
		DependencyIndexes: file_api_account_v1_account_proto_depIdxs,
		EnumInfos:         file_api_account_v1_account_proto_enumTypes,
		MessageInfos:      file_api_account_v1_account_proto_msgTypes,
	}.Build()
	File_api_account_v1_account_proto = out.File
	file_api_account_v1_account_proto_goTypes = nil
	file_api_account_v1_account_proto_depIdxs = nil
}
