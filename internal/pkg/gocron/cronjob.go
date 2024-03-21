package gocron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/mailer"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

func ScheduleDeleteInvalidOtp(duration time.Duration, cr repository.CacheRepository, ref string) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	t := gocron.NewTask(func() {
		cr.Delete("otp:" + ref)
		cr.Delete("user-ref:" + ref)
	})

	s.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(duration))), t)
	s.Start()
	return nil
}

func ScheduleSendNotification(time time.Time, mailer mailer.EmailService, data model.EmailNotification) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	t := gocron.NewTask(func() {
		if err := mailer.SendNotification(data); err != nil {
			fmt.Println(err.Error())
		}
	})

	s.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), t)
	s.Start()
	return nil
}

func ScheduleDeleteOrganizerDenied(time time.Time, userRepository repository.UserRepository, user entity.User) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	t := gocron.NewTask(func() {
		if err := userRepository.Delete(user); err != nil {
			fmt.Println(err.Error())
		}
	})

	s.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), t)
	s.Start()
	return nil
}
