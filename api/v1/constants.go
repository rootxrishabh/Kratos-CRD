package v1

// ClusterPhase represents the phase of cluster provisioning
type ClusterPhase string

const (
	PhasePending  ClusterPhase = "Pending"  // Cluster creation is pending
	PhaseCreating ClusterPhase = "Creating" // Cluster is being created
	PhaseRunning  ClusterPhase = "Running"  // Cluster is up and running
	PhaseFailed   ClusterPhase = "Failed"   // Cluster creation failed
	PhaseDeleting ClusterPhase = "Deleting" // Cluster is being deleted
	PhaseDeleted  ClusterPhase = "Deleted"  // Cluster is deleted
)
