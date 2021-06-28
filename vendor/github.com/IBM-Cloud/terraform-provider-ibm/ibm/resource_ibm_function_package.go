// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	funcPkgNamespace    = "namespace"
	funcPkgName         = "name"
	funcPkgUsrDefAnnots = "user_defined_annotations"
	funcPkgUsrDefParams = "user_defined_parameters"
	funcPkgBindPkgName  = "bind_package_name"
)

func resourceIBMFunctionPackage() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFunctionPackageCreate,
		Read:     resourceIBMFunctionPackageRead,
		Update:   resourceIBMFunctionPackageUpdate,
		Delete:   resourceIBMFunctionPackageDelete,
		Exists:   resourceIBMFunctionPackageExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			funcPkgNamespace: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IBM Cloud function namespace.",
				ValidateFunc: InvokeValidator("ibm_function_package", funcPkgNamespace),
			},
			funcPkgName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of package.",
				ValidateFunc: InvokeValidator("ibm_function_package", funcPkgName),
			},
			"publish": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Package visibilty.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the item.",
			},
			funcPkgUsrDefAnnots: {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Annotation values in KEY VALUE format.",
				Default:          "[]",
				ValidateFunc:     InvokeValidator("ibm_function_package", funcPkgUsrDefAnnots),
				DiffSuppressFunc: suppressEquivalentJSON,
				StateFunc: func(v interface{}) string {
					json, _ := normalizeJSONString(v)
					return json
				},
			},
			funcPkgUsrDefParams: {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Parameters values in KEY VALUE format. Parameter bindings included in the context.TODO() passed to the package.",
				ValidateFunc:     InvokeValidator("ibm_function_package", funcPkgUsrDefParams),
				Default:          "[]",
				DiffSuppressFunc: suppressEquivalentJSON,
				StateFunc: func(v interface{}) string {
					json, _ := normalizeJSONString(v)
					return json
				},
			},
			"annotations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All annotations set on package by user and those set by the IBM Cloud Function backend/API.",
			},
			"parameters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All parameters set on package by user and those set by the IBM Cloud Function backend/API.",
			},
			funcPkgBindPkgName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Name of package to be binded.",
				ValidateFunc: InvokeValidator("ibm_function_package", funcPkgBindPkgName),
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					if strings.Compare(n, o) == 0 {
						return true
					}
					return false
				},
			},
			"package_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMFuncPackageValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcPkgName,
			ValidateFunctionIdentifier: ValidateRegexp,
			Type:                       TypeString,
			Regexp:                     `\A([\w]|[\w][\w@ .-]*[\w@.-]+)\z`,
			Required:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcPkgNamespace,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcPkgUsrDefAnnots,
			ValidateFunctionIdentifier: ValidateJSONString,
			Type:                       TypeString,
			Default:                    "[]",
			Optional:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcPkgBindPkgName,
			ValidateFunctionIdentifier: ValidateBindedPackageName,
			Type:                       TypeString,
			Optional:                   true})

	ibmFuncPackageResourceValidator := ResourceValidator{ResourceName: "ibm_function_package", Schema: validateSchema}
	return &ibmFuncPackageResourceValidator
}

func resourceIBMFunctionPackageCreate(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	namespace := d.Get("namespace").(string)
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	packageService := wskClient.Packages

	name := d.Get("name").(string)

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(name); err != nil {
		return NewQualifiedNameError(name, err)
	}

	payload := whisk.Package{
		Name:      qualifiedName.GetEntityName(),
		Namespace: qualifiedName.GetNamespace(),
	}

	userDefinedAnnotations := d.Get("user_defined_annotations").(string)
	payload.Annotations, err = expandAnnotations(userDefinedAnnotations)
	if err != nil {
		return err
	}

	userDefinedParameters := d.Get("user_defined_parameters").(string)
	payload.Parameters, err = expandParameters(userDefinedParameters)
	if err != nil {
		return err
	}

	if publish, ok := d.GetOk("publish"); ok {
		p := publish.(bool)
		payload.Publish = &p
	}

	if v, ok := d.GetOk("bind_package_name"); ok {
		var BindingQualifiedName = new(QualifiedName)
		if BindingQualifiedName, err = NewQualifiedName(v.(string)); err != nil {
			return NewQualifiedNameError(v.(string), err)
		}
		BindingPayload := whisk.Binding{
			Name:      BindingQualifiedName.GetEntityName(),
			Namespace: BindingQualifiedName.GetNamespace(),
		}
		payload.Binding = &BindingPayload
	}

	log.Println("[INFO] Creating IBM CLoud Function package")
	result, _, err := packageService.Insert(&payload, false)
	if err != nil {
		return fmt.Errorf("Error creating IBM CLoud Function package: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", namespace, result.Name))

	return resourceIBMFunctionPackageRead(d, meta)
}

func resourceIBMFunctionPackageRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := ""
	packageID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		packageID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		packageID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, packageID))
	}

	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}
	packageService := wskClient.Packages

	pkg, _, err := packageService.Get(packageID)
	if err != nil {
		return fmt.Errorf("Error retrieving IBM Cloud Function package %s : %s", packageID, err)
	}
	d.Set("package_id", pkg.Name)
	d.Set("name", pkg.Name)
	d.Set("namespace", namespace)
	d.Set("publish", pkg.Publish)
	d.Set("version", pkg.Version)
	annotations, err := flattenAnnotations(pkg.Annotations)
	if err != nil {
		return err
	}
	d.Set("annotations", annotations)
	parameters, err := flattenParameters(pkg.Parameters)
	if err != nil {
		return err
	}
	d.Set("parameters", parameters)
	if isEmpty(*pkg.Binding) {

		d.Set("user_defined_annotations", annotations)
		d.Set("user_defined_parameters", parameters)

	} else {
		d.Set("bind_package_name", fmt.Sprintf("/%s/%s", pkg.Binding.Namespace, pkg.Binding.Name))
		c, err := whisk.NewClient(http.DefaultClient, &whisk.Config{
			Namespace:         pkg.Binding.Namespace,
			AuthToken:         wskClient.AuthToken,
			Host:              wskClient.Host,
			AdditionalHeaders: wskClient.AdditionalHeaders,
		})
		bindedPkg, _, err := c.Packages.Get(pkg.Binding.Name)
		if err != nil {
			return fmt.Errorf("Error retrieving Binded IBM Cloud Function package %s : %s", pkg.Binding.Name, err)
		}

		userAnnotations, err := flattenAnnotations(filterInheritedAnnotations(bindedPkg.Annotations, pkg.Annotations))
		if err != nil {
			return err
		}
		d.Set("user_defined_annotations", userAnnotations)

		userParameters, err := flattenParameters(filterInheritedParameters(bindedPkg.Parameters, pkg.Parameters))
		if err != nil {
			return err
		}
		d.Set("user_defined_parameters", userParameters)
	}

	return nil
}

func resourceIBMFunctionPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]

	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	packageService := wskClient.Packages

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(d.Get("name").(string)); err != nil {
		return NewQualifiedNameError(d.Get("name").(string), err)
	}

	payload := whisk.Package{
		Name:      qualifiedName.GetEntityName(),
		Namespace: qualifiedName.GetNamespace(),
	}
	ischanged := false
	if d.HasChange("publish") {
		p := d.Get("publish").(bool)
		payload.Publish = &p
		ischanged = true
	}

	if d.HasChange("user_defined_parameters") {
		var err error
		payload.Parameters, err = expandParameters(d.Get("user_defined_parameters").(string))
		if err != nil {
			return err
		}
		ischanged = true
	}

	if d.HasChange("user_defined_annotations") {
		var err error
		payload.Annotations, err = expandAnnotations(d.Get("user_defined_annotations").(string))
		if err != nil {
			return err
		}
		ischanged = true
	}

	if ischanged {
		log.Println("[INFO] Update IBM Cloud Function Package")
		_, _, err = packageService.Insert(&payload, true)
		if err != nil {
			return fmt.Errorf("Error updating IBM Cloud Function Package: %s", err)
		}
	}

	return resourceIBMFunctionPackageRead(d, meta)
}

func resourceIBMFunctionPackageDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]
	packageID := parts[1]

	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	packageService := wskClient.Packages

	_, err = packageService.Delete(packageID)
	if err != nil {
		return fmt.Errorf("Error deleting IBM Cloud Function Package: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMFunctionPackageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return false, err
	}

	namespace := ""
	packageID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		packageID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		packageID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, packageID))
	}

	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return false, err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return false, err

	}

	packageService := wskClient.Packages

	pkg, resp, err := packageService.Get(packageID)
	if err != nil {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error communicating with IBM Cloud Function Client : %s", err)
	}

	return pkg.Name == packageID, nil
}
