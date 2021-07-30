package connectivity

import (
	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm-health/pkg/osm"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")
	osmVersion := osm.ControllerVersion("v0.9")

	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Err(err).Msg("Error creating Kubernetes client")
	}

	var srcConfigGetter, dstConfigGetter envoy.ConfigGetter

	srcConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(fromPod)
	if err != nil {
		log.Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", fromPod.Namespace, fromPod.Name)
	}

	dstConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(toPod)
	if err != nil {
		log.Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", toPod.Namespace, toPod.Name)
	}

	outcomes := common.Run(
		// Check source Pod's namespace
		namespace.IsInjectEnabled(client, fromPod.Namespace),
		namespace.IsMonitoredBy(client, fromPod.Namespace, meshName),
		pod.HasEnvoySidecar(fromPod),

		// Check destination Pod's namespace
		namespace.IsInjectEnabled(client, toPod.Namespace),
		namespace.IsMonitoredBy(client, toPod.Namespace, meshName),
		pod.HasEnvoySidecar(toPod),

		envoy.HasListener(srcConfigGetter, osmVersion),
		envoy.HasListener(dstConfigGetter, osmVersion),
	)

	common.Print(outcomes...)

	return common.Result{
		SMIPolicy: common.SMIPolicy{
			HasPolicy:                  false,
			ValidPolicy:                false,
			SourceToDestinationAllowed: false,
		},
		Successful: false,
	}
}
