// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/rglujing/gf.

package g_test

import (
	"github.com/rglujing/gf/frame/g"
	"github.com/rglujing/gf/net/ghttp"
)

func ExampleServer() {
	// A hello world example.
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}
