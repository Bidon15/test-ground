package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"
)

var testcases = map[string]interface{}{
	"one":  run.InitializedTestCaseFn(oneToMany),
	"many": run.InitializedTestCaseFn(manyToMany),
}

func main() {
	run.InvokeMap(testcases)
}

func oneToMany(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	client := initCtx.SyncClient
	id := initCtx.GlobalSeq

	for i := 1; i <= runenv.TestInstanceCount; i++ {
		topic := fmt.Sprintf("topic-%d", i)
		ft := sync.NewTopic(topic, "")

		if strings.Contains(topic, strconv.Itoa(int(id))) {
			seq, err := client.Publish(ctx, ft, fmt.Sprintf("%d has published to the %s topic", id, topic))
			if err != nil {
				return err
			}
			runenv.RecordMessage("Seq number - %d", seq)
		} else {
			fch := make(chan string)
			_, err := client.Subscribe(ctx, ft, fch)
			if err != nil {
				return err
			}

			for i := 1; i <= runenv.TestInstanceCount-1; i++ {
				f := <-fch
				runenv.RecordMessage("%d has received the message --> %s", id, f)
			}
		}
	}
	return nil
}

func manyToMany(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	client := initCtx.SyncClient
	id := initCtx.GlobalSeq

	for i := 1; i <= runenv.TestInstanceCount; i++ {
		topic := fmt.Sprintf("topic-%d", i)
		ft := sync.NewTopic(topic, "")

		seq, err := client.Publish(ctx, ft, fmt.Sprintf("%d has published to the %s topic", id, topic))
		if err != nil {
			return err
		}
		runenv.RecordMessage("Seq number - %d", seq)

		fch := make(chan string)
		_, err = client.Subscribe(ctx, ft, fch)
		if err != nil {
			return err
		}

		for i := 1; i <= runenv.TestInstanceCount; i++ {
			f := <-fch
			runenv.RecordMessage("%d has received the message --> %s", id, f)
		}
	}
	return nil
}
