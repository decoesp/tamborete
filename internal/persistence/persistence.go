package persistence

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/decoesp/tamborete/internal/database"
)

type Persistence struct {
	filePath string
	db       *database.Database
	mu       sync.Mutex
}

func New(filePath string, db *database.Database) *Persistence {
	return &Persistence{
		filePath: filePath,
		db:       db,
	}
}

func (p *Persistence) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data := struct {
		Strings map[string]string   `json:"strings"`
		Lists   map[string][]string `json:"lists"`
	}{
		Strings: p.db.GetStrings(),
		Lists:   p.db.GetLists(),
	}

	file, err := os.Create(p.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (p *Persistence) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.Open(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	data := struct {
		Strings map[string]string   `json:"strings"`
		Lists   map[string][]string `json:"lists"`
	}{}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	p.db.Lock()
	defer p.db.Unlock()
	p.db.SetStrings(data.Strings)
	p.db.SetLists(data.Lists)
	return nil
}
