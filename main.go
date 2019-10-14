// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/google/pprof/profile"
	"google.golang.org/api/option"
	gtransport "google.golang.org/api/transport/grpc"
	pb "google.golang.org/genproto/googleapis/devtools/cloudprofiler/v2"
)

var (
	client pb.ProfilerServiceClient

	project string
	zone    string
	target  string
	version string

	input    string
	keepTime bool
)

const (
	apiAddr = "cloudprofiler.googleapis.com:443"
	scope   = "https://www.googleapis.com/auth/monitoring.write"
)

func main() {
	ctx := context.Background()
	flag.StringVar(&project, "project", "", "")
	flag.StringVar(&zone, "zone", "", "")
	flag.StringVar(&target, "target", "", "")
	flag.StringVar(&version, "version", "", "")
	flag.StringVar(&input, "i", "pprof.out", "")
	flag.BoolVar(&keepTime, "keep-time", false, "")
	flag.Usage = usageAndExit
	flag.Parse()

	// TODO(jbd): Automatically detect input. Don't convert if pprof.

	if project == "" {
		id, err := metadata.ProjectID()
		if err != nil {
			log.Fatalf("Cannot resolve the GCP project from the metadata server: %v", err)
		}
		project = id
	}
	if zone == "" {
		// Ignore error. If we cannot resolve the instance name,
		// it would be too aggressive to fatal exit.
		zone, _ = metadata.Zone()
	}

	if target == "" {
		target = input
	}

	conn, err := gtransport.Dial(ctx,
		option.WithEndpoint(apiAddr),
		option.WithScopes(scope))
	if err != nil {
		log.Fatal(err)
	}
	client = pb.NewProfilerServiceClient(conn)

	pprofBytes, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("Cannot convert perf data to pprof: %v", err)
	}

	if err := upload(ctx, pprofBytes); err != nil {
		log.Fatalf("Cannot upload to Stackdriver Profiler: %v", err)
	}
	fmt.Printf("https://console.cloud.google.com/profiler/%s;type=%s?project=%s\n", url.PathEscape(target), pb.ProfileType_CPU, project)
}

func upload(ctx context.Context, payload []byte) error {
	// Reset time, otherwise old profiles wont be shown
	// at Cloud profiler due to data retention limits.
	if !keepTime {
		var err error
		payload, err = resetTime(payload)
		if err != nil {
			log.Printf("Cannot reset the profile's time: %v", err)
		}
	}

	req := &pb.CreateOfflineProfileRequest{
		Parent: "projects/" + project,
		Profile: &pb.Profile{
			// TODO(jbd): Guess the profile type from the input.
			ProfileType: pb.ProfileType_CPU,
			Deployment: &pb.Deployment{
				ProjectId: project,
				Target:    target,
				Labels: map[string]string{
					"zone":    zone,
					"version": version,
				},
			},
			ProfileBytes: payload,
		},
	}

	// TODO(jbd): Is there a way without having
	// to load the profile all in memory?
	_, err := client.CreateOfflineProfile(ctx, req)
	return err
}

// TODO(jbd): Make it optional.
func resetTime(pprofBytes []byte) ([]byte, error) {
	p, err := profile.ParseData(pprofBytes)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse the profile: %v", err)
	}
	p.TimeNanos = time.Now().UnixNano()

	var buf bytes.Buffer
	if err := p.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TODO(jbd): Check binary dependencies and install if not available.

const usageText = `pprof-upload [-i pprof.out]

Other options:
-project    Google Cloud project name, tries to automatically
            resolve if none is set.
-zone       Google Cloud zone, tries to automatically resolve if
		    none is set.
-target     Target profile name to upload data to.
-version    Version of the profiled program.
-keep-time  When set, keeps the original time info from the profile file.
			Due to data retention limits, Stackdriver Profiler won't
            show data older than 30 days. By default, false.
`

func usageAndExit() {
	fmt.Println(usageText)
	os.Exit(1)
}
