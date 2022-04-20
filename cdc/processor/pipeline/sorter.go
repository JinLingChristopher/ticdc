// Copyright 2020 PingCAP, Inc.
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

package pipeline

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/failpoint"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/entry"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/redo"
	"github.com/pingcap/tiflow/cdc/sorter"
	"github.com/pingcap/tiflow/cdc/sorter/leveldb"
	"github.com/pingcap/tiflow/cdc/sorter/memory"
	"github.com/pingcap/tiflow/cdc/sorter/unified"
	"github.com/pingcap/tiflow/pkg/actor"
	"github.com/pingcap/tiflow/pkg/actor/message"
	"github.com/pingcap/tiflow/pkg/config"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/pipeline"
	pmessage "github.com/pingcap/tiflow/pkg/pipeline/message"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	flushMemoryMetricsDuration = time.Second * 5
)

type sorterNode struct {
	sorter sorter.EventSorter

	tableID   model.TableID
	tableName string // quoted schema and table, used in metircs only

	// for per-table flow control
	flowController tableFlowController

	mounter entry.Mounter

	eg     *errgroup.Group
	cancel context.CancelFunc

	// The latest resolved ts that sorter has received.
	resolvedTs model.Ts

	// The latest barrier ts that sorter has received.
	barrierTs model.Ts

	replConfig *config.ReplicaConfig

	// isTableActorMode identify if the sorter node is run is actor mode, todo: remove it after GA
	isTableActorMode bool
}

func newSorterNode(
	tableName string, tableID model.TableID, startTs model.Ts,
	flowController tableFlowController, mounter entry.Mounter, sorter sorter.EventSorter,
	replConfig *config.ReplicaConfig,
) *sorterNode {
	return &sorterNode{
		tableName:      tableName,
		tableID:        tableID,
		flowController: flowController,
		mounter:        mounter,
		sorter:         sorter,
		resolvedTs:     startTs,
		barrierTs:      startTs,
		replConfig:     replConfig,
	}
}

func (n *sorterNode) Init(ctx pipeline.NodeContext) error {
	wg := errgroup.Group{}
	return n.start(ctx, false, &wg)
}

func createSorter(ctx pipeline.NodeContext, tableName string, tableID model.TableID) (sorter.EventSorter, error) {
	sortEngine := ctx.ChangefeedVars().Info.Engine
	switch sortEngine {
	case model.SortInMemory:
		return memory.NewEntrySorter(), nil
	case model.SortUnified, model.SortInFile /* `file` becomes an alias of `unified` for backward compatibility */ :
		if sortEngine == model.SortInFile {
			log.Warn("File sorter is obsolete and replaced by unified sorter. Please revise your changefeed settings",
				zap.String("changefeed", ctx.ChangefeedVars().ID), zap.String("tableName", tableName))
		}

		if config.GetGlobalServerConfig().Debug.EnableDBSorter {
			startTs := ctx.ChangefeedVars().Info.StartTs
			ssystem := ctx.GlobalVars().SorterSystem
			dbActorID := ssystem.DBActorID(uint64(tableID))
			compactScheduler := ctx.GlobalVars().SorterSystem.CompactScheduler()
			levelSorter, err := leveldb.NewSorter(
				ctx, tableID, startTs, ssystem.DBRouter, dbActorID,
				ssystem.WriterSystem, ssystem.WriterRouter,
				ssystem.ReaderSystem, ssystem.ReaderRouter,
				compactScheduler, config.GetGlobalServerConfig().Debug.DB)
			if err != nil {
				return nil, err
			}
			return levelSorter, nil
		}
		// Sorter dir has been set and checked when server starts.
		// See https://github.com/pingcap/tiflow/blob/9dad09/cdc/server.go#L275
		sortDir := config.GetGlobalServerConfig().Sorter.SortDir
		unifiedSorter, err := unified.NewUnifiedSorter(sortDir, ctx.ChangefeedVars().ID, tableName, tableID)
		if err != nil {
			return nil, err
		}
		return unifiedSorter, nil
	default:
		return nil, cerror.ErrUnknownSortEngine.GenWithStackByArgs(sortEngine)
	}
}

func (n *sorterNode) output(ctx context.Context, tableActorID actor.ID, tableActorRouter *actor.Router[pmessage.Message]) (result pmessage.Message, ok bool, err error) {
	var msg *model.PolymorphicEvent
	// We must call `sorter.Output` before receiving resolved events.
	// Skip calling `sorter.Output` and caching output channel may fail
	// to receive any events.
	output := n.sorter.Output()
	select {
	case <-ctx.Done():
		return pmessage.Message{}, false, nil
	case msg, ok = <-output:
		if !ok {
			// sorter output channel closed.
			return pmessage.Message{}, false, nil
		}
	default:
		return pmessage.Message{}, false, nil
	}

	if msg == nil || msg.RawKV == nil {
		log.Panic("unexpected empty msg", zap.Any("msg", msg))
	}

	if msg.RawKV.OpType == model.OpTypeResolved {
		// handle OpTypeResolved
		//if msg.CRTs < lastSentResolvedTs {
		//	continue
		//}
		if n.isTableActorMode {
			msg := message.ValueMessage(pmessage.TickMessage())
			_ = tableActorRouter.Send(tableActorID, msg)
		}
		return pmessage.PolymorphicEventMessage(msg), true, nil
	}

	if err := n.mounter.DecodeEvent(ctx, msg); err != nil {
		return pmessage.Message{}, false, errors.Trace(err)
	}

	return pmessage.PolymorphicEventMessage(msg), true, nil
}

func (n *sorterNode) start(
	ctx pipeline.NodeContext, isTableActorMode bool, eg *errgroup.Group) error {
	n.isTableActorMode = isTableActorMode
	n.eg = eg
	stdCtx, cancel := context.WithCancel(ctx)
	n.cancel = cancel

	failpoint.Inject("ProcessorAddTableError", func() {
		failpoint.Return(errors.New("processor add table injected error"))
	})
	n.eg.Go(func() error {
		ctx.Throw(errors.Trace(n.sorter.Run(stdCtx)))
		return nil
	})
	return nil
}

// Receive receives the message from the previous node
func (n *sorterNode) Receive(ctx pipeline.NodeContext) error {
	_, err := n.TryHandleDataMessage(ctx, ctx.Message())
	return err
}

// handleRawEvent process the raw kv event,send it to sorter
func (n *sorterNode) handleRawEvent(ctx context.Context, event *model.PolymorphicEvent) {
	rawKV := event.RawKV
	if rawKV != nil && rawKV.OpType == model.OpTypeResolved {
		// Puller resolved ts should not fall back.
		resolvedTs := rawKV.CRTs
		oldResolvedTs := atomic.SwapUint64(&n.resolvedTs, resolvedTs)
		if oldResolvedTs > resolvedTs {
			log.Panic("resolved ts regression",
				zap.Int64("tableID", n.tableID),
				zap.Uint64("resolvedTs", resolvedTs),
				zap.Uint64("oldResolvedTs", oldResolvedTs))
		}
		atomic.StoreUint64(&n.resolvedTs, rawKV.CRTs)

		if resolvedTs > n.BarrierTs() &&
			!redo.IsConsistentEnabled(n.replConfig.Consistent.Level) {
			// Do not send resolved ts events that is larger than
			// barrier ts.
			// When DDL puller stall, resolved events that outputted by
			// sorter may pile up in memory, as they have to wait DDL.
			//
			// Disabled if redolog is on, it requires sink reports
			// resolved ts, conflicts to this change.
			// TODO: Remove redolog check once redolog decouples for global
			//       resolved ts.
			event = model.NewResolvedPolymorphicEvent(0, n.BarrierTs())
		}
	}
	n.sorter.AddEntry(ctx, event)
}

func (n *sorterNode) TryHandleDataMessage(
	ctx context.Context, msg pmessage.Message,
) (bool, error) {
	switch msg.Tp {
	case pmessage.MessageTypePolymorphicEvent:
		n.handleRawEvent(ctx, msg.PolymorphicEvent)
		return true, nil
	case pmessage.MessageTypeBarrier:
		n.updateBarrierTs(msg.BarrierTs)
		fallthrough
	default:
		ctx.(pipeline.NodeContext).SendToNextNode(msg)
		return true, nil
	}
}

func (n *sorterNode) updateBarrierTs(barrierTs model.Ts) {
	if barrierTs > n.BarrierTs() {
		atomic.StoreUint64(&n.barrierTs, barrierTs)
	}
}

func (n *sorterNode) releaseResource(_ context.Context, changefeedID string) {
	defer tableMemoryHistogram.DeleteLabelValues(changefeedID)
	// Since the flowController is implemented by `Cond`, it is not cancelable by a context
	// the flowController will be blocked in a background goroutine,
	// We need to abort the flowController manually in the nodeRunner
	n.flowController.Abort()
}

func (n *sorterNode) Destroy(ctx pipeline.NodeContext) error {
	n.cancel()
	n.releaseResource(ctx, ctx.ChangefeedVars().ID)
	return n.eg.Wait()
}

func (n *sorterNode) ResolvedTs() model.Ts {
	return atomic.LoadUint64(&n.resolvedTs)
}

// BarrierTs returns the sorter barrierTs
func (n *sorterNode) BarrierTs() model.Ts {
	return atomic.LoadUint64(&n.barrierTs)
}
