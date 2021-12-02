// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/rglujing/gf.

package gview_test

import (
	"context"
	"github.com/rglujing/gf/debug/gdebug"
	"github.com/rglujing/gf/frame/g"
	"github.com/rglujing/gf/os/gfile"
	"github.com/rglujing/gf/os/gview"
	"github.com/rglujing/gf/test/gtest"
	"testing"
)

func Test_Encode_Parse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetPath(gdebug.TestDataPath("tpl"))
		v.SetAutoEncode(true)
		result, err := v.Parse(context.TODO(), "encode.tpl", g.Map{
			"title": "<b>my title</b>",
		})
		t.Assert(err, nil)
		t.Assert(result, "<div>&lt;b&gt;my title&lt;/b&gt;</div>")
	})
}

func Test_Encode_ParseContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		tplContent := gfile.GetContents(gdebug.TestDataPath("tpl", "encode.tpl"))
		v.SetAutoEncode(true)
		result, err := v.ParseContent(context.TODO(), tplContent, g.Map{
			"title": "<b>my title</b>",
		})
		t.Assert(err, nil)
		t.Assert(result, "<div>&lt;b&gt;my title&lt;/b&gt;</div>")
	})
}
