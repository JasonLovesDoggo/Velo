// Package gocker is a wrapper around the Docker golang client. It's called gocker to avoid naming conflicts with the Docker SDK.
package gocker

import (
	docker "github.com/docker/docker/client"
)

var client *docker.Client

func GetClient() *docker.Client {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		// todo: handle error properly once i learn what the errors can be...
		panic(err)
		return nil
	}
	return client

}
