package globaldns

import (
	"fmt"
	"strings"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"k8s.io/api/core/v1"
)

func gatherIngressEndpoints(ingressEps []v1.LoadBalancerIngress) []string {
	endpoints := []string{}
	for _, ep := range ingressEps {
		if ep.IP != "" {
			endpoints = append(endpoints, ep.IP)
		} else if ep.Hostname != "" {
			endpoints = append(endpoints, ep.Hostname)
		}
	}
	return endpoints
}

func getMultiClusterAppName(multiClusterAppName string) (string, error) {
	split := strings.SplitN(multiClusterAppName, ":", 2)
	if len(split) != 2 {
		return "", fmt.Errorf("error in splitting multiclusterapp ID %v", multiClusterAppName)
	}
	mcappName := split[1]
	return mcappName, nil
}

func ifEndpointsDiffer(endpointsOne []string, endpointsTwo []string) bool {
	if len(endpointsOne) != len(endpointsTwo) {
		return true
	}

	mapEndpointsOne := make(map[string]bool)
	for _, ep := range endpointsOne {
		mapEndpointsOne[ep] = true
	}

	for _, ep := range endpointsTwo {
		if !mapEndpointsOne[ep] {
			return true
		}
	}
	return false
}

func reconcileGlobalDNSEndpoints(globalDNS *v3.GlobalDNS) {
	//aggregate all clusterEndpoints and form the final DNS endpoints[]
	var reconciledEps []string
	originalEps := globalDNS.Status.Endpoints

	for _, clusterEndpoints := range globalDNS.Status.ClusterEndpoints {
		reconciledEps = append(reconciledEps, clusterEndpoints...)
	}

	//update the DNS endpoints if different
	if ifEndpointsDiffer(originalEps, reconciledEps) {
		globalDNS.Status.Endpoints = reconciledEps
	}
}
