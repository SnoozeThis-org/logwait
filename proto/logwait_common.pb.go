// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: logwait_common.proto

package proto

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

type Filter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field  string `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Regexp string `protobuf:"bytes,2,opt,name=regexp,proto3" json:"regexp,omitempty"`
}

func (x *Filter) Reset() {
	*x = Filter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logwait_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Filter) ProtoMessage() {}

func (x *Filter) ProtoReflect() protoreflect.Message {
	mi := &file_logwait_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Filter.ProtoReflect.Descriptor instead.
func (*Filter) Descriptor() ([]byte, []int) {
	return file_logwait_common_proto_rawDescGZIP(), []int{0}
}

func (x *Filter) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *Filter) GetRegexp() string {
	if x != nil {
		return x.Regexp
	}
	return ""
}

type Observable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Signature string    `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	Filters   []*Filter `protobuf:"bytes,3,rep,name=filters,proto3" json:"filters,omitempty"`
}

func (x *Observable) Reset() {
	*x = Observable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logwait_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Observable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Observable) ProtoMessage() {}

func (x *Observable) ProtoReflect() protoreflect.Message {
	mi := &file_logwait_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Observable.ProtoReflect.Descriptor instead.
func (*Observable) Descriptor() ([]byte, []int) {
	return file_logwait_common_proto_rawDescGZIP(), []int{1}
}

func (x *Observable) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Observable) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *Observable) GetFilters() []*Filter {
	if x != nil {
		return x.Filters
	}
	return nil
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
		mi := &file_logwait_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RejectObservable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RejectObservable) ProtoMessage() {}

func (x *RejectObservable) ProtoReflect() protoreflect.Message {
	mi := &file_logwait_common_proto_msgTypes[2]
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
	return file_logwait_common_proto_rawDescGZIP(), []int{2}
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

var File_logwait_common_proto protoreflect.FileDescriptor

var file_logwait_common_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f,
	0x7a, 0x65, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x22, 0x36,
	0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x67, 0x65, 0x78, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x67, 0x65, 0x78, 0x70, 0x22, 0x74, 0x0a, 0x0a, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x62, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6e, 0x6f, 0x6f, 0x7a, 0x65,
	0x74, 0x68, 0x69, 0x73, 0x2e, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x2e, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x22, 0x3c, 0x0a, 0x10,
	0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x6e, 0x6f, 0x6f, 0x7a, 0x65, 0x54,
	0x68, 0x69, 0x73, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x77, 0x61, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_logwait_common_proto_rawDescOnce sync.Once
	file_logwait_common_proto_rawDescData = file_logwait_common_proto_rawDesc
)

func file_logwait_common_proto_rawDescGZIP() []byte {
	file_logwait_common_proto_rawDescOnce.Do(func() {
		file_logwait_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_logwait_common_proto_rawDescData)
	})
	return file_logwait_common_proto_rawDescData
}

var file_logwait_common_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_logwait_common_proto_goTypes = []interface{}{
	(*Filter)(nil),           // 0: com.snoozethis.logwait.Filter
	(*Observable)(nil),       // 1: com.snoozethis.logwait.Observable
	(*RejectObservable)(nil), // 2: com.snoozethis.logwait.RejectObservable
}
var file_logwait_common_proto_depIdxs = []int32{
	0, // 0: com.snoozethis.logwait.Observable.filters:type_name -> com.snoozethis.logwait.Filter
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_logwait_common_proto_init() }
func file_logwait_common_proto_init() {
	if File_logwait_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_logwait_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Filter); i {
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
		file_logwait_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Observable); i {
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
		file_logwait_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_logwait_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_logwait_common_proto_goTypes,
		DependencyIndexes: file_logwait_common_proto_depIdxs,
		MessageInfos:      file_logwait_common_proto_msgTypes,
	}.Build()
	File_logwait_common_proto = out.File
	file_logwait_common_proto_rawDesc = nil
	file_logwait_common_proto_goTypes = nil
	file_logwait_common_proto_depIdxs = nil
}
