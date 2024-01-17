package service

import (
	"errors"
	"github.com/quipham98/cdn-testing/internal/processor"
	"github.com/quipham98/cdn-testing/internal/report"
	"github.com/trustwallet/assets-go-libs/file"
	"github.com/trustwallet/assets-go-libs/validation"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	fileService      *file.Service
	processorService *processor.Service
	reportService    *report.Service
	paths            []string
}

func NewService(fs *file.Service, cs *processor.Service, rs *report.Service, paths []string) *Service {
	return &Service{
		fileService:      fs,
		processorService: cs,
		reportService:    rs,
		paths:            paths,
	}
}

func (s *Service) RunJob(job func(*file.AssetFile)) {
	for _, path := range s.paths {
		f := s.fileService.GetAssetFile(path)
		s.reportService.IncTotalFiles()
		job(f)
	}

	reportMsg := s.reportService.GetReport()
	if s.reportService.IsFailed() {
		log.Fatal(reportMsg)
	} else {
		log.Info(reportMsg)
	}
}

func (s *Service) Check(f *file.AssetFile) {
	validators := s.processorService.GetValidator(f)

	for _, validator := range validators {
		if err := validator.Run(f); err != nil {
			s.handleError(err, f, validator.Name)
		}
	}
}

func (s *Service) handleError(err error, info *file.AssetFile, valName string) {
	errors := UnwrapComposite(err)

	for _, err := range errors {
		log.WithFields(log.Fields{
			"type":       info.Type(),
			"chain":      info.Chain().Handle,
			"asset":      info.Asset(),
			"path":       info.Path(),
			"validation": valName,
		}).Error(err)

		s.reportService.IncErrors()
	}
}

func UnwrapComposite(err error) []error {
	var compErr *validation.ErrComposite
	ok := errors.As(err, &compErr)
	if !ok {
		return []error{err}
	}

	var errors []error
	for _, e := range compErr.GetErrors() {
		errors = append(errors, UnwrapComposite(e)...)
	}

	return errors
}
