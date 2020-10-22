package meta

import (
	"fmt"
)

type InstanceInfo struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	Comment string `json:"comment,omitempty"`
}

func (i *InstanceInfo) GetMName() string {
	return fmt.Sprintf("%s%s", i.Name, InstanceSuffix)
}

func (i *InstanceInfo) GetPath() string {
	return fmt.Sprintf("%s%s%s", MetaPath, PathSeparator, i.GetMName())
}

func (i *InstanceInfo) Store() error {
	return StoreMetadata(i)
}

func (i *InstanceInfo) Load() error {
	return LoadMetadata(i)
}
