package main

import (
	"encoding/json"
	"fmt"
	. "github.com/advancemg/vimb-loader/pkg/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"net/http"
)

func main() {
	//os.RemoveAll("/Users/eminshakh/data/vimb")
	//var bucket = "vimb"
	pool, err := New()
	if err != nil {
		panic(err)
	}
	dc, err := pool.Minio("/Users/eminshakh/data")
	if err != nil {
		println(err.Error())
	}
	defer dc.Close()
	//err = s3.CreateBucket(bucket)
	//if err != nil {
	//	fmt.Println("CreateBucket: ", err.Error())
	//}
	var input = []byte(`{"SellingDirectionID": "21","InclProgAttr": "1","InclForecast": "1","AudRatDec": "9","StartDate": "20210309","EndDate": "20210309","LightMode": "0","CnlList": {"Cnl": "1018566"},"ProtocolVersion": "2"}`)
	var js GetProgramBreaks
	json.Unmarshal(input, &js)
	err = js.UploadToS3()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (req *InternalRequest) Minio(mountPath string) (*dockertest.Resource, error) {
	resCurrent, exist := req.Pool.ContainerByName("minio-server")
	if exist {
		err := req.Pool.Purge(resCurrent)
		if err != nil {
			return nil, err
		}
	}
	dc, err := req.Pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "minio-server",
		Repository: "minio/minio",
		Tag:        "latest",
		Cmd:        []string{"server", "/data"},
		//ExposedPorts: []string{"9000", "9001"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9000/tcp": []docker.PortBinding{{HostPort: "9000"}},
		},
		Env:       []string{"MINIO_ACCESS_KEY=admin", "MINIO_SECRET_KEY=adminadmin"},
		Mounts:    []string{fmt.Sprintf("%s:/data", mountPath)},
		Name:      "minio-server",
		NetworkID: req.Network.ID,
	})
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("localhost:%s", "9000")
	if err := req.Pool.Retry(func() error {
		url := fmt.Sprintf("http://%s/minio/health/live", endpoint)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return dc, nil
}

type InternalRequest struct {
	Network *docker.Network
	Pool    *dockertest.Pool
}

func New() (*InternalRequest, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	_, err = pool.Client.PruneNetworks(docker.PruneNetworksOptions{})
	if err != nil {
		log.Fatalf("could not prune networks to docker: %s", err)
	}
	network, err := pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: "testing-net"})
	if err != nil {
		log.Fatalf("could not create a network to zookeeper and kafka: %s", err)
	}

	return &InternalRequest{
		Network: network,
		Pool:    pool,
	}, nil
}
