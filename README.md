# maven-stage-resource
Concourse resource to copy maven resources from one repository to another.

## Source Configuration

* `source_url`: *Required*. The source maven repository location.
* `target_url`: *Required*. The target maven repository location.
* `artifact`: *Required*. The artifact coordinates in the form of groupId:artifactId:type\[:classifier\]
* `username`: *Optional* The username used to authenticate.
* `password`: *Optional*. The password used to authenticate.
* `verbose`: *Optional*. True to write intensive log.

## Check

Checks for new versions of the artifact by retrieving the maven-metadata.xml from the source repository.

## Get

Download the source artifact from repositry and puts it to target.

* `version`: *Optional* Make abbilty to provide version instead of getting the latest.

## Put

Undefined.

## Pipeline example

```yaml
---
resource_types:
  - name: maven-stage-resource
    type: docker-image
    source:
      repository: lorands/maven-stage-resource
resources:
  - name: stage-to-uat
    type: maven-stage-resource
    source:
      source_url: https://mynexus.example.com/repository/develop
      target_url: https://mynexus.example.com/repository/uat
      username: myUser
      password: myPass
jobs:
  - name: merge-dev-to-uat
    plan:
    - get: dev-artifact

```





