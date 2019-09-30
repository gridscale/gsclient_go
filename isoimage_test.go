package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetISOImageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareISOImageHTTPGetList())
	})
	res, err := client.GetISOImageList(emptyCtx)
	assert.Nil(t, err, "GetISOImageList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockISOImage()), fmt.Sprintf("%v", res))
}

func TestClient_GetISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareISOImageHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetISOImage(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetISOImage returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockISOImage("active")), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockISOImage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateISOImage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		uri := path.Join(apiISOBase)
		var isFailed bool
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprintf(writer, prepareISOImageHTTPCreateResponse())
			}
		})
		if clientTest {
			httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
			mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			response, err := client.CreateISOImage(
				emptyCtx,
				ISOImageCreateRequest{
					Name:         "Test",
					SourceURL:    "http://example.org",
					Labels:       []string{"label"},
					LocationUUID: "aa-bb-cc",
				})
			if test.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateISOImage returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockISOImageCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
		server.Close()
	}
}

func TestClient_UpdateISOImage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiISOBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			if isFailed {
				writer.WriteHeader(400)
			} else {
				if request.Method == http.MethodPatch {
					fmt.Fprintf(writer, "")
				} else if request.Method == http.MethodGet {
					fmt.Fprint(writer, prepareISOImageHTTPGet("active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.UpdateISOImage(
					emptyCtx,
					test.testUUID,
					ISOImageUpdateRequest{
						Name:   "test",
						Labels: []string{},
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateISOImage returned an error %v", err)
				}
			}
		}
		server.Close()
	}

func TestClient_DeleteISOImage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiISOBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			if isFailed {
				writer.WriteHeader(400)
			} else {
				if request.Method == http.MethodDelete {
					fmt.Fprintf(writer, "")
				} else if request.Method == http.MethodGet {
					writer.WriteHeader(404)
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.DeleteISOImage(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteISOImage returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_UpdateISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetISOImageEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetISOImageEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
}

	err := client.UpdateISOImage(dummyUUID, ISOImageUpdateRequest{
		Name:   "test",
		Labels: []string{},
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetISOImagesByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetISOImagesByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockISOImage("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_DeleteISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	res, err := client.GetDeletedISOImages(emptyCtx)
	assert.Nil(t, err, "GetDeletedISOImages returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockISOImage("deleted")), fmt.Sprintf("%v", res))
}

func TestClient_GetISOImageEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareISOImageHTTPGetEventList())
	})
	err := client.waitForISOImageActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForISOImageActive returned an error %v", err)
}

func TestClient_waitForISOImageDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(404)
	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForISOImageDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForISOImageDeleted returned an error %v", err)
		}
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockISOImageEvent()), fmt.Sprintf("%v", res))
}

func getMockISOImage() ISOImage {
	mock := ISOImage{Properties: ISOImageProperties{
		ObjectUUID: dummyUUID,
		Relations: ISOImageRelation{
			Servers: []ServerinISOImage{
				{
					Bootdevice: true,
					CreateTime: dummyTime,
					ObjectName: "test",
					ObjectUUID: "abc-def",
				},
			},
		},
		Description:     "description",
		LocationName:    "locationName",
		SourceURL:       "url",
		Labels:          []string{"label"},
		LocationIata:    "iata",
		LocationUUID:    "locUUID",
		Status:          "active",
		CreateTime:      dummyTime,
		Name:            "test",
		Version:         "1.0",
		LocationCountry: "Country",
		UsageInMinutes:  10,
		Private:         false,
		ChangeTime:      dummyTime,
		Capacity:        10,
		CurrentPrice:    9.99,
	}}
	return mock
}

func prepareISOImageHTTPGetList() string {
	iso := getMockISOImage()
	res, _ := json.Marshal(iso.Properties)
	return fmt.Sprintf(`{"isoimages": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareISOImageHTTPGet() string {
	iso := getMockISOImage()
	res, _ := json.Marshal(iso)
	return string(res)
}

func getMockISOImageCreateResponse() ISOImageCreateResponse {
	mock := ISOImageCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareISOImageHTTPCreateResponse() string {
	res, _ := json.Marshal(getMockISOImageCreateResponse())
	return string(res)
}

func getMockISOImageEvent() ISOImageEvent {
	mock := ISOImageEvent{Properties: ISOImageEventProperties{
		ObjectType:    "type",
		RequestUUID:   dummyRequestUUID,
		ObjectUUID:    dummyUUID,
		Activity:      "activity",
		RequestType:   "request type",
		RequestStatus: "active",
		Change:        "change description",
		Timestamp:     dummyTime,
		UserUUID:      "user-id",
	}}
	return mock
}

func prepareISOImageHTTPGetEventList() string {
	res, _ := json.Marshal(getMockISOImageEvent().Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}
