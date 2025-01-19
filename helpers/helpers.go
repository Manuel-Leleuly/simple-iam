package helpers

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Manuel-Leleuly/simple-iam/constants"
	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
)

func GetPagination(fullUrl string) (*models.Pagination, error) {
	url, err := url.Parse(fullUrl)
	if err != nil {
		return nil, err
	}

	next := url.Path
	prev := url.Path
	queryParams := url.Query()
	for k, v := range queryParams {
		if k == "offset" {
			continue
		}

		// TODO: find a better way to check the question mark
		if strings.Contains(next, "?") {
			next += "&" + k + "=" + strings.Join(v, ",")
			prev += "&" + k + "=" + strings.Join(v, ",")
		} else {
			next += "?" + k + "=" + strings.Join(v, ",")
			prev += "?" + k + "=" + strings.Join(v, ",")
		}
	}

	selectedOffset := constants.DEFAULT_OFFSET
	if offset, ok := queryParams["offset"]; ok {
		selectedOffset, err = strconv.Atoi(offset[0])
		if err != nil {
			return nil, err
		}
	}

	selectedLimit := constants.DEFAULT_LIMIT
	if limit, ok := queryParams["limit"]; ok {
		selectedLimit, err = strconv.Atoi(limit[0])
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(next, "?") {
		next += "&offset=" + strconv.Itoa(selectedOffset+selectedLimit)
	} else {
		next += "?offset=" + strconv.Itoa(selectedOffset+selectedLimit)
	}

	if selectedOffset > selectedLimit {
		if strings.Contains(prev, "?") {
			prev += "&offset" + strconv.Itoa(selectedOffset-selectedLimit)
		} else {
			prev += "?offset" + strconv.Itoa(selectedOffset-selectedLimit)
		}
	}

	// reset prev to empty string if offset is 0
	if selectedOffset == 0 {
		prev = ""
	}

	// TODO: improve this
	// check if there is actually a next data
	var user models.User

	result := initializers.DB.Offset(selectedOffset + 1).First(&user)
	if result.Error != nil || user.Id == "" {
		next = ""
	}

	return &models.Pagination{Next: next, Prev: prev}, nil
}
