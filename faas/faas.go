// Copyright 2021 Nitric Technologies Pty Ltd.
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

package faas

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type HandlerBuilder interface {
	Http(string, ...HttpMiddleware) HandlerBuilder
	Event(...EventMiddleware) HandlerBuilder
	BucketNotification(...BucketNotificationMiddleware) HandlerBuilder
	Websocket(...WebsocketMiddleware) HandlerBuilder
	Default(...TriggerMiddleware) HandlerBuilder
	WithApiWorkerOpts(ApiWorkerOptions) HandlerBuilder
	WithRateWorkerOpts(RateWorkerOptions) HandlerBuilder
	WithCronWorkerOpts(CronWorkerOptions) HandlerBuilder
	WithSubscriptionWorkerOpts(SubscriptionWorkerOptions) HandlerBuilder
	WithBucketNotificationWorkerOptions(BucketNotificationWorkerOptions) HandlerBuilder
	WithWebsocketWorkerOptions(WebsocketWorkerOptions) HandlerBuilder
	Start() error
	String() string
}

type HandlerProvider interface {
	GetHttp(method string) HttpMiddleware
	GetEvent() EventMiddleware
	GetBucketNotification() BucketNotificationMiddleware
	GetWebsocket() WebsocketMiddleware
	GetDefault() TriggerMiddleware
}

type faasClientImpl struct {
	apiWorkerOpts                ApiWorkerOptions
	rateWorkerOpts               RateWorkerOptions
	cronWorkerOpts               CronWorkerOptions
	subscriptionWorkerOpts       SubscriptionWorkerOptions
	bucketNotificationWorkerOpts BucketNotificationWorkerOptions
	websocketWorkerOpts          WebsocketWorkerOptions

	http               map[string]HttpMiddleware
	event              EventMiddleware
	bucketNotification BucketNotificationMiddleware
	websocket          WebsocketMiddleware
	trig               TriggerMiddleware
}

func (f *faasClientImpl) String() string {
	out := []string{}

	if f.apiWorkerOpts.ApiName != "" {
		methods := []string{}
		for k := range f.http {
			methods = append(methods, k)
		}
		sort.Strings(methods)
		out = append(out, fmt.Sprintf("Api:%s, path:%s methods:[%s]", f.apiWorkerOpts.ApiName, f.apiWorkerOpts.Path, strings.Join(methods, ",")))
	}
	if f.rateWorkerOpts.Frequency != "" {
		out = append(out, fmt.Sprintf("Rate:%d, Freq:%d", f.rateWorkerOpts.Rate, f.rateWorkerOpts.Rate))
	}
	if f.cronWorkerOpts.Cron != "" {
		out = append(out, fmt.Sprintf("Cron:%s", f.cronWorkerOpts.Cron))
	}
	if f.subscriptionWorkerOpts.Topic != "" {
		out = append(out, fmt.Sprintf("Subscribe:%s", f.subscriptionWorkerOpts.Topic))
	}
	if f.bucketNotificationWorkerOpts.Bucket != "" {
		out = append(out, fmt.Sprintf("Bucket:%s, Type:%s, Filter:%s", f.bucketNotificationWorkerOpts.Bucket, f.bucketNotificationWorkerOpts.NotificationType, f.bucketNotificationWorkerOpts.NotificationPrefixFilter))
	}
	if f.websocketWorkerOpts.Socket != "" {
		out = append(out, fmt.Sprintf("Socket:%s, Type:%s", f.websocketWorkerOpts.Socket, f.websocketWorkerOpts.EventType))
	}

	return strings.Join(out, "\n")
}

func (f *faasClientImpl) Http(method string, mwares ...HttpMiddleware) HandlerBuilder {
	f.http[method] = ComposeHttpMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetHttp(method string) HttpMiddleware {
	if _, ok := f.http[method]; !ok {
		return nil
	}
	return f.http[method]
}

func (f *faasClientImpl) Event(mwares ...EventMiddleware) HandlerBuilder {
	f.event = ComposeEventMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetEvent() EventMiddleware {
	return f.event
}

func (f *faasClientImpl) BucketNotification(mwares ...BucketNotificationMiddleware) HandlerBuilder {
	f.bucketNotification = ComposeBucketNotificationMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetBucketNotification() BucketNotificationMiddleware {
	return f.bucketNotification
}

func (f *faasClientImpl) Websocket(mwares ...WebsocketMiddleware) HandlerBuilder {
	f.websocket = ComposeWebsocketMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetWebsocket() WebsocketMiddleware {
	return f.websocket
}

func (f *faasClientImpl) Default(mwares ...TriggerMiddleware) HandlerBuilder {
	f.trig = ComposeTriggerMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetDefault() TriggerMiddleware {
	return f.trig
}

func (f *faasClientImpl) Start() error {
	// Fail if no handlers were provided
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return errors.NewWithCause(
			codes.Unavailable,
			"faas.Start: Unable to reach FaasServiceServer",
			err,
		)
	}

	fsc := v1.NewFaasServiceClient(conn)

	return f.startWithClient(fsc)
}

func (f *faasClientImpl) startWithClient(fsc v1.FaasServiceClient) error {
	// Fail if no handlers were provided
	if len(f.http) == 0 && f.event == nil && f.trig == nil && f.websocket == nil && f.bucketNotification == nil {
		return fmt.Errorf("no valid handlers provided")
	}

	if stream, err := fsc.TriggerStream(context.TODO()); err == nil {
		initRequest := &v1.InitRequest{}

		if len(f.http) > 0 {
			methods := []string{}
			for k := range f.http {
				methods = append(methods, k)
			}
			sort.Strings(methods)

			sec := map[string]*v1.ApiWorkerScopes{}
			for k, v := range f.apiWorkerOpts.Security {
				sec[k] = &v1.ApiWorkerScopes{
					Scopes: v,
				}
			}

			initRequest.Worker = &v1.InitRequest_Api{
				Api: &v1.ApiWorker{
					Api:     f.apiWorkerOpts.ApiName,
					Path:    f.apiWorkerOpts.Path,
					Methods: methods,
					Options: &v1.ApiWorkerOptions{
						SecurityDisabled: f.apiWorkerOpts.SecurityDisabled,
						Security:         sec,
					},
				},
			}
		}
		if f.bucketNotificationWorkerOpts.Bucket != "" {
			notificationType, err := f.bucketNotificationWorkerOpts.notificationTypeToWire()
			if err != nil {
				return err
			}

			initRequest.Worker = &v1.InitRequest_BucketNotification{
				BucketNotification: &v1.BucketNotificationWorker{
					Bucket: f.bucketNotificationWorkerOpts.Bucket,
					Config: &v1.BucketNotificationConfig{
						NotificationType:         notificationType,
						NotificationPrefixFilter: f.bucketNotificationWorkerOpts.NotificationPrefixFilter,
					},
				},
			}
		}
		if f.rateWorkerOpts.Rate > 0 {
			initRequest.Worker = &v1.InitRequest_Schedule{
				Schedule: &v1.ScheduleWorker{
					Key: f.rateWorkerOpts.Description,
					Cadence: &v1.ScheduleWorker_Rate{
						Rate: &v1.ScheduleRate{
							Rate: fmt.Sprintf("%d %s", f.rateWorkerOpts.Rate, string(f.rateWorkerOpts.Frequency)),
						},
					},
				},
			}
		}
		if f.cronWorkerOpts.Cron != "" {
			initRequest.Worker = &v1.InitRequest_Schedule{
				Schedule: &v1.ScheduleWorker{
					Key: f.cronWorkerOpts.Description,
					Cadence: &v1.ScheduleWorker_Cron{
						Cron: &v1.ScheduleCron{
							Cron: f.cronWorkerOpts.Cron,
						},
					},
				},
			}
		}
		if f.subscriptionWorkerOpts.Topic != "" {
			initRequest.Worker = &v1.InitRequest_Subscription{
				Subscription: &v1.SubscriptionWorker{
					Topic: f.subscriptionWorkerOpts.Topic,
				},
			}
		}
		if f.websocketWorkerOpts.Socket != "" {
			evtType, err := f.websocketWorkerOpts.eventTypeToWire()
			if err != nil {
				return err
			}

			initRequest.Worker = &v1.InitRequest_Websocket{
				Websocket: &v1.WebsocketWorker{
					Socket: f.websocketWorkerOpts.Socket,
					Event:  evtType,
				},
			}
		}

		// Let the membrane know the function is ready for initialization
		err := stream.Send(&v1.ClientMessage{
			Content: &v1.ClientMessage_InitRequest{
				InitRequest: initRequest,
			},
		})
		if err != nil {
			return err
		}

		errChan := make(chan error)

		// Start faasLoop in a go routine
		go faasLoop(stream, f, errChan)

		return <-errChan
	} else {
		return err
	}
}

// Creates a new HandlerBuilder
func New() HandlerBuilder {
	return &faasClientImpl{http: map[string]HttpMiddleware{}}
}

func (f *faasClientImpl) WithApiWorkerOpts(opts ApiWorkerOptions) HandlerBuilder {
	f.apiWorkerOpts = opts
	return f
}

func (f *faasClientImpl) WithRateWorkerOpts(opts RateWorkerOptions) HandlerBuilder {
	f.rateWorkerOpts = opts
	return f
}

func (f *faasClientImpl) WithCronWorkerOpts(opts CronWorkerOptions) HandlerBuilder {
	f.cronWorkerOpts = opts
	return f
}

func (f *faasClientImpl) WithSubscriptionWorkerOpts(opts SubscriptionWorkerOptions) HandlerBuilder {
	f.subscriptionWorkerOpts = opts
	return f
}

func (f *faasClientImpl) WithBucketNotificationWorkerOptions(opts BucketNotificationWorkerOptions) HandlerBuilder {
	f.bucketNotificationWorkerOpts = opts
	return f
}

func (f *faasClientImpl) WithWebsocketWorkerOptions(opts WebsocketWorkerOptions) HandlerBuilder {
	f.websocketWorkerOpts = opts
	return f
}
