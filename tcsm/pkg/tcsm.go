package pkg

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cast"
)

type State map[string]interface{}

func getCurrentState(workspaceID string, c *tfe.Client, ctx context.Context) (State, error) {
	var state State

	// Get workspace
	ws, err := c.Workspaces.ReadByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	// Get current state
	sv, err := c.StateVersions.ReadCurrent(ctx, ws.ID)
	if err != nil {
		return nil, err
	}

	// Download state file into memory
	resp, err := http.Get(sv.DownloadURL)
	if err != nil {
		return nil, err
	}

	// Decode JSON body into the custom type
	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func getSpecificState(stateVersion string, c *tfe.Client, ctx context.Context) (State, error) {
	var state State

	// Get current state
	sv, err := c.StateVersions.Read(ctx, stateVersion)
	if err != nil {
		return nil, err
	}

	// Download state file into memory
	resp, err := http.Get(sv.DownloadURL)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func prepareState(state State) (*tfe.StateVersionCreateOptions, error) {
	// Generate JSON state
	jsonState, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	// Base64 encode JSON state
	base64EncodedState := base64.StdEncoding.EncodeToString(jsonState)

	// Create state version payload
	opts := tfe.StateVersionCreateOptions{
		MD5:     tfe.String(fmt.Sprintf("%x", md5.Sum(jsonState))),
		Serial:  tfe.Int64(state["serial"].(int64)),
		State:   tfe.String(base64EncodedState),
		Lineage: tfe.String(state["lineage"].(string)),
	}

	return &opts, nil
}

func uploadState(opts *tfe.StateVersionCreateOptions, ctx context.Context,
	c *tfe.Client, workspaceID string) error {
	// Lock workspace
	lockOptions := tfe.WorkspaceLockOptions{
		Reason: tfe.String("Locking workspace in order to perform rollback."),
	}

	_, err := c.Workspaces.Lock(ctx, workspaceID, lockOptions)
	if err != nil {
		return err
	}

	// Create new state version
	_, err = c.StateVersions.Create(ctx, workspaceID, *opts)
	if err != nil {
		_, _ = c.Workspaces.Unlock(ctx, workspaceID)
		return (err)
	}

	// Unlock workspace
	_, err = c.Workspaces.Unlock(ctx, workspaceID)

	return nil
}

func RollbackToSpecificVersion(stateVersion string, ctx context.Context,
	c *tfe.Client, workspaceID string) error {
	var state State

	currentState, err := getCurrentState(workspaceID, c, ctx)
	if err != nil {
		panic(err)
	}

	specificState, err := getSpecificState(stateVersion, c, ctx)
	if err != nil {
		panic(err)
	}

	state = specificState
	state["serial"] = cast.ToInt64(currentState["serial"]) + 1

	opts, err := prepareState(state)
	if err != nil {
		return err
	}

	err = uploadState(opts, ctx, c, workspaceID)
	if err != nil {
		return err
	}

	return nil
}
