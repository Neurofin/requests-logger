import os
import shutil
import boto3
import logging
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from detect import pdf_to_png, load_and_extract_features, compute_similarity_matrix, similarity_to_distance_matrix, cluster_signatures, save_clusters
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

app = FastAPI()

class S3PdfUri(BaseModel):
    s3_uri: str

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
    try:
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
        logging.info(f"Downloaded {s3_uri} to {local_path}")
        return local_path
    except Exception as e:
        logging.error(f"Failed to download {s3_uri}: {e}")
        raise

def upload_directory_to_s3(directory_path, bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name):
    try:
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
                logging.info(f"Uploaded {local_path} to {s3_uri}")
        return uris
    except Exception as e:
        logging.error(f"Failed to upload directory {directory_path} to S3: {e}")
        raise

def process_pdf_from_s3(s3_pdf_uri):
    temp_dir = 'temp'
    final_results_dir = os.path.join(temp_dir, 'final_results')
    clustering_results_dir = os.path.join(temp_dir, 'clustering_results')
    bucket_name = 'clustering-results123'
    s3_uris_file = os.path.join(temp_dir, 's3_uris.txt')

    aws_access_key_id = os.getenv('AWS_ACCESS_KEY_ID')
    aws_secret_access_key = os.getenv('AWS_SECRET_ACCESS_KEY')
    region_name = os.getenv('AWS_REGION')

    if not all([aws_access_key_id, aws_secret_access_key, region_name]):
        logging.error("AWS credentials or region are not set in environment variables.")
        return []

    os.makedirs(temp_dir, exist_ok=True)

    try:
        pdf_path = download_pdf_from_s3(s3_pdf_uri, temp_dir, aws_access_key_id, aws_secret_access_key, region_name)
        pdf_to_png(pdf_path, temp_dir, final_results_dir)

        pdf_base_name = os.path.splitext(os.path.basename(pdf_path))[0]
        input_folder = os.path.join(final_results_dir, pdf_base_name)

        if not os.path.exists(clustering_results_dir):
            os.makedirs(clustering_results_dir)

        images, filenames, features = load_and_extract_features(input_folder)
        if len(features) == 0:
            logging.warning("No features extracted. Exiting.")
            return []

        similarity_matrix = compute_similarity_matrix(features)
        distance_matrix = similarity_to_distance_matrix(similarity_matrix)
        labels = cluster_signatures(distance_matrix)
        save_clusters(images, labels, filenames, clustering_results_dir, pdf_base_name)

        logging.info(f"Clustering completed. Results saved in {os.path.join(clustering_results_dir, pdf_base_name)}")

        s3_folder = pdf_base_name
        s3_uris = load_s3_uris(s3_uris_file)
        new_uris = upload_directory_to_s3(os.path.join(clustering_results_dir, pdf_base_name), bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name)
        s3_uris.extend(new_uris)
        save_s3_uris(s3_uris_file, s3_uris)

        logging.info(f"Clustering results uploaded to S3. URIs: {new_uris}")
        logging.info(f"All S3 URIs: {s3_uris}")

        shutil.rmtree(temp_dir)
        logging.info(f"Temporary directory {temp_dir} removed successfully.")

        return s3_uris
    except Exception as e:
        logging.error(f"Error processing PDF from S3: {e}")
        if os.path.exists(temp_dir):
            try:
                shutil.rmtree(temp_dir)
                logging.info(f"Temporary directory {temp_dir} removed successfully.")
            except Exception as cleanup_error:
                logging.error(f"Error removing temporary directory {temp_dir}: {cleanup_error}")
        raise

@app.post("/get-signatures")
async def process_pdf_endpoint(s3_pdf_uri: S3PdfUri):
    try:
        s3_uris = process_pdf_from_s3(s3_pdf_uri.s3_uri)
        if not s3_uris:
            raise HTTPException(status_code=500, detail="Failed to process PDF or no features extracted.")
        return {"s3_uris": s3_uris}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# if __name__ == "__main__":
#     import uvicorn
#     uvicorn.run(app, host="0.0.0.0", port=8000)







# import os
# import shutil
# import boto3
# from detect import pdf_to_png, load_and_extract_features, compute_similarity_matrix, similarity_to_distance_matrix, cluster_signatures, save_clusters
# from dotenv import load_dotenv
# load_dotenv()
# def load_s3_uris(file_path):
#     if os.path.exists(file_path):
#         with open(file_path, 'r') as file:
#             return [line.strip() for line in file.readlines()]
#     return []

# def save_s3_uris(file_path, uris):
#     with open(file_path, 'w') as file:
#         for uri in uris:
#             file.write(f"{uri}\n")

# def download_pdf_from_s3(s3_uri, download_path, aws_access_key_id, aws_secret_access_key, region_name):
#     s3_client = boto3.client(
#         's3',
#         aws_access_key_id=aws_access_key_id,
#         aws_secret_access_key=aws_secret_access_key,
#         region_name=region_name
#     )
#     bucket_name, key = s3_uri.replace("s3://", "").split("/", 1)
#     local_path = os.path.join(download_path, os.path.basename(key))
#     os.makedirs(os.path.dirname(local_path), exist_ok=True)
#     s3_client.download_file(bucket_name, key, local_path)
#     print(f"Downloaded {s3_uri} to {local_path}")
#     return local_path

# def upload_directory_to_s3(directory_path, bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name):
#     s3_client = boto3.client(
#         's3',
#         aws_access_key_id=aws_access_key_id,
#         aws_secret_access_key=aws_secret_access_key,
#         region_name=region_name
#     )
#     uris = []
#     for root, dirs, files in os.walk(directory_path):
#         for file in files:
#             local_path = os.path.join(root, file)
#             relative_path = os.path.relpath(local_path, directory_path)
#             s3_path = os.path.join(s3_folder, relative_path)
#             s3_client.upload_file(local_path, bucket_name, s3_path)
#             s3_uri = f"s3://{bucket_name}/{s3_path}"
#             uris.append(s3_uri)
#             print(f"Uploaded {local_path} to {s3_uri}")
#     return uris

# def process_pdf_from_s3(s3_pdf_uri):
#     # Define paths
#     temp_dir = 'temp'
#     final_results_dir = 'temp/final_results'
#     clustering_results_dir = 'temp/clustering_results'
#     bucket_name = 'clustering-results123'
#     s3_uris_file = 'temp/s3_uris.txt'

#     # AWS credentials (retrieved from environment variables)
#     aws_access_key_id = os.getenv('AWS_ACCESS_KEY_ID')
#     aws_secret_access_key = os.getenv('AWS_SECRET_ACCESS_KEY')
#     region_name = os.getenv('AWS_REGION')

#     # Ensure temporary directory exists
#     os.makedirs(temp_dir, exist_ok=True)

#     # Download PDF from S3
#     pdf_path = download_pdf_from_s3(s3_pdf_uri, temp_dir, aws_access_key_id, aws_secret_access_key, region_name)

#     # Perform PDF to PNG conversion and signature extraction
#     pdf_to_png(pdf_path, temp_dir, final_results_dir)

#     # Extract the base name of the PDF without extension to match the output folder name
#     pdf_base_name = os.path.splitext(os.path.basename(pdf_path))[0]

#     # Define the input folder for clustering
#     input_folder = os.path.join(final_results_dir, pdf_base_name)

#     if not os.path.exists(clustering_results_dir):
#         os.makedirs(clustering_results_dir)

#     # Extract input folder name
#     input_folder_name = os.path.basename(os.path.normpath(input_folder))

#     # Load images and extract features
#     images, filenames, features = load_and_extract_features(input_folder)

#     # Check if features are extracted
#     if len(features) == 0:
#         print("No features extracted. Exiting.")
#         return []
#     else:
#         # Compute similarity matrix
#         similarity_matrix = compute_similarity_matrix(features)

#         # Convert similarity matrix to distance matrix
#         distance_matrix = similarity_to_distance_matrix(similarity_matrix)

#         # Cluster signatures
#         labels = cluster_signatures(distance_matrix)

#         # Save clusters
#         save_clusters(images, labels, filenames, clustering_results_dir, input_folder_name)

#         print(f"Clustering completed. Results saved in {os.path.join(clustering_results_dir, input_folder_name)}")

#         # Determine S3 folder name
#         s3_folder = input_folder_name

#         # Load existing S3 URIs
#         s3_uris = load_s3_uris(s3_uris_file)

#         # Upload clustering results to S3 and collect URIs
#         new_uris = upload_directory_to_s3(os.path.join(clustering_results_dir, input_folder_name), bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name)
#         s3_uris.extend(new_uris)

#         # Save updated S3 URIs
#         save_s3_uris(s3_uris_file, s3_uris)

#         print(f"Clustering results uploaded to S3. URIs: {new_uris}")
#         print(f"All S3 URIs: {s3_uris}")

#         # Clean up temporary directory
#         try:
#             shutil.rmtree(temp_dir)
#             print(f"Temporary directory {temp_dir} removed successfully.")
#         except Exception as e:
#             print(f"Error removing temporary directory {temp_dir}: {e}")

#         return s3_uris

# if __name__ == "__main__":
#     # Get S3 PDF URI from user input
#     s3_pdf_uri = 's3://testing-neurofin/GPA.PDF'
    
#     # Process the PDF from S3 and get the list of S3 URIs
#     s3_uris = process_pdf_from_s3(s3_pdf_uri)
#     print(f"List of S3 URIs: {s3_uris}")
# import os
# import shutil
# import boto3
# import logging
# from detect import pdf_to_png, load_and_extract_features, compute_similarity_matrix, similarity_to_distance_matrix, cluster_signatures, save_clusters
# from dotenv import load_dotenv

# # Load environment variables
# load_dotenv()

# # Configure logging
# logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# def load_s3_uris(file_path):
#     if os.path.exists(file_path):
#         with open(file_path, 'r') as file:
#             return [line.strip() for line in file.readlines()]
#     return []

# def save_s3_uris(file_path, uris):
#     with open(file_path, 'w') as file:
#         for uri in uris:
#             file.write(f"{uri}\n")

# def download_pdf_from_s3(s3_uri, download_path, aws_access_key_id, aws_secret_access_key, region_name):
#     try:
#         s3_client = boto3.client(
#             's3',
#             aws_access_key_id=aws_access_key_id,
#             aws_secret_access_key=aws_secret_access_key,
#             region_name=region_name
#         )
#         bucket_name, key = s3_uri.replace("s3://", "").split("/", 1)
#         local_path = os.path.join(download_path, os.path.basename(key))
#         os.makedirs(os.path.dirname(local_path), exist_ok=True)
#         s3_client.download_file(bucket_name, key, local_path)
#         logging.info(f"Downloaded {s3_uri} to {local_path}")
#         return local_path
#     except Exception as e:
#         logging.error(f"Failed to download {s3_uri}: {e}")
#         raise

# def upload_directory_to_s3(directory_path, bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name):
#     try:
#         s3_client = boto3.client(
#             's3',
#             aws_access_key_id=aws_access_key_id,
#             aws_secret_access_key=aws_secret_access_key,
#             region_name=region_name
#         )
#         uris = []
#         for root, dirs, files in os.walk(directory_path):
#             for file in files:
#                 local_path = os.path.join(root, file)
#                 relative_path = os.path.relpath(local_path, directory_path)
#                 s3_path = os.path.join(s3_folder, relative_path)
#                 s3_client.upload_file(local_path, bucket_name, s3_path)
#                 s3_uri = f"s3://{bucket_name}/{s3_path}"
#                 uris.append(s3_uri)
#                 logging.info(f"Uploaded {local_path} to {s3_uri}")
#         return uris
#     except Exception as e:
#         logging.error(f"Failed to upload directory {directory_path} to S3: {e}")
#         raise

# def process_pdf_from_s3(s3_pdf_uri):
#     temp_dir = 'temp'
#     final_results_dir = os.path.join(temp_dir, 'final_results')
#     clustering_results_dir = os.path.join(temp_dir, 'clustering_results')
#     bucket_name = 'clustering-results123'
#     s3_uris_file = os.path.join(temp_dir, 's3_uris.txt')

#     aws_access_key_id = os.getenv('AWS_ACCESS_KEY_ID')
#     aws_secret_access_key = os.getenv('AWS_SECRET_ACCESS_KEY')
#     region_name = os.getenv('AWS_REGION')

#     if not all([aws_access_key_id, aws_secret_access_key, region_name]):
#         logging.error("AWS credentials or region are not set in environment variables.")
#         return []

#     os.makedirs(temp_dir, exist_ok=True)

#     try:
#         pdf_path = download_pdf_from_s3(s3_pdf_uri, temp_dir, aws_access_key_id, aws_secret_access_key, region_name)
#         pdf_to_png(pdf_path, temp_dir, final_results_dir)

#         pdf_base_name = os.path.splitext(os.path.basename(pdf_path))[0]
#         input_folder = os.path.join(final_results_dir, pdf_base_name)

#         if not os.path.exists(clustering_results_dir):
#             os.makedirs(clustering_results_dir)

#         images, filenames, features = load_and_extract_features(input_folder)
#         if len(features) == 0:
#             logging.warning("No features extracted. Exiting.")
#             return []

#         similarity_matrix = compute_similarity_matrix(features)
#         distance_matrix = similarity_to_distance_matrix(similarity_matrix)
#         labels = cluster_signatures(distance_matrix)
#         save_clusters(images, labels, filenames, clustering_results_dir, pdf_base_name)

#         logging.info(f"Clustering completed. Results saved in {os.path.join(clustering_results_dir, pdf_base_name)}")

#         s3_folder = pdf_base_name
#         s3_uris = load_s3_uris(s3_uris_file)
#         new_uris = upload_directory_to_s3(os.path.join(clustering_results_dir, pdf_base_name), bucket_name, s3_folder, aws_access_key_id, aws_secret_access_key, region_name)
#         s3_uris.extend(new_uris)
#         save_s3_uris(s3_uris_file, s3_uris)

#         logging.info(f"Clustering results uploaded to S3. URIs: {new_uris}")
#         logging.info(f"All S3 URIs: {s3_uris}")

#         shutil.rmtree(temp_dir)
#         logging.info(f"Temporary directory {temp_dir} removed successfully.")

#         return s3_uris
#     except Exception as e:
#         logging.error(f"Error processing PDF from S3: {e}")
#         if os.path.exists(temp_dir):
#             try:
#                 shutil.rmtree(temp_dir)
#                 logging.info(f"Temporary directory {temp_dir} removed successfully.")
#             except Exception as cleanup_error:
#                 logging.error(f"Error removing temporary directory {temp_dir}: {cleanup_error}")
#         return []

# if __name__ == "__main__":
#     s3_pdf_uri = 's3://testing-neurofin/GPA.PDF'
#     s3_uris = process_pdf_from_s3(s3_pdf_uri)
#     logging.info(f"List of S3 URIs: {s3_uris}")