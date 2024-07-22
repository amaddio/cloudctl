## Required Configuration Files
As with the classic `kubectcl` the `cloudctl` expects you to have a `.kube/config` file in your home folder. 

## Development Dependencies 
```sh
go get -u github.com/spf13/cobra@latest
go get -u k8s.io/client-go@latest
go get -u k8s.io/apimachinery@latest
go get k8s.io/apimachinery/pkg/apis/meta/v1@v0.30.2
```