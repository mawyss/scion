// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.4
// source: proto/control_plane/v1/seg_extensions.proto

package control_plane

import (
	proto "github.com/golang/protobuf/proto"
	experimental "github.com/scionproto/scion/go/pkg/proto/control_plane/experimental"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type LinkType int32

const (
	LinkType_LINK_TYPE_UNSPECIFIED LinkType = 0
	LinkType_LINK_TYPE_DIRECT      LinkType = 1
	LinkType_LINK_TYPE_MULTI_HOP   LinkType = 2
	LinkType_LINK_TYPE_OPEN_NET    LinkType = 3
)

// Enum value maps for LinkType.
var (
	LinkType_name = map[int32]string{
		0: "LINK_TYPE_UNSPECIFIED",
		1: "LINK_TYPE_DIRECT",
		2: "LINK_TYPE_MULTI_HOP",
		3: "LINK_TYPE_OPEN_NET",
	}
	LinkType_value = map[string]int32{
		"LINK_TYPE_UNSPECIFIED": 0,
		"LINK_TYPE_DIRECT":      1,
		"LINK_TYPE_MULTI_HOP":   2,
		"LINK_TYPE_OPEN_NET":    3,
	}
)

func (x LinkType) Enum() *LinkType {
	p := new(LinkType)
	*p = x
	return p
}

func (x LinkType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LinkType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_control_plane_v1_seg_extensions_proto_enumTypes[0].Descriptor()
}

func (LinkType) Type() protoreflect.EnumType {
	return &file_proto_control_plane_v1_seg_extensions_proto_enumTypes[0]
}

func (x LinkType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LinkType.Descriptor instead.
func (LinkType) EnumDescriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{0}
}

type PathSegmentExtensions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StaticInfo *StaticInfoExtension `protobuf:"bytes,1,opt,name=static_info,json=staticInfo,proto3" json:"static_info,omitempty"`
	HiddenPath *HiddenPathExtension `protobuf:"bytes,2,opt,name=hidden_path,json=hiddenPath,proto3" json:"hidden_path,omitempty"`
	Digests    *DigestExtension     `protobuf:"bytes,1000,opt,name=digests,proto3" json:"digests,omitempty"`
}

func (x *PathSegmentExtensions) Reset() {
	*x = PathSegmentExtensions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PathSegmentExtensions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PathSegmentExtensions) ProtoMessage() {}

func (x *PathSegmentExtensions) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PathSegmentExtensions.ProtoReflect.Descriptor instead.
func (*PathSegmentExtensions) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{0}
}

func (x *PathSegmentExtensions) GetStaticInfo() *StaticInfoExtension {
	if x != nil {
		return x.StaticInfo
	}
	return nil
}

func (x *PathSegmentExtensions) GetHiddenPath() *HiddenPathExtension {
	if x != nil {
		return x.HiddenPath
	}
	return nil
}

func (x *PathSegmentExtensions) GetDigests() *DigestExtension {
	if x != nil {
		return x.Digests
	}
	return nil
}

type HiddenPathExtension struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsHidden bool `protobuf:"varint,1,opt,name=is_hidden,json=isHidden,proto3" json:"is_hidden,omitempty"`
}

func (x *HiddenPathExtension) Reset() {
	*x = HiddenPathExtension{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HiddenPathExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HiddenPathExtension) ProtoMessage() {}

func (x *HiddenPathExtension) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HiddenPathExtension.ProtoReflect.Descriptor instead.
func (*HiddenPathExtension) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{1}
}

func (x *HiddenPathExtension) GetIsHidden() bool {
	if x != nil {
		return x.IsHidden
	}
	return false
}

type StaticInfoExtension struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latency      *LatencyInfo               `protobuf:"bytes,1,opt,name=latency,proto3" json:"latency,omitempty"`
	Bandwidth    *BandwidthInfo             `protobuf:"bytes,2,opt,name=bandwidth,proto3" json:"bandwidth,omitempty"`
	Geo          map[uint64]*GeoCoordinates `protobuf:"bytes,3,rep,name=geo,proto3" json:"geo,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	LinkType     map[uint64]LinkType        `protobuf:"bytes,4,rep,name=link_type,json=linkType,proto3" json:"link_type,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3,enum=proto.control_plane.v1.LinkType"`
	InternalHops map[uint64]uint32          `protobuf:"bytes,5,rep,name=internal_hops,json=internalHops,proto3" json:"internal_hops,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Note         string                     `protobuf:"bytes,6,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *StaticInfoExtension) Reset() {
	*x = StaticInfoExtension{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StaticInfoExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticInfoExtension) ProtoMessage() {}

func (x *StaticInfoExtension) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticInfoExtension.ProtoReflect.Descriptor instead.
func (*StaticInfoExtension) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{2}
}

func (x *StaticInfoExtension) GetLatency() *LatencyInfo {
	if x != nil {
		return x.Latency
	}
	return nil
}

func (x *StaticInfoExtension) GetBandwidth() *BandwidthInfo {
	if x != nil {
		return x.Bandwidth
	}
	return nil
}

func (x *StaticInfoExtension) GetGeo() map[uint64]*GeoCoordinates {
	if x != nil {
		return x.Geo
	}
	return nil
}

func (x *StaticInfoExtension) GetLinkType() map[uint64]LinkType {
	if x != nil {
		return x.LinkType
	}
	return nil
}

func (x *StaticInfoExtension) GetInternalHops() map[uint64]uint32 {
	if x != nil {
		return x.InternalHops
	}
	return nil
}

func (x *StaticInfoExtension) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

type LatencyInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Intra map[uint64]uint32 `protobuf:"bytes,1,rep,name=intra,proto3" json:"intra,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Inter map[uint64]uint32 `protobuf:"bytes,2,rep,name=inter,proto3" json:"inter,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *LatencyInfo) Reset() {
	*x = LatencyInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LatencyInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatencyInfo) ProtoMessage() {}

func (x *LatencyInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatencyInfo.ProtoReflect.Descriptor instead.
func (*LatencyInfo) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{3}
}

func (x *LatencyInfo) GetIntra() map[uint64]uint32 {
	if x != nil {
		return x.Intra
	}
	return nil
}

func (x *LatencyInfo) GetInter() map[uint64]uint32 {
	if x != nil {
		return x.Inter
	}
	return nil
}

type BandwidthInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Intra map[uint64]uint64 `protobuf:"bytes,1,rep,name=intra,proto3" json:"intra,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Inter map[uint64]uint64 `protobuf:"bytes,2,rep,name=inter,proto3" json:"inter,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *BandwidthInfo) Reset() {
	*x = BandwidthInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BandwidthInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BandwidthInfo) ProtoMessage() {}

func (x *BandwidthInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BandwidthInfo.ProtoReflect.Descriptor instead.
func (*BandwidthInfo) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{4}
}

func (x *BandwidthInfo) GetIntra() map[uint64]uint64 {
	if x != nil {
		return x.Intra
	}
	return nil
}

func (x *BandwidthInfo) GetInter() map[uint64]uint64 {
	if x != nil {
		return x.Inter
	}
	return nil
}

type GeoCoordinates struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float32 `protobuf:"fixed32,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float32 `protobuf:"fixed32,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Address   string  `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *GeoCoordinates) Reset() {
	*x = GeoCoordinates{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeoCoordinates) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeoCoordinates) ProtoMessage() {}

func (x *GeoCoordinates) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeoCoordinates.ProtoReflect.Descriptor instead.
func (*GeoCoordinates) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{5}
}

func (x *GeoCoordinates) GetLatitude() float32 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *GeoCoordinates) GetLongitude() float32 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *GeoCoordinates) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type DigestExtension struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epic *DigestExtension_Digest `protobuf:"bytes,1000,opt,name=epic,proto3" json:"epic,omitempty"`
}

func (x *DigestExtension) Reset() {
	*x = DigestExtension{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DigestExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DigestExtension) ProtoMessage() {}

func (x *DigestExtension) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DigestExtension.ProtoReflect.Descriptor instead.
func (*DigestExtension) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{6}
}

func (x *DigestExtension) GetEpic() *DigestExtension_Digest {
	if x != nil {
		return x.Epic
	}
	return nil
}

type PathSegmentUnsignedExtensions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epic *experimental.EPICDetachedExtension `protobuf:"bytes,1000,opt,name=epic,proto3" json:"epic,omitempty"`
}

func (x *PathSegmentUnsignedExtensions) Reset() {
	*x = PathSegmentUnsignedExtensions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PathSegmentUnsignedExtensions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PathSegmentUnsignedExtensions) ProtoMessage() {}

func (x *PathSegmentUnsignedExtensions) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PathSegmentUnsignedExtensions.ProtoReflect.Descriptor instead.
func (*PathSegmentUnsignedExtensions) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{7}
}

func (x *PathSegmentUnsignedExtensions) GetEpic() *experimental.EPICDetachedExtension {
	if x != nil {
		return x.Epic
	}
	return nil
}

type DigestExtension_Digest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digest []byte `protobuf:"bytes,1,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *DigestExtension_Digest) Reset() {
	*x = DigestExtension_Digest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DigestExtension_Digest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DigestExtension_Digest) ProtoMessage() {}

func (x *DigestExtension_Digest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_control_plane_v1_seg_extensions_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DigestExtension_Digest.ProtoReflect.Descriptor instead.
func (*DigestExtension_Digest) Descriptor() ([]byte, []int) {
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP(), []int{6, 0}
}

func (x *DigestExtension_Digest) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

var File_proto_control_plane_v1_seg_extensions_proto protoreflect.FileDescriptor

var file_proto_control_plane_v1_seg_extensions_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x67, 0x5f, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x41, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x78, 0x70, 0x65, 0x72,
	0x69, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x67, 0x5f, 0x64,
	0x65, 0x74, 0x61, 0x63, 0x68, 0x65, 0x64, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf7, 0x01, 0x0a, 0x15, 0x50, 0x61, 0x74,
	0x68, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x4c, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x4c, 0x0a, 0x0b, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48,
	0x69, 0x64, 0x64, 0x65, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x0a, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x12, 0x42,
	0x0a, 0x07, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x73, 0x18, 0xe8, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x64, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x73, 0x22, 0x32, 0x0a, 0x13, 0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x50, 0x61, 0x74, 0x68,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f,
	0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73,
	0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x22, 0xb1, 0x05, 0x0a, 0x13, 0x53, 0x74, 0x61, 0x74, 0x69,
	0x63, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3d,
	0x0a, 0x07, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x43, 0x0a,
	0x09, 0x62, 0x61, 0x6e, 0x64, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x6e, 0x64, 0x77, 0x69,
	0x64, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x62, 0x61, 0x6e, 0x64, 0x77, 0x69, 0x64,
	0x74, 0x68, 0x12, 0x46, 0x0a, 0x03, 0x67, 0x65, 0x6f, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x34, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49,
	0x6e, 0x66, 0x6f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x6f,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x67, 0x65, 0x6f, 0x12, 0x56, 0x0a, 0x09, 0x6c, 0x69,
	0x6e, 0x6b, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x39, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49, 0x6e, 0x66,
	0x6f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x54,
	0x79, 0x70, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6c, 0x69, 0x6e, 0x6b, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x62, 0x0a, 0x0d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x68,
	0x6f, 0x70, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x48,
	0x6f, 0x70, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x48, 0x6f, 0x70, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x1a, 0x5e, 0x0a, 0x08, 0x47, 0x65,
	0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x6f, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x5d, 0x0a, 0x0d, 0x4c, 0x69,
	0x6e, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x36, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3f, 0x0a, 0x11, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x48, 0x6f, 0x70, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8d, 0x02, 0x0a, 0x0b, 0x4c,
	0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x44, 0x0a, 0x05, 0x69, 0x6e,
	0x74, 0x72, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x49,
	0x6e, 0x74, 0x72, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x6e, 0x74, 0x72, 0x61,
	0x12, 0x44, 0x0a, 0x05, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79,
	0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x05, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x1a, 0x38, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x72, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0x38, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x93, 0x02, 0x0a, 0x0d, 0x42,
	0x61, 0x6e, 0x64, 0x77, 0x69, 0x64, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x46, 0x0a, 0x05,
	0x69, 0x6e, 0x74, 0x72, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x6e, 0x64, 0x77, 0x69, 0x64, 0x74, 0x68, 0x49, 0x6e,
	0x66, 0x6f, 0x2e, 0x49, 0x6e, 0x74, 0x72, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69,
	0x6e, 0x74, 0x72, 0x61, 0x12, 0x46, 0x0a, 0x05, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x6e,
	0x64, 0x77, 0x69, 0x64, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x1a, 0x38, 0x0a, 0x0a,
	0x49, 0x6e, 0x74, 0x72, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x38, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x64, 0x0a, 0x0e, 0x47, 0x65, 0x6f, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x78, 0x0a, 0x0f, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x43, 0x0a, 0x04, 0x65, 0x70, 0x69,
	0x63, 0x18, 0xe8, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x52, 0x04, 0x65, 0x70, 0x69, 0x63, 0x1a, 0x20,
	0x0a, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74,
	0x22, 0x70, 0x0a, 0x1d, 0x50, 0x61, 0x74, 0x68, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x55,
	0x6e, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x4f, 0x0a, 0x04, 0x65, 0x70, 0x69, 0x63, 0x18, 0xe8, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x3a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e,
	0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x50, 0x49, 0x43, 0x44, 0x65, 0x74, 0x61, 0x63,
	0x68, 0x65, 0x64, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x65, 0x70,
	0x69, 0x63, 0x2a, 0x6c, 0x0a, 0x08, 0x4c, 0x69, 0x6e, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19,
	0x0a, 0x15, 0x4c, 0x49, 0x4e, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x4c, 0x49, 0x4e,
	0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x10, 0x01, 0x12,
	0x17, 0x0a, 0x13, 0x4c, 0x49, 0x4e, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x55, 0x4c,
	0x54, 0x49, 0x5f, 0x48, 0x4f, 0x50, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x4c, 0x49, 0x4e, 0x4b,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x50, 0x45, 0x4e, 0x5f, 0x4e, 0x45, 0x54, 0x10, 0x03,
	0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x63, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x63, 0x69, 0x6f, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_control_plane_v1_seg_extensions_proto_rawDescOnce sync.Once
	file_proto_control_plane_v1_seg_extensions_proto_rawDescData = file_proto_control_plane_v1_seg_extensions_proto_rawDesc
)

func file_proto_control_plane_v1_seg_extensions_proto_rawDescGZIP() []byte {
	file_proto_control_plane_v1_seg_extensions_proto_rawDescOnce.Do(func() {
		file_proto_control_plane_v1_seg_extensions_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_control_plane_v1_seg_extensions_proto_rawDescData)
	})
	return file_proto_control_plane_v1_seg_extensions_proto_rawDescData
}

var file_proto_control_plane_v1_seg_extensions_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_control_plane_v1_seg_extensions_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_proto_control_plane_v1_seg_extensions_proto_goTypes = []interface{}{
	(LinkType)(0),                         // 0: proto.control_plane.v1.LinkType
	(*PathSegmentExtensions)(nil),         // 1: proto.control_plane.v1.PathSegmentExtensions
	(*HiddenPathExtension)(nil),           // 2: proto.control_plane.v1.HiddenPathExtension
	(*StaticInfoExtension)(nil),           // 3: proto.control_plane.v1.StaticInfoExtension
	(*LatencyInfo)(nil),                   // 4: proto.control_plane.v1.LatencyInfo
	(*BandwidthInfo)(nil),                 // 5: proto.control_plane.v1.BandwidthInfo
	(*GeoCoordinates)(nil),                // 6: proto.control_plane.v1.GeoCoordinates
	(*DigestExtension)(nil),               // 7: proto.control_plane.v1.DigestExtension
	(*PathSegmentUnsignedExtensions)(nil), // 8: proto.control_plane.v1.PathSegmentUnsignedExtensions
	nil,                                   // 9: proto.control_plane.v1.StaticInfoExtension.GeoEntry
	nil,                                   // 10: proto.control_plane.v1.StaticInfoExtension.LinkTypeEntry
	nil,                                   // 11: proto.control_plane.v1.StaticInfoExtension.InternalHopsEntry
	nil,                                   // 12: proto.control_plane.v1.LatencyInfo.IntraEntry
	nil,                                   // 13: proto.control_plane.v1.LatencyInfo.InterEntry
	nil,                                   // 14: proto.control_plane.v1.BandwidthInfo.IntraEntry
	nil,                                   // 15: proto.control_plane.v1.BandwidthInfo.InterEntry
	(*DigestExtension_Digest)(nil),        // 16: proto.control_plane.v1.DigestExtension.Digest
	(*experimental.EPICDetachedExtension)(nil), // 17: proto.control_plane.experimental.v1.EPICDetachedExtension
}
var file_proto_control_plane_v1_seg_extensions_proto_depIdxs = []int32{
	3,  // 0: proto.control_plane.v1.PathSegmentExtensions.static_info:type_name -> proto.control_plane.v1.StaticInfoExtension
	2,  // 1: proto.control_plane.v1.PathSegmentExtensions.hidden_path:type_name -> proto.control_plane.v1.HiddenPathExtension
	7,  // 2: proto.control_plane.v1.PathSegmentExtensions.digests:type_name -> proto.control_plane.v1.DigestExtension
	4,  // 3: proto.control_plane.v1.StaticInfoExtension.latency:type_name -> proto.control_plane.v1.LatencyInfo
	5,  // 4: proto.control_plane.v1.StaticInfoExtension.bandwidth:type_name -> proto.control_plane.v1.BandwidthInfo
	9,  // 5: proto.control_plane.v1.StaticInfoExtension.geo:type_name -> proto.control_plane.v1.StaticInfoExtension.GeoEntry
	10, // 6: proto.control_plane.v1.StaticInfoExtension.link_type:type_name -> proto.control_plane.v1.StaticInfoExtension.LinkTypeEntry
	11, // 7: proto.control_plane.v1.StaticInfoExtension.internal_hops:type_name -> proto.control_plane.v1.StaticInfoExtension.InternalHopsEntry
	12, // 8: proto.control_plane.v1.LatencyInfo.intra:type_name -> proto.control_plane.v1.LatencyInfo.IntraEntry
	13, // 9: proto.control_plane.v1.LatencyInfo.inter:type_name -> proto.control_plane.v1.LatencyInfo.InterEntry
	14, // 10: proto.control_plane.v1.BandwidthInfo.intra:type_name -> proto.control_plane.v1.BandwidthInfo.IntraEntry
	15, // 11: proto.control_plane.v1.BandwidthInfo.inter:type_name -> proto.control_plane.v1.BandwidthInfo.InterEntry
	16, // 12: proto.control_plane.v1.DigestExtension.epic:type_name -> proto.control_plane.v1.DigestExtension.Digest
	17, // 13: proto.control_plane.v1.PathSegmentUnsignedExtensions.epic:type_name -> proto.control_plane.experimental.v1.EPICDetachedExtension
	6,  // 14: proto.control_plane.v1.StaticInfoExtension.GeoEntry.value:type_name -> proto.control_plane.v1.GeoCoordinates
	0,  // 15: proto.control_plane.v1.StaticInfoExtension.LinkTypeEntry.value:type_name -> proto.control_plane.v1.LinkType
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_proto_control_plane_v1_seg_extensions_proto_init() }
func file_proto_control_plane_v1_seg_extensions_proto_init() {
	if File_proto_control_plane_v1_seg_extensions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PathSegmentExtensions); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HiddenPathExtension); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StaticInfoExtension); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LatencyInfo); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BandwidthInfo); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeoCoordinates); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DigestExtension); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PathSegmentUnsignedExtensions); i {
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
		file_proto_control_plane_v1_seg_extensions_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DigestExtension_Digest); i {
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
			RawDescriptor: file_proto_control_plane_v1_seg_extensions_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_control_plane_v1_seg_extensions_proto_goTypes,
		DependencyIndexes: file_proto_control_plane_v1_seg_extensions_proto_depIdxs,
		EnumInfos:         file_proto_control_plane_v1_seg_extensions_proto_enumTypes,
		MessageInfos:      file_proto_control_plane_v1_seg_extensions_proto_msgTypes,
	}.Build()
	File_proto_control_plane_v1_seg_extensions_proto = out.File
	file_proto_control_plane_v1_seg_extensions_proto_rawDesc = nil
	file_proto_control_plane_v1_seg_extensions_proto_goTypes = nil
	file_proto_control_plane_v1_seg_extensions_proto_depIdxs = nil
}
