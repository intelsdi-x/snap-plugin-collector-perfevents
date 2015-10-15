//
// +build unit

package perfevents

import (
	"testing"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPSUtilPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create PerfEvents Collector", t, func() {
		pfCol := NewPerfeventsCollector()
		Convey("So pfCol should not be nil", func() {
			So(pfCol, ShouldNotBeNil)
		})
		Convey("So pfCol should be of Psutil type", func() {
			So(pfCol, ShouldHaveSameTypeAs, &Perfevents{})
		})
		configPolicy, err := pfCol.GetConfigPolicy()
		Convey("pfCol.GetConfigPolicy() should return a config policy", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So getting config policy should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
		})
	})
}
