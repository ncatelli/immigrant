package mock

import (
	"errors"

	"github.com/PacketFire/immigrant/pkg/core"
)

type MockDriver struct {
	Revisions []core.Revision
}

func (this *MockDriver) Init(config map[string]string) error {
	return nil
}

func (this *MockDriver) Migrate(r core.Revision) error {
	this.Revisions = append(this.Revisions, r)
	return nil
}

func (this *MockDriver) Rollback(r core.Revision) error {
	if len(this.Revisions) == 0 {
		return errors.New("No revisions applied.")
	} else {
		this.Revisions = this.Revisions[:len(this.Revisions)-1]
	}

	return nil
}

func (this *MockDriver) State() *core.Revision {
	return &this.Revisions[len(this.Revisions)-1]
}

func Close() {
}