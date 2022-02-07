package mysqlstore

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/suite"
	"os"
	"time"
)

type MysqlCfg struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

type MysqlStoreTestSuite struct {
	suite.Suite
	s      *MysqlStore
	docker *client.Client
	dbId   string
	dbCfg  *MysqlCfg
}

func (suite *MysqlStoreTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.dbCfg = &MysqlCfg{
		Username: "autotester",
		Password: "autotest123",
		Host:     "localhost",
		Port:     0,
		DbName:   "autotest",
	}
	// Load a docker-based mysql server to test against
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		suite.FailNowf("Failed to init docker client: %s", err.Error())
	}
	suite.docker = cli

	_, err = suite.docker.ImagePull(ctx, "mariadb", types.ImagePullOptions{})
	if err != nil {
		suite.FailNowf("Failed to pull mariadb docker image: %s", err.Error())
	}

	startResp, err := suite.docker.ContainerCreate(ctx, &container.Config{
		Env: []string{
			"MARIADB_RANDOM_ROOT_PASSWORD=yes",
			fmt.Sprintf("MARIADB_DATABASE=%s", suite.dbCfg.DbName),
			fmt.Sprintf("MARIADB_USER=%s", suite.dbCfg.Username),
			fmt.Sprintf("MARIADB_PASSWORD=%s", suite.dbCfg.Password),
		},
		Image: "mariadb",
	}, nil, nil, nil, "autotest-db")
	if err != nil {
		suite.FailNowf("Failed to create mariadb container: %s", err.Error())
	}
	suite.dbId = startResp.ID

	if err := suite.docker.ContainerStart(ctx, startResp.ID, types.ContainerStartOptions{}); err != nil {
		suite.FailNowf("Failed to start mariadb container: %s", err.Error())
	}

	statusCh, errCh := suite.docker.ContainerWait(ctx, startResp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			suite.FailNowf("Failed to get mariadb container running: %s", err.Error())
		}
	case <-statusCh:
	}
	out, err := suite.docker.ContainerLogs(ctx, startResp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		suite.FailNowf("Failed to get mariadb container logs: %s", err.Error())
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	resp, err := suite.docker.ContainerInspect(ctx, startResp.ID)
	if err != nil {
		suite.FailNowf("Failed to inspect mariadb container: %s", err.Error())
	}
	fmt.Printf("%+v\n", resp.NetworkSettings.Ports)
}

func (suite *MysqlStoreTestSuite) TearDownSuite() {
	if len(suite.dbId) == 0 {
		return
	}
	ctx := context.Background()
	expiration := time.Second * 1
	err := suite.docker.ContainerStop(ctx, suite.dbId, &expiration)
	if err != nil {
		fmt.Printf("Failed to halt db container: %s\n", err.Error())
	}
}
