// Copyright 2018 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mediastore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/aws/aws-sdk-go-v2/service/mediastoredata"
)

func (u *msUploader) getDataClientForContainer(ctx context.Context, name string) (*mediastoredata.Client, error) {
	v, ok := u.containers.Load(name)
	if !ok {
		client, err := u.newDataClient(ctx, name)
		if err != nil {
			return nil, err
		}
		v = client
		u.containers.Store(name, v)
	}
	return v.(*mediastoredata.Client), nil
}

func (u *msUploader) newDataClient(ctx context.Context, containerName string) (*mediastoredata.Client, error) {
	resp, err := u.client.DescribeContainer(ctx, &mediastore.DescribeContainerInput{
		ContainerName: aws.String(containerName),
	})
	if err != nil {
		return nil, err
	}
	client := mediastoredata.NewFromConfig(u.config, func(options *mediastoredata.Options) {
		options.BaseEndpoint = resp.Container.Endpoint
	})
	return client, nil
}
