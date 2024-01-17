package processor

import (
	"bytes"
	"fmt"
	"os"

	"github.com/quipham98/cdn-testing/internal/config"
	"github.com/trustwallet/assets-go-libs/file"
	"github.com/trustwallet/assets-go-libs/image"
	"github.com/trustwallet/assets-go-libs/path"
	"github.com/trustwallet/assets-go-libs/validation"
	"github.com/trustwallet/assets-go-libs/validation/info"
)

func (s *Service) ValidateJSON(f *file.AssetFile) error {
	file, err := os.Open(f.Path())
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return err
	}

	err = validation.ValidateJSON(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateRootFolder(f *file.AssetFile) error {
	dirFiles, err := file.ReadDir(f.Path())
	if err != nil {
		return err
	}

	err = validation.ValidateAllowedFiles(dirFiles, config.Default.ValidatorsSettings.RootFolder.AllowedFiles)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateChainFolder(f *file.AssetFile) error {
	file, err := os.Open(f.Path())
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	var compErr = validation.NewErrComposite()

	err = validation.ValidateLowercase(fileInfo.Name())
	if err != nil {
		compErr.Append(err)
	}

	dirFiles, err := file.ReadDir(0)
	if err != nil {
		return err
	}

	err = validation.ValidateAllowedFiles(dirFiles, config.Default.ValidatorsSettings.ChainFolder.AllowedFiles)
	if err != nil {
		compErr.Append(err)
	}

	if compErr.Len() > 0 {
		return compErr
	}

	return nil
}

func validatePngImageDimensionForCI(path string) error {
	imgWidth, imgHeight, err := image.GetPNGImageDimensions(path)
	if err != nil {
		return err
	}

	if imgWidth > config.Default.ValidatorsSettings.Image.MaxW || imgHeight > config.Default.ValidatorsSettings.Image.MaxH || imgWidth < config.Default.ValidatorsSettings.Image.MinW || imgHeight < config.Default.ValidatorsSettings.Image.MinH || imgWidth != imgHeight {
		return fmt.Errorf("%s: max - %dx%d, min - %dx%d; given %dx%d",
			"invalid file dimension", config.Default.ValidatorsSettings.Image.MaxW, config.Default.ValidatorsSettings.Image.MaxH, config.Default.ValidatorsSettings.Image.MinW, config.Default.ValidatorsSettings.Image.MinH, imgWidth, imgHeight)
	}

	return nil
}

func (s *Service) ValidateImage(f *file.AssetFile) error {
	var compErr = validation.NewErrComposite()

	err := validation.ValidateLogoFileSize(f.Path())
	if err != nil {
		compErr.Append(err)
	}

	// TODO: Replace it with validation.ValidatePngImageDimension when "assets" repo is fixed.
	// Read comments in ValidatePngImageDimensionForCI.
	err = validatePngImageDimensionForCI(f.Path())
	if err != nil {
		compErr.Append(err)
	}

	if compErr.Len() > 0 {
		return compErr
	}

	return nil
}

func (s *Service) ValidateAssetFolder(f *file.AssetFile) error {
	dirFiles, err := file.ReadDir(f.Path())
	if err != nil {
		return err
	}

	var compErr = validation.NewErrComposite()

	err = validation.ValidateAllowedFiles(dirFiles, config.Default.ValidatorsSettings.AssetFolder.AllowedFiles)
	if err != nil {
		compErr.Append(err)
	}

	err = validation.ValidateAssetAddress(f.Chain(), f.Asset())
	if err != nil {
		compErr.Append(err)
	}

	errInfo := validation.ValidateHasFiles(dirFiles, []string{"info.json"})
	errLogo := validation.ValidateHasFiles(dirFiles, []string{"logo.png"})

	if errLogo != nil || errInfo != nil {
		assetInfoPath := path.GetAssetInfoPath(f.Chain().Handle, f.Asset())

		var infoJson info.AssetModel
		if err = file.ReadJSONFile(assetInfoPath, &infoJson); err != nil {
			return err
		}

		if infoJson.GetStatus() != "spam" && infoJson.GetStatus() != "abandoned" {
			compErr.Append(fmt.Errorf("%w: logo.png for non-spam assest", validation.ErrMissingFile))
		}
	}

	if compErr.Len() > 0 {
		return compErr
	}

	return nil
}

func (s *Service) ValidateInfoFolder(f *file.AssetFile) error {
	dirFiles, err := file.ReadDir(f.Path())
	if err != nil {
		return err
	}

	err = validation.ValidateHasFiles(dirFiles, config.Default.ValidatorsSettings.ChainInfoFolder.HasFiles)
	if err != nil {
		return err
	}

	return nil
}
