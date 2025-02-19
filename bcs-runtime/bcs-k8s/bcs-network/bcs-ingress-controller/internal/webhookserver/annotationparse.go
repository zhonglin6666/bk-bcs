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

package webhookserver

import (
	"fmt"
	"strings"
)

const (
	DelimiterSemicolon = ";"
	DelimiterEnter     = "\n"
)

type annotationPort struct {
	poolNamespace string
	poolName      string
	protocol      string
	portIntOrStr  string
}

func parserAnnotation(value string) ([]*annotationPort, error) {
	delimiter := DelimiterSemicolon
	if strings.Contains(value, DelimiterEnter) {
		delimiter = DelimiterEnter
	}
	var retPorts []*annotationPort
	lines := strings.Split(value, delimiter)
	for _, line := range lines {
		realLine := strings.TrimSpace(line)
		fields := strings.Fields(realLine)
		var poolKey string
		var protocol string
		var portStr string
		if len(fields) == 3 {
			poolKey = fields[0]
			protocol = fields[1]
			portStr = fields[2]
		} else if len(fields) == 2 {
			poolKey = fields[0]
			portStr = fields[1]
		} else {
			return nil, fmt.Errorf("annotation line %s is invalid", line)
		}
		var poolName string
		var poolNamespace string
		if strings.Contains(poolKey, ".") {
			poolKeyStrs := strings.Split(poolKey, ".")
			if len(poolKeyStrs) != 2 {
				return nil, fmt.Errorf("invalid poolKey %s", poolKey)
			}
			poolName = poolKeyStrs[0]
			poolNamespace = poolKeyStrs[1]
		} else {
			poolName = poolKey
		}

		if !isProtocolValid(protocol) {
			return nil, fmt.Errorf("protocol %s is invalid", protocol)
		}
		retPorts = append(retPorts, &annotationPort{
			poolNamespace: poolNamespace,
			poolName:      poolName,
			protocol:      protocol,
			portIntOrStr:  portStr,
		})
	}
	return retPorts, nil
}
