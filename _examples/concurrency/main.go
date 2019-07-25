package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	"go.uber.org/zap"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	trialchan := make(chan goptuna.FrozenTrial, 8)
	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionIgnoreObjectiveErr(true),
		goptuna.StudyOptionSetTrialNotifyChannel(trialchan),
	)
	if err != nil {
		log.Fatal("failed to create study", zap.Error(err))
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := study.Optimize(objective, 100)
			if err != nil {
				log.Println("error", err)
			}
		}()
	}

	var wgnotify sync.WaitGroup
	wgnotify.Add(1)
	go func() {
		defer wgnotify.Done()
		for t := range trialchan {
			log.Println("trial", t)
		}
	}()

	wg.Wait()
	close(trialchan)
	wgnotify.Wait()

	params, err := study.GetBestParams()
	fmt.Println("best params", params, err)
}