// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

// Package leveldb provides functions to trace the syndtr/goleveldb package (https://github.com/syndtr/goleveldb).
package leveldb // import "github.com/DataDog/dd-trace-go/contrib/syndtr/goleveldb/v2/leveldb"

import (
	"context"
	"math"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/instrumentation"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const componentName = "syndtr/goleveldb/leveldb"

var instr *instrumentation.Instrumentation

func init() {
	instr = instrumentation.Load(instrumentation.PackageSyndtrGoLevelDB)
}

// A DB wraps a leveldb.DB and traces all queries.
type DB struct {
	*leveldb.DB
	cfg *config
}

// Open calls leveldb.Open and wraps the resulting DB.
func Open(stor storage.Storage, o *opt.Options, opts ...Option) (*DB, error) {
	db, err := leveldb.Open(stor, o)
	if err != nil {
		return nil, err
	}
	return WrapDB(db, opts...), nil
}

// OpenFile calls leveldb.OpenFile and wraps the resulting DB.
func OpenFile(path string, o *opt.Options, opts ...Option) (*DB, error) {
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}
	return WrapDB(db, opts...), nil
}

// WrapDB wraps a leveldb.DB so that queries are traced.
func WrapDB(db *leveldb.DB, opts ...Option) *DB {
	cfg := newConfig(opts...)
	instr.Logger().Debug("contrib/syndtr/goleveldb/leveldb: Wrapping DB: %#v", cfg)
	return &DB{
		DB:  db,
		cfg: cfg,
	}
}

// WithContext returns a new DB with the context set to ctx.
func (db *DB) WithContext(ctx context.Context) *DB {
	newcfg := *db.cfg
	newcfg.ctx = ctx
	return &DB{
		DB:  db.DB,
		cfg: &newcfg,
	}
}

// CompactRange calls DB.CompactRange and traces the result.
func (db *DB) CompactRange(r util.Range) error {
	span := startSpan(db.cfg, "CompactRange")
	err := db.DB.CompactRange(r)
	span.Finish(tracer.WithError(err))
	return err
}

// Delete calls DB.Delete and traces the result.
func (db *DB) Delete(key []byte, wo *opt.WriteOptions) error {
	span := startSpan(db.cfg, "Delete")
	err := db.DB.Delete(key, wo)
	span.Finish(tracer.WithError(err))
	return err
}

// Get calls DB.Get and traces the result.
func (db *DB) Get(key []byte, ro *opt.ReadOptions) (value []byte, err error) {
	span := startSpan(db.cfg, "Get")
	value, err = db.DB.Get(key, ro)
	span.Finish(tracer.WithError(err))
	return value, err
}

// GetSnapshot calls DB.GetSnapshot and returns a wrapped Snapshot.
func (db *DB) GetSnapshot() (*Snapshot, error) {
	snap, err := db.DB.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return WrapSnapshot(snap, withConfig(db.cfg)), nil
}

// Has calls DB.Has and traces the result.
func (db *DB) Has(key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	span := startSpan(db.cfg, "Has")
	ret, err = db.DB.Has(key, ro)
	span.Finish(tracer.WithError(err))
	return ret, err
}

// NewIterator calls DB.NewIterator and returns a wrapped Iterator.
func (db *DB) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return WrapIterator(db.DB.NewIterator(slice, ro), withConfig(db.cfg))
}

// OpenTransaction calls DB.OpenTransaction and returns a wrapped Transaction.
func (db *DB) OpenTransaction() (*Transaction, error) {
	tr, err := db.DB.OpenTransaction()
	if err != nil {
		return nil, err
	}
	return WrapTransaction(tr, withConfig(db.cfg)), nil
}

// Put calls DB.Put and traces the result.
func (db *DB) Put(key, value []byte, wo *opt.WriteOptions) error {
	span := startSpan(db.cfg, "Put")
	err := db.DB.Put(key, value, wo)
	span.Finish(tracer.WithError(err))
	return err
}

// Write calls DB.Write and traces the result.
func (db *DB) Write(batch *leveldb.Batch, wo *opt.WriteOptions) error {
	span := startSpan(db.cfg, "Write")
	err := db.DB.Write(batch, wo)
	span.Finish(tracer.WithError(err))
	return err
}

// A Snapshot wraps a leveldb.Snapshot and traces all queries.
type Snapshot struct {
	*leveldb.Snapshot
	cfg *config
}

// WrapSnapshot wraps a leveldb.Snapshot so that queries are traced.
func WrapSnapshot(snap *leveldb.Snapshot, opts ...Option) *Snapshot {
	return &Snapshot{
		Snapshot: snap,
		cfg:      newConfig(opts...),
	}
}

// WithContext returns a new Snapshot with the context set to ctx.
func (snap *Snapshot) WithContext(ctx context.Context) *Snapshot {
	newcfg := *snap.cfg
	newcfg.ctx = ctx
	return &Snapshot{
		Snapshot: snap.Snapshot,
		cfg:      &newcfg,
	}
}

// Get calls Snapshot.Get and traces the result.
func (snap *Snapshot) Get(key []byte, ro *opt.ReadOptions) (value []byte, err error) {
	span := startSpan(snap.cfg, "Get")
	value, err = snap.Snapshot.Get(key, ro)
	span.Finish(tracer.WithError(err))
	return value, err
}

// Has calls Snapshot.Has and traces the result.
func (snap *Snapshot) Has(key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	span := startSpan(snap.cfg, "Has")
	ret, err = snap.Snapshot.Has(key, ro)
	span.Finish(tracer.WithError(err))
	return ret, err
}

// NewIterator calls Snapshot.NewIterator and returns a wrapped Iterator.
func (snap *Snapshot) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return WrapIterator(snap.Snapshot.NewIterator(slice, ro), withConfig(snap.cfg))
}

// A Transaction wraps a leveldb.Transaction and traces all queries.
type Transaction struct {
	*leveldb.Transaction
	cfg *config
}

// WrapTransaction wraps a leveldb.Transaction so that queries are traced.
func WrapTransaction(tr *leveldb.Transaction, opts ...Option) *Transaction {
	return &Transaction{
		Transaction: tr,
		cfg:         newConfig(opts...),
	}
}

// WithContext returns a new Transaction with the context set to ctx.
func (tr *Transaction) WithContext(ctx context.Context) *Transaction {
	newcfg := *tr.cfg
	newcfg.ctx = ctx
	return &Transaction{
		Transaction: tr.Transaction,
		cfg:         &newcfg,
	}
}

// Commit calls Transaction.Commit and traces the result.
func (tr *Transaction) Commit() error {
	span := startSpan(tr.cfg, "Commit")
	err := tr.Transaction.Commit()
	span.Finish(tracer.WithError(err))
	return err
}

// Get calls Transaction.Get and traces the result.
func (tr *Transaction) Get(key []byte, ro *opt.ReadOptions) ([]byte, error) {
	span := startSpan(tr.cfg, "Get")
	value, err := tr.Transaction.Get(key, ro)
	span.Finish(tracer.WithError(err))
	return value, err
}

// Has calls Transaction.Has and traces the result.
func (tr *Transaction) Has(key []byte, ro *opt.ReadOptions) (bool, error) {
	span := startSpan(tr.cfg, "Has")
	ret, err := tr.Transaction.Has(key, ro)
	span.Finish(tracer.WithError(err))
	return ret, err
}

// NewIterator calls Transaction.NewIterator and returns a wrapped Iterator.
func (tr *Transaction) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return WrapIterator(tr.Transaction.NewIterator(slice, ro), withConfig(tr.cfg))
}

// An Iterator wraps a leveldb.Iterator and traces until Release is called.
type Iterator struct {
	iterator.Iterator
	span *tracer.Span
}

// WrapIterator wraps a leveldb.Iterator so that queries are traced.
func WrapIterator(it iterator.Iterator, opts ...Option) *Iterator {
	return &Iterator{
		Iterator: it,
		span:     startSpan(newConfig(opts...), "Iterator"),
	}
}

// Release calls Iterator.Release and traces the result.
func (it *Iterator) Release() {
	err := it.Error()
	it.Iterator.Release()
	it.span.Finish(tracer.WithError(err))
}

func startSpan(cfg *config, name string) *tracer.Span {
	opts := []tracer.StartSpanOption{
		tracer.SpanType(ext.SpanTypeLevelDB),
		tracer.ServiceName(cfg.serviceName),
		tracer.ResourceName(name),
		tracer.Tag(ext.Component, componentName),
		tracer.Tag(ext.SpanKind, ext.SpanKindClient),
		tracer.Tag(ext.DBSystem, ext.DBSystemLevelDB),
	}
	if !math.IsNaN(cfg.analyticsRate) {
		opts = append(opts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
	}
	span, _ := tracer.StartSpanFromContext(cfg.ctx, cfg.spanName, opts...)
	return span
}
