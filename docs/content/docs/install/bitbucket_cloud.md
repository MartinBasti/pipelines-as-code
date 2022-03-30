---
title: Bitbucket Cloud
weight: 14
---
# Install Pipelines-As-Code for Bitbucket Cloud

Pipelines-As-Code has a full support on Bitbucket Cloud on
<https://bitbucket.org> as Webhook.

Following the [infrastructure installation](install.md#install-pipelines-as-code-infrastructure) :

* You will have to generate an app password for Pipelines-as-Code Bitbucket API
  operations. Follow this guide to create an app password :

<https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/>

Add those permissions to the token :

[image](https://user-images.githubusercontent.com/98980/154526912-75c52ded-45e9-42d4-8c09-908b86eb57b4.png)

Make sure you note somewhere the generated token or otherwise you will have to
recreate it.

* Go to you **"Repository setting"** tab on your **Repository** and click on the
  **WebHooks** tab and **"Add webhook"** button.

* Set a **Title** (i.e: Pipelines as Code)

* Set the URL to the event listener public URL. On OpenShift you can get the public URL of the Pipelines-as-Code
  controller like this :

  ```shell
  echo https://$(oc get route -n pipelines-as-code pipelines-as-code-controller -o jsonpath='{.spec.host}')
  ```

* [Refer to this screenshot](/images/bitbucket-cloud-create-webhook.png) on how to configure the Webhook. The
  individual events to select are :
  * Repository -> Push
  * Pull Request -> Created
  * Pull Request -> Updated
  * Pull Request -> Comment created
  * Pull Request -> Comment updated

* You are now able to create a Repository CRD. The repository CRD will have a Secret and Username that contains the App
  Password as generated and Pipelines as Code will know how to use it for Bitbucket API operations.

  * First create the secret with the app password in the `target-namespace` :

  ```shell
  kubectl -n target-namespace create secret generic bitbucket-cloud-token \
          --from-literal token="TOKEN_AS_GENERATED_PREVIOUSLY"
  ```

  * And now create Repository CRD with the secret field referencing it.

  * Here is an example of a Repository CRD :

```yaml
---
apiVersion: "pipelinesascode.tekton.dev/v1alpha1"
kind: Repository
metadata:
  name: my-repo
  namespace: target-namespace
spec:
  url: "https://bitbucket.com/workspace/repo"
  branch: "main"
  git_provider:
    user: "yourbitbucketusername"
    secret:
      name: "bitbucket-cloud-token"
      # Set this if you have a different key in your secret
      # key: "token"
```

## Bitbucket Cloud Notes

* `git_provider.secret` cannot reference a secret in another namespace,
  Pipelines as code assumes always it will be the same namespace as where the
  repository has been created.

* `tkn-pac create` and `bootstrap` is not supported on Bitbucket Server.

{{< hint info >}}
You can only reference user by `ACCOUNT_ID` in owner file, see here for the
reasoning :

<https://developer.atlassian.com/cloud/bitbucket/bitbucket-api-changes-gdpr/#introducing-atlassian-account-id-and-nicknames>
{{< /hint >}}

{{< hint danger >}}
* There is no Webhook secret support in Bitbucket Cloud. To be able to secure
  the payload and not let a user hijack the CI, Pipelines-as-Code will fetch the
  ip addresses list from <https://ip-ranges.atlassian.com/> and make sure the
  webhook only comes from the Bitbucket Cloud IPS.
* If you want to add some ips address or networks you can add them to the
  key **bitbucket-cloud-additional-source-ip** in the pipelines-as-code
  configmap in the pipelines-as-code namespace.  You can added multiple
  network or ips separated by a comma.

* If you want to disable this behavior you can set the key
  **bitbucket-cloud-check-source-ip** to false in the pipelines-as-code
  configmap in the pipelines-as-code namespace.
{{< /hint >}}