// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017-present Datadog, Inc.

//go:build !clusterchecks && !kubeapiserver

package v1

import (
	workloadmeta "github.com/DataDog/datadog-agent/comp/core/workloadmeta/def"
	"github.com/gorilla/mux"
)

func installCloudFoundryMetadataEndpoints(_ *mux.Router) {}

func installKubernetesMetadataEndpoints(_ *mux.Router, _ workloadmeta.Component) {}
