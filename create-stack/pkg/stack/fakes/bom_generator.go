package fakes

import "sync"

type BOMGenerator struct {
	AttachCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			CnbImageTag string
			Files       []string
		}
		Returns struct {
			Err error
		}
		Stub func(string, []string) error
	}
	GenerateCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			ImageTag string
		}
		Returns struct {
			OutputPaths []string
			Err         error
		}
		Stub func(string) ([]string, error)
	}
}

func (f *BOMGenerator) Attach(param1 string, param2 []string) error {
	f.AttachCall.Lock()
	defer f.AttachCall.Unlock()
	f.AttachCall.CallCount++
	f.AttachCall.Receives.CnbImageTag = param1
	f.AttachCall.Receives.Files = param2
	if f.AttachCall.Stub != nil {
		return f.AttachCall.Stub(param1, param2)
	}
	return f.AttachCall.Returns.Err
}
func (f *BOMGenerator) Generate(param1 string) ([]string, error) {
	f.GenerateCall.Lock()
	defer f.GenerateCall.Unlock()
	f.GenerateCall.CallCount++
	f.GenerateCall.Receives.ImageTag = param1
	if f.GenerateCall.Stub != nil {
		return f.GenerateCall.Stub(param1)
	}
	return f.GenerateCall.Returns.OutputPaths, f.GenerateCall.Returns.Err
}
