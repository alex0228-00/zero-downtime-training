package src

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type Tester func(old, new *App, rq *require.Assertions)

type TestManager struct {
	docker *DockerDeployer
	logger *zap.Logger
	rq     *require.Assertions

	current *App
	tests   []Tester
}

func NewTestManager(docker *DockerDeployer, rq *require.Assertions) *TestManager {
	return &TestManager{
		docker: docker,
		logger: DefaultLogger.With(
			zap.String("componenet", "TestManager"),
		),
		rq:    rq,
		tests: []Tester{},
	}
}

func (mngr *TestManager) PrepareForTesting() {
	ctx := context.Background()
	mngr.rq.NoError(mngr.docker.CreateNetwork(ctx))
	mngr.rq.NoError(mngr.docker.DeployMysql(ctx))

	mngr.logger.Info("Waiting for MySQL to start...")
	time.Sleep(5 * time.Second)
}

func (mngr *TestManager) DeployFirstVersionAndTest(app *App) {
	mngr.logger.Info(
		"Deploying first version",
		zap.String("version", app.tag),
	)

	mngr.rq.NoError(app.Deploy())
	mngr.testCURD(app)

	mngr.current = app
}

func (mngr *TestManager) RunZeroDowntimeDeploymentAndTest(newVersion *App) {
	mngr.logger.Info(
		"Running zero downtime deployment",
		zap.String("newVersion", newVersion.tag),
	)

	mngr.rq.NoError(newVersion.Deploy())

	mngr.testCURD(newVersion)

	mngr.logger.Info(
		"Running tests for zero downtime deployment",
		zap.String("oldVersion", mngr.current.tag),
		zap.String("newVersion", newVersion.tag),
	)
	for _, test := range mngr.tests {
		test(mngr.current, newVersion, mngr.rq)
	}

	mngr.logger.Info(
		"Stopping old version",
		zap.String("version", mngr.current.tag),
	)
	mngr.rq.NoError(mngr.current.Stop())
	mngr.current = newVersion
}

func (mngr *TestManager) testCURD(app *App) {
	mngr.logger.Info("Running CURD tests", zap.String("version", app.tag))

	asset := &Asset{
		ID:     uuid.New().String(),
		Name:   app.tag,
		Source: fmt.Sprintf("source-%s", app.tag),
	}

	// create
	mngr.logger.Info(
		"Test creating new asset",
		zap.String("version", app.tag),
		zap.String("asset id", asset.ID),
	)
	created, err := app.client.CreateAsset(asset)
	mngr.rq.NoError(err)

	// read
	mngr.logger.Info(
		"Test reading asset",
		zap.String("version", app.tag),
		zap.String("asset id", created.ID),
	)
	read, err := app.client.ReadAsset(created.ID)
	mngr.rq.NoError(err)
	mngr.rq.EqualValues(created, read)

	// update
	source := fmt.Sprintf("source-%s-updated", app.tag)
	mngr.logger.Info(
		"Test updating asset source",
		zap.String("version", app.tag),
		zap.String("asset id", created.ID),
		zap.String("new source", source),
	)
	mngr.rq.NoError(app.client.UpdateAssetSourceByID(created.ID, source))

	read, err = app.client.ReadAsset(created.ID)
	mngr.rq.NoError(err)
	mngr.rq.Equal(read.Source, source)

	// delete
	mngr.logger.Info(
		"Test deleting asset",
		zap.String("version", app.tag),
		zap.String("asset id", created.ID),
	)
	mngr.rq.NoError(app.client.DeleteAsset(created.ID))
}
