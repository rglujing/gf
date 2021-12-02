// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/rglujing/gf.

package gdb_test

import (
	"context"

	"github.com/rglujing/gf/database/gdb"
	"github.com/rglujing/gf/frame/g"
)

func Example_transaction() {
	db.Transaction(context.TODO(), func(ctx context.Context, tx *gdb.TX) error {
		// user
		result, err := tx.Insert("user", g.Map{
			"passport": "john",
			"password": "12345678",
			"nickname": "JohnGuo",
		})
		if err != nil {
			return err
		}
		// user_detail
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		_, err = tx.Insert("user_detail", g.Map{
			"uid":       id,
			"site":      "https://johng.cn",
			"true_name": "GuoQiang",
		})
		if err != nil {
			return err
		}
		return nil
	})
}
