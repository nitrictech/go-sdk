package constants

import "google.golang.org/grpc"

// DefaultOptions - Provides option defaults for creating a gRPC service connection with the Nitric Membrane
func DefaultOptions() []grpc.DialOption {
	return []grpc.DialOption{
		// TODO: Look at authentication config with membrane
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(NitricDialTimeout()),
	}
}
