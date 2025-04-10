# -*- mode: Python -*-

load("ext://cert_manager", "deploy_cert_manager")

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
    "capi_version": "v1.9.5",
    "caaph_version": "v0.2.5",
    "cert_manager_version": "v1.16.2",
    "kubernetes_version": "v1.30.2",
    "aks_kubernetes_version": "v1.30.2",
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

# Helper function to update settings and environment variables
def update_settings_and_env(source_settings):
    source_dict = settings.get(source_settings, {})
    os.environ.update(source_dict)

    # Update settings with lowercase values
    for key, value in source_dict.items():
        # print("key: %s, value: %s" % (key, value))
        if key.lower() in settings:
            settings[key.lower()] = str(value)

# Update environment variables with AKS as management cluster settings
if "aks_as_mgmt_settings" in settings:
    update_settings_and_env("aks_as_mgmt_settings")

# kustomize_substitutions takes precedence over aks_as_mgmt_settings if both are set
if "kustomize_substitutions" in settings:
    update_settings_and_env("kustomize_substitutions")

# Pretty print settings
def pretty_print_dict(d, indent = 0):
    for key, value in d.items():
        indent_str = " " * indent
        value_type = str(type(value))

        if value_type == "dict":
            print(indent_str + str(key) + ":")
            pretty_print_dict(value, indent + 2)
        elif value_type == "list":
            print(indent_str + str(key) + ":")
            for item in value:
                if str(type(item)) == "dict":
                    pretty_print_dict(item, indent + 2)
                else:
                    print(indent_str + "  - " + str(item))
        else:
            print(indent_str + str(key) + ": " + str(value))

os_arch = str(local("go env GOARCH")).rstrip("\n")

# set os_arch to amd64 if using AKS as management cluster
if "aks_as_mgmt_settings" in settings and os_arch != "amd64":
    YELLOW = "\033[1;33m"
    RESET = "\033[0m"
    print("\n" + YELLOW + "ARCHITECTURE OVERRIDE: Using AKS as management cluster requires CAPZ to be built for amd64 architecture." + RESET)
    print(YELLOW + "Ignoring local GOARCH output and forcing CAPZ's os_arch to amd64" + RESET + "\n")
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
FROM golang:1.22 AS tilt-helper
# Support live reloading with Tilt
RUN wget --output-document /restart.sh --quiet https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/restart.sh  && \
    wget --output-document /start.sh --quiet https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/start.sh && \
    chmod +x /start.sh && chmod +x /restart.sh && \
    touch /process.txt && chmod 0777 /process.txt `# pre-create PID file to allow even non-root users to run the image`
"""

tilt_dockerfile_header = """
FROM gcr.io/distroless/base:debug AS tilt
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
    if settings.get("container_args"):
        capz_container_args = settings.get("container_args").get("capz-controller-manager")
        yaml_dict = decode_yaml_stream(yaml)
        append_arg_for_container_in_deployment(yaml_dict, "capz-controller-manager", "capz-system", "cluster-api-azure-controller", capz_container_args)
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
    if os.getenv("REGISTRY", "") != "":
        registry = os.getenv("REGISTRY", "")
        print("\nUsing REGISTRY: " + registry + "\n")
        image = registry + "/cluster-api-azure-controller"
    else:
        image = "gcr.io/k8s-staging-cluster-api-azure/cluster-api-azure-controller"

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

    if "total_nodes" not in settings:
        settings["total_nodes"] = {}
    for template in template_list:
        deploy_worker_templates(template, substitutions)

    delete_all_workload_clusters = kubectl_cmd + " delete clusters --all --wait=false;"

    if "aks_as_mgmt_settings" in settings and os.getenv("SUBSCRIPTION_TYPE", "") == "corporate":
        delete_all_workload_clusters += clear_aks_vnet_peerings()

    local_resource(
        name = "delete-all-workload-clusters",
        cmd = ["sh", "-ec", delete_all_workload_clusters],
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

    total_nodes = 0
    for substitution in substitutions:
        value = substitutions[substitution]
        if substitution == "CONTROL_PLANE_MACHINE_COUNT" or substitution == "WORKER_MACHINE_COUNT":
            count = yaml.count(substitution)
            total_nodes += int(value) * count

        # NOTE: ENV Variables of type ${VAR:=default} are not replaced below.
        yaml = yaml.replace("${" + substitution + "}", value)

    yaml = shlex.quote(yaml)
    flavor_name = os.path.basename(flavor)
    if total_nodes > 0:
        settings["total_nodes"][flavor_name] = total_nodes

    # Flavor command is built from here
    flavor_cmd = "RANDOM=$(bash -c 'echo $RANDOM'); "

    if "aks_as_mgmt_settings" in settings and os.getenv("SUBSCRIPTION_TYPE", "") == "corporate" and "aks" not in flavor_name:
        apiserver_lb_private_ip = os.getenv("AZURE_INTERNAL_LB_PRIVATE_IP", "")
        if "windows-apiserver-ilb" in flavor and apiserver_lb_private_ip == "":
            flavor_cmd += "export AZURE_INTERNAL_LB_PRIVATE_IP=\"40.0.11.100\"; "
        elif "apiserver-ilb" in flavor and apiserver_lb_private_ip == "":
            flavor_cmd += "export AZURE_INTERNAL_LB_PRIVATE_IP=\"30.0.11.100\"; "

    flavor_cmd += "export CLUSTER_NAME=" + flavor.replace("windows", "win") + "-$RANDOM; echo " + yaml + "> ./.tiltbuild/" + flavor + "; cat ./.tiltbuild/" + flavor + " | " + envsubst_cmd + " | " + kubectl_cmd + " apply -f -; "

    if "aks_as_mgmt_settings" in settings and os.getenv("SUBSCRIPTION_TYPE", "") == "corporate" and "aks" not in flavor_name:
        flavor_cmd += peer_vnets()

    # wait for kubeconfig to be available
    flavor_cmd += '''
    echo "Waiting for kubeconfig to be available";
    until ''' + kubectl_cmd + """ get secret ${CLUSTER_NAME}-kubeconfig > /dev/null 2>&1; do sleep 5; done;
    """ + kubectl_cmd + ''' get secret ${CLUSTER_NAME}-kubeconfig -o jsonpath={.data.value} | base64 --decode > ./${CLUSTER_NAME}.kubeconfig;
    chmod 600 ./${CLUSTER_NAME}.kubeconfig;
    echo "Kubeconfig for ${CLUSTER_NAME} created and saved in the local";
    echo "Waiting for ${CLUSTER_NAME} API Server to be accessible";
    until ''' + kubectl_cmd + ''' --kubeconfig=./${CLUSTER_NAME}.kubeconfig get nodes > /dev/null 2>&1; do sleep 5; done;
    echo "API Server of ${CLUSTER_NAME} is accessible";
    '''

    # copy the kubeadm configmap to the calico-system namespace.
    # This is a workaround needed for the calico-node-windows daemonset to be able to run in the calico-system namespace.
    if "windows" in flavor_name:
        flavor_cmd += """
        until """ + kubectl_cmd + """ --kubeconfig ./${CLUSTER_NAME}.kubeconfig get configmap kubeadm-config --namespace=kube-system > /dev/null 2>&1; do sleep 5; done;
        """ + kubectl_cmd + """ --kubeconfig ./${CLUSTER_NAME}.kubeconfig create namespace calico-system --dry-run=client -o yaml |         """ + kubectl_cmd + """ --kubeconfig ./${CLUSTER_NAME}.kubeconfig apply -f -;
        """ + kubectl_cmd + """ --kubeconfig ./${CLUSTER_NAME}.kubeconfig get configmap kubeadm-config --namespace=kube-system -o yaml |         sed 's/namespace: kube-system/namespace: calico-system/' |         """ + kubectl_cmd + """ --kubeconfig ./${CLUSTER_NAME}.kubeconfig apply -f -;
        """

    if "aks_as_mgmt_settings" in settings and os.getenv("SUBSCRIPTION_TYPE", "") == "corporate" and "aks" not in flavor_name:
        flavor_cmd += create_private_dns_zone()

    flavor_cmd += get_addons(flavor_name)
    if settings.get("total_nodes").get(flavor_name, 0) > 0:
        flavor_cmd += check_nodes_ready(flavor_name)
    flavor_cmd += "echo \"Cluster ${CLUSTER_NAME} created, don't forget to delete\"; "

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

    addon_cmd = "export CIDRS=$(" + kubectl_cmd + " get cluster ${CLUSTER_NAME} -o jsonpath='{.spec.clusterNetwork.pods.cidrBlocks[*]}'); "
    addon_cmd += "export CIDR_LIST=$(bash -c 'echo $CIDRS' | tr ' ' ','); "
    addon_cmd += helm_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig install --repo https://raw.githubusercontent.com/kubernetes-sigs/cloud-provider-azure/master/helm/repo cloud-provider-azure --generate-name --set infra.clusterName=${CLUSTER_NAME} --set cloudControllerManager.clusterCIDR=${CIDR_LIST}"
    if "flatcar" in flavor_name:  # append caCetDir location to the cloud-provider-azure helm install command for flatcar flavor
        addon_cmd += " --set-string cloudControllerManager.caCertDir=/usr/share/ca-certificates"
    addon_cmd += "; "

    if "azure-cni-v1" in flavor_name:
        addon_cmd += kubectl_cmd + " apply -f ./templates/addons/azure-cni-v1.yaml --kubeconfig ./${CLUSTER_NAME}.kubeconfig; "
    else:
        # install calico
        if "ipv6" in flavor_name:
            calico_values = "./templates/addons/calico-ipv6/values.yaml"
        elif "dual-stack" in flavor_name:
            calico_values = "./templates/addons/calico-dual-stack/values.yaml"
        else:
            calico_values = "./templates/addons/calico/values.yaml"
        addon_cmd += helm_cmd + " --kubeconfig ./${CLUSTER_NAME}.kubeconfig install --repo https://docs.tigera.io/calico/charts --version ${CALICO_VERSION} calico tigera-operator -f " + calico_values + " --namespace tigera-operator --create-namespace; "

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

def peer_vnets():
    # TODO: check for az cli to be installed in local
    peering_cmd = '''
    echo "--------Peering VNETs--------";
    az network vnet wait --resource-group ${AKS_RESOURCE_GROUP} --name ${AKS_MGMT_VNET_NAME} --created --timeout 180;
    export MGMT_VNET_ID=$(az network vnet show --resource-group ${AKS_RESOURCE_GROUP} --name ${AKS_MGMT_VNET_NAME} --query id --output tsv);
    echo "1/4 ${AKS_MGMT_VNET_NAME} found ";

    az network vnet wait --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME}-vnet --created --timeout 180;
    export WORKLOAD_VNET_ID=$(az network vnet show --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME}-vnet --query id --output tsv);
    echo "2/4 ${CLUSTER_NAME}-vnet found ";

    az network vnet peering create --name mgmt-to-${CLUSTER_NAME} --resource-group ${AKS_RESOURCE_GROUP} --vnet-name ${AKS_MGMT_VNET_NAME} --remote-vnet \"${WORKLOAD_VNET_ID}\" --allow-vnet-access true --allow-forwarded-traffic true --only-show-errors --output none;
    az network vnet peering wait --name mgmt-to-${CLUSTER_NAME} --resource-group ${AKS_RESOURCE_GROUP} --vnet-name ${AKS_MGMT_VNET_NAME} --created --timeout 300 --only-show-errors --output none;
    echo "3/4 mgmt-to-${CLUSTER_NAME} peering created in ${AKS_MGMT_VNET_NAME}";

    az network vnet peering create --name ${CLUSTER_NAME}-to-mgmt --resource-group ${CLUSTER_NAME} --vnet-name ${CLUSTER_NAME}-vnet --remote-vnet \"${MGMT_VNET_ID}\" --allow-vnet-access true --allow-forwarded-traffic true --only-show-errors --output none;
    az network vnet peering wait --name ${CLUSTER_NAME}-to-mgmt --resource-group ${CLUSTER_NAME} --vnet-name ${CLUSTER_NAME}-vnet --created --timeout 300 --only-show-errors --output none;
    echo "4/4 ${CLUSTER_NAME}-to-mgmt peering created in ${CLUSTER_NAME}-vnet";
    '''

    return peering_cmd

def create_private_dns_zone():
    create_private_dns_zone_cmd = '''
    echo "--------Creating private DNS zone--------";
    az network private-dns zone create --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --only-show-errors --output none;
    az network private-dns zone wait --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --created --timeout 300 --only-show-errors --output none;
    echo "1/4 ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com private DNS zone created in ${CLUSTER_NAME}";

    az network private-dns link vnet create --resource-group ${CLUSTER_NAME} --zone-name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --name ${CLUSTER_NAME}-to-mgmt --virtual-network \"${WORKLOAD_VNET_ID}\" --registration-enabled false --only-show-errors --output none;
    az network private-dns link vnet wait --resource-group ${CLUSTER_NAME} --zone-name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --name ${CLUSTER_NAME}-to-mgmt --created --timeout 300 --only-show-errors --output none;
    echo "2/4 workload cluster vnet ${CLUSTER_NAME}-vnet linked with private DNS zone";

    az network private-dns link vnet create --resource-group ${CLUSTER_NAME} --zone-name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --name mgmt-to-${CLUSTER_NAME} --virtual-network \"${MGMT_VNET_ID}\" --registration-enabled false --only-show-errors --output none;
    az network private-dns link vnet wait --resource-group ${CLUSTER_NAME} --zone-name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --name mgmt-to-${CLUSTER_NAME} --created --timeout 300 --only-show-errors --output none;
    echo "3/4 management cluster vnet ${AKS_MGMT_VNET_NAME} linked with private DNS zone";

    az network private-dns record-set a add-record --resource-group ${CLUSTER_NAME} --zone-name ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com --record-set-name \"@\" --ipv4-address ${AZURE_INTERNAL_LB_PRIVATE_IP} --only-show-errors --output none;
    echo "4/4 private DNS zone record @ created to point ${CLUSTER_NAME}-${APISERVER_LB_DNS_SUFFIX}.${AZURE_LOCATION}.cloudapp.azure.com to ${AZURE_INTERNAL_LB_PRIVATE_IP}";
    '''

    return create_private_dns_zone_cmd

def check_nodes_ready(flavor_name):
    total_nodes = settings.get("total_nodes", {}).get(flavor_name, 0)
    check_nodes_ready_cmd = '''
    echo "--------Checking if all nodes are available and ready--------";
    echo "Waiting for ''' + str(total_nodes) + ''' nodes to be available...";
    TIMEOUT=600  # 10 minutes timeout
    START_TIME=$(date +%s);
    while true; do
        NODE_COUNT=$(''' + kubectl_cmd + ''' get nodes --kubeconfig=./${CLUSTER_NAME}.kubeconfig --no-headers | wc -l);
        if [ "$NODE_COUNT" -eq ''' + str(total_nodes) + ''' ]; then
            break;
        fi;
        CURRENT_TIME=$(date +%s);
        ELAPSED_TIME=$((CURRENT_TIME - START_TIME));
        if [ $ELAPSED_TIME -gt $TIMEOUT ]; then
            echo "Timeout waiting for nodes to be available. Found $NODE_COUNT nodes, expected ''' + str(total_nodes) + '''";
            exit 1;
        fi;
        echo "Found $NODE_COUNT nodes, waiting for ''' + str(total_nodes) + ''' nodes... (${ELAPSED_TIME}s elapsed)";
        sleep 10;
    done;

    echo "All ''' + str(total_nodes) + ''' nodes are available, waiting for them to be ready...";
    ''' + kubectl_cmd + """ wait --for=condition=ready node --all --kubeconfig=./${CLUSTER_NAME}.kubeconfig --timeout=600s;

    READY_NODES=$(""" + kubectl_cmd + ''' get nodes --kubeconfig=./${CLUSTER_NAME}.kubeconfig -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}' | tr ' ' '\\n' | grep -c "True");
    if [ "$READY_NODES" -eq ''' + str(total_nodes) + ''' ]; then
        echo "All ''' + str(total_nodes) + ''' nodes are ready!";
    else
        echo "Expected ''' + str(total_nodes) + ''' nodes to be ready but got $READY_NODES ready nodes";
        exit 1;
    fi;
    '''
    return check_nodes_ready_cmd

def clear_aks_vnet_peerings():
    delete_peering_cmd = '''
    echo "--------Clearing AKS MGMT VNETs Peerings--------";
    az network vnet wait --resource-group ${AKS_RESOURCE_GROUP} --name ${AKS_MGMT_VNET_NAME} --created --timeout 180;
    echo "VNet ${AKS_MGMT_VNET_NAME} found ";

    PEERING_NAMES=$(az network vnet peering list --resource-group ${AKS_RESOURCE_GROUP} --vnet-name ${AKS_MGMT_VNET_NAME} --query "[].name" --output tsv);
    for PEERING_NAME in ${PEERING_NAMES}; do echo "Deleting peering: ${PEERING_NAME}"; az network vnet peering delete --name ${PEERING_NAME} --resource-group ${AKS_RESOURCE_GROUP} --vnet-name ${AKS_MGMT_VNET_NAME}; done;
    echo "All VNETs Peerings deleted in ${AKS_MGMT_VNET_NAME}";
    '''

    return delete_peering_cmd

# allow_tcp_udp_ports is a helper function to allow TCP ports: 443,6443 and UDP ports: 53 on the management cluster when using aks as management cluster.
def allow_tcp_udp_ports():
    allow_tcp_udp_ports_cmd = '''
    echo "--------Allowing TCP ports: 443,6443 and UDP port: 53 on the management cluster's API server--------";

    # Define ports to allow
    TCP_PORTS="443 6443"
    UDP_PORTS="53"
    TIMEOUT=3000
    SLEEP_INTERVAL=10

    echo "Waiting for NSG rules with prefix 'NRMS-Rule-101' to appear...";

    # Process NSGs in each resource group
    for RG in ${AKS_NODE_RESOURCE_GROUP} ${AKS_RESOURCE_GROUP}; do
        echo "Processing NSGs in resource group '$RG'...";

        # Wait for NSGs to appear in the resource group
        RG_START_TIME=$(date +%s);
        while true; do
            NSG_LIST=$(az network nsg list --resource-group "$RG" --query "[].name" --output tsv);
            if [ -n "$NSG_LIST" ]; then break; fi;
            if [ $(($(date +%s) - RG_START_TIME)) -ge $TIMEOUT ]; then
                echo "Timeout waiting for NSGs in resource group '$RG'";
                continue 2;
            fi;
            echo "No NSGs found in '$RG' yet, waiting...";
            sleep $SLEEP_INTERVAL;
        done;

        # Process each NSG in the resource group
        for NSG in $NSG_LIST; do
            echo "Checking for NRMS-Rule-101 rules in NSG '$NSG' in resource group '$RG'...";

            # Wait for NRMS rules to appear
            RULE_START_TIME=$(date +%s);
            while true; do
                TCP_RULE_FOUND=$(az network nsg rule list --resource-group "$RG" --nsg-name "$NSG" --query "[?starts_with(name, 'NRMS-Rule-101')].name" --output tsv);
                UDP_RULE_FOUND=$(az network nsg rule list --resource-group "$RG" --nsg-name "$NSG" --query "[?starts_with(name, 'NRMS-Rule-103')].name" --output tsv);
                if [ -n "$TCP_RULE_FOUND" ] && [ -n "$UDP_RULE_FOUND" ]; then
                    echo "! --- Found NRMS rules in NSG '$NSG': --- !";
                    echo "! --- TCP Rule: $TCP_RULE_FOUND --- !";
                    echo "! --- UDP Rule: $UDP_RULE_FOUND --- !";
                    break;
                fi;
                if [ $(($(date +%s) - RULE_START_TIME)) -ge $TIMEOUT ]; then
                    echo "Timeout waiting for NRMS rules in NSG '$NSG' in RG '$RG'. Skipping NSG.";
                    continue 2;
                fi;
                echo "Waiting for NRMS rules in NSG '$NSG'...";
                if [ -z "$TCP_RULE_FOUND" ]; then echo "TCP Rule (NRMS-Rule-101) not found"; fi;
                if [ -z "$UDP_RULE_FOUND" ]; then echo "UDP Rule (NRMS-Rule-103) not found"; fi;
                sleep $SLEEP_INTERVAL;
            done;

            # Update TCP rule
            if az network nsg rule show --resource-group "$RG" --nsg-name "$NSG" --name "NRMS-Rule-101" --output none 2>/dev/null; then
                echo " - Updating NRMS-Rule-101 (TCP) in NSG '$NSG' of RG '$RG'...";
                az network nsg rule update \
                    --resource-group "$RG" \
                    --nsg-name "$NSG" \
                    --name "NRMS-Rule-101" \
                    --access Allow \
                    --direction Inbound \
                    --protocol "TCP" \
                    --destination-port-ranges $TCP_PORTS \
                    --destination-address-prefixes "*" \
                    --source-address-prefixes "*" \
                    --source-port-ranges "*" \
                    --only-show-errors --output none || echo "Failed to update NRMS-Rule-101";

                # Update UDP rule
                echo " - Updating NRMS-Rule-103 (UDP) in NSG '$NSG' of RG '$RG'...";
                az network nsg rule update \
                    --resource-group "$RG" \
                    --nsg-name "$NSG" \
                    --name "NRMS-Rule-103" \
                    --access Allow \
                    --direction Inbound \
                    --protocol "UDP" \
                    --destination-port-ranges $UDP_PORTS \
                    --destination-address-prefixes "*" \
                    --source-address-prefixes "*" \
                    --source-port-ranges "*" \
                    --only-show-errors --output none || echo "Failed to update NRMS-Rule-103";
            fi;
        done;
    done;

    echo "\nNSG NRMS rule check and modification complete.";
    '''

    local_resource(
        name = "allow required ports on mgmt cluster",
        cmd = ["sh", "-ec", allow_tcp_udp_ports_cmd],
        auto_init = False,
        trigger_mode = TRIGGER_MODE_MANUAL,
        labels = ["cluster-api"],
        allow_parallel = True,
    )

##############################
# Actual work happens here
##############################

validate_auth()

include_user_tilt_files()

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

if "aks_as_mgmt_settings" in settings and os.getenv("SUBSCRIPTION_TYPE", "") == "corporate":
    allow_tcp_udp_ports()

print("\n\n=== Active Tilt Configuration Settings ===")
pretty_print_dict(settings)
print("=======================================\n")
