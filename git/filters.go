package git

import (
	"fmt"
	"regexp"
)

var (
	valueRegexp = regexp.MustCompile("\\Agit[\\-\\s]media")
)

type Filter struct {
	Name  string
	Value string
}

func (f *Filter) Install() error {
	key := fmt.Sprintf("filter.lfs.%s", f.Name)

	currentValue := Config.Find(key)
	if force || f.shouldReset(currentValue) {
		Config.UnsetGlobal(key)
		Config.SetGlobal(key, f.Value)

		return nil
	} else if currentVal != f.Value {
		return fmt.Errorf("The %s filter should be \"%s\" but is \"%s\"",
			filterName, f.Value, existing)
	}

	return nil
}

func (f *Filter) shouldReset(value string) bool {
	if len(value) == 0 {
		return true
	}
	return valueRegexp.MatchString(value)
}

type Filters []*Filter

func (fs *Filters) Setup() error {
	for _, f := range fs {
		if err := f.Install(); err != nil {
			return err
		}
	}

	return nil

}

func (fs *Filters) Teardown() {
	Config.UnsetGlobalSection("filters.lfs")
}
