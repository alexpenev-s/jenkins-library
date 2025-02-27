// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type kubernetesDeployOptions struct {
	AdditionalParameters       []string `json:"additionalParameters,omitempty"`
	APIServer                  string   `json:"apiServer,omitempty"`
	AppTemplate                string   `json:"appTemplate,omitempty"`
	ChartPath                  string   `json:"chartPath,omitempty"`
	ContainerRegistryPassword  string   `json:"containerRegistryPassword,omitempty"`
	ContainerRegistryURL       string   `json:"containerRegistryUrl,omitempty"`
	ContainerRegistryUser      string   `json:"containerRegistryUser,omitempty"`
	ContainerRegistrySecret    string   `json:"containerRegistrySecret,omitempty"`
	CreateDockerRegistrySecret bool     `json:"createDockerRegistrySecret,omitempty"`
	DeploymentName             string   `json:"deploymentName,omitempty"`
	DeployTool                 string   `json:"deployTool,omitempty"`
	ForceUpdates               bool     `json:"forceUpdates,omitempty"`
	HelmDeployWaitSeconds      int      `json:"helmDeployWaitSeconds,omitempty"`
	HelmValues                 []string `json:"helmValues,omitempty"`
	Image                      string   `json:"image,omitempty"`
	IngressHosts               []string `json:"ingressHosts,omitempty"`
	KeepFailedDeployments      bool     `json:"keepFailedDeployments,omitempty"`
	KubeConfig                 string   `json:"kubeConfig,omitempty"`
	KubeContext                string   `json:"kubeContext,omitempty"`
	KubeToken                  string   `json:"kubeToken,omitempty"`
	Namespace                  string   `json:"namespace,omitempty"`
	TillerNamespace            string   `json:"tillerNamespace,omitempty"`
	DockerConfigJSON           string   `json:"dockerConfigJSON,omitempty"`
}

// KubernetesDeployCommand Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.
func KubernetesDeployCommand() *cobra.Command {
	const STEP_NAME = "kubernetesDeploy"

	metadata := kubernetesDeployMetadata()
	var stepConfig kubernetesDeployOptions
	var startTime time.Time
	var logCollector *log.CollectorHook

	var createKubernetesDeployCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.",
		Long: `Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.

!!! note "Deployment supports multiple deployment tools"
    Currently the following are supported:

    * [Helm](https://helm.sh/) command line tool and [Helm Charts](https://docs.helm.sh/developing_charts/#charts).
    * [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and ` + "`" + `kubectl apply` + "`" + ` command.

## Helm
Following helm command will be executed by default:

` + "`" + `` + "`" + `` + "`" + `
helm upgrade <deploymentName> <chartPath> --install --force --namespace <namespace> --wait --timeout <helmDeployWaitSeconds> --set "image.repository=<yourRegistry>/<yourImageName>,image.tag=<yourImageTag>,secret.dockerconfigjson=<dockerSecret>,ingress.hosts[0]=<ingressHosts[0]>,,ingress.hosts[1]=<ingressHosts[1]>,...
` + "`" + `` + "`" + `` + "`" + `

* ` + "`" + `yourRegistry` + "`" + ` will be retrieved from ` + "`" + `containerRegistryUrl` + "`" + `
* ` + "`" + `yourImageName` + "`" + `, ` + "`" + `yourImageTag` + "`" + ` will be retrieved from ` + "`" + `image` + "`" + `
* ` + "`" + `dockerSecret` + "`" + ` will be calculated with a call to ` + "`" + `kubectl create secret docker-registry regsecret --docker-server=<yourRegistry> --docker-username=<containerRegistryUser> --docker-password=<containerRegistryPassword> --dry-run=true --output=json'` + "`" + ``,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.ContainerRegistryPassword)
			log.RegisterSecret(stepConfig.ContainerRegistryUser)
			log.RegisterSecret(stepConfig.KubeConfig)
			log.RegisterSecret(stepConfig.KubeToken)
			log.RegisterSecret(stepConfig.DockerConfigJSON)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetryData.ErrorCategory = log.GetErrorCategory().String()
				telemetry.Send(&telemetryData)
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunk.Send(&telemetryData, logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunk.Initialize(GeneralConfig.CorrelationID,
					GeneralConfig.HookConfig.SplunkConfig.Dsn,
					GeneralConfig.HookConfig.SplunkConfig.Token,
					GeneralConfig.HookConfig.SplunkConfig.Index,
					GeneralConfig.HookConfig.SplunkConfig.SendLogs)
			}
			kubernetesDeploy(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addKubernetesDeployFlags(createKubernetesDeployCmd, &stepConfig)
	return createKubernetesDeployCmd
}

func addKubernetesDeployFlags(cmd *cobra.Command, stepConfig *kubernetesDeployOptions) {
	cmd.Flags().StringSliceVar(&stepConfig.AdditionalParameters, "additionalParameters", []string{}, "Defines additional parameters for \"helm install\" or \"kubectl apply\" command.")
	cmd.Flags().StringVar(&stepConfig.APIServer, "apiServer", os.Getenv("PIPER_apiServer"), "Defines the Url of the API Server of the Kubernetes cluster.")
	cmd.Flags().StringVar(&stepConfig.AppTemplate, "appTemplate", os.Getenv("PIPER_appTemplate"), "Defines the filename for the kubernetes app template (e.g. k8s_apptemplate.yaml)")
	cmd.Flags().StringVar(&stepConfig.ChartPath, "chartPath", os.Getenv("PIPER_chartPath"), "Defines the chart path for deployments using helm. It is a mandatory parameter when `deployTool:helm` or `deployTool:helm3`.")
	cmd.Flags().StringVar(&stepConfig.ContainerRegistryPassword, "containerRegistryPassword", os.Getenv("PIPER_containerRegistryPassword"), "Password for container registry access - typically provided by the CI/CD environment.")
	cmd.Flags().StringVar(&stepConfig.ContainerRegistryURL, "containerRegistryUrl", os.Getenv("PIPER_containerRegistryUrl"), "http(s) url of the Container registry where the image to deploy is located.")
	cmd.Flags().StringVar(&stepConfig.ContainerRegistryUser, "containerRegistryUser", os.Getenv("PIPER_containerRegistryUser"), "Username for container registry access - typically provided by the CI/CD environment.")
	cmd.Flags().StringVar(&stepConfig.ContainerRegistrySecret, "containerRegistrySecret", `regsecret`, "Name of the container registry secret used for pulling containers from the registry.")
	cmd.Flags().BoolVar(&stepConfig.CreateDockerRegistrySecret, "createDockerRegistrySecret", false, "Only for `deployTool:kubectl`: Toggle to turn on `containerRegistrySecret` creation.")
	cmd.Flags().StringVar(&stepConfig.DeploymentName, "deploymentName", os.Getenv("PIPER_deploymentName"), "Defines the name of the deployment. It is a mandatory parameter when `deployTool:helm` or `deployTool:helm3`.")
	cmd.Flags().StringVar(&stepConfig.DeployTool, "deployTool", `kubectl`, "Defines the tool which should be used for deployment.")
	cmd.Flags().BoolVar(&stepConfig.ForceUpdates, "forceUpdates", true, "Helm only: force resource updates with helm parameter `--force`")
	cmd.Flags().IntVar(&stepConfig.HelmDeployWaitSeconds, "helmDeployWaitSeconds", 300, "Number of seconds before helm deploy returns.")
	cmd.Flags().StringSliceVar(&stepConfig.HelmValues, "helmValues", []string{}, "List of helm values as YAML file reference or URL (as per helm parameter description for `-f` / `--values`)")
	cmd.Flags().StringVar(&stepConfig.Image, "image", os.Getenv("PIPER_image"), "Full name of the image to be deployed.")
	cmd.Flags().StringSliceVar(&stepConfig.IngressHosts, "ingressHosts", []string{}, "(Deprecated) List of ingress hosts to be exposed via helm deployment.")
	cmd.Flags().BoolVar(&stepConfig.KeepFailedDeployments, "keepFailedDeployments", false, "Defines whether a failed deployment will be purged")
	cmd.Flags().StringVar(&stepConfig.KubeConfig, "kubeConfig", os.Getenv("PIPER_kubeConfig"), "Defines the path to the \"kubeconfig\" file.")
	cmd.Flags().StringVar(&stepConfig.KubeContext, "kubeContext", os.Getenv("PIPER_kubeContext"), "Defines the context to use from the \"kubeconfig\" file.")
	cmd.Flags().StringVar(&stepConfig.KubeToken, "kubeToken", os.Getenv("PIPER_kubeToken"), "Contains the id_token used by kubectl for authentication. Consider using kubeConfig parameter instead.")
	cmd.Flags().StringVar(&stepConfig.Namespace, "namespace", `default`, "Defines the target Kubernetes namespace for the deployment.")
	cmd.Flags().StringVar(&stepConfig.TillerNamespace, "tillerNamespace", os.Getenv("PIPER_tillerNamespace"), "Defines optional tiller namespace for deployments using helm.")
	cmd.Flags().StringVar(&stepConfig.DockerConfigJSON, "dockerConfigJSON", os.Getenv("PIPER_dockerConfigJSON"), "Path to the file `.docker/config.json` - this is typically provided by your CI/CD system. You can find more details about the Docker credentials in the [Docker documentation](https://docs.docker.com/engine/reference/commandline/login/).")

	cmd.MarkFlagRequired("containerRegistryUrl")
	cmd.MarkFlagRequired("deployTool")
	cmd.MarkFlagRequired("image")
}

// retrieve step metadata
func kubernetesDeployMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "kubernetesDeploy",
			Aliases:     []config.Alias{{Name: "deployToKubernetes", Deprecated: true}},
			Description: "Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "kubeConfigFileCredentialsId", Description: "Jenkins 'Secret file' credentials ID containing kubeconfig file. Details can be found in the [Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).", Type: "jenkins", Aliases: []config.Alias{{Name: "kubeCredentialsId", Deprecated: true}}},
					{Name: "kubeTokenCredentialsId", Description: "Jenkins 'Secret text' credentials ID containing token to authenticate to Kubernetes. This is an alternative way to using a kubeconfig file. Details can be found in the [Kubernetes documentation](https://kubernetes.io/docs/reference/access-authn-authz/authentication/).", Type: "jenkins", Aliases: []config.Alias{{Name: "k8sTokenCredentialsId", Deprecated: true}}},
					{Name: "dockerCredentialsId", Type: "jenkins"},
					{Name: "dockerConfigJsonCredentialsId", Description: "Jenkins 'Secret file' credentials ID containing Docker config.json (with registry credential(s)).", Type: "jenkins"},
				},
				Resources: []config.StepResources{
					{Name: "deployDescriptor", Type: "stash"},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "additionalParameters",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmDeploymentParameters"}},
						Default:     []string{},
					},
					{
						Name:        "apiServer",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "k8sAPIServer"}},
						Default:     os.Getenv("PIPER_apiServer"),
					},
					{
						Name:        "appTemplate",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "k8sAppTemplate"}},
						Default:     os.Getenv("PIPER_appTemplate"),
					},
					{
						Name:        "chartPath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmChartPath"}},
						Default:     os.Getenv("PIPER_chartPath"),
					},
					{
						Name: "containerRegistryPassword",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "dockerCredentialsId",
								Param: "password",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_containerRegistryPassword"),
					},
					{
						Name: "containerRegistryUrl",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "container/registryUrl",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "dockerRegistryUrl"}},
						Default:   os.Getenv("PIPER_containerRegistryUrl"),
					},
					{
						Name: "containerRegistryUser",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "dockerCredentialsId",
								Param: "username",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_containerRegistryUser"),
					},
					{
						Name:        "containerRegistrySecret",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `regsecret`,
					},
					{
						Name:        "createDockerRegistrySecret",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "deploymentName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmDeploymentName"}},
						Default:     os.Getenv("PIPER_deploymentName"),
					},
					{
						Name:        "deployTool",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     `kubectl`,
					},
					{
						Name:        "forceUpdates",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     true,
					},
					{
						Name:        "helmDeployWaitSeconds",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "int",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     300,
					},
					{
						Name:        "helmValues",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name: "image",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "container/imageNameTag",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "deployImage"}},
						Default:   os.Getenv("PIPER_image"),
					},
					{
						Name:        "ingressHosts",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name:        "keepFailedDeployments",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name: "kubeConfig",
						ResourceRef: []config.ResourceReference{
							{
								Name: "kubeConfigFileCredentialsId",
								Type: "secret",
							},

							{
								Name:    "kubeConfigFileSecretName",
								Type:    "vaultSecretFile",
								Default: "kube-config",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_kubeConfig"),
					},
					{
						Name:        "kubeContext",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_kubeContext"),
					},
					{
						Name: "kubeToken",
						ResourceRef: []config.ResourceReference{
							{
								Name: "kubeTokenCredentialsId",
								Type: "secret",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_kubeToken"),
					},
					{
						Name:        "namespace",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmDeploymentNamespace"}, {Name: "k8sDeploymentNamespace"}},
						Default:     `default`,
					},
					{
						Name:        "tillerNamespace",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "helmTillerNamespace"}},
						Default:     os.Getenv("PIPER_tillerNamespace"),
					},
					{
						Name: "dockerConfigJSON",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/dockerConfigJSON",
							},

							{
								Name: "dockerConfigJsonCredentialsId",
								Type: "secret",
							},

							{
								Name: "dockerConfigFileVaultSecretName",
								Type: "vaultSecretFile",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_dockerConfigJSON"),
					},
				},
			},
			Containers: []config.Container{
				{Image: "dtzar/helm-kubectl:3.4.1", WorkingDir: "/config", Options: []config.Option{{Name: "-u", Value: "0"}}, Conditions: []config.Condition{{ConditionRef: "strings-equal", Params: []config.Param{{Name: "deployTool", Value: "helm3"}}}}},
				{Image: "dtzar/helm-kubectl:2.17.0", WorkingDir: "/config", Options: []config.Option{{Name: "-u", Value: "0"}}, Conditions: []config.Condition{{ConditionRef: "strings-equal", Params: []config.Param{{Name: "deployTool", Value: "helm"}}}}},
				{Image: "dtzar/helm-kubectl:2.17.0", WorkingDir: "/config", Options: []config.Option{{Name: "-u", Value: "0"}}, Conditions: []config.Condition{{ConditionRef: "strings-equal", Params: []config.Param{{Name: "deployTool", Value: "kubectl"}}}}},
			},
		},
	}
	return theMetaData
}
