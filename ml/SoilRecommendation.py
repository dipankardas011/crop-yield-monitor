from skimage import io, color, transform
import os
import numpy as np
from PIL import Image
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestClassifier
from sklearn.preprocessing import LabelEncoder
from sklearn.metrics import accuracy_score


def load_and_preprocess(folder_path):
    images = []
    labels = []
    if os.path.isdir(folder_path):
        for file_or_dir in os.listdir(folder_path):
            file_or_dir_path = os.path.join(folder_path, file_or_dir)
            if os.path.isdir(file_or_dir_path):
                for file in os.listdir(file_or_dir_path):
                    image_path = os.path.join(file_or_dir_path, file)
                    print("-- Loading --")  # Debugging line
                    image = Image.open(image_path)
                    image = image.convert('L')  # Convert to grayscale
                    image = image.resize((64, 64))  # Resize for consistency
                    images.append(np.array(image).flatten())  # Flatten the image
                    labels.append(file_or_dir)
            elif os.path.isfile(file_or_dir_path):
                print("-- Loading --")
                #print("Loading single image:", file_or_dir_path)  # Debugging line
                image = Image.open(file_or_dir_path)
                image = image.convert('L')  # Convert to grayscale
                image = image.resize((64, 64))  # Resize for consistency
                images.append(np.array(image).flatten())  # Flatten the image
                labels.append(os.path.basename(os.path.dirname(file_or_dir_path)))
    else:
        raise ValueError(f"Invalid path: {folder_path}")
    return np.array(images), np.array(labels)

def predict_single_image(image_path, classifier, label_encoder):
    # Load and preprocess the test image
    test_image = Image.open(image_path)
    test_image = test_image.convert('L')  # Convert to grayscale
    test_image = test_image.resize((64, 64))  # Resize for consistency
    test_image = np.array(test_image).flatten().reshape(1, -1)

    # Predict the class label
    predicted_label_encoded = classifier.predict(test_image)[0]

    # Convert the predicted label back to the original class
    predicted_label = label_encoder.inverse_transform([predicted_label_encoded])[0]

    return predicted_label

# Define paths for each class folder
gravel_path = '/app/train/Gravel'
sand_path = '/app/train/Sand'
silt_path = '/app/train/Silt'

# Load and preprocess data for each class
gravel_images, gravel_labels = load_and_preprocess(gravel_path)
sand_images, sand_labels = load_and_preprocess(sand_path)
silt_images, silt_labels = load_and_preprocess(silt_path)

# Combine data from different classes
all_images = np.concatenate([gravel_images, sand_images, silt_images], axis=0)
all_labels = np.concatenate([gravel_labels, sand_labels, silt_labels], axis=0)

# Convert labels to numerical values using LabelEncoder
label_encoder = LabelEncoder()
all_labels_encoded = label_encoder.fit_transform(all_labels)

# Train a Random Forest classifier
rf_classifier = RandomForestClassifier(n_estimators=100, random_state=42)
rf_classifier.fit(all_images, all_labels_encoded)
