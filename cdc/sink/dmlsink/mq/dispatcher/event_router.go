// Copyright 2022 PingCAP, Inc.
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

package dispatcher

import (
	"strings"

	"github.com/pingcap/log"
	filter "github.com/pingcap/tidb/util/table-filter"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/sink/dmlsink/mq/dispatcher/partition"
	"github.com/pingcap/tiflow/cdc/sink/dmlsink/mq/dispatcher/topic"
	"github.com/pingcap/tiflow/pkg/config"
	"github.com/pingcap/tiflow/pkg/errors"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/sink"
	"go.uber.org/zap"
)

// EventRouter is a router, it determines which topic and which partition
// an event should be dispatched to.
type EventRouter struct {
	defaultTopic string

	rules []struct {
		partitionDispatcher partition.Dispatcher
		topicDispatcher     topic.Dispatcher
		filter.Filter
	}
}

// NewEventRouter creates a new EventRouter.
func NewEventRouter(
	cfg *config.ReplicaConfig, protocol config.Protocol, defaultTopic, scheme string,
) (*EventRouter, error) {
	// If an event does not match any dispatching rules in the config file,
	// it will be dispatched by the default partition dispatcher and
	// static topic dispatcher because it matches *.* rule.
	ruleConfigs := append(cfg.Sink.DispatchRules, &config.DispatchRule{
		Matcher:       []string{"*.*"},
		PartitionRule: "default",
		TopicRule:     "",
	})

	rules := make([]struct {
		partitionDispatcher partition.Dispatcher
		topicDispatcher     topic.Dispatcher
		filter.Filter
	}, 0, len(ruleConfigs))

	for _, ruleConfig := range ruleConfigs {
		f, err := filter.Parse(ruleConfig.Matcher)
		if err != nil {
			return nil, cerror.WrapError(cerror.ErrFilterRuleInvalid, err, ruleConfig.Matcher)
		}
		if !cfg.CaseSensitive {
			f = filter.CaseInsensitive(f)
		}

		d := getPartitionDispatcher(ruleConfig.PartitionRule, scheme)
		t, err := getTopicDispatcher(ruleConfig.TopicRule, defaultTopic, protocol, scheme)
		if err != nil {
			return nil, err
		}
		rules = append(rules, struct {
			partitionDispatcher partition.Dispatcher
			topicDispatcher     topic.Dispatcher
			filter.Filter
		}{partitionDispatcher: d, topicDispatcher: t, Filter: f})
	}

	return &EventRouter{
		defaultTopic: defaultTopic,
		rules:        rules,
	}, nil
}

// GetTopicForRowChange returns the target topic for row changes.
func (s *EventRouter) GetTopicForRowChange(row *model.RowChangedEvent) string {
	topicDispatcher, _ := s.matchDispatcher(row.Table.Schema, row.Table.Table)
	return topicDispatcher.Substitute(row.Table.Schema, row.Table.Table)
}

// GetTopicForDDL returns the target topic for DDL.
func (s *EventRouter) GetTopicForDDL(ddl *model.DDLEvent) string {
	var schema, table string
	if ddl.PreTableInfo != nil {
		if ddl.PreTableInfo.TableName.Table == "" {
			return s.defaultTopic
		}
		schema = ddl.PreTableInfo.TableName.Schema
		table = ddl.PreTableInfo.TableName.Table
	} else {
		if ddl.TableInfo.TableName.Table == "" {
			return s.defaultTopic
		}
		schema = ddl.TableInfo.TableName.Schema
		table = ddl.TableInfo.TableName.Table
	}

	topicDispatcher, _ := s.matchDispatcher(schema, table)
	return topicDispatcher.Substitute(schema, table)
}

// GetPartitionForRowChange returns the target partition for row changes.
func (s *EventRouter) GetPartitionForRowChange(
	row *model.RowChangedEvent,
	partitionNum int32,
) (int32, string) {
	_, partitionDispatcher := s.matchDispatcher(
		row.Table.Schema, row.Table.Table,
	)

	return partitionDispatcher.DispatchRowChangedEvent(
		row, partitionNum,
	)
}

// GetActiveTopics returns a list of the corresponding topics
// for the tables that are actively synchronized.
func (s *EventRouter) GetActiveTopics(activeTables []model.TableName) []string {
	topics := make([]string, 0)
	topicsMap := make(map[string]bool, len(activeTables))
	for _, table := range activeTables {
		topicDispatcher, _ := s.matchDispatcher(table.Schema, table.Table)
		topicName := topicDispatcher.Substitute(table.Schema, table.Table)
		if topicName == s.defaultTopic {
			log.Debug("topic name corresponding to the table is the same as the default topic name",
				zap.String("table", table.String()),
				zap.String("defaultTopic", s.defaultTopic),
				zap.String("topicDispatcherExpression", topicDispatcher.String()),
			)
		}
		if !topicsMap[topicName] {
			topicsMap[topicName] = true
			topics = append(topics, topicName)
		}
	}

	// We also need to add the default topic.
	if !topicsMap[s.defaultTopic] {
		topics = append(topics, s.defaultTopic)
	}

	return topics
}

// GetDefaultTopic returns the default topic name.
func (s *EventRouter) GetDefaultTopic() string {
	return s.defaultTopic
}

// matchDispatcher returns the target topic dispatcher and partition dispatcher if a
// row changed event matches a specific table filter.
func (s *EventRouter) matchDispatcher(
	schema, table string,
) (topic.Dispatcher, partition.Dispatcher) {
	for _, rule := range s.rules {
		if !rule.MatchTable(schema, table) {
			continue
		}
		return rule.topicDispatcher, rule.partitionDispatcher
	}
	log.Panic("the dispatch rule must cover all tables")
	return nil, nil
}

// getPartitionDispatcher returns the partition dispatcher for a specific partition rule.
func getPartitionDispatcher(rule string, scheme string) partition.Dispatcher {
	switch strings.ToLower(rule) {
	case "default":
		return partition.NewDefaultDispatcher()
	case "ts":
		return partition.NewTsDispatcher()
	case "table":
		return partition.NewTableDispatcher()
	case "index-value", "rowid":
		log.Warn("rowid is deprecated, please use index-value instead.")
		return partition.NewIndexValueDispatcher()
	default:
	}

	if sink.IsPulsarScheme(scheme) {
		return partition.NewKeyDispatcher(rule)
	}

	log.Warn("the partition dispatch rule is not default/ts/table/index-value," +
		" use the default rule instead.")
	return partition.NewDefaultDispatcher()
}

// getTopicDispatcher returns the topic dispatcher for a specific topic rule (aka topic expression).
func getTopicDispatcher(
	rule string, defaultTopic string, protocol config.Protocol, schema string,
) (topic.Dispatcher, error) {
	if rule == "" {
		return topic.NewStaticTopicDispatcher(defaultTopic), nil
	}

	if topic.IsHardCode(rule) {
		return topic.NewStaticTopicDispatcher(rule), nil
	}

	// check if this rule is a valid topic expression
	topicExpr := topic.Expression(rule)

	var err error
	// validate the topic expression for pulsar sink
	if sink.IsPulsarScheme(schema) {
		err = topicExpr.PulsarValidate()
		if err != nil {
			return nil, errors.Trace(err)
		}
	} else {
		// validate the topic expression for kafka sink
		switch protocol {
		case config.ProtocolAvro:
			err = topicExpr.ValidateForAvro()
		default:
			err = topicExpr.Validate()
		}
	}
	if err != nil {
		return nil, err
	}

	return topic.NewDynamicTopicDispatcher(topicExpr), nil
}
