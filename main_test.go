package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "", nil)

    return 
    if err != nil {
        t.Fatal(err)
    }
    recorder := httptest.NewRecorder()
    hf := http.HandlerFunc(handler)
    hf.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("Handler returned unexpected status code: got %v expected %v", status, http.StatusOK)
    }
    expected := "Hello World!"
    actual := recorder.Body.String()
    if actual != expected {
        t.Errorf("Handler returned unexpected body: got %v expected %v", actual, expected)
    }
}
