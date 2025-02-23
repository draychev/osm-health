package main

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/connectivity"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
)

const connectivityPodToURLDesc = `
Checks connectivity between a Kubernetes pod and a host name (or URL)
`

const connectivityPodToURLExample = `$ osm-health connectivity pod-to-url namespace-a/pod-a https://contoso.com/store`

func newConnectivityPodToURLCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "pod-to-url SOURCE_POD DESTINATION_URL",
		Short:   "Checks connectivity between a Kubernetes pod and a given URL",
		Example: connectivityPodToURLExample,
		Long:    connectivityPodToURLDesc,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.Errorf("provide both SOURCE_POD and DESTINATION_URL")
			}

			srcPod, err := pod.FromString(args[0])
			if err != nil {
				return errors.New("invalid SOURCE_POD")
			}

			dstURL, err := url.Parse(args[1])
			if err != nil {
				return errors.New("invalid DESTINATION_URL")
			}

			connectivity.PodToURL(srcPod, dstURL, settings.Namespace())
			return nil
		},
	}
}
