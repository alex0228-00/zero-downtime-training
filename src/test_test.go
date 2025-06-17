package src

import (
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

/*
 * Zero Downtime Deployment -> Rolling deployment
 */
func TestZeroDowntimeDeployment(t *testing.T) {
	rq := require.New(t)

	docker, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	rq.NoError(err)

	mngr := &TestManager{rq: rq, docker: docker}
	mngr.Init(t.Context())

	// original version
	v1 := NewApp("8081", "v1", docker)
	mngr.DeployFirstVersion(v1)

	// Stage 1: Change database schema
	// v2 -> startup:  add new schema
	//       read:     old schema
	//       write:    both schemas
	v2 := NewApp("8082", "v2", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v2)

	// Stage 2: Database migration
	// v3 -> startup:  data migration
	//       read:     new schema
	//       write:    both schemas
	v3 := NewApp("8083", "v3", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v3)

	// Stage 3: Remove dependency on the old structure
	// v4 -> startup:
	//       read:     new schema
	//       write:    new schemas
	v4 := NewApp("8084", "v4", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v4)

	// Stage 4: clean up
	// v5 -> startup:  remove old field
	//       read:     new schema
	//       write:    new schemas
	v5 := NewApp("8085", "v5", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v5)
}
