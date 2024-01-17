package processor

import (
	"github.com/trustwallet/assets-go-libs/file"
)

type Service struct {
	fileService *file.Service
}

func NewService(fileProvider *file.Service) *Service {
	return &Service{
		fileService: fileProvider,
	}
}

func (s *Service) GetValidator(f *file.AssetFile) []Validator {
	switch f.Type() {
	case file.TypeAssetFolder:
		return []Validator{
			{Name: "Each asset folder has valid asset address and contains only allowed files", Run: s.ValidateAssetFolder},
		}
	case file.TypeChainFolder:
		return []Validator{
			{Name: "Chain folders are lowercase and contains only allowed files", Run: s.ValidateChainFolder},
		}
	case file.TypeChainInfoFolder:
		return []Validator{
			{Name: "Chain Info Folder (has files)", Run: s.ValidateInfoFolder},
		}
	case file.TypeRootFolder:
		return []Validator{
			{Name: "Root folder contains only allowed files", Run: s.ValidateRootFolder},
		}

	case file.TypeAssetLogoFile, file.TypeChainLogoFile:
		return []Validator{
			{Name: "Logos size and dimension are valid", Run: s.ValidateImage},
		}
	}

	return nil
}
