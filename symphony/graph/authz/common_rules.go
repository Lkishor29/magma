// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
)

func cudBasedCheck(cud *models.Cud, m ent.Mutation) bool {
	var permission *models.BasicPermissionRule
	switch {
	case m.Op().Is(ent.OpCreate):
		permission = cud.Create
	case m.Op().Is(ent.OpUpdateOne | ent.OpUpdate):
		permission = cud.Update
	case m.Op().Is(ent.OpDeleteOne | ent.OpDelete):
		permission = cud.Delete
	default:
		return false
	}
	return permission.IsAllowed == models2.PermissionValueYes
}

func cudBasedRule(cud *models.Cud, m ent.Mutation) error {
	if cudBasedCheck(cud, m) {
		return privacy.Allow
	}
	return privacy.Skip
}

// AllowWritePermissionsRule grants write permission.
func AllowWritePermissionsRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		if FromContext(ctx).CanWrite {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AlwaysDenyIfNoPermissionRule denies access if no permissions is present on context.
func AlwaysDenyIfNoPermissionRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, _ ent.Mutation) error {
		if FromContext(ctx) == nil {
			return privacy.Deny
		}
		return privacy.Skip
	})
}
