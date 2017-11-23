# Download / Upload a static 'fixed' Resource

_You can use this concourse-ci resource to download resources which you want to explicitly fetch using a specific version._

Upfront, static means, that you consume one exact file, a specific version of something. It should not change and it not designed to do so.
Same goes from upload, it should upload into the same bucket, but it should rather be used as initial upload, not as an "replace".
E.g. you can upload to a Sonatype nexus raw repository.

It's not like a usual/classic concourse-ci resource where you want to utilize `check` to pull in the most recent release of the packages
named, but rather you want to pick on specific release.
This can be seen as a simimlar dependency as using a `package.json` or a `pom.xml` where you pick a specific version you know you are compatible with.

In addition, with `out`, you can upload a file into an endpoint which understands multipart upload.
## Deeper thoughts on why and what we use static downloads

The reason to not use a `curl/wget` inside your task instead or build script  (instead of using this resource) is that you want the resource/package to be interchangeable in the pipeline**s** - that means, you want to

 a) sometimes pick a specific version - the stable-build-pipeline
 b) you want to test-integrate the most recent build / release - the bloody-edge-pipeline
 c) <your special case here>

That said, we use this resource with a simple rule to apply, which ensures that b) actually can be achieved:

**The dependency you download with static-download should be build in one of your (other) pipelines and released from there - it should not be an external dependency**

The reason behind that is, that to achieve b) we want to reference the result of this pipeline with `passed` or its latest `semver` in our `bloody-edge-pipeine` build.

Lets make this a example case.

So we have a `stable-build-pipeline`

```yaml
jobs:
  - name: build-stable
    # this get we want to be interchangeable with latest/from SCM resources
    - get: static-download-our-gem
    - task: build-app-with-gem
      file: ....
      input_mapping:
        docker-sync: static-download-our-gem
resources:
  - name: static-download-our-gem
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync-0.5.0.gem

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

 **So the concept of resources is still ensured and it is quiet important to have something like static-downloads for specific builds**

## Source Configuration

* `uri`: *Required.* The location of the file to download - we download using a simple curl request
* `authentication`: *Optional.* Your basic-auth data `authentication.user` and `authentication.password`
* `extract`: *Optional.* if `true`, `gunzip| tar xf` will be used to extract the download
* `skip_ssl_validation`: *Optional.* Skip SSL validation.

```yaml
resources:
  - name: static-our-gem
    type: static
    source:
      uri: https://rubygems.org/downloads/docker-sync-0.5.0.gem
      authentication:
        user: eugenmayer
        password: verysecret
```
## Behavior

### `check`: Not implemented.

Is not implemented but does work in pipelines using a pseudo version. **It does never use `trigger: true` on it!**

### `in`: Download and extract the archive.

Fetches a URL and if `extract` is true also unpacks it

### `out`: Not implemented.

Currently this uses a multipart upload of your file. The filename should be yet part of the URL itself like `https://myendpoint/thisfile.tgz`, the upload is implemented using

`curl --upload-file <filepath>`

It uses basic authentication, if set.

It can for example be used to upload into a Sonatype nexus repository of type `raw`, when you need `docker` or `maven`, be sure to use the corresponding resources for those

## Examples

```yaml
jobs:
  - name: build-stable
    - get: static-download-our-gem
    - task: build-app-with-gem
      file: ....

# define the resource we want to consume      
resources:
  - name: static-download-our-gem
    type: static-download
    source:
      uri: https://rubygems.org/downloads/docker-sync-0.5.0.gem

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
