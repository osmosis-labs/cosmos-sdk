// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/evidence/v1beta1/evidence.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Equivocation implements the Evidence interface and defines evidence of double
// signing misbehavior.
type Equivocation struct {
	Height           int64     `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Time             time.Time `protobuf:"bytes,2,opt,name=time,proto3,stdtime" json:"time"`
	Power            int64     `protobuf:"varint,3,opt,name=power,proto3" json:"power,omitempty"`
	ConsensusAddress string    `protobuf:"bytes,4,opt,name=consensus_address,json=consensusAddress,proto3" json:"consensus_address,omitempty" yaml:"consensus_address"`
}

func (m *Equivocation) Reset()      { *m = Equivocation{} }
func (*Equivocation) ProtoMessage() {}
func (*Equivocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd143e71a177f0dd, []int{0}
}
func (m *Equivocation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Equivocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Equivocation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Equivocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Equivocation.Merge(m, src)
}
func (m *Equivocation) XXX_Size() int {
	return m.Size()
}
func (m *Equivocation) XXX_DiscardUnknown() {
	xxx_messageInfo_Equivocation.DiscardUnknown(m)
}

var xxx_messageInfo_Equivocation proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Equivocation)(nil), "cosmos.evidence.v1beta1.Equivocation")
}

func init() {
	proto.RegisterFile("cosmos/evidence/v1beta1/evidence.proto", fileDescriptor_dd143e71a177f0dd)
}

var fileDescriptor_dd143e71a177f0dd = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x3f, 0x4f, 0x02, 0x31,
	0x18, 0xc6, 0x5b, 0x41, 0xa2, 0x27, 0x83, 0x5e, 0x88, 0x5e, 0x88, 0x69, 0x09, 0x83, 0x61, 0xe1,
	0x1a, 0x74, 0x31, 0x6c, 0x92, 0x38, 0x18, 0x37, 0xe2, 0xe4, 0x62, 0xee, 0x4f, 0x2d, 0x8d, 0xdc,
	0xbd, 0x27, 0xed, 0xa1, 0x7c, 0x03, 0x47, 0x46, 0x47, 0x46, 0x3f, 0x0a, 0x9b, 0x8c, 0x4e, 0x68,
	0x8e, 0xc5, 0xd9, 0x4f, 0x60, 0xb8, 0x02, 0x0e, 0x4e, 0xed, 0xf3, 0xe4, 0xf7, 0xfe, 0x92, 0xf7,
	0xb5, 0x4e, 0x02, 0x50, 0x11, 0x28, 0xc6, 0x87, 0x32, 0xe4, 0x71, 0xc0, 0xd9, 0xb0, 0xe5, 0x73,
	0xed, 0xb5, 0x36, 0x85, 0x9b, 0x0c, 0x40, 0x83, 0x7d, 0x64, 0x38, 0x77, 0x53, 0xaf, 0xb8, 0x6a,
	0x45, 0x80, 0x80, 0x9c, 0x61, 0xcb, 0x9f, 0xc1, 0xab, 0x54, 0x00, 0x88, 0x3e, 0x67, 0x79, 0xf2,
	0xd3, 0x7b, 0xa6, 0x65, 0xc4, 0x95, 0xf6, 0xa2, 0xc4, 0x00, 0xf5, 0x77, 0x6c, 0x95, 0x2f, 0x1f,
	0x53, 0x39, 0x84, 0xc0, 0xd3, 0x12, 0x62, 0xfb, 0xd0, 0x2a, 0xf5, 0xb8, 0x14, 0x3d, 0xed, 0xe0,
	0x1a, 0x6e, 0x14, 0xba, 0xab, 0x64, 0x9f, 0x5b, 0xc5, 0xe5, 0xac, 0xb3, 0x55, 0xc3, 0x8d, 0xbd,
	0xd3, 0xaa, 0x6b, 0xc4, 0xee, 0x5a, 0xec, 0xde, 0xac, 0xc5, 0x9d, 0x9d, 0xe9, 0x9c, 0xa2, 0xf1,
	0x27, 0xc5, 0xdd, 0x7c, 0xc2, 0xae, 0x58, 0xdb, 0x09, 0x3c, 0xf1, 0x81, 0x53, 0xc8, 0x85, 0x26,
	0xd8, 0x57, 0xd6, 0x41, 0x00, 0xb1, 0xe2, 0xb1, 0x4a, 0xd5, 0x9d, 0x17, 0x86, 0x03, 0xae, 0x94,
	0x53, 0xac, 0xe1, 0xc6, 0x6e, 0xe7, 0xf8, 0x67, 0x4e, 0x9d, 0x91, 0x17, 0xf5, 0xdb, 0xf5, 0x7f,
	0x48, 0xbd, 0xbb, 0xbf, 0xe9, 0x2e, 0x4c, 0xd5, 0x2e, 0xbf, 0x4c, 0x28, 0x7a, 0x9d, 0x50, 0xf4,
	0x3d, 0xa1, 0xa8, 0x73, 0xfd, 0x96, 0x11, 0x3c, 0xcd, 0x08, 0x9e, 0x65, 0x04, 0x7f, 0x65, 0x04,
	0x8f, 0x17, 0x04, 0xcd, 0x16, 0x04, 0x7d, 0x2c, 0x08, 0xba, 0x6d, 0x0a, 0xa9, 0x7b, 0xa9, 0xef,
	0x06, 0x10, 0xb1, 0xd5, 0xc9, 0xcd, 0xd3, 0x54, 0xe1, 0x03, 0x7b, 0xfe, 0xbb, 0xbf, 0x1e, 0x25,
	0x5c, 0xf9, 0xa5, 0x7c, 0xbf, 0xb3, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x61, 0x33, 0x3a, 0x69,
	0x9f, 0x01, 0x00, 0x00,
}

func (m *Equivocation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Equivocation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Equivocation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ConsensusAddress) > 0 {
		i -= len(m.ConsensusAddress)
		copy(dAtA[i:], m.ConsensusAddress)
		i = encodeVarintEvidence(dAtA, i, uint64(len(m.ConsensusAddress)))
		i--
		dAtA[i] = 0x22
	}
	if m.Power != 0 {
		i = encodeVarintEvidence(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x18
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintEvidence(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x12
	if m.Height != 0 {
		i = encodeVarintEvidence(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvidence(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvidence(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Equivocation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovEvidence(uint64(m.Height))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovEvidence(uint64(l))
	if m.Power != 0 {
		n += 1 + sovEvidence(uint64(m.Power))
	}
	l = len(m.ConsensusAddress)
	if l > 0 {
		n += 1 + l + sovEvidence(uint64(l))
	}
	return n
}

func sovEvidence(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvidence(x uint64) (n int) {
	return sovEvidence(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Equivocation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvidence
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Equivocation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Equivocation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvidence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvidence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsensusAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvidence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvidence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConsensusAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvidence(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvidence
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvidence(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvidence
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvidence
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvidence
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvidence
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvidence        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvidence          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvidence = fmt.Errorf("proto: unexpected end of group")
)
