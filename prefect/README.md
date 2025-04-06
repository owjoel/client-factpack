# Running Prefect Worker Locally

Install Deps

```zsh
pip install -r requirements.txt
```


To initialize Prefect

```zsh
prefect init
```

To deploy a workflow:

```zsh
prefect deploy
```

To start the agent locally

```zsh
prefect worker start --pool <worker-pool-name>
```

To trigger with API call

```
POST https://api.prefect.cloud/api/accounts/<account-id>/workspaces/<workspace-id>/deployments/<deployment-id>/create_flow_run

Body:
    {
        "parameters": {
            "target": "Jeff Bezos"
	    }
    }
```
