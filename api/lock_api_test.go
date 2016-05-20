package api_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/github/git-lfs/api"
	"github.com/github/git-lfs/api/schema"
)

var LockService api.LockService

func TestSuccessfullyObtainingALock(t *testing.T) {
	got, body := LockService.Lock(new(api.LockRequest))

	AssertSchema(t, &api.RequestSchema{
		Method: http.MethodPost,
		Path:   "/locks",
		Body:   new(api.LockRequest),
		Into:   body,
	}, got)
}

func TestLockRequest(t *testing.T) {
	schema.Validate(t, schema.LockRequestSchema, &api.LockRequest{
		Path:               "/path/to/lock",
		LatestRemoteCommit: "deadbeef",
		Committer: api.Committer{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	})
}

func TestLockResponseWithLockedLock(t *testing.T) {
	schema.Validate(t, schema.LockResponseSchema, &api.LockResponse{
		Lock: &api.Lock{
			Id:   "some-lock-id",
			Path: "/lock/path",
			Committer: api.Committer{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			LockedAt: time.Now(),
		},
	})
}

func TestLockResponseWithUnlockedLock(t *testing.T) {
	schema.Validate(t, schema.LockResponseSchema, &api.LockResponse{
		Lock: &api.Lock{
			Id:   "some-lock-id",
			Path: "/lock/path",
			Committer: api.Committer{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			LockedAt:   time.Now(),
			UnlockedAt: time.Now(),
		},
	})
}

func TestLockResponseWithError(t *testing.T) {
	schema.Validate(t, schema.LockResponseSchema, &api.LockResponse{
		Err: "some error",
	})
}

func TestLockResponseWithCommitNeeded(t *testing.T) {
	schema.Validate(t, schema.LockResponseSchema, &api.LockResponse{
		CommitNeeded: "deadbeef",
	})
}

func TestLockResponseInvalidWithCommitAndError(t *testing.T) {
	schema.Refute(t, schema.LockResponseSchema, &api.LockResponse{
		Err:          "some error",
		CommitNeeded: "deadbeef",
	})
}

func TestUnlockResponseWithLock(t *testing.T) {
	schema.Validate(t, schema.UnlockResponseSchema, &api.UnlockResponse{
		Lock: &api.Lock{
			Id: "some-lock-id",
		},
	})
}

func TestUnlockResponseWithError(t *testing.T) {
	schema.Validate(t, schema.UnlockResponseSchema, &api.UnlockResponse{
		Err: "some-error",
	})
}

func TestUnlockResponseDoesNotAllowLockAndError(t *testing.T) {
	schema.Refute(t, schema.UnlockResponseSchema, &api.UnlockResponse{
		Lock: &api.Lock{
			Id: "some-lock-id",
		},
		Err: "some-error",
	})
}
