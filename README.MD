
## Required Configuration Files

As with the classic `kubectcl` the `cloudctl` expects you to have a `.kube/config` file in your home folder. 

## Create Test Kubernetes Cluster 

1. create test cluster here: https://labs.play-with-k8s.com/
1. create control-plane (a.k.a master) node in cluster: `kubeadm init --apiserver-advertise-address $(hostname -i) --pod-network-cidr 10.5.0.0/16`
1. create network with:
   2. `kubectl apply -f https://raw.githubusercontent.com/cloudnativelabs/kube-router/master/daemonset/kubeadm-kuberouter.yaml` or
   3. `kubectl apply -n kube-system -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"`
4. ... wait until network is ready. Check with: `kubectl get nodes`
5. Create a worker node: `kubeadm join 192.168.0.13:6443 --token lsixxy.mufm4mv64xi4tuf7 --discovery-token-ca-cert-hash sha256:3f6c6aeebe8fe347dfd2df7e45963cbaddca1d84bfc492bf1690879c1accd6a5`
   
       sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
      sudo chown $(id -u):$(id -g) $HOME/.kube/config
- create network for services: `kubectl apply -n kube-system -f "https://cloud.weave.works/k8s/net?
  k8s-version=$(kubectl version | base64 | tr -d '\n')"`
- 
- Initialize cluster networking: `kubectl apply -f https://raw.githubusercontent.com/cloudnativelabs/kube-router/master/daemonset/kubeadm-kuberouter.yaml`


kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/master/content/en/examples/application/nginx-app.yaml

## Dependencies 

```sh
go get -u github.com/spf13/cobra@latest
go get -u k8s.io/client-go@latest
go get -u k8s.io/apimachinery@latest
go get k8s.io/apimachinery/pkg/apis/meta/v1@v0.30.2
```