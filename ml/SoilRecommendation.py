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
gravel_path = 'D:/soil classification/soilPhotos/Gravel'
sand_path = 'D:/soil classification/soilPhotos/Sand'
silt_path = 'D:/soil classification/soilPhotos/Silt'

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

# Specify the path of the test image
test_image_path = 'D:/soil classification/test/Gravel/4.jpg'

# Make a prediction for the test image
predicted_class = predict_single_image(test_image_path, rf_classifier, label_encoder)

# Print the predicted class
print(f"The test image belongs to the class: {predicted_class}")

if predicted_class == "Gravel":
    print("\nGravelly soil tends to be well-draining but may not retain water and nutrients as effectively. Crops that can tolerate well-draining conditions and do not require high fertility may be suitable.\n\nRecommended crops: Grapevines, certain varieties of strawberries, drought-tolerant vegetables (e.g., tomatoes, peppers), and herbs like rosemary.")
    
if predicted_class == "Sand":
    print("\nSandy soil drains quickly but may lack fertility and moisture retention. Crops that thrive in well-draining conditions and can tolerate lower fertility are suitable for sandy soil.\n\nRecommended crops: Carrots, radishes, potatoes, asparagus, melons, and drought-tolerant plants like lavender.")
if predicted_class == "Silt":
    print("\nSilt has finer particles and better moisture retention than sand. It is fertile and suitable for a wide range of crops.\n\nRecommended crops: Beans, peas, leafy greens (e.g., lettuce, spinach), broccoli, cabbage, and most garden vegetables. Silt soil is versatile and supports various crops.")