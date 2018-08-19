package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func Test_addCronTask(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	r, err := inst.NewRequest("GET", "/add-cron-task", nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}

	w := httptest.NewRecorder()
	addCronTask(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
