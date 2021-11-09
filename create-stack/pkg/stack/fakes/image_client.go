package fakes

import (
	"sync"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type ImageClient struct {
	BuildCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Tag            string
			DockerfilePath string
			WithBuildKit   bool
			Secrets        map[string]string
			BuildArgs      []string
		}
		Returns struct {
			Error error
		}
		Stub func(string, string, bool, map[string]string, ...string) error
	}
	PullCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Tag      string
			Keychain authn.Keychain
		}
		Returns struct {
			Image v1.Image
			Error error
		}
		Stub func(string, authn.Keychain) (v1.Image, error)
	}
	PushCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Tag string
		}
		Returns struct {
			String string
			Error  error
		}
		Stub func(string) (string, error)
	}
	SetLabelCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Tag   string
			Key   string
			Value string
		}
		Returns struct {
			Error error
		}
		Stub func(string, string, string) error
	}
}

func (f *ImageClient) Build(param1 string, param2 string, param3 bool, param4 map[string]string, param5 ...string) error {
	f.BuildCall.mutex.Lock()
	defer f.BuildCall.mutex.Unlock()
	f.BuildCall.CallCount++
	f.BuildCall.Receives.Tag = param1
	f.BuildCall.Receives.DockerfilePath = param2
	f.BuildCall.Receives.WithBuildKit = param3
	f.BuildCall.Receives.Secrets = param4
	f.BuildCall.Receives.BuildArgs = param5
	if f.BuildCall.Stub != nil {
		return f.BuildCall.Stub(param1, param2, param3, param4, param5...)
	}
	return f.BuildCall.Returns.Error
}
func (f *ImageClient) Pull(param1 string, param2 authn.Keychain) (v1.Image, error) {
	f.PullCall.mutex.Lock()
	defer f.PullCall.mutex.Unlock()
	f.PullCall.CallCount++
	f.PullCall.Receives.Tag = param1
	f.PullCall.Receives.Keychain = param2
	if f.PullCall.Stub != nil {
		return f.PullCall.Stub(param1, param2)
	}
	return f.PullCall.Returns.Image, f.PullCall.Returns.Error
}
func (f *ImageClient) Push(param1 string) (string, error) {
	f.PushCall.mutex.Lock()
	defer f.PushCall.mutex.Unlock()
	f.PushCall.CallCount++
	f.PushCall.Receives.Tag = param1
	if f.PushCall.Stub != nil {
		return f.PushCall.Stub(param1)
	}
	return f.PushCall.Returns.String, f.PushCall.Returns.Error
}
func (f *ImageClient) SetLabel(param1 string, param2 string, param3 string) error {
	f.SetLabelCall.mutex.Lock()
	defer f.SetLabelCall.mutex.Unlock()
	f.SetLabelCall.CallCount++
	f.SetLabelCall.Receives.Tag = param1
	f.SetLabelCall.Receives.Key = param2
	f.SetLabelCall.Receives.Value = param3
	if f.SetLabelCall.Stub != nil {
		return f.SetLabelCall.Stub(param1, param2, param3)
	}
	return f.SetLabelCall.Returns.Error
}
