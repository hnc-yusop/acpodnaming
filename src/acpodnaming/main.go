package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"

	// for Kubernetes
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

type myValidServerhandler struct {
}

func (gs *myValidServerhandler) serve(w http.ResponseWriter, r *http.Request) {
	var Body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			Body = data
		}
	}

	if len(Body) == 0 {
		glog.Error("Unable to retrive Body from API")
		http.Error(w, "Empty Body", http.StatusBadRequest)
		return
	}

	glog.Info("Received Request")

	if r.URL.Path != "/validate" {
		glog.Error("Not a Validataion String")
		http.Error(w, "Not a Validataion String", http.StatusBadRequest)
		return
	}

	arRequest := &v1beta1.AdmissionReview{}
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

	arResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result:  &metav1.Status{Status: "Failure", Message: "the pod Naming is NOT up to the name standard", Code: 401},
			Allowed: false,
		},
	}

	podnamingReg := regexp.MustCompile(`kuku`)

	if podnamingReg.MatchString(string(pod.Name)) {
		fmt.Printf("the pod %s is up to the name standard", pod.Name)
		arResponse.Response.Allowed = true
	}

	arResponse.APIVersion = "admission.k8s.io/v1"
	arResponse.Kind = arRequest.Kind
	arResponse.Response.UID = arRequest.Request.UID

	resp, err := json.Marshal(arResponse)

	if err != nil {
		glog.Error("Can't encode response:", err)
		http.Error(w, fmt.Sprintf("couldn't encode response: %v", err), http.StatusInternalServerError)
	}

	if _, err := w.Write(resp); err != nil {
		glog.Error("Can't write response", err)
		http.Error(w, fmt.Sprintf("cloud not write response: %v", err), http.StatusInternalServerError)
	}
}
