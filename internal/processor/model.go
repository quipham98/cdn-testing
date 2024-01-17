package processor

import (
	"github.com/trustwallet/assets-go-libs/file"
)

type (
	Validator struct {
		Name string
		Run  func(f *file.AssetFile) error
	}

	Fixer struct {
		Name string
		Run  func(f *file.AssetFile) error
	}

	Updater struct {
		Name string
		Run  func() error
	}
)
