import logging
import boto3
from botocore.exceptions import ClientError
from botocore.config import Config
from config import config

def create_presigned_url(bucket_name, object_name, expiration=3600):
    s3_client = boto3.client('s3',
        aws_access_key_id=config.AWS_ACCESS_KEY_ID,
        aws_secret_access_key=config.AWS_SECRET_ACCESS_KEY,
        region_name=config.AWS_DEFAULT_REGION,
        config=Config(signature_version='s3v4')
    )

    try:
        response = s3_client.generate_presigned_url(
            'get_object',
            Params={'Bucket': bucket_name, 'Key': object_name},
            ExpiresIn=expiration,
        )
    except ClientError as e:
        return None
    
    
    return response

# upload a file to s3
def upload_file_to_s3(bucket_name, file_path, object_name):
    s3_client = boto3.client('s3',
        aws_access_key_id=config.AWS_ACCESS_KEY_ID,
        aws_secret_access_key=config.AWS_SECRET_ACCESS_KEY,
        region_name=config.AWS_DEFAULT_REGION,
    )
    
    
    s3_client.upload_file(file_path, bucket_name, object_name)
    
