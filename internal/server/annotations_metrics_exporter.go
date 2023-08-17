package server

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	annotationsTagsPath = "/api/annotations/tags"
	methodGet           = "GET"
	limit               = "100"
)

type AnnotationsTagsMetricExporter struct {
	Logger log.Logger
	Opts   Opts
	Client *GrafanaClient
}

type AnnotationTags struct {
	Result struct {
		Tags []struct {
			Tag   string `json:"tag"`
			Count int    `json:"count"`
		} `json:"tags"`
	} `json:"result"`
}

func NewAnnotationsTagsMetricExporter(logger log.Logger, opts Opts) (*AnnotationsTagsMetricExporter, error) {
	cli, err := NewGrafanaClient(opts)
	if err != nil {
		return nil, err
	}

	return &AnnotationsTagsMetricExporter{
		Logger: logger,
		Opts:   opts,
		Client: cli,
	}, nil
}

func (a *AnnotationsTagsMetricExporter) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(a.Opts.AnnotationAPIPollSeconds) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				a.collectAnnotations(ctx)
			case <-ctx.Done():
				_ = level.Info(a.Logger).Log("msg", "stopped polling for runner metrics")
				return
			}
		}
	}()
}

func (a *AnnotationsTagsMetricExporter) collectAnnotations(ctx context.Context) {
	params := url.Values{}
	params.Add("limit", limit)
	result := AnnotationTags{}
	err := a.Client.Request(ctx, methodGet, annotationsTagsPath, params, &result)

	if err != nil {
		level.Error(a.Logger).Log("Error performing request:", err)
	}

	if len(result.Result.Tags) == 0 {
		annotationsTagsTotal.Reset()
	} else {
		for _, tag := range result.Result.Tags {
			annotationsTagsTotal.WithLabelValues(tag.Tag).Set(float64(tag.Count))
		}
	}
}
