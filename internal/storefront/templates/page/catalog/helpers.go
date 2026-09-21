package catalog

import (
	"fmt"
	"net/url"

	"github.com/dv1x3r/chinook-go/internal/storefront/services"
	"github.com/dv1x3r/chinook-go/internal/storefront/templates/component/table"
)

// filtersVals sends the checked column filters the table model keeps along with every table request
const filtersVals = "js:{ filters: JSON.stringify(Alpine.$data(document.getElementById('hx-table-model'))._model.filters) }"

func genreFilters(genres []services.Genre) []table.Filter {
	filters := make([]table.Filter, 0, len(genres))
	for _, genre := range genres {
		filters = append(filters, table.Filter{Label: genre.Name, Value: fmt.Sprint(genre.ID), Icon: "fa fa-tag"})
	}
	return filters
}

func mediaFilters(mediaTypes []services.MediaType) []table.Filter {
	filters := make([]table.Filter, 0, len(mediaTypes))
	for _, media := range mediaTypes {
		filters = append(filters, table.Filter{Label: media.Name, Value: fmt.Sprint(media.ID), Icon: "fa fa-file-audio"})
	}
	return filters
}

// length turns the stored milliseconds into the usual minutes and seconds
func length(milliseconds int) string {
	seconds := milliseconds / 1000
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}

func size(bytes int) string {
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}

func plural(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func searchLink(props Props, by string, value string) string {
	query := url.Values{}
	query.Set("search", value)
	query.Set("search-by", by)
	return props.Path + "?" + query.Encode()
}
