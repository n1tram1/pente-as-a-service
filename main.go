package main

import (
    "fmt"
    "net/http"

    "github.com/n1tram1/pente-as-a-service/nominatim"

    "github.com/gin-gonic/gin"
)

type climb struct {
    Name string   `json:"name"`
    Grade float64 `json:"grade"`
}

func findClimbsIn(location string) []climb {
    bbox, err := nominatim.GetBbox(location)
    if err != nil {
        return nil
    }

    fmt.Printf("found bbox for %v: %+v", location, bbox)

    return []climb{}
}

func getClimbs(c *gin.Context) {
    location := c.Param("location")

    climbs := findClimbsIn(location)
    if climbs != nil {
        c.IndentedJSON(http.StatusOK, climbs)
        return
    }


    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "location not found"})
}

func main() {
    router := gin.Default()
    router.GET("/climbs/:location", getClimbs)

    router.Run("localhost:8080")
}
