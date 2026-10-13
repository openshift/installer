# Migrate the Image Registry from Cinder to Swift

This document describes:
* [Migrating existing images to a new Swift container](#migrate-the-images-from-the-volume-to-swift)
* [Setting Swift as the Image Registry backend](#switch-to-swift-as-the-image-registry-storage-backend)
* [Aborting the migration to return to the volume backend](#abort-the-migration-and-return-to-the-volume-backend)

The persistent volume backend does not support multiple replicas by default,
because they typically only allow one read-write client at once. Object storage
does allow for multiple Image Registry replicas.

The migration as described only incur minimal downtime; however, the registry
will be set to read-only mode during data migration.

## Migrate the images from the volume to Swift

As the first step, create the Swift container that will consitute the new Image
Registry backend. In order for it to be garbage-collected by the OpenShift
installer in case of `destroy`, it has to have a specific name and specific
properties (the suffix "aaa" here is arbitrary):

```shell
infrastructure_name="$(oc get infrastructure/cluster -o jsonpath='{.status.infrastructureName}')"
container_name="${infrastructure_name}-image-registry-aaa"
openstack container create "$container_name"
openstack container set --property Name="$container_name" --property Openshiftclusterid="$infrastructure_name" "$container_name"
```

Then, build the container image that will perform the migration from within the
cluster.

The container will mount two volumes:
* one volume with the OpenStack credentials
* the current Image Registry storage, mounted read-only mode.

The example below reads the credentials and makes them available to the Swift
client. Then the Swift client moves data from the registry volume to the new
Swift container.

Build and upload the image to a registry your OpenShift instance has access to.
It can also be OpenShift's Image Registry itself. We will use
`quay.io/shiftstack/image-registry-migrator:v1` in our example.

```Containerfile
# Containerfile
FROM registry.access.redhat.com/ubi9/ubi:latest

RUN dnf -y install python3-pip jq

RUN pip install python-keystoneclient python-swiftclient yq

RUN true \
	&& echo $'export OS_AUTH_URL="$(yq -r \'.clouds.openstack.auth.auth_url // empty\' /etc/openstack/clouds.yaml)"'                                                 >> /migrate.sh \
	&& echo $'export OS_IDENTITY_API_VERSION="$(yq -r \'.clouds.openstack.identity_api_version // empty\' /etc/openstack/clouds.yaml)"'                              >> /migrate.sh \
	&& echo $'export OS_REGION_NAME="$(yq -r \'.clouds.openstack.region_name // empty\' /etc/openstack/clouds.yaml)"'                                                >> /migrate.sh \
	&& echo $'export OS_AUTH_TYPE="$(yq -r \'.clouds.openstack.auth_type // empty\' /etc/openstack/clouds.yaml)"'                                                    >> /migrate.sh \
	&& echo $''                                                                                                                                                      >> /migrate.sh \
	&& echo $'if [ "$OS_AUTH_TYPE" == "v3applicationcredential" ]; then'                                                                                             >> /migrate.sh \
	&& echo $'	export OS_APPLICATION_CREDENTIAL_ID="$(yq -r \'.clouds.openstack.auth.application_credential_id // empty\' /etc/openstack/clouds.yaml)"'         >> /migrate.sh \
	&& echo $'	export OS_APPLICATION_CREDENTIAL_NAME="$(yq -r \'.clouds.openstack.auth.application_credential_name // empty\' /etc/openstack/clouds.yaml)"'     >> /migrate.sh \
	&& echo $'	export OS_APPLICATION_CREDENTIAL_SECRET="$(yq -r \'.clouds.openstack.auth.application_credential_secret // empty\' /etc/openstack/clouds.yaml)"' >> /migrate.sh \
	&& echo $'else'                                                                                                                                                  >> /migrate.sh \
	&& echo $'	export OS_USERNAME="$(yq -r \'.clouds.openstack.auth.username // empty\' /etc/openstack/clouds.yaml)"'                                           >> /migrate.sh \
	&& echo $'	export OS_PASSWORD="$(yq -r \'.clouds.openstack.auth.password // empty\' /etc/openstack/clouds.yaml)"'                                           >> /migrate.sh \
	&& echo $'	export OS_USER_DOMAIN_NAME="$(yq -r \'.clouds.openstack.auth.user_domain_name // empty\' /etc/openstack/clouds.yaml)"'                           >> /migrate.sh \
	&& echo $'	export OS_PROJECT_ID="$(yq -r \'.clouds.openstack.auth.project_id // empty\' /etc/openstack/clouds.yaml)"'                                       >> /migrate.sh \
	&& echo $'	export OS_PROJECT_NAME="$(yq -r \'.clouds.openstack.auth.project_name // empty\' /etc/openstack/clouds.yaml)"'                                   >> /migrate.sh \
	&& echo $'fi'                                                                                                                                                    >> /migrate.sh \
	&& echo $''                                                                                                                                                      >> /migrate.sh \
	&& echo $'swift upload "$REGISTRY_STORAGE_SWIFT_CONTAINER" /files'                                                                                               >> /migrate.sh \
	&& chmod +x /migrate.sh
	
CMD [ "/migrate.sh" ]
```

Here is the definition of the pod that will run the container. Replace the
example image name `quay.io/shiftstack/image-registry-migrator:v1` in the last
line with your own image. Also replace `openshift-lrgv9-image-registry-aaa`
with the `$container_name` above:

```yaml
# image-migration-pod.yaml
kind: Pod
apiVersion: v1
metadata:
  name: image-registry-migration
  namespace: openshift-image-registry
  labels:
    app: image-registry-migration
  annotations:
    openshift.io/scc: anyuid
spec:
  volumes:
    - name: old-registry
      persistentVolumeClaim:
        claimName: image-registry-storage
    - name: openstack
      secret:
        secretName: installer-cloud-credentials
  restartPolicy: OnFailure
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  containers:
    - name: migrate-ir
      imagePullPolicy: Always
      volumeMounts:
        - name: old-registry
          readOnly: true
          mountPath: /files
        - name: openstack
          mountPath: "/etc/openstack"
          readOnly: true
      image: quay.io/shiftstack/image-registry-migrator:v1
      env:
        - name: REGISTRY_STORAGE_SWIFT_CONTAINER
          value: openshift-lrgv9-image-registry-aaa
```

**In order to avoid corruption of existing data during the migration, switch
the Image Registry to read-only mode.** Note that this action restarts the one
Image Registry replica available at this point; this means a brief downtime.

```shell
oc patch configs.imageregistry.operator.openshift.io/cluster -p '{"spec":{"readOnly":true}}' --type=merge
```

Now run the container in your cluster:

```shell
oc apply -f image-migration-pod.yaml
```

Check for any errors:

```shell
oc -n openshift-image-registry get pod image-registry-migration
```

When its status is `Completed`, the pod can be removed:

```shell
oc -n openshift-image-registry delete pod image-registry-migration
```

## Switch to Swift as the Image Registry storage backend

This command sets Swift as the image-registry backend and sets two replicas, as
Swift allows concurrent access to its resources. Replace
`openshift-lrgv9-image-registry-aaa` with the `$container_name` above:

```shell
oc patch configs.imageregistry.operator.openshift.io/cluster \
    --type='json' \
    --patch='[
        {"op": "replace", "path": "/spec/rolloutStrategy", "value": "RollingUpdate"},
        {"op": "replace", "path": "/spec/replicas", "value": 2},
        {"op": "replace", "path": "/spec/storage", "value": {"managementState": "Managed", "swift": {"container": "openshift-lrgv9-image-registry-aaa"}}}
    ]'
```

Wait for the operator to recreate the image-registry instances.

```shell
oc wait --for=condition=Progressing=false clusteroperator/image-registry --timeout=5m
```

After the operation is complete and successful, re-enable write operations as
needed:

```shell
oc patch configs.imageregistry.operator.openshift.io/cluster -p '{"spec":{"readOnly":false}}' --type=merge
```

After the successful migration of all images has been confirmed, the old volume
claim can be released:

```shell
oc -n openshift-image-registry delete pvc image-registry-storage
```

## Abort the migration and return to the volume backend

This command rolls back to the previous volume backend. Use it to return to a
stable state if a successful migration can't be confirmed.

```shell
oc patch configs.imageregistry.operator.openshift.io/cluster \
    --type='json' \
    --patch='[
        {"op": "replace", "path": "/spec/rolloutStrategy", "value": "Recreate"},
        {"op": "replace", "path": "/spec/replicas", "value": 1},
        {"op": "replace", "path": "/spec/storage", "value": {"managementState": "Unmanaged", "pvc": {"claim":"image-registry-storage"}}}
    ]'
```
