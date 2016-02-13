# s3-tool
[![Go Report Card](https://goreportcard.com/badge/github.com/artemnikitin/s3-tool)](https://goreportcard.com/report/github.com/artemnikitin/s3-tool)    
CLI tool for AWS S3, written in Go. Work in progress...

##### AWS Credentials

Set environment variables     
```
export AWS_ACCESS_KEY_ID=<key>       
export AWS_SECRET_ACCESS_KEY=<secret>
```     

##### Install
Get it via    
``` 
go get github.com/artemnikitin/s3-tool 
``` 

##### Commands
Current list of supported commands:
```
presigned
download
```

- ```presigned``` generate pre-signed URL for downloading file from S3.   
Requires parameters:
    - ```bucket``` specified bucket in S3
    - ```key``` specified key in S3      
Example:   
```
s3tool presigned -bucket=mybucket -key=my-file.png
```

- ```download``` download file from S3.   
Requires parameters:
    - ```bucket``` specified bucket in S3
    - ```key``` specified key in S3  
    - ```url``` if specified, then ```s3-tool``` will try to download file be pre-signed URL ignoring other parameters   
Example:   
```
s3-tool download -bucket=mybucket -key=my-file.png    
s3-tool download -url=https://
```

##### Optional parameters
- ```region``` set S3 region, by default region will be set to ```us-east-1```       
Example:    
``` 
s3-tool -bucket=bucket-name -key=my-key -region=region-name 
```    

##### Run
Run it like:   
```
s3-tool command -parameter=blabla ...
```

You can specify parameter ```-log=true``` for logging AWS requests and responses.

##### TODO  
1. Integrate s3-uploader as command(s)
2. Add more commands
3. Alternative ways to authenticate in AWS
