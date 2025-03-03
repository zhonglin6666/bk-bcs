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

package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	"github.com/Tencent/bk-bcs/bcs-common/pkg/odm/drivers"
	"github.com/Tencent/bk-bcs/bcs-common/pkg/odm/operator"

	cmproto "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/api/clustermanager"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/common"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/remote/auth"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/store"
	storeopt "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/store/options"
)

// ListAction list action for cluster
type ListAction struct {
	ctx         context.Context
	model       store.ClusterManagerModel
	req         *cmproto.ListClusterReq
	resp        *cmproto.ListClusterResp
	clusterList []*cmproto.Cluster
}

// NewListAction create list action for cluster
func NewListAction(model store.ClusterManagerModel) *ListAction {
	return &ListAction{
		model: model,
	}
}

func (la *ListAction) validate() error {
	if err := la.req.Validate(); err != nil {
		return err
	}

	if len(la.req.Environment) > 0 {
		_, ok := EnvironmentLookup[la.req.Environment]
		if !ok {
			return fmt.Errorf("request Environment invalid, must be [test/debug/prod]")
		}
	}

	if len(la.req.EngineType) > 0 {
		_, ok := EngineTypeLookup[la.req.EngineType]
		if !ok {
			return fmt.Errorf("request EngineType invalid, must be [k8s/mesos]")
		}
	}

	if len(la.req.ClusterType) > 0 {
		_, ok := ClusterTypeLookup[la.req.ClusterType]
		if !ok {
			return fmt.Errorf("request ClusterType invalid, must be [federation/signal]")
		}
	}

	return nil
}

func (la *ListAction) getSharedCluster() error {
	condM := make(operator.M)
	condM["isshared"] = true

	condCluster := operator.NewLeafCondition(operator.Eq, condM)
	clusterList, err := la.model.ListCluster(la.ctx, condCluster, &storeopt.ListOption{})
	if err != nil && !errors.Is(err, drivers.ErrTableRecordNotFound) {
		return err
	}

	clusterIDs := make([]string, 0)
	for i := range clusterList {
		la.clusterList = append(la.clusterList, &clusterList[i])
		clusterIDs = append(clusterIDs, clusterList[i].ClusterID)
	}

	if la.resp.ClusterPerm == nil {
		la.resp.ClusterPerm = make(map[string]*cmproto.Permission)
	}

	for _, clusterID := range clusterIDs {
		if _, ok := la.resp.ClusterPerm[clusterID]; !ok {
			la.resp.ClusterPerm[clusterID] = &cmproto.Permission{
				Policy: auth.GetSharedClusterPerm(),
			}
		}
	}

	return nil
}

func (la *ListAction) listCluster() error {
	getSharedCluster := true
	condM := make(operator.M)
	if len(la.req.ClusterName) != 0 {
		condM["clustername"] = la.req.ClusterName
	}
	if len(la.req.Provider) != 0 {
		condM["provider"] = la.req.Provider
	}
	if len(la.req.Region) != 0 {
		condM["region"] = la.req.Region
	}
	if len(la.req.VpcID) != 0 {
		condM["vpcid"] = la.req.VpcID
	}
	if len(la.req.ProjectID) != 0 {
		condM["projectid"] = la.req.ProjectID
	}
	if len(la.req.BusinessID) != 0 {
		condM["businessid"] = la.req.BusinessID
	}
	if len(la.req.Environment) != 0 {
		condM["environment"] = la.req.Environment
	}
	if len(la.req.EngineType) != 0 {
		condM["enginetype"] = la.req.EngineType
	}
	if len(la.req.SystemID) != 0 {
		condM["systemid"] = la.req.SystemID
	}
	if len(la.req.ExtraClusterID) != 0 {
		condM["extraclusterid"] = la.req.ExtraClusterID
		getSharedCluster = false
	}
	if len(la.req.ClusterID) != 0 {
		condM["clusterid"] = la.req.ClusterID
	}

	condCluster := operator.NewLeafCondition(operator.Eq, condM)
	condStatus := operator.NewLeafCondition(operator.Ne, operator.M{"status": common.StatusDeleted})

	branchCond := operator.NewBranchCondition(operator.And, condCluster, condStatus)
	clusterList, err := la.model.ListCluster(la.ctx, branchCond, &storeopt.ListOption{})
	if err != nil && !errors.Is(err, drivers.ErrTableRecordNotFound) {
		return err
	}

	clusterIDList := make([]string, 0)
	for i := range clusterList {
		if clusterList[i].IsShared {
			clusterList[i].IsShared = false
		}
		la.clusterList = append(la.clusterList, &clusterList[i])
		clusterIDList = append(clusterIDList, clusterList[i].ClusterID)
	}

	// return cluster extraInfo
	la.returnClusterExtraInfo(clusterList)

	// projectID / operator get user perm
	if la.req.ProjectID != "" && la.req.Operator != "" {
		perms, err := GetUserPermListByProjectAndCluster(UserInfo{
			ProjectID: la.req.ProjectID,
			UserID:    la.req.Operator,
		}, clusterIDList, false)
		if err != nil {
			blog.Errorf("listCluster GetUserPermListByProjectAndCluster failed: %v", err.Error())
			return err
		}

		la.resp.ClusterPerm = perms

		createPerm := GetClusterCreatePerm(UserInfo{
			ProjectID: la.req.ProjectID,
			UserID:    la.req.Operator,
		})
		la.resp.Permissions = createPerm
	}

	// default return shared cluster
	if getSharedCluster {
		err = la.getSharedCluster()
		if err != nil {
			blog.Errorf("ListCluster getSharedCluster failed: %v", err)
			return err
		}
	}

	return nil
}

// GetCloudProviderEngine get cloud engineType
func (la *ListAction) GetCloudProviderEngine(cls cmproto.Cluster) string {
	cloud, err := la.model.GetCloud(la.ctx, cls.Provider)
	if err != nil {
		blog.Errorf("listCluster GetCloudProviderEngine failed: %v", err)
		return ""
	}

	return cloud.GetEngineType()
}

func (la *ListAction) returnClusterExtraInfo(clusterList []cmproto.Cluster) {
	if la.resp.ClusterExtraInfo == nil {
		la.resp.ClusterExtraInfo = make(map[string]*cmproto.ExtraInfo)
	}

	// cluster extra info
	for i := range clusterList {
		cond := operator.NewLeafCondition(operator.Eq, operator.M{"clusterid": clusterList[i].ClusterID})
		nodes, err := la.model.ListNode(la.ctx, cond, &storeopt.ListOption{})
		if err != nil {
			continue
		}
		canDeleted := func() bool {
			if len(nodes) > 0 {
				return false
			}
			return true
		}()

		la.resp.ClusterExtraInfo[clusterList[i].ClusterID] = &cmproto.ExtraInfo{
			CanDeleted:   canDeleted,
			ProviderType: la.GetCloudProviderEngine(clusterList[i]),
		}
	}

	return
}

func (la *ListAction) setResp(code uint32, msg string) {
	la.resp.Code = code
	la.resp.Message = msg
	la.resp.Result = (code == common.BcsErrClusterManagerSuccess)
	la.resp.Data = la.clusterList
}

// Handle handle list cluster request
func (la *ListAction) Handle(ctx context.Context, req *cmproto.ListClusterReq, resp *cmproto.ListClusterResp) {
	if req == nil || resp == nil {
		blog.Errorf("list cluster failed, req or resp is empty")
		return
	}
	la.ctx = ctx
	la.req = req
	la.resp = resp

	if err := la.validate(); err != nil {
		la.setResp(common.BcsErrClusterManagerInvalidParameter, err.Error())
		return
	}
	if err := la.listCluster(); err != nil {
		la.setResp(common.BcsErrClusterManagerDBOperation, err.Error())
		return
	}
	la.setResp(common.BcsErrClusterManagerSuccess, common.BcsErrClusterManagerSuccessStr)
	return
}

// ListNodesInClusterAction list action for cluster
type ListNodesInClusterAction struct {
	ctx   context.Context
	model store.ClusterManagerModel
	req   *cmproto.ListNodesInClusterRequest
	resp  *cmproto.ListNodesInClusterResponse
	nodes []*cmproto.Node
}

// NewListNodesInClusterAction create list action for cluster
func NewListNodesInClusterAction(model store.ClusterManagerModel) *ListNodesInClusterAction {
	return &ListNodesInClusterAction{
		model: model,
	}
}

func (la *ListNodesInClusterAction) validate() error {
	if err := la.req.Validate(); err != nil {
		return err
	}
	return nil
}

func (la *ListNodesInClusterAction) listNodes() error {
	condM := make(operator.M)
	if len(la.req.ClusterID) != 0 {
		condM["clusterid"] = la.req.ClusterID
	}
	if len(la.req.Region) != 0 {
		condM["region"] = la.req.Region
	}
	if len(la.req.VpcID) != 0 {
		condM["vpcid"] = la.req.VpcID
	}
	if len(la.req.NodeGroupID) != 0 {
		condM["nodegroupid"] = la.req.NodeGroupID
	}
	if len(la.req.InstanceType) != 0 {
		condM["instancetype"] = la.req.InstanceType
	}
	if len(la.req.Status) != 0 {
		condM["status"] = la.req.Status
	}
	cond := operator.NewLeafCondition(operator.Eq, condM)
	nodes, err := la.model.ListNode(la.ctx, cond, &storeopt.ListOption{})
	if err != nil {
		blog.Errorf("list nodes in cluster %s failed, %s", la.req.ClusterID, err.Error())
		return err
	}
	for i := range nodes {
		if !la.req.ShowPwd {
			nodes[i].Passwd = ""
		}
		la.nodes = append(la.nodes, nodes[i])
	}

	return nil
}

func (la *ListNodesInClusterAction) setResp(code uint32, msg string) {
	la.resp.Code = code
	la.resp.Message = msg
	la.resp.Result = (code == common.BcsErrClusterManagerSuccess)
	la.resp.Data = la.nodes
}

// Handle handle list cluster request
func (la *ListNodesInClusterAction) Handle(ctx context.Context,
	req *cmproto.ListNodesInClusterRequest, resp *cmproto.ListNodesInClusterResponse) {
	if req == nil || resp == nil {
		blog.Errorf("list cluster nodes failed, req or resp is empty")
		return
	}
	la.ctx = ctx
	la.req = req
	la.resp = resp

	if err := la.validate(); err != nil {
		la.setResp(common.BcsErrClusterManagerInvalidParameter, err.Error())
		return
	}
	if err := la.listNodes(); err != nil {
		la.setResp(common.BcsErrClusterManagerDBOperation, err.Error())
		return
	}
	la.setResp(common.BcsErrClusterManagerSuccess, common.BcsErrClusterManagerSuccessStr)
	return
}
