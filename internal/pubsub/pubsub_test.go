package pubsub

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPubSub(t *testing.T) {
	convey.Convey("Subject: pubsub", t, func() {
		ps := NewPubSub[string](1)

		convey.Convey("Subscribe/unsubscribe for one subscription", func() {
			subCh := ps.Subscribe()
			ps.Publish("test")

			convey.So(len(subCh), convey.ShouldEqual, 1)
			convey.So(<-subCh, convey.ShouldEqual, "test")

			ps.Unsubscribe(subCh)

			ps.Publish("test")
			convey.So(len(subCh), convey.ShouldEqual, 0)
		})

		convey.Convey("Subscribe/unsubscribe for several subscriptions", func() {
			subCh := ps.Subscribe()
			subCh1 := ps.Subscribe()
			subCh2 := ps.Subscribe()
			ps.Publish("test")

			convey.So(<-subCh, convey.ShouldEqual, "test")
			convey.So(<-subCh1, convey.ShouldEqual, "test")
			convey.So(<-subCh2, convey.ShouldEqual, "test")

			ps.Unsubscribe(subCh1)

			ps.Publish("test")
			convey.So(len(subCh1), convey.ShouldEqual, 0)

			convey.So(<-subCh, convey.ShouldEqual, "test")
			convey.So(<-subCh2, convey.ShouldEqual, "test")
		})
	})
}
