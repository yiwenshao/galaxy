/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
package galaxy

import (
	t020 "github.com/containernetworking/cni/pkg/types/020"
	"github.com/containernetworking/cni/pkg/types/current"
	"gotest.tools/assert"
	"net"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/galaxy/pkg/api/galaxy/constant"
)

func TestParseExtendedCNIArgs(t *testing.T) {
	m, err := parseExtendedCNIArgs(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{
		constant.ExtendedCNIArgsAnnotation: `{"request_ip_range":[["10.0.0.2~10.0.0.30"],["10.0.0.200~10.0.0.238"]],"common":{"ipinfos":[{"ip":"10.0.0.3/24","vlan":0,"gateway":"10.0.0.1"},{"ip":"10.0.0.200/24","vlan":0,"gateway":"10.0.0.1"}]}}`,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if val, ok := m["ipinfos"]; !ok {
		t.Fatal()
	} else if string(val) != `[{"ip":"10.0.0.3/24","vlan":0,"gateway":"10.0.0.1"},{"ip":"10.0.0.200/24","vlan":0,"gateway":"10.0.0.1"}]` {
		t.Fatal()
	}
}

func TestConvertResult(t *testing.T) {
	t1 := &t020.Result{
		IP4: &t020.IPConfig{
			IP: net.IPNet{
				IP: net.IPv4(127, 0, 0, 1),
			},
		},
	}
	res, _ := convertResult(t1)
	assert.Equal(t, res.IP4.IP.IP.String(), "127.0.0.1")

	t2 := &current.Result{
		IPs: []*current.IPConfig{
			{
				Version: "4",
				Address: net.IPNet{
					IP: net.IPv4(127, 0, 0, 1),
				},
			},
		},
	}
	res, _ = convertResult(t2)
	assert.Equal(t, res.IP4.IP.IP.String(), "127.0.0.1")

}
