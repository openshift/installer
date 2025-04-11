# -*- mode: Python -*-

# Pre-requisite make targets "install-tools" and "kind-create" ensure that the below tools are already installed.
envsubst_cmd = "./hack/tools/bin/envsubst"
kubectl_cmd = "./hack/tools/bin/kubectl"
helm_cmd = "./hack/tools/bin/helm"
kind_cmd = "./hack/tools/bin/kind"
tools_bin = "./hack/tools/bin"

#Add tools to path
os.putenv("PATH", os.getenv("PATH") + ":" + tools_bin)

update_settings(k8s_upsert_timeout_secs = 60)  # on first tilt up, often can take longer than 30 seconds

# Default settings for tilt
settings = {
    "allowed_contexts": [
        "kind-capz",
    ],
    "deploy_cert_manager": True,
    "preload_images_for_kind": True,
    "kind_cluster_name": "capz",
    "capi_version": "v1.8.10",
    "caaph_version": "v0.2.5",
    "cert_manager_version": "v1.16.1",
    "kubernetes_version": "v1.28.3",
    "aks_kubernetes_version": "v1.28.3",
    "flatcar_version": "3374.2.1",
    "azure_location": "eastus",
    "control_plane_machine_count": "1",
    "az_control_plane_machine_type": "Standard_B2s",
    "worker_machine_count": "2",
    "az_node_machine_type": "Standard_B2s",
    "cluster_class_name": "default",
}

# Auth keys that need to be loaded from the environment
keys = ["AZURE_SUBSCRIPTION_ID", "AZURE_TENANT_ID", "AZURE_CLIENT_ID"]

# Get global settings from tilt-settings.yaml or tilt-settings.json
tilt_file = "./tilt-settings.yaml" if os.path.exists("./tilt-settings.yaml") else "./tilt-settings.json"
settings.update(read_yaml(tilt_file, default = {}))

if settings.get("trigger_mode") == "manual":
    trigger_mode(TRIGGER_MODE_MANUAL)

if "allowed_contexts" in settings:
    allow_k8s_contexts(settings.get("allowed_contexts"))

if "default_registry" in settings:
    default_registry(settings.get("default_registry"))

os_arch = str(local("go env GOARCH")).rstrip("\n")
if "aks" in settings.get("kustomize_substitutions", {}).get("MGMT_CLUSTER_NAME", ""):
    print("Using AKS as management cluster, setting os_arch to amd64")
    os_arch = "amd64"

# deploy CAPI
def deploy_capi():
    version = settings.get("capi_version")
    capi_uri = "https://github.com/kubernetes-sigs/cluster-api/releases/download/{}/cluster-api-components.yaml".format(version)
    cmd = "curl --retry 3 -sSL {} | {} | {} apply -f -".format(capi_uri, envsubst_cmd, kubectl_cmd)
    local(cmd, quiet = True)
    if settings.get("extra_args"):
        extra_args = settings.get("extra_args")
        if extra_args.get("core"):
            core_extra_args = extra_args.get("core")
            for namespace in ["capi-system", "capi-webhook-system"]:
                patch_args_with_extra_args(namespace, "capi-controller-manager", core_extra_args)
        if extra_args.get("kubeadm-bootstrap"):
            kb_extra_args = extra_args.get("kubeadm-bootstrap")
            patch_args_with_extra_args("capi-kubeadm-bootstrap-system", "capi-kubeadm-bootstrap-controller-manager", kb_extra_args)

# deploy CAAPH
def deploy_caaph():
    version = settings.get("caaph_version")

    caaph_uri = "https://github.com/kubernetes-sigs/cluster-api-addon-provider-helm/releases/download/{}/addon-components.yaml".format(version)
    cmd = "curl --retry 3 -sSL {} | {} | {} apply -f -".format(caaph_uri, envsubst_cmd, kubectl_cmd)
    local(cmd, quiet = True)
    if settings.get("extra_args"):
        extra_args = settings.get("extra_args")
        if extra_args.get("helm"):
            core_extra_args = extra_args.get("helm")
            for namespace in ["caaph-system", "caaph-webhook-system"]:
                patch_args_with_extra_args(namespace, "caaph-controller-manager", core_extra_args)

def patch_args_with_extra_args(namespace, name, extra_args):
    args_str = str(local("{} get deployments {} -n {} -o jsonpath={{.spec.template.spec.containers[1].args}}".format(kubectl_cmd, name, namespace)))
    args_to_add = [arg for arg in extra_args if arg not in args_str]
    if args_to_add:
        args = args_str[1:-1].split()
        args.extend(args_to_add)
        patch = [{
            "op": "replace",
            "path": "/spec/template/spec/containers/1/args",
            "value": args,
        }]
        local("{} patch deployment {} -n {} --type json -p='{}'".format(kubectl_cmd, name, namespace, str(encode_json(patch)).replace("\n", "")))

# Users may define their own Tilt customizations in tilt.d. This directory is excluded from git and these files will
# not be checked in to version control.
def include_user_tilt_files():
    user_tiltfiles = listdir("tilt.d")
    for f in user_tiltfiles:
        include(f)

def append_arg_for_container_in_deployment(yaml_stream, name, namespace, contains_image_name, args):
    for item in yaml_stream:
        if item["kind"] == "Deployment" and item.get("metadata").get("name") == name and item.get("metadata").get("namespace") == namespace:
            containers = item.get("spec").get("template").get("spec").get("containers")
            for container in containers:
                if contains_image_name in container.get("image"):
                    container.get("args").extend(args)

def fixup_yaml_empty_arrays(yaml_str):
    yaml_str = yaml_str.replace("conditions: null", "conditions: []")
    return yaml_str.replace("storedVersions: null", "storedVersions: []")

def validate_auth():
    substitutions = settings.get("kustomize_substitutions", {})
    os.environ.update(substitutions)
    for sub in substitutions:
        if sub[-4:] == "_B64":
            os.environ[sub[:-4]] = base64_decode(os.environ[sub])
    missing = [k for k in keys if not os.environ.get(k)]
    if missing:
        fail("missing kustomize_substitutions keys {} in tilt-setting.json".format(missing))

tilt_helper_dockerfile_header = """
# Tilt image
FROM golang:1.22 as tilt-helper
# Support live reloading with Tilt
RUN wget --output-document /restart.sh --quiet https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/restart.sh  && \
    wget --output-document /start.sh --quiet https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/start.sh && \
    chmod +x /start.sh && chmod +x /restart.sh && \
    touch /process.txt && chmod 0777 /process.txt `# pre-create PID file to allow even non-root users to run the image`
"""

tilt_dockerfile_header = """
FROM gcr.io/distroless/base:debug as tilt
WORKDIR /tilt
RUN ["/busybox/chmod", "0777", "."]
COPY --from=tilt-helper /process.txt .
COPY --from=tilt-helper /start.sh .
COPY --from=tilt-helper /restart.sh .
COPY manager .
"""

# Install the OpenTelemetry helm chart
def observability():
    instrumentation_key = os.getenv("AZURE_INSTRUMENTATION_KEY", "")
    if instrumentation_key == "":
        warn("AZURE_INSTRUMENTATION_KEY is not set, so traces won't be exported to Application Insights")
        trace_links = []
    else:
        trace_links = [link("https://ms.portal.azure.com/#blade/HubsExtension/BrowseResource/resourceType/microsoft.insights%2Fcomponents", "App Insights")]
    k8s_yaml(helm(
        "./hack/observability/opentelemetry/chart",
        name = "opentelemetry-collector",
        namespace = "capz-system",
        values = ["./hack/observability/opentelemetry/values.yaml"],
        set = ["config.exporters.azuremonitor.instrumentation_key=" + instrumentation_key],
    ))
    k8s_yaml(helm(
        "./hack/observability/jaeger/chart",
        name = "jaeger-all-in-one",
        namespace = "capz-system",
        set = [
            "crd.install=false",
            "rbac.create=false",
            "resources.limits.cpu=200m",
            "resources.limits.memory=256Mi",
        ],
    ))

    k8s_yaml(helm(
        "./hack/observability/cluster-api-visualizer/chart",
        name = "visualizer",
        namespace = "capz-system",
    ))

    k8s_resource(
        workload = "jaeger-all-in-one",
        new_name = "traces: jaeger-all-in-one",
        port_forwards = [port_forward(16686, name = "View traces", link_path = "/search?service=capz")],
        links = trace_links,
        labels = ["observability"],
    )
    k8s_resource(
        workload = "prometheus-operator",
        new_name = "metrics: prometheus-operator",
        port_forwards = [port_forward(local_port = 9090, container_port = 9090, name = "View metrics")],
        extra_pod_selectors = [{"app": "prometheus"}],
        labels = ["observability"],
    )
    k8s_resource(workload = "opentelemetry-collector", labels = ["observability"])
    k8s_resource(workload = "opentelemetry-collector-agent", labels = ["observability"])
    k8s_resource(
        workload = "capi-visualizer",
        new_name = "visualizer",
        port_forwards = [port_forward(local_port = 8000, container_port = 8081, name = "View visualization")],
        labels = ["observability"],
    )

    k8s_resource(workload = "capz-controller-manager", labels = ["cluster-api"])
    k8s_resource(workload = "azureserviceoperator-controller-manager", labels = ["cluster-api"])

# Build CAPZ and add feature gates
def capz():
    # Apply the kustomized yaml for this provider
    yaml = str(kustomizesub("./hack/observability"))  # build an observable kind deployment by default

    # add extra_args if they are defined
    if settings.get("extra_args"):
        azure_extra_args = settings.get("extra_args").get("azure")
        yaml_dict = decode_yaml_stream(yaml)
        append_arg_for_container_in_deployment(yaml_dict, "capz-controller-manager", "capz-system", "cluster-api-azure-controller", azure_extra_args)
        yaml = str(encode_yaml_stream(yaml_dict))
        yaml = fixup_yaml_empty_arrays(yaml)

    # Forge the build command
    ldflags = "-extldflags \"-static\" " + str(local("hack/version.sh")).rstrip("\n")
    build_env = "CGO_ENABLED=0 GOOS=linux GOARCH={arch}".format(arch = os_arch)
    build_cmd = "{build_env} go build -ldflags '{ldflags}' -o .tiltbuild/manager".format(
        build_env = build_env,
        ldflags = ldflags,
    )

    # Set up a local_resource build of the provider's manager binary.
    local_resource(
        "manager",
        cmd = "mkdir -p .tiltbuild; " + build_cmd,
        deps = ["api", "azure", "config", "controllers", "exp", "feature", "pkg", "util", "go.mod", "go.sum", "main.go"],
        labels = ["cluster-api"],
    )

    dockerfile_contents = "\n".join([
        tilt_helper_dockerfile_header,
        tilt_dockerfile_header,
    ])

    entrypoint = ["sh", "/tilt/start.sh", "/tilt/manager"]
    extra_args = settings.get("extra_args")
    if extra_args:
        entrypoint.extend(extra_args)

    # use the user REGISTRY if set, otherwise use the default
    if settings.get("kustomize_substitutions", {}).get("REGISTRY", "") != "":
        registry = settings.get("kustomize_substitutions", {}).get("REGISTRY", "")
        print("Using REGISTRY: " + registry + " from tilt-settings.yaml")
        image = registry + "/cluster-api-azure-controller"
    else:
        image = "gcr.io/cluster-api-provider-azure/cluster-api-azure-controller"

    # Set up an image build for the provider. The live update configuration syncs the output from the local_resource
    # build into the container.
    docker_build(
        ref = image,
        context = "./.tiltbuild/",
        dockerfile_contents = dockerfile_contents,
        target = "tilt",
        entrypoint = entrypoint,
        only = "manager",
        live_update = [
            sync(".tiltbuild/manager", "/tilt/manager"),
            run("sh /tilt/restart.sh"),
        ],
        ignore = ["templates"],
    )

    k8s_yaml(blob(yaml))

def create_identity_secret():
    #create secret for identity password
    local(kubectl_cmd + " delete secret cluster-identity-secret --ignore-not-found=true")

    os.putenv("AZURE_CLUSTER_IDENTITY_SECRET_NAME", "cluster-identity-secret")
    os.putenv("AZURE_CLUSTER_IDENTITY_SECRET_NAMESPACE", "default")
    os.putenv("CLUSTER_IDENTITY_NAME", "cluster-identity-ci")
    os.putenv("ASO_CREDENTIAL_SECRET_NAME", "aso-credentials")

    local("cat templates/flavors/aks-aso/credentials.yaml | " + envsubst_cmd + " | " + kubectl_cmd + " apply -f -", quiet = True, echo_off = True)

def create_crs():
    # create config maps
    local(kubectl_cmd + " delete configmaps csi-proxy-addon --ignore-not-found=true")
    local(kubectl_cmd + " create configmap csi-proxy-addon --from-file=templates/addons/windows/csi-proxy/csi-proxy.yaml")

    # need to set version for kube-proxy on windows.
    os.putenv("KUBERNETES_VERSION", settings.get("kubernetes_version", {}))
    local(kubectl_cmd + " create configmap calico-windows-addon --from-file=templates/addons/windows/calico/ --dry-run=client -o yaml | " + envsubst_cmd + " | " + kubectl_cmd + " apply -f -")

    # set up crs
    local(kubectl_cmd + " apply -f templates/addons/windows/calico-resource-set.yaml")
    local(kubectl_cmd + " apply -f templates/addons/windows/csi-proxy/csi-proxy-resource-set.yaml")

# create flavor resources from cluster-template files in the templates directory
def flavors():
    substitutions = settings.get("kustomize_substitutions", {})

    az_key_b64_name = "AZURE_SSH_PUBLIC_KEY_B64"
    az_key_name = "AZURE_SSH_PUBLIC_KEY"
    default_key_path = "$HOME/.ssh/id_rsa.pub"

    if substitutions.get(az_key_b64_name):
        os.environ.update({az_key_b64_name: substitutions.get(az_key_b64_name)})
        os.environ.update({az_key_name: base64_decode(substitutions.get(az_key_b64_name))})
    else:
        print("{} was not specified in tilt-settings.json, attempting to load {}".format(az_key_b64_name, default_key_path))
        os.environ.update({az_key_b64_name: base64_encode_file(default_key_path)})
        os.environ.update({az_key_name: read_file_from_path(default_key_path)})

    template_list = [item for item in listdir("./templates")]
    template_list = [template for template in template_list if os.path.basename(template).endswith("yaml")]

    for template in template_list:
        deploy_worker_templates(template, substitutions)

    local_resource(
        name = "delete-all-workload-clusters",
        cmd = kubectl_cmd + " delete clusters --all --wait=false",
        auto_init = False,
        trigger_mode = TRIGGER_MODE_MANUAL,
        labels = ["flavors"],
    )

def deploy_worker_templates(template, substitutions):
    # validate template exists
    if not os.path.exists(template):
        fail(template + " not found")

    yaml = str(read_file(template))
    flavor = os.path.basename(template).replace("cluster-template-", "").replace(".yaml", "")

    # for the base cluster-template, flavor is "default"
    flavor = os.path.basename(flavor).replace("cluster-template", "default")

    # azure account and ssh replacements
    for substitution in substitutions:
        value = substitutions[substitution]
        yaml = yaml.replace("${" + substitution + "}", value)

    # if metadata defined for worker-templates in tilt_settings
    if "worker-templates" in settings:
        # first priority replacements defined per template
        if "flavors" in settings.get("worker-templates", {}):
            substitutions = settings.get("worker-templates").get("flavors").get(flavor, {})
            for substitution in substitutions:
                value = substitutions[substitution]
                yaml = yaml.replace("${" + substitution + "}", value)

        # second priority replacements defined common to templates
        if "metadata" in settings.get("worker-templates", {}):
            substitutions = settings.get("worker-templates").get("metadata", {})
            for substitution in substitutions:
                value = substitutions[substitution]
                yaml = yaml.replace("${" + substitution + "}", value)

    # programmatically define any remaining vars
    # "windows" can not be for cluster name because it sets the dns to trademarked name during reconciliation
    substitutions = {
        "AZURE_LOCATION": settings.get("azure_location"),
        "AZURE_VNET_NAME": "${CLUSTER_NAME}-vnet",
        "AZURE_RESOURCE_GROUP": "${CLUSTER_NAME}-rg",
        "CONTROL_PLANE_MACHINE_COUNT": settings.get("control_plane_machine_count"),
        "KUBERNETES_VERSION": settings.get("kubernetes_version"),
        "AZURE_CONTROL_PLANE_MACHINE_TYPE": settings.get("az_control_plane_machine_type"),
        "WORKER_MACHINE_COUNT": settings.get("worker_machine_count"),
        "AZURE_NODE_MACHINE_TYPE": settings.get("az_node_machine_type"),
        "FLATCAR_VERSION": settings.get("flatcar_version"),
        "CLUSTER_CLASS_NAME": settings.get("cluster_class_name"),
    }

    if "aks" in flavor:
        # AKS version support is usually a bit behind CAPI version, so use an older version
        substitutions["KUBERNETES_VERSION"] = settings.get("aks_kubernetes_version")

    for substitution in substitutions:
        value = substitutions[substitution]
        yaml = yaml.replace("${" + substitution + "}", value)

    yaml = shlex.quote(yaml)
    flavor_name = os.path.basename(flavor)
    flavor_cmd = "RANDOM=$(bash -c 'echo $RANDOM'); export CLUSTER_NAME=" + flavor.replace("windows", "win") + "-$RANDOM; make generate-flavors; echo " + yaml + "> ./.tiltbuild/" + flavor + "; cat ./.tiltbuild/" + flavor + " | " + envsubst_cmd + " | " + kubectl_cmd + " apply -f -; echo \"Cluster \'$CLUSTER_NAME\' created, don't forget to delete\""

    # wait for kubeconfig to be available
    flavor_cmd += "; until " + kubectl_cmd + " get secret ${CLUSTER_NAME}-kubeconfig > /dev/null 2>&1; do sleep 5; done; " + kubectl_cmd + " get secret ${CLUSTER_NAME}-kubeconfig -o jsonpath={.data.value} | base64 --decode > ./${CLUSTER_NAME}.kubeconfig; chmod 600 ./${CLUSTER_NAME}.kubeconfig; until " + kubectl_cmd + " --kubeconfig=./${CLUSTER_NAME}.kubeconfig get nodes > /dev/null 2>&1; do sleep 5; done"

    # copy the kubeadm configmap to the calico-system namespace.
    # This is a workaround needed for the calico-node-windows daemonset to be able to run in the calico-system namespace.
    if "windows" in flavor_name:
        flavor_cmd += "; until " + kubectl_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig get configmap kubeadm-config --namespace=kube-system > /dev/null 2>&1; do sleep 5; done"
        flavor_cmd += "; " + kubectl_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig create namespace calico-system --dry-run=client -o yaml | " + kubectl_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig apply -f -; " + kubectl_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig get configmap kubeadm-config --namespace=kube-system -o yaml | sed 's/namespace: kube-system/namespace: calico-system/' | " + kubectl_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig apply -f -"

    flavor_cmd += get_addons(flavor_name)

    local_resource(
        name = flavor_name,
        cmd = ["sh", "-ec", flavor_cmd],
        auto_init = False,
        trigger_mode = TRIGGER_MODE_MANUAL,
        labels = ["flavors"],
        allow_parallel = True,
    )

def get_addons(flavor_name):
    # do not install calico and out of tree cloud provider for aks workload cluster
    if "aks" in flavor_name:
        return ""

    addon_cmd = "; export CIDRS=$(" + kubectl_cmd + " get cluster ${CLUSTER_NAME} -o jsonpath='{.spec.clusterNetwork.pods.cidrBlocks[*]}')"
    addon_cmd += "; export CIDR_LIST=$(bash -c 'echo $CIDRS' | tr ' ' ',')"
    addon_cmd += "; " + helm_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig install --repo https://raw.githubusercontent.com/kubernetes-sigs/cloud-provider-azure/master/helm/repo cloud-provider-azure --generate-name --set infra.clusterName=${CLUSTER_NAME} --set cloudControllerManager.clusterCIDR=${CIDR_LIST}"
    if "flatcar" in flavor_name:  # append caCetDir location to the cloud-provider-azure helm install command for flatcar flavor
        addon_cmd += " --set-string cloudControllerManager.caCertDir=/usr/share/ca-certificates"

    if "azure-cni-v1" in flavor_name:
        addon_cmd += "; " + kubectl_cmd + " apply -f ./templates/addons/azure-cni-v1.yaml --kubeconfig ./${CLUSTER_NAME}.kubeconfig"
    else:
        # install calico
        if "ipv6" in flavor_name:
            calico_values = "./templates/addons/calico-ipv6/values.yaml"
        elif "dual-stack" in flavor_name:
            calico_values = "./templates/addons/calico-dual-stack/values.yaml"
        else:
            calico_values = "./templates/addons/calico/values.yaml"
        addon_cmd += "; " + helm_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig install --repo https://docs.tigera.io/calico/charts --version ${CALICO_VERSION} calico tigera-operator -f " + calico_values + " --namespace tigera-operator --create-namespace"

    return addon_cmd

def base64_encode(to_encode):
    encode_blob = local("echo '{}' | tr -d '\n' | base64 | tr -d '\n'".format(to_encode), quiet = True, echo_off = True)
    return str(encode_blob)

def base64_encode_file(path_to_encode):
    encode_blob = local("cat {} | tr -d '\n' | base64 | tr -d '\n'".format(path_to_encode), quiet = True)
    return str(encode_blob)

def read_file_from_path(path_to_read):
    str_blob = local("cat {} | tr -d '\n'".format(path_to_read), quiet = True)
    return str(str_blob)

def base64_decode(to_decode):
    decode_blob = local("echo '{}' | base64 --decode".format(to_decode), quiet = True, echo_off = True)
    return str(decode_blob)

def kustomizesub(folder):
    yaml = local("hack/kustomize-sub.sh {}".format(folder), quiet = True)
    return yaml

def waitforsystem():
    local(kubectl_cmd + " wait --for=condition=ready --timeout=300s pod --all -n capi-kubeadm-bootstrap-system")
    local(kubectl_cmd + " wait --for=condition=ready --timeout=300s pod --all -n capi-kubeadm-control-plane-system")
    local(kubectl_cmd + " wait --for=condition=ready --timeout=300s pod --all -n capi-system")

##############################
# Actual work happens here
##############################

validate_auth()

include_user_tilt_files()

load("ext://cert_manager", "deploy_cert_manager")

if settings.get("deploy_cert_manager"):
    deploy_cert_manager(version = settings.get("cert_manager_version"))

deploy_capi()

deploy_caaph()

create_identity_secret()

capz()

observability()

waitforsystem()

create_crs()

flavors()
