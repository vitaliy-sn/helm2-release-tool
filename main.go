package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"k8s.io/api/core/v1"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
	// yaml "gopkg.in/yaml.v2"
)

func main() {
	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	var data string

	for j := 0; j < len(output); j++ {
		data = data + string(output[j])
	}

	action := os.Args[1]

	switch action {
	case "info":
		releaseInfo(data)
	case "get-manifests":
		getManifest(data)
	case "set-status-deployed":
		setReleaseStatus(data)
	case "set-release-name":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "ERROR: new release name not specify\n")
		}
		setReleaseName(data, os.Args[2])
	case "set-release-namespace":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "ERROR: new release namespace not specify\n")
		}
		setReleaseNamespace(data, os.Args[2])
	default:
		fmt.Fprintf(os.Stderr, "ERROR: invalid action: %s\nYou need to specify action: info or set-status-deployed or set-new-release-name or set-new-release-namespace\n", action)
		os.Exit(1)
	}

}

func releaseInfo(data string) {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, _, _ := decoder.Decode([]byte(data), nil, nil)
	cm := obj.(*v1.ConfigMap)

	release, _ := DecodeRelease(cm.Data["release"])
	fmt.Printf("Name: %v\nNamespace: %v\nStatus: %v\n", release.Name, release.Namespace, release.Info.Status.Code)
}

func getManifest(data string) {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, _, _ := decoder.Decode([]byte(data), nil, nil)
	cm := obj.(*v1.ConfigMap)

	release, _ := DecodeRelease(cm.Data["release"])
	fmt.Printf("%v\n", release.Manifest)
}

func setReleaseStatus(data string) {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, _, _ := decoder.Decode([]byte(data), nil, nil)
	cm := obj.(*v1.ConfigMap)

	release, _ := DecodeRelease(cm.Data["release"])
	cm.ObjectMeta.Labels["STATUS"] = "DEPLOYED"
	// https://github.com/helm/helm/blob/release-2.16/pkg/proto/hapi/release/status.pb.go#L47
	release.Info.Status.Code = 1

	encoded, _ := EncodeRelease(release)

	cm.Data["release"] = encoded

	jsonEncoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	_ = jsonEncoder.Encode(cm, os.Stdout)
}

func setReleaseName(data, newReleaseName string) {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, _, _ := decoder.Decode([]byte(data), nil, nil)
	cm := obj.(*v1.ConfigMap)

	decoded, _ := DecodeRelease(cm.Data["release"])
	decoded.Name = newReleaseName
	encoded, _ := EncodeRelease(decoded)

	cm.Data["release"] = encoded

	currentName := cm.ObjectMeta.Labels["NAME"]

	cm.Name = strings.Replace(cm.Name, currentName, newReleaseName, 1)
	cm.ObjectMeta.Labels["NAME"] = newReleaseName
	cm.SelfLink = strings.Replace(cm.SelfLink, currentName, newReleaseName, 1)

	jsonEncoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	_ = jsonEncoder.Encode(cm, os.Stdout)
}

func setReleaseNamespace(data, newReleaseNamespace string) {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, _, _ := decoder.Decode([]byte(data), nil, nil)
	cm := obj.(*v1.ConfigMap)

	decoded, _ := DecodeRelease(cm.Data["release"])
	decoded.Namespace = newReleaseNamespace
	encoded, _ := EncodeRelease(decoded)

	cm.Data["release"] = encoded

	jsonEncoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	_ = jsonEncoder.Encode(cm, os.Stdout)
}
