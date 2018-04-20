# portieris
Portieris is a Kubernetes Admission Controller that enforces image trust in your Kubernetes cluster. It hooks into the Kubernetes API server webhooks for a mutating and validating admission control. The service receives workload creation and update events and is able to verify provenance via a Notary server associated with the image registry.
