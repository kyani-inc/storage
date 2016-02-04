# DynamoDB Storage provider
This is a simple abstraction layer between Amazon's DynamoDB and our local storage interface.

### Configuration
To run the tests, create a .env file in the dynamodb directory with the following vars:

```
AWS_ACCESS=my-access-id
AWS_SECRET=my-access-secret
AWS_REGION=us-east-1
DYNAMO_DB_TABLE=my_test_table
DYNAMO_DB_ENDPOINT=https://dynamodb.us-east-1.amazonaws.com
```

### Running Local
Amazon provides a convenient java executable for running a mock DynamoDB server locally. Read all
about it here:
http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Tools.DynamoDBLocal.html.
