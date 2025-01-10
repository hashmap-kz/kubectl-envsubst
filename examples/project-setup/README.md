## A brief example of usage

### Prerequisites:

You can use a Kind cluster for testing. If you haven’t installed it yet, check out the [Kind installation guide](https://kind.sigs.k8s.io/).

### Deploy:

Execute scripts one by one:
```
# prepare kind cluster
bash 00-setup-kind.sh

# deploy for each environment
bash 01-deploy-dev.sh
bash 02-deploy-stage.sh
bash 03-deploy-prod.sh
```

### Verification:

Check the result:
- dev: http://localhost:32501
- stage: http://localhost:32502
- prod: http://localhost:32503

_If your Kind cluster is running on a remote machine, replace ‘localhost’ with the machine’s IP address._
