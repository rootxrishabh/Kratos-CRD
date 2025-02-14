package cloud

import (
	"context"
	"fmt"

	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"

	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "operator.kratos.io/kratos/api/v1"
)

// createCluster creates a new GKE cluster
func CreateCluster(ctx context.Context, gkeService *container.Service, kratos *operatorv1.Kratos) error {
	log := log.FromContext(ctx)
	log.Info("Creating GKE cluster", "ClusterName", kratos.Spec.ClusterName)

	req := &container.CreateClusterRequest{
		Cluster: &container.Cluster{
			Name:             kratos.Spec.ClusterName,
			InitialNodeCount: int64(kratos.Spec.NodePools[0].NodeCount),
			Location:         kratos.Spec.Region,
			NetworkConfig:    &container.NetworkConfig{Network: kratos.Spec.Networking.VPCName, Subnetwork: kratos.Spec.Networking.SubnetName},
		},
	}

	_, err := gkeService.Projects.Locations.Clusters.Create(
		fmt.Sprintf("projects/%s/locations/%s", kratos.Spec.ProjectID, kratos.Spec.Region), req).Context(ctx).Do()

	if err != nil {
		return err
	}

	log.Info("GKE cluster created successfully")
	return nil
}

// updateCluster updates an existing GKE cluster
func UpdateCluster(ctx context.Context, gkeService *container.Service, kratos *operatorv1.Kratos) error {
	log := log.FromContext(ctx)
	log.Info("Updating GKE cluster", "ClusterName", kratos.Spec.ClusterName)

	_, err := gkeService.Projects.Locations.Clusters.NodePools.SetSize(
		fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/default-pool",
			kratos.Spec.ProjectID, kratos.Spec.Region, kratos.Spec.ClusterName),
		&container.SetNodePoolSizeRequest{NodeCount: int64(kratos.Spec.NodePools[0].NodeCount)}).Context(ctx).Do()

	if err != nil {
		return err
	}

	log.Info("GKE cluster updated successfully")
	return nil
}

// deleteCluster deletes an existing GKE cluster
func DeleteCluster(ctx context.Context, gkeService *container.Service, kratos *operatorv1.Kratos) error {
	log := log.FromContext(ctx)
	log.Info("Deleting GKE cluster", "ClusterName", kratos.Spec.ClusterName)

	_, err := gkeService.Projects.Locations.Clusters.Delete(
		fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
			kratos.Spec.ProjectID, kratos.Spec.Region, kratos.Spec.ClusterName)).Context(ctx).Do()

	if err != nil {
		return err
	}

	log.Info("GKE cluster deleted successfully")
	return nil
}

// getGKEClient initializes the GKE API client
func GetGKEClient(serviceAccountKey string) (*container.Service, error) {
	ctx := context.Background()
	return container.NewService(ctx, option.WithCredentialsJSON([]byte(serviceAccountKey)))
}

// clusterExists checks if the GKE cluster already exists
func ClusterExists(ctx context.Context, gkeService *container.Service, kratos *operatorv1.Kratos) (bool, error) {
	_, err := gkeService.Projects.Locations.Clusters.Get(fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		kratos.Spec.ProjectID, kratos.Spec.Region, kratos.Spec.ClusterName)).Context(ctx).Do()

	if err != nil {
		if isNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isNotFoundError(err error) bool {
	return err != nil && err.Error() == "not found"
}
