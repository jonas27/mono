// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.15.8
// source: vocab.proto

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

type VocabAdditional struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info string `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *VocabAdditional) Reset() {
	*x = VocabAdditional{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vocab_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VocabAdditional) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VocabAdditional) ProtoMessage() {}

func (x *VocabAdditional) ProtoReflect() protoreflect.Message {
	mi := &file_vocab_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VocabAdditional.ProtoReflect.Descriptor instead.
func (*VocabAdditional) Descriptor() ([]byte, []int) {
	return file_vocab_proto_rawDescGZIP(), []int{0}
}

func (x *VocabAdditional) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

type Vocab struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Word        string   `protobuf:"bytes,1,opt,name=word,proto3" json:"word,omitempty"`
	Description string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Translation string   `protobuf:"bytes,3,opt,name=translation,proto3" json:"translation,omitempty"`
	Info        []string `protobuf:"bytes,4,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *Vocab) Reset() {
	*x = Vocab{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vocab_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vocab) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vocab) ProtoMessage() {}

func (x *Vocab) ProtoReflect() protoreflect.Message {
	mi := &file_vocab_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vocab.ProtoReflect.Descriptor instead.
func (*Vocab) Descriptor() ([]byte, []int) {
	return file_vocab_proto_rawDescGZIP(), []int{1}
}

func (x *Vocab) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *Vocab) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Vocab) GetTranslation() string {
	if x != nil {
		return x.Translation
	}
	return ""
}

func (x *Vocab) GetInfo() []string {
	if x != nil {
		return x.Info
	}
	return nil
}

type VocabListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vocab      []*Vocab `protobuf:"bytes,1,rep,name=vocab,proto3" json:"vocab,omitempty"`
	TotalCount int32    `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
}

func (x *VocabListResponse) Reset() {
	*x = VocabListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vocab_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VocabListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VocabListResponse) ProtoMessage() {}

func (x *VocabListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vocab_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VocabListResponse.ProtoReflect.Descriptor instead.
func (*VocabListResponse) Descriptor() ([]byte, []int) {
	return file_vocab_proto_rawDescGZIP(), []int{2}
}

func (x *VocabListResponse) GetVocab() []*Vocab {
	if x != nil {
		return x.Vocab
	}
	return nil
}

func (x *VocabListResponse) GetTotalCount() int32 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

type VocabListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageNumber int32       `protobuf:"varint,1,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   int32       `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Pagination *Pagination `protobuf:"bytes,3,opt,name=pagination,proto3,oneof" json:"pagination,omitempty"`
}

func (x *VocabListRequest) Reset() {
	*x = VocabListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vocab_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VocabListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VocabListRequest) ProtoMessage() {}

func (x *VocabListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vocab_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VocabListRequest.ProtoReflect.Descriptor instead.
func (*VocabListRequest) Descriptor() ([]byte, []int) {
	return file_vocab_proto_rawDescGZIP(), []int{3}
}

func (x *VocabListRequest) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *VocabListRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *VocabListRequest) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

type Pagination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start          int32 `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	End            int32 `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
	ResultsPerPage int32 `protobuf:"varint,3,opt,name=results_per_page,json=resultsPerPage,proto3" json:"results_per_page,omitempty"`
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vocab_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_vocab_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_vocab_proto_rawDescGZIP(), []int{4}
}

func (x *Pagination) GetStart() int32 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *Pagination) GetEnd() int32 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *Pagination) GetResultsPerPage() int32 {
	if x != nil {
		return x.ResultsPerPage
	}
	return 0
}

var File_vocab_proto protoreflect.FileDescriptor

var file_vocab_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x76, 0x6f, 0x63, 0x61, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x76,
	0x6f, 0x63, 0x61, 0x62, 0x22, 0x25, 0x0a, 0x0f, 0x56, 0x6f, 0x63, 0x61, 0x62, 0x41, 0x64, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x73, 0x0a, 0x05, 0x56,
	0x6f, 0x63, 0x61, 0x62, 0x12, 0x12, 0x0a, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f,
	0x22, 0x58, 0x0a, 0x11, 0x56, 0x6f, 0x63, 0x61, 0x62, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x05, 0x76, 0x6f, 0x63, 0x61, 0x62, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x76, 0x6f, 0x63, 0x61, 0x62, 0x2e, 0x56, 0x6f, 0x63,
	0x61, 0x62, 0x52, 0x05, 0x76, 0x6f, 0x63, 0x61, 0x62, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x97, 0x01, 0x0a, 0x10, 0x56,
	0x6f, 0x63, 0x61, 0x62, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x36, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x76, 0x6f, 0x63, 0x61, 0x62, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x5e, 0x0a, 0x0a, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x50, 0x65, 0x72,
	0x50, 0x61, 0x67, 0x65, 0x32, 0x51, 0x0a, 0x0c, 0x56, 0x6f, 0x63, 0x61, 0x62, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x6f, 0x63, 0x61,
	0x62, 0x73, 0x12, 0x17, 0x2e, 0x76, 0x6f, 0x63, 0x61, 0x62, 0x2e, 0x56, 0x6f, 0x63, 0x61, 0x62,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x76, 0x6f,
	0x63, 0x61, 0x62, 0x2e, 0x56, 0x6f, 0x63, 0x61, 0x62, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vocab_proto_rawDescOnce sync.Once
	file_vocab_proto_rawDescData = file_vocab_proto_rawDesc
)

func file_vocab_proto_rawDescGZIP() []byte {
	file_vocab_proto_rawDescOnce.Do(func() {
		file_vocab_proto_rawDescData = protoimpl.X.CompressGZIP(file_vocab_proto_rawDescData)
	})
	return file_vocab_proto_rawDescData
}

var file_vocab_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_vocab_proto_goTypes = []interface{}{
	(*VocabAdditional)(nil),   // 0: vocab.VocabAdditional
	(*Vocab)(nil),             // 1: vocab.Vocab
	(*VocabListResponse)(nil), // 2: vocab.VocabListResponse
	(*VocabListRequest)(nil),  // 3: vocab.VocabListRequest
	(*Pagination)(nil),        // 4: vocab.Pagination
}
var file_vocab_proto_depIdxs = []int32{
	1, // 0: vocab.VocabListResponse.vocab:type_name -> vocab.Vocab
	4, // 1: vocab.VocabListRequest.pagination:type_name -> vocab.Pagination
	3, // 2: vocab.VocabService.ListVocabs:input_type -> vocab.VocabListRequest
	2, // 3: vocab.VocabService.ListVocabs:output_type -> vocab.VocabListResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_vocab_proto_init() }
func file_vocab_proto_init() {
	if File_vocab_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vocab_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VocabAdditional); i {
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
		file_vocab_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vocab); i {
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
		file_vocab_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VocabListResponse); i {
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
		file_vocab_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VocabListRequest); i {
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
		file_vocab_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pagination); i {
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
	file_vocab_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vocab_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vocab_proto_goTypes,
		DependencyIndexes: file_vocab_proto_depIdxs,
		MessageInfos:      file_vocab_proto_msgTypes,
	}.Build()
	File_vocab_proto = out.File
	file_vocab_proto_rawDesc = nil
	file_vocab_proto_goTypes = nil
	file_vocab_proto_depIdxs = nil
}
