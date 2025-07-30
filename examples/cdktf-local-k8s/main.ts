import { Construct } from "constructs";
import { App, TerraformStack } from "cdktf";
import { KubernetesProvider } from "@cdktf/provider-kubernetes/lib/provider";
import { Deployment } from "@cdktf/provider-kubernetes/lib/deployment";
import { Service } from "@cdktf/provider-kubernetes/lib/service";

class MyStack extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);

    new KubernetesProvider(this, "kubernetes", {});

    const appName = "nginx-app";
    const appLabels = { app: appName };

    new Deployment(this, "nginx-deployment", {
      metadata: {
        name: appName,
      },
      spec: {
        replicas: "2",
        selector: {
          matchLabels: appLabels,
        },
        template: {
          metadata: {
            labels: appLabels,
          },
          spec: {
            container: [
              {
                name: appName,
                image: "nginx:1.19.10",
                port: [{ containerPort: 80 }],
              },
            ],
          },
        },
      },
    });

    new Service(this, "nginx-service", {
      metadata: {
        name: appName,
      },
      spec: {
        selector: appLabels,
        port: [{ port: 80, targetPort: "80" }],
        type: "LoadBalancer",
      },
    });
  }
}

const app = new App();
new MyStack(app, "cdktf-local-k8s");
app.synth();
