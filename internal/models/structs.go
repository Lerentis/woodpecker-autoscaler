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

type JobInformation struct {
	ID           int               `json:"id"`
	Data         string            `json:"data"`
	Labels       map[string]string `json:"labels"`
	Dependencies string            `json:"dependencies,omitempty"`
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
	Pending       []JobInformation `json:"pending,omitempty"`
	WaitingOnDeps string           `json:"-"` // dont need those
	Running       []JobInformation `json:"running,omitempty"`
	Stats         Stats            `json:"stats"`
	Paused        bool             `json:"paused"`
}

/*[
  {
    "id": 2,
    "created": 1693567407,
    "updated": 1694013270,
    "name": "",
    "owner_id": -1,
    "token": "redacted",
    "last_contact": 1694013270,
    "platform": "linux/arm64",
    "backend": "kubernetes",
    "capacity": 4,
    "version": "next-971534929c",
    "no_schedule": false
  }
]*/

type Agent struct {
	ID          int64  `json:"id"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
	Name        string `json:"name"`
	OwnerID     int64  `json:"owner_id"`
	Token       string `json:"token"`
	LastContact int64  `json:"last_contact"`
	Platform    string `json:"platform"`
	Backend     string `json:"backend"`
	Capacity    int32  `json:"capacity"`
	Version     string `json:"version"`
	NoSchedule  bool   `json:"no_schedule"`
}

type AgentList struct {
	Agents []Agent
}
