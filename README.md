# Download / Upload a static 'fixed' Resource

_You can use this concourse-ci resource to download resources which you want to explicitly fetch using a specific version._

Upfront, static means, that you consume `in` one exact file, a specific version of something. It should not change and it not designed to do so.
The `out` operation is ment to multipart upload something into an endpoint, where a usual resource integration does not make sense.
E.g. you can upload to a Sonatype nexus raw repository, where versioning an upload is very hard to implement.
In contrary when you want to deploy a jar into the very same nexus, you rather want to use the `maven-resource` instead.

It's not like a usual/classic concourse-ci resource where you want to utilize `check` to pull in the most recent release of the packages
named, but rather you want to pick on specific release.
This can be seen as a similar dependency as using a `package.json` or a `pom.xml` where you pick a specific version you know you are compatible with.

In addition, with `out`, you can upload a file into an endpoint which understands multipart upload.

## Deeper thoughts on why and what we use static downloads

The reason to not use a `curl/wget` inside your task instead or build script  (instead of using this resource) is that you want the resource/package to be interchangeable in the pipeline**s** - that means, you want to

 a) sometimes pick a specific version - the stable-build-pipeline
 b) you want to test-integrate the most recent build / release - the bloody-edge-pipeline
 c) <your special case here>

That said, we use this resource with a simple rule to apply, which ensures that b) actually can be achieved:

**The dependency you download with static-download should be build in one of your (other) pipelines and released from there - it should not be an external dependency**

The reason behind that is, that to achieve b) we want to reference the result of this pipeline with `passed` or its latest `semver` in our `bloody-edge-pipeine` build.

Lets make create a example

So we have a `stable-build-pipeline`

```yaml
jobs:
  - name: build-stable
    # this get we want to be interchangeable with latest/from SCM resources
    - get: static-our-gem
    - task: build-app-with-gem
      file: ....
      input_mapping:
        docker-sync: static-download-our-gem
resources:
  - name: static-our-gem
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync-0.5.0.gem
      version_static: 0.5.0

resource_types:
  - name:                       static
    type:                       docker-image
    source:
      repository:               eugenmayer/concourse-static-resource
      tag:                      latest
```

and now our `bloody-edge-pipeline`

```yaml
jobs:
  - name: build-stable
    # now we replaced the gem download with a s3 latest version fetch
    - get: staging-artifact-storage
    - task: build-app-with-gem
      file: ....
      input_mapping:
        docker-sync: staging-artifact-storage
resources:
    - name: staging-artifact-storage
      type: s3
      source:
        bucket: ci-artifacts
        versioned_file: staging/artifacts/docker-sync/docker-sync.gem
        access_key_id: ((aws_s3_artifacts.user))
        secret_access_key: ((aws_s3_artifacts.password))
        region_name: eu-central-1
        private: true
 ```

And as you see, we can use the same task, just with a different resource, this time from "the latest staging release" from s3

Now even the `build-from-source-pipeline`

You can do the same with `get: from_scm` then do a intermidiate `task: build-gem` and then put the result back into
 `- task: build-app-with-gem` - lets do that.

 ```yaml
 jobs:
   - name: build-stable
     # get the gem source
     - get: source_scm
     # build the gem from source
     - task: build-gem
       file: build-gem-using-gem-build
       output_mapping:
         artifact: gem_build_from_scm # thats our gem
     # now use the very same task again, consuming the gem build from source this time
     - task: build-app-with-gem
       input_mapping:
         docker-sync: gem_build_from_scm
       file: ....
 resources:
     - name: source_scm_for_docker_sync_gem
       type: git
       source:
         uri: https://github.com/EugenMayer/docker-sync

  ```

**So the concept of resources is still ensured and it is quiet important to have something like static-resources for specific builds**

## Source Configuration

* `uri`: *Required.* The location of the file to download - we download using a simple curl request
* `version_static`: *Required.* The version of your downloaded artefact - we do not parse it from URL since its not always possible
  it will be use to replace the placeholder `<version>` if it is provided
* `authentication`: *Optional.* Your basic-auth data `authentication.user` and `authentication.password`
* `extract`: *Optional.* if `true`, `gunzip| tar xf` will be used to extract the download
* `skip_ssl_validation`: *Optional.* Skip SSL validation.

```yaml
resources:
  - name: static-our-gem
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync-<version>.gem
      version_static: 0.5.0
      authentication:
        user: eugenmayer
        password: verysecret
```

You can also use a static URL without a placeholder, e.g. if you have no version string in the URL anyway

```yaml
resources:
  - name: static-our-gem
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync.gem
      # you still have to set this to make sure we know, which version that is
      version_static: 0.5.0
      authentication:
        user: eugenmayer
        password: verysecret
```
## Behavior

### `check`: Pseudo checks for a new version, always returning the static version you provided

It will always return the value you have set in `version_static` as a new version to keep the `check` to `in` handover
properly working.

Do never use this with `trigger: true`

### `in`: Download and extract the archive.

Fetches the URL and if `extract` is true also unpacks it

### `out`: Uploads your file with the (dynamic) version you give using curl multipart to any endpoint

Accepts `params`:

 - `source_filepath`: **required** the filepath to the source file from your resources. You can use a glob here like input/artifact-*.gz
 - `version_filepath`: *optional* A filepath to a file including the version. Needs to be provided if `URI` includes a placeholder (<version>

Currently this uses a multipart upload of your file using curl utilizing

`curl --upload-file <filepath>`

It uses basic authentication, if set.

It can for example be used to upload into a Sonatype nexus repository of type `raw`, when you need `docker` or `maven`, be sure to use the corresponding resources for those

## Examples

This example illustrates how we download docker-sync-0.5.0.gem from an arbitrary source, build whatever app we want
and then upload that app back to our nexus `raw` repository using our `out`

```yaml
jobs:
  - name: build-stable
    - get: some-semver-version
    - get: static-download-our-gem
    - task: build-app-with-gem
      file: ....
      output_mapping:
        artifact_gem: our-app

    # this will upload the tgz provided by the build-app-with-gem to our nexus
    # using the filename myapp-0.0.1.tgz if some-semver-version/number is 0.0.1
    - put: upload-to-nexus
      params:
        source_filepath: our-app/app-*.tgz
        version_filepath: some-semver-version/number

# define the resource we want to consume      
resources:
  - name: upload-to-nexus
    type: static
    source:
      uri: https://nexus.sonatype.com/repository/raw-artifacts/docker-sync/myapp-<version>.tgz
      version_static: 0.5.0

  - name: static-our-gem
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync-<version>.gem
      version_static: 0.5.0

  - name: static-no-version-url
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync.gem
      version_static: 0.5.0

# we need this to register the custom resource_type
resource_types:
  - name:                       static
    type:                       docker-image
    source:
      repository:               eugenmayer/concourse-static-resource
      tag:                      latest
```

## Docker hub

You find this image at https://hub.docker.com/r/eugenmayer/concourse-static-resource/

## Development

### Running the tests

Just run

```
docker build .
```

The build runs the integrations tests in a multi-stag build

### Contributing

Sure, just go for it - happy to merge pull requests.
