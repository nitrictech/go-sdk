package batch

type HandlerOption func(opts *handlerOptions)

// HandlerOptions defines the resource requirements for a job
type handlerOptions struct {
	// Cpus is the number of CPUs/vCPUs to allocate to the job
	cpus *float32
	// Memory is the amount of memory in MiB to allocate to the job
	memory *int64
	// Gpus is the number of GPUs to allocate to the job
	gpus *int64
}

// WithCpus - Set the number of CPUs/vCPUs to allocate to job handler instances
func WithCpus(cpus float32) HandlerOption {
	return func(opts *handlerOptions) {
		opts.cpus = &cpus
	}
}

// WithMemory - Set the amount of memory in MiB to allocate to job handler instances
func WithMemory(mib int64) HandlerOption {
	return func(opts *handlerOptions) {
		opts.memory = &mib
	}
}

// WithGpus - Set the number of GPUs to allocate to job handler instances
func WithGpus(gpus int64) HandlerOption {
	return func(opts *handlerOptions) {
		opts.gpus = &gpus
	}
}
