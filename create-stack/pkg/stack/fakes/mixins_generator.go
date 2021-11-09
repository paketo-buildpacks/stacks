package fakes

import "sync"

type MixinsGenerator struct {
	GetMixinsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			BuildPackages []string
			RunPackages   []string
		}
		Returns struct {
			BuildMixins []string
			RunMixins   []string
		}
		Stub func([]string, []string) ([]string, []string)
	}
}

func (f *MixinsGenerator) GetMixins(param1 []string, param2 []string) ([]string, []string) {
	f.GetMixinsCall.mutex.Lock()
	defer f.GetMixinsCall.mutex.Unlock()
	f.GetMixinsCall.CallCount++
	f.GetMixinsCall.Receives.BuildPackages = param1
	f.GetMixinsCall.Receives.RunPackages = param2
	if f.GetMixinsCall.Stub != nil {
		return f.GetMixinsCall.Stub(param1, param2)
	}
	return f.GetMixinsCall.Returns.BuildMixins, f.GetMixinsCall.Returns.RunMixins
}
