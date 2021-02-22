package faas_test

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.nitric.io/go-sdk/v1/faas"
)

type MockFunctionSpy struct {
	loggedRequests []*faas.NitricRequest
	mockResponse   *faas.NitricResponse
}

func (m *MockFunctionSpy) reset() {
	m.loggedRequests = make([]*faas.NitricRequest, 0)
}

func (m *MockFunctionSpy) handler(r *faas.NitricRequest) *faas.NitricResponse {
	if m.loggedRequests == nil {
		m.loggedRequests = make([]*faas.NitricRequest, 0)
	}

	m.loggedRequests = append(m.loggedRequests, r)

	return m.mockResponse
}

type MockHttpOptions struct {
	requestId   string
	sourceType  string
	source      string
	payloadType string
	body        []byte
}

func createMockHttpRequest(opts MockHttpOptions) *http.Request {
	request, _ := http.NewRequest("POST", "http://0.0.0.0:9001", bytes.NewReader(opts.body))

	request.Header.Add("x-nitric-request-id", opts.requestId)
	request.Header.Add("x-nitric-source-type", opts.sourceType)
	request.Header.Add("x-nitric-source", opts.source)
	request.Header.Add("x-nitric-payloadTyp", opts.payloadType)

	return request
}

var _ = Describe("Faas", func() {
	Context("Start", func() {
		mockFunction := &MockFunctionSpy{
			mockResponse: &faas.NitricResponse{
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Status: 200,
				Body:   []byte("Hello"),
			},
		}

		BeforeEach(func() {
			mockFunction.reset()
		})

		go (func() {
			faas.Start(mockFunction.handler)
		})()

		When("Function is called with a Request payload", func() {
			BeforeEach(func() {
				request := createMockHttpRequest(MockHttpOptions{
					requestId:   "1234",
					sourceType:  "REQUEST",
					source:      "test-source",
					payloadType: "test-payload",
					body:        []byte("Test"),
				})

				http.DefaultClient.Do(request)
			})

			It("Should recieve the correct request", func() {
				By("Recieving a single request")
				Expect(mockFunction.loggedRequests).To(HaveLen(1))

				recievedRequest := mockFunction.loggedRequests[0]
				recievedContext := recievedRequest.GetContext()

				By("Having the provided request id")
				Expect((&recievedContext).GetRequestID()).To(BeEquivalentTo("1234"))

				//By("Having the provided payload type")
				//Expect((&recievedContext).GetPayloadType()).To(BeEquivalentTo("test-payload"))

				By("Having the correct source type")
				Expect((&recievedContext).GetSourceType()).To(BeEquivalentTo(faas.Request))

				By("Having the provided source")
				Expect((&recievedContext).GetSource()).To(BeEquivalentTo("test-source"))
			})
		})

		When("The Function is called with a Subscription payload", func() {
			BeforeEach(func() {
				request := createMockHttpRequest(MockHttpOptions{
					requestId:   "1234",
					sourceType:  "SUBSCRIPTION",
					source:      "test-source",
					payloadType: "test-payload",
					body:        []byte("Test"),
				})

				http.DefaultClient.Do(request)
			})

			It("Should have the supplied sourceType", func() {
				By("Recieving a single request")
				Expect(mockFunction.loggedRequests).To(HaveLen(1))

				recievedRequest := mockFunction.loggedRequests[0]
				recievedContext := recievedRequest.GetContext()

				Expect((&recievedContext).GetSourceType()).To(BeEquivalentTo(faas.Subscription))
			})
		})

		When("The Function is called with an unknown source type", func() {
			BeforeEach(func() {
				request := createMockHttpRequest(MockHttpOptions{
					requestId:   "1234",
					sourceType:  "fake-source",
					source:      "test-source",
					payloadType: "test-payload",
					body:        []byte("Test"),
				})

				http.DefaultClient.Do(request)
			})

			It("Should have the supplied sourceType", func() {
				By("Recieving a single request")
				Expect(mockFunction.loggedRequests).To(HaveLen(1))

				recievedRequest := mockFunction.loggedRequests[0]
				recievedContext := recievedRequest.GetContext()

				Expect((&recievedContext).GetSourceType()).To(BeEquivalentTo(faas.Unknown))
			})
		})
	})
})
