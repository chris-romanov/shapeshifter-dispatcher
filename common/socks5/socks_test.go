/*
 * Copyright (c) 2015, Yawning Angel <yawning at torproject dot org>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package socks5

import (
	"io"
	"net"
	"testing"
)

func tcpAddrsEqual(a, b *net.TCPAddr) bool {
	return a.IP.Equal(b.IP) && a.Port == b.Port
}

// TestAuthInvalidVersion tests auth negotiation with an invalid version.
func TestAuthInvalidVersion(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// VER = 03, NMETHODS = 01, METHODS = [00]
	_, hexErr := c.WriteHex("030100")
	if hexErr != nil {
		t.Error("NegotiateAuth(InvalidVersion) could not be decoded")
	}
	if _, err := req.NegotiateAuth(false); err == nil {
		t.Error("NegotiateAuth(InvalidVersion) succeeded")
	}
}

// TestAuthInvalidNMethods tests auth negotiation with no methods.
func TestAuthInvalidNMethods(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()
	var err error
	var method byte

	// VER = 05, NMETHODS = 00
	_, hexErr := c.WriteHex("0500")
	if hexErr != nil {
		t.Error("NegotiateAuth(No Methods) could not be decoded")
	}
	if method, err = req.NegotiateAuth(false); err != nil {
		t.Error("NegotiateAuth(No Methods) failed:", err)
	}
	if method != authNoAcceptableMethods {
		t.Error("NegotiateAuth(No Methods) picked unexpected method:", method)
	}
	if msg := c.ReadHex(); msg != "05ff" {
		t.Error("NegotiateAuth(No Methods) invalid response:", msg)
	}
}

// TestAuthNoneRequired tests auth negotiation with NO AUTHENTICATION REQUIRED.
func TestAuthNoneRequired(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()
	var err error
	var method byte

	// VER = 05, NMETHODS = 01, METHODS = [00]
	_, hexErr := c.WriteHex("050100")
	if hexErr != nil {
		t.Error("NegotiateAuth(None) could not be decoded")
	}
	if method, err = req.NegotiateAuth(false); err != nil {
		t.Error("NegotiateAuth(None) failed:", err)
	}
	if method != authNoneRequired {
		t.Error("NegotiateAuth(None) unexpected method:", method)
	}
	if msg := c.ReadHex(); msg != "0500" {
		t.Error("NegotiateAuth(None) invalid response:", msg)
	}
}

// TestAuthJsonParameterBlock tests auth negotiation with jsonParameterBlock.
func TestAuthJsonParameterBlock(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()
	var err error
	var method byte

	// VER = 05, NMETHODS = 01, METHODS = [09]
	//Method 9 is the json parameter block authentication
	_, hexErr := c.WriteHex("050109")
	if hexErr != nil{
		t.Error("NegotiateAuth(jsonParameterBlock) could not be decoded")
	}
	if method, err = req.NegotiateAuth(false); err != nil {
		t.Error("NegotiateAuth(jsonParameterBlock) failed:", err)
	}
	if method != AuthJsonParameterBlock {
		t.Error("NegotiateAuth(jsonParameterBlock) unexpected method:", method)
	}
	if msg := c.ReadHex(); msg != "0509" {
		t.Error("NegotiateAuth(jsonParameterBlock) invalid response:", msg)
	}
}

// TestAuthBoth tests auth negotiation containing both NO AUTHENTICATION
// REQUIRED and jsonParameterBlock.
func TestAuthBoth(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()
	var err error
	var method byte

	// VER = 05, NMETHODS = 02, METHODS = [00, 09]
	_, hexErr := c.WriteHex("05020009")
	if hexErr != nil {
		t.Error("NegotiateAuth(Both) could not be decoded")
	}
	if method, err = req.NegotiateAuth(true); err != nil {
		t.Error("NegotiateAuth(Both) failed:", err)
	}
	if method != AuthJsonParameterBlock {
		t.Error("NegotiateAuth(Both) unexpected method:", method)
	}
	if msg := c.ReadHex(); msg != "0509" {
		t.Error("NegotiateAuth(Both) invalid response:", msg)
	}
}

// TestAuthUnsupported tests auth negotiation with a unsupported method.
func TestAuthUnsupported(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()
	var err error
	var method byte

	// VER = 05, NMETHODS = 01, METHODS = [01] (GSSAPI)
	_, hexErr := c.WriteHex("050101")
	if hexErr != nil {
		t.Error("NegotiateAuth(Unknown) could not be decoded")
	}
	if method, err = req.NegotiateAuth(false); err != nil {
		t.Error("NegotiateAuth(Unknown) failed:", err)
	}
	if method != authNoAcceptableMethods {
		t.Error("NegotiateAuth(Unknown) picked unexpected method:", method)
	}
	if msg := c.ReadHex(); msg != "05ff" {
		t.Error("NegotiateAuth(Unknown) invalid response:", msg)
	}
}

// TestAuthUnsupported2 tests auth negotiation with supported and unsupported
// methods.
func TestAuthUnsupported2(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()
	var err error
	var method byte

	// VER = 05, NMETHODS = 03, METHODS = [00,01,09]
	_, hexErr := c.WriteHex("0503000109")
	if hexErr != nil {
		t.Error("NegotiateAuth(Unknown2) could not be decoded")
	}
	if method, err = req.NegotiateAuth(true); err != nil {
		t.Error("NegotiateAuth(Unknown2) failed:", err)
	}
	if method != AuthJsonParameterBlock {
		t.Error("NegotiateAuth(Unknown2) picked unexpected method:", method)
	}
	if msg := c.ReadHex(); msg != "0509" {
		t.Error("NegotiateAuth(Unknown2) invalid response:", msg)
	}
}

// TestRFC1928InvalidVersion tests RFC1929 auth with an invalid version.
func TestRFC1928InvalidVersion(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// VER = 03,  NMETHODS = 03, METHODS = [00,01,09], JLEN = 2, JSON = "{}"
	_, hexErr := c.WriteHex("0303000109000000027b7d")
	if hexErr != nil {
		t.Error("authenticate(InvalidVersion) could not be decoded")
	}
	if _, err := req.NegotiateAuth(true); err == nil {
		t.Error("failed to detect incorrect socks version:", err)
	}
}

// TestPT2Success tests PT2.1 jsonParameterBlock auth with valid pt args.
func TestPT2Success(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// JLEN = 2, JSON = "{}"
	_, hexErr := c.WriteHex("000000027b7d")
	if hexErr != nil {
		t.Error("authenticate(Success) could not be decoded")
	}
	if err := req.authenticate(AuthJsonParameterBlock); err != nil {
		t.Error("authenticate(Success) failed:", err)
	}
	if msg := c.ReadHex(); msg != "" {
		t.Error("authenticate(Success) invalid response:", msg)
	}
	if req.Args == nil {
		t.Error("RFC1929 k,v parse failure:")
	}
}

// TestPT2Fail tests PT2.1 jsonParameterBlock auth with invalid pt args.
func TestPT2Fail(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// JLEN = 2, JSON = "{}"
	_, hexErr := c.WriteHex("000000027d7b")
	if hexErr != nil {
		t.Error("authenticate(Success) could not be decoded")
	}
	if err := req.authenticate(AuthJsonParameterBlock); err == nil {
		t.Error("authenticate(Success) failed:", err)
	}
}
// TestRequestInvalidHdr tests SOCKS5 requests with invalid VER/CMD/RSV/ATYPE
func TestRequestInvalidHdr(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// VER = 03, CMD = 01, RSV = 00, ATYPE = 01, DST.ADDR = 127.0.0.1, DST.PORT = 9050
	_, hexErr := c.WriteHex("030100017f000001235a")
	if hexErr != nil {
		t.Error("readCommand(InvalidVer) could not be decoded")
	}
	if err := req.readCommand(); err == nil {
		t.Error("readCommand(InvalidVer) succeeded")
	}
	if msg := c.ReadHex(); msg != "05010001000000000000" {
		t.Error("readCommand(InvalidVer) invalid response:", msg)
	}
	c.reset(req)

	// VER = 05, CMD = 05, RSV = 00, ATYPE = 01, DST.ADDR = 127.0.0.1, DST.PORT = 9050
	_, hexErr2 := c.WriteHex("050500017f000001235a")
	if hexErr2 != nil {
		t.Error("readCommand(InvalidCmd) could not be decoded")
	}
	if err := req.readCommand(); err == nil {
		t.Error("readCommand(InvalidCmd) succeeded")
	}
	if msg := c.ReadHex(); msg != "05070001000000000000" {
		t.Error("readCommand(InvalidCmd) invalid response:", msg)
	}
	c.reset(req)

	// VER = 05, CMD = 01, RSV = 30, ATYPE = 01, DST.ADDR = 127.0.0.1, DST.PORT = 9050
	_, hexErr3 := c.WriteHex("050130017f000001235a")
	if hexErr3 != nil {
		t.Error("readCommand(InvalidRsv) could not be decoded")
	}
	if err := req.readCommand(); err == nil {
		t.Error("readCommand(InvalidRsv) succeeded")
	}
	if msg := c.ReadHex(); msg != "05010001000000000000" {
		t.Error("readCommand(InvalidRsv) invalid response:", msg)
	}
	c.reset(req)

	// VER = 05, CMD = 01, RSV = 01, ATYPE = 05, DST.ADDR = 127.0.0.1, DST.PORT = 9050
	_, hexErr4 := c.WriteHex("050100057f000001235a")
	if hexErr4 != nil {
		t.Error("readCommand(InvalidAtype) could not be decoded")
	}
	if err := req.readCommand(); err == nil {
		t.Error("readCommand(InvalidAtype) succeeded")
	}
	if msg := c.ReadHex(); msg != "05080001000000000000" {
		t.Error("readCommand(InvalidAtype) invalid response:", msg)
	}
	c.reset(req)
}

// TestRequestIPv4 tests IPv4 SOCKS5 requests.
func TestRequestIPv4(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// VER = 05, CMD = 01, RSV = 00, ATYPE = 01, DST.ADDR = 127.0.0.1, DST.PORT = 9050
	_, hexErr := c.WriteHex("050100017f000001235a")
	if hexErr != nil {
		t.Error("readCommand(IPv4) could not be decoded")
	}
	if err := req.readCommand(); err != nil {
		t.Error("readCommand(IPv4) failed:", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", req.Target)
	if err != nil {
		t.Error("net.ResolveTCPAddr failed:", err)
	}
	if !tcpAddrsEqual(addr, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9050}) {
		t.Error("Unexpected target:", addr)
	}
}

// TestRequestIPv6 tests IPv4 SOCKS5 requests.
func TestRequestIPv6(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// VER = 05, CMD = 01, RSV = 00, ATYPE = 04, DST.ADDR = 0102:0304:0506:0708:090a:0b0c:0d0e:0f10, DST.PORT = 9050
	_, hexErr := c.WriteHex("050100040102030405060708090a0b0c0d0e0f10235a")
	if hexErr != nil {
		t.Error("readCommand(IPv6) could not be decoded")
	}
	if err := req.readCommand(); err != nil {
		t.Error("readCommand(IPv6) failed:", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", req.Target)
	if err != nil {
		t.Error("net.ResolveTCPAddr failed:", err)
	}
	if !tcpAddrsEqual(addr, &net.TCPAddr{IP: net.ParseIP("0102:0304:0506:0708:090a:0b0c:0d0e:0f10"), Port: 9050}) {
		t.Error("Unexpected target:", addr)
	}
}

// TestRequestFQDN tests FQDN (DOMAINNAME) SOCKS5 requests.
func TestRequestFQDN(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	// VER = 05, CMD = 01, RSV = 00, ATYPE = 04, DST.ADDR = example.com, DST.PORT = 9050
	_, hexErr := c.WriteHex("050100030b6578616d706c652e636f6d235a")
	if hexErr != nil {
		t.Error("readCommand(FQDN) could not be decoded")
	}
	if err := req.readCommand(); err != nil {
		t.Error("readCommand(FQDN) failed:", err)
	}
	if req.Target != "example.com:9050" {
		t.Error("Unexpected target:", req.Target)
	}
}

// TestResponseNil tests nil address SOCKS5 responses.
func TestResponseNil(t *testing.T) {
	c := new(TestReadWriter)
	req := c.ToRequest()

	if err := req.Reply(ReplySucceeded); err != nil {
		t.Error("Reply(ReplySucceeded) failed:", err)
	}
	if msg := c.ReadHex(); msg != "05000001000000000000" {
		t.Error("Reply(ReplySucceeded) invalid response:", msg)
	}
}

var _ io.ReadWriter = (*TestReadWriter)(nil)
