package fakes

import "sync"

type PackageFinder struct {
	GetBuildPackageMetadataCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Image string
		}
		Returns struct {
			Metadata string
			Err      error
		}
		Stub func(string) (string, error)
	}
	GetBuildPackagesListCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Image string
		}
		Returns struct {
			List []string
			Err  error
		}
		Stub func(string) ([]string, error)
	}
	GetRunPackageMetadataCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Image string
		}
		Returns struct {
			Metadata string
			Err      error
		}
		Stub func(string) (string, error)
	}
	GetRunPackagesListCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Image string
		}
		Returns struct {
			List []string
			Err  error
		}
		Stub func(string) ([]string, error)
	}
}

func (f *PackageFinder) GetBuildPackageMetadata(param1 string) (string, error) {
	f.GetBuildPackageMetadataCall.Lock()
	defer f.GetBuildPackageMetadataCall.Unlock()
	f.GetBuildPackageMetadataCall.CallCount++
	f.GetBuildPackageMetadataCall.Receives.Image = param1
	if f.GetBuildPackageMetadataCall.Stub != nil {
		return f.GetBuildPackageMetadataCall.Stub(param1)
	}
	return f.GetBuildPackageMetadataCall.Returns.Metadata, f.GetBuildPackageMetadataCall.Returns.Err
}
func (f *PackageFinder) GetBuildPackagesList(param1 string) ([]string, error) {
	f.GetBuildPackagesListCall.Lock()
	defer f.GetBuildPackagesListCall.Unlock()
	f.GetBuildPackagesListCall.CallCount++
	f.GetBuildPackagesListCall.Receives.Image = param1
	if f.GetBuildPackagesListCall.Stub != nil {
		return f.GetBuildPackagesListCall.Stub(param1)
	}
	return f.GetBuildPackagesListCall.Returns.List, f.GetBuildPackagesListCall.Returns.Err
}
func (f *PackageFinder) GetRunPackageMetadata(param1 string) (string, error) {
	f.GetRunPackageMetadataCall.Lock()
	defer f.GetRunPackageMetadataCall.Unlock()
	f.GetRunPackageMetadataCall.CallCount++
	f.GetRunPackageMetadataCall.Receives.Image = param1
	if f.GetRunPackageMetadataCall.Stub != nil {
		return f.GetRunPackageMetadataCall.Stub(param1)
	}
	return f.GetRunPackageMetadataCall.Returns.Metadata, f.GetRunPackageMetadataCall.Returns.Err
}
func (f *PackageFinder) GetRunPackagesList(param1 string) ([]string, error) {
	f.GetRunPackagesListCall.Lock()
	defer f.GetRunPackagesListCall.Unlock()
	f.GetRunPackagesListCall.CallCount++
	f.GetRunPackagesListCall.Receives.Image = param1
	if f.GetRunPackagesListCall.Stub != nil {
		return f.GetRunPackagesListCall.Stub(param1)
	}
	return f.GetRunPackagesListCall.Returns.List, f.GetRunPackagesListCall.Returns.Err
}
