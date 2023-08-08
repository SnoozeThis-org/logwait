// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: observer.proto

package logwait

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

type RegisterObserverRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogInstanceToken string `protobuf:"bytes,1,opt,name=log_instance_token,json=logInstanceToken,proto3" json:"log_instance_token,omitempty"`
}

func (x *RegisterObserverRequest) Reset() {
	*x = RegisterObserverRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_observer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterObserverRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterObserverRequest) ProtoMessage() {}

func (x *RegisterObserverRequest) ProtoReflect() protoreflect.Message {
	mi := &file_observer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterObserverRequest.ProtoReflect.Descriptor instead.
func (*RegisterObserverRequest) Descriptor() ([]byte, []int) {
	return file_observer_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterObserverRequest) GetLogInstanceToken() string {
	if x != nil {
		return x.LogInstanceToken
	}
	return ""
}

type RegisterObserverResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActiveObservables []*Observable `protobuf:"bytes,1,rep,name=active_observables,json=activeObservables,proto3" json:"active_observables,omitempty"`
	// The url you've registered with SnoozeThis. The Observer should use this URL as the base for filter-URLs to ensure the domain is understood by SnoozeThis.
	RegisteredUrl string `protobuf:"bytes,2,opt,name=registered_url,json=registeredUrl,proto3" json:"registered_url,omitempty"`
}

func (x *RegisterObserverResponse) Reset() {
	*x = RegisterObserverResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_observer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterObserverResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterObserverResponse) ProtoMessage() {}

func (x *RegisterObserverResponse) ProtoReflect() protoreflect.Message {
	mi := &file_observer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterObserverResponse.ProtoReflect.Descriptor instead.
func (*RegisterObserverResponse) Descriptor() ([]byte, []int) {
	return file_observer_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterObserverResponse) GetActiveObservables() []*Observable {
	if x != nil {
		return x.ActiveObservables
	}
	return nil
}

func (x *RegisterObserverResponse) GetRegisteredUrl() string {
	if x != nil {
		return x.RegisteredUrl
	}
	return ""
}

type RejectObservable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RejectObservable) Reset() {
	*x = RejectObservable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_observer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RejectObservable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RejectObservable) ProtoMessage() {}

func (x *RejectObservable) ProtoReflect() protoreflect.Message {
	mi := &file_observer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RejectObservable.ProtoReflect.Descriptor instead.
func (*RejectObservable) Descriptor() ([]byte, []int) {
	return file_observer_proto_rawDescGZIP(), []int{2}
}

func (x *RejectObservable) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RejectObservable) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ObserverToSnoozeThis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Msg:
	//	*ObserverToSnoozeThis_Register
	//	*ObserverToSnoozeThis_ObservedObservable
	//	*ObserverToSnoozeThis_RejectObservable
	Msg isObserverToSnoozeThis_Msg `protobuf_oneof:"msg"`
}

func (x *ObserverToSnoozeThis) Reset() {
	*x = ObserverToSnoozeThis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_observer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ObserverToSnoozeThis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObserverToSnoozeThis) ProtoMessage() {}

func (x *ObserverToSnoozeThis) ProtoReflect() protoreflect.Message {
	mi := &file_observer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObserverToSnoozeThis.ProtoReflect.Descriptor instead.
func (*ObserverToSnoozeThis) Descriptor() ([]byte, []int) {
	return file_observer_proto_rawDescGZIP(), []int{3}
}

func (m *ObserverToSnoozeThis) GetMsg() isObserverToSnoozeThis_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *ObserverToSnoozeThis) GetRegister() *RegisterObserverRequest {
	if x, ok := x.GetMsg().(*ObserverToSnoozeThis_Register); ok {
		return x.Register
	}
	return nil
}

func (x *ObserverToSnoozeThis) GetObservedObservable() string {
	if x, ok := x.GetMsg().(*ObserverToSnoozeThis_ObservedObservable); ok {
		return x.ObservedObservable
	}
	return ""
}

func (x *ObserverToSnoozeThis) GetRejectObservable() *RejectObservable {
	if x, ok := x.GetMsg().(*ObserverToSnoozeThis_RejectObservable); ok {
		return x.RejectObservable
	}
	return nil
}

type isObserverToSnoozeThis_Msg interface {
	isObserverToSnoozeThis_Msg()
}

type ObserverToSnoozeThis_Register struct {
	Register *RegisterObserverRequest `protobuf:"bytes,1,opt,name=register,proto3,oneof"`
}

type ObserverToSnoozeThis_ObservedObservable struct {
	ObservedObservable string `protobuf:"bytes,2,opt,name=observed_observable,json=observedObservable,proto3,oneof"`
}

type ObserverToSnoozeThis_RejectObservable struct {
	RejectObservable *RejectObservable `protobuf:"bytes,3,opt,name=reject_observable,json=rejectObservable,proto3,oneof"`
}

func (*ObserverToSnoozeThis_Register) isObserverToSnoozeThis_Msg() {}

func (*ObserverToSnoozeThis_ObservedObservable) isObserverToSnoozeThis_Msg() {}

func (*ObserverToSnoozeThis_RejectObservable) isObserverToSnoozeThis_Msg() {}

type SnoozeThisToObserver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Msg:
	//	*SnoozeThisToObserver_Register
	//	*SnoozeThisToObserver_NewObservable
	//	*SnoozeThisToObserver_CancelObservable
	Msg isSnoozeThisToObserver_Msg `protobuf_oneof:"msg"`
}

func (x *SnoozeThisToObserver) Reset() {
	*x = SnoozeThisToObserver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_observer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SnoozeThisToObserver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SnoozeThisToObserver) ProtoMessage() {}

func (x *SnoozeThisToObserver) ProtoReflect() protoreflect.Message {
	mi := &file_observer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SnoozeThisToObserver.ProtoReflect.Descriptor instead.
func (*SnoozeThisToObserver) Descriptor() ([]byte, []int) {
	return file_observer_proto_rawDescGZIP(), []int{4}
}

func (m *SnoozeThisToObserver) GetMsg() isSnoozeThisToObserver_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *SnoozeThisToObserver) GetRegister() *RegisterObserverResponse {
	if x, ok := x.GetMsg().(*SnoozeThisToObserver_Register); ok {
		return x.Register
	}
	return nil
}

func (x *SnoozeThisToObserver) GetNewObservable() *Observable {
	if x, ok := x.GetMsg().(*SnoozeThisToObserver_NewObservable); ok {
		return x.NewObservable
	}
	return nil
}

func (x *SnoozeThisToObserver) GetCancelObservable() string {
	if x, ok := x.GetMsg().(*SnoozeThisToObserver_CancelObservable); ok {
		return x.CancelObservable
	}
	return ""
}

type isSnoozeThisToObserver_Msg interface {
	isSnoozeThisToObserver_Msg()
}

type SnoozeThisToObserver_Register struct {
	Register *RegisterObserverResponse `protobuf:"bytes,1,opt,name=register,proto3,oneof"`
}

type SnoozeThisToObserver_NewObservable struct {
	NewObservable *Observable `protobuf:"bytes,2,opt,name=new_observable,json=newObservable,proto3,oneof"`
}

type SnoozeThisToObserver_CancelObservable struct {
	CancelObservable string `protobuf:"bytes,3,opt,name=cancel_observable,json=cancelObservable,proto3,oneof"`
}

func (*SnoozeThisToObserver_Register) isSnoozeThisToObserver_Msg() {}

func (*SnoozeThisToObserver_NewObservable) isSnoozeThisToObserver_Msg() {}

func (*SnoozeThisToObserver_CancelObservable) isSnoozeThisToObserver_Msg() {}

var File_observer_proto protoreflect.FileDescriptor

var file_observer_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68, 0x69, 0x73,
	0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x17, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2c, 0x0a, 0x12, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x6c,
	0x6f, 0x67, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x94, 0x01, 0x0a, 0x18, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x12,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73,
	0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69,
	0x74, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x11, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12,
	0x25, 0x0a, 0x0e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x65, 0x64, 0x55, 0x72, 0x6c, 0x22, 0x3c, 0x0a, 0x10, 0x52, 0x65, 0x6a, 0x65, 0x63, 0x74,
	0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xf8, 0x01, 0x0a, 0x14, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x54, 0x6f, 0x53, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x54, 0x68, 0x69, 0x73, 0x12, 0x4d, 0x0a,
	0x08, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68, 0x69, 0x73,
	0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x13,
	0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x12, 0x6f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x57, 0x0a, 0x11, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77,
	0x61, 0x69, 0x74, 0x2e, 0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x62, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x10, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x4f, 0x62,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x22,
	0xe9, 0x01, 0x0a, 0x14, 0x53, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x54, 0x68, 0x69, 0x73, 0x54, 0x6f,
	0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x4e, 0x0a, 0x08, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77,
	0x61, 0x69, 0x74, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x08,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x4b, 0x0a, 0x0e, 0x6e, 0x65, 0x77, 0x5f,
	0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68, 0x69,
	0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x62, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x6e, 0x65, 0x77, 0x4f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x2d, 0x0a, 0x11, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x5f,
	0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x10, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x62, 0x6c, 0x65, 0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x32, 0x87, 0x01, 0x0a, 0x14,
	0x53, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x54, 0x68, 0x69, 0x73, 0x4c, 0x6f, 0x67, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x6f, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x12, 0x2c, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65,
	0x74, 0x68, 0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x2e, 0x4f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x6f, 0x53, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x54, 0x68, 0x69,
	0x73, 0x1a, 0x2c, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x74, 0x68,
	0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x2e, 0x53, 0x6e, 0x6f, 0x6f, 0x7a,
	0x65, 0x54, 0x68, 0x69, 0x73, 0x54, 0x6f, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22,
	0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x54, 0x68, 0x69, 0x73, 0x2d, 0x6f,
	0x72, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_observer_proto_rawDescOnce sync.Once
	file_observer_proto_rawDescData = file_observer_proto_rawDesc
)

func file_observer_proto_rawDescGZIP() []byte {
	file_observer_proto_rawDescOnce.Do(func() {
		file_observer_proto_rawDescData = protoimpl.X.CompressGZIP(file_observer_proto_rawDescData)
	})
	return file_observer_proto_rawDescData
}

var file_observer_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_observer_proto_goTypes = []interface{}{
	(*RegisterObserverRequest)(nil),  // 0: com.snoozethis.logwait.RegisterObserverRequest
	(*RegisterObserverResponse)(nil), // 1: com.snoozethis.logwait.RegisterObserverResponse
	(*RejectObservable)(nil),         // 2: com.snoozethis.logwait.RejectObservable
	(*ObserverToSnoozeThis)(nil),     // 3: com.snoozethis.logwait.ObserverToSnoozeThis
	(*SnoozeThisToObserver)(nil),     // 4: com.snoozethis.logwait.SnoozeThisToObserver
	(*Observable)(nil),               // 5: com.snoozethis.logwait.Observable
}
var file_observer_proto_depIdxs = []int32{
	5, // 0: com.snoozethis.logwait.RegisterObserverResponse.active_observables:type_name -> com.snoozethis.logwait.Observable
	0, // 1: com.snoozethis.logwait.ObserverToSnoozeThis.register:type_name -> com.snoozethis.logwait.RegisterObserverRequest
	2, // 2: com.snoozethis.logwait.ObserverToSnoozeThis.reject_observable:type_name -> com.snoozethis.logwait.RejectObservable
	1, // 3: com.snoozethis.logwait.SnoozeThisToObserver.register:type_name -> com.snoozethis.logwait.RegisterObserverResponse
	5, // 4: com.snoozethis.logwait.SnoozeThisToObserver.new_observable:type_name -> com.snoozethis.logwait.Observable
	3, // 5: com.snoozethis.logwait.SnoozeThisLogService.Communicate:input_type -> com.snoozethis.logwait.ObserverToSnoozeThis
	4, // 6: com.snoozethis.logwait.SnoozeThisLogService.Communicate:output_type -> com.snoozethis.logwait.SnoozeThisToObserver
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_observer_proto_init() }
func file_observer_proto_init() {
	if File_observer_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_observer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterObserverRequest); i {
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
		file_observer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterObserverResponse); i {
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
		file_observer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RejectObservable); i {
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
		file_observer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ObserverToSnoozeThis); i {
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
		file_observer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SnoozeThisToObserver); i {
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
	file_observer_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*ObserverToSnoozeThis_Register)(nil),
		(*ObserverToSnoozeThis_ObservedObservable)(nil),
		(*ObserverToSnoozeThis_RejectObservable)(nil),
	}
	file_observer_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*SnoozeThisToObserver_Register)(nil),
		(*SnoozeThisToObserver_NewObservable)(nil),
		(*SnoozeThisToObserver_CancelObservable)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_observer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_observer_proto_goTypes,
		DependencyIndexes: file_observer_proto_depIdxs,
		MessageInfos:      file_observer_proto_msgTypes,
	}.Build()
	File_observer_proto = out.File
	file_observer_proto_rawDesc = nil
	file_observer_proto_goTypes = nil
	file_observer_proto_depIdxs = nil
}
