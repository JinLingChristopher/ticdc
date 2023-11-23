// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package testkit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pingcap/log"
	ticonfig "github.com/pingcap/tidb/config"
	tiddl "github.com/pingcap/tidb/ddl"
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/meta"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/session"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/store/mockstore"
	"github.com/pingcap/tidb/testkit"
	"github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/filter"
	"github.com/stretchr/testify/require"
	"github.com/tikv/client-go/v2/oracle"
	"go.uber.org/zap"
)

type TestKit struct {
	*testkit.TestKit
	t       *testing.T
	storage kv.Storage
	domain  *domain.Domain
}

// New return a new testkit
func New(t *testing.T) *TestKit {
	store, err := mockstore.NewMockStore()
	require.NoError(t, err)
	ticonfig.UpdateGlobal(func(conf *ticonfig.Config) {
		conf.AlterPrimaryKey = true
	})
	session.SetSchemaLease(0)
	session.DisableStats4Test()
	domain, err := session.BootstrapSession(store)
	require.NoError(t, err)
	domain.SetStatsUpdating(true)
	tk := testkit.NewTestKit(t, store)
	return &TestKit{
		t:       t,
		TestKit: tk,
		storage: store,
		domain:  domain,
	}
}

// DDL2Job executes the DDL stmt and returns the DDL job
func (tk *TestKit) DDL2Job(ddl string) *model.Job {
	tk.MustExec(ddl)
	jobs, err := tiddl.GetLastNHistoryDDLJobs(tk.GetCurrentMeta(), 1)
	require.Nil(tk.t, err)
	require.Len(tk.t, jobs, 1)
	// Set State from Synced to Done.
	// Because jobs are put to history queue after TiDB alter its state from
	// Done to Synced.
	jobs[0].State = model.JobStateDone
	res := jobs[0]
	if res.Type != model.ActionRenameTables {
		return res
	}

	// the RawArgs field in job fetched from tidb snapshot meta is incorrent,
	// so we manually construct `job.RawArgs` to do the workaround.
	// we assume the old schema name is same as the new schema name here.
	// for example, "ALTER TABLE RENAME test.t1 TO test.t1, test.t2 to test.t22", schema name is "test"
	schema := strings.Split(strings.Split(strings.Split(res.Query, ",")[1], " ")[1], ".")[0]
	tableNum := len(res.BinlogInfo.MultipleTableInfos)
	oldSchemaIDs := make([]int64, tableNum)
	for i := 0; i < tableNum; i++ {
		oldSchemaIDs[i] = res.SchemaID
	}
	oldTableIDs := make([]int64, tableNum)
	for i := 0; i < tableNum; i++ {
		oldTableIDs[i] = res.BinlogInfo.MultipleTableInfos[i].ID
	}
	newTableNames := make([]model.CIStr, tableNum)
	for i := 0; i < tableNum; i++ {
		newTableNames[i] = res.BinlogInfo.MultipleTableInfos[i].Name
	}
	oldSchemaNames := make([]model.CIStr, tableNum)
	for i := 0; i < tableNum; i++ {
		oldSchemaNames[i] = model.NewCIStr(schema)
	}
	newSchemaIDs := oldSchemaIDs

	args := []interface{}{
		oldSchemaIDs, newSchemaIDs,
		newTableNames, oldTableIDs, oldSchemaNames,
	}
	rawArgs, err := json.Marshal(args)
	require.NoError(tk.t, err)
	res.RawArgs = rawArgs
	return res
}

// DDL2Jobs executes the DDL statement and return the corresponding DDL jobs.
// It is mainly used for "DROP TABLE" and "DROP VIEW" statement because
// multiple jobs will be generated after executing these two types of
// DDL statements.
func (tk *TestKit) DDL2Jobs(ddl string, jobCnt int) []*model.Job {
	tk.MustExec(ddl)
	jobs, err := tiddl.GetLastNHistoryDDLJobs(tk.GetCurrentMeta(), jobCnt)
	require.Nil(tk.t, err)
	require.Len(tk.t, jobs, jobCnt)
	// Set State from Synced to Done.
	// Because jobs are put to history queue after TiDB alter its state from
	// Done to Synced.
	for i := range jobs {
		jobs[i].State = model.JobStateDone
	}
	return jobs
}

// Storage returns the tikv storage
func (tk *TestKit) Storage() kv.Storage {
	return tk.storage
}

// GetCurrentMeta return the current meta snapshot
func (tk *TestKit) GetCurrentMeta() *meta.Meta {
	ver, err := tk.storage.CurrentVersion(oracle.GlobalTxnScope)
	require.Nil(tk.t, err)
	return meta.NewSnapshotMeta(tk.storage.GetSnapshot(ver))
}

// Close closes the helper
func (tk *TestKit) Close() {
	tk.domain.Close()
	tk.storage.Close() //nolint:errcheck
}

func (tk *TestKit) GetAllHistoryDDLJob(f filter.Filter) ([]*model.Job, error) {
	s, err := session.CreateSession(tk.storage)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if s != nil {
		defer s.Close()
	}

	store := domain.GetDomain(s.(sessionctx.Context)).Store()
	txn, err := store.Begin()
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer txn.Rollback() //nolint:errcheck
	txnMeta := meta.NewMeta(txn)

	jobs, err := tiddl.GetAllHistoryDDLJobs(txnMeta)
	res := make([]*model.Job, 0)
	if err != nil {
		return nil, errors.Trace(err)
	}
	for i, job := range jobs {
		ignoreSchema := f.ShouldIgnoreSchema(job.SchemaName)
		ignoreTable := f.ShouldIgnoreTable(job.SchemaName, job.TableName)
		if ignoreSchema || ignoreTable {
			log.Info("Ignore ddl job", zap.Stringer("job", job))
			continue
		}
		// Set State from Synced to Done.
		// Because jobs are put to history queue after TiDB alter its state from
		// Done to Synced.
		jobs[i].State = model.JobStateDone
		res = append(res, job)
	}
	return jobs, nil
}
