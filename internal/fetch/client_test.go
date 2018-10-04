package fetch

import (
	"fmt"
	"testing"
	"time"

	"github.com/sand8080/d-data-transfer/internal/env"
)

func TestNewClient(t *testing.T) {
	NewClient(env.DataSourceUrl(), env.DataSourceAuthKey())
}

func TestClient_GetReports(t *testing.T) {
	cli := NewClient(env.DataSourceUrl(), env.DataSourceAuthKey())

	dateFrom := time.Now().Add(-24 * 3 * time.Hour)
	dateTo := dateFrom.Add(24 * time.Hour)

	resp, err := cli.GetReports(dateFrom, dateTo)
	fmt.Printf("resp: %v, err: %v\n", string(resp), err)
}
