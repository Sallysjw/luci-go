// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package tsmon

import (
	"golang.org/x/net/context"

	"github.com/luci/gae/service/module"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/tsmon"
	"github.com/luci/luci-go/common/tsmon/memstats"
	"github.com/luci/luci-go/common/tsmon/metric"
)

var (
	defaultVersion = metric.NewCallbackString(
		"appengine/default_version",
		"Name of the version currently marked as default.")
)

// collectGlobalMetrics populates service-global metrics.
//
// Called by tsmon from inside /housekeeping cron handler. Metrics reported must
// not depend on the state of the particular process that happens to report
// them.
func collectGlobalMetrics(c context.Context) {
	version, err := module.Get(c).DefaultVersion("")
	if err != nil {
		logging.Errorf(c, "Error getting default appengine version: %s", err)
		defaultVersion.Set(c, "(unknown)")
	} else {
		defaultVersion.Set(c, version)
	}
}

// collectProcessMetrics populates per-process metrics.
//
// It is called by each individual process right before flushing the metrics.
func collectProcessMetrics(c context.Context, s *tsmonSettings) {
	if s.ReportMemStats {
		memstats.Report(c)
	}
}

func init() {
	tsmon.RegisterGlobalCallback(collectGlobalMetrics, defaultVersion)
}
