# Provider ClickHouse <img src="icon.svg" alt="Provider ClickHouse Icon" style="height: 1em; vertical-align: middle; margin-left: 1em">

`provider-clickhouse` is a [Crossplane](https://crossplane.io/) provider for
[ClickHouse Cloud](https://clickhouse.com/cloud) built with
[Upjet](https://github.com/crossplane/upjet) on top of the
[ClickHouse/clickhouse](https://registry.terraform.io/providers/ClickHouse/clickhouse/latest)
Terraform provider. It exposes XRM-conformant managed resources for the
ClickHouse Cloud API (services, private endpoints, transparent data encryption,
organization settings, ...).

## Getting Started

### Prerequisites

- A Kubernetes cluster with [Crossplane](https://crossplane.io/) installed
- A ClickHouse Cloud organization with an API key (`token_key` + `token_secret`)

### Authentication

The provider authenticates against the ClickHouse Cloud API using an API key
stored in a Kubernetes Secret. The secret contains a JSON object with the
organization ID and the key pair.

#### 1. Create a Secret with ClickHouse Cloud credentials

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: clickhouse-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "organization_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
      "token_key": "my-token-key",
      "token_secret": "my-token-secret"
    }
```

Optional: override the API endpoint by adding `"api_url": "https://api.clickhouse.cloud/v1"`.

#### 2. Create a ProviderConfig referencing the Secret

**Namespaced** (secret must be in the same namespace as the managed resources):

```yaml
apiVersion: clickhouse.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: crossplane-system
spec:
  credentials:
    source: Secret
    secretRef:
      name: clickhouse-creds
      key: credentials
```

**Cluster-scoped** (for cluster-wide access):

```yaml
apiVersion: clickhouse.m.crossplane.io/v1beta1
kind: ClusterProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: clickhouse-creds
      namespace: crossplane-system
      key: credentials
```

The `credentials.source` field supports: `Secret`, `InjectedIdentity`, `Environment`, `Filesystem`, and `None`.

See example manifests under [`examples/namespaced/`](examples/namespaced/) and [`examples/cluster/`](examples/cluster/).

## Managed Resources

The provider exposes resources across three API groups:

- `clickhouse.clickhouse.m.crossplane.io` — `Service`
- `service.clickhouse.m.crossplane.io` — `PrivateEndpointsAttachment`, `TransparentDataEncryptionKeyAssociation`
- `organization.clickhouse.m.crossplane.io` — `Settings`

## Developing

Run code-generation pipeline:

```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster (out of cluster):

```console
make run
```

or (deploying in-cluster):

```console
make local-deploy
```

Run e2e tests (locally in a KinD cluster):

```console
# UPTEST_SKIP_DELETE=true
make e2e
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/lansweeper-oss/provider-clickhouse/issues).

## Licensing

`provider-clickhouse` is under the Apache 2.0 [license](LICENSE).
