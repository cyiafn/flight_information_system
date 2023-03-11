package worker_pools

// Load follows the worker pool pattern to concurrently process jobs with a max defined amount of workers
func Load[request any, response any](function func(request) response, jobs []request, maxWorkers int) []response {
	// jobsChan is a buffered channel of all jobs
	jobsChan := make(chan request, len(jobs))
	// resultsChan is a buffered channel of all results completed by the job
	resultsChan := make(chan response, len(jobs))

	for w := 1; w <= maxWorkers; w++ {
		// start each worker
		go genericWorker(function, jobsChan, resultsChan)
	}
	// for each job, send the job to the jobsChannel
	for _, job := range jobs {
		jobsChan <- job
	}
	// signal that no more jobs are coming
	close(jobsChan)

	// grab all responses
	resp := make([]response, len(jobs))
	for a := 0; a < len(jobs); a++ {
		newRes := <-resultsChan
		resp[a] = newRes
	}
	close(resultsChan)

	return resp
}

// genericWorker processes all jobs in the jobs channel and sends the output to results sequentially. It will automatically terminate when there are no more jobs
func genericWorker[request any, response any](function func(request) response, jobs <-chan request, results chan<- response) {
	for j := range jobs {
		resp := function(j)

		results <- resp
	}
}
