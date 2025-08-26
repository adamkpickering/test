These instructions start at:
https://www.elastic.co/docs/deploy-manage/deploy/cloud-on-k8s/install-using-yaml-manifest-quickstart

Deploy elasticsearch operator:
1. kubectl create -f elastic-operator-crds.yaml
2. kubectl apply -f elastic-operator.yaml

Deploy elasticsearch and test its API:
1. kubectl apply -f elasticsearch.yaml
2. Ensure `kubectl get elasticsearch` is green (>90% disk usage means failure)
3. let password = (kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
4. kubectl port-forward service/quickstart-es-http 9200
5. curl -u "elastic:$PASSWORD" -k "https://localhost:9200"

Deploy kibana:
1. kubectl apply -f kibana.yaml
2. Ensure `kubectl get kibana` is green
2. kubectl get secret quickstart-es-elastic-user -o=jsonpath='{.data.elastic}' | base64 --decode; echo
3. kubectl port-forward service/quickstart-kb-http 5601
4. In browser, go to localhost:5601 and log in with user "elastic" and password from above command.

Create log generator, Flow and Output:
1. helm upgrade --install --wait log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
2. kubectl apply -f output-flow.yaml

In kibana, sanity test:
1. Stack Management > Data Views
2. Create data view
3. Use "fluentd" as index pattern (should see it on right)
4. Analytics > Discover
5. Select the data view you just created in upper left corner; data should be there
