// Copyright 2020 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	errInternalError  = runtime.NewError("internal server error", 13) // INTERNAL
	errMarshal        = runtime.NewError("cannot marshal type", 13)   // INTERNAL
	errNoInputAllowed = runtime.NewError("no input allowed", 3)       // INVALID_ARGUMENT
	errNoUserIdFound  = runtime.NewError("no user ID in context", 3)  // INVALID_ARGUMENT
	errUnmarshal      = runtime.NewError("cannot unmarshal type", 13) // INTERNAL
)

const (
	rpcIdRewards   = "rewards"
	rpcIdFindMatch = "find_match"
	rpcMyFunc      = "MyFunc"
	rpcMyFunc2     = "MyFunc2"
)

func rpcMyFunch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	fmt.Println("xxxxxxxxxxxxxxx   sssssssssssssssssssss")
	// fmt.Println("xxxxxxxxxxxxxxxxxxxxx")
	userID := "32b9a76c-e10b-4765-8bbb-48f32c8fd569"
	changeset := map[string]int64{
		"coins": 10, // Add 10 coins to the user's wallet.
		"gems":  -5, // Remove 5 gems from the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_result": "won",
	}
	updated, previous, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)

	fmt.Println(updated)
	fmt.Println(previous)
	fmt.Println(err)

	return "", nil
}

func rpcMyFunch2(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	fmt.Println("xxxxxxxxxxxxxxx   sssssssssssssssssssss")
	// fmt.Println("xxxxxxxxxxxxxxxxxxxxx")
	userID := "93b917a8-f7e0-41c3-aaee-332360c29b9d"
	changeset := map[string]int64{
		"coins": 10, // Add 10 coins to the user's wallet.
		"gems":  5,  // Remove 5 gems from the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_result": "won",
	}
	updated, previous, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)

	fmt.Println(updated)
	fmt.Println(previous)
	fmt.Println(err)

	return "", nil
}

func MyAccessSessionVars(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) error {
	vars, ok := ctx.Value(runtime.RUNTIME_CTX_VARS).(map[string]string)

	if !ok {
		logger.Info("User session does not contain any key-value pairs set")
		return nil
	}

	logger.Info("User session contains key-value pairs set by both the client and the before authentication hook: %v", vars)
	return nil
}

// noinspection GoUnusedExportedFunction
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	fmt.Println("0------------------------------")
	initStart := time.Now()

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}

	if err := initializer.RegisterRpc(rpcIdRewards, rpcRewards); err != nil {
		return err
	}

	if err := initializer.RegisterRpc(rpcIdFindMatch, rpcFindMatch(marshaler, unmarshaler)); err != nil {
		return err
	}
	fmt.Println("1------------------------------")
	if err := initializer.RegisterRpc(rpcMyFunc, rpcMyFunch); err != nil {

		return err
	}
	fmt.Println("2------------------------------")
	if err := initializer.RegisterRpc(rpcMyFunc2, rpcMyFunch2); err != nil {

		return err
	}
	fmt.Println("3------------------------------")
	if err := initializer.RegisterBeforeGetAccount(MyAccessSessionVars); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	fmt.Println("4------------------------------")

	if err := initializer.RegisterMatch(moduleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &MatchHandler{
			marshaler:        marshaler,
			unmarshaler:      unmarshaler,
			tfServingAddress: "http://tf:8501/v1/models/ttt:predict",
		}, nil
	}); err != nil {
		return err
	}

	if err := registerSessionEvents(db, nk, initializer); err != nil {
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
