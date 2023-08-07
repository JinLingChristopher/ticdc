// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"context"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/retry"
	"go.uber.org/zap"
)

type saramaAdminClient struct {
	brokerEndpoints []string
	config          *sarama.Config
	changefeed      model.ChangeFeedID

	mu     sync.Mutex
	client sarama.Client
	admin  sarama.ClusterAdmin
}

const (
	defaultRetryBackoff  = 20
	defaultRetryMaxTries = 3
)

func newAdminClient(
	brokerEndpoints []string,
	config *sarama.Config,
	changefeed model.ChangeFeedID,
) (ClusterAdminClient, error) {
	client, err := sarama.NewClient(brokerEndpoints, config)
	if err != nil {
		return nil, cerror.Trace(err)
	}

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &saramaAdminClient{
		client:          client,
		admin:           admin,
		brokerEndpoints: brokerEndpoints,
		config:          config,
		changefeed:      changefeed,
	}, nil
}

func (a *saramaAdminClient) reset() error {
	newClient, err := sarama.NewClient(a.brokerEndpoints, a.config)
	if err != nil {
		return cerror.Trace(err)
	}
	newAdmin, err := sarama.NewClusterAdminFromClient(newClient)
	if err != nil {
		return cerror.Trace(err)
	}

	_ = a.admin.Close()
	a.client = newClient
	a.admin = newAdmin
	log.Info("kafka admin client is reset",
		zap.String("namespace", a.changefeed.Namespace),
		zap.String("changefeed", a.changefeed.ID))
	return errors.New("retry after reset")
}

func (a *saramaAdminClient) queryClusterWithRetry(ctx context.Context, query func() error) error {
	err := retry.Do(ctx, func() error {
		a.mu.Lock()
		defer a.mu.Unlock()
		err := query()
		if err == nil {
			return nil
		}

		log.Warn("query kafka cluster meta failed, retry it",
			zap.String("namespace", a.changefeed.Namespace),
			zap.String("changefeed", a.changefeed.ID),
			zap.Error(err))

		if cerror.Is(err, syscall.EPIPE) || cerror.Is(err, net.ErrClosed) || cerror.Is(err, io.EOF) {
			return a.reset()
		}
		return err
	}, retry.WithBackoffBaseDelay(defaultRetryBackoff), retry.WithMaxTries(defaultRetryMaxTries))
	return err
}

func (a *saramaAdminClient) GetAllBrokers(_ context.Context) ([]Broker, error) {
	brokers := a.client.Brokers()
	result := make([]Broker, 0, len(brokers))
	for _, broker := range brokers {
		result = append(result, Broker{
			ID: broker.ID(),
		})
	}

	return result, nil
}

func (a *saramaAdminClient) GetBrokerConfig(
	_ context.Context,
	configName string,
) (string, error) {
	controller, err := a.client.Controller()
	if err != nil {
		return "", errors.Trace(err)
	}

	configEntries, err := a.admin.DescribeConfig(sarama.ConfigResource{
		Type:        sarama.BrokerResource,
		Name:        strconv.Itoa(int(controller.ID())),
		ConfigNames: []string{configName},
	})

	if err != nil {
		return "", errors.Trace(err)
	}

	// For compatibility with KOP, we checked all return values.
	// 1. Kafka only returns requested configs.
	// 2. Kop returns all configs.
	for _, entry := range configEntries {
		if entry.Name == configName {
			return entry.Value, nil
		}
	}

	log.Warn("Kafka config item not found",
		zap.String("namespace", a.changefeed.Namespace),
		zap.String("changefeed", a.changefeed.ID),
		zap.String("configName", configName))
	return "", cerror.ErrKafkaConfigNotFound.GenWithStack(
		"cannot find the `%s` from the broker's configuration", configName)
}

func (a *saramaAdminClient) GetTopicConfig(
	_ context.Context, topicName string, configName string,
) (string, error) {
	configEntries, err := a.admin.DescribeConfig(sarama.ConfigResource{
		Type:        sarama.TopicResource,
		Name:        topicName,
		ConfigNames: []string{configName},
	})
	if err != nil {
		return "", errors.Trace(err)
	}

	// For compatibility with KOP, we checked all return values.
	// 1. Kafka only returns requested configs.
	// 2. Kop returns all configs.
	for _, entry := range configEntries {
		if entry.Name == configName {
			log.Info("Kafka config item found",
				zap.String("namespace", a.changefeed.Namespace),
				zap.String("changefeed", a.changefeed.ID),
				zap.String("configName", configName),
				zap.String("configValue", entry.Value))
			return entry.Value, nil
		}
	}

	log.Warn("Kafka config item not found",
		zap.String("namespace", a.changefeed.Namespace),
		zap.String("changefeed", a.changefeed.ID),
		zap.String("configName", configName))
	return "", cerror.ErrKafkaConfigNotFound.GenWithStack(
		"cannot find the `%s` from the topic's configuration", configName)
}

func (a *saramaAdminClient) GetTopicsMeta(
	_ context.Context, topics []string, _ bool,
) (map[string]TopicDetail, error) {
	result := make(map[string]TopicDetail, len(topics))
	for _, topic := range topics {
		partitions, err := a.client.Partitions(topic)
		if err != nil {
			return nil, errors.Trace(err)
		}
		result[topic] = TopicDetail{
			Name:          topic,
			NumPartitions: int32(len(partitions)),
		}
	}

	return result, nil
}

func (a *saramaAdminClient) CreateTopic(
	ctx context.Context,
	detail *TopicDetail,
	validateOnly bool,
) error {
	request := &sarama.TopicDetail{
		NumPartitions:     detail.NumPartitions,
		ReplicationFactor: detail.ReplicationFactor,
	}

	err := a.admin.CreateTopic(detail.Name, request, validateOnly)
	// Ignore the already exists error because it's not harmful.
	if err != nil && !strings.Contains(err.Error(), sarama.ErrTopicAlreadyExists.Error()) {
		return err
	}
	return nil
}

func (a *saramaAdminClient) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.admin.Close(); err != nil {
		log.Warn("close admin client meet error",
			zap.String("namespace", a.changefeed.Namespace),
			zap.String("changefeed", a.changefeed.ID),
			zap.Error(err))
	}
}
