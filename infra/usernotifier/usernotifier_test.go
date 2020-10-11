package usernotifier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-users-example/domain/users"
)

type notifier interface {
	users.ChangeNotifier
	Listen() chan *users.ChangeEvent
}

func runTestSuite(t *testing.T, n notifier) {
	t.Run("notification sent are correctly received", func(t *testing.T) {
		e := &users.ChangeEvent{
			Op: users.CreateOp,
		}
		l := n.Listen()
		require.NoError(t, n.Notify(e))

		select {
		case <-time.NewTimer(3 * time.Second).C:
			t.Fatal("didn't receive event, time out after 3sec")
		case receivedEvt := <-l:
			require.Equal(t, e.Op, receivedEvt.Op)
		}
	})
}
