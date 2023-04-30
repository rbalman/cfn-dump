### Extract CFN Stack Dependency
Extract CloudFormation stacks dependencies based on stack's exports and imports into a json file.  You can provide stack name pattern using --pattern flag


**Sample Output**

Above command will create json file in the working directory with these naming convention:
  -  `cfn-dependency.json` when stack prefix is provided
  - `global-cfn-dependency.json` when stack prefix is empty

```json
{
  "rbalman-dev-rds-network": { //StackName and its list of exports
    "Exports": {
      "balman:rds:SecurityGroupId": [ //Export's name and stack list that imports it
        "rbalman-rds-db"
      ],
      "balman:rds:Subnet": [
        "rbalman-rds-subnet"
      ]
    }
  },
  "rbalman-rds-db": {
    "Exports": {
      "rbalman:RDSEndpoint": null, //Import list is null as no stack has imported this export
      "rbalman:RDSEndpoint": [
        "rbalman-token-management-service"
      ]
    }
  }
}
```


**Dependency Document Format**

**NOTE:** CFN stacks with zero exports will not be listed in the dependency json file.
```json
{
  "<stack-name1>": {
    "Exports": {
      "<first-export1>": null, //Import list is null as no stack has imported this export
      "<second-export2>": [
        "<first-importing-stack>",
        "<second-importing-stack>"
      ]
    }
  },
  "<stack-name2>": {
    "Exports": {
      "<first-export1>": null, //Import list is null as no stack has imported this export
      "<second-export2>": [
        "<first-importing-stack>",
        "<second-importing-stack>"
      ]
    }
  }
}
```


**Syntax**
```
export AWS_PROFILE=dev AWS_REGION=us-east-1
cfnd dump --pattern "dev-iam"
```
