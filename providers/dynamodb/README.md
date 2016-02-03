# DynamoDB Storage provider
This is a simple abstraction layer between Amazon's DynamoDB and our local storage interface.

### Configuration
To run the tests, create a .env file in the dynamodb directory with the following vars:

```
AWS_REGION=us-west-2
DYNAMO_DB_TABLE=test_table
DYNAMO_DB_ENDPOINT=dynamodb.serveraddress.aws.com
```

### Running Local
Amazon provides a convenient java executable for running a mock DynamoDB server locally. Read all
about it here:
http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Tools.DynamoDBLocal.html.
