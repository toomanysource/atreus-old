// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: api/publish/service/v1/publish.proto

package v1

import (
	v1 "Atreus/api/entity/v1"
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

type PublishActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // 用户鉴权token
	Data  []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`   // 视频数据
	Title string `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"` // 视频标题
}

func (x *PublishActionRequest) Reset() {
	*x = PublishActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_publish_service_v1_publish_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishActionRequest) ProtoMessage() {}

func (x *PublishActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_publish_service_v1_publish_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishActionRequest.ProtoReflect.Descriptor instead.
func (*PublishActionRequest) Descriptor() ([]byte, []int) {
	return file_api_publish_service_v1_publish_proto_rawDescGZIP(), []int{0}
}

func (x *PublishActionRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *PublishActionRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PublishActionRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

type PublishActionReplay struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`     // 返回状态描述
}

func (x *PublishActionReplay) Reset() {
	*x = PublishActionReplay{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_publish_service_v1_publish_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishActionReplay) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishActionReplay) ProtoMessage() {}

func (x *PublishActionReplay) ProtoReflect() protoreflect.Message {
	mi := &file_api_publish_service_v1_publish_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishActionReplay.ProtoReflect.Descriptor instead.
func (*PublishActionReplay) Descriptor() ([]byte, []int) {
	return file_api_publish_service_v1_publish_proto_rawDescGZIP(), []int{1}
}

func (x *PublishActionReplay) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *PublishActionReplay) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

type PublishListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // 用户id
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`                  // 用户鉴权token
}

func (x *PublishListRequest) Reset() {
	*x = PublishListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_publish_service_v1_publish_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishListRequest) ProtoMessage() {}

func (x *PublishListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_publish_service_v1_publish_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishListRequest.ProtoReflect.Descriptor instead.
func (*PublishListRequest) Descriptor() ([]byte, []int) {
	return file_api_publish_service_v1_publish_proto_rawDescGZIP(), []int{2}
}

func (x *PublishListRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *PublishListRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type PublishListReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32     `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"` // 状态码，0-成功，其他值-失败
	StatusMsg  string    `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`     // 返回状态描述
	VideoList  *v1.Video `protobuf:"bytes,3,opt,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`     // 用户发布的视频列表
}

func (x *PublishListReply) Reset() {
	*x = PublishListReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_publish_service_v1_publish_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishListReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishListReply) ProtoMessage() {}

func (x *PublishListReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_publish_service_v1_publish_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishListReply.ProtoReflect.Descriptor instead.
func (*PublishListReply) Descriptor() ([]byte, []int) {
	return file_api_publish_service_v1_publish_proto_rawDescGZIP(), []int{3}
}

func (x *PublishListReply) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *PublishListReply) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

func (x *PublishListReply) GetVideoList() *v1.Video {
	if x != nil {
		return x.VideoList
	}
	return nil
}

var File_api_publish_service_v1_publish_proto protoreflect.FileDescriptor

var file_api_publish_service_v1_publish_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1a, 0x61, 0x70, 0x69, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x14, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73,
	0x68, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x22, 0x55,
	0x0a, 0x13, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x70, 0x6c, 0x61, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x4d, 0x73, 0x67, 0x22, 0x43, 0x0a, 0x12, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x83, 0x01, 0x0a, 0x10, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x73, 0x67, 0x12,
	0x2f, 0x0a, 0x0a, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x09, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x4c, 0x69, 0x73, 0x74,
	0x32, 0xd1, 0x01, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12, 0x60, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x26,
	0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x64,
	0x0a, 0x0d, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x28, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c,
	0x61, 0x79, 0x22, 0x00, 0x42, 0x1b, 0x5a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_publish_service_v1_publish_proto_rawDescOnce sync.Once
	file_api_publish_service_v1_publish_proto_rawDescData = file_api_publish_service_v1_publish_proto_rawDesc
)

func file_api_publish_service_v1_publish_proto_rawDescGZIP() []byte {
	file_api_publish_service_v1_publish_proto_rawDescOnce.Do(func() {
		file_api_publish_service_v1_publish_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_publish_service_v1_publish_proto_rawDescData)
	})
	return file_api_publish_service_v1_publish_proto_rawDescData
}

var file_api_publish_service_v1_publish_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_publish_service_v1_publish_proto_goTypes = []interface{}{
	(*PublishActionRequest)(nil), // 0: publish.service.v1.PublishActionRequest
	(*PublishActionReplay)(nil),  // 1: publish.service.v1.PublishActionReplay
	(*PublishListRequest)(nil),   // 2: publish.service.v1.PublishListRequest
	(*PublishListReply)(nil),     // 3: publish.service.v1.PublishListReply
	(*v1.Video)(nil),             // 4: entity.v1.Video
}
var file_api_publish_service_v1_publish_proto_depIdxs = []int32{
	4, // 0: publish.service.v1.PublishListReply.video_list:type_name -> entity.v1.Video
	2, // 1: publish.service.v1.Publish.GetPublishList:input_type -> publish.service.v1.PublishListRequest
	0, // 2: publish.service.v1.Publish.PublishAction:input_type -> publish.service.v1.PublishActionRequest
	3, // 3: publish.service.v1.Publish.GetPublishList:output_type -> publish.service.v1.PublishListReply
	1, // 4: publish.service.v1.Publish.PublishAction:output_type -> publish.service.v1.PublishActionReplay
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_publish_service_v1_publish_proto_init() }
func file_api_publish_service_v1_publish_proto_init() {
	if File_api_publish_service_v1_publish_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_publish_service_v1_publish_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishActionRequest); i {
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
		file_api_publish_service_v1_publish_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishActionReplay); i {
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
		file_api_publish_service_v1_publish_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishListRequest); i {
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
		file_api_publish_service_v1_publish_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishListReply); i {
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
			RawDescriptor: file_api_publish_service_v1_publish_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_publish_service_v1_publish_proto_goTypes,
		DependencyIndexes: file_api_publish_service_v1_publish_proto_depIdxs,
		MessageInfos:      file_api_publish_service_v1_publish_proto_msgTypes,
	}.Build()
	File_api_publish_service_v1_publish_proto = out.File
	file_api_publish_service_v1_publish_proto_rawDesc = nil
	file_api_publish_service_v1_publish_proto_goTypes = nil
	file_api_publish_service_v1_publish_proto_depIdxs = nil
}
