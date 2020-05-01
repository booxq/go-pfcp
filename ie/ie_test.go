// Copyright 2019-2020 go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/wmnsk/go-pfcp/ie"
)

func TestIEs(t *testing.T) {
	cases := []struct {
		description string
		structured  *ie.IE
		serialized  []byte
	}{
		{
			"Cause",
			ie.NewCause(ie.CauseRequestAccepted),
			[]byte{0x00, 0x13, 0x00, 0x01, 0x01},
		}, {
			"SourceInterface",
			ie.NewSourceInterface(ie.SrcInterfaceAccess),
			[]byte{0x00, 0x14, 0x00, 0x01, 0x00},
		}, {
			"FTEID/TEID/IPv4", // TODO: add other forms
			ie.NewFTEID(0x11111111, net.ParseIP("127.0.0.1"), nil, nil),
			[]byte{0x00, 0x15, 0x00, 0x09, 0x01, 0x11, 0x11, 0x11, 0x11, 0x7f, 0x00, 0x00, 0x01},
		}, {
			"FTEID/TEID/IPv6", // TODO: add other forms
			ie.NewFTEID(0x11111111, nil, net.ParseIP("2001::1"), nil),
			[]byte{0x00, 0x15, 0x00, 0x15, 0x02, 0x11, 0x11, 0x11, 0x11, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		}, {
			"NetworkInstance",
			ie.NewNetworkInstance("some.instance.example"),
			[]byte{0x00, 0x16, 0x00, 0x15, 0x73, 0x6f, 0x6d, 0x65, 0x2e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65},
		}, {
			"SDFFilter",
			ie.NewSDFFilter("aaaaaaaa", "bb", "cccc", "ddd", 0xffffffff),
			[]byte{
				0x00, 0x17, 0x00, 0x19,
				0x1f, 0x00, // Flags & Spare octet
				0x00, 0x08, 0x61, 0x61, 0x61, 0x61, 0x61, 0x61, 0x61, 0x61, // FD
				0x62, 0x62, // TTC
				0x63, 0x63, 0x63, 0x63, // SPI
				0x64, 0x64, 0x64, // FL
				0xff, 0xff, 0xff, 0xff, // BID
			},
		}, {
			"ApplicationID",
			ie.NewApplicationID("https://github.com/wmnsk/go-pfcp/"),
			[]byte{0x00, 0x18, 0x00, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6d, 0x6e, 0x73, 0x6b, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2f},
		}, {
			"GateStatus/OpenOpen",
			ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen),
			[]byte{0x00, 0x19, 0x00, 0x01, 0x00},
		}, {
			"GateStatus/OpenClosed",
			ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusClosed),
			[]byte{0x00, 0x19, 0x00, 0x01, 0x01},
		}, {
			"GateStatus/ClosedOpen",
			ie.NewGateStatus(ie.GateStatusClosed, ie.GateStatusOpen),
			[]byte{0x00, 0x19, 0x00, 0x01, 0x04},
		}, {
			"GateStatus/ClosedClosed",
			ie.NewGateStatus(ie.GateStatusClosed, ie.GateStatusClosed),
			[]byte{0x00, 0x19, 0x00, 0x01, 0x05},
		}, {
			"MBR",
			ie.NewMBR(0x11111111, 0x22222222),
			[]byte{0x00, 0x1a, 0x00, 0x08, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22},
		}, {
			"GBR",
			ie.NewGBR(0x11111111, 0x22222222),
			[]byte{0x00, 0x1b, 0x00, 0x08, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22},
		}, {
			"QERCorrelationID",
			ie.NewQERCorrelationID(0x11111111),
			[]byte{0x00, 0x1c, 0x00, 0x04, 0x11, 0x11, 0x11, 0x11},
		}, {
			"Precedence",
			ie.NewPrecedence(0x11111111),
			[]byte{0x00, 0x1d, 0x00, 0x04, 0x11, 0x11, 0x11, 0x11},
		}, {
			"TransportLevelMarking",
			ie.NewTransportLevelMarking(0x1111),
			[]byte{0x00, 0x1e, 0x00, 0x02, 0x11, 0x11},
		}, {
			"VolumeThreshold/TOTAL",
			ie.NewVolumeThreshold(0x01, 0x1111111111111111, 0, 0),
			[]byte{0x00, 0x1f, 0x00, 0x09, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		}, {
			"VolumeThreshold/ULDL",
			ie.NewVolumeThreshold(0x06, 0, 0x1111111111111111, 0x2222222222222222),
			[]byte{0x00, 0x1f, 0x00, 0x11, 0x06, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
		}, {
			"VolumeThreshold/ALL",
			ie.NewVolumeThreshold(0x07, 0x3333333333333333, 0x1111111111111111, 0x2222222222222222),
			[]byte{0x00, 0x1f, 0x00, 0x19, 0x07, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
		}, {
			"TimeThreshold",
			ie.NewTimeThreshold(0x11111111),
			[]byte{0x00, 0x20, 0x00, 0x04, 0x11, 0x11, 0x11, 0x11},
		}, {
			"MonitoringTime",
			ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			[]byte{0x00, 0x21, 0x00, 0x04, 0xdf, 0xd5, 0x2c, 0x00},
		}, {
			"SubsequentVolumeThreshold/TOTAL",
			ie.NewSubsequentVolumeThreshold(0x01, 0x1111111111111111, 0, 0),
			[]byte{0x00, 0x22, 0x00, 0x09, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		}, {
			"SubsequentVolumeThreshold/ULDL",
			ie.NewSubsequentVolumeThreshold(0x06, 0, 0x1111111111111111, 0x2222222222222222),
			[]byte{0x00, 0x22, 0x00, 0x11, 0x06, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
		}, {
			"SubsequentVolumeThreshold/ALL",
			ie.NewSubsequentVolumeThreshold(0x07, 0x3333333333333333, 0x1111111111111111, 0x2222222222222222),
			[]byte{0x00, 0x22, 0x00, 0x19, 0x07, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
		}, {
			"SubsequentTimeThreshold",
			ie.NewSubsequentTimeThreshold(0x11111111),
			[]byte{0x00, 0x23, 0x00, 0x04, 0x11, 0x11, 0x11, 0x11},
		}, {
			"InactivityDetectionTime",
			ie.NewInactivityDetectionTime(0x11111111),
			[]byte{0x00, 0x24, 0x00, 0x04, 0x11, 0x11, 0x11, 0x11},
		}, {
			"ReportingTriggers",
			ie.NewReportingTriggers(0x1122),
			[]byte{0x00, 0x25, 0x00, 0x02, 0x11, 0x22},
		}, {
			"RedirectInformation/URL/1",
			ie.NewRedirectInformation(ie.RedirectAddrURL, "https://github.com/wmnsk/go-pfcp/"),
			[]byte{0x00, 0x26, 0x00, 0x26, 0x02, 0x00, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6d, 0x6e, 0x73, 0x6b, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2f, 0x00, 0x00},
		}, {
			"RedirectInformation/URL/2",
			ie.NewRedirectInformation(ie.RedirectAddrURL, "https://github.com/wmnsk/go-pfcp/", "https://github.com/wmnsk/go-pfcp/"),
			[]byte{
				0x00, 0x26, 0x00, 0x47, 0x02,
				0x00, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6d, 0x6e, 0x73, 0x6b, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2f,
				0x00, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6d, 0x6e, 0x73, 0x6b, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2f,
			},
		}, {
			"RedirectInformation/IPv4IPv6",
			ie.NewRedirectInformation(ie.RedirectAddrIPv4AndIPv6, "127.0.0.1", "2001::1"),
			[]byte{0x00, 0x26, 0x00, 0x15, 0x04, 0x00, 0x09, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x00, 0x07, 0x32, 0x30, 0x30, 0x31, 0x3a, 0x3a, 0x31},
		}, {
			"ReportType",
			ie.NewReportType(1, 1, 1, 1),
			[]byte{0x00, 0x27, 0x00, 0x01, 0x0f},
		}, {
			"OffendingIE",
			ie.NewOffendingIE(ie.Cause),
			[]byte{0x00, 0x28, 0x00, 0x02, 0x00, 0x13},
		}, {
			"ForwardingPolicy",
			ie.NewForwardingPolicy("go-pfcp"),
			[]byte{0x00, 0x29, 0x00, 0x08, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70},
		}, {
			"DestinationInterface",
			ie.NewDestinationInterface(ie.DstInterfaceAccess),
			[]byte{0x00, 0x2a, 0x00, 0x01, 0x00},
		}, {
			"UPFunctionFeatures/Normal",
			ie.NewUPFunctionFeatures(0x01, 0x02),
			[]byte{0x00, 0x2b, 0x00, 0x02, 0x01, 0x02},
		}, {
			"UPFunctionFeatures/HasAdditional",
			ie.NewUPFunctionFeatures(0x01, 0x02, 0x03, 0x04),
			[]byte{0x00, 0x2b, 0x00, 0x04, 0x01, 0x02, 0x03, 0x04},
		}, {
			"ApplyAction",
			ie.NewApplyAction(0x04),
			[]byte{0x00, 0x2c, 0x00, 0x01, 0x04},
		}, {
			"DownlinkDataServiceInformation/HasPPI",
			ie.NewDownlinkDataServiceInformation(true, false, 0xff, 0),
			[]byte{0x00, 0x2d, 0x00, 0x02, 0x01, 0xff},
		}, {
			"DownlinkDataServiceInformation/HasQFI",
			ie.NewDownlinkDataServiceInformation(false, true, 0, 0xff),
			[]byte{0x00, 0x2d, 0x00, 0x02, 0x02, 0xff},
		}, {
			"DownlinkDataServiceInformation/HasALL",
			ie.NewDownlinkDataServiceInformation(true, true, 0xff, 0xff),
			[]byte{0x00, 0x2d, 0x00, 0x03, 0x03, 0xff, 0xff},
		}, {
			"DownlinkDataNotificationDelay",
			ie.NewDownlinkDataNotificationDelay(100 * time.Millisecond),
			[]byte{0x00, 0x2e, 0x00, 0x01, 0x02},
		}, {
			"DLBufferingDuration/20hr",
			ie.NewDLBufferingDuration(20 * time.Hour),
			[]byte{0x00, 0x2f, 0x00, 0x01, 0x82},
		}, {
			"DLBufferingDuration/30sec",
			ie.NewDLBufferingDuration(30 * time.Second),
			[]byte{0x00, 0x2f, 0x00, 0x01, 0x0f},
		}, {
			"DLBufferingDuration/15min",
			ie.NewDLBufferingDuration(15 * time.Minute),
			[]byte{0x00, 0x2f, 0x00, 0x01, 0x2f},
		}, {
			"DLBufferingSuggestedPacketCount/uint8",
			ie.NewDLBufferingSuggestedPacketCount(0xff),
			[]byte{0x00, 0x30, 0x00, 0x01, 0xff},
		}, {
			"DLBufferingSuggestedPacketCount/uint16",
			ie.NewDLBufferingSuggestedPacketCount(0xffff),
			[]byte{0x00, 0x30, 0x00, 0x02, 0xff, 0xff},
		}, {
			"SxSMReqFlags",
			ie.NewSxSMReqFlags(0x03),
			[]byte{0x00, 0x31, 0x00, 0x01, 0x03},
		}, {
			"SxSRRspFlags",
			ie.NewSxSRRspFlags(0x01),
			[]byte{0x00, 0x32, 0x00, 0x01, 0x01},
		}, {
			"LoadControlInformation",
			ie.NewLoadControlInformation(ie.NewSequenceNumber(0xffffffff), ie.NewMetric(0x01)),
			[]byte{0x00, 0x33, 0x00, 0x0d, 0x00, 0x34, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff, 0x00, 0x35, 0x00, 0x01, 0x01},
		}, {
			"SequenceNumber",
			ie.NewSequenceNumber(0xffffffff),
			[]byte{0x00, 0x34, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"Metric",
			ie.NewMetric(0x01),
			[]byte{0x00, 0x35, 0x00, 0x01, 0x01},
		}, {
			"OverloadControlInformation",
			ie.NewOverloadControlInformation(
				ie.NewSequenceNumber(0xffffffff),
				ie.NewMetric(0x01),
				ie.NewTimer(20*time.Hour),
				ie.NewOCIFlags(0x01),
			),
			[]byte{
				0x00, 0x36, 0x00, 0x17,
				0x00, 0x34, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff,
				0x00, 0x35, 0x00, 0x01, 0x01,
				0x00, 0x37, 0x00, 0x01, 0x82,
				0x00, 0x6e, 0x00, 0x01, 0x01,
			},
		}, {
			"Timer/20hr",
			ie.NewTimer(20 * time.Hour),
			[]byte{0x00, 0x37, 0x00, 0x01, 0x82},
		}, {
			"Timer/30sec",
			ie.NewTimer(30 * time.Second),
			[]byte{0x00, 0x37, 0x00, 0x01, 0x0f},
		}, {
			"Timer/15min",
			ie.NewTimer(15 * time.Minute),
			[]byte{0x00, 0x37, 0x00, 0x01, 0x2f},
		}, {
			"PacketDetectionRuleID",
			ie.NewPacketDetectionRuleID(0xffff),
			[]byte{0x00, 0x38, 0x00, 0x02, 0xff, 0xff},
		}, {
			"FSEID/SEID/IPv4", // TODO: add other forms
			ie.NewFSEID(0x1111111122222222, net.ParseIP("127.0.0.1"), nil, nil),
			[]byte{0x00, 0x39, 0x00, 0x0d, 0x02, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x7f, 0x00, 0x00, 0x01},
		}, {
			"ApplicationIDsPFDs",
			ie.NewApplicationIDsPFDs(
				ie.NewApplicationID("https://github.com/wmnsk/go-pfcp/"),
				ie.NewPFDContext(ie.NewPFDContents("aa", "bb", "cc", "dd", "ee", []string{"11", "22"}, []string{"33", "44"}, []string{"55", "66"})),
			),
			[]byte{
				0x00, 0x3a, 0x00, 0x61,
				// ApplicationID
				0x00, 0x18, 0x00, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6d, 0x6e, 0x73, 0x6b, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2f,
				// PFDContext
				0x00, 0x3b, 0x00, 0x38,
				0x00, 0x3d, 0x00, 0x34,
				0xff, 0x00, // Flags & Spare octet
				0x00, 0x02, 0x61, 0x61, // FD
				0x00, 0x02, 0x62, 0x62, // URL
				0x00, 0x02, 0x63, 0x63, // DN
				0x00, 0x02, 0x64, 0x64, // CP
				0x00, 0x02, 0x65, 0x65, // DNP
				0x00, 0x08, 0x00, 0x02, 0x31, 0x31, 0x00, 0x02, 0x32, 0x32, // AFD
				0x00, 0x08, 0x00, 0x02, 0x33, 0x33, 0x00, 0x02, 0x34, 0x34, // AURL
				0x00, 0x08, 0x00, 0x02, 0x35, 0x35, 0x00, 0x02, 0x36, 0x36, // ADNP
			},
		}, {
			"PFDContext",
			ie.NewPFDContext(ie.NewPFDContents("aa", "bb", "cc", "dd", "ee", []string{"11", "22"}, []string{"33", "44"}, []string{"55", "66"})),
			[]byte{
				0x00, 0x3b, 0x00, 0x38,
				0x00, 0x3d, 0x00, 0x34,
				0xff, 0x00, // Flags & Spare octet
				0x00, 0x02, 0x61, 0x61, // FD
				0x00, 0x02, 0x62, 0x62, // URL
				0x00, 0x02, 0x63, 0x63, // DN
				0x00, 0x02, 0x64, 0x64, // CP
				0x00, 0x02, 0x65, 0x65, // DNP
				0x00, 0x08, 0x00, 0x02, 0x31, 0x31, 0x00, 0x02, 0x32, 0x32, // AFD
				0x00, 0x08, 0x00, 0x02, 0x33, 0x33, 0x00, 0x02, 0x34, 0x34, // AURL
				0x00, 0x08, 0x00, 0x02, 0x35, 0x35, 0x00, 0x02, 0x36, 0x36, // ADNP
			},
		}, {
			"NodeID/IPv4", // TODO: add IPv6
			ie.NewNodeID("127.0.0.1", "", ""),
			[]byte{0x00, 0x3c, 0x00, 0x05, 0x00, 0x7f, 0x00, 0x00, 0x01},
		}, {
			"NodeID/FQDN", // TODO: add IPv6
			ie.NewNodeID("", "", "go-pfcp.epc.3gppnetwork.org"),
			[]byte{0x00, 0x3c, 0x00, 0x1c, 0x02, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2e, 0x65, 0x70, 0x63, 0x2e, 0x33, 0x67, 0x70, 0x70, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x6f, 0x72, 0x67},
		}, {
			"PFDContents",
			ie.NewPFDContents("aa", "bb", "cc", "dd", "ee", []string{"11", "22"}, []string{"33", "44"}, []string{"55", "66"}),
			[]byte{
				0x00, 0x3d, 0x00, 0x34,
				0xff, 0x00, // Flags & Spare octet
				0x00, 0x02, 0x61, 0x61, // FD
				0x00, 0x02, 0x62, 0x62, // URL
				0x00, 0x02, 0x63, 0x63, // DN
				0x00, 0x02, 0x64, 0x64, // CP
				0x00, 0x02, 0x65, 0x65, // DNP
				0x00, 0x08, 0x00, 0x02, 0x31, 0x31, 0x00, 0x02, 0x32, 0x32, // AFD
				0x00, 0x08, 0x00, 0x02, 0x33, 0x33, 0x00, 0x02, 0x34, 0x34, // AURL
				0x00, 0x08, 0x00, 0x02, 0x35, 0x35, 0x00, 0x02, 0x36, 0x36, // ADNP
			},
		}, {
			"MeasurementMethod",
			ie.NewMeasurementMethod(1, 1, 1),
			[]byte{0x00, 0x3e, 0x00, 0x01, 0x07},
		}, {
			"UsageReportTrigger",
			ie.NewUsageReportTrigger(0xff, 0xff, 0xff),
			[]byte{0x00, 0x3f, 0x00, 0x03, 0xff, 0xff, 0xff},
		}, {
			"MeasurementPeriod",
			ie.NewMeasurementPeriod(10 * time.Second),
			[]byte{0x00, 0x40, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"FQCSID/IPv4/SingleCSID",
			ie.NewFQCSID("127.0.0.1", 1),
			[]byte{0x00, 0x41, 0x00, 0x07, 0x01, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x01},
		}, {
			"FQCSID/IPv4/MultiCSIDs",
			ie.NewFQCSID("127.0.0.1", 1, 2),
			[]byte{0x00, 0x41, 0x00, 0x09, 0x02, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x02},
		}, {
			"FQCSID/IPv6/SingleCSID",
			ie.NewFQCSID("2001::1", 1),
			[]byte{0x00, 0x41, 0x00, 0x13, 0x11, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01},
		}, {
			"FQCSID/other",
			ie.NewFQCSID("12304501", 1),
			[]byte{0x00, 0x41, 0x00, 0x07, 0x21, 0x12, 0x30, 0x45, 0x01, 0x00, 0x01},
		}, {
			"DurationMeasurement",
			ie.NewDurationMeasurement(10 * time.Second),
			[]byte{0x00, 0x43, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"ApplicationDetectionInformation",
			ie.NewApplicationDetectionInformation(
				ie.NewApplicationID("https://github.com/wmnsk/go-pfcp/"),
				ie.NewApplicationInstanceID("go-pfcp"),
				ie.NewFlowInformation(ie.FlowDirectionDownlink, "go-pfcp"),
				ie.NewPacketDetectionRuleID(0xffff),
			),
			[]byte{
				0x00, 0x44, 0x00, 0x44,
				0x00, 0x18, 0x00, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6d, 0x6e, 0x73, 0x6b, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x2f,
				0x00, 0x5b, 0x00, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70,
				0x00, 0x5c, 0x00, 0x0a, 0x01, 0x00, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70,
				0x00, 0x38, 0x00, 0x02, 0xff, 0xff,
			},
		}, {
			"TimeOfFirstPacket",
			ie.NewTimeOfFirstPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			[]byte{0x00, 0x45, 0x00, 0x04, 0xdf, 0xd5, 0x2c, 0x00},
		}, {
			"TimeOfLastPacket",
			ie.NewTimeOfLastPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			[]byte{0x00, 0x46, 0x00, 0x04, 0xdf, 0xd5, 0x2c, 0x00},
		}, {
			"QuotaHoldingTime",
			ie.NewQuotaHoldingTime(10 * time.Second),
			[]byte{0x00, 0x47, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"DroppedDLTrafficThreshold/Packets",
			ie.NewDroppedDLTrafficThreshold(true, false, 0xffffffffffffffff, 0),
			[]byte{0x00, 0x48, 0x00, 0x09, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"DroppedDLTrafficThreshold/Bytes",
			ie.NewDroppedDLTrafficThreshold(false, true, 0, 0xffffffffffffffff),
			[]byte{0x00, 0x48, 0x00, 0x09, 0x02, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"DroppedDLTrafficThreshold/Both",
			ie.NewDroppedDLTrafficThreshold(true, true, 0xffffffffffffffff, 0xffffffffffffffff),
			[]byte{0x00, 0x48, 0x00, 0x11, 0x03, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"VolumeQuota/Total",
			ie.NewVolumeQuota(true, false, false, 0xffffffffffffffff, 0, 0),
			[]byte{0x00, 0x49, 0x00, 0x09, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"VolumeQuota/UL",
			ie.NewVolumeQuota(false, true, false, 0, 0xffffffffffffffff, 0),
			[]byte{0x00, 0x49, 0x00, 0x09, 0x02, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"VolumeQuota/DL",
			ie.NewVolumeQuota(false, false, true, 0, 0, 0xffffffffffffffff),
			[]byte{0x00, 0x49, 0x00, 0x09, 0x04, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"VolumeQuota/ALL",
			ie.NewVolumeQuota(true, true, true, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff),
			[]byte{0x00, 0x49, 0x00, 0x19, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		}, {
			"TimeQuota",
			ie.NewTimeQuota(10 * time.Second),
			[]byte{0x00, 0x4a, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"StartTime",
			ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			[]byte{0x00, 0x4b, 0x00, 0x04, 0xdf, 0xd5, 0x2c, 0x00},
		}, {
			"EndTime",
			ie.NewEndTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			[]byte{0x00, 0x4c, 0x00, 0x04, 0xdf, 0xd5, 0x2c, 0x00},
		}, {
			"QueryURR",
			ie.NewQueryURR(
				ie.NewURRID(0xffffffff),
			),
			[]byte{
				0x00, 0x4d, 0x00, 0x08,
				0x00, 0x51, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff,
			},
		}, {
			"URRID",
			ie.NewURRID(0xffffffff),
			[]byte{0x00, 0x51, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"LinkedURRID",
			ie.NewLinkedURRID(0xffffffff),
			[]byte{0x00, 0x52, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"OuterHeaderCreation/TEID/IPv4", // TODO: add other types
			ie.NewOuterHeaderCreation(0x0100, 0x11223344, "127.0.0.1", "", 0, 0, 0),
			[]byte{0x00, 0x54, 0x00, 0x0a, 0x01, 0x00, 0x11, 0x22, 0x33, 0x44, 0x7f, 0x00, 0x00, 0x01},
		}, {
			"OuterHeaderCreation/IPv4/Port", // TODO: add other types
			ie.NewOuterHeaderCreation(0x0400, 0, "127.0.0.1", "", 2152, 0, 0),
			[]byte{0x00, 0x54, 0x00, 0x08, 0x04, 0x00, 0x7f, 0x00, 0x00, 0x01, 0x08, 0x68},
		}, {
			"OuterHeaderCreation/IPv4", // TODO: add other types
			ie.NewOuterHeaderCreation(0x1000, 0, "127.0.0.1", "", 0, 0, 0),
			[]byte{0x00, 0x54, 0x00, 0x06, 0x10, 0x00, 0x7f, 0x00, 0x00, 0x01},
		}, {
			"BARID",
			ie.NewBARID(0xff),
			[]byte{0x00, 0x58, 0x00, 0x01, 0xff},
		}, {
			"CPFunctionFeatures",
			ie.NewCPFunctionFeatures(0x3f),
			[]byte{0x00, 0x59, 0x00, 0x01, 0x3f},
		}, {
			"UsageInformation",
			ie.NewUsageInformation(1, 1, 1, 1),
			[]byte{0x00, 0x5a, 0x00, 0x01, 0x0f},
		}, {
			"ApplicationInstanceID",
			ie.NewApplicationInstanceID("go-pfcp"),
			[]byte{0x00, 0x5b, 0x00, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70},
		}, {
			"FlowInformation",
			ie.NewFlowInformation(ie.FlowDirectionDownlink, "go-pfcp"),
			[]byte{0x00, 0x5c, 0x00, 0x0a, 0x01, 0x00, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70},
		}, {
			"UEIPAddress/IPv4", // TODO: add other types
			ie.NewUEIPAddress(0x02, "127.0.0.1", "", 0),
			[]byte{0x00, 0x5d, 0x00, 0x05, 0x02, 0x7f, 0x00, 0x00, 0x01},
		}, {
			"PacketRate/Both", // TODO: add other types
			ie.NewPacketRate(0x03, ie.TimeUnitMinute, 0x1122, ie.TimeUnitMinute, 0x3344),
			[]byte{0x00, 0x5e, 0x00, 0x07, 0x03, 0x00, 0x11, 0x22, 0x00, 0x33, 0x44},
		}, {
			"OuterHeaderRemoval",
			ie.NewOuterHeaderRemoval(0x01, 0x02),
			[]byte{0x00, 0x5f, 0x00, 0x02, 0x01, 0x02},
		}, {
			"RecoveryTimeStamp",
			ie.NewRecoveryTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			[]byte{0x00, 0x60, 0x00, 0x04, 0xdf, 0xd5, 0x2c, 0x00},
		}, {
			"DLFlowLevelMarking",
			ie.NewDLFlowLevelMarking(0x03, 0x1122, 0x3344),
			[]byte{0x00, 0x61, 0x00, 0x05, 0x03, 0x11, 0x22, 0x33, 0x44},
		}, {
			"HeaderEnrichment",
			ie.NewHeaderEnrichment(ie.HeaderTypeHTTP, "name", "value"),
			[]byte{0x00, 0x62, 0x00, 0x0c, 0x00, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65},
		}, {
			"MeasurementInformation",
			ie.NewMeasurementInformation(0x1f),
			[]byte{0x00, 0x64, 0x00, 0x01, 0x1f},
		}, {
			"NodeReportType",
			ie.NewNodeReportType(0x01),
			[]byte{0x00, 0x65, 0x00, 0x01, 0x01},
		}, {
			"RemoteGTPUPeer/IPv4/DI/NI", // TODO: add other types
			ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, "some.instance.example"),
			[]byte{
				0x00, 0x67, 0x00, 0x1f,
				0x0e,
				0x7f, 0x00, 0x00, 0x01,
				0x00, 0x01, 0x00,
				0x00, 0x15, 0x73, 0x6f, 0x6d, 0x65, 0x2e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
			},
		}, {
			"URSEQN",
			ie.NewURSEQN(0xffffffff),
			[]byte{0x00, 0x68, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"ActivatePredefinedRules",
			ie.NewActivatePredefinedRules("go-pfcp"),
			[]byte{0x00, 0x6a, 0x00, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70},
		}, {
			"DeactivatePredefinedRules",
			ie.NewDeactivatePredefinedRules("go-pfcp"),
			[]byte{0x00, 0x6b, 0x00, 0x07, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70},
		}, {
			"FARID",
			ie.NewFARID(0xffffffff),
			[]byte{0x00, 0x6c, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"QERID",
			ie.NewQERID(0xffffffff),
			[]byte{0x00, 0x6d, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"OCIFlags",
			ie.NewOCIFlags(0x01),
			[]byte{0x00, 0x6e, 0x00, 0x01, 0x01},
		}, {
			"SxAssociationReleaseRequest",
			ie.NewSxAssociationReleaseRequest(1, 1),
			[]byte{0x00, 0x6f, 0x00, 0x01, 0x03},
		}, {
			"GracefulReleasePeriod/20hr",
			ie.NewGracefulReleasePeriod(20 * time.Hour),
			[]byte{0x00, 0x70, 0x00, 0x01, 0x82},
		}, {
			"GracefulReleasePeriod/30sec",
			ie.NewGracefulReleasePeriod(30 * time.Second),
			[]byte{0x00, 0x70, 0x00, 0x01, 0x0f},
		}, {
			"GracefulReleasePeriod/15min",
			ie.NewGracefulReleasePeriod(15 * time.Minute),
			[]byte{0x00, 0x70, 0x00, 0x01, 0x2f},
		}, {
			"PDNType",
			ie.NewPDNType(ie.PDNTypeIPv4),
			[]byte{0x00, 0x71, 0x00, 0x01, 0x01},
		}, {
			"FailedRuleID/PDR",
			ie.NewFailedRuleID(ie.RuleIDTypePDR, 0xffff),
			[]byte{0x00, 0x72, 0x00, 0x03, 0x00, 0xff, 0xff},
		}, {
			"FailedRuleID/FAR",
			ie.NewFailedRuleID(ie.RuleIDTypeFAR, 0xffffffff),
			[]byte{0x00, 0x72, 0x00, 0x05, 0x01, 0xff, 0xff, 0xff, 0xff},
		}, {
			"FailedRuleID/BAR",
			ie.NewFailedRuleID(ie.RuleIDTypeBAR, 0xff),
			[]byte{0x00, 0x72, 0x00, 0x02, 0x04, 0xff},
		}, {
			"TimeQuotaMechanism",
			ie.NewTimeQuotaMechanism(ie.BTITCTP, 10*time.Second),
			[]byte{0x00, 0x73, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"UserPlaneIPResourceInformation/IPv4/TEIDRange/NI/SI", // TODO: add other types
			ie.NewUserPlaneIPResourceInformation(0x71, 15, "127.0.0.1", "", "some.instance.example", ie.SrcInterfaceAccess),
			[]byte{
				0x00, 0x74, 0x00, 0x1c,
				0x71,
				0x0f,
				0x7f, 0x00, 0x00, 0x01,
				0x73, 0x6f, 0x6d, 0x65, 0x2e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
				0x00,
			},
		}, {
			"UserPlaneInactivityTimer",
			ie.NewUserPlaneInactivityTimer(10 * time.Second),
			[]byte{0x00, 0x75, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"Multiplier",
			ie.NewMultiplier(0xffffffffffffffff, 0x11223344),
			[]byte{0x00, 0x77, 0x00, 0x0c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x11, 0x22, 0x33, 0x44},
		}, {
			"AggregatedURRID",
			ie.NewAggregatedURRID(0xffffffff),
			[]byte{0x00, 0x78, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"SubsequentVolumeQuota/TOTAL",
			ie.NewSubsequentVolumeQuota(0x01, 0x1111111111111111, 0, 0),
			[]byte{0x00, 0x79, 0x00, 0x09, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		}, {
			"SubsequentVolumeQuota/ULDL",
			ie.NewSubsequentVolumeQuota(0x06, 0, 0x1111111111111111, 0x2222222222222222),
			[]byte{0x00, 0x79, 0x00, 0x11, 0x06, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
		}, {
			"SubsequentVolumeQuota/ALL",
			ie.NewSubsequentVolumeQuota(0x07, 0x3333333333333333, 0x1111111111111111, 0x2222222222222222),
			[]byte{0x00, 0x79, 0x00, 0x19, 0x07, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
		}, {
			"SubsequentTimeQuota",
			ie.NewSubsequentTimeQuota(10 * time.Second),
			[]byte{0x00, 0x7a, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"RQI",
			ie.NewRQI(0x01),
			[]byte{0x00, 0x7b, 0x00, 0x01, 0x01},
		}, {
			"QFI",
			ie.NewQFI(0x01),
			[]byte{0x00, 0x7c, 0x00, 0x01, 0x01},
		}, {
			"QueryURRReference",
			ie.NewQueryURRReference(0xffffffff),
			[]byte{0x00, 0x7d, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"AdditionalUsageReportsInformation",
			ie.NewAdditionalUsageReportsInformation(0x00ff),
			[]byte{0x00, 0x7e, 0x00, 0x02, 0x80, 0xff},
		}, {
			"TrafficEndpointID",
			ie.NewTrafficEndpointID(0x01),
			[]byte{0x00, 0x83, 0x00, 0x01, 0x01},
		}, {
			"Ethertype",
			ie.NewEthertype(0xffff),
			[]byte{0x00, 0x88, 0x00, 0x02, 0xff, 0xff},
		}, {
			"Proxying",
			ie.NewProxying(1, 1),
			[]byte{0x00, 0x89, 0x00, 0x01, 0x03},
		}, {
			"EthernetFilterID",
			ie.NewEthernetFilterID(0xffffffff),
			[]byte{0x00, 0x8a, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"EthernetFilterProperties",
			ie.NewEthernetFilterProperties(0x01),
			[]byte{0x00, 0x8b, 0x00, 0x01, 0x01},
		}, {
			"SuggestedBufferingPacketsCount",
			ie.NewSuggestedBufferingPacketsCount(0x01),
			[]byte{0x00, 0x8c, 0x00, 0x01, 0x01},
		}, {
			"UserID",
			ie.NewUserID(0x0f, "123451234567890", "123451234567890", "123451234567890", "go-pfcp@github.com"),
			[]byte{
				0x00, 0x8d, 0x00, 0x2f, 0x0f,
				0x08, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
				0x08, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
				0x08, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
				0x12, 0x67, 0x6f, 0x2d, 0x70, 0x66, 0x63, 0x70, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
			},
		}, {
			"EthernetPDUSessionInformation",
			ie.NewEthernetPDUSessionInformation(0x01),
			[]byte{0x00, 0x8e, 0x00, 0x01, 0x01},
		}, {
			"EthernetInactivityTimer",
			ie.NewEthernetInactivityTimer(10 * time.Second),
			[]byte{0x00, 0x92, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a},
		}, {
			"SubsequentEventQuota",
			ie.NewSubsequentEventQuota(0xffffffff),
			[]byte{0x00, 0x96, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		}, {
			"SubsequentEventThreshold",
			ie.NewSubsequentEventThreshold(0xffffffff),
			[]byte{0x00, 0x97, 0x00, 0x04, 0xff, 0xff, 0xff, 0xff},
		},
	}

	for _, c := range cases {
		t.Run("marshal/"+c.description, func(t *testing.T) {
			got, err := c.structured.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, c.serialized); diff != "" {
				t.Error(diff)
			}
		})

		t.Run("unmarshal/"+c.description, func(t *testing.T) {
			got, err := ie.Parse(c.serialized)
			if err != nil {
				t.Fatal(err)
			}

			opt := cmp.AllowUnexported(*got, *c.structured)
			if diff := cmp.Diff(got, c.structured, opt); diff != "" {
				t.Error(diff)
			}
		})
	}
}
