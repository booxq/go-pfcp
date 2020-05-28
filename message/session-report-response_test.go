// Copyright 2019-2020 go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"net"
	"testing"
	"time"

	"github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"

	"github.com/wmnsk/go-pfcp/internal/testutil"
)

func TestSessionReportResponse(t *testing.T) {
	cases := []testutil.TestCase{
		{
			Description: "Normal",
			Structured: message.NewSessionReportResponse(
				mp, fo, seid, seq, pri,
				ie.NewCause(ie.CauseRequestAccepted),
				ie.NewOffendingIE(ie.Cause),
				ie.NewUpdateBARWithinSessionReportResponse(
					ie.NewBARID(0xff),
					ie.NewDownlinkDataNotificationDelay(100*time.Millisecond),
					ie.NewDLBufferingDuration(30*time.Second),
					ie.NewDLBufferingSuggestedPacketCount(0xffff),
					ie.NewSuggestedBufferingPacketsCount(0x01),
				),
				ie.NewPFCPSRRspFlags(0x01),
				ie.NewFSEID(0x1111111122222222, net.ParseIP("127.0.0.1"), nil, nil),
				ie.NewFTEID(0x11111111, net.ParseIP("127.0.0.1"), nil, nil),
				ie.NewAlternativeSMFIPAddress(net.ParseIP("127.0.0.1"), nil),
			),
			Serialized: []byte{
				0x21, 0x39, 0x00, 0x61, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x11, 0x22, 0x33, 0x00,
				0x00, 0x13, 0x00, 0x01, 0x01,
				0x00, 0x28, 0x00, 0x02, 0x00, 0x13,
				0x00, 0x0c, 0x00, 0x1a,
				0x00, 0x58, 0x00, 0x01, 0xff,
				0x00, 0x2e, 0x00, 0x01, 0x02,
				0x00, 0x2f, 0x00, 0x01, 0x0f,
				0x00, 0x30, 0x00, 0x02, 0xff, 0xff,
				0x00, 0x8c, 0x00, 0x01, 0x01,
				0x00, 0x32, 0x00, 0x01, 0x01,
				0x00, 0x39, 0x00, 0x0d, 0x02, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x7f, 0x00, 0x00, 0x01,
				0x00, 0x15, 0x00, 0x09, 0x01, 0x11, 0x11, 0x11, 0x11, 0x7f, 0x00, 0x00, 0x01,
				0x00, 0xb2, 0x00, 0x05, 0x02, 0x7f, 0x00, 0x00, 0x01,
			},
		},
	}

	testutil.Run(t, cases, func(b []byte) (testutil.Serializable, error) {
		v, err := message.ParseSessionReportResponse(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}