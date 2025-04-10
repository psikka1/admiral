package clusters

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"

	v1 "github.com/istio-ecosystem/admiral/admiral/pkg/apis/admiral/v1alpha1"
	"github.com/istio-ecosystem/admiral/admiral/pkg/controller/admiral"
	"github.com/istio-ecosystem/admiral/admiral/pkg/controller/common"
)

type ClientConnectionConfigHandler struct {
	RemoteRegistry *RemoteRegistry
	ClusterID      string
}

type ClientConnectionConfigCache interface {
	GetFromIdentity(identity string, environment string) (*v1.ClientConnectionConfig, error)
	Put(clientConnectionSettings *v1.ClientConnectionConfig) error
	Delete(identity string, environment string) error
}

type clientConnectionSettingsCache struct {
	identityCache map[string]*v1.ClientConnectionConfig
	mutex         *sync.RWMutex
}

func NewClientConnectionConfigCache() ClientConnectionConfigCache {
	return &clientConnectionSettingsCache{
		identityCache: make(map[string]*v1.ClientConnectionConfig),
		mutex:         &sync.RWMutex{},
	}
}

func (c *clientConnectionSettingsCache) GetFromIdentity(identity string,
	environment string) (*v1.ClientConnectionConfig, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.identityCache[common.ConstructKeyWithEnvAndIdentity(environment, identity)], nil
}

func (c *clientConnectionSettingsCache) Put(clientConnectionSettings *v1.ClientConnectionConfig) error {
	if clientConnectionSettings.Name == "" {
		return errors.New(
			"skipped adding to clientConnectionSettingsCache, missing name in clientConnectionSettings")
	}
	defer c.mutex.Unlock()
	c.mutex.Lock()
	var clientConnectionSettingsIdentity = common.GetClientConnectionConfigIdentity(clientConnectionSettings)
	var clientConnectionSettingsEnv = common.GetClientConnectionConfigEnv(clientConnectionSettings)

	key := common.ConstructKeyWithEnvAndIdentity(clientConnectionSettingsEnv, clientConnectionSettingsIdentity)
	c.identityCache[key] = clientConnectionSettings
	return nil
}

func (c *clientConnectionSettingsCache) Delete(identity string, environment string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	key := common.ConstructKeyWithEnvAndIdentity(environment, identity)
	if _, ok := c.identityCache[key]; ok {
		delete(c.identityCache, key)
		return nil
	}
	return fmt.Errorf("clientConnectionSettings with key %s not found in clientConnectionSettingsCache", key)
}

func (c *ClientConnectionConfigHandler) Added(ctx context.Context,
	clientConnectionSettings *v1.ClientConnectionConfig) error {
	err := HandleEventForClientConnectionConfig(
		ctx, admiral.Add, clientConnectionSettings, c.RemoteRegistry, c.ClusterID, modifyServiceEntryForNewServiceOrPod)
	if err != nil {
		return fmt.Errorf(
			LogErrFormat, common.Add, common.ClientConnectionConfig, clientConnectionSettings.Name, c.ClusterID, err.Error())
	}
	return nil
}

func (c *ClientConnectionConfigHandler) Updated(
	ctx context.Context, clientConnectionSettings *v1.ClientConnectionConfig) error {
	err := HandleEventForClientConnectionConfig(
		ctx, admiral.Update, clientConnectionSettings, c.RemoteRegistry, c.ClusterID, modifyServiceEntryForNewServiceOrPod)
	if err != nil {
		return fmt.Errorf(
			LogErrFormat, common.Update, common.ClientConnectionConfig, clientConnectionSettings.Name, c.ClusterID, err.Error())
	}
	return nil
}

func (c *ClientConnectionConfigHandler) Deleted(
	ctx context.Context, clientConnectionSettings *v1.ClientConnectionConfig) error {
	err := HandleEventForClientConnectionConfig(
		ctx, admiral.Update, clientConnectionSettings, c.RemoteRegistry, c.ClusterID, modifyServiceEntryForNewServiceOrPod)
	if err != nil {
		return fmt.Errorf(
			LogErrFormat, common.Delete, common.ClientConnectionConfig, clientConnectionSettings.Name, c.ClusterID, err.Error())
	}
	return nil
}

func HandleEventForClientConnectionConfig(
	ctx context.Context, event admiral.EventType, clientConnectionSettings *v1.ClientConnectionConfig,
	registry *RemoteRegistry, clusterName string, modifySE ModifySEFunc) error {

	identity := common.GetClientConnectionConfigIdentity(clientConnectionSettings)
	if len(identity) <= 0 {
		return fmt.Errorf(
			LogFormat, "Event", common.ClientConnectionConfig, clientConnectionSettings.Name, clusterName,
			"skipped as label "+common.GetAdmiralCRDIdentityLabel()+" was not found, namespace="+clientConnectionSettings.Namespace)
	}

	env := common.GetClientConnectionConfigEnv(clientConnectionSettings)
	if len(env) <= 0 {
		return fmt.Errorf(
			LogFormat, "Event", common.ClientConnectionConfig, clientConnectionSettings.Name, clusterName,
			"skipped as env "+env+" was not found, namespace="+clientConnectionSettings.Namespace)
	}

	ctx = context.WithValue(ctx, common.ClusterName, clusterName)
	ctx = context.WithValue(ctx, common.EventResourceType, common.ClientConnectionConfig)

	_ = callRegistryForClientConnectionConfig(ctx, event, registry, clusterName, clientConnectionSettings)

	_, err := modifySE(ctx, admiral.Update, env, identity, registry)

	return err
}

func callRegistryForClientConnectionConfig(ctx context.Context, event admiral.EventType, registry *RemoteRegistry, clusterName string, clientConnectionSettings *v1.ClientConnectionConfig) error {
	var err error
	if common.IsAdmiralStateSyncerMode() && common.IsStateSyncerCluster(clusterName) && registry.RegistryClient != nil {
		switch event {
		case admiral.Add:
			err = registry.RegistryClient.PutCustomData(clusterName, clientConnectionSettings.Namespace, clientConnectionSettings.Name, common.ClientConnectionConfig, ctx.Value("txId").(string), clientConnectionSettings)
		case admiral.Update:
			err = registry.RegistryClient.PutCustomData(clusterName, clientConnectionSettings.Namespace, clientConnectionSettings.Name, common.ClientConnectionConfig, ctx.Value("txId").(string), clientConnectionSettings)
		case admiral.Delete:
			err = registry.RegistryClient.DeleteCustomData(clusterName, clientConnectionSettings.Namespace, clientConnectionSettings.Name, common.ClientConnectionConfig, ctx.Value("txId").(string))
		}
		if err != nil {
			err = fmt.Errorf(LogFormat, event, common.ClientConnectionConfig, clientConnectionSettings.Name, clusterName, "failed to "+string(event)+" "+common.ClientConnectionConfig+" with err: "+err.Error())
			log.Error(err)
		}
	}
	return err
}
