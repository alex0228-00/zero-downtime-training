package src

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type Tester func(old, new *App, rq *require.Assertions)

type TestManager struct {
	docker           *client.Client
	networkID        string
	mysqlContainreID string

	rq      *require.Assertions
	current *App

	tests []Tester
}

func (mngr *TestManager) Init(ctx context.Context) {
	mngr.createDockerNetworkIfNotExist(ctx)
	mngr.deployMysql(ctx)
}

func (mngr *TestManager) createDockerNetworkIfNotExist(ctx context.Context) {
	networks, err := mngr.docker.NetworkList(ctx, network.ListOptions{})
	mngr.rq.NoError(err)

	for _, nw := range networks {
		if nw.Name == Network {
			mngr.networkID = nw.ID
			return
		}
	}

	resp, err := mngr.docker.NetworkCreate(ctx, Network, network.CreateOptions{})
	mngr.rq.NoError(err)

	mngr.networkID = resp.ID
}

func (mngr *TestManager) deployMysql(ctx context.Context) {
	containers, err := mngr.docker.ContainerList(ctx, container.ListOptions{})
	mngr.rq.NoError(err)

	for _, c := range containers {
		if strings.Contains(c.Names[0], MysqlContainerName) {
			mngr.mysqlContainreID = c.ID
			return
		}
	}

	mngr.rq.NoError(mngr.pullImageIfNotExist(ctx, "mysql:latest"))

	response, err := mngr.docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "mysql:latest",
			ExposedPorts: nat.PortSet{
				"3306/tcp": struct{}{},
			},
			Env: []string{
				EncodeDockerEnv("MYSQL_ROOT_PASSWORD", "rootpwd"),
				EncodeDockerEnv("MYSQL_DATABASE", DbSchema),
				EncodeDockerEnv("MYSQL_USER", DbUser),
				EncodeDockerEnv("MYSQL_PASSWORD", DbPassword),
			},
		},
		nil,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				Network: {},
			},
		},
		nil,
		MysqlContainerName,
	)
	mngr.rq.NoError(err)

	mngr.rq.NoError(mngr.docker.ContainerStart(ctx, response.ID, container.StartOptions{}))
	mngr.mysqlContainreID = response.ID
}

func (mngr *TestManager) pullImageIfNotExist(ctx context.Context, imageName string) error {
	images, err := mngr.docker.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	found := false
	for _, img := range images {
		if img.RepoTags[0] == imageName {
			found = true
		}
	}

	if !found {
		reader, err := mngr.docker.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image %s: %w", imageName, err)
		}
		defer reader.Close()

		_, err = io.Copy(io.Discard, reader)
		return fmt.Errorf("failed to read image pull response: %w", err)
	}
	return nil
}

func (mngr *TestManager) DeployFirstVersion(app *App) {
	mngr.rq.NoError(app.Deploy())
	mngr.testCURD(app)

	mngr.current = app
}

func (mngr *TestManager) RunZeroDowntimeDeploymentAndTest(newVersion *App) {
	mngr.rq.NoError(newVersion.Deploy())
	mngr.testCURD(newVersion)

	for _, test := range mngr.tests {
		test(mngr.current, newVersion, mngr.rq)
	}

	mngr.rq.NoError(mngr.current.Stop())
	mngr.current = newVersion
}

func (mngr *TestManager) testCURD(app *App) {
	asset := &Asset{
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
