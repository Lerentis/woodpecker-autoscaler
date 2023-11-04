# woodpecker-autoscaler

[![status-badge](https://woodpecker.uploadfilter24.eu/api/badges/8/status.svg)](https://woodpecker.uploadfilter24.eu/repos/8)

Dynamically spawns woodpecker ci build agents in hetzner cloud.

## Installing with Helm

The deployment will use helm and the chart in `chart/woodpecker-autoscaler`.  
You will need a [hcloud api token](https://docs.hetzner.com/cloud/api/getting-started/generating-api-token/), a [woodpecker agent secret](https://woodpecker-ci.org/docs/administration/agent-config#woodpecker_agent_secret), a woodpecker api token and a definition of the build agents to be created.
Expose these information to the floater as described in this example:

```yaml
env:
  - name: WOODPECKER_AUTOSCALER_LOGLEVEL
    value: "Info"
  - name: WOODPECKER_AUTOSCALER_CHECK_INTERVAL
    value: "15"
  - name: WOODPECKER_AUTOSCALER_WOODPECKER_LABEL_SELECTOR
    value: "uploadfilter24.eu/instance-role=Woodpecker"
  - name: WOODPECKER_AUTOSCALER_WOODPECKER_INSTANCE
    value: "define_it"
  - name: WOODPECKER_AUTOSCALER_WOODPECKER_AGENT_SECRET
    value: "define_it"
  - name: WOODPECKER_AUTOSCALER_WOODPECKER_API_TOKEN
    value: "define_it"
  - name: WOODPECKER_AUTOSCALER_HCLOUD_TOKEN
    value: "define_it"
  - name: WOODPECKER_AUTOSCALER_HCLOUD_INSTANCE_TYPE
    value: "cpx21"
  - name: WOODPECKER_AUTOSCALER_HCLOUD_REGION
    value: "define_it"
  - name: WOODPECKER_AUTOSCALER_HCLOUD_DATACENTER
    value: "define_it"
  - name: WOODPECKER_AUTOSCALER_HCLOUD_SSH_KEY
    value: "define_it"
```

you can also create a secret manually with these information and reference the existing secret like this in the `values.yaml`:

```yaml
externalConfigSecret:
  enabled: true
  name: "my-existing-secret"
```

Now you are able to deploy:

```bash
kubectl create namespace woodpecker-autoscaler
cd chart/woodpecker-autoscaler
helm upgrade --install -f values.yaml -n woodpecker-autoscaler woodpecker-autoscaler ./
```

## Installing Manually

Download the binary from the release section and place it somewhere on your system; `/usr/bin/woodpecker-autoscaler` for example.
Create a systemd service like the one in this template:

```systemd
[Unit]
Description=Dynamically spawn woodpecker ci build agents in hetzner cloud

[Service]
Type=simple
Nice=10
ExecStart=/usr/bin/woodpecker-autoscaler
EnvironmentFile="/etc/default/woodpecker-autoscaler.env"
```

Now place the environment variable configuration in the specified file:

```bash
WOODPECKER_AUTOSCALER_LOGLEVEL=Info
WOODPECKER_AUTOSCALER_CHECK_INTERVAL=15
WOODPECKER_AUTOSCALER_WOODPECKER_LABEL_SELECTOR="uploadfilter24.eu/instance-role=Woodpecker"
WOODPECKER_AUTOSCALER_WOODPECKER_INSTANCE="define_it"
WOODPECKER_AUTOSCALER_WOODPECKER_AGENT_SECRET="define_it"
WOODPECKER_AUTOSCALER_WOODPECKER_API_TOKEN="define_it"
WOODPECKER_AUTOSCALER_HCLOUD_TOKEN="define_it"
WOODPECKER_AUTOSCALER_HCLOUD_INSTANCE_TYPE=cpx21
WOODPECKER_AUTOSCALER_HCLOUD_REGION="define_it"
WOODPECKER_AUTOSCALER_HCLOUD_DATACENTER="define_it"
WOODPECKER_AUTOSCALER_HCLOUD_SSH_KEY="define_it"
```

Now reload the systemd daemons and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable woodpecker-autoscaler
sudo systemctl start woodpecker-autoscaler
```

> Made with ♡ by the folkes at uploadfilter24.eu :)