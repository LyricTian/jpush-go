package jpush

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testJob struct {
	payload  int
	callback func(int)
}

func (t *testJob) Call() {
	t.payload++
	t.callback(t.payload)
}

func TestQueue(t *testing.T) {
	Convey("test queue", t, func() {
		q := NewQueue(2)
		q.Start()

		var data int
		q.Push(&testJob{
			payload: 0,
			callback: func(result int) {
				data += result
			},
		})
		q.Push(&testJob{
			payload: 0,
			callback: func(result int) {
				data += result
			},
		})
		q.Push(&testJob{
			payload: 0,
			callback: func(result int) {
				data += result
			},
		})
		q.Stop()
		So(data, ShouldEqual, 3)
	})
}
