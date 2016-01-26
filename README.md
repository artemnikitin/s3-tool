# s3-tool
CLI tool for AWS S3, written in Go

##### Dependency

Only depends on AWS SDK. Install it via    
```
go get github.com/aws/aws-sdk-go/...
```

##### AWS Credentials

Currently assumes that you will have credentials settled as environmental variables.   
```
export AWS_ACCESS_KEY_ID=<key>
export AWS_SECRET_ACCESS_KEY=<secret>
```
##### Running
Get it via    
``` 
go get github.com/artemnikitin/s3-tool 
``` 

You can specify parameter ```-log=true``` for logging AWS requests and responses.

##### TODO  
1. Alternative ways to authenticate in AWS