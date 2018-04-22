// Copyright 2018 IBM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"fmt"
	"io"
	"log"
	"sort"
	"text/tabwriter"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DumpEvents returns a reader that will have events for a given namespace written to
func (f *Framework) DumpEvents(namespace string) io.Reader {
	fmt.Printf("Dumping events for namespace: %v\n", namespace)
	events, err := f.KubeClient.CoreV1().Events(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("error retrieving events from %q: %v", namespace, err)
	}

	sort.Sort(evs(events.Items))

	reader, writer := io.Pipe()
	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	go func() {
		fmt.Fprint(w, "\nTIME\tNAME\tKIND\tREASON\tSOURCE\tMESSAGE\n")
		for _, event := range events.Items {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
				event.FirstTimestamp,
				event.Name,
				event.InvolvedObject.Kind,
				event.Reason,
				event.Source,
				event.Message,
			)
		}
		w.Flush()
		writer.Close()
	}()
	return reader
}

type evs []corev1.Event

func (e evs) Len() int {
	return len(e)
}

func (e evs) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e evs) Less(i, j int) bool {
	return e[i].LastTimestamp.UnixNano() < e[j].LastTimestamp.UnixNano()
}

// DumpPolicies returns a reader that will have all cluster and image policies present in it
func (f *Framework) DumpPolicies(namespace string) io.Reader {
	fmt.Printf("Dumping cluster policies and policies for namespace: %v\n", namespace)

	clusterImagePolicies, err := f.ListClusterImagePolicies()
	if err != nil {
		log.Printf("error listing ClusterImagePolicies: %v", err)
	}
	imagePolicies, err := f.ListImagePolicies(namespace)
	if err != nil {
		log.Printf("error listing ImagePolicies in %q: %v", namespace, err)
	}

	reader, writer := io.Pipe()
	go func() {
		fmt.Fprint(writer, "\nClusterImagePolicies Present:\n")
		for _, clusterImagePolicy := range clusterImagePolicies.Items {
			fmt.Fprintf(writer, "- %v\n", clusterImagePolicy.Name)
		}
		fmt.Fprint(writer, "\nImagePolicies Present:\n")
		for _, imagePolicy := range imagePolicies.Items {
			fmt.Fprintf(writer, "- %v\n", imagePolicy.Name)
		}
		writer.Close()
	}()
	return reader

}
