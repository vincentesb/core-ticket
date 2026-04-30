package dto

import (
	"errors"
	"strings"
)

type SortDto struct {
	Sort      string   `json:"sort,omitempty" form:"sort"`
	ValidSort []string `json:"-"`
}

func (sd *SortDto) SetDefaultSort(sort string) {
	if sd.Sort == "" {
		sd.Sort = sort
	}
}

func (sd *SortDto) ValidateSort() error {
	for _, vs := range sd.ValidSort {
		if strings.TrimPrefix(sd.Sort, "-") == vs {
			return nil
		}
	}
	return errors.New("invalid sort value")
}
