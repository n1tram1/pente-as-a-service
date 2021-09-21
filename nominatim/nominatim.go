package nominatim

import (
    "net/http"
    "net/url"
    "encoding/json"
    "io"
    "errors"
    "strconv"
    "fmt"
)

type SearchResult struct {
    Address struct {
        Bakery        string `json:"bakery"`
        CityDistrict  string `json:"city_district"`
        Continent     string `json:"continent"`
        Country       string `json:"country"`
        CountryCode   string `json:"country_code"`
        Footway       string `json:"footway"`
        Neighbourhood string `json:"neighbourhood"`
        Postcode      string `json:"postcode"`
        State         string `json:"state"`
        Suburb        string `json:"suburb"`
    } `json:"address"`
    Boundingbox []string `json:"boundingbox"`
    Class       string   `json:"class"`
    DisplayName string   `json:"display_name"`
    Icon        string   `json:"icon"`
    Importance  float64  `json:"importance"`
    Lat         string   `json:"lat"`
    Licence     string   `json:"licence"`
    Lon         string   `json:"lon"`
    OsmID       int   `json:"osm_id"`
    OsmType     string   `json:"osm_type"`
    PlaceID     int   `json:"place_id"`
    Type        string   `json:"type"`
}

type SearchResults struct {
    results []SearchResult
}

type BoundingBox struct {
    SouthLatitude float64
    NorthLatitude float64
    WestLongitude float64
    EastLatitude float64
}

func search(location string) (*SearchResult, error) {
    params := url.Values{}

    params.Add("addressdetails", "1")
    params.Add("q", location)
    params.Add("format", "json")
    params.Add("limit", "1") // TODO: should we limit to 1 result ?

    resp, err := http.Get("https://nominatim.openstreetmap.org/?" + params.Encode())
    if err != nil {
        return nil, err
    }

    body, err := io.ReadAll(resp.Body)

    search_results := []SearchResult{}
    err = json.Unmarshal(body, &search_results)
    if err != nil {
        return nil, err
    }

    return &search_results[0], nil
}

func GetBbox(location string) (*BoundingBox, error) {
    search_result, err := search(location)
    fmt.Printf("search returning: %+v, %+v\n", search_result, err)
    if err != nil {
        return nil, err
    }

    if len(search_result.Boundingbox) != 4 {
        // TODO: return an error
        return nil, errors.New("Boundingbox did not countains 4 points")
    }

    south_lat, err := strconv.ParseFloat(search_result.Boundingbox[0], 64)
    if err != nil {
        return nil, err
    }

    north_latitude, err := strconv.ParseFloat(search_result.Boundingbox[1], 64)
    if err != nil {
        return nil, err
    }

    west_longitude, err := strconv.ParseFloat(search_result.Boundingbox[2], 64)
    if err != nil {
        return nil, err
    }

    east_longitude, err := strconv.ParseFloat(search_result.Boundingbox[3], 64)
    if err != nil {
        return nil, err
    }

    return  &BoundingBox {
        south_lat,
        north_latitude,
        west_longitude,
        east_longitude,
    }, nil
}
