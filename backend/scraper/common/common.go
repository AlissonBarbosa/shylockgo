package common

import (
  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
)

func GetProvider() (*gophercloud.ProviderClient, error)  {
  opts, err := openstack.AuthOptionsFromEnv()
  if err != nil {
    return nil, err
  }
  provider, err := openstack.AuthenticatedClient(opts)
  if err != nil {
    return nil, err
  }
  return provider, nil
}
