package worker_pools

func Load[request any, response any](function func(request) response, jobs []request, maxWorkers int) []response {
	jobsChan := make(chan request, len(jobs))
	resultsChan := make(chan response, len(jobs))

	for w := 1; w <= maxWorkers; w++ {
		go genericWorker(function, jobsChan, resultsChan)
	}
	for _, job := range jobs {
		jobsChan <- job
	}
	close(jobsChan)

	resp := make([]response, len(jobs))
	for a := 0; a < len(jobs); a++ {
		newRes := <-resultsChan
		resp[a] = newRes
	}
	close(resultsChan)

	return resp
}

func genericWorker[request any, response any](function func(request) response, jobs <-chan request, results chan<- response) {
	for j := range jobs {
		resp := function(j)

		results <- resp
	}
}
