package gocron

import (
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
)

func ScheduleDeleteInvalidOtp(duration time.Duration, cr repository.CacheRepository, ref string) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	t := gocron.NewTask(func() {
		// val, _ := cr.Get("otp:" + ref)
		// fmt.Println(string(val))
		cr.Delete("otp:" + ref)
		cr.Delete("user-ref:" + ref)
	})

	s.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(duration))), t)
	s.Start()
	return nil
}
