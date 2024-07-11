import os
import cv2 # type: ignore
import numpy as np
from keras.applications import VGG16
from keras.applications.vgg16 import preprocess_input
from sklearn.metrics.pairwise import cosine_similarity
from sklearn.cluster import DBSCAN
import shutil
import os
import shutil

# Load the pre-trained VGG16 model
vgg_model = VGG16(weights='imagenet', include_top=False, input_shape=(100, 200, 3))

# Function to preprocess image for VGG16 input
def preprocess_image_vgg(image):
    resized_image = cv2.resize(image, (200, 100))
    if resized_image.ndim == 2:
        resized_image = cv2.cvtColor(resized_image, cv2.COLOR_GRAY2RGB)
    processed_image = preprocess_input(resized_image)
    return processed_image

# Function to extract features using VGG16 model
def extract_features_vgg(image):
    processed_image = preprocess_image_vgg(image)
    processed_image = np.expand_dims(processed_image, axis=0)
    features = vgg_model.predict(processed_image)
    return features.flatten()

# Load images from the folder and extract features
# def load_and_extract_features(folder):
#     images = []
#     filenames = []
#     features = []
#     for filename in os.listdir(folder):
#         img_path = os.path.join(folder, filename)
#         img = cv2.imread(img_path)
#         if img is not None:
#             images.append(img)
#             filenames.append(filename)
#             features.append(extract_features_vgg(img))
#     return images, filenames, np.array(features)
# Load images from the folder and extract features
def load_and_extract_features(folder):
    images = []
    filenames = []
    features = []
    for filename in os.listdir(folder):
        img_path = os.path.join(folder, filename)
        img = cv2.imread(img_path)
        if img is not None:
            images.append(img)
            filenames.append(filename)
            features.append(extract_features_vgg(img))
        else:
            print(f"Failed to load image: {img_path}")
    print(f"Number of images loaded: {len(images)}")
    print(f"Number of features extracted: {len(features)}")
    return images, filenames, np.array(features)




# Compute cosine similarity matrix
def compute_similarity_matrix(features):

    similarity_matrix = cosine_similarity(features)
    return similarity_matrix

# Convert similarity matrix to distance matrix
def similarity_to_distance_matrix(similarity_matrix):
    similarity_matrix = np.clip(similarity_matrix, -1, 1)
    return 1 - similarity_matrix

# Cluster signatures based on the distance matrix
def cluster_signatures(distance_matrix, eps=0.5, min_samples=2):
    clustering = DBSCAN(metric='precomputed', eps=eps, min_samples=min_samples)
    labels = clustering.fit_predict(distance_matrix)
    return labels

# Save clusters to the specified directory
def save_clusters(images, labels, filenames, results_path, input_folder_name):
    unique_labels = set(labels)
    cluster_mapping = {label: idx for idx, label in enumerate(sorted(unique_labels))}
    output_folder = os.path.join(results_path, input_folder_name)

    # Check if the folder already exists, if so, create a new folder with a modified name
    if os.path.exists(output_folder):
        # Find a unique folder name by appending (2), (3), etc.
        i = 2
        while True:
            new_folder_name = f"{input_folder_name} ({i})"
            new_output_folder = os.path.join(results_path, new_folder_name)
            if not os.path.exists(new_output_folder):
                output_folder = new_output_folder
                break
            i += 1

    os.makedirs(output_folder, exist_ok=True)

    for label in unique_labels:
        cluster_idx = cluster_mapping[label]
        cluster_folder = os.path.join(output_folder, f'{cluster_idx}')
        if not os.path.exists(cluster_folder):
            os.makedirs(cluster_folder)

        indices = [i for i, lbl in enumerate(labels) if lbl == label]
        for index in indices:
            filename = filenames[index]
            img_path = os.path.join(cluster_folder, filename)
            cv2.imwrite(img_path, images[index])



import time
from pathlib import Path

import cv2
import torch

from yolo_components.attempt_load import attempt_load
import math
import torch.nn as nn
import torchvision

from yolo_components.attempt_load import Colors, plot_one_box, xyxy2xywh, save_one_box, scale_coords, non_max_suppression, apply_classifier

import re

def increment_path(path, exist_ok=False, sep='', mkdir=False):
    # Increment file or directory path, i.e. runs/exp --> runs/exp{sep}2, runs/exp{sep}3, ... etc.
    path = Path(path)  # os-agnostic
    if path.exists() and not exist_ok:
        suffix = path.suffix
        path = path.with_suffix('')
        dirs = glob.glob(f"{path}{sep}*")  # similar paths
        matches = [re.search(rf"%s{sep}(\d+)" % path.stem, d) for d in dirs]
        i = [int(m.groups()[0]) for m in matches if m]  # indices
        n = max(i) + 1 if i else 2  # increment number
        path = Path(f"{path}{sep}{n}{suffix}")  # update path
    dir = path if path.suffix == '' else path.parent  # directory
    if not dir.exists() and mkdir:
        dir.mkdir(parents=True, exist_ok=True)  # make directory
    return path

def make_divisible(x, divisor):
    # Returns x evenly divisible by divisor
    return math.ceil(x / divisor) * divisor
def check_img_size(img_size, s=32):
    # Verify img_size is a multiple of stride s
    new_size = make_divisible(img_size, int(s))  # ceil gs-multiple
    if new_size != img_size:
        print('WARNING: --img-size %g must be multiple of max stride %g, updating to %g' % (img_size, s, new_size))
    return new_size

def load_classifier(name='resnet101', n=2):
    # Loads a pretrained model reshaped to n-class output
    model = torchvision.models.__dict__[name](pretrained=True)

    # ResNet model properties
    # input_size = [3, 224, 224]
    # input_space = 'RGB'
    # input_range = [0, 1]
    # mean = [0.485, 0.456, 0.406]
    # std = [0.229, 0.224, 0.225]

    # Reshape output to n classes
    filters = model.fc.weight.shape[1]
    model.fc.bias = nn.Parameter(torch.zeros(n), requires_grad=True)
    model.fc.weight = nn.Parameter(torch.zeros(n, filters), requires_grad=True)
    model.fc.out_features = n
    return model

import logging
def set_logging(rank=-1, verbose=True):
    logging.basicConfig(
        format="%(message)s",
        level=logging.INFO if (verbose and rank in [-1, 0]) else logging.WARN)
    
def check_imshow():
    # Check if environment supports image displays
    return False

img_formats = ['bmp', 'jpg', 'jpeg', 'png', 'tif', 'tiff', 'dng', 'webp', 'mpo']  # acceptable image suffixes
vid_formats = ['mov', 'avi', 'mp4', 'mpg', 'mpeg', 'm4v', 'wmv', 'mkv']  # acceptable video suffixes
def letterbox(img, new_shape=(640, 640), color=(114, 114, 114), auto=True, scaleFill=False, scaleup=True, stride=32):
    # Resize and pad image while meeting stride-multiple constraints
    shape = img.shape[:2]  # current shape [height, width]
    if isinstance(new_shape, int):
        new_shape = (new_shape, new_shape)

    # Scale ratio (new / old)
    r = min(new_shape[0] / shape[0], new_shape[1] / shape[1])
    if not scaleup:  # only scale down, do not scale up (for better test mAP)
        r = min(r, 1.0)

    # Compute padding
    ratio = r, r  # width, height ratios
    new_unpad = int(round(shape[1] * r)), int(round(shape[0] * r))
    dw, dh = new_shape[1] - new_unpad[0], new_shape[0] - new_unpad[1]  # wh padding
    if auto:  # minimum rectangle
        dw, dh = np.mod(dw, stride), np.mod(dh, stride)  # wh padding
    elif scaleFill:  # stretch
        dw, dh = 0.0, 0.0
        new_unpad = (new_shape[1], new_shape[0])
        ratio = new_shape[1] / shape[1], new_shape[0] / shape[0]  # width, height ratios

    dw /= 2  # divide padding into 2 sides
    dh /= 2

    if shape[::-1] != new_unpad:  # resize
        img = cv2.resize(img, new_unpad, interpolation=cv2.INTER_LINEAR)
    top, bottom = int(round(dh - 0.1)), int(round(dh + 0.1))
    left, right = int(round(dw - 0.1)), int(round(dw + 0.1))
    img = cv2.copyMakeBorder(img, top, bottom, left, right, cv2.BORDER_CONSTANT, value=color)  # add border
    return img, ratio, (dw, dh)
class LoadImages:  # for inference
    def __init__(self, path, img_size=640, stride=32):
        p = str(Path(path).absolute())  # os-agnostic absolute path
        if '*' in p:
            files = sorted(glob.glob(p, recursive=True))  # glob
        elif os.path.isdir(p):
            files = sorted(glob.glob(os.path.join(p, '*.*')))  # dir
        elif os.path.isfile(p):
            files = [p]  # files
        else:
            raise Exception(f'ERROR: {p} does not exist')

        images = [x for x in files if x.split('.')[-1].lower() in img_formats]
        videos = [x for x in files if x.split('.')[-1].lower() in vid_formats]
        ni, nv = len(images), len(videos)

        self.img_size = img_size
        self.stride = stride
        self.files = images + videos
        self.nf = ni + nv  # number of files
        self.video_flag = [False] * ni + [True] * nv
        self.mode = 'image'
        if any(videos):
            self.new_video(videos[0])  # new video
        else:
            self.cap = None
        assert self.nf > 0, f'No images or videos found in {p}. ' \
                            f'Supported formats are:\nimages: {img_formats}\nvideos: {vid_formats}'

    def __iter__(self):
        self.count = 0
        return self

    def __next__(self):
        if self.count == self.nf:
            raise StopIteration
        path = self.files[self.count]

        if self.video_flag[self.count]:
            # Read video
            self.mode = 'video'
            ret_val, img0 = self.cap.read()
            if not ret_val:
                self.count += 1
                self.cap.release()
                if self.count == self.nf:  # last video
                    raise StopIteration
                else:
                    path = self.files[self.count]
                    self.new_video(path)
                    ret_val, img0 = self.cap.read()

            self.frame += 1
            print(f'video {self.count + 1}/{self.nf} ({self.frame}/{self.nframes}) {path}: ', end='')

        else:
            # Read image
            self.count += 1
            img0 = cv2.imread(path)  # BGR
            assert img0 is not None, 'Image Not Found ' + path
            print(f'image {self.count}/{self.nf} {path}: ', end='')

        # Padded resize
        img = letterbox(img0, self.img_size, stride=self.stride)[0]

        # Convert
        img = img[:, :, ::-1].transpose(2, 0, 1)  # BGR to RGB, to 3x416x416
        img = np.ascontiguousarray(img)

        return path, img, img0, self.cap

    def new_video(self, path):
        self.frame = 0
        self.cap = cv2.VideoCapture(path)
        self.nframes = int(self.cap.get(cv2.CAP_PROP_FRAME_COUNT))

    def __len__(self):
        return self.nf  # number of files

# from SOURCE.yolo_files.models.experimental import attempt_load

def detect(image_path):
    opt = {
    'weights': 'best.pt',
    'source': image_path,
    'img_size': 640,
    'conf_thres': 0.25,
    'iou_thres': 0.45,
    'device': '',
    'view_img': False,
    'save_txt': True,
    'save_conf': True,
    'save_crop': True,
    'nosave': True,
    'classes': 1,
    'agnostic_nms': False,
    'augment': False,
    'update': False,
    'project': 'temp/results/yolov5/',
    'name': 'exp',
    'exist_ok': False,
    'line_thickness': 3,
    'hide_labels': False,
    'hide_conf': False,

}


    source, weights, view_img, save_txt, imgsz = opt['source'], opt['weights'], opt['view_img'], opt['save_txt'], opt['img_size']
    save_img = not opt['nosave'] and not source.endswith('.txt')  # save inference images
    # webcam = source.isnumeric() or source.endswith('.txt') or source.lower().startswith(
    #     ('rtsp://', 'rtmp://', 'http://', 'https://'))

    # Directories
    save_dir = increment_path(Path(opt['project']) / opt['name'], exist_ok=opt['exist_ok'])  # increment run
    (save_dir / 'labels' if save_txt else save_dir).mkdir(parents=True, exist_ok=True)  # make dir

    # Initialize
    set_logging()
    # device = select_device(opt['device'])
    device = torch.device('cpu')
    half = device.type != 'cpu'  # half precision only supported on CUDA

    # Load model
    model = attempt_load(weights, map_location=device)  # load FP32 model
    stride = int(model.stride.max())  # model stride
    imgsz = check_img_size(imgsz, s=stride)  # check img_size
    names = model.module.names if hasattr(model, 'module') else model.names  # get class names
    if half:
        model.half()  # to FP16

    # Second-stage classifier
    classify = False
    if classify:
        modelc = load_classifier(name='resnet101', n=2)  # initialize
        modelc.load_state_dict(torch.load('weights/resnet101.pt', map_location=device)['model']).to(device).eval()

    # Set Dataloader
    vid_path, vid_writer = None, None
    # if webcam:
    #     view_img = check_imshow()
    #     cudnn.benchmark = True  # set True to speed up constant image size inference
    #     dataset = LoadStreams(source, img_size=imgsz, stride=stride)
    # else:
    dataset = LoadImages(source, img_size=imgsz, stride=stride)

    # Run inference
    if device.type != 'cpu':
        model(torch.zeros(1, 3, imgsz, imgsz).to(device).type_as(next(model.parameters())))  # run once
    t0 = time.time()
    for path, img, im0s, vid_cap in dataset:
        img = torch.from_numpy(img).to(device)
        img = img.half() if half else img.float()  # uint8 to fp16/32
        img /= 255.0  # 0 - 255 to 0.0 - 1.0
        if img.ndimension() == 3:
            img = img.unsqueeze(0)

        # Inference
        t1 = time.time()
        pred = model(img, augment=opt['augment'])[0]

        # Apply NMS
        pred = non_max_suppression(pred, opt['conf_thres'], opt['iou_thres'], classes=opt['classes'], agnostic=opt['agnostic_nms'])
        t2 = time.time()

        # Apply Classifier
        if classify:
            pred = apply_classifier(pred, modelc, img, im0s)

        # Process detections
        for i, det in enumerate(pred):  # detections per image
            # if webcam:  # batch_size >= 1
            #     p, s, im0, frame = path[i], '%g: ' % i, im0s[i].copy(), dataset.count
            # else:
            p, s, im0, frame = path, '', im0s.copy(), getattr(dataset, 'frame', 0)

            p = Path(p)  # to Path
            save_path = str(save_dir / p.name)  # img.jpg
            txt_path = str(save_dir / 'labels' / p.stem) + ('' if dataset.mode == 'image' else f'_{frame}')  # img.txt
            s += '%gx%g ' % img.shape[2:]  # print string
            gn = torch.tensor(im0.shape)[[1, 0, 1, 0]]  # normalization gain whwh
            if len(det):
                # Rescale boxes from img_size to im0 size
                det[:, :4] = scale_coords(img.shape[2:], det[:, :4], im0.shape).round()

                # Print results
                for c in det[:, -1].unique():
                    n = (det[:, -1] == c).sum()  # detections per class
                    s += f"{n} {names[int(c)]}{'s' * (n > 1)}, "  # add to string

                # Write results
                for *xyxy, conf, cls in reversed(det):
                    if save_txt:  # Write to file
                        xywh = (xyxy2xywh(torch.tensor(xyxy).view(1, 4)) / gn).view(-1).tolist()  # normalized xywh
                        line = (cls, *xywh, conf) if opt['save_conf'] else (cls, *xywh)  # label format
                        with open(txt_path + '.txt', 'a') as f:
                            f.write(('%g ' * len(line)).rstrip() % line + '\n')

                    if save_img or opt['save_crop'] or view_img:  # Add bbox to image
                        c = int(cls)  # integer class
                        label = None if opt['hide_labels'] else (names[c] if opt['hide_conf'] else f'{names[c]} {conf:.2f}')

                        plot_one_box(xyxy, im0, label=label, color=Colors()(c, True), line_thickness=opt['line_thickness'])
                        if opt['save_crop']:
                            save_one_box(xyxy, im0s, file=save_dir / 'crops' / names[c] / f'{p.stem}.jpg', BGR=True)

            # Print time (inference + NMS)
            print(f'{s}Done. ({t2 - t1:.3f}s)')

            # Stream results
            if view_img:
                cv2.imshow(str(p), im0)
                cv2.waitKey(1)  # 1 millisecond

            # Save results (image with detections)
            if save_img:
                if dataset.mode == 'image':
                    cv2.imwrite(save_path, im0)
                else:  # 'video' or 'stream'
                    if vid_path != save_path:  # new video
                        vid_path = save_path
                        if isinstance(vid_writer, cv2.VideoWriter):
                            vid_writer.release()  # release previous video writer
                        if vid_cap:  # video
                            fps = vid_cap.get(cv2.CAP_PROP_FPS)
                            w = int(vid_cap.get(cv2.CAP_PROP_FRAME_WIDTH))
                            h = int(vid_cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
                        else:  # stream
                            fps, w, h = 30, im0.shape[1], im0.shape[0]
                            save_path += '.mp4'
                        vid_writer = cv2.VideoWriter(save_path, cv2.VideoWriter_fourcc(*'mp4v'), fps, (w, h))
                    vid_writer.write(im0)

    if save_txt or save_img:
        s = f"\n{len(list(save_dir.glob('labels/*.txt')))} labels saved to {save_dir / 'labels'}" if save_txt else ''
        print(f"Results saved to {save_dir}{s}")

    print(f'Done. ({time.time() - t0:.3f}s)')
    return 'Success'





import cv2
import shutil
import glob
import os
from pathlib import Path
# from SOURCE.yolo_files import detect
import fitz  # PyMuPDF

YOLO_RESULT = 'temp/results/yolov5/'
YOLO_OP = 'crops/DLSignature/'

# def clear_yolo_results():
#     # Clear the YOLO results directory
#     for folder in glob.glob(os.path.join(YOLO_RESULT, '*/')):
#         shutil.rmtree(folder)
#     print(f"Cleared YOLO results directory: {YOLO_RESULT}")

def pdf_to_png(pdf_path, temp_dir, final_results_dir):
    # Open the PDF
    pdf_document = fitz.open(pdf_path)

    # Get the filename without extension
    filename = Path(pdf_path).stem

    # Create a subfolder for the current PDF within the final results directory
    output_signature_folder = os.path.join(final_results_dir, filename)

    # Check if the folder already exists, if so, create a new folder with a modified name
    if os.path.exists(output_signature_folder):
        # Find a unique folder name by appending (2), (3), etc.
        i = 2
        while True:
            new_folder_name = f"{filename} ({i})"
            new_output_signature_folder = os.path.join(final_results_dir, new_folder_name)
            if not os.path.exists(new_output_signature_folder):
                output_signature_folder = new_output_signature_folder
                break
            i += 1

    os.makedirs(output_signature_folder, exist_ok=True)
    print(f"Output signature folder created: {output_signature_folder}")

    # Clear YOLO results directory before starting detection
    # clear_yolo_results()

    # Iterate through each page
    for page_number in range(len(pdf_document)):
        # Get the page
        page = pdf_document.load_page(page_number)

        # Render the page as a PNG image
        pix = page.get_pixmap()

        # Construct the output file name
        output_filename = os.path.join(temp_dir, f"{filename}_page_{page_number + 1}.png")

        # Save the PNG image
        pix.save(output_filename)
        print(f"Saved PNG image: {output_filename}")

        # Perform signature detection on the PNG image
        signature_detection(output_filename, output_signature_folder)

    # Close the PDF
    pdf_document.close()

def signature_detection(image_path, output_folder):
    '''
    Performs signature detection on a PNG image and saves the detected signatures in the output folder.

    Args:
    - image_path: Path to the input PNG image.
    - output_folder: Path to the folder where detected signatures will be saved.
    '''
    # Ensure the output folder exists, create if it doesn't.
    if not os.path.exists(output_folder):
        os.makedirs(output_folder)

    # Call YOLOv5 detection on the input image.
    detect(image_path)

    # Get the path of the latest detection results.
    latest_detection = max(glob.glob(os.path.join(YOLO_RESULT, '*/')), key=os.path.getmtime)
    print(f"Latest YOLO detection results folder: {latest_detection}")

    # Copy detected signatures to the output folder.
    yolo_op_folder = os.path.join(latest_detection, YOLO_OP)
    signature_files = glob.glob(os.path.join(yolo_op_folder, '*.jpg'))
    print(f"Detected signature files: {signature_files}")

    if not signature_files:
        print(f"No signature files detected for image: {image_path}")

    for idx, signature_file in enumerate(signature_files):
        # Construct the destination path for the signature file.
        input_filename = os.path.splitext(os.path.basename(image_path))[0]
        destination_filename = f"{input_filename}_signature_{idx+1}.png"
        destination_path = os.path.join(output_folder, destination_filename)
        # Copy the signature file to the output folder.
        shutil.copyfile(signature_file, destination_path)
        print(f"Copied signature file to: {destination_path}")




