package helpers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Manuel-Leleuly/simple-iam/constants"
	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
)

func GetFullUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host + c.Request.URL.String()
}

func GetPagination(fullUrl string) (*models.Pagination, error) {
	selectedUrl, err := url.Parse(fullUrl)
	fmt.Println("full url:", fullUrl)
	if err != nil {
		return nil, err
	}

	next := selectedUrl.Path
	prev := selectedUrl.Path
	queryParams := selectedUrl.Query()

	nextQueryParams := url.Values{}
	prevQueryParams := url.Values{}

	for k, v := range queryParams {
		if k == "offset" {
			continue
		}

		nextQueryParams.Add(k, strings.Join(v, ","))
		prevQueryParams.Add(k, strings.Join(v, ","))
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

	nextQueryParams.Add("offset", strconv.Itoa(selectedOffset+selectedLimit))

	if selectedOffset > selectedLimit {
		prevQueryParams.Add("offset", strconv.Itoa(selectedOffset-selectedLimit))
	}

	if len(nextQueryParams) > 0 {
		next = next + "?" + nextQueryParams.Encode()
	}

	if len(prevQueryParams) > 0 {
		prev = prev + "?" + prevQueryParams.Encode()
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
