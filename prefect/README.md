# Running Prefect Worker Locally

Install Deps

```zsh
cd prefect

pip install -r requirements.txt
```

To start the agent locally

```zsh
prefect worker start --pool justin-local
```

To initialize Prefect

```zsh
prefect init
```

To deploy a workflow:

```zsh
prefect deploy
```


To trigger with API call

```zsh
POST https://api.prefect.cloud/api/accounts/<account-id>/workspaces/<workspace-id>/deployments/<deployment-id>/create_flow_run

Body:
    {
        "parameters": {
            "target": "Jeff Bezos"
	    }
    }
```
