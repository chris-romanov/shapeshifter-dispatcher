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
	"fmt"

	"github.com/OperatorFoundation/shapeshifter-ipc/v2"
)

func (req *Request) authPT2() (err error) {
	// The client sends a PT 2.0 authentication request.
	//  uint32_t u
	//  uint8_t data[u]

	// Read the authentication data.
	var u uint32
	if u, err = req.readUint32(); err != nil {
		return
	}
	if u == 0 {
		err = fmt.Errorf("PT 2.0 authentication data with 0 length")
		return
	}
	var data []byte
	if data, err = req.readBytes(int(u)); err != nil {
		return
	}

	var result = string(data)

	// Parse the authentication data according to the PT 2.0 specification
	if req.Args, err = pt.ParsePT2ClientParameters(result); err != nil {
		fmt.Println("Error parsing PT2 client parameters", err)
		return
	}

	return
}
