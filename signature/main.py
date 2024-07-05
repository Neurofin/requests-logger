import os
import shutil
import boto3
from detect import pdf_to_png, load_and_extract_features, compute_similarity_matrix, similarity_to_distance_matrix, cluster_signatures, save_clusters
from dotenv import load_dotenv
load_dotenv()
def load_s3_uris(file_path):
    if os.path.exists(file_path):
        with open(file_path, 'r') as file:
            return [line.strip() for line in file.readlines()]
    return []

def save_s3_uris(file_path, uris):
    with open(file_path, 'w') as file:
        for uri in uris:
            file.write(f"{uri}\n")

def download_pdf_from_s3(s3_uri, download_path, aws_access_key_id, aws_secret_access_key, region_name):
    s3_client = boto3.client(
        's3',
        aws_access_key_id=aws_access_key_id,
        aws_secret_access_key=aws_secret_access_key,
        region_name=region_name
    )
    bucket_name, key = s3_uri.replace("s3://", "").split("/", 1)
    local_path = os.path.join(download_path, os.path.basename(key))
    os.makedirs(os.path.dirname(local_path), exist_ok=True)
    s3_client.download_file(bucket_name, key, local_path)
    print(f"Downloaded {s3_uri} to {local_path}")
    return local_path

def upload_directory_to_s3(directory_path, bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name):
    s3_client = boto3.client(
        's3',
        aws_access_key_id=aws_access_key_id,
        aws_secret_access_key=aws_secret_access_key,
        region_name=region_name
    )
    uris = []
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            local_path = os.path.join(root, file)
            relative_path = os.path.relpath(local_path, directory_path)
            s3_path = os.path.join(s3_folder, relative_path)
            s3_client.upload_file(local_path, bucket_name, s3_path)
            s3_uri = f"s3://{bucket_name}/{s3_path}"
            uris.append(s3_uri)
            print(f"Uploaded {local_path} to {s3_uri}")
    return uris

def process_pdf_from_s3(s3_pdf_uri):
    # Define paths
    temp_dir = 'temp'
    final_results_dir = 'temp/final_results'
    clustering_results_dir = 'temp/clustering_results'
    bucket_name = 'clustering-results123'
    s3_uris_file = 'temp/s3_uris.txt'

    # AWS credentials (retrieved from environment variables)
    aws_access_key_id = os.getenv('AWS_ACCESS_KEY_ID')
    aws_secret_access_key = os.getenv('AWS_SECRET_ACCESS_KEY')
    region_name = os.getenv('AWS_REGION')

    # Ensure temporary directory exists
    os.makedirs(temp_dir, exist_ok=True)

    # Download PDF from S3
    pdf_path = download_pdf_from_s3(s3_pdf_uri, temp_dir, aws_access_key_id, aws_secret_access_key, region_name)

    # Perform PDF to PNG conversion and signature extraction
    pdf_to_png(pdf_path, temp_dir, final_results_dir)

    # Extract the base name of the PDF without extension to match the output folder name
    pdf_base_name = os.path.splitext(os.path.basename(pdf_path))[0]

    # Define the input folder for clustering
    input_folder = os.path.join(final_results_dir, pdf_base_name)

    if not os.path.exists(clustering_results_dir):
        os.makedirs(clustering_results_dir)

    # Extract input folder name
    input_folder_name = os.path.basename(os.path.normpath(input_folder))

    # Load images and extract features
    images, filenames, features = load_and_extract_features(input_folder)

    # Check if features are extracted
    if len(features) == 0:
        print("No features extracted. Exiting.")
        return []
    else:
        # Compute similarity matrix
        similarity_matrix = compute_similarity_matrix(features)

        # Convert similarity matrix to distance matrix
        distance_matrix = similarity_to_distance_matrix(similarity_matrix)

        # Cluster signatures
        labels = cluster_signatures(distance_matrix)

        # Save clusters
        save_clusters(images, labels, filenames, clustering_results_dir, input_folder_name)

        print(f"Clustering completed. Results saved in {os.path.join(clustering_results_dir, input_folder_name)}")

        # Determine S3 folder name
        s3_folder = input_folder_name

        # Load existing S3 URIs
        s3_uris = load_s3_uris(s3_uris_file)

        # Upload clustering results to S3 and collect URIs
        new_uris = upload_directory_to_s3(os.path.join(clustering_results_dir, input_folder_name), bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name)
        s3_uris.extend(new_uris)

        # Save updated S3 URIs
        save_s3_uris(s3_uris_file, s3_uris)

        print(f"Clustering results uploaded to S3. URIs: {new_uris}")
        print(f"All S3 URIs: {s3_uris}")

        # Clean up temporary directory
        try:
            shutil.rmtree(temp_dir)
            print(f"Temporary directory {temp_dir} removed successfully.")
        except Exception as e:
            print(f"Error removing temporary directory {temp_dir}: {e}")

        return s3_uris

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class Query(BaseModel):
    docPath: str

@app.get('/sign')
async def root(query: Query):
    s3_uris = process_pdf_from_s3(query.docPath)
    return { 'message': "Success!", 'data': s3_uris }
