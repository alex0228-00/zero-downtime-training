package src

import (
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

/*
 * Zero Downtime Deployment -> Rolling deployment
 *
 * Assets                                    Assets
 * +----+------+--------+                    +----+------+-----------+
 * | id | name | source |     migration      | id | name | source id |
 * +----+------+--------+     -------->      +----+------+-----------+
 *                                           Sources
 * 										     +-----------+------+
 *                                           | source id | name |
 *                                           +-----------+------+
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
	//
	// Assets
	// +----+------+--------+
	// | id | name | source |
	// +----+------+--------+
	//                 RW
	//
	v1 := NewApp("8081", "v1", docker)
	mngr.DeployFirstVersion(v1)

	// Stage 1: Change database schema
	//
	// Assets
	//						add new column
	// +----+------+--------+-----------+
	// | id | name | source | source id |
	// +----+------+--------+-----------+
	//              RW(v1/v2)    W(v2)
	//
	// Sources
	// +-----------+------+
	// | source id | name |
	// +-----------+------+
	//     W(v2)     W(v2)
	//
	v2 := NewApp("8082", "v2", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v2)

	// Stage 2: Database migration
	//
	// Assets
	//					      migration
	//                  ----------------------
	//                  |                    |
	//                  |                    |
	// +----+------+--------+-----------+    |
	// | id | name | source | source id |  <-|
	// +----+------+--------+-----------+    |
	//               R(v2)      R(v3)        |
	//              W(v2/v3)   W(v2/v3)      |
	//                                       |
	// Sources                               |
	// +-----------+------+                  |
	// | source id | name |       <-----------
	// +-----------+------+
	//    R(v3)      R(v3)
	//   W(v2/v3)   W(v2/v3)
	//
	v3 := NewApp("8083", "v3", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v3)

	// Stage 3: Remove dependency on the old structure
	//
	// Assets
	// +----+------+--------+-----------+
	// | id | name | source | source id |
	// +----+------+--------+-----------+
	//               W(v3)    RW(v3/v4)
	// Sources
	// +-----------+------+
	// | source id | name |
	// +-----------+------+
	//   RW(v3/v4)  RW(v3/v4)
	//
	v4 := NewApp("8084", "v4", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v4)

	// Stage 4: clean up
	//
	// Assets
	//             remove column
	//             +--------+
	//             | source |
	//             +--------+
	//             |        |
	// +----+------+        +-----------+
	// | id | name |        | source id |
	// +----+------+        +-----------+
	//                        RW(v4/v5)
	// Sources
	// +-----------+------+
	// | source id | name |
	// +-----------+------+
	//   RW(v4/v5)  RW(v4/v5)
	//
	v5 := NewApp("8085", "v5", docker)
	mngr.RunZeroDowntimeDeploymentAndTest(v5)
}
