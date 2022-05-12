// Code generated by protoc-gen-go. DO NOT EDIT.
// source: flow/entities/block.proto

package entities

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Block struct {
	Id                       []byte                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ParentId                 []byte                  `protobuf:"bytes,2,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	Height                   uint64                  `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Timestamp                *timestamppb.Timestamp  `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	CollectionGuarantees     []*CollectionGuarantee  `protobuf:"bytes,5,rep,name=collection_guarantees,json=collectionGuarantees,proto3" json:"collection_guarantees,omitempty"`
	BlockSeals               []*BlockSeal            `protobuf:"bytes,6,rep,name=block_seals,json=blockSeals,proto3" json:"block_seals,omitempty"`
	Signatures               [][]byte                `protobuf:"bytes,7,rep,name=signatures,proto3" json:"signatures,omitempty"`
	ExecutionReceiptMetaList []*ExecutionReceiptMeta `protobuf:"bytes,8,rep,name=execution_receipt_metaList,json=executionReceiptMetaList,proto3" json:"execution_receipt_metaList,omitempty"`
	ExecutionResultList      []*ExecutionResult      `protobuf:"bytes,9,rep,name=execution_result_list,json=executionResultList,proto3" json:"execution_result_list,omitempty"`
	BlockHeader              *BlockHeader            `protobuf:"bytes,10,opt,name=block_header,json=blockHeader,proto3" json:"block_header,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}                `json:"-"`
	XXX_unrecognized         []byte                  `json:"-"`
	XXX_sizecache            int32                   `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1afc3335f2172fc, []int{0}
}

func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (m *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(m, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Block) GetParentId() []byte {
	if m != nil {
		return m.ParentId
	}
	return nil
}

func (m *Block) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Block) GetTimestamp() *timestamppb.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Block) GetCollectionGuarantees() []*CollectionGuarantee {
	if m != nil {
		return m.CollectionGuarantees
	}
	return nil
}

func (m *Block) GetBlockSeals() []*BlockSeal {
	if m != nil {
		return m.BlockSeals
	}
	return nil
}

func (m *Block) GetSignatures() [][]byte {
	if m != nil {
		return m.Signatures
	}
	return nil
}

func (m *Block) GetExecutionReceiptMetaList() []*ExecutionReceiptMeta {
	if m != nil {
		return m.ExecutionReceiptMetaList
	}
	return nil
}

func (m *Block) GetExecutionResultList() []*ExecutionResult {
	if m != nil {
		return m.ExecutionResultList
	}
	return nil
}

func (m *Block) GetBlockHeader() *BlockHeader {
	if m != nil {
		return m.BlockHeader
	}
	return nil
}

func init() {
	proto.RegisterType((*Block)(nil), "flow.entities.Block")
}

func init() { proto.RegisterFile("flow/entities/block.proto", fileDescriptor_c1afc3335f2172fc) }

var fileDescriptor_c1afc3335f2172fc = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0x4f, 0x8f, 0xd3, 0x30,
	0x10, 0xc5, 0xd5, 0x3f, 0x5b, 0xb6, 0xd3, 0xc2, 0xc1, 0xb0, 0xc8, 0x14, 0x54, 0xa2, 0x85, 0x43,
	0x4e, 0x0e, 0x5a, 0x2e, 0x70, 0xe0, 0xb2, 0x08, 0x01, 0x12, 0x48, 0xc8, 0x20, 0x21, 0x71, 0x89,
	0x9c, 0x64, 0x9a, 0x5a, 0xb8, 0x71, 0x15, 0x4f, 0x04, 0x1f, 0x8b, 0x8f, 0x88, 0x62, 0x37, 0x6d,
	0x53, 0xca, 0x25, 0x8a, 0xdf, 0xfb, 0xf9, 0x59, 0x79, 0xe3, 0xc0, 0xa3, 0x95, 0xb1, 0xbf, 0x12,
	0xac, 0x48, 0x93, 0x46, 0x97, 0x64, 0xc6, 0xe6, 0x3f, 0xc5, 0xb6, 0xb6, 0x64, 0xd9, 0xdd, 0xd6,
	0x12, 0x9d, 0xb5, 0x78, 0x5a, 0x5a, 0x5b, 0x1a, 0x4c, 0xbc, 0x99, 0x35, 0xab, 0x84, 0xf4, 0x06,
	0x1d, 0xa9, 0xcd, 0x36, 0xf0, 0x8b, 0x65, 0x3f, 0x2a, 0xb7, 0xc6, 0x60, 0x4e, 0xda, 0x56, 0xe7,
	0x7d, 0x7f, 0x54, 0xea, 0x50, 0x99, 0x9d, 0xff, 0xbc, 0xef, 0xe3, 0x6f, 0xcc, 0x9b, 0x76, 0x7b,
	0x5a, 0xa3, 0x6b, 0x0c, 0xed, 0xa8, 0xe8, 0x5c, 0xca, 0x1a, 0x55, 0x81, 0x75, 0x20, 0xae, 0xff,
	0x8c, 0xe1, 0xe2, 0xb6, 0x95, 0xd9, 0x3d, 0x18, 0xea, 0x82, 0x0f, 0xa2, 0x41, 0x3c, 0x97, 0x43,
	0x5d, 0xb0, 0xc7, 0x30, 0xdd, 0xaa, 0x1a, 0x2b, 0x4a, 0x75, 0xc1, 0x87, 0x5e, 0xbe, 0x0c, 0xc2,
	0xc7, 0x82, 0x3d, 0x84, 0xc9, 0x1a, 0x75, 0xb9, 0x26, 0x3e, 0x8a, 0x06, 0xf1, 0x58, 0xee, 0x56,
	0xec, 0x15, 0x4c, 0xf7, 0x5f, 0xca, 0xc7, 0xd1, 0x20, 0x9e, 0xdd, 0x2c, 0x44, 0xe8, 0x42, 0x74,
	0x5d, 0x88, 0x6f, 0x1d, 0x21, 0x0f, 0x30, 0xfb, 0x0e, 0x57, 0x87, 0x12, 0xd2, 0xb2, 0x51, 0xb5,
	0xaa, 0x08, 0xd1, 0xf1, 0x8b, 0x68, 0x14, 0xcf, 0x6e, 0xae, 0x45, 0xaf, 0x60, 0xf1, 0x76, 0xcf,
	0xbe, 0xef, 0x50, 0xf9, 0x20, 0xff, 0x57, 0x74, 0xec, 0x35, 0xcc, 0x0e, 0xed, 0x39, 0x3e, 0xf1,
	0x71, 0xfc, 0x24, 0xce, 0x57, 0xf0, 0x15, 0x95, 0x91, 0x90, 0x75, 0xaf, 0x8e, 0x2d, 0x01, 0x9c,
	0x2e, 0x2b, 0x45, 0x4d, 0x8d, 0x8e, 0xdf, 0x89, 0x46, 0xf1, 0x5c, 0x1e, 0x29, 0x4c, 0xc1, 0xe2,
	0xb8, 0xf8, 0x1c, 0xf5, 0x96, 0xd2, 0x0d, 0x92, 0xfa, 0xa4, 0x1d, 0xf1, 0x4b, 0x7f, 0xd2, 0xb3,
	0x93, 0x93, 0xde, 0x75, 0x1b, 0x64, 0xe0, 0x3f, 0x23, 0x29, 0xc9, 0xf1, 0x8c, 0xda, 0x86, 0x30,
	0x09, 0x57, 0xa7, 0xb3, 0x4d, 0x4d, 0x9b, 0x3e, 0xf5, 0xe9, 0xcb, 0xff, 0xa7, 0xb7, 0xa8, 0xbc,
	0x8f, 0x7d, 0xc1, 0x67, 0xbe, 0x81, 0xf9, 0xf1, 0x4d, 0xe0, 0xb0, 0x9b, 0xd3, 0x99, 0x4a, 0x3e,
	0x78, 0x42, 0x86, 0x06, 0xc3, 0xe2, 0xf6, 0x0b, 0x3c, 0xb1, 0x75, 0x29, 0x6c, 0xe5, 0xf9, 0xfd,
	0x54, 0xbb, 0x8d, 0x3f, 0x5e, 0x94, 0x9a, 0xd6, 0x4d, 0x26, 0x72, 0xbb, 0x49, 0x02, 0x94, 0xf8,
	0xc7, 0xfe, 0x5f, 0x28, 0x6d, 0xd2, 0xbb, 0x97, 0xd9, 0xc4, 0x5b, 0x2f, 0xff, 0x06, 0x00, 0x00,
	0xff, 0xff, 0xf2, 0x4c, 0x6f, 0x01, 0x60, 0x03, 0x00, 0x00,
}
