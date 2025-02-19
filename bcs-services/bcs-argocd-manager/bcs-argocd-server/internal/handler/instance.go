/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"context"
	actions "github.com/Tencent/bk-bcs/bcs-services/bcs-argocd-manager/bcs-argocd-server/internal/action/instance"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-argocd-manager/pkg/sdk/instance"
)

// InstanceHandler handler that implements the micro handler interface
type InstanceHandler struct{}

// NewInstanceHandler return a new InstanceHandler instance
func NewInstanceHandler() *InstanceHandler {
	return &InstanceHandler{}
}

// CreateArgocdInstance create argocd instance
func (handler *InstanceHandler) CreateArgocdInstance(ctx context.Context,
	request *instance.CreateArgocdInstanceRequest, response *instance.CreateArgocdInstanceResponse) error {
	action := actions.NewCreateArgocdInstanceAction()
	return action.Handle(ctx, request, response)
}

// UpdateArgocdInstance update argocd instance
func (handler *InstanceHandler) UpdateArgocdInstance(ctx context.Context,
	request *instance.UpdateArgocdInstanceRequest, response *instance.UpdateArgocdInstanceResponse) error {
	action := actions.NewUpdateArgocdInstanceAction()
	return action.Handle(ctx, request, response)
}

// DeleteArgocdInstance delete argocd instance
func (handler *InstanceHandler) DeleteArgocdInstance(ctx context.Context,
	request *instance.DeleteArgocdInstanceRequest, response *instance.DeleteArgocdInstanceResponse) error {
	action := actions.NewDeleteArgocdInstanceAction()
	return action.Handle(ctx, request, response)
}

// GetArgocdInstance get argocd instance
func (handler *InstanceHandler) GetArgocdInstance(ctx context.Context,
	request *instance.GetArgocdInstanceRequest, response *instance.GetArgocdInstanceResponse) error {
	action := actions.NewGetArgocdInstanceAction()
	return action.Handle(ctx, request, response)
}

// ListArgocdInstance list argocd instance
func (handler *InstanceHandler) ListArgocdInstances(ctx context.Context,
	request *instance.ListArgocdInstancesRequest, response *instance.ListArgocdInstancesResponse) error {
	action := actions.NewListArgocdInstancesAction()
	return action.Handle(ctx, request, response)
}
