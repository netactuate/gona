package gona

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func checkNKEDrop1Request(t *testing.T, r *http.Request, method, path, body string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Errorf("request %s %s, want %s %s", r.Method, r.URL.Path, method, path)
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if body == "" {
		if len(raw) > 0 && string(raw) != "null" {
			t.Errorf("unexpected body %s", raw)
		}
		return
	}
	var got, want interface{}
	if err = json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	json.Unmarshal([]byte(body), &want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("body %s, want %s", raw, body)
	}
}
func TestListNKENodePoolsContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/nke/clusters/1/node-pools", ``)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"meta":{"total":1,"limit":1000,"offset":0},"data":[{"nodePoolId":2,"name":"nah-test","mode":"managed","nodes":{"ready":0,"total":0}}]}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			_, err := c.ListNKENodePools(1)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestGetNKENodePoolContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/nke/clusters/1/node-pools", ``)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"meta":{"total":1,"limit":1000,"offset":0},"data":[{"nodePoolId":2,"name":"nah-test","mode":"managed","nodes":{"ready":0,"total":0}}]}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			_, err := c.GetNKENodePool(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestCreateNKENodePoolContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "POST", "/nke/clusters/1/node-pools", `{"name":"nah-test","mode":"managed","minimumNodes":0}`)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"nodePoolId":2}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			_, err := c.CreateNKENodePool(1, &CreateNKENodePoolRequest{Name: "nah-test", Mode: "managed", MinimumNodes: &zero})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestUpdateNKENodePoolContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "PATCH", "/nke/clusters/1/node-pools/2", `{"minimumNodes":0,"autoscalingEnabled":false,"labels":{},"taints":[]}`)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"nodePoolId":2}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.UpdateNKENodePool(1, 2, &UpdateNKENodePoolRequest{MinimumNodes: &zero, AutoscalingEnabled: &disabled, Labels: &labels, Taints: &taints})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestSetNKENodePoolSizeContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "PUT", "/nke/clusters/1/node-pools/2/size", `{"desiredNodes":0}`)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"nodePoolId":2}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.SetNKENodePoolSize(1, 2, 0)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestDeleteNKENodePoolContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "DELETE", "/nke/clusters/1/node-pools/2", ``)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"nodePoolId":2}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.DeleteNKENodePool(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestRemoveNKEWorkerNodeContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					checkNKEDrop1Request(t, r, "GET", "/nke/clusters/1/worker-nodes/2", "")
					w.Write([]byte(`{"code":200,"data":{"workerNodeId":2,"poolId":4}}`))
					return
				}
				checkNKEDrop1Request(t, r, "POST", "/nke/clusters/1/worker-nodes/2/remove", `{"nodePoolId":4}`)
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"nodePoolId":2}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.RemoveNKEWorkerNode(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestListNKEStorageAvailabilityContract(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	_ = zero
	_ = disabled
	_ = labels
	_ = taints
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "404"}[missing], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/nke/storage-availability", ``)
				if r.URL.Query().Get("locationId") != "3" {
					t.Errorf("missing locationId: %s", r.URL)
				}
				if missing {
					w.WriteHeader(404)
					w.Write([]byte(`{"code":404,"data":{}}`))
					return
				}
				w.Write([]byte(`{"code":200,"data":{"meta":{"total":1,"limit":1000,"offset":0},"data":[{"locationId":3,"storageAvailable":true}]}}`))
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			_, err := c.ListNKEStorageAvailability(3)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404 not classified: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
