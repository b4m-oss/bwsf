package infra

import (
	"errors"

	"bwsf/src/core"
)

// ErrAPINotImplemented is returned by the API backend stub until later Issue #53 steps land.
var ErrAPINotImplemented = errors.New(
	"API backend is not implemented yet (Issue #53). " +
		"The Step 0 feasibility spike lives under app/spike/. " +
		"Use `bwsf backend --set bw` to switch back to the Bitwarden CLI backend",
)

// ApiBwClient is a stub BwClient for backend "api".
// All methods return ErrAPINotImplemented until the real API adapter is wired.
type ApiBwClient struct{}

// NewApiBwClient creates an ApiBwClient stub.
func NewApiBwClient() *ApiBwClient {
	return &ApiBwClient{}
}

func (c *ApiBwClient) GetDotenvsFolderID() (string, error) {
	return "", ErrAPINotImplemented
}

func (c *ApiBwClient) DotenvsFolderExists() (bool, error) {
	return false, ErrAPINotImplemented
}

func (c *ApiBwClient) CreateDotenvsFolder() error {
	return ErrAPINotImplemented
}

func (c *ApiBwClient) ListItemsInFolder(folderID string) ([]core.Item, error) {
	return nil, ErrAPINotImplemented
}

func (c *ApiBwClient) GetItemByName(folderID, name string) (*core.FullItem, error) {
	return nil, ErrAPINotImplemented
}

func (c *ApiBwClient) GetItemByID(id string) (*core.FullItem, error) {
	return nil, ErrAPINotImplemented
}

func (c *ApiBwClient) CreateNoteItem(folderID, name, notes string) error {
	return ErrAPINotImplemented
}

func (c *ApiBwClient) UpdateNoteItem(id, notes string) error {
	return ErrAPINotImplemented
}

func (c *ApiBwClient) Login(email, password, serverURL string) error {
	return ErrAPINotImplemented
}

func (c *ApiBwClient) Unlock(masterPassword string) error {
	return ErrAPINotImplemented
}
