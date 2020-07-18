package main

import (
	  "fmt"
	  "encoding/json"
	  "io/ioutil"
	  "net/http"
	  "github.com/golang/glog"

	  // for Kubernetes 
	  "k8s.io/api/admission/v1beta1"
	  "k8s.io/api/core/v1"
	  matav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type myValidServerhandler struct {

}

func (gs *myValidServerhandler) serve(w http.ResponseWriter, r *http.Request) {
	var Body []byte
	if r.Body != nil {
		if data , err := ioutil.ReadAll(r.Body); err == nil {
			Body = data
		}
	}

	if len(Body) == 0 {
		glog.Error("Unable to retrive Body from API")
		http.Error(w,"Empty Body", http.StatusBadRequest)
		return
	}

	glog.Info("Received Request")

	if r.URL.Path != "/validate" {
		glog.Error("Not a Validataion String")
		http.Error(w,"Not a Validataion String", http.StatusBadRequest)
		return
	}

	arRequest := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(Body, arRequest); err != nil {
		glog.Error("incorrect Body")
		http.Error(w, "incorrect Body", http.StatusBadRequest)
		return
	}

	raw := arRequest.Request.Object.Raw
	pod := v1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		glog.Error("Error Deserializing Pod")
		return
	}

	if pod.Name == "my-valid-pod" {
		return
	}

	arResponse := v1beta1.AdmissionReview {
		Response: &v1beta1.AdmissionResponse {
			Allowed: false,
			Result: &matav1.Status{
				Message: "we are just making sure the Pod has the right Name",
			},
		},
	}

	resp, err := json.Marshal(arResponse)
	if err != nil {
		glog.Error("Can't encode response:", err)
		http.Error(w, fmt.Sprintf("couldn't encode response: %v", err), http.StatusInternalServerError)
	}

	glog.Infof("Ready to write  response ...")
	if _, err := w.Write(resp); err != nil {
		glog.Error("Can't write response", err)
		http.Error(w, fmt.Sprintf("cloud not write response: %v", err), http.StatusInternalServerError)
	}
}