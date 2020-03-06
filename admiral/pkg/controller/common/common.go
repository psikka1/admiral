package common

import (
	"github.com/istio-ecosystem/admiral/admiral/pkg/apis/admiral/v1"
	log "github.com/sirupsen/logrus"
	k8sAppsV1 "k8s.io/api/apps/v1"
	k8sV1 "k8s.io/api/core/v1"
	"sort"
	"strings"
)

const (
	NamespaceKubeSystem        = "kube-system"
	NamespaceIstioSystem       = "istio-system"
	Env                        = "env"
	Http                       = "http"
	DefaultMtlsPort            = 15443
	DefaultHttpPort            = 80
	Sep                        = "."
	Dash                       = "-"
	Slash                      = "/"
	DotLocalDomainSuffix       = ".svc.cluster.local"
	Mesh                       = "mesh"
	MulticlusterIngressGateway = "istio-multicluster-ingressgateway"
	LocalAddressPrefix         = "240.0"
	NodeRegionLabel            = "failure-domain.beta.kubernetes.io/region"
	SpiffePrefix               = "spiffe://"
	SidecarEnabledPorts        = "traffic.sidecar.istio.io/includeInboundPorts"
	Default                    = "default"
)

type Event int

const (
	Add    Event = 0
	Update Event = 1
	Delete Event = 2
)

type ResourceType string

const (
	VirtualService  ResourceType = "VirtualService"
	DestinationRule  ResourceType = "DestinationRule"
	ServiceEntry  ResourceType = "ServiceEntry"
)

func GetPodGlobalIdentifier(pod *k8sV1.Pod) string {
	identity := pod.Labels[GetWorkloadIdentifier()]
	if len(identity) == 0 {
		identity = pod.Annotations[GetWorkloadIdentifier()]
	}
	return identity
}

func GetDeploymentGlobalIdentifier(deployment *k8sAppsV1.Deployment) string {
	identity := deployment.Spec.Template.Labels[GetWorkloadIdentifier()]
	if len(identity) == 0 {
		//TODO can this be removed now? This was for backward compatibility
		identity = deployment.Spec.Template.Annotations[GetWorkloadIdentifier()]
	}
	return identity
}

// GetCname returns cname in the format <env>.<service identity>.global, Ex: stage.Admiral.services.registry.global
func GetCname(deployment *k8sAppsV1.Deployment, identifier string, nameSuffix string) string {
	var environment = GetEnv(deployment)
	alias := GetValueForKeyFromDeployment(identifier, deployment)
	if len(alias) == 0 {
		log.Warnf("%v label missing on deployment %v in namespace %v. Falling back to annotation to create cname.", identifier, deployment.Name, deployment.Namespace)
		alias = deployment.Spec.Template.Annotations[identifier]
	}
	if len(alias) == 0 {
		log.Errorf("Unable to get cname for deployment with name %v in namespace %v as it doesn't have the %v annotation", deployment.Name, deployment.Namespace, identifier)
		return ""
	}
	return environment + Sep + alias + Sep + nameSuffix
}

func GetEnv(deployment *k8sAppsV1.Deployment) string {
	var environment = deployment.Spec.Template.Labels[Env]
	if len(environment) == 0 {
		environment = deployment.Spec.Template.Annotations[Env]
	}
	if len(environment) == 0 {
		splitNamespace := strings.Split(deployment.Namespace, Dash)
		if len(splitNamespace) > 1 {
			environment = splitNamespace[len(splitNamespace)-1]
		}
	}
	if len(environment) == 0 {
		environment = Default
	}
	return environment
}

// GetSAN returns SAN for a service entry in the format spiffe://<domain>/<identifier>, Ex: spiffe://subdomain.domain.com/Admiral.platform.mesh.server
func GetSAN(domain string, deployment *k8sAppsV1.Deployment, identifier string) string {
	identifierVal := GetValueForKeyFromDeployment(identifier, deployment)
	if len(identifierVal) == 0 {
		log.Errorf("Unable to get SAN for deployment with name %v in namespace %v as it doesn't have the %v annotation or label", deployment.Name, deployment.Namespace, identifier)
		return ""
	}
	if len(domain) > 0 {
		return SpiffePrefix + domain + Slash + identifierVal
	} else {
		return SpiffePrefix + identifierVal
	}
}

func GetNodeLocality(node *k8sV1.Node) string {
	region, _ := node.Labels[NodeRegionLabel]
	return region
}

func GetValueForKeyFromDeployment(key string, deployment *k8sAppsV1.Deployment) string {
	value := deployment.Spec.Template.Labels[key]
	if len(value) == 0 {
		log.Warnf("%v label missing on deployment %v in namespace %v. Falling back to annotation.", key, deployment.Name, deployment.Namespace)
		value = deployment.Spec.Template.Annotations[key]
	}
	return value
}

func MatchDeploymentsToGTP(gtp *v1.GlobalTrafficPolicy, deployments []k8sAppsV1.Deployment) *k8sAppsV1.Deployment{
	if gtp == nil || gtp.Name == "" {
		log.Warn("Nil or empty GlobalTrafficPolicy provided for deployment match. Returning nil.")
		return nil
	}

	if len(deployments) == 0 {
		return nil
	}

	//If one is found, return it.
	if len(deployments) == 1 {
		return &deployments[0]
	}

	var envMatchedDeployments []k8sAppsV1.Deployment

	for _, deployment := range deployments {
		if deployment.Spec.Template.Labels[Env] == gtp.Labels[Env] {
			envMatchedDeployments = append(envMatchedDeployments, deployment)
		}
	}

	//if one matches the environment from the gtp, return it
	if len(envMatchedDeployments) == 1 {
		return &envMatchedDeployments[0]
	}

	//if no deployments match the environment, we follow the same logic as if multiple did.
	if len(envMatchedDeployments) == 0 {
		envMatchedDeployments = deployments
	}

	sort.Slice(envMatchedDeployments, func(i, j int) bool {
		iTime := envMatchedDeployments[i].CreationTimestamp.Nanosecond()
		jTime := envMatchedDeployments[j].CreationTimestamp.Nanosecond()
		return iTime<jTime
	})

	//return most recently created gtp
	return &envMatchedDeployments[0]

}

func MatchGTPsToDeployment(gtpList []v1.GlobalTrafficPolicy, deployment *k8sAppsV1.Deployment) *v1.GlobalTrafficPolicy{
	if deployment == nil || deployment.Name == "" {
		log.Warn("Nil or empty GlobalTrafficPolicy provided for deployment match. Returning nil.")
		return nil
	}

	//If one is found, return it.
	if len(gtpList) == 1 {
		return &gtpList[0]
	}

	if len(gtpList) == 0 {
		return nil
	}

	var envMatchedGTPList []v1.GlobalTrafficPolicy

	for _, gtp := range gtpList {
		if gtp.Labels[Env] == deployment.Spec.Template.Labels[Env] {
			envMatchedGTPList = append(envMatchedGTPList, gtp)
		}
	}

	//if one matches the environment from the gtp, return it
	if len(envMatchedGTPList) == 1 {
		return &envMatchedGTPList[0]
	}

	//if no GTPs match the environment, we follow the same logic as if multiple did.
	if len(envMatchedGTPList) == 0 {
		envMatchedGTPList = gtpList
	}

	sort.Slice(envMatchedGTPList, func(i, j int) bool {
		iTime := envMatchedGTPList[i].CreationTimestamp.Nanosecond()
		jTime := envMatchedGTPList[j].CreationTimestamp.Nanosecond()
		return iTime<jTime
	})

	//return most recently created gtp
	return &envMatchedGTPList[0]

}