# docker-gs-ping Helm Chart

A simple Go server/microservice example for [Docker's Go Language Guide](https://docs.docker.com/language/golang/).

## Usage

[Helm](https://helm.sh) must be installed to use the charts. Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

    helm repo add docker-gs-ping https://unleftie.github.io/docker-gs-ping

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages. You can then run `helm search repo
{alias}` to see the charts.

To install the docker-gs-ping chart:

    helm install docker-gs-ping docker-gs-ping/docker-gs-ping

To uninstall the chart:

    helm delete docker-gs-ping

## License

[Apache-2.0 License](LICENSE)
