package main

import (
	"context"
	"time"

	"github.com/hatchet-dev/hatchet/pkg/client"
	"github.com/hatchet-dev/hatchet/pkg/cmdutils"
	"github.com/hatchet-dev/hatchet/pkg/worker"
	"github.com/joho/godotenv"
)

type userCreateEvent struct {
	Username string            `json:"username"`
	UserId   string            `json:"user_id"`
	Data     map[string]string `json:"data"`
}

type stepOneOutput struct {
	Message string `json:"message"`
}

func StepOne(ctx context.Context, input *userCreateEvent) (result *stepOneOutput, err error) {
	return &stepOneOutput{
		Message: "Username is: " + input.Username,
	}, nil
}

func StepTwo(ctx context.Context, input *stepOneOutput) (result *stepOneOutput, err error) {
	return &stepOneOutput{
		Message: "Above message is: " + input.Message,
	}, nil
}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	client, err := client.New()

	if err != nil {
		panic(err)
	}

	// Create a worker. This automatically reads in a TemporalClient from .env and workflow files from the .hatchet
	// directory, but this can be customized with the `worker.WithTemporalClient` and `worker.WithWorkflowFiles` options.
	w, err := worker.NewWorker(
		worker.WithClient(
			client,
		),
	)

	if err != nil {
		panic(err)
	}

	err = w.On(worker.Event("user:create"), &worker.WorkflowJob{
		Name:        "test-job",
		Description: "This is a test job.",
		Steps: []worker.WorkflowStep{
			{
				Function: StepOne,
			},
			{
				Function: StepTwo,
			},
		},
	})

	if err != nil {
		panic(err)
	}

	// err = worker.RegisterAction("echo:echo", func(ctx context.Context, input *actionInput) (result any, err error) {
	// 	return map[string]interface{}{
	// 		"message": input.Message,
	// 	}, nil
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// err = worker.RegisterAction("echo:object", func(ctx context.Context, input *actionInput) (result any, err error) {
	// 	return nil, nil
	// })

	// if err != nil {
	// 	panic(err)
	// }

	interruptCtx, cancel := cmdutils.InterruptContextFromChan(cmdutils.InterruptChan())
	defer cancel()

	go func() {
		err = w.Start(interruptCtx)

		if err != nil {
			panic(err)
		}

		cancel()
	}()

	testEvent := userCreateEvent{
		Username: "echo-test",
		UserId:   "1234",
		Data: map[string]string{
			"test": "test",
		},
	}

	// push an event
	err = client.Event().Push(
		context.Background(),
		"user:create",
		testEvent,
	)

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-interruptCtx.Done():
			return
		default:
			time.Sleep(time.Second)
		}
	}
}