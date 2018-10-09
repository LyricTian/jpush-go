package jpush

import (
	"context"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	appKey       = "b1ccd0dd04ec36b66c75e99f"
	masterSecret = "ed431429270144d3ed53555b"
)

func TestPush(t *testing.T) {
	Convey("test client push", t, func() {
		cli := NewClient(2,
			SetAppKey(appKey),
			SetMasterSecret(masterSecret),
			SetCIDCount(2),
		)

		pushID, err := cli.GetPushID(context.Background())
		So(err, ShouldBeNil)
		So(pushID, ShouldNotBeEmpty)

		payload := &Payload{
			Platform: NewPlatform().All(),
			Audience: NewAudience().All(),
			Notification: &Notification{
				Alert: "推送通知测试",
			},
			Options: &Options{
				SendNO: 1,
			},
			CID: pushID,
		}
		err = cli.Push(context.Background(), payload, func(ctx context.Context, result *PushResult, err error) {
			Convey("async callback", t, func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result.SendNO, ShouldEqual, strconv.Itoa(payload.Options.SendNO))
			})
		})
		So(err, ShouldBeNil)
		cli.Terminate()
	})
}

func TestPushValidate(t *testing.T) {
	Convey("test client push validate", t, func() {
		cli := NewClient(2,
			SetAppKey(appKey),
			SetMasterSecret(masterSecret),
			SetCIDCount(2),
		)

		payload := &Payload{
			Platform: NewPlatform().All(),
			Audience: NewAudience().All(),
			Notification: &Notification{
				Alert: "推送通知测试2",
			},
			Options: &Options{
				SendNO: 2,
			},
		}
		err := cli.PushValidate(context.Background(), payload, func(ctx context.Context, result *PushResult, err error) {
			Convey("async callback", t, func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result.SendNO, ShouldEqual, strconv.Itoa(payload.Options.SendNO))
			})
		})
		So(err, ShouldBeNil)
		cli.Terminate()
	})
}
