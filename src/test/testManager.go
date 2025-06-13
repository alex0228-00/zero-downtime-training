package test

import (
	"fmt"

	"zero-downtime-training/src"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type Tester func(old, new *App, rq *require.Assertions)

type TestManager struct {
	rq      *require.Assertions
	current *App

	tests []Tester
}

func (mngr *TestManager) DeployAppAndTest(app *App) {
	mngr.rq.NoError(app.Deploy())
	mngr.testCURD(app)
}

func (mngr *TestManager) RunZeroDowntimeDeploymentAndTest(newVersion *App) {
	mngr.DeployAppAndTest(newVersion)

	for _, test := range mngr.tests {
		test(mngr.current, newVersion, mngr.rq)
	}

	mngr.rq.NoError(mngr.current.Stop())
	mngr.current = newVersion
}

func (mngr *TestManager) testCURD(app *App) {
	asset := &src.Asset{
		ID:     uuid.New().String(),
		Name:   app.tag,
		Source: fmt.Sprintf("s-%s", app.tag),
	}

	// create
	mngr.rq.NoError(app.client.CreateAsset(asset))

	// read
	read, err := app.client.ReadAsset(asset.ID)
	mngr.rq.NoError(err)
	mngr.rq.EqualValues(asset, read)

	// update
	source := fmt.Sprintf("s-%s-updated", app.tag)
	mngr.rq.NoError(app.client.UpdateAssetSourceByID(asset.ID, source))

	read, err = app.client.ReadAsset(asset.ID)
	mngr.rq.NoError(err)
	mngr.rq.Equal(read.Source, source)

	// delete
	mngr.rq.NoError(app.client.DeleteAsset(asset.ID))
}
