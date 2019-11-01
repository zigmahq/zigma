// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types.proto

package dht

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MessageType int32

const (
	MessageType_NOOP       MessageType = 0
	MessageType_PING       MessageType = 1
	MessageType_STORE      MessageType = 2
	MessageType_FIND_NODE  MessageType = 3
	MessageType_FIND_VALUE MessageType = 4
)

var MessageType_name = map[int32]string{
	0: "NOOP",
	1: "PING",
	2: "STORE",
	3: "FIND_NODE",
	4: "FIND_VALUE",
}

var MessageType_value = map[string]int32{
	"NOOP":       0,
	"PING":       1,
	"STORE":      2,
	"FIND_NODE":  3,
	"FIND_VALUE": 4,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

type ConnectionType int32

const (
	ConnectionType_NOT_CONNECTED  ConnectionType = 0
	ConnectionType_CONNECTED      ConnectionType = 1
	ConnectionType_CONNECTABLE    ConnectionType = 2
	ConnectionType_CANNOT_CONNECT ConnectionType = 3
)

var ConnectionType_name = map[int32]string{
	0: "NOT_CONNECTED",
	1: "CONNECTED",
	2: "CONNECTABLE",
	3: "CANNOT_CONNECT",
}

var ConnectionType_value = map[string]int32{
	"NOT_CONNECTED":  0,
	"CONNECTED":      1,
	"CONNECTABLE":    2,
	"CANNOT_CONNECT": 3,
}

func (x ConnectionType) String() string {
	return proto.EnumName(ConnectionType_name, int32(x))
}

func (ConnectionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

type Node struct {
	Id                   []byte         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Hash                 []byte         `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	PeerId               []byte         `protobuf:"bytes,3,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	Addrs                [][]byte       `protobuf:"bytes,4,rep,name=addrs,proto3" json:"addrs,omitempty"`
	Connection           ConnectionType `protobuf:"varint,5,opt,name=connection,proto3,enum=dht.ConnectionType" json:"connection,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}
func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Node) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Node) GetPeerId() []byte {
	if m != nil {
		return m.PeerId
	}
	return nil
}

func (m *Node) GetAddrs() [][]byte {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func (m *Node) GetConnection() ConnectionType {
	if m != nil {
		return m.Connection
	}
	return ConnectionType_NOT_CONNECTED
}

type Message struct {
	Id                   []byte      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 MessageType `protobuf:"varint,2,opt,name=type,proto3,enum=dht.MessageType" json:"type,omitempty"`
	Sender               []byte      `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver             []byte      `protobuf:"bytes,4,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Payload              []byte      `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	Hash                 []byte      `protobuf:"bytes,6,opt,name=hash,proto3" json:"hash,omitempty"`
	Sig                  []byte      `protobuf:"bytes,7,opt,name=sig,proto3" json:"sig,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Message) GetType() MessageType {
	if m != nil {
		return m.Type
	}
	return MessageType_NOOP
}

func (m *Message) GetSender() []byte {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *Message) GetReceiver() []byte {
	if m != nil {
		return m.Receiver
	}
	return nil
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Message) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Message) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func init() {
	proto.RegisterEnum("dht.MessageType", MessageType_name, MessageType_value)
	proto.RegisterEnum("dht.ConnectionType", ConnectionType_name, ConnectionType_value)
	proto.RegisterType((*Node)(nil), "dht.Node")
	proto.RegisterType((*Message)(nil), "dht.Message")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 347 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xcf, 0x4e, 0xc2, 0x40,
	0x10, 0xc6, 0xe9, 0x1f, 0x5a, 0x18, 0xa0, 0xae, 0xa3, 0xd1, 0x8d, 0x27, 0x42, 0x3c, 0x10, 0x0e,
	0x1c, 0xe0, 0x09, 0xb0, 0x54, 0x43, 0xc4, 0x2d, 0xa9, 0x55, 0x8f, 0xa4, 0xb2, 0x1b, 0xda, 0xc4,
	0xb4, 0x4d, 0xb7, 0x31, 0xe1, 0x25, 0x7c, 0x19, 0x5f, 0xd0, 0x74, 0xf9, 0x63, 0x8d, 0xb7, 0xf9,
	0x7d, 0x33, 0x3b, 0xfb, 0xed, 0xb7, 0xd0, 0x29, 0x77, 0xb9, 0x90, 0xe3, 0xbc, 0xc8, 0xca, 0x0c,
	0x0d, 0x1e, 0x97, 0x83, 0x2f, 0x0d, 0x4c, 0x96, 0x71, 0x81, 0x0e, 0xe8, 0x09, 0xa7, 0x5a, 0x5f,
	0x1b, 0x76, 0x03, 0x3d, 0xe1, 0x88, 0x60, 0xc6, 0x91, 0x8c, 0xa9, 0xae, 0x14, 0x55, 0xe3, 0x35,
	0xd8, 0xb9, 0x10, 0xc5, 0x3a, 0xe1, 0xd4, 0x50, 0xb2, 0x55, 0xe1, 0x82, 0xe3, 0x25, 0x34, 0x23,
	0xce, 0x0b, 0x49, 0xcd, 0xbe, 0x31, 0xec, 0x06, 0x7b, 0xc0, 0x29, 0xc0, 0x26, 0x4b, 0x53, 0xb1,
	0x29, 0x93, 0x2c, 0xa5, 0xcd, 0xbe, 0x36, 0x74, 0x26, 0x17, 0x63, 0x1e, 0x97, 0x63, 0xf7, 0x24,
	0x87, 0xbb, 0x5c, 0x04, 0xb5, 0xb1, 0xc1, 0xb7, 0x06, 0xf6, 0x93, 0x90, 0x32, 0xda, 0xfe, 0xf7,
	0x74, 0x0b, 0x66, 0xf5, 0x00, 0xe5, 0xc9, 0x99, 0x10, 0xb5, 0xea, 0x30, 0xab, 0xf6, 0xa8, 0x2e,
	0x5e, 0x81, 0x25, 0x45, 0xca, 0x45, 0x71, 0x34, 0xb9, 0x27, 0xbc, 0x81, 0x56, 0x21, 0x36, 0x22,
	0xf9, 0x14, 0x05, 0x35, 0x55, 0xe7, 0xc4, 0x48, 0xc1, 0xce, 0xa3, 0xdd, 0x47, 0x16, 0x71, 0xe5,
	0xb3, 0x1b, 0x1c, 0xf1, 0x94, 0x83, 0x55, 0xcb, 0x81, 0x80, 0x21, 0x93, 0x2d, 0xb5, 0x95, 0x54,
	0x95, 0xa3, 0x47, 0xe8, 0xd4, 0x8c, 0x60, 0x0b, 0x4c, 0xe6, 0xfb, 0x2b, 0xd2, 0xa8, 0xaa, 0xd5,
	0x82, 0x3d, 0x10, 0x0d, 0xdb, 0xd0, 0x7c, 0x0e, 0xfd, 0xc0, 0x23, 0x3a, 0xf6, 0xa0, 0x7d, 0xbf,
	0x60, 0xf3, 0x35, 0xf3, 0xe7, 0x1e, 0x31, 0xd0, 0x01, 0x50, 0xf8, 0x3a, 0x5b, 0xbe, 0x78, 0xc4,
	0x1c, 0xbd, 0x81, 0xf3, 0x37, 0x20, 0x3c, 0x87, 0x1e, 0xf3, 0xc3, 0xb5, 0xeb, 0x33, 0xe6, 0xb9,
	0xa1, 0x37, 0x27, 0x8d, 0x6a, 0xc7, 0x2f, 0x6a, 0x78, 0x06, 0x9d, 0x03, 0xce, 0xee, 0x96, 0xd5,
	0x1d, 0x08, 0x8e, 0x3b, 0x63, 0xb5, 0x53, 0xc4, 0x78, 0xb7, 0xd4, 0xc7, 0x4f, 0x7f, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xee, 0x82, 0x88, 0x7e, 0x07, 0x02, 0x00, 0x00,
}
