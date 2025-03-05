package agent

type Agent struct {
	client  *OrchestratorClient
	workers int
}

func NewAgent(client *OrchestratorClient, workers int) *Agent {
	return &Agent{
		client:  client,
		workers: workers,
	}
}

func (a *Agent) Run() {
	for i := 0; i < a.workers; i++ {
		w := NewWorker(a.client)
		go w.Start()
	}
	select {}
}
