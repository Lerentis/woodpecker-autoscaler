package models

/*
	{
	  "pending": [
	    {
	      "id": "146",
	      "data": "REDACTED",
	      "labels": {
	        "repo": "REDACTED",
	        "type": "picus"
	      },
	      "dependencies": null,
	      "run_on": null,
	      "dep_status": {},
	      "agent_id": 0
	    }
	  ],
	  "waiting_on_deps": null,
	  "running": null,
	  "stats": {
	    "worker_count": 40,
	    "pending_count": 1,
	    "waiting_on_deps_count": 0,
	    "running_count": 0,
	    "completed_count": 0
	  },
	  "paused": false
	}
*/

type PendingInformation struct {
	ID           int               `json:"id"`
	Data         string            `json:"data"`
	Labels       map[string]string `json:"labels"`
	Dependencies string            `json:"dependencies"`
	RunOn        string            `json:"run_on"`
	DepStatus    string            `json:"-"` // dont need those
	AgentId      int               `json:"agent_id"`
}

type Stats struct {
	WorkerCount        int `json:"worker_count"`
	PendingCount       int `json:"pending_count"`
	WaitingOnDepsCount int `json:"waiting_on_deps_count"`
	RunningCount       int `json:"running_count"`
	CompletedCount     int `json:"completed_count"`
}

type QueueInfo struct {
	Pending       []PendingInformation `json:"pending"`
	WaitingOnDeps string               `json:"-"` // dont need those
	Running       int                  `json:"running"`
	Stats         Stats                `json:"stats"`
	Paused        bool                 `json:"paused"`
}
