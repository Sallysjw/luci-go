// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package iface

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/luci/luci-go/common/ts_mon/field"
	"github.com/luci/luci-go/common/ts_mon/target"
	"github.com/luci/luci-go/common/ts_mon/types"
	"golang.org/x/net/context"

	pb "github.com/luci/luci-go/common/ts_mon/ts_mon_proto"

	. "github.com/smartystreets/goconvey/convey"
)

type fakeStore struct {
	cells []types.Cell
}

func (s *fakeStore) Register(types.Metric) error                                     { return nil }
func (s *fakeStore) Unregister(string)                                               {}
func (s *fakeStore) Get(context.Context, string, []interface{}) (interface{}, error) { return nil, nil }
func (s *fakeStore) Set(context.Context, string, []interface{}, interface{}) error   { return nil }
func (s *fakeStore) Incr(context.Context, string, []interface{}, interface{}) error  { return nil }
func (s *fakeStore) GetAll(context.Context) []types.Cell                             { return s.cells }
func (s *fakeStore) ResetForUnittest()                                               {}

type fakeMonitor struct {
	chunkSize int
	cells     [][]types.Cell
}

func (m *fakeMonitor) ChunkSize() int {
	return m.chunkSize
}

func (m *fakeMonitor) Send(cells []types.Cell, t target.Target) error {
	m.cells = append(m.cells, cells)
	return nil
}

func TestFlush(t *testing.T) {
	ctx := context.Background()

	Target = (*target.Task)(&pb.Task{
		ServiceName: proto.String("test"),
	})

	Convey("Sends a metric", t, func() {
		s := fakeStore{
			cells: []types.Cell{
				{
					MetricName: "foo",
					Fields:     []field.Field{},
					ValueType:  types.StringType,
					FieldVals:  []interface{}{},
					ResetTime:  time.Unix(1234, 1000),
					Value:      "bar",
				},
			},
		}
		Store = &s

		m := fakeMonitor{
			chunkSize: 42,
			cells:     [][]types.Cell{},
		}
		Monitor = &m

		Flush(ctx)

		So(len(m.cells), ShouldEqual, 1)
		So(len(m.cells[0]), ShouldEqual, 1)
		So(m.cells[0][0], ShouldResemble, types.Cell{
			MetricName: "foo",
			Fields:     []field.Field{},
			ValueType:  types.StringType,
			FieldVals:  []interface{}{},
			ResetTime:  time.Unix(1234, 1000),
			Value:      "bar",
		})
	})

	Convey("Splits up ChunkSize metrics", t, func() {
		s := fakeStore{
			cells: make([]types.Cell, 43),
		}
		Store = &s

		m := fakeMonitor{
			chunkSize: 42,
			cells:     [][]types.Cell{},
		}
		Monitor = &m

		for i := 0; i < 43; i++ {
			s.cells[i] = types.Cell{
				MetricName: "foo",
				Fields:     []field.Field{},
				ValueType:  types.StringType,
				FieldVals:  []interface{}{},
				ResetTime:  time.Unix(1234, 1000),
				Value:      "bar",
			}
		}

		Flush(ctx)

		So(len(m.cells), ShouldEqual, 2)
		So(len(m.cells[0]), ShouldEqual, 42)
		So(len(m.cells[1]), ShouldEqual, 1)
	})

	Convey("Doesn't split metrics when ChunkSize is 0", t, func() {
		s := fakeStore{
			cells: make([]types.Cell, 43),
		}
		Store = &s

		m := fakeMonitor{
			chunkSize: 0,
			cells:     [][]types.Cell{},
		}
		Monitor = &m

		for i := 0; i < 43; i++ {
			s.cells[i] = types.Cell{
				MetricName: "foo",
				Fields:     []field.Field{},
				ValueType:  types.StringType,
				FieldVals:  []interface{}{},
				ResetTime:  time.Unix(1234, 1000),
				Value:      "bar",
			}
		}

		Flush(ctx)

		So(len(m.cells), ShouldEqual, 1)
		So(len(m.cells[0]), ShouldEqual, 43)
	})

	Convey("No Monitor configured", t, func() {
		Monitor = nil

		err := Flush(ctx)
		So(err, ShouldNotBeNil)
	})
}