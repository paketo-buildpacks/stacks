package fakes

import "sync"

type PackageFinder struct {
	GetBuildPackageMetadataCall struct {
		mutex     sync.Mutex
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
		mutex     sync.Mutex
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
		mutex     sync.Mutex
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
		mutex     sync.Mutex
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
	f.GetBuildPackageMetadataCall.mutex.Lock()
	defer f.GetBuildPackageMetadataCall.mutex.Unlock()
	f.GetBuildPackageMetadataCall.CallCount++
	f.GetBuildPackageMetadataCall.Receives.Image = param1
	if f.GetBuildPackageMetadataCall.Stub != nil {
		return f.GetBuildPackageMetadataCall.Stub(param1)
	}
	return f.GetBuildPackageMetadataCall.Returns.Metadata, f.GetBuildPackageMetadataCall.Returns.Err
}
func (f *PackageFinder) GetBuildPackagesList(param1 string) ([]string, error) {
	f.GetBuildPackagesListCall.mutex.Lock()
	defer f.GetBuildPackagesListCall.mutex.Unlock()
	f.GetBuildPackagesListCall.CallCount++
	f.GetBuildPackagesListCall.Receives.Image = param1
	if f.GetBuildPackagesListCall.Stub != nil {
		return f.GetBuildPackagesListCall.Stub(param1)
	}
	return f.GetBuildPackagesListCall.Returns.List, f.GetBuildPackagesListCall.Returns.Err
}
func (f *PackageFinder) GetRunPackageMetadata(param1 string) (string, error) {
	f.GetRunPackageMetadataCall.mutex.Lock()
	defer f.GetRunPackageMetadataCall.mutex.Unlock()
	f.GetRunPackageMetadataCall.CallCount++
	f.GetRunPackageMetadataCall.Receives.Image = param1
	if f.GetRunPackageMetadataCall.Stub != nil {
		return f.GetRunPackageMetadataCall.Stub(param1)
	}
	return f.GetRunPackageMetadataCall.Returns.Metadata, f.GetRunPackageMetadataCall.Returns.Err
}
func (f *PackageFinder) GetRunPackagesList(param1 string) ([]string, error) {
	f.GetRunPackagesListCall.mutex.Lock()
	defer f.GetRunPackagesListCall.mutex.Unlock()
	f.GetRunPackagesListCall.CallCount++
	f.GetRunPackagesListCall.Receives.Image = param1
	if f.GetRunPackagesListCall.Stub != nil {
		return f.GetRunPackagesListCall.Stub(param1)
	}
	return f.GetRunPackagesListCall.Returns.List, f.GetRunPackagesListCall.Returns.Err
}
